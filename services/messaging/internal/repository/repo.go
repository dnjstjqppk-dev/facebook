package repository

import (
	"context"
	"errors"
	"sort"
	"sync"
	"time"

	"github.com/dnjstjqppk-dev/facebook/pkg/domain"
)

var (
	ErrConversationNotFound = errors.New("conversation not found")
	ErrMessageNotFound      = errors.New("message not found")
	ErrCallNotFound         = errors.New("call session not found")
)

type MessagingRepository interface {
	SaveConversation(ctx context.Context, conv *domain.Conversation) error
	GetConversation(ctx context.Context, id string) (*domain.Conversation, error)
	ListUserConversations(ctx context.Context, userID string) ([]*domain.Conversation, error)

	SaveMessage(ctx context.Context, msg *domain.Message) error
	GetMessage(ctx context.Context, id string) (*domain.Message, error)
	ListMessages(ctx context.Context, convID string, limit int) ([]*domain.Message, error)

	MarkDelivered(ctx context.Context, msgID, userID string, t time.Time) error
	MarkRead(ctx context.Context, convID, msgID, userID string, t time.Time) error
	UnsendMessage(ctx context.Context, msgID string) error
	AddMessageReaction(ctx context.Context, msgID, userID string, r domain.ReactionType) error

	UpdatePresence(ctx context.Context, userID string, status domain.PresenceStatus, t time.Time) error
	GetPresence(ctx context.Context, userID string) (*domain.UserPresence, error)

	SaveCallSession(ctx context.Context, call *domain.CallSession) error
	GetCallSession(ctx context.Context, id string) (*domain.CallSession, error)
}

type InMemoryMessagingRepository struct {
	mu            sync.RWMutex
	conversations map[string]*domain.Conversation
	messages      map[string]*domain.Message
	convMessages  map[string][]string // convID -> []msgID
	presence      map[string]*domain.UserPresence
	calls         map[string]*domain.CallSession
}

func NewInMemoryMessagingRepository() *InMemoryMessagingRepository {
	return &InMemoryMessagingRepository{
		conversations: make(map[string]*domain.Conversation),
		messages:      make(map[string]*domain.Message),
		convMessages:  make(map[string][]string),
		presence:      make(map[string]*domain.UserPresence),
		calls:         make(map[string]*domain.CallSession),
	}
}

func (r *InMemoryMessagingRepository) SaveConversation(ctx context.Context, conv *domain.Conversation) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.conversations[conv.ID] = conv
	return nil
}

func (r *InMemoryMessagingRepository) GetConversation(ctx context.Context, id string) (*domain.Conversation, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	c, exists := r.conversations[id]
	if !exists {
		return nil, ErrConversationNotFound
	}
	return c, nil
}

func (r *InMemoryMessagingRepository) ListUserConversations(ctx context.Context, userID string) ([]*domain.Conversation, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var results []*domain.Conversation
	for _, c := range r.conversations {
		if _, isMember := c.Members[userID]; isMember {
			results = append(results, c)
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].LastMessageAt.After(results[j].LastMessageAt)
	})
	return results, nil
}

func (r *InMemoryMessagingRepository) SaveMessage(ctx context.Context, msg *domain.Message) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.messages[msg.ID] = msg
	r.convMessages[msg.ConversationID] = append(r.convMessages[msg.ConversationID], msg.ID)

	if c, exists := r.conversations[msg.ConversationID]; exists {
		c.LastMessageSnippet = msg.Content
		c.LastMessageAt = msg.CreatedAt
		c.UpdatedAt = msg.CreatedAt
	}
	return nil
}

func (r *InMemoryMessagingRepository) GetMessage(ctx context.Context, id string) (*domain.Message, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	m, exists := r.messages[id]
	if !exists {
		return nil, ErrMessageNotFound
	}
	return m, nil
}

func (r *InMemoryMessagingRepository) ListMessages(ctx context.Context, convID string, limit int) ([]*domain.Message, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	msgIDs, exists := r.convMessages[convID]
	if !exists {
		return []*domain.Message{}, nil
	}

	var results []*domain.Message
	for _, id := range msgIDs {
		if m, ok := r.messages[id]; ok {
			results = append(results, m)
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].CreatedAt.Before(results[j].CreatedAt)
	})

	if limit > 0 && len(results) > limit {
		results = results[len(results)-limit:]
	}
	return results, nil
}

func (r *InMemoryMessagingRepository) MarkDelivered(ctx context.Context, msgID, userID string, t time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	msg, exists := r.messages[msgID]
	if !exists {
		return ErrMessageNotFound
	}
	if msg.DeliveredTo == nil {
		msg.DeliveredTo = make(map[string]time.Time)
	}
	msg.DeliveredTo[userID] = t
	return nil
}

func (r *InMemoryMessagingRepository) MarkRead(ctx context.Context, convID, msgID, userID string, t time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	msg, exists := r.messages[msgID]
	if !exists {
		return ErrMessageNotFound
	}
	if msg.ReadBy == nil {
		msg.ReadBy = make(map[string]time.Time)
	}
	msg.ReadBy[userID] = t
	return nil
}

func (r *InMemoryMessagingRepository) UnsendMessage(ctx context.Context, msgID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	msg, exists := r.messages[msgID]
	if !exists {
		return ErrMessageNotFound
	}
	msg.IsUnsent = true
	msg.Content = "[This message was unsent]"
	return nil
}

func (r *InMemoryMessagingRepository) AddMessageReaction(ctx context.Context, msgID, userID string, reaction domain.ReactionType) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	msg, exists := r.messages[msgID]
	if !exists {
		return ErrMessageNotFound
	}
	if msg.Reactions == nil {
		msg.Reactions = make(map[string]domain.ReactionType)
	}
	msg.Reactions[userID] = reaction
	return nil
}

func (r *InMemoryMessagingRepository) UpdatePresence(ctx context.Context, userID string, status domain.PresenceStatus, t time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.presence[userID] = &domain.UserPresence{
		UserID:     userID,
		Status:     status,
		LastSeenAt: t,
	}
	return nil
}

func (r *InMemoryMessagingRepository) GetPresence(ctx context.Context, userID string) (*domain.UserPresence, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	p, exists := r.presence[userID]
	if !exists {
		return &domain.UserPresence{
			UserID:     userID,
			Status:     domain.PresenceOffline,
			LastSeenAt: time.Now().UTC(),
		}, nil
	}
	return p, nil
}

func (r *InMemoryMessagingRepository) SaveCallSession(ctx context.Context, call *domain.CallSession) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.calls[call.ID] = call
	return nil
}

func (r *InMemoryMessagingRepository) GetCallSession(ctx context.Context, id string) (*domain.CallSession, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	call, exists := r.calls[id]
	if !exists {
		return nil, ErrCallNotFound
	}
	return call, nil
}

