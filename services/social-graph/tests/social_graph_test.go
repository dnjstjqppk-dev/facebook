package tests

import (
	"context"
	"testing"

	"github.com/dnjstjqppk-dev/facebook/pkg/domain"
	"github.com/dnjstjqppk-dev/facebook/services/social-graph/internal/repository"
	"github.com/dnjstjqppk-dev/facebook/services/social-graph/internal/service"
)

func setupTestSocialService() *service.SocialGraphService {
	repo := repository.NewInMemorySocialGraphRepository()
	return service.NewSocialGraphService(repo)
}

func TestSendFriendRequest(t *testing.T) {
	svc := setupTestSocialService()
	ctx := context.Background()

	aliceID := "user-alice-1"
	bobID := "user-bob-2"

	rel, err := svc.SendFriendRequest(ctx, aliceID, bobID)
	if err != nil {
		t.Fatalf("expected friend request to succeed, got: %v", err)
	}

	if rel.Status != domain.RelStatusPending {
		t.Errorf("expected status PENDING, got: %s", rel.Status)
	}

	// Verify pending requests for Bob
	pending, err := svc.GetPendingRequests(ctx, bobID)
	if err != nil {
		t.Fatalf("failed to get pending requests: %v", err)
	}
	if len(pending) != 1 || pending[0].UserID != aliceID {
		t.Errorf("expected 1 pending request from Alice, got: %d", len(pending))
	}
}

func TestSelfFriendshipRejection(t *testing.T) {
	svc := setupTestSocialService()
	ctx := context.Background()

	aliceID := "user-alice-1"
	_, err := svc.SendFriendRequest(ctx, aliceID, aliceID)
	if err != service.ErrSelfFriendship {
		t.Fatalf("expected ErrSelfFriendship, got: %v", err)
	}
}

func TestDuplicateFriendRequest(t *testing.T) {
	svc := setupTestSocialService()
	ctx := context.Background()

	aliceID := "user-alice-1"
	bobID := "user-bob-2"

	_, err := svc.SendFriendRequest(ctx, aliceID, bobID)
	if err != nil {
		t.Fatalf("first request failed: %v", err)
	}

	_, err = svc.SendFriendRequest(ctx, aliceID, bobID)
	if err != service.ErrPendingRequestExists {
		t.Fatalf("expected ErrPendingRequestExists, got: %v", err)
	}
}

func TestAcceptFriendRequestBidirectional(t *testing.T) {
	svc := setupTestSocialService()
	ctx := context.Background()

	aliceID := "user-alice-1"
	bobID := "user-bob-2"

	// Alice sends request to Bob
	_, err := svc.SendFriendRequest(ctx, aliceID, bobID)
	if err != nil {
		t.Fatalf("send request failed: %v", err)
	}

	// Bob accepts Alice's request
	accepted, err := svc.AcceptFriendRequest(ctx, bobID, aliceID)
	if err != nil {
		t.Fatalf("accept request failed: %v", err)
	}
	if accepted.Status != domain.RelStatusAccepted {
		t.Errorf("expected accepted status, got: %s", accepted.Status)
	}

	// Verify Alice sees Bob as friend
	aliceFriends, err := svc.GetFriends(ctx, aliceID)
	if err != nil || len(aliceFriends) != 1 {
		t.Errorf("expected Alice to have 1 friend, got: %d", len(aliceFriends))
	}

	// Verify Bob sees Alice as friend
	bobFriends, err := svc.GetFriends(ctx, bobID)
	if err != nil || len(bobFriends) != 1 {
		t.Errorf("expected Bob to have 1 friend, got: %d", len(bobFriends))
	}

	// Verify pending request list is now empty
	pending, _ := svc.GetPendingRequests(ctx, bobID)
	if len(pending) != 0 {
		t.Errorf("expected 0 pending requests, got: %d", len(pending))
	}
}

func TestUnfriend(t *testing.T) {
	svc := setupTestSocialService()
	ctx := context.Background()

	aliceID := "user-alice-1"
	bobID := "user-bob-2"

	_, _ = svc.SendFriendRequest(ctx, aliceID, bobID)
	_, _ = svc.AcceptFriendRequest(ctx, bobID, aliceID)

	// Alice unfriends Bob
	err := svc.Unfriend(ctx, aliceID, bobID)
	if err != nil {
		t.Fatalf("unfriend failed: %v", err)
	}

	aliceFriends, _ := svc.GetFriends(ctx, aliceID)
	bobFriends, _ := svc.GetFriends(ctx, bobID)
	if len(aliceFriends) != 0 || len(bobFriends) != 0 {
		t.Errorf("expected both users to have 0 friends after unfriend")
	}
}

func TestBlockUserAndDenyFriendRequest(t *testing.T) {
	svc := setupTestSocialService()
	ctx := context.Background()

	aliceID := "user-alice-1"
	bobID := "user-bob-2"

	// Alice blocks Bob
	err := svc.BlockUser(ctx, aliceID, bobID)
	if err != nil {
		t.Fatalf("block user failed: %v", err)
	}

	// Bob attempts to send friend request to Alice
	_, err = svc.SendFriendRequest(ctx, bobID, aliceID)
	if err != service.ErrBlockedUser {
		t.Fatalf("expected ErrBlockedUser, got: %v", err)
	}
}
