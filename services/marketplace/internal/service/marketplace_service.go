package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"strings"
	"time"

	"github.com/dnjstjqppk-dev/facebook/pkg/domain"
	"github.com/dnjstjqppk-dev/facebook/services/marketplace/internal/repository"
)

var (
	ErrProhibitedItem = errors.New("listing contains prohibited goods violating commerce safety policy")
	ErrNotSeller      = errors.New("only seller can accept offers for this item")
	ErrInvalidOffer   = errors.New("offer amount must be greater than zero")
)

var prohibitedKeywords = []string{
	"gun", "weapon", "firearm", "ammo", "drug", "narcotic", "prescription", "counterfeit", "replica rolex",
}

type MarketplaceService struct {
	repo repository.MarketplaceRepository
}

func NewMarketplaceService(repo repository.MarketplaceRepository) *MarketplaceService {
	return &MarketplaceService{repo: repo}
}

func (s *MarketplaceService) ListItem(ctx context.Context, sellerID, title, description, category, location, currency string, priceCents int64, mediaURLs []string) (*domain.MarketplaceItem, error) {
	// Policy Enforcement
	contentLower := strings.ToLower(title + " " + description)
	for _, kw := range prohibitedKeywords {
		if strings.Contains(contentLower, kw) {
			return nil, ErrProhibitedItem
		}
	}

	if currency == "" {
		currency = "USD"
	}

	now := time.Now().UTC()
	item := &domain.MarketplaceItem{
		ID:          generateUUID(),
		SellerID:    sellerID,
		Title:       title,
		Description: description,
		Category:    category,
		PriceCents:  priceCents,
		Currency:    currency,
		Location:    location,
		MediaURLs:   mediaURLs,
		Status:      domain.ItemStatusAvailable,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.repo.SaveItem(ctx, item); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *MarketplaceService) MakeOffer(ctx context.Context, buyerID, itemID string, offerAmountCents int64) (*domain.ItemOffer, error) {
	if offerAmountCents <= 0 {
		return nil, ErrInvalidOffer
	}

	item, err := s.repo.GetItem(ctx, itemID)
	if err != nil {
		return nil, err
	}

	if item.Status != domain.ItemStatusAvailable {
		return nil, errors.New("item is no longer available for offers")
	}

	offer := &domain.ItemOffer{
		ID:               generateUUID(),
		ItemID:           itemID,
		BuyerID:          buyerID,
		OfferAmountCents: offerAmountCents,
		Status:           domain.OfferPending,
		CreatedAt:        time.Now().UTC(),
	}

	if err := s.repo.SaveOffer(ctx, offer); err != nil {
		return nil, err
	}
	return offer, nil
}

func (s *MarketplaceService) AcceptOffer(ctx context.Context, sellerID, offerID string) error {
	offer, err := s.repo.GetOffer(ctx, offerID)
	if err != nil {
		return err
	}

	item, err := s.repo.GetItem(ctx, offer.ItemID)
	if err != nil {
		return err
	}

	if item.SellerID != sellerID {
		return ErrNotSeller
	}

	offer.Status = domain.OfferAccepted
	item.Status = domain.ItemStatusPending
	item.UpdatedAt = time.Now().UTC()

	_ = s.repo.SaveOffer(ctx, offer)
	return s.repo.SaveItem(ctx, item)
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

