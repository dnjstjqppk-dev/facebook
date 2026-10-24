package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/dnjstjqppk-dev/facebook/pkg/domain"
	"github.com/dnjstjqppk-dev/facebook/services/social-graph/internal/repository"
)

var (
	ErrSelfFriendship       = errors.New("cannot initiate friendship with yourself")
	ErrBlockedUser          = errors.New("action not allowed: user is blocked or has blocked you")
	ErrAlreadyFriends       = errors.New("users are already friends")
	ErrPendingRequestExists = errors.New("a friend request is already pending")
	ErrRequestNotFound      = errors.New("pending friend request not found")
)

type SocialGraphService struct {
	repo repository.SocialGraphRepository
}

func NewSocialGraphService(repo repository.SocialGraphRepository) *SocialGraphService {
	return &SocialGraphService{repo: repo}
}

func (s *SocialGraphService) SendFriendRequest(ctx context.Context, requesterID, targetID string) (*domain.Relationship, error) {
	if requesterID == targetID {
		return nil, ErrSelfFriendship
	}

	blocked, err := s.repo.IsBlocked(ctx, requesterID, targetID)
	if err != nil {
		return nil, err
	}
	if blocked {
		return nil, ErrBlockedUser
	}

	// Check existing friendship
	existing, err := s.repo.GetRelationship(ctx, requesterID, targetID, domain.RelFriend)
	if err == nil {
		if existing.Status == domain.RelStatusAccepted {
			return nil, ErrAlreadyFriends
		}
		if existing.Status == domain.RelStatusPending {
			return nil, ErrPendingRequestExists
		}
	}

	// Check if target already sent a request to requester (mutual collision -> auto accept)
	reverseReq, err := s.repo.GetRelationship(ctx, targetID, requesterID, domain.RelFriend)
	if err == nil && reverseReq.Status == domain.RelStatusPending {
		return s.AcceptFriendRequest(ctx, requesterID, targetID)
	}

	now := time.Now().UTC()
	rel := &domain.Relationship{
		ID:               generateUUID(),
		UserID:           requesterID,
		TargetUserID:     targetID,
		RelationshipType: domain.RelFriend,
		Status:           domain.RelStatusPending,
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	if err := s.repo.SaveRelationship(ctx, rel); err != nil {
		return nil, err
	}

	return rel, nil
}

func (s *SocialGraphService) AcceptFriendRequest(ctx context.Context, recipientID, senderID string) (*domain.Relationship, error) {
	req, err := s.repo.GetRelationship(ctx, senderID, recipientID, domain.RelFriend)
	if err != nil || req.Status != domain.RelStatusPending {
		return nil, ErrRequestNotFound
	}

	now := time.Now().UTC()
	req.Status = domain.RelStatusAccepted
	req.UpdatedAt = now

	if err := s.repo.SaveRelationship(ctx, req); err != nil {
		return nil, err
	}

	// Establish bidirectional relationship
	reciprocalRel := &domain.Relationship{
		ID:               generateUUID(),
		UserID:           recipientID,
		TargetUserID:     senderID,
		RelationshipType: domain.RelFriend,
		Status:           domain.RelStatusAccepted,
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	if err := s.repo.SaveRelationship(ctx, reciprocalRel); err != nil {
		return nil, err
	}

	return req, nil
}

func (s *SocialGraphService) RejectFriendRequest(ctx context.Context, recipientID, senderID string) error {
	req, err := s.repo.GetRelationship(ctx, senderID, recipientID, domain.RelFriend)
	if err != nil || req.Status != domain.RelStatusPending {
		return ErrRequestNotFound
	}

	return s.repo.DeleteRelationship(ctx, senderID, recipientID, domain.RelFriend)
}

func (s *SocialGraphService) Unfriend(ctx context.Context, userA, userB string) error {
	_ = s.repo.DeleteRelationship(ctx, userA, userB, domain.RelFriend)
	_ = s.repo.DeleteRelationship(ctx, userB, userA, domain.RelFriend)
	return nil
}

func (s *SocialGraphService) BlockUser(ctx context.Context, blockerID, targetID string) error {
	if blockerID == targetID {
		return errors.New("cannot block yourself")
	}

	// Remove any active friendships
	_ = s.Unfriend(ctx, blockerID, targetID)

	now := time.Now().UTC()
	blockRel := &domain.Relationship{
		ID:               generateUUID(),
		UserID:           blockerID,
		TargetUserID:     targetID,
		RelationshipType: domain.RelBlock,
		Status:           domain.RelStatusAccepted,
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	return s.repo.SaveRelationship(ctx, blockRel)
}

func (s *SocialGraphService) GetFriends(ctx context.Context, userID string) ([]*domain.Relationship, error) {
	return s.repo.GetRelationshipsByUserAndType(ctx, userID, domain.RelFriend, domain.RelStatusAccepted)
}

func (s *SocialGraphService) GetPendingRequests(ctx context.Context, userID string) ([]*domain.Relationship, error) {
	return s.repo.GetPendingReceivedRequests(ctx, userID)
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
