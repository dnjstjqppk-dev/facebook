package service

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/dnjstjqppk-dev/facebook/pkg/domain"
	"github.com/dnjstjqppk-dev/facebook/services/admin/internal/repository"
)

type AdminService struct {
	repo repository.AdminRepository
}

func NewAdminService(repo repository.AdminRepository) *AdminService {
	return &AdminService{repo: repo}
}

// 98. Platform Super Admin Audit & Kill Switches
func (s *AdminService) RecordAudit(adminID, action, resource, ip, details string) error {
	log := &domain.AuditLog{
		ID:        fmt.Sprintf("audit_%d", time.Now().UnixNano()),
		AdminID:   adminID,
		Action:    action,
		Resource:  resource,
		IPAddress: ip,
		Details:   details,
		Timestamp: time.Now(),
	}
	return s.repo.SaveAuditLog(log)
}

func (s *AdminService) ListAuditLogs(limit int) ([]domain.AuditLog, error) {
	return s.repo.ListAuditLogs(limit)
}

func (s *AdminService) SetKillSwitch(adminID, featureKey string, enabled bool, reason string) (*domain.KillSwitch, error) {
	ks := &domain.KillSwitch{
		FeatureKey: featureKey,
		Enabled:    enabled,
		Reason:     reason,
		UpdatedBy:  adminID,
		UpdatedAt:  time.Now(),
	}
	if err := s.repo.SetKillSwitch(ks); err != nil {
		return nil, err
	}
	_ = s.RecordAudit(adminID, "UPDATE_KILL_SWITCH", featureKey, "127.0.0.1", fmt.Sprintf("Enabled: %v, Reason: %s", enabled, reason))
	return ks, nil
}

func (s *AdminService) CheckFeatureEnabled(featureKey string) (bool, error) {
	ks, err := s.repo.GetKillSwitch(featureKey)
	if err != nil {
		return false, err
	}
	return ks.Enabled, nil
}

// 99. Developer Center & Graph API Gateway Token Management
func (s *AdminService) RegisterDeveloperApp(ownerID, appName string, scopes []string, tier string) (*domain.DeveloperApp, error) {
	rawSecret := make([]byte, 16)
	_, _ = rand.Read(rawSecret)
	secret := hex.EncodeToString(rawSecret)

	rawID := make([]byte, 6)
	_, _ = rand.Read(rawID)
	idSuffix := hex.EncodeToString(rawID)

	app := &domain.DeveloperApp{
		AppID:     fmt.Sprintf("app_%d_%s", time.Now().UnixNano(), idSuffix),
		OwnerID:   ownerID,
		AppName:   appName,
		AppSecret: secret,
		Scopes:    scopes,
		RateTier:  tier,
		CreatedAt: time.Now(),
	}
	if err := s.repo.SaveDeveloperApp(app); err != nil {
		return nil, err
	}
	return app, nil
}

func (s *AdminService) ValidateDeveloperApp(appID, appSecret string) (*domain.DeveloperApp, error) {
	app, err := s.repo.GetDeveloperApp(appID)
	if err != nil {
		return nil, err
	}
	if app.AppSecret != appSecret {
		return nil, domain.ErrAppUnauthorized
	}
	return app, nil
}
