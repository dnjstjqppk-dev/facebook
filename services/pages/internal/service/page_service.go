package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/dnjstjqppk-dev/facebook/pkg/domain"
	"github.com/dnjstjqppk-dev/facebook/services/pages/internal/repository"
)

var (
	ErrNotPageAdmin = errors.New("operation requires page administrator role")
)

type PageService struct {
	repo repository.PageRepository
}

func NewPageService(repo repository.PageRepository) *PageService {
	return &PageService{repo: repo}
}

func (s *PageService) CreatePage(ctx context.Context, creatorID, name, category, description string) (*domain.Page, error) {
	now := time.Now().UTC()
	page := &domain.Page{
		ID:             generateUUID(),
		Name:           name,
		Category:       category,
		Description:    description,
		CreatedBy:      creatorID,
		Roles:          map[string]domain.PageRole{creatorID: domain.PageRoleAdmin},
		FollowersCount: 0,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	if err := s.repo.SavePage(ctx, page); err != nil {
		return nil, err
	}
	return page, nil
}

func (s *PageService) AssignRole(ctx context.Context, adminID, pageID, targetUserID string, role domain.PageRole) error {
	page, err := s.repo.GetPage(ctx, pageID)
	if err != nil {
		return err
	}

	if page.Roles[adminID] != domain.PageRoleAdmin {
		return ErrNotPageAdmin
	}

	page.Roles[targetUserID] = role
	page.UpdatedAt = time.Now().UTC()
	return s.repo.SavePage(ctx, page)
}

func (s *PageService) FollowPage(ctx context.Context, userID, pageID string) error {
	return s.repo.IncrementFollowers(ctx, pageID)
}

func (s *PageService) GetPageInsights(ctx context.Context, requesterID, pageID string) (*domain.PageInsights, error) {
	page, err := s.repo.GetPage(ctx, pageID)
	if err != nil {
		return nil, err
	}

	// Only assigned page team can view insights
	if _, hasRole := page.Roles[requesterID]; !hasRole {
		return nil, ErrNotPageAdmin
	}

	return s.repo.GetPageInsights(ctx, pageID)
}

func generateUUID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return hex.EncodeToString(b[0:4]) + "-" +
		hex.EncodeToString(b[4:6]) + "-" +
		hex.EncodeToString(b[6:8]) + "-" +
		hex.EncodeToString(b[8:10]) + "-" +
		hex.EncodeToString(b[10:16])
}
