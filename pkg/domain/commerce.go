package domain

import "time"

// ============================================================================
// Marketplace
// ============================================================================

type ItemStatus string

const (
	ItemStatusAvailable ItemStatus = "AVAILABLE"
	ItemStatusPending   ItemStatus = "PENDING"
	ItemStatusSold      ItemStatus = "SOLD"
)

type MarketplaceItem struct {
	ID          string     `json:"id"`
	SellerID    string     `json:"seller_id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Category    string     `json:"category"`
	PriceCents  int64      `json:"price_cents"`
	Currency    string     `json:"currency"`
	Location    string     `json:"location"`
	MediaURLs   []string   `json:"media_urls,omitempty"`
	Status      ItemStatus `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type OfferStatus string

const (
	OfferPending  OfferStatus = "PENDING"
	OfferAccepted OfferStatus = "ACCEPTED"
	OfferDeclined OfferStatus = "DECLINED"
)

type ItemOffer struct {
	ID               string      `json:"id"`
	ItemID           string      `json:"item_id"`
	BuyerID          string      `json:"buyer_id"`
	OfferAmountCents int64       `json:"offer_amount_cents"`
	Status           OfferStatus `json:"status"`
	CreatedAt        time.Time   `json:"created_at"`
}

// ============================================================================
// Payments & Escrow Ledger
// ============================================================================

type EscrowStatus string

const (
	EscrowCreated   EscrowStatus = "CREATED"
	EscrowHeld      EscrowStatus = "FUNDS_HELD"
	EscrowShipped   EscrowStatus = "SHIPPED"
	EscrowDelivered EscrowStatus = "DELIVERED"
	EscrowReleased  EscrowStatus = "RELEASED"
	EscrowRefunded  EscrowStatus = "REFUNDED"
	EscrowDisputed  EscrowStatus = "DISPUTED"
)

type EscrowOrder struct {
	ID             string       `json:"id"`
	ItemID         string       `json:"item_id"`
	BuyerID        string       `json:"buyer_id"`
	SellerID       string       `json:"seller_id"`
	AmountCents    int64        `json:"amount_cents"`
	Currency       string       `json:"currency"`
	Status         EscrowStatus `json:"status"`
	TrackingNumber string       `json:"tracking_number,omitempty"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
}

type LedgerEntryType string

const (
	LedgerCredit LedgerEntryType = "CREDIT"
	LedgerDebit  LedgerEntryType = "DEBIT"
)

type LedgerEntry struct {
	ID             string          `json:"id"`
	AccountID      string          `json:"account_id"`
	EntryType      LedgerEntryType `json:"entry_type"`
	AmountCents    int64           `json:"amount_cents"`
	Currency       string          `json:"currency"`
	ReferenceID    string          `json:"reference_id"`
	Description    string          `json:"description"`
	CreatedAt      time.Time       `json:"created_at"`
}

// ============================================================================
// Creator Economy: Stars & Subscriptions
// ============================================================================

type CreatorWallet struct {
	UserID            string    `json:"user_id"`
	StarBalance       int64     `json:"star_balance"`
	FiatEarningsCents int64     `json:"fiat_earnings_cents"`
	Currency          string    `json:"currency"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type StarTip struct {
	ID        string    `json:"id"`
	SenderID  string    `json:"sender_id"`
	CreatorID string    `json:"creator_id"`
	Stars     int64     `json:"stars"`
	Message   string    `json:"message,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type FanSubscription struct {
	ID               string    `json:"id"`
	SubscriberID     string    `json:"subscriber_id"`
	CreatorID        string    `json:"creator_id"`
	TierMonthlyCents int64     `json:"tier_monthly_cents"`
	IsActive         bool      `json:"is_active"`
	CreatedAt        time.Time `json:"created_at"`
	RenewsAt         time.Time `json:"renews_at"`
}
