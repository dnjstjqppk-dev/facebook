package repository

import (
	"context"
	"errors"
	"sync"

	"github.com/dnjstjqppk-dev/facebook/pkg/domain"
)

var (
	ErrPageNotFound = errors.New("page not found")
)

type PageRepository interface {
	SavePage(ctx context.Context, page *domain.Page) error
	GetPage(ctx context.Context, id string) (*domain.Page, error)
	IncrementFollowers(ctx context.Context, pageID string) error
	GetPageInsights(ctx context.Context, pageID string) (*domain.PageInsights, error)
}

type InMemoryPageRepository struct {
	mu    sync.RWMutex
	pages map[string]*domain.Page
}

func NewInMemoryPageRepository() *InMemoryPageRepository {
	return &InMemoryPageRepository{
		pages: make(map[string]*domain.Page),
	}
}

func (r *InMemoryPageRepository) SavePage(ctx context.Context, page *domain.Page) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.pages[page.ID] = page
	return nil
}

func (r *InMemoryPageRepository) GetPage(ctx context.Context, id string) (*domain.Page, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, exists := r.pages[id]
	if !exists {
		return nil, ErrPageNotFound
	}
	return p, nil
}

func (r *InMemoryPageRepository) IncrementFollowers(ctx context.Context, pageID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	p, exists := r.pages[pageID]
	if !exists {
		return ErrPageNotFound
	}
	p.FollowersCount++
	return nil
}

func (r *InMemoryPageRepository) GetPageInsights(ctx context.Context, pageID string) (*domain.PageInsights, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, exists := r.pages[pageID]
	if !exists {
		return nil, ErrPageNotFound
	}

	return &domain.PageInsights{
		PageID:         pageID,
		FollowersCount: p.FollowersCount,
		TotalPostReach: p.FollowersCount * 8, // simulated reach
		EngagementRate: 4.8,
		Period:         "Last 28 Days",
	}, nil
}
