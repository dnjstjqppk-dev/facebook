package tests

import (
	"context"
	"testing"

	"github.com/dnjstjqppk-dev/facebook/pkg/domain"
	"github.com/dnjstjqppk-dev/facebook/services/groups/internal/repository"
	"github.com/dnjstjqppk-dev/facebook/services/groups/internal/service"
)

func setupTestGroupService() (*service.GroupService, repository.GroupRepository) {
	repo := repository.NewInMemoryGroupRepository()
	svc := service.NewGroupService(repo)
	return svc, repo
}

func TestCreateGroup_PublicAndPrivate(t *testing.T) {
	svc, _ := setupTestGroupService()
	ctx := context.Background()

	// 1. Public group
	pubGroup, err := svc.CreateGroup(ctx, "admin-1", "Go Developers Community", "All about Go", domain.GroupPrivacyPublic, false)
	if err != nil {
		t.Fatalf("failed to create public group: %v", err)
	}
	if pubGroup.Privacy != domain.GroupPrivacyPublic || pubGroup.Members["admin-1"] != domain.GroupRoleAdmin {
		t.Errorf("public group mismatch: %v", pubGroup)
	}

	// 2. Private group with post approval
	privGroup, err := svc.CreateGroup(ctx, "admin-2", "Secret Project Beta", "Internal group", domain.GroupPrivacyPrivate, true)
	if err != nil {
		t.Fatalf("failed to create private group: %v", err)
	}
	if privGroup.Privacy != domain.GroupPrivacyPrivate || !privGroup.RequirePostApproval {
		t.Errorf("private group settings mismatch: %v", privGroup)
	}
}

func TestJoinPrivateGroup_PendingApproval(t *testing.T) {
	svc, repo := setupTestGroupService()
	ctx := context.Background()

	group, _ := svc.CreateGroup(ctx, "admin-1", "Exclusive Club", "Members only", domain.GroupPrivacyPrivate, false)

	// User requests to join
	joined, err := svc.JoinGroup(ctx, "user-applicant", group.ID)
	if err != nil {
		t.Fatalf("join group error: %v", err)
	}
	if joined {
		t.Errorf("private group join should be queued as pending, not immediately joined")
	}

	g, _ := repo.GetGroup(ctx, group.ID)
	if len(g.PendingMembers) != 1 || g.PendingMembers[0] != "user-applicant" {
		t.Errorf("expected user in pending members list")
	}

	// Admin approves member
	err = svc.ApproveMember(ctx, "admin-1", group.ID, "user-applicant")
	if err != nil {
		t.Fatalf("failed to approve member: %v", err)
	}

	g, _ = repo.GetGroup(ctx, group.ID)
	if len(g.PendingMembers) != 0 || g.Members["user-applicant"] != domain.GroupRoleMember {
		t.Errorf("expected user to be moved from pending to approved members")
	}
}

func TestPostModerationQueue_PendingApproval(t *testing.T) {
	svc, repo := setupTestGroupService()
	ctx := context.Background()

	group, _ := svc.CreateGroup(ctx, "admin-1", "Tech Forum", "Moderated forum", domain.GroupPrivacyPublic, true)
	_, _ = svc.JoinGroup(ctx, "member-1", group.ID)

	// Normal member submits post
	post, err := svc.CreateGroupPost(ctx, "member-1", group.ID, "Is Go better than Rust?")
	if err != nil {
		t.Fatalf("failed to submit group post: %v", err)
	}
	if post.Status != domain.GroupPostPending {
		t.Errorf("expected post status PENDING_APPROVAL, got %s", post.Status)
	}

	// Admin approves post
	err = svc.ModeratePost(ctx, "admin-1", group.ID, post.ID, true)
	if err != nil {
		t.Fatalf("admin moderation failed: %v", err)
	}

	p, _ := repo.GetGroupPost(ctx, post.ID)
	if p.Status != domain.GroupPostApproved {
		t.Errorf("expected post to be APPROVED, got: %s", p.Status)
	}
}

func TestGroupRules_AdminOnly(t *testing.T) {
	svc, repo := setupTestGroupService()
	ctx := context.Background()

	group, _ := svc.CreateGroup(ctx, "admin-1", "Photography Club", "Photo rules", domain.GroupPrivacyPublic, false)
	_, _ = svc.JoinGroup(ctx, "member-1", group.ID)

	rules := []domain.GroupRule{
		{Title: "Be Respectful", Description: "No spam or hate speech."},
		{Title: "Original Content Only", Description: "Post only your own photos."},
	}

	// Member attempts to set rules (should fail)
	err := svc.SetRules(ctx, "member-1", group.ID, rules)
	if err != service.ErrNotGroupAdmin {
		t.Errorf("expected ErrNotGroupAdmin, got: %v", err)
	}

	// Admin sets rules
	err = svc.SetRules(ctx, "admin-1", group.ID, rules)
	if err != nil {
		t.Fatalf("admin set rules failed: %v", err)
	}

	g, _ := repo.GetGroup(ctx, group.ID)
	if len(g.Rules) != 2 || g.Rules[0].Title != "Be Respectful" {
		t.Errorf("rules not saved correctly: %v", g.Rules)
	}
}
