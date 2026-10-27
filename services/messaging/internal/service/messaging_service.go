package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/dnjstjqppk-dev/facebook/pkg/domain"
	"github.com/dnjstjqppk-dev/facebook/services/messaging/internal/repository"
)

var (
	ErrNotConversationMember = errors.New("user is not a member of this conversation")
	ErrCannotUnsend          = errors.New("cannot unsend message: not original author")
	ErrInvalidCallOperation  = errors.New("invalid call operation or participant")
)

type MessagingService struct {
	repo repository.MessagingRepository
}

func NewMessagingService(repo repository.MessagingRepository) *MessagingService {
	return &MessagingService{repo: repo}
}

func (s *MessagingService) Create1on1Conversation(ctx context.Context, userA, userB string, isE2EE bool) (*domain.Conversation, error) {
	// Check existing 1:1 conversation between the two users
	existingList, _ := s.repo.ListUserConversations(ctx, userA)
	for _, c := range existingList {
		if c.Type == domain.ConvOneToOne {
			if _, hasB := c.Members[userB]; hasB {
				return c, nil
			}
		}
	}

	now := time.Now().UTC()
	conv := &domain.Conversation{
		ID:        generateUUID(),
		Type:      domain.ConvOneToOne,
		CreatorID: userA,
		Members: map[string]domain.MemberRole{
			userA: domain.MemberRoleMember,
			userB: domain.MemberRoleMember,
		},
		IsE2EE:    isE2EE,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.repo.SaveConversation(ctx, conv); err != nil {
		return nil, err
	}
	return conv, nil
}

func (s *MessagingService) CreateGroupConversation(ctx context.Context, creatorID, title string, memberIDs []string) (*domain.Conversation, error) {
	now := time.Now().UTC()
	members := make(map[string]domain.MemberRole)
	members[creatorID] = domain.MemberRoleAdmin

	for _, mid := range memberIDs {
		if mid != creatorID {
			members[mid] = domain.MemberRoleMember
		}
	}

	conv := &domain.Conversation{
		ID:        generateUUID(),
		Type:      domain.ConvGroup,
		Title:     title,
		CreatorID: creatorID,
		Members:   members,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.repo.SaveConversation(ctx, conv); err != nil {
		return nil, err
	}
	return conv, nil
}

func (s *MessagingService) CreateBroadcastChannel(ctx context.Context, creatorID, title string) (*domain.Conversation, error) {
	now := time.Now().UTC()
	conv := &domain.Conversation{
		ID:        generateUUID(),
		Type:      domain.ConvBroadcastChannel,
		Title:     title,
		CreatorID: creatorID,
		Members: map[string]domain.MemberRole{
			creatorID: domain.MemberRoleAdmin,
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.repo.SaveConversation(ctx, conv); err != nil {
		return nil, err
	}
	return conv, nil
}

func (s *MessagingService) SendMessage(ctx context.Context, senderID, convID, content, mediaURL, replyToID string) (*domain.Message, error) {
	conv, err := s.repo.GetConversation(ctx, convID)
	if err != nil {
		return nil, err
	}

	// Membership check
	if _, isMember := conv.Members[senderID]; !isMember {
		return nil, ErrNotConversationMember
	}

	now := time.Now().UTC()
	msg := &domain.Message{
		ID:               generateUUID(),
		ConversationID:   convID,
		SenderID:         senderID,
		Content:          content,
		MediaURL:         mediaURL,
		ReplyToMessageID: replyToID,
		IsUnsent:         false,
		IsE2EE:           conv.IsE2EE,
		Reactions:        make(map[string]domain.ReactionType),
		DeliveredTo:      make(map[string]time.Time),
		ReadBy:           make(map[string]time.Time),
		CreatedAt:        now,
	}

	// Sender automatically counts as delivered & read
	msg.DeliveredTo[senderID] = now
	msg.ReadBy[senderID] = now

	if err := s.repo.SaveMessage(ctx, msg); err != nil {
		return nil, err
	}
	return msg, nil
}

func (s *MessagingService) MarkDelivered(ctx context.Context, msgID, userID string) error {
	return s.repo.MarkDelivered(ctx, msgID, userID, time.Now().UTC())
}

func (s *MessagingService) MarkRead(ctx context.Context, convID, msgID, userID string) error {
	return s.repo.MarkRead(ctx, convID, msgID, userID, time.Now().UTC())
}

func (s *MessagingService) UnsendMessage(ctx context.Context, msgID, requesterID string) error {
	msg, err := s.repo.GetMessage(ctx, msgID)
	if err != nil {
		return err
	}

	if msg.SenderID != requesterID {
		return ErrCannotUnsend
	}

	return s.repo.UnsendMessage(ctx, msgID)
}

func (s *MessagingService) ReactToMessage(ctx context.Context, msgID, userID string, reaction domain.ReactionType) error {
	return s.repo.AddMessageReaction(ctx, msgID, userID, reaction)
}

func (s *MessagingService) SetPresence(ctx context.Context, userID string, status domain.PresenceStatus) error {
	return s.repo.UpdatePresence(ctx, userID, status, time.Now().UTC())
}

func (s *MessagingService) GetPresence(ctx context.Context, userID string) (*domain.UserPresence, error) {
	return s.repo.GetPresence(ctx, userID)
}

func (s *MessagingService) InitiateCall(ctx context.Context, callerID string, receiverIDs []string, callType domain.CallType) (*domain.CallSession, error) {
	now := time.Now().UTC()
	call := &domain.CallSession{
		ID:          generateUUID(),
		CallerID:    callerID,
		ReceiverIDs: receiverIDs,
		CallType:    callType,
		Status:      domain.CallInitiated,
		CreatedAt:   now,
	}

	if err := s.repo.SaveCallSession(ctx, call); err != nil {
		return nil, err
	}
	return call, nil
}

func (s *MessagingService) AnswerCall(ctx context.Context, callID, receiverID string) (*domain.CallSession, error) {
	call, err := s.repo.GetCallSession(ctx, callID)
	if err != nil {
		return nil, err
	}

	isReceiver := false
	for _, id := range call.ReceiverIDs {
		if id == receiverID {
			isReceiver = true
			break
		}
	}
	if !isReceiver {
		return nil, ErrInvalidCallOperation
	}

	now := time.Now().UTC()
	call.Status = domain.CallConnected
	call.StartedAt = &now

	if err := s.repo.SaveCallSession(ctx, call); err != nil {
		return nil, err
	}
	return call, nil
}

func (s *MessagingService) EndCall(ctx context.Context, callID, userID string) (*domain.CallSession, error) {
	call, err := s.repo.GetCallSession(ctx, callID)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	call.Status = domain.CallEnded
	call.EndedAt = &now

	if err := s.repo.SaveCallSession(ctx, call); err != nil {
		return nil, err
	}
	return call, nil
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

