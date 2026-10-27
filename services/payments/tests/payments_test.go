package tests

import (
	"context"
	"testing"

	"github.com/dnjstjqppk-dev/facebook/pkg/domain"
	"github.com/dnjstjqppk-dev/facebook/services/payments/internal/repository"
	"github.com/dnjstjqppk-dev/facebook/services/payments/internal/service"
)

func setupTestPaymentService() (*service.PaymentService, repository.PaymentRepository) {
	repo := repository.NewInMemoryPaymentRepository()
	svc := service.NewPaymentService(repo)
	return svc, repo
}

func TestEscrowLifecycle_FullDoubleEntryLedger(t *testing.T) {
	svc, repo := setupTestPaymentService()
	ctx := context.Background()

	// 1. Create order
	order, err := svc.CreateEscrowOrder(ctx, "buyer-alex", "seller-sarah", "item-vintage-camera", 10000, "USD")
	if err != nil {
		t.Fatalf("create escrow order failed: %v", err)
	}
	if order.Status != domain.EscrowCreated {
		t.Errorf("expected status CREATED, got: %s", order.Status)
	}

	// 2. Fund escrow
	funded, err := svc.FundEscrow(ctx, order.ID)
	if err != nil || funded.Status != domain.EscrowHeld {
		t.Fatalf("fund escrow failed: %v", err)
	}

	// Verify buyer debited and system escrow credited
	buyerEntries, _ := repo.GetLedgerEntriesByAccount(ctx, "buyer-alex")
	if len(buyerEntries) != 1 || buyerEntries[0].EntryType != domain.LedgerDebit || buyerEntries[0].AmountCents != 10000 {
		t.Errorf("buyer ledger debit entry mismatch: %v", buyerEntries)
	}

	escrowEntries, _ := repo.GetLedgerEntriesByAccount(ctx, service.SystemEscrowAccount)
	if len(escrowEntries) != 1 || escrowEntries[0].EntryType != domain.LedgerCredit || escrowEntries[0].AmountCents != 10000 {
		t.Errorf("system escrow credit entry mismatch: %v", escrowEntries)
	}

	// 3. Seller ships order
	shipped, err := svc.ShipOrder(ctx, order.ID, "seller-sarah", "FEDEX-99887766")
	if err != nil || shipped.Status != domain.EscrowShipped || shipped.TrackingNumber != "FEDEX-99887766" {
		t.Fatalf("ship order failed: %v", err)
	}

	// 4. Buyer confirms delivery and releases funds
	released, err := svc.ReleaseFunds(ctx, order.ID, "buyer-alex")
	if err != nil || released.Status != domain.EscrowReleased {
		t.Fatalf("release funds failed: %v", err)
	}

	// Verify seller credited and system escrow debited
	sellerEntries, _ := repo.GetLedgerEntriesByAccount(ctx, "seller-sarah")
	if len(sellerEntries) != 1 || sellerEntries[0].EntryType != domain.LedgerCredit || sellerEntries[0].AmountCents != 10000 {
		t.Errorf("seller ledger credit entry mismatch: %v", sellerEntries)
	}

	escrowEntriesFinal, _ := repo.GetLedgerEntriesByAccount(ctx, service.SystemEscrowAccount)
	if len(escrowEntriesFinal) != 2 || escrowEntriesFinal[1].EntryType != domain.LedgerDebit {
		t.Errorf("system escrow debit entry on release mismatch: %v", escrowEntriesFinal)
	}
}

func TestEscrowDispute(t *testing.T) {
	svc, _ := setupTestPaymentService()
	ctx := context.Background()

	order, _ := svc.CreateEscrowOrder(ctx, "buyer-1", "seller-1", "item-1", 5000, "USD")
	_, _ = svc.FundEscrow(ctx, order.ID)

	disputed, err := svc.RaiseDispute(ctx, order.ID, "buyer-1")
	if err != nil || disputed.Status != domain.EscrowDisputed {
		t.Fatalf("raise dispute failed: %v", err)
	}
}

func TestCreator_StarsAndSubscriptions(t *testing.T) {
	svc, repo := setupTestPaymentService()
	ctx := context.Background()

	// 1. Tip 500 stars
	wallet, err := svc.TipStars(ctx, "fan-1", "creator-pro", 500)
	if err != nil || wallet.StarBalance != 500 || wallet.FiatEarningsCents != 500 {
		t.Fatalf("tip stars failed: %v", err)
	}

	w, _ := repo.GetWallet(ctx, "creator-pro")
	if w.StarBalance != 500 {
		t.Errorf("star balance mismatch in repository")
	}

	// 2. Monthly Fan Subscription ($4.99 / 499 cents)
	sub, err := svc.SubscribeToCreator(ctx, "fan-1", "creator-pro", 499)
	if err != nil || !sub.IsActive || sub.TierMonthlyCents != 499 {
		t.Fatalf("fan subscription failed: %v", err)
	}
}

