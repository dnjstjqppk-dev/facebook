package repository

import (
	"sync"

	"github.com/dnjstjqppk-dev/facebook/pkg/domain"
)

type AdminRepository interface {
	SaveAuditLog(log *domain.AuditLog) error
	ListAuditLogs(limit int) ([]domain.AuditLog, error)
	SetKillSwitch(ks *domain.KillSwitch) error
	GetKillSwitch(featureKey string) (*domain.KillSwitch, error)
	SaveDeveloperApp(app *domain.DeveloperApp) error
	GetDeveloperApp(appID string) (*domain.DeveloperApp, error)
}

type InMemoryAdminRepo struct {
	mu           sync.RWMutex
	auditLogs    []domain.AuditLog
	killSwitches map[string]*domain.KillSwitch
	devApps      map[string]*domain.DeveloperApp
}

func NewInMemoryAdminRepo() *InMemoryAdminRepo {
	return &InMemoryAdminRepo{
		killSwitches: make(map[string]*domain.KillSwitch),
		devApps:      make(map[string]*domain.DeveloperApp),
	}
}

func (r *InMemoryAdminRepo) SaveAuditLog(log *domain.AuditLog) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.auditLogs = append(r.auditLogs, *log)
	return nil
}

func (r *InMemoryAdminRepo) ListAuditLogs(limit int) ([]domain.AuditLog, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if limit <= 0 || limit > len(r.auditLogs) {
		limit = len(r.auditLogs)
	}
	res := make([]domain.AuditLog, limit)
	copy(res, r.auditLogs[len(r.auditLogs)-limit:])
	return res, nil
}

func (r *InMemoryAdminRepo) SetKillSwitch(ks *domain.KillSwitch) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.killSwitches[ks.FeatureKey] = ks
	return nil
}

func (r *InMemoryAdminRepo) GetKillSwitch(featureKey string) (*domain.KillSwitch, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	ks, ok := r.killSwitches[featureKey]
	if !ok {
		// By default all features are enabled unless explicitly disabled
		return &domain.KillSwitch{
			FeatureKey: featureKey,
			Enabled:    true,
		}, nil
	}
	return ks, nil
}

func (r *InMemoryAdminRepo) SaveDeveloperApp(app *domain.DeveloperApp) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.devApps[app.AppID] = app
	return nil
}

func (r *InMemoryAdminRepo) GetDeveloperApp(appID string) (*domain.DeveloperApp, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	app, ok := r.devApps[appID]
	if !ok {
		return nil, domain.ErrAppUnauthorized
	}
	return app, nil
}

