package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/dnjstjqppk-dev/facebook/pkg/domain"
	"github.com/dnjstjqppk-dev/facebook/services/payments/internal/repository"
)

var (
	ErrInvalidEscrowState = errors.New("invalid escrow transition state")
	ErrUnauthorizedParty  = errors.New("user is not authorized for this order action")
)

const SystemEscrowAccount = "SYSTEM_ESCROW_RESERVE"

type PaymentService struct {
	repo repository.PaymentRepository
}

func NewPaymentService(repo repository.PaymentRepository) *PaymentService {
	return &PaymentService{repo: repo}
}

func (s *PaymentService) CreateEscrowOrder(ctx context.Context, buyerID, sellerID, itemID string, amountCents int64, currency string) (*domain.EscrowOrder, error) {
	if currency == "" {
		currency = "USD"
	}
	now := time.Now().UTC()
	order := &domain.EscrowOrder{
		ID:          generateUUID(),
		ItemID:      itemID,
		BuyerID:     buyerID,
		SellerID:    sellerID,
		AmountCents: amountCents,
		Currency:    currency,
		Status:      domain.EscrowCreated,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.repo.SaveEscrowOrder(ctx, order); err != nil {
		return nil, err
	}
	return order, nil
}

func (s *PaymentService) FundEscrow(ctx context.Context, orderID string) (*domain.EscrowOrder, error) {
	order, err := s.repo.GetEscrowOrder(ctx, orderID)
	if err != nil {
		return nil, err
	}
	if order.Status != domain.EscrowCreated {
		return nil, ErrInvalidEscrowState
	}

	now := time.Now().UTC()
	order.Status = domain.EscrowHeld
	order.UpdatedAt = now

	// Double-entry bookkeeping:
	// 1. Debit buyer account
	_ = s.repo.AppendLedgerEntry(ctx, &domain.LedgerEntry{
		ID:          generateUUID(),
		AccountID:   order.BuyerID,
		EntryType:   domain.LedgerDebit,
		AmountCents: order.AmountCents,
		Currency:    order.Currency,
		ReferenceID: order.ID,
		Description: "Escrow funding payment",
		CreatedAt:   now,
	})

	// 2. Credit system escrow reserve
	_ = s.repo.AppendLedgerEntry(ctx, &domain.LedgerEntry{
		ID:          generateUUID(),
		AccountID:   SystemEscrowAccount,
		EntryType:   domain.LedgerCredit,
		AmountCents: order.AmountCents,
		Currency:    order.Currency,
		ReferenceID: order.ID,
		Description: "Held in escrow for order " + order.ID,
		CreatedAt:   now,
	})

	_ = s.repo.SaveEscrowOrder(ctx, order)
	return order, nil
}

func (s *PaymentService) ShipOrder(ctx context.Context, orderID, sellerID, trackingNumber string) (*domain.EscrowOrder, error) {
	order, err := s.repo.GetEscrowOrder(ctx, orderID)
	if err != nil {
		return nil, err
	}
	if order.SellerID != sellerID {
		return nil, ErrUnauthorizedParty
	}
	if order.Status != domain.EscrowHeld {
		return nil, ErrInvalidEscrowState
	}

	order.Status = domain.EscrowShipped
	order.TrackingNumber = trackingNumber
	order.UpdatedAt = time.Now().UTC()

	_ = s.repo.SaveEscrowOrder(ctx, order)
	return order, nil
}

func (s *PaymentService) ReleaseFunds(ctx context.Context, orderID, buyerID string) (*domain.EscrowOrder, error) {
	order, err := s.repo.GetEscrowOrder(ctx, orderID)
	if err != nil {
		return nil, err
	}
	if order.BuyerID != buyerID {
		return nil, ErrUnauthorizedParty
	}
	if order.Status != domain.EscrowShipped && order.Status != domain.EscrowHeld {
		return nil, ErrInvalidEscrowState
	}

	now := time.Now().UTC()
	order.Status = domain.EscrowReleased
	order.UpdatedAt = now

	// Double-entry bookkeeping:
	// 1. Debit system escrow reserve
	_ = s.repo.AppendLedgerEntry(ctx, &domain.LedgerEntry{
		ID:          generateUUID(),
		AccountID:   SystemEscrowAccount,
		EntryType:   domain.LedgerDebit,
		AmountCents: order.AmountCents,
		Currency:    order.Currency,
		ReferenceID: order.ID,
		Description: "Escrow release payout to seller",
		CreatedAt:   now,
	})

	// 2. Credit seller account
	_ = s.repo.AppendLedgerEntry(ctx, &domain.LedgerEntry{
		ID:          generateUUID(),
		AccountID:   order.SellerID,
		EntryType:   domain.LedgerCredit,
		AmountCents: order.AmountCents,
		Currency:    order.Currency,
		ReferenceID: order.ID,
		Description: "Item sale revenue for order " + order.ID,
		CreatedAt:   now,
	})

	_ = s.repo.SaveEscrowOrder(ctx, order)
	return order, nil
}

func (s *PaymentService) RaiseDispute(ctx context.Context, orderID, userID string) (*domain.EscrowOrder, error) {
	order, err := s.repo.GetEscrowOrder(ctx, orderID)
	if err != nil {
		return nil, err
	}
	if order.BuyerID != userID && order.SellerID != userID {
		return nil, ErrUnauthorizedParty
	}

	order.Status = domain.EscrowDisputed
	order.UpdatedAt = time.Now().UTC()
	_ = s.repo.SaveEscrowOrder(ctx, order)
	return order, nil
}

func (s *PaymentService) TipStars(ctx context.Context, senderID, creatorID string, stars int64) (*domain.CreatorWallet, error) {
	wallet, err := s.repo.GetWallet(ctx, creatorID)
	if err != nil {
		return nil, err
	}

	wallet.StarBalance += stars
	wallet.FiatEarningsCents += (stars * 1) // 1 star = $0.01 (1 cent) standard model
	wallet.UpdatedAt = time.Now().UTC()

	if err := s.repo.SaveWallet(ctx, wallet); err != nil {
		return nil, err
	}
	return wallet, nil
}

func (s *PaymentService) SubscribeToCreator(ctx context.Context, subscriberID, creatorID string, tierCents int64) (*domain.FanSubscription, error) {
	now := time.Now().UTC()
	sub := &domain.FanSubscription{
		ID:               generateUUID(),
		SubscriberID:     subscriberID,
		CreatorID:        creatorID,
		TierMonthlyCents: tierCents,
		IsActive:         true,
		CreatedAt:        now,
		RenewsAt:         now.AddDate(0, 1, 0),
	}

	if err := s.repo.SaveSubscription(ctx, sub); err != nil {
		return nil, err
	}
	return sub, nil
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
