package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/dnjstjqppk-dev/facebook/pkg/domain"
	"github.com/dnjstjqppk-dev/facebook/services/feed/internal/repository"
)

var (
	ErrEmptyPostContent    = errors.New("post content cannot be empty without media")
	ErrUnauthorizedPost    = errors.New("not authorized to modify this post")
	ErrUnauthorizedComment = errors.New("not authorized to modify this comment")
	ErrInvalidReaction     = errors.New("invalid reaction type")
	ErrMaxNestingExceeded  = errors.New("maximum 2-tier comment nesting exceeded")
)

type FriendLister interface {
	GetFriendIDs(ctx context.Context, userID string) ([]string, error)
}

type FeedService struct {
	repo         repository.FeedRepository
	friendLister FriendLister
}

func NewFeedService(repo repository.FeedRepository, friendLister FriendLister) *FeedService {
	return &FeedService{
		repo:         repo,
		friendLister: friendLister,
	}
}

func (s *FeedService) CreatePost(ctx context.Context, authorID string, req *domain.CreatePostRequest) (*domain.Post, error) {
	if strings.TrimSpace(req.Content) == "" && len(req.MediaURLs) == 0 {
		return nil, ErrEmptyPostContent
	}

	privacy := req.Privacy
	if privacy == "" {
		privacy = domain.PrivacyFriends
	}

	postID := generateUUID()
	now := time.Now().UTC()

	post := &domain.Post{
		ID:             postID,
		AuthorID:       authorID,
		Content:        req.Content,
		MediaURLs:      req.MediaURLs,
		Feeling:        req.Feeling,
		Location:       req.Location,
		TaggedUserIDs:  req.TaggedUserIDs,
		Privacy:        privacy,
		OriginalPostID: req.OriginalPostID,
		ReactionCounts: make(map[domain.ReactionType]int),
		CommentCount:   0,
		ShareCount:     0,
		EditHistory:    []domain.EditHistoryEntry{},
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	if err := s.repo.CreatePost(ctx, post); err != nil {
		return nil, err
	}

	// Fan-out to author's own timeline
	_ = s.repo.AddToTimeline(ctx, authorID, postID, now)

	// Fan-out-on-write to friends' timelines if privacy is PUBLIC or FRIENDS
	if s.friendLister != nil && privacy != domain.PrivacyOnlyMe {
		friendIDs, err := s.friendLister.GetFriendIDs(ctx, authorID)
		if err == nil {
			for _, fid := range friendIDs {
				_ = s.repo.AddToTimeline(ctx, fid, postID, now)
			}
		}
	}

	return post, nil
}

func (s *FeedService) UpdatePost(ctx context.Context, authorID, postID string, req *domain.UpdatePostRequest) (*domain.Post, error) {
	post, err := s.repo.GetPost(ctx, postID)
	if err != nil {
		return nil, err
	}

	if post.AuthorID != authorID {
		return nil, ErrUnauthorizedPost
	}

	now := time.Now().UTC()
	// Record edit history if content changed
	if post.Content != req.Content {
		post.EditHistory = append(post.EditHistory, domain.EditHistoryEntry{
			PreviousContent: post.Content,
			EditedAt:        now,
		})
		post.Content = req.Content
	}

	if req.Privacy != "" {
		post.Privacy = req.Privacy
	}
	if req.MediaURLs != nil {
		post.MediaURLs = req.MediaURLs
	}
	post.UpdatedAt = now

	if err := s.repo.UpdatePost(ctx, post); err != nil {
		return nil, err
	}
	return post, nil
}

func (s *FeedService) DeletePost(ctx context.Context, authorID, postID string) error {
	post, err := s.repo.GetPost(ctx, postID)
	if err != nil {
		return err
	}

	if post.AuthorID != authorID {
		return ErrUnauthorizedPost
	}

	return s.repo.DeletePost(ctx, postID)
}

func (s *FeedService) ReactToPost(ctx context.Context, userID, postID string, reactionType domain.ReactionType) error {
	switch reactionType {
	case domain.ReactionLike, domain.ReactionLove, domain.ReactionCare,
		domain.ReactionHaha, domain.ReactionWow, domain.ReactionSad, domain.ReactionAngry:
	default:
		return ErrInvalidReaction
	}

	reaction := &domain.Reaction{
		PostID:       postID,
		UserID:       userID,
		ReactionType: reactionType,
		CreatedAt:    time.Now().UTC(),
	}

	return s.repo.SaveReaction(ctx, reaction)
}

func (s *FeedService) RemoveReaction(ctx context.Context, userID, postID string) error {
	return s.repo.DeleteReaction(ctx, postID, userID)
}

func (s *FeedService) AddComment(ctx context.Context, authorID, postID string, req *domain.AddCommentRequest) (*domain.Comment, error) {
	if strings.TrimSpace(req.Content) == "" {
		return nil, errors.New("comment content cannot be empty")
	}

	// Enforce 2-tier maximum nesting (Post -> Root Comment -> Reply)
	parentID := req.ParentCommentID
	if parentID != "" {
		parent, err := s.repo.GetComment(ctx, parentID)
		if err != nil {
			return nil, err
		}
		// If the parent is already a reply, collapse under the root comment
		if parent.ParentCommentID != "" {
			parentID = parent.ParentCommentID
		}
	}

	now := time.Now().UTC()
	comment := &domain.Comment{
		ID:              generateUUID(),
		PostID:          postID,
		AuthorID:        authorID,
		ParentCommentID: parentID,
		Content:         req.Content,
		Mentions:        req.Mentions,
		ReactionCounts:  make(map[domain.ReactionType]int),
		RepliesCount:    0,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	if err := s.repo.CreateComment(ctx, comment); err != nil {
		return nil, err
	}
	return comment, nil
}

func (s *FeedService) GetComments(ctx context.Context, postID string) ([]*domain.Comment, error) {
	return s.repo.GetCommentsByPost(ctx, postID)
}

func (s *FeedService) GetFeed(ctx context.Context, userID, cursor string, limit int, mode string) (*domain.FeedResponse, error) {
	if limit <= 0 || limit > 50 {
		limit = 10
	}

	timeline, err := s.repo.GetTimeline(ctx, userID)
	if err != nil {
		return nil, err
	}

	// If user timeline has fewer entries, also merge public/recent posts
	seen := make(map[string]bool)
	var candidatePosts []*domain.Post

	for _, entry := range timeline {
		if seen[entry.PostID] {
			continue
		}
		seen[entry.PostID] = true
		p, err := s.repo.GetPost(ctx, entry.PostID)
		if err == nil && p.DeletedAt == nil {
			candidatePosts = append(candidatePosts, p)
		}
	}

	now := time.Now().UTC()
	var feedItems []*domain.FeedItem

	for _, p := range candidatePosts {
		item := &domain.FeedItem{
			Post: p,
		}
		if uReact, ok := s.repo.GetUserReaction(ctx, p.ID, userID); ok {
			item.UserReaction = uReact
		}

		if mode == "ranked" {
			// Edge-rank score formula:
			// Engagement = (likes * 1) + (comments * 2) + (shares * 3)
			totalReactions := 0
			for _, count := range p.ReactionCounts {
				totalReactions += count
			}
			engagement := float64(totalReactions)*1.0 + float64(p.CommentCount)*2.0 + float64(p.ShareCount)*3.0
			ageHours := now.Sub(p.CreatedAt).Hours()
			if ageHours < 0 {
				ageHours = 0
			}
			gravity := 1.5
			item.RankingScore = (engagement + 1.0) / math.Pow(ageHours+2.0, gravity)
		} else {
			// Chronological: timestamp as score
			item.RankingScore = float64(p.CreatedAt.Unix())
		}
		feedItems = append(feedItems, item)
	}

	// Sort items according to score descending
	sort.Slice(feedItems, func(i, j int) bool {
		return feedItems[i].RankingScore > feedItems[j].RankingScore
	})

	// Cursor decoding and pagination
	startIndex := 0
	if cursor != "" {
		if decoded, err := base64.URLEncoding.DecodeString(cursor); err == nil {
			parts := strings.Split(string(decoded), "_")
			if len(parts) == 2 {
				lastID := parts[1]
				for idx, item := range feedItems {
					if item.Post.ID == lastID {
						startIndex = idx + 1
						break
					}
				}
			}
		}
	}

	endIndex := startIndex + limit
	hasMore := false
	if endIndex < len(feedItems) {
		hasMore = true
	} else {
		endIndex = len(feedItems)
	}

	var pagedItems []*domain.FeedItem
	if startIndex < len(feedItems) {
		pagedItems = feedItems[startIndex:endIndex]
	}

	var nextCursor string
	if hasMore && len(pagedItems) > 0 {
		lastItem := pagedItems[len(pagedItems)-1]
		cursorRaw := fmt.Sprintf("%d_%s", lastItem.Post.CreatedAt.UnixMilli(), lastItem.Post.ID)
		nextCursor = base64.URLEncoding.EncodeToString([]byte(cursorRaw))
	}

	return &domain.FeedResponse{
		Items:      pagedItems,
		NextCursor: nextCursor,
		HasMore:    hasMore,
	}, nil
}

func (s *FeedService) SavePost(ctx context.Context, userID, collectionName, postID string) error {
	if _, err := s.repo.GetPost(ctx, postID); err != nil {
		return err
	}
	if collectionName == "" {
		collectionName = "All Saved"
	}
	return s.repo.SaveToCollection(ctx, userID, collectionName, postID)
}

func (s *FeedService) GetSavedCollections(ctx context.Context, userID string) ([]*domain.SavedCollection, error) {
	return s.repo.GetSavedCollections(ctx, userID)
}

func (s *FeedService) GetMemories(ctx context.Context, userID string) ([]*domain.Post, error) {
	now := time.Now().UTC()
	return s.repo.GetMemories(ctx, userID, now.Month(), now.Day(), now.Year())
}

func generateUUID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return hex.EncodeToString(b[0:4]) + "-" +
		hex.EncodeToString(b[4:6]) + "-" +
		hex.EncodeToString(b[6:8]) + "-" +
		hex.EncodeToString(b[8:10]) + "-" +
		hex.EncodeToString(b[10:16])
}

func formatInt(i int) string {
	return strconv.Itoa(i)
}
