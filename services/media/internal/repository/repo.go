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
	ErrMediaNotFound  = errors.New("media file not found")
	ErrStoryNotFound  = errors.New("story not found")
	ErrAlbumNotFound  = errors.New("album not found")
	ErrStreamNotFound = errors.New("live stream not found")
)

type MediaRepository interface {
	SaveMediaFile(ctx context.Context, media *domain.MediaFile) error
	GetMediaFile(ctx context.Context, id string) (*domain.MediaFile, error)

	SaveStory(ctx context.Context, story *domain.Story) error
	GetStory(ctx context.Context, id string) (*domain.Story, error)
	GetActiveStories(ctx context.Context, authorIDs []string, now time.Time) ([]*domain.Story, error)
	RecordStoryView(ctx context.Context, storyID, viewerID string, t time.Time) (bool, error)
	GetStoryViewers(ctx context.Context, storyID string) ([]*domain.StoryView, error)

	SaveReel(ctx context.Context, reel *domain.Reel) error
	ListReels(ctx context.Context, limit int) ([]*domain.Reel, error)

	SaveAlbum(ctx context.Context, album *domain.MediaAlbum) error
	GetAlbum(ctx context.Context, id string) (*domain.MediaAlbum, error)

	SaveLiveStream(ctx context.Context, stream *domain.LiveStream) error
	GetLiveStream(ctx context.Context, id string) (*domain.LiveStream, error)
}

type InMemoryMediaRepository struct {
	mu         sync.RWMutex
	mediaFiles map[string]*domain.MediaFile
	stories    map[string]*domain.Story
	storyViews map[string]map[string]*domain.StoryView // story_id -> viewer_id -> view
	reels      map[string]*domain.Reel
	albums     map[string]*domain.MediaAlbum
	streams    map[string]*domain.LiveStream
}

func NewInMemoryMediaRepository() *InMemoryMediaRepository {
	return &InMemoryMediaRepository{
		mediaFiles: make(map[string]*domain.MediaFile),
		stories:    make(map[string]*domain.Story),
		storyViews: make(map[string]map[string]*domain.StoryView),
		reels:      make(map[string]*domain.Reel),
		albums:     make(map[string]*domain.MediaAlbum),
		streams:    make(map[string]*domain.LiveStream),
	}
}

func (r *InMemoryMediaRepository) SaveMediaFile(ctx context.Context, media *domain.MediaFile) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.mediaFiles[media.ID] = media
	return nil
}

func (r *InMemoryMediaRepository) GetMediaFile(ctx context.Context, id string) (*domain.MediaFile, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	m, exists := r.mediaFiles[id]
	if !exists {
		return nil, ErrMediaNotFound
	}
	return m, nil
}

func (r *InMemoryMediaRepository) SaveStory(ctx context.Context, story *domain.Story) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.stories[story.ID] = story
	return nil
}

func (r *InMemoryMediaRepository) GetStory(ctx context.Context, id string) (*domain.Story, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	s, exists := r.stories[id]
	if !exists {
		return nil, ErrStoryNotFound
	}
	return s, nil
}

func (r *InMemoryMediaRepository) GetActiveStories(ctx context.Context, authorIDs []string, now time.Time) ([]*domain.Story, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	authorMap := make(map[string]bool)
	for _, id := range authorIDs {
		authorMap[id] = true
	}

	var results []*domain.Story
	for _, s := range r.stories {
		if (len(authorMap) == 0 || authorMap[s.AuthorID]) && !s.IsExpired(now) {
			results = append(results, s)
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].CreatedAt.After(results[j].CreatedAt)
	})
	return results, nil
}

func (r *InMemoryMediaRepository) RecordStoryView(ctx context.Context, storyID, viewerID string, t time.Time) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	story, exists := r.stories[storyID]
	if !exists {
		return false, ErrStoryNotFound
	}

	if _, ok := r.storyViews[storyID]; !ok {
		r.storyViews[storyID] = make(map[string]*domain.StoryView)
	}

	// Idempotent check
	if _, alreadyViewed := r.storyViews[storyID][viewerID]; alreadyViewed {
		return false, nil
	}

	r.storyViews[storyID][viewerID] = &domain.StoryView{
		StoryID:  storyID,
		ViewerID: viewerID,
		ViewedAt: t,
	}
	story.ViewCount++
	return true, nil
}

func (r *InMemoryMediaRepository) GetStoryViewers(ctx context.Context, storyID string) ([]*domain.StoryView, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	viewMap, exists := r.storyViews[storyID]
	if !exists {
		return []*domain.StoryView{}, nil
	}

	var results []*domain.StoryView
	for _, v := range viewMap {
		results = append(results, v)
	}
	sort.Slice(results, func(i, j int) bool {
		return results[i].ViewedAt.After(results[j].ViewedAt)
	})
	return results, nil
}

func (r *InMemoryMediaRepository) SaveReel(ctx context.Context, reel *domain.Reel) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.reels[reel.ID] = reel
	return nil
}

func (r *InMemoryMediaRepository) ListReels(ctx context.Context, limit int) ([]*domain.Reel, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var results []*domain.Reel
	for _, reel := range r.reels {
		results = append(results, reel)
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].CreatedAt.After(results[j].CreatedAt)
	})

	if limit > 0 && len(results) > limit {
		results = results[:limit]
	}
	return results, nil
}

func (r *InMemoryMediaRepository) SaveAlbum(ctx context.Context, album *domain.MediaAlbum) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.albums[album.ID] = album
	return nil
}

func (r *InMemoryMediaRepository) GetAlbum(ctx context.Context, id string) (*domain.MediaAlbum, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	a, exists := r.albums[id]
	if !exists {
		return nil, ErrAlbumNotFound
	}
	return a, nil
}

func (r *InMemoryMediaRepository) SaveLiveStream(ctx context.Context, stream *domain.LiveStream) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.streams[stream.ID] = stream
	return nil
}

func (r *InMemoryMediaRepository) GetLiveStream(ctx context.Context, id string) (*domain.LiveStream, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	st, exists := r.streams[id]
	if !exists {
		return nil, ErrStreamNotFound
	}
	return st, nil
}
