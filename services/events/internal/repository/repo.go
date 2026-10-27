package repository

import (
	"context"
	"errors"
	"sync"

	"github.com/dnjstjqppk-dev/facebook/pkg/domain"
)

var (
	ErrEventNotFound = errors.New("event not found")
)

type EventRepository interface {
	SaveEvent(ctx context.Context, event *domain.Event) error
	GetEvent(ctx context.Context, id string) (*domain.Event, error)
	ListEventsByUser(ctx context.Context, userID string) ([]*domain.Event, error)
}

type InMemoryEventRepository struct {
	mu     sync.RWMutex
	events map[string]*domain.Event
}

func NewInMemoryEventRepository() *InMemoryEventRepository {
	return &InMemoryEventRepository{
		events: make(map[string]*domain.Event),
	}
}

func (r *InMemoryEventRepository) SaveEvent(ctx context.Context, event *domain.Event) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.events[event.ID] = event
	return nil
}

func (r *InMemoryEventRepository) GetEvent(ctx context.Context, id string) (*domain.Event, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	e, exists := r.events[id]
	if !exists {
		return nil, ErrEventNotFound
	}
	return e, nil
}

func (r *InMemoryEventRepository) ListEventsByUser(ctx context.Context, userID string) ([]*domain.Event, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var results []*domain.Event
	for _, e := range r.events {
		if e.HostID == userID {
			results = append(results, e)
		} else if status, ok := e.RSVPs[userID]; ok && (status == domain.RSVPGoing || status == domain.RSVPInterested) {
			results = append(results, e)
		}
	}
	return results, nil
}

