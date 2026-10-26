package repository

import (
	"sync"

	"github.com/dnjstjqppk-dev/facebook/pkg/domain"
)

type AdRepository interface {
	SaveAccount(acc *domain.AdAccount) error
	GetAccount(id string) (*domain.AdAccount, error)
	SaveCampaign(cmp *domain.AdCampaign) error
	GetCampaign(id string) (*domain.AdCampaign, error)
	SaveAdSet(adSet *domain.AdSet) error
	GetAdSet(id string) (*domain.AdSet, error)
	ListAdSetsByCampaign(cmpID string) ([]domain.AdSet, error)
	SaveAd(ad *domain.Ad) error
	GetAd(id string) (*domain.Ad, error)
	ListActiveAds() ([]domain.Ad, error)
	RecordImpression(adID string) error
	RecordClick(adID string) error
	RecordSpend(adID string, amount int64) error
	GetMetrics(adID string) (*domain.AdMetrics, error)
	SaveAdLibraryEntry(entry *domain.AdLibraryEntry) error
	SearchAdLibrary(query string) ([]domain.AdLibraryEntry, error)
	SavePartnerRole(role *domain.BusinessPartnerRole) error
	GetPartnerRoles(businessID string) ([]domain.BusinessPartnerRole, error)
}

type InMemoryAdRepo struct {
	mu           sync.RWMutex
	accounts     map[string]*domain.AdAccount
	campaigns    map[string]*domain.AdCampaign
	adSets       map[string]*domain.AdSet
	ads          map[string]*domain.Ad
	metrics      map[string]*domain.AdMetrics
	adLibrary    map[string]*domain.AdLibraryEntry
	partnerRoles map[string][]domain.BusinessPartnerRole
}

func NewInMemoryAdRepo() *InMemoryAdRepo {
	return &InMemoryAdRepo{
		accounts:     make(map[string]*domain.AdAccount),
		campaigns:    make(map[string]*domain.AdCampaign),
		adSets:       make(map[string]*domain.AdSet),
		ads:          make(map[string]*domain.Ad),
		metrics:      make(map[string]*domain.AdMetrics),
		adLibrary:    make(map[string]*domain.AdLibraryEntry),
		partnerRoles: make(map[string][]domain.BusinessPartnerRole),
	}
}

func (r *InMemoryAdRepo) SaveAccount(acc *domain.AdAccount) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.accounts[acc.ID] = acc
	return nil
}

func (r *InMemoryAdRepo) GetAccount(id string) (*domain.AdAccount, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	acc, ok := r.accounts[id]
	if !ok {
		return nil, domain.ErrAdAccountNotFound
	}
	return acc, nil
}

func (r *InMemoryAdRepo) SaveCampaign(cmp *domain.AdCampaign) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.campaigns[cmp.ID] = cmp
	return nil
}

func (r *InMemoryAdRepo) GetCampaign(id string) (*domain.AdCampaign, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	cmp, ok := r.campaigns[id]
	if !ok {
		return nil, domain.ErrAdAccountNotFound
	}
	return cmp, nil
}

func (r *InMemoryAdRepo) SaveAdSet(adSet *domain.AdSet) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.adSets[adSet.ID] = adSet
	return nil
}

func (r *InMemoryAdRepo) GetAdSet(id string) (*domain.AdSet, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	adSet, ok := r.adSets[id]
	if !ok {
		return nil, domain.ErrAdAccountNotFound
	}
	return adSet, nil
}

func (r *InMemoryAdRepo) ListAdSetsByCampaign(cmpID string) ([]domain.AdSet, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var res []domain.AdSet
	for _, as := range r.adSets {
		if as.CampaignID == cmpID {
			res = append(res, *as)
		}
	}
	return res, nil
}

func (r *InMemoryAdRepo) SaveAd(ad *domain.Ad) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.ads[ad.ID] = ad
	return nil
}

func (r *InMemoryAdRepo) GetAd(id string) (*domain.Ad, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	ad, ok := r.ads[id]
	if !ok {
		return nil, domain.ErrAdAccountNotFound
	}
	return ad, nil
}

func (r *InMemoryAdRepo) ListActiveAds() ([]domain.Ad, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var res []domain.Ad
	for _, ad := range r.ads {
		if ad.Status == domain.AdStatusActive {
			res = append(res, *ad)
		}
	}
	return res, nil
}

func (r *InMemoryAdRepo) RecordImpression(adID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	m, ok := r.metrics[adID]
	if !ok {
		m = &domain.AdMetrics{AdID: adID}
		r.metrics[adID] = m
	}
	m.Impressions++
	return nil
}

func (r *InMemoryAdRepo) RecordClick(adID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	m, ok := r.metrics[adID]
	if !ok {
		m = &domain.AdMetrics{AdID: adID}
		r.metrics[adID] = m
	}
	m.Clicks++
	return nil
}

func (r *InMemoryAdRepo) RecordSpend(adID string, amount int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	m, ok := r.metrics[adID]
	if !ok {
		m = &domain.AdMetrics{AdID: adID}
		r.metrics[adID] = m
	}
	m.Spend += amount
	return nil
}

func (r *InMemoryAdRepo) GetMetrics(adID string) (*domain.AdMetrics, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	m, ok := r.metrics[adID]
	if !ok {
		return &domain.AdMetrics{AdID: adID}, nil
	}
	return m, nil
}

func (r *InMemoryAdRepo) SaveAdLibraryEntry(entry *domain.AdLibraryEntry) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.adLibrary[entry.AdID] = entry
	return nil
}

func (r *InMemoryAdRepo) SearchAdLibrary(query string) ([]domain.AdLibraryEntry, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var res []domain.AdLibraryEntry
	for _, e := range r.adLibrary {
		res = append(res, *e)
	}
	return res, nil
}

func (r *InMemoryAdRepo) SavePartnerRole(role *domain.BusinessPartnerRole) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.partnerRoles[role.BusinessID] = append(r.partnerRoles[role.BusinessID], *role)
	return nil
}

func (r *InMemoryAdRepo) GetPartnerRoles(businessID string) ([]domain.BusinessPartnerRole, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.partnerRoles[businessID], nil
}
