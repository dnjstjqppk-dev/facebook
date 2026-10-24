package repository

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/dnjstjqppk-dev/facebook/pkg/domain"
)

var (
	ErrRelationshipNotFound = errors.New("relationship record not found")
	ErrAlreadyExists        = errors.New("relationship already exists")
)

type SocialGraphRepository interface {
	SaveRelationship(ctx context.Context, rel *domain.Relationship) error
	GetRelationship(ctx context.Context, userID, targetID string, relType domain.RelationshipType) (*domain.Relationship, error)
	DeleteRelationship(ctx context.Context, userID, targetID string, relType domain.RelationshipType) error
	GetRelationshipsByUserAndType(ctx context.Context, userID string, relType domain.RelationshipType, status domain.RelationshipStatus) ([]*domain.Relationship, error)
	GetPendingReceivedRequests(ctx context.Context, targetID string) ([]*domain.Relationship, error)
	IsBlocked(ctx context.Context, userA, userB string) (bool, error)
}

type InMemorySocialGraphRepository struct {
	mu            sync.RWMutex
	relationships map[string]*domain.Relationship // key: userID:targetID:relType
}

func NewInMemorySocialGraphRepository() *InMemorySocialGraphRepository {
	return &InMemorySocialGraphRepository{
		relationships: make(map[string]*domain.Relationship),
	}
}

func makeKey(userID, targetID string, relType domain.RelationshipType) string {
	return userID + ":" + targetID + ":" + string(relType)
}

func (r *InMemorySocialGraphRepository) SaveRelationship(ctx context.Context, rel *domain.Relationship) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	key := makeKey(rel.UserID, rel.TargetUserID, rel.RelationshipType)
	rel.UpdatedAt = time.Now().UTC()
	r.relationships[key] = rel
	return nil
}

func (r *InMemorySocialGraphRepository) GetRelationship(ctx context.Context, userID, targetID string, relType domain.RelationshipType) (*domain.Relationship, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	key := makeKey(userID, targetID, relType)
	rel, exists := r.relationships[key]
	if !exists {
		return nil, ErrRelationshipNotFound
	}
	return rel, nil
}

func (r *InMemorySocialGraphRepository) DeleteRelationship(ctx context.Context, userID, targetID string, relType domain.RelationshipType) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	key := makeKey(userID, targetID, relType)
	delete(r.relationships, key)
	return nil
}

func (r *InMemorySocialGraphRepository) GetRelationshipsByUserAndType(ctx context.Context, userID string, relType domain.RelationshipType, status domain.RelationshipStatus) ([]*domain.Relationship, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var results []*domain.Relationship
	for _, rel := range r.relationships {
		if rel.UserID == userID && rel.RelationshipType == relType && (status == "" || rel.Status == status) {
			results = append(results, rel)
		}
	}
	return results, nil
}

func (r *InMemorySocialGraphRepository) GetPendingReceivedRequests(ctx context.Context, targetID string) ([]*domain.Relationship, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var results []*domain.Relationship
	for _, rel := range r.relationships {
		if rel.TargetUserID == targetID && rel.RelationshipType == domain.RelFriend && rel.Status == domain.RelStatusPending {
			results = append(results, rel)
		}
	}
	return results, nil
}

func (r *InMemorySocialGraphRepository) IsBlocked(ctx context.Context, userA, userB string) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	key1 := makeKey(userA, userB, domain.RelBlock)
	key2 := makeKey(userB, userA, domain.RelBlock)

	_, blocked1 := r.relationships[key1]
	_, blocked2 := r.relationships[key2]

	return blocked1 || blocked2, nil
}
