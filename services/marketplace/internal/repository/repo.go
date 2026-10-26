package repository

import (
	"context"
	"errors"
	"sync"

	"github.com/dnjstjqppk-dev/facebook/pkg/domain"
)

var (
	ErrItemNotFound  = errors.New("marketplace item not found")
	ErrOfferNotFound = errors.New("offer not found")
)

type MarketplaceRepository interface {
	SaveItem(ctx context.Context, item *domain.MarketplaceItem) error
	GetItem(ctx context.Context, id string) (*domain.MarketplaceItem, error)
	ListItemsByCategory(ctx context.Context, category string) ([]*domain.MarketplaceItem, error)
	SaveOffer(ctx context.Context, offer *domain.ItemOffer) error
	GetOffer(ctx context.Context, id string) (*domain.ItemOffer, error)
}

type InMemoryMarketplaceRepository struct {
	mu     sync.RWMutex
	items  map[string]*domain.MarketplaceItem
	offers map[string]*domain.ItemOffer
}

func NewInMemoryMarketplaceRepository() *InMemoryMarketplaceRepository {
	return &InMemoryMarketplaceRepository{
		items:  make(map[string]*domain.MarketplaceItem),
		offers: make(map[string]*domain.ItemOffer),
	}
}

func (r *InMemoryMarketplaceRepository) SaveItem(ctx context.Context, item *domain.MarketplaceItem) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.items[item.ID] = item
	return nil
}

func (r *InMemoryMarketplaceRepository) GetItem(ctx context.Context, id string) (*domain.MarketplaceItem, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	item, exists := r.items[id]
	if !exists {
		return nil, ErrItemNotFound
	}
	return item, nil
}

func (r *InMemoryMarketplaceRepository) ListItemsByCategory(ctx context.Context, category string) ([]*domain.MarketplaceItem, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var results []*domain.MarketplaceItem
	for _, item := range r.items {
		if (category == "" || item.Category == category) && item.Status == domain.ItemStatusAvailable {
			results = append(results, item)
		}
	}
	return results, nil
}

func (r *InMemoryMarketplaceRepository) SaveOffer(ctx context.Context, offer *domain.ItemOffer) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.offers[offer.ID] = offer
	return nil
}

func (r *InMemoryMarketplaceRepository) GetOffer(ctx context.Context, id string) (*domain.ItemOffer, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	offer, exists := r.offers[id]
	if !exists {
		return nil, ErrOfferNotFound
	}
	return offer, nil
}
