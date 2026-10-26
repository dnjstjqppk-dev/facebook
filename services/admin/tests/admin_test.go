package tests

import (
	"testing"

	"github.com/dnjstjqppk-dev/facebook/pkg/domain"
	"github.com/dnjstjqppk-dev/facebook/services/admin/internal/repository"
	"github.com/dnjstjqppk-dev/facebook/services/admin/internal/service"
)

func TestAdminAuditAndKillSwitch(t *testing.T) {
	repo := repository.NewInMemoryAdminRepo()
	svc := service.NewAdminService(repo)

	// Step 98: Super Admin Audit Log
	if err := svc.RecordAudit("super_admin_01", "USER_LOCKDOWN", "user_malicious_99", "10.0.0.1", "Suspended due to zero-day exploit attempt"); err != nil {
		t.Fatalf("failed to record audit: %v", err)
	}

	logs, err := svc.ListAuditLogs(10)
	if err != nil || len(logs) != 1 {
		t.Fatalf("expected 1 audit log, got %d (err: %v)", len(logs), err)
	}

	// Step 98: System Kill Switch
	enabled, err := svc.CheckFeatureEnabled("marketplace_payments")
	if err != nil || !enabled {
		t.Fatalf("expected default enabled, got %v", enabled)
	}

	// Trigger emergency kill switch
	ks, err := svc.SetKillSwitch("super_admin_01", "marketplace_payments", false, "Zero-day vulnerability mitigation in escrow ledger")
	if err != nil || ks.Enabled {
		t.Fatalf("expected kill switch disabled, got: %+v", ks)
	}

	enabledAfter, err := svc.CheckFeatureEnabled("marketplace_payments")
	if err != nil || enabledAfter {
		t.Fatalf("expected feature to be disabled after kill switch triggered")
	}

	// Check that kill switch action was audited
	logsAfter, _ := svc.ListAuditLogs(10)
	if len(logsAfter) != 2 {
		t.Errorf("expected 2 audit logs, got %d", len(logsAfter))
	}
}

func TestDeveloperCenter_GraphAPICredentials(t *testing.T) {
	repo := repository.NewInMemoryAdminRepo()
	svc := service.NewAdminService(repo)

	// Step 99: Developer App Registration
	app, err := svc.RegisterDeveloperApp("developer_dev_10", "Social Insights Analytics Bot", []string{"pages_read_engagement", "ads_read"}, "ENTERPRISE")
	if err != nil {
		t.Fatalf("failed to register app: %v", err)
	}

	if app.AppID == "" || app.AppSecret == "" {
		t.Fatalf("invalid app credentials: %+v", app)
	}

	// Validate valid credentials
	validApp, err := svc.ValidateDeveloperApp(app.AppID, app.AppSecret)
	if err != nil || validApp.AppID != app.AppID {
		t.Fatalf("failed to validate correct credentials: %v", err)
	}

	// Reject invalid credentials
	_, err = svc.ValidateDeveloperApp(app.AppID, "invalid_secret_payload")
	if err != domain.ErrAppUnauthorized {
		t.Errorf("expected ErrAppUnauthorized, got: %v", err)
	}
}
