package tests

import (
	"testing"

	"github.com/dnjstjqppk-dev/facebook/pkg/domain"
	"github.com/dnjstjqppk-dev/facebook/services/advertising/internal/repository"
	"github.com/dnjstjqppk-dev/facebook/services/advertising/internal/service"
)

func TestAdCampaignHierarchyAndTargeting(t *testing.T) {
	repo := repository.NewInMemoryAdRepo()
	svc := service.NewAdvertisingService(repo)

	// Step 91: 3-Tier Hierarchy Creation
	acc, err := svc.CreateAdAccount("user_biz_1", "biz_portfolio_1", 100000, "USD")
	if err != nil {
		t.Fatalf("failed to create ad account: %v", err)
	}

	cmp, err := svc.CreateCampaign(acc.ID, "Summer Sale 2026", domain.AdObjectiveConversions, 5000)
	if err != nil {
		t.Fatalf("failed to create campaign: %v", err)
	}

	adSet, err := svc.CreateAdSet(cmp.ID, "Tech Geeks KR", 150, domain.AdTargeting{
		MinAge:    20,
		MaxAge:    35,
		Genders:   []string{"ALL"},
		Locations: []string{"KR"},
		Interests: []string{"gadgets", "software"},
	})
	if err != nil {
		t.Fatalf("failed to create ad set: %v", err)
	}

	ad, err := svc.CreateAd(adSet.ID, "NextGen AI Glasses", "Future on your eyes", "Experience next level AR.", "https://cdn.fb.test/glasses.jpg", "SHOP_NOW", "https://store.example.com")
	if err != nil {
		t.Fatalf("failed to create ad: %v", err)
	}

	if ad.ID == "" || ad.Status != domain.AdStatusActive {
		t.Errorf("expected active ad, got: %v", ad)
	}

	// Step 92: Targeting Engine Evaluation
	matchAudience := domain.UserAudienceProfile{
		UserID:    "user_25_kr",
		Age:       28,
		Gender:    "MALE",
		Location:  "KR",
		Interests: []string{"software", "gaming"},
	}
	if !svc.MatchesTargeting(matchAudience, adSet.Targeting) {
		t.Errorf("expected audience to match targeting")
	}

	mismatchAudience := domain.UserAudienceProfile{
		UserID:    "user_50_us",
		Age:       52,
		Gender:    "FEMALE",
		Location:  "US",
		Interests: []string{"cooking"},
	}
	if svc.MatchesTargeting(mismatchAudience, adSet.Targeting) {
		t.Errorf("expected audience to NOT match targeting")
	}
}

func TestAdAuction_VickreySecondPrice(t *testing.T) {
	repo := repository.NewInMemoryAdRepo()
	svc := service.NewAdvertisingService(repo)

	targeting := domain.AdTargeting{
		MinAge:    18,
		MaxAge:    65,
		Genders:   []string{"ALL"},
		Locations: []string{"KR"},
		Interests: []string{"shopping"},
	}

	ad1 := domain.Ad{ID: "ad_1", Name: "High Bid Ad"}
	ad2 := domain.Ad{ID: "ad_2", Name: "Low Bid Ad"}

	adSet1 := domain.AdSet{ID: "set_1", BidAmount: 200, Targeting: targeting}
	adSet2 := domain.AdSet{ID: "set_2", BidAmount: 120, Targeting: targeting}

	candidates := []domain.AuctionCandidate{
		{
			Ad:           ad1,
			AdSet:        adSet1,
			BidAmount:    200,
			EstimatedCTR: 0.05,
			QualityScore: 8.0,
		},
		{
			Ad:           ad2,
			AdSet:        adSet2,
			BidAmount:    120,
			EstimatedCTR: 0.04,
			QualityScore: 7.0,
		},
	}

	audience := domain.UserAudienceProfile{
		UserID:    "buyer_1",
		Age:       30,
		Gender:    "FEMALE",
		Location:  "KR",
		Interests: []string{"shopping"},
	}

	// Step 93: Real-Time Auction
	result, err := svc.RunAuction(candidates, audience)
	if err != nil {
		t.Fatalf("failed to run auction: %v", err)
	}

	if result.WinningAd.ID != "ad_1" {
		t.Errorf("expected ad_1 to win, got: %s", result.WinningAd.ID)
	}

	// Vickrey second price: ad_2's bid (120) + 1 = 121
	if result.PriceCharged != 121 {
		t.Errorf("expected Vickrey second price 121, got: %d", result.PriceCharged)
	}

	// Step 94: Real-time tracking of impression & click
	if err := svc.TrackImpression(result.WinningAd.ID, result.PriceCharged); err != nil {
		t.Fatalf("failed to track impression: %v", err)
	}
	if err := svc.TrackClick(result.WinningAd.ID); err != nil {
		t.Fatalf("failed to track click: %v", err)
	}

	metrics, err := svc.GetMetrics(result.WinningAd.ID)
	if err != nil {
		t.Fatalf("failed to get metrics: %v", err)
	}
	if metrics.Impressions != 1 || metrics.Clicks != 1 || metrics.Spend != 121 {
		t.Errorf("unexpected metrics: %+v", metrics)
	}
}

func TestAdLibraryAndBusinessSuite(t *testing.T) {
	repo := repository.NewInMemoryAdRepo()
	svc := service.NewAdvertisingService(repo)

	// Step 95: Ad Library Registration & Public Transparency Search
	_, _ = svc.CreateAd("adset_lib_1", "Public Election Transparency Ad", "Vote on May 10", "Know your rights.", "https://cdn.test/civic.png", "LEARN_MORE", "https://civic.org")

	entries, err := svc.SearchAdLibrary("Election")
	if err != nil {
		t.Fatalf("failed to search ad library: %v", err)
	}
	if len(entries) == 0 {
		t.Errorf("expected at least 1 ad library entry")
	}

	// Step 96: Meta Business Suite Partner Delegation
	role, err := svc.DelegatePartnerRole("biz_portfolio_100", "partner_agency_abc", "ADVERTISER")
	if err != nil {
		t.Fatalf("failed to delegate partner role: %v", err)
	}
	if role.Role != "ADVERTISER" {
		t.Errorf("expected role ADVERTISER, got %s", role.Role)
	}

	roles, err := svc.GetPartnerRoles("biz_portfolio_100")
	if err != nil || len(roles) != 1 {
		t.Errorf("expected 1 partner role, got %d (err: %v)", len(roles), err)
	}
}
