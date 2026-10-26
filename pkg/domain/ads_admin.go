package domain

import (
	"errors"
	"time"
)

// Advertising Objectives
const (
	AdObjectiveAwareness   = "AWARENESS"
	AdObjectiveTraffic     = "TRAFFIC"
	AdObjectiveEngagement  = "ENGAGEMENT"
	AdObjectiveConversions = "CONVERSIONS"
)

// Ad Status
const (
	AdStatusDraft    = "DRAFT"
	AdStatusActive   = "ACTIVE"
	AdStatusPaused   = "PAUSED"
	AdStatusArchived = "ARCHIVED"
	AdStatusRejected = "REJECTED"
)

// Targeting Definition
type AdTargeting struct {
	MinAge    int      `json:"min_age"`
	MaxAge    int      `json:"max_age"`
	Genders   []string `json:"genders"`   // "ALL", "MALE", "FEMALE"
	Locations []string `json:"locations"` // Country/City codes e.g. ["KR", "US"]
	Interests []string `json:"interests"` // e.g. ["technology", "gaming"]
}

// 91 & 92. Ad Account & 3-Tier Hierarchy (Campaign -> AdSet -> Ad)
type AdAccount struct {
	ID        string    `json:"id"`
	OwnerID   string    `json:"owner_id"`
	BusinessID string   `json:"business_id,omitempty"`
	Balance   int64     `json:"balance"` // In cents or microunits
	Currency  string    `json:"currency"`
	Status    string    `json:"status"` // ACTIVE, DISABLED
	CreatedAt time.Time `json:"created_at"`
}

type AdCampaign struct {
	ID          string    `json:"id"`
	AccountID   string    `json:"account_id"`
	Name        string    `json:"name"`
	Objective   string    `json:"objective"`
	DailyBudget int64     `json:"daily_budget"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type AdSet struct {
	ID         string      `json:"id"`
	CampaignID string      `json:"campaign_id"`
	Name       string      `json:"name"`
	BidAmount  int64       `json:"bid_amount"` // In cents
	Targeting  AdTargeting `json:"targeting"`
	Status     string      `json:"status"`
	CreatedAt  time.Time   `json:"created_at"`
}

type Ad struct {
	ID             string    `json:"id"`
	AdSetID        string    `json:"ad_set_id"`
	Name           string    `json:"name"`
	Headline       string    `json:"headline"`
	Body           string    `json:"body"`
	MediaURL       string    `json:"media_url"`
	CallToAction   string    `json:"call_to_action"` // "LEARN_MORE", "SHOP_NOW", "SIGN_UP"
	DestinationURL string    `json:"destination_url"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
}

// 93. Ad Auction & Bidding Engine (Vickrey Second-Price Auction)
type UserAudienceProfile struct {
	UserID    string   `json:"user_id"`
	Age       int      `json:"age"`
	Gender    string   `json:"gender"`
	Location  string   `json:"location"`
	Interests []string `json:"interests"`
}

type AuctionCandidate struct {
	Ad           Ad      `json:"ad"`
	AdSet        AdSet   `json:"ad_set"`
	BidAmount    int64   `json:"bid_amount"`
	EstimatedCTR float64 `json:"estimated_ctr"`
	QualityScore float64 `json:"quality_score"` // 1.0 to 10.0
	ComputedECPM float64 `json:"computed_ecpm"`
}

type AuctionResult struct {
	WinningAd   Ad      `json:"winning_ad"`
	WinningBid  int64   `json:"winning_bid"`
	PriceCharged int64  `json:"price_charged"` // Second price or cleared price
	ECPM        float64 `json:"ecpm"`
}

// 94. Ad Performance Analytics
type AdMetrics struct {
	AdID        string    `json:"ad_id"`
	Impressions int64     `json:"impressions"`
	Clicks      int64     `json:"clicks"`
	Spend       int64     `json:"spend"`
	Conversions int64     `json:"conversions"`
	LastUpdated time.Time `json:"last_updated"`
}

// 95. Public Ad Transparency Library Entry
type AdLibraryEntry struct {
	AdID             string      `json:"ad_id"`
	AdvertiserName   string      `json:"advertiser_name"`
	TargetingSummary AdTargeting `json:"targeting_summary"`
	TotalSpendRange  string      `json:"total_spend_range"` // e.g. "$1k - $5k"
	ImpressionsRange string      `json:"impressions_range"`
	IsPolitical      bool        `json:"is_political"`
	ActiveDates      string      `json:"active_dates"`
}

// 96. Meta Business Suite Partner Delegation
type BusinessPartnerRole struct {
	BusinessID string    `json:"business_id"`
	PartnerID  string    `json:"partner_id"`
	Role       string    `json:"role"` // "ADMIN", "ANALYST", "ADVERTISER"
	AssignedAt time.Time `json:"assigned_at"`
}

// 97. Moderation, Trust & Safety
const (
	ReportReasonSpam       = "SPAM"
	ReportReasonHateSpeech = "HATE_SPEECH"
	ReportReasonHarassment = "HARASSMENT"
	ReportReasonNudity     = "NUDITY"
	ReportReasonViolence   = "VIOLENCE"

	ReportStatusPending     = "PENDING"
	ReportStatusReviewed    = "REVIEWED"
	ReportStatusActionTaken = "ACTION_TAKEN"
	ReportStatusDismissed   = "DISMISSED"

	ModerationActionNone            = "NONE"
	ModerationActionContentRemoved  = "CONTENT_REMOVED"
	ModerationActionUserWarned      = "USER_WARNED"
	ModerationActionUserSuspended   = "USER_SUSPENDED"
	ModerationActionPermanentBan    = "PERMANENT_BAN"
)

type ModerationReport struct {
	ID          string    `json:"id"`
	ReporterID  string    `json:"reporter_id"`
	TargetType  string    `json:"target_type"` // POST, COMMENT, USER, AD
	TargetID    string    `json:"target_id"`
	Reason      string    `json:"reason"`
	Status      string    `json:"status"`
	ReviewerID  string    `json:"reviewer_id,omitempty"`
	ActionTaken string    `json:"action_taken"`
	Notes       string    `json:"notes,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	ResolvedAt  *time.Time `json:"resolved_at,omitempty"`
}

type UserStrike struct {
	ID             string    `json:"id"`
	UserID         string    `json:"user_id"`
	ViolationCount int       `json:"violation_count"`
	CurrentTier    string    `json:"current_tier"` // WARNING, TEMP_MUTE, SUSPENSION, BAN
	ExpiresAt      time.Time `json:"expires_at"`
}

// 98. Platform Super Admin & Kill Switch
type AuditLog struct {
	ID        string    `json:"id"`
	AdminID   string    `json:"admin_id"`
	Action    string    `json:"action"`
	Resource  string    `json:"resource"`
	IPAddress string    `json:"ip_address"`
	Details   string    `json:"details"`
	Timestamp time.Time `json:"timestamp"`
}

type KillSwitch struct {
	FeatureKey string    `json:"feature_key"` // e.g. "marketplace_payments", "live_streaming"
	Enabled    bool      `json:"enabled"`
	Reason     string    `json:"reason"`
	UpdatedBy  string    `json:"updated_by"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// 99. Developer Center & Graph API Token
type DeveloperApp struct {
	AppID       string    `json:"app_id"`
	OwnerID     string    `json:"owner_id"`
	AppName     string    `json:"app_name"`
	AppSecret   string    `json:"app_secret"`
	Scopes      []string  `json:"scopes"`
	RateTier    string    `json:"rate_tier"` // "STANDARD", "ENTERPRISE"
	CreatedAt   time.Time `json:"created_at"`
}

var (
	ErrAdAccountNotFound   = errors.New("ad account not found")
	ErrInsufficientBudget  = errors.New("insufficient budget or balance")
	ErrReportNotFound      = errors.New("moderation report not found")
	ErrKillSwitchTriggered = errors.New("service feature is currently disabled by system kill switch")
	ErrAppUnauthorized     = errors.New("developer app credentials unauthorized")
)
