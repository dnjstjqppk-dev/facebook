package repository

import (
	"context"
	"errors"
	"sort"
	"sync"
	"time"

	"github.com/dnjstjqppk-dev/facebook/pkg/domain"
)

var (
	ErrPostNotFound    = errors.New("post not found")
	ErrCommentNotFound = errors.New("comment not found")
)

type TimelineEntry struct {
	PostID    string
	Timestamp time.Time
}

type FeedRepository interface {
	CreatePost(ctx context.Context, post *domain.Post) error
	GetPost(ctx context.Context, id string) (*domain.Post, error)
	UpdatePost(ctx context.Context, post *domain.Post) error
	DeletePost(ctx context.Context, id string) error
	ListPostsByAuthor(ctx context.Context, authorID string) ([]*domain.Post, error)

	SaveReaction(ctx context.Context, r *domain.Reaction) error
	DeleteReaction(ctx context.Context, postID, userID string) error
	GetUserReaction(ctx context.Context, postID, userID string) (domain.ReactionType, bool)

	CreateComment(ctx context.Context, comment *domain.Comment) error
	GetCommentsByPost(ctx context.Context, postID string) ([]*domain.Comment, error)
	GetComment(ctx context.Context, commentID string) (*domain.Comment, error)

	AddToTimeline(ctx context.Context, userID, postID string, t time.Time) error
	GetTimeline(ctx context.Context, userID string) ([]TimelineEntry, error)

	SaveToCollection(ctx context.Context, userID, collectionName, postID string) error
	GetSavedCollections(ctx context.Context, userID string) ([]*domain.SavedCollection, error)
	GetMemories(ctx context.Context, userID string, month time.Month, day, currentYear int) ([]*domain.Post, error)
}

type InMemoryFeedRepository struct {
	mu          sync.RWMutex
	posts       map[string]*domain.Post
	reactions   map[string]map[string]*domain.Reaction // post_id -> user_id -> Reaction
	comments    map[string]*domain.Comment             // comment_id -> Comment
	postComments map[string][]string                   // post_id -> []comment_id
	timelines   map[string][]TimelineEntry             // user_id -> []TimelineEntry
	saved       map[string]map[string]*domain.SavedCollection // user_id -> collection_name -> SavedCollection
}

func NewInMemoryFeedRepository() *InMemoryFeedRepository {
	return &InMemoryFeedRepository{
		posts:        make(map[string]*domain.Post),
		reactions:    make(map[string]map[string]*domain.Reaction),
		comments:     make(map[string]*domain.Comment),
		postComments: make(map[string][]string),
		timelines:    make(map[string][]TimelineEntry),
		saved:        make(map[string]map[string]*domain.SavedCollection),
	}
}

func (r *InMemoryFeedRepository) CreatePost(ctx context.Context, post *domain.Post) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.posts[post.ID] = post
	return nil
}

func (r *InMemoryFeedRepository) GetPost(ctx context.Context, id string) (*domain.Post, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	p, exists := r.posts[id]
	if !exists || p.DeletedAt != nil {
		return nil, ErrPostNotFound
	}
	return p, nil
}

func (r *InMemoryFeedRepository) UpdatePost(ctx context.Context, post *domain.Post) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.posts[post.ID]; !exists {
		return ErrPostNotFound
	}
	r.posts[post.ID] = post
	return nil
}

func (r *InMemoryFeedRepository) DeletePost(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	p, exists := r.posts[id]
	if !exists || p.DeletedAt != nil {
		return ErrPostNotFound
	}
	now := time.Now().UTC()
	p.DeletedAt = &now
	return nil
}

func (r *InMemoryFeedRepository) ListPostsByAuthor(ctx context.Context, authorID string) ([]*domain.Post, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var results []*domain.Post
	for _, p := range r.posts {
		if p.AuthorID == authorID && p.DeletedAt == nil {
			results = append(results, p)
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].CreatedAt.After(results[j].CreatedAt)
	})
	return results, nil
}

func (r *InMemoryFeedRepository) SaveReaction(ctx context.Context, reaction *domain.Reaction) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	p, exists := r.posts[reaction.PostID]
	if !exists || p.DeletedAt != nil {
		return ErrPostNotFound
	}

	if _, ok := r.reactions[reaction.PostID]; !ok {
		r.reactions[reaction.PostID] = make(map[string]*domain.Reaction)
	}

	// Decrement previous reaction type if changed
	if prev, ok := r.reactions[reaction.PostID][reaction.UserID]; ok {
		if prev.ReactionType == reaction.ReactionType {
			return nil // Same reaction, idempotent
		}
		p.ReactionCounts[prev.ReactionType]--
	}

	r.reactions[reaction.PostID][reaction.UserID] = reaction
	if p.ReactionCounts == nil {
		p.ReactionCounts = make(map[domain.ReactionType]int)
	}
	p.ReactionCounts[reaction.ReactionType]++
	return nil
}

func (r *InMemoryFeedRepository) DeleteReaction(ctx context.Context, postID, userID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	p, exists := r.posts[postID]
	if !exists || p.DeletedAt != nil {
		return ErrPostNotFound
	}

	userMap, ok := r.reactions[postID]
	if !ok {
		return nil
	}

	if prev, ok := userMap[userID]; ok {
		p.ReactionCounts[prev.ReactionType]--
		if p.ReactionCounts[prev.ReactionType] < 0 {
			p.ReactionCounts[prev.ReactionType] = 0
		}
		delete(userMap, userID)
	}
	return nil
}

func (r *InMemoryFeedRepository) GetUserReaction(ctx context.Context, postID, userID string) (domain.ReactionType, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if userMap, ok := r.reactions[postID]; ok {
		if react, ok := userMap[userID]; ok {
			return react.ReactionType, true
		}
	}
	return "", false
}

func (r *InMemoryFeedRepository) CreateComment(ctx context.Context, comment *domain.Comment) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	p, exists := r.posts[comment.PostID]
	if !exists || p.DeletedAt != nil {
		return ErrPostNotFound
	}

	r.comments[comment.ID] = comment
	r.postComments[comment.PostID] = append(r.postComments[comment.PostID], comment.ID)
	p.CommentCount++

	if comment.ParentCommentID != "" {
		if parent, exists := r.comments[comment.ParentCommentID]; exists {
			parent.RepliesCount++
		}
	}
	return nil
}

func (r *InMemoryFeedRepository) GetCommentsByPost(ctx context.Context, postID string) ([]*domain.Comment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	ids, ok := r.postComments[postID]
	if !ok {
		return []*domain.Comment{}, nil
	}

	var results []*domain.Comment
	for _, id := range ids {
		if c, exists := r.comments[id]; exists && c.DeletedAt == nil {
			results = append(results, c)
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].CreatedAt.Before(results[j].CreatedAt)
	})
	return results, nil
}

func (r *InMemoryFeedRepository) GetComment(ctx context.Context, commentID string) (*domain.Comment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	c, exists := r.comments[commentID]
	if !exists || c.DeletedAt != nil {
		return nil, ErrCommentNotFound
	}
	return c, nil
}

func (r *InMemoryFeedRepository) AddToTimeline(ctx context.Context, userID, postID string, t time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.timelines[userID] = append(r.timelines[userID], TimelineEntry{
		PostID:    postID,
		Timestamp: t,
	})
	return nil
}

func (r *InMemoryFeedRepository) GetTimeline(ctx context.Context, userID string) ([]TimelineEntry, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	entries := r.timelines[userID]
	result := make([]TimelineEntry, len(entries))
	copy(result, entries)

	sort.Slice(result, func(i, j int) bool {
		return result[i].Timestamp.After(result[j].Timestamp)
	})
	return result, nil
}

func (r *InMemoryFeedRepository) SaveToCollection(ctx context.Context, userID, collectionName, postID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.saved[userID]; !exists {
		r.saved[userID] = make(map[string]*domain.SavedCollection)
	}

	col, exists := r.saved[userID][collectionName]
	if !exists {
		col = &domain.SavedCollection{
			ID:        userID + "-" + collectionName,
			UserID:    userID,
			Name:      collectionName,
			PostIDs:   []string{},
			CreatedAt: time.Now().UTC(),
		}
		r.saved[userID][collectionName] = col
	}

	for _, pid := range col.PostIDs {
		if pid == postID {
			return nil // already saved
		}
	}
	col.PostIDs = append(col.PostIDs, postID)
	return nil
}

func (r *InMemoryFeedRepository) GetSavedCollections(ctx context.Context, userID string) ([]*domain.SavedCollection, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var results []*domain.SavedCollection
	if userCols, exists := r.saved[userID]; exists {
		for _, col := range userCols {
			results = append(results, col)
		}
	}
	return results, nil
}

func (r *InMemoryFeedRepository) GetMemories(ctx context.Context, userID string, month time.Month, day, currentYear int) ([]*domain.Post, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var results []*domain.Post
	for _, p := range r.posts {
		if p.AuthorID == userID && p.DeletedAt == nil {
			t := p.CreatedAt
			if t.Month() == month && t.Day() == day && t.Year() < currentYear {
				results = append(results, p)
			}
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].CreatedAt.After(results[j].CreatedAt)
	})
	return results, nil
}
