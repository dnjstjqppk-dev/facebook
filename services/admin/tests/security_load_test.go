package tests

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/dnjstjqppk-dev/facebook/pkg/auth"
	"github.com/dnjstjqppk-dev/facebook/pkg/domain"
	"github.com/dnjstjqppk-dev/facebook/services/admin/internal/repository"
	"github.com/dnjstjqppk-dev/facebook/services/admin/internal/service"
)

// Step 100: End-to-End Security Penetration & Load Verification
func TestSecurityPenetration_TokenTamperingAndAuthBoundary(t *testing.T) {
	jwtMgr := auth.NewJWTManager("production-ultra-secure-key-998877", 1*time.Second, 1*time.Minute)

	user := &domain.User{
		ID:    "user_legit_01",
		Email: "user@facebook.local",
		Role:  domain.RoleUserStandard,
	}

	// 1. Token Tampering Defense
	tokens, err := jwtMgr.GenerateTokens(user)
	if err != nil {
		t.Fatalf("failed to generate tokens: %v", err)
	}

	tamperedToken := tokens.AccessToken + "malicious_bits"
	_, err = jwtMgr.VerifyToken(tamperedToken)
	if err == nil {
		t.Fatalf("security violation: tampered token was accepted!")
	}

	// 2. Token Expiration Defense
	shortLivedMgr := auth.NewJWTManager("short-key", -1*time.Second, -1*time.Second)
	expiredTokens, _ := shortLivedMgr.GenerateTokens(user)
	_, err = shortLivedMgr.VerifyToken(expiredTokens.AccessToken)
	if err != auth.ErrExpiredToken {
		t.Fatalf("security violation: expired token was accepted! err: %v", err)
	}

	// 3. Admin Kill Switch Boundary Enforcement
	repo := repository.NewInMemoryAdminRepo()
	svc := service.NewAdminService(repo)

	// Set kill switch
	_, err = svc.SetKillSwitch("admin_007", "live_streaming", false, "DDoS incident mitigation")
	if err != nil {
		t.Fatalf("failed to set kill switch: %v", err)
	}

	enabled, err := svc.CheckFeatureEnabled("live_streaming")
	if err != nil || enabled {
		t.Errorf("expected live_streaming to be disabled, got %v", enabled)
	}
}

func TestConcurrentLoadSimulation_AdminAuditAndTokens(t *testing.T) {
	repo := repository.NewInMemoryAdminRepo()
	svc := service.NewAdminService(repo)

	var wg sync.WaitGroup
	var successfulOps int64
	numWorkers := 50
	opsPerWorker := 30

	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for i := 0; i < opsPerWorker; i++ {
				// 1. Record Audit Log concurrently
				err := svc.RecordAudit("admin_worker", "INSPECT_TRAFFIC", "gateway_edge", "10.0.1.20", "Load spike observation")
				if err == nil {
					atomic.AddInt64(&successfulOps, 1)
				}

				// 2. Developer App token verification concurrently
				app, err := svc.RegisterDeveloperApp("dev_owner", "Load Tester App", []string{"read_public"}, "STANDARD")
				if err == nil && app.AppID != "" {
					val, vErr := svc.ValidateDeveloperApp(app.AppID, app.AppSecret)
					if vErr == nil && val.AppID == app.AppID {
						atomic.AddInt64(&successfulOps, 1)
					}
				}
			}
		}(w)
	}

	wg.Wait()

	expectedOps := int64(numWorkers * opsPerWorker * 2)
	if successfulOps != expectedOps {
		t.Fatalf("expected %d ops under concurrent load, got %d", expectedOps, successfulOps)
	}

	logs, err := svc.ListAuditLogs(100)
	if err != nil || len(logs) == 0 {
		t.Fatalf("failed to list audit logs after high concurrency test")
	}
}
