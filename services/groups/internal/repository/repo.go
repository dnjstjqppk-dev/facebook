package repository

import (
	"context"
	"errors"
	"sync"

	"github.com/dnjstjqppk-dev/facebook/pkg/domain"
)

var (
	ErrGroupNotFound     = errors.New("group not found")
	ErrGroupPostNotFound = errors.New("group post not found")
)

type GroupRepository interface {
	SaveGroup(ctx context.Context, group *domain.Group) error
	GetGroup(ctx context.Context, id string) (*domain.Group, error)
	SaveGroupPost(ctx context.Context, post *domain.GroupPost) error
	GetGroupPost(ctx context.Context, id string) (*domain.GroupPost, error)
	ListGroupPosts(ctx context.Context, groupID string, status domain.GroupPostStatus) ([]*domain.GroupPost, error)
}

type InMemoryGroupRepository struct {
	mu         sync.RWMutex
	groups     map[string]*domain.Group
	groupPosts map[string]*domain.GroupPost
}

func NewInMemoryGroupRepository() *InMemoryGroupRepository {
	return &InMemoryGroupRepository{
		groups:     make(map[string]*domain.Group),
		groupPosts: make(map[string]*domain.GroupPost),
	}
}

func (r *InMemoryGroupRepository) SaveGroup(ctx context.Context, group *domain.Group) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.groups[group.ID] = group
	return nil
}

func (r *InMemoryGroupRepository) GetGroup(ctx context.Context, id string) (*domain.Group, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	g, exists := r.groups[id]
	if !exists {
		return nil, ErrGroupNotFound
	}
	return g, nil
}

func (r *InMemoryGroupRepository) SaveGroupPost(ctx context.Context, post *domain.GroupPost) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.groupPosts[post.ID] = post
	return nil
}

func (r *InMemoryGroupRepository) GetGroupPost(ctx context.Context, id string) (*domain.GroupPost, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, exists := r.groupPosts[id]
	if !exists {
		return nil, ErrGroupPostNotFound
	}
	return p, nil
}

func (r *InMemoryGroupRepository) ListGroupPosts(ctx context.Context, groupID string, status domain.GroupPostStatus) ([]*domain.GroupPost, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var results []*domain.GroupPost
	for _, p := range r.groupPosts {
		if p.GroupID == groupID && (status == "" || p.Status == status) {
			results = append(results, p)
		}
	}
	return results, nil
}
