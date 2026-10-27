package tests

import (
	"testing"

	"github.com/dnjstjqppk-dev/facebook/pkg/domain"
	"github.com/dnjstjqppk-dev/facebook/services/moderation/internal/repository"
	"github.com/dnjstjqppk-dev/facebook/services/moderation/internal/service"
)

func TestModerationReportLifecycle(t *testing.T) {
	repo := repository.NewInMemoryModerationRepo()
	svc := service.NewModerationService(repo)

	// Step 97: Submit content report
	rep, err := svc.SubmitReport("reporter_user_1", "POST", "post_spam_99", domain.ReportReasonSpam)
	if err != nil {
		t.Fatalf("failed to submit report: %v", err)
	}

	pending, err := svc.ListPendingReports()
	if err != nil || len(pending) != 1 {
		t.Fatalf("expected 1 pending report, got %d", len(pending))
	}

	// Reviewer takes action: remove content
	reviewed, err := svc.ReviewReport(rep.ID, "mod_admin_42", domain.ModerationActionContentRemoved, "Confirmed spam bot advertising crypto scam.")
	if err != nil {
		t.Fatalf("failed to review report: %v", err)
	}

	if reviewed.Status != domain.ReportStatusActionTaken || reviewed.ActionTaken != domain.ModerationActionContentRemoved {
		t.Errorf("unexpected reviewed state: %+v", reviewed)
	}
}

func TestModerationUserStrikeEscalation(t *testing.T) {
	repo := repository.NewInMemoryModerationRepo()
	svc := service.NewModerationService(repo)

	userID := "bad_actor_007"

	// 1st strike -> WARNING
	s1, err := svc.ApplyStrike(userID)
	if err != nil || s1.CurrentTier != "WARNING" || s1.ViolationCount != 1 {
		t.Fatalf("expected strike 1 WARNING, got: %+v (err: %v)", s1, err)
	}

	// 2nd strike -> TEMP_MUTE
	s2, err := svc.ApplyStrike(userID)
	if err != nil || s2.CurrentTier != "TEMP_MUTE" || s2.ViolationCount != 2 {
		t.Fatalf("expected strike 2 TEMP_MUTE, got: %+v", s2)
	}

	// 3rd strike -> SUSPENSION
	s3, err := svc.ApplyStrike(userID)
	if err != nil || s3.CurrentTier != "SUSPENSION" || s3.ViolationCount != 3 {
		t.Fatalf("expected strike 3 SUSPENSION, got: %+v", s3)
	}

	// 4th strike -> PERMANENT_BAN
	s4, err := svc.ApplyStrike(userID)
	if err != nil || s4.CurrentTier != "PERMANENT_BAN" || s4.ViolationCount != 4 {
		t.Fatalf("expected strike 4 PERMANENT_BAN, got: %+v", s4)
	}
}

