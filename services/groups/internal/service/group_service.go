package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/dnjstjqppk-dev/facebook/pkg/domain"
	"github.com/dnjstjqppk-dev/facebook/services/groups/internal/repository"
)

var (
	ErrNotGroupAdmin = errors.New("operation requires group administrator or moderator privilege")
	ErrNotMember     = errors.New("user is not an approved member of this group")
)

type GroupService struct {
	repo repository.GroupRepository
}

func NewGroupService(repo repository.GroupRepository) *GroupService {
	return &GroupService{repo: repo}
}

func (s *GroupService) CreateGroup(ctx context.Context, creatorID, name, description string, privacy domain.GroupPrivacy, requirePostApproval bool) (*domain.Group, error) {
	now := time.Now().UTC()
	group := &domain.Group{
		ID:                  generateUUID(),
		Name:                name,
		Description:         description,
		Privacy:             privacy,
		CreatedBy:           creatorID,
		Members:             map[string]domain.GroupRole{creatorID: domain.GroupRoleAdmin},
		PendingMembers:      []string{},
		Rules:               []domain.GroupRule{},
		RequirePostApproval: requirePostApproval,
		CreatedAt:           now,
		UpdatedAt:           now,
	}

	if err := s.repo.SaveGroup(ctx, group); err != nil {
		return nil, err
	}
	return group, nil
}

func (s *GroupService) JoinGroup(ctx context.Context, userID, groupID string) (bool, error) {
	group, err := s.repo.GetGroup(ctx, groupID)
	if err != nil {
		return false, err
	}

	if _, isMember := group.Members[userID]; isMember {
		return true, nil // already a member
	}

	// Public group: instant join
	if group.Privacy == domain.GroupPrivacyPublic {
		group.Members[userID] = domain.GroupRoleMember
		_ = s.repo.SaveGroup(ctx, group)
		return true, nil
	}

	// Private group: queued in pending approval list
	for _, pid := range group.PendingMembers {
		if pid == userID {
			return false, nil // already pending
		}
	}
	group.PendingMembers = append(group.PendingMembers, userID)
	_ = s.repo.SaveGroup(ctx, group)
	return false, nil
}

func (s *GroupService) ApproveMember(ctx context.Context, adminID, groupID, targetUserID string) error {
	group, err := s.repo.GetGroup(ctx, groupID)
	if err != nil {
		return err
	}

	role, isMember := group.Members[adminID]
	if !isMember || (role != domain.GroupRoleAdmin && role != domain.GroupRoleModerator) {
		return ErrNotGroupAdmin
	}

	// Remove from pending, add to members
	var newPending []string
	found := false
	for _, pid := range group.PendingMembers {
		if pid == targetUserID {
			found = true
		} else {
			newPending = append(newPending, pid)
		}
	}
	if !found {
		return errors.New("user not in pending approval list")
	}

	group.PendingMembers = newPending
	group.Members[targetUserID] = domain.GroupRoleMember
	return s.repo.SaveGroup(ctx, group)
}

func (s *GroupService) CreateGroupPost(ctx context.Context, authorID, groupID, content string) (*domain.GroupPost, error) {
	group, err := s.repo.GetGroup(ctx, groupID)
	if err != nil {
		return nil, err
	}

	role, isMember := group.Members[authorID]
	if !isMember {
		return nil, ErrNotMember
	}

	status := domain.GroupPostApproved
	if group.RequirePostApproval && role == domain.GroupRoleMember {
		status = domain.GroupPostPending
	}

	post := &domain.GroupPost{
		ID:        generateUUID(),
		GroupID:   groupID,
		AuthorID:  authorID,
		Content:   content,
		Status:    status,
		CreatedAt: time.Now().UTC(),
	}

	if err := s.repo.SaveGroupPost(ctx, post); err != nil {
		return nil, err
	}
	return post, nil
}

func (s *GroupService) ModeratePost(ctx context.Context, adminID, groupID, postID string, approve bool) error {
	group, err := s.repo.GetGroup(ctx, groupID)
	if err != nil {
		return err
	}

	role, isMember := group.Members[adminID]
	if !isMember || (role != domain.GroupRoleAdmin && role != domain.GroupRoleModerator) {
		return ErrNotGroupAdmin
	}

	post, err := s.repo.GetGroupPost(ctx, postID)
	if err != nil {
		return err
	}

	if approve {
		post.Status = domain.GroupPostApproved
	} else {
		post.Status = domain.GroupPostRejected
	}

	return s.repo.SaveGroupPost(ctx, post)
}

func (s *GroupService) SetRules(ctx context.Context, adminID, groupID string, rules []domain.GroupRule) error {
	group, err := s.repo.GetGroup(ctx, groupID)
	if err != nil {
		return err
	}

	if group.Members[adminID] != domain.GroupRoleAdmin {
		return ErrNotGroupAdmin
	}

	group.Rules = rules
	return s.repo.SaveGroup(ctx, group)
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
