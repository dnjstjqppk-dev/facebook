package repository

import (
	"sync"

	"github.com/dnjstjqppk-dev/facebook/pkg/domain"
)

type ModerationRepository interface {
	SaveReport(rep *domain.ModerationReport) error
	GetReport(id string) (*domain.ModerationReport, error)
	ListPendingReports() ([]domain.ModerationReport, error)
	SaveUserStrike(strike *domain.UserStrike) error
	GetUserStrike(userID string) (*domain.UserStrike, error)
}

type InMemoryModerationRepo struct {
	mu      sync.RWMutex
	reports map[string]*domain.ModerationReport
	strikes map[string]*domain.UserStrike
}

func NewInMemoryModerationRepo() *InMemoryModerationRepo {
	return &InMemoryModerationRepo{
		reports: make(map[string]*domain.ModerationReport),
		strikes: make(map[string]*domain.UserStrike),
	}
}

func (r *InMemoryModerationRepo) SaveReport(rep *domain.ModerationReport) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.reports[rep.ID] = rep
	return nil
}

func (r *InMemoryModerationRepo) GetReport(id string) (*domain.ModerationReport, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	rep, ok := r.reports[id]
	if !ok {
		return nil, domain.ErrReportNotFound
	}
	return rep, nil
}

func (r *InMemoryModerationRepo) ListPendingReports() ([]domain.ModerationReport, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var res []domain.ModerationReport
	for _, rep := range r.reports {
		if rep.Status == domain.ReportStatusPending {
			res = append(res, *rep)
		}
	}
	return res, nil
}

func (r *InMemoryModerationRepo) SaveUserStrike(strike *domain.UserStrike) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.strikes[strike.UserID] = strike
	return nil
}

func (r *InMemoryModerationRepo) GetUserStrike(userID string) (*domain.UserStrike, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	strike, ok := r.strikes[userID]
	if !ok {
		return &domain.UserStrike{
			UserID:         userID,
			ViolationCount: 0,
			CurrentTier:    "CLEAN",
		}, nil
	}
	return strike, nil
}

