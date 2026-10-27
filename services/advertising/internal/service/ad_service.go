package service

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/dnjstjqppk-dev/facebook/pkg/domain"
	"github.com/dnjstjqppk-dev/facebook/services/advertising/internal/repository"
)

type AdvertisingService struct {
	repo repository.AdRepository
}

func NewAdvertisingService(repo repository.AdRepository) *AdvertisingService {
	return &AdvertisingService{repo: repo}
}

// 91. 3-Tier Hierarchy: Create Ad Account, Campaign, AdSet, and Ad
func (s *AdvertisingService) CreateAdAccount(ownerID, businessID string, initialBalance int64, currency string) (*domain.AdAccount, error) {
	if initialBalance < 0 {
		return nil, errors.New("initial balance cannot be negative")
	}
	acc := &domain.AdAccount{
		ID:         fmt.Sprintf("act_%d", time.Now().UnixNano()),
		OwnerID:    ownerID,
		BusinessID: businessID,
		Balance:    initialBalance,
		Currency:   currency,
		Status:     "ACTIVE",
		CreatedAt:  time.Now(),
	}
	if err := s.repo.SaveAccount(acc); err != nil {
		return nil, err
	}
	return acc, nil
}

func (s *AdvertisingService) CreateCampaign(accountID, name, objective string, dailyBudget int64) (*domain.AdCampaign, error) {
	if dailyBudget <= 0 {
		return nil, domain.ErrInsufficientBudget
	}
	cmp := &domain.AdCampaign{
		ID:          fmt.Sprintf("cmp_%d", time.Now().UnixNano()),
		AccountID:   accountID,
		Name:        name,
		Objective:   objective,
		DailyBudget: dailyBudget,
		Status:      domain.AdStatusActive,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	if err := s.repo.SaveCampaign(cmp); err != nil {
		return nil, err
	}
	return cmp, nil
}

func (s *AdvertisingService) CreateAdSet(campaignID, name string, bidAmount int64, targeting domain.AdTargeting) (*domain.AdSet, error) {
	if bidAmount <= 0 {
		return nil, errors.New("bid amount must be positive")
	}
	adSet := &domain.AdSet{
		ID:         fmt.Sprintf("adset_%d", time.Now().UnixNano()),
		CampaignID: campaignID,
		Name:       name,
		BidAmount:  bidAmount,
		Targeting:  targeting,
		Status:     domain.AdStatusActive,
		CreatedAt:  time.Now(),
	}
	if err := s.repo.SaveAdSet(adSet); err != nil {
		return nil, err
	}
	return adSet, nil
}

func (s *AdvertisingService) CreateAd(adSetID, name, headline, body, mediaURL, cta, destURL string) (*domain.Ad, error) {
	if headline == "" || destURL == "" {
		return nil, errors.New("headline and destination URL are required")
	}
	ad := &domain.Ad{
		ID:             fmt.Sprintf("ad_%d", time.Now().UnixNano()),
		AdSetID:        adSetID,
		Name:           name,
		Headline:       headline,
		Body:           body,
		MediaURL:       mediaURL,
		CallToAction:   cta,
		DestinationURL: destURL,
		Status:         domain.AdStatusActive,
		CreatedAt:      time.Now(),
	}
	if err := s.repo.SaveAd(ad); err != nil {
		return nil, err
	}

	// Auto-register in Public Ad Library (Phase 10: Step 95)
	adSet, _ := s.repo.GetAdSet(adSetID)
	var targeting domain.AdTargeting
	if adSet != nil {
		targeting = adSet.Targeting
	}
	_ = s.repo.SaveAdLibraryEntry(&domain.AdLibraryEntry{
		AdID:             ad.ID,
		AdvertiserName:   name,
		TargetingSummary: targeting,
		TotalSpendRange:  "< $500",
		ImpressionsRange: "1k - 10k",
		IsPolitical:      false,
		ActiveDates:      time.Now().Format("2006-01-02"),
	})

	return ad, nil
}

// 92. Targeting Match Engine
func (s *AdvertisingService) MatchesTargeting(audience domain.UserAudienceProfile, targeting domain.AdTargeting) bool {
	// Age validation
	if targeting.MinAge > 0 && audience.Age < targeting.MinAge {
		return false
	}
	if targeting.MaxAge > 0 && audience.Age > targeting.MaxAge {
		return false
	}

	// Gender validation
	if len(targeting.Genders) > 0 {
		matchedGender := false
		for _, g := range targeting.Genders {
			if strings.EqualFold(g, "ALL") || strings.EqualFold(g, audience.Gender) {
				matchedGender = true
				break
			}
		}
		if !matchedGender {
			return false
		}
	}

	// Location validation
	if len(targeting.Locations) > 0 {
		matchedLoc := false
		for _, loc := range targeting.Locations {
			if strings.EqualFold(loc, audience.Location) {
				matchedLoc = true
				break
			}
		}
		if !matchedLoc {
			return false
		}
	}

	// Interests validation
	if len(targeting.Interests) > 0 {
		matchedInterest := false
		for _, requiredInterest := range targeting.Interests {
			for _, userInterest := range audience.Interests {
				if strings.EqualFold(requiredInterest, userInterest) {
					matchedInterest = true
					break
				}
			}
			if matchedInterest {
				break
			}
		}
		if !matchedInterest {
			return false
		}
	}

	return true
}

// 93. Real-Time Ad Auction & Bidding Engine (eCPM Vickrey Second-Price)
func (s *AdvertisingService) RunAuction(candidates []domain.AuctionCandidate, audience domain.UserAudienceProfile) (*domain.AuctionResult, error) {
	var eligible []domain.AuctionCandidate
	for _, c := range candidates {
		if s.MatchesTargeting(audience, c.AdSet.Targeting) {
			// eCPM formula: (Bid * estimatedCTR * 1000) + (QualityScore * 50)
			c.ComputedECPM = (float64(c.BidAmount) * c.EstimatedCTR * 10.0) + (c.QualityScore * 50.0)
			eligible = append(eligible, c)
		}
	}

	if len(eligible) == 0 {
		return nil, errors.New("no eligible ads match target audience")
	}

	// Sort eligible by ComputedECPM descending
	sort.Slice(eligible, func(i, j int) bool {
		return eligible[i].ComputedECPM > eligible[j].ComputedECPM
	})

	winner := eligible[0]
	priceCharged := winner.BidAmount
	if len(eligible) > 1 {
		// Vickrey second price: clearing price is competitor's bid + 1 cent
		secondBid := eligible[1].BidAmount
		if secondBid < winner.BidAmount {
			priceCharged = secondBid + 1
		}
	}

	return &domain.AuctionResult{
		WinningAd:    winner.Ad,
		WinningBid:   winner.BidAmount,
		PriceCharged: priceCharged,
		ECPM:         winner.ComputedECPM,
	}, nil
}

// 94. Ad Metrics Tracking
func (s *AdvertisingService) TrackImpression(adID string, priceCharged int64) error {
	_ = s.repo.RecordImpression(adID)
	return s.repo.RecordSpend(adID, priceCharged)
}

func (s *AdvertisingService) TrackClick(adID string) error {
	return s.repo.RecordClick(adID)
}

func (s *AdvertisingService) GetMetrics(adID string) (*domain.AdMetrics, error) {
	return s.repo.GetMetrics(adID)
}

// 95. Public Ad Transparency Library
func (s *AdvertisingService) SearchAdLibrary(query string) ([]domain.AdLibraryEntry, error) {
	return s.repo.SearchAdLibrary(query)
}

// 96. Meta Business Suite Partner Delegation
func (s *AdvertisingService) DelegatePartnerRole(businessID, partnerID, role string) (*domain.BusinessPartnerRole, error) {
	pRole := &domain.BusinessPartnerRole{
		BusinessID: businessID,
		PartnerID:  partnerID,
		Role:       role,
		AssignedAt: time.Now(),
	}
	if err := s.repo.SavePartnerRole(pRole); err != nil {
		return nil, err
	}
	return pRole, nil
}

func (s *AdvertisingService) GetPartnerRoles(businessID string) ([]domain.BusinessPartnerRole, error) {
	return s.repo.GetPartnerRoles(businessID)
}

