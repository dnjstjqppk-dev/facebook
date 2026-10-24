package repository

import (
	"context"
	"errors"
	"sync"

	"github.com/dnjstjqppk-dev/facebook/pkg/domain"
)

var (
	ErrProfileNotFound = errors.New("profile not found")
)

type ProfileRepository interface {
	SaveProfile(ctx context.Context, p *domain.Profile) error
	GetProfileByUserID(ctx context.Context, userID string) (*domain.Profile, error)
	GetProfileByUsername(ctx context.Context, username string) (*domain.Profile, error)
}

type InMemoryProfileRepository struct {
	mu         sync.RWMutex
	byUserID   map[string]*domain.Profile
	byUsername map[string]string // username -> userID
}

func NewInMemoryProfileRepository() *InMemoryProfileRepository {
	return &InMemoryProfileRepository{
		byUserID:   make(map[string]*domain.Profile),
		byUsername: make(map[string]string),
	}
}

func (r *InMemoryProfileRepository) SaveProfile(ctx context.Context, p *domain.Profile) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.byUserID[p.UserID] = p
	r.byUsername[p.Username] = p.UserID
	return nil
}

func (r *InMemoryProfileRepository) GetProfileByUserID(ctx context.Context, userID string) (*domain.Profile, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	p, exists := r.byUserID[userID]
	if !exists {
		return nil, ErrProfileNotFound
	}
	return p, nil
}

func (r *InMemoryProfileRepository) GetProfileByUsername(ctx context.Context, username string) (*domain.Profile, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	userID, exists := r.byUsername[username]
	if !exists {
		return nil, ErrProfileNotFound
	}
	return r.byUserID[userID], nil
}
