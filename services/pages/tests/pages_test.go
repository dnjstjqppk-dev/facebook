package tests

import (
	"context"
	"testing"

	"github.com/dnjstjqppk-dev/facebook/pkg/domain"
	"github.com/dnjstjqppk-dev/facebook/services/pages/internal/repository"
	"github.com/dnjstjqppk-dev/facebook/services/pages/internal/service"
)

func setupTestPageService() (*service.PageService, repository.PageRepository) {
	repo := repository.NewInMemoryPageRepository()
	svc := service.NewPageService(repo)
	return svc, repo
}

func TestCreatePage_WithAdmin(t *testing.T) {
	svc, _ := setupTestPageService()
	ctx := context.Background()

	page, err := svc.CreatePage(ctx, "brand-owner", "Acme Corp Official", "Business/Technology", "Innovating the future")
	if err != nil {
		t.Fatalf("failed to create page: %v", err)
	}

	if page.Name != "Acme Corp Official" || page.Roles["brand-owner"] != domain.PageRoleAdmin {
		t.Errorf("page attributes mismatch: %v", page)
	}
}

func TestAssignRole_EditorAndModerator(t *testing.T) {
	svc, repo := setupTestPageService()
	ctx := context.Background()

	page, _ := svc.CreatePage(ctx, "admin-1", "Fashion Brand", "Retail", "Trending fashion")

	// Admin assigns Editor to user-2
	err := svc.AssignRole(ctx, "admin-1", page.ID, "user-editor", domain.PageRoleEditor)
	if err != nil {
		t.Fatalf("assign role failed: %v", err)
	}

	p, _ := repo.GetPage(ctx, page.ID)
	if p.Roles["user-editor"] != domain.PageRoleEditor {
		t.Errorf("expected user-editor role EDITOR, got: %s", p.Roles["user-editor"])
	}

	// Non-admin attempting to assign role fails
	err = svc.AssignRole(ctx, "user-editor", page.ID, "user-3", domain.PageRoleAdmin)
	if err != service.ErrNotPageAdmin {
		t.Errorf("expected ErrNotPageAdmin, got: %v", err)
	}
}

func TestPageFollowers_Increment(t *testing.T) {
	svc, repo := setupTestPageService()
	ctx := context.Background()

	page, _ := svc.CreatePage(ctx, "creator-1", "Gaming News", "Media", "Daily gaming updates")

	_ = svc.FollowPage(ctx, "fan-1", page.ID)
	_ = svc.FollowPage(ctx, "fan-2", page.ID)

	p, _ := repo.GetPage(ctx, page.ID)
	if p.FollowersCount != 2 {
		t.Errorf("expected 2 followers, got %d", p.FollowersCount)
	}
}

func TestPageInsights_AccessControl(t *testing.T) {
	svc, _ := setupTestPageService()
	ctx := context.Background()

	page, _ := svc.CreatePage(ctx, "business-admin", "Mega Store", "Retail", "Everything for you")
	_ = svc.FollowPage(ctx, "customer-1", page.ID)

	// Admin views insights
	insights, err := svc.GetPageInsights(ctx, "business-admin", page.ID)
	if err != nil || insights.FollowersCount != 1 {
		t.Fatalf("expected insights retrieval to succeed: %v", err)
	}

	// Unauthorized customer attempts to view private page insights
	_, err = svc.GetPageInsights(ctx, "customer-1", page.ID)
	if err != service.ErrNotPageAdmin {
		t.Errorf("expected ErrNotPageAdmin for outsider, got: %v", err)
	}
}
