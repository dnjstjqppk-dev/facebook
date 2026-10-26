package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/dnjstjqppk-dev/facebook/pkg/domain"
	"github.com/dnjstjqppk-dev/facebook/services/moderation/internal/repository"
)

type ModerationService struct {
	repo repository.ModerationRepository
}

func NewModerationService(repo repository.ModerationRepository) *ModerationService {
	return &ModerationService{repo: repo}
}

// 97. Submit Content Report
func (s *ModerationService) SubmitReport(reporterID, targetType, targetID, reason string) (*domain.ModerationReport, error) {
	if targetID == "" || reporterID == "" {
		return nil, errors.New("target and reporter are required")
	}
	rep := &domain.ModerationReport{
		ID:          fmt.Sprintf("rep_%d", time.Now().UnixNano()),
		ReporterID:  reporterID,
		TargetType:  targetType,
		TargetID:    targetID,
		Reason:      reason,
		Status:      domain.ReportStatusPending,
		ActionTaken: domain.ModerationActionNone,
		CreatedAt:   time.Now(),
	}
	if err := s.repo.SaveReport(rep); err != nil {
		return nil, err
	}
	return rep, nil
}

func (s *ModerationService) ListPendingReports() ([]domain.ModerationReport, error) {
	return s.repo.ListPendingReports()
}

func (s *ModerationService) ReviewReport(reportID, reviewerID, action, notes string) (*domain.ModerationReport, error) {
	rep, err := s.repo.GetReport(reportID)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	rep.ReviewerID = reviewerID
	rep.ActionTaken = action
	rep.Notes = notes
	rep.ResolvedAt = &now

	if action == domain.ModerationActionNone {
		rep.Status = domain.ReportStatusDismissed
	} else {
		rep.Status = domain.ReportStatusActionTaken
	}

	if err := s.repo.SaveReport(rep); err != nil {
		return nil, err
	}
	return rep, nil
}

// User Strikes Escalation Policy
func (s *ModerationService) ApplyStrike(userID string) (*domain.UserStrike, error) {
	strike, err := s.repo.GetUserStrike(userID)
	if err != nil {
		return nil, err
	}

	strike.ViolationCount++
	now := time.Now()

	switch strike.ViolationCount {
	case 1:
		strike.CurrentTier = "WARNING"
		strike.ExpiresAt = now.Add(30 * 24 * time.Hour)
	case 2:
		strike.CurrentTier = "TEMP_MUTE"
		strike.ExpiresAt = now.Add(24 * time.Hour)
	case 3:
		strike.CurrentTier = "SUSPENSION"
		strike.ExpiresAt = now.Add(7 * 24 * time.Hour)
	default:
		strike.CurrentTier = "PERMANENT_BAN"
		strike.ExpiresAt = now.Add(100 * 365 * 24 * time.Hour)
	}

	if err := s.repo.SaveUserStrike(strike); err != nil {
		return nil, err
	}
	return strike, nil
}

func (s *ModerationService) GetUserStrike(userID string) (*domain.UserStrike, error) {
	return s.repo.GetUserStrike(userID)
}
