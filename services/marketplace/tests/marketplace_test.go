package tests

import (
	"context"
	"testing"

	"github.com/dnjstjqppk-dev/facebook/pkg/domain"
	"github.com/dnjstjqppk-dev/facebook/services/marketplace/internal/repository"
	"github.com/dnjstjqppk-dev/facebook/services/marketplace/internal/service"
)

func setupTestMarketplaceService() (*service.MarketplaceService, repository.MarketplaceRepository) {
	repo := repository.NewInMemoryMarketplaceRepository()
	svc := service.NewMarketplaceService(repo)
	return svc, repo
}

func TestListItem_ProhibitedPolicyEnforcement(t *testing.T) {
	svc, _ := setupTestMarketplaceService()
	ctx := context.Background()

	// Prohibited firearm attempt
	_, err := svc.ListItem(ctx, "seller-bad", "Hunting Rifle Firearm", "Almost new rifle", "Sporting Goods", "Texas", "USD", 45000, nil)
	if err != service.ErrProhibitedItem {
		t.Errorf("expected ErrProhibitedItem, got: %v", err)
	}

	// Counterfeit attempt
	_, err = svc.ListItem(ctx, "seller-bad", "Replica Rolex Luxury Watch", "1:1 exact copy", "Jewelry", "New York", "USD", 15000, nil)
	if err != service.ErrProhibitedItem {
		t.Errorf("expected ErrProhibitedItem for counterfeit goods, got: %v", err)
	}
}

func TestListItem_Valid(t *testing.T) {
	svc, repo := setupTestMarketplaceService()
	ctx := context.Background()

	item, err := svc.ListItem(ctx, "seller-1", "MacBook Pro 16 M3 Max", "Excellent condition, 36GB RAM", "Electronics", "Seoul, Gangnam", "KRW", 320000000, []string{"https://cdn.example.com/macbook.jpg"})
	if err != nil {
		t.Fatalf("failed to list item: %v", err)
	}

	if item.Status != domain.ItemStatusAvailable || item.SellerID != "seller-1" {
		t.Errorf("item fields mismatch: %v", item)
	}

	retrieved, _ := repo.GetItem(ctx, item.ID)
	if retrieved.Title != "MacBook Pro 16 M3 Max" {
		t.Errorf("item title mismatch in storage")
	}
}

func TestMakeAndAcceptOffer(t *testing.T) {
	svc, repo := setupTestMarketplaceService()
	ctx := context.Background()

	item, _ := svc.ListItem(ctx, "seller-1", "Vintage Acoustic Guitar", "Warm tone", "Musical Instruments", "Austin", "USD", 50000, nil)

	// Buyer makes offer $450 (45000 cents)
	offer, err := svc.MakeOffer(ctx, "buyer-1", item.ID, 45000)
	if err != nil {
		t.Fatalf("failed to make offer: %v", err)
	}
	if offer.Status != domain.OfferPending {
		t.Errorf("expected offer status PENDING")
	}

	// Non-seller tries to accept offer
	err = svc.AcceptOffer(ctx, "imposter-99", offer.ID)
	if err != service.ErrNotSeller {
		t.Errorf("expected ErrNotSeller, got: %v", err)
	}

	// Seller accepts offer
	err = svc.AcceptOffer(ctx, "seller-1", offer.ID)
	if err != nil {
		t.Fatalf("seller failed to accept offer: %v", err)
	}

	it, _ := repo.GetItem(ctx, item.ID)
	if it.Status != domain.ItemStatusPending {
		t.Errorf("expected item status PENDING after offer accepted, got %s", it.Status)
	}
}

