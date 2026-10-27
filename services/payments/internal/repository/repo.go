package repository

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/dnjstjqppk-dev/facebook/pkg/domain"
)

var (
	ErrOrderNotFound        = errors.New("escrow order not found")
	ErrSubscriptionNotFound = errors.New("subscription not found")
)

type PaymentRepository interface {
	SaveEscrowOrder(ctx context.Context, order *domain.EscrowOrder) error
	GetEscrowOrder(ctx context.Context, id string) (*domain.EscrowOrder, error)

	AppendLedgerEntry(ctx context.Context, entry *domain.LedgerEntry) error
	GetLedgerEntriesByAccount(ctx context.Context, accountID string) ([]*domain.LedgerEntry, error)

	SaveWallet(ctx context.Context, wallet *domain.CreatorWallet) error
	GetWallet(ctx context.Context, userID string) (*domain.CreatorWallet, error)

	SaveSubscription(ctx context.Context, sub *domain.FanSubscription) error
	GetSubscription(ctx context.Context, id string) (*domain.FanSubscription, error)
}

type InMemoryPaymentRepository struct {
	mu            sync.RWMutex
	orders        map[string]*domain.EscrowOrder
	ledger        map[string][]*domain.LedgerEntry // accountID -> entries
	wallets       map[string]*domain.CreatorWallet // userID -> wallet
	subscriptions map[string]*domain.FanSubscription
}

func NewInMemoryPaymentRepository() *InMemoryPaymentRepository {
	return &InMemoryPaymentRepository{
		orders:        make(map[string]*domain.EscrowOrder),
		ledger:        make(map[string][]*domain.LedgerEntry),
		wallets:       make(map[string]*domain.CreatorWallet),
		subscriptions: make(map[string]*domain.FanSubscription),
	}
}

func (r *InMemoryPaymentRepository) SaveEscrowOrder(ctx context.Context, order *domain.EscrowOrder) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.orders[order.ID] = order
	return nil
}

func (r *InMemoryPaymentRepository) GetEscrowOrder(ctx context.Context, id string) (*domain.EscrowOrder, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	order, exists := r.orders[id]
	if !exists {
		return nil, ErrOrderNotFound
	}
	return order, nil
}

func (r *InMemoryPaymentRepository) AppendLedgerEntry(ctx context.Context, entry *domain.LedgerEntry) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.ledger[entry.AccountID] = append(r.ledger[entry.AccountID], entry)
	return nil
}

func (r *InMemoryPaymentRepository) GetLedgerEntriesByAccount(ctx context.Context, accountID string) ([]*domain.LedgerEntry, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.ledger[accountID], nil
}

func (r *InMemoryPaymentRepository) SaveWallet(ctx context.Context, wallet *domain.CreatorWallet) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.wallets[wallet.UserID] = wallet
	return nil
}

func (r *InMemoryPaymentRepository) GetWallet(ctx context.Context, userID string) (*domain.CreatorWallet, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	w, exists := r.wallets[userID]
	if !exists {
		return &domain.CreatorWallet{
			UserID:            userID,
			StarBalance:       0,
			FiatEarningsCents: 0,
			Currency:          "USD",
			UpdatedAt:         time.Now().UTC(),
		}, nil
	}
	return w, nil
}

func (r *InMemoryPaymentRepository) SaveSubscription(ctx context.Context, sub *domain.FanSubscription) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.subscriptions[sub.ID] = sub
	return nil
}

func (r *InMemoryPaymentRepository) GetSubscription(ctx context.Context, id string) (*domain.FanSubscription, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	sub, exists := r.subscriptions[id]
	if !exists {
		return nil, ErrSubscriptionNotFound
	}
	return sub, nil
}

