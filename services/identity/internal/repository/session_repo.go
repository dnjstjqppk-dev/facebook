package repository

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/dnjstjqppk-dev/facebook/services/identity/internal/domain"
)

var (
	ErrSessionNotFound = errors.New("session not found")
	ErrSessionRevoked  = errors.New("session has been revoked")
	ErrSessionExpired  = errors.New("session has expired")
)

type SessionRepository interface {
	SaveSession(ctx context.Context, session *domain.UserSession) error
	GetSessionByID(ctx context.Context, sessionID string) (*domain.UserSession, error)
	ListUserSessions(ctx context.Context, userID string) ([]*domain.UserSession, error)
	RevokeSession(ctx context.Context, sessionID, reason string) error
	RevokeAllUserSessionsExcept(ctx context.Context, userID, currentSessionID, reason string) error
	UpdateLastActive(ctx context.Context, sessionID string, t time.Time) error
}

type InMemorySessionRepository struct {
	mu       sync.RWMutex
	sessions map[string]*domain.UserSession // key: session_id
}

func NewInMemorySessionRepository() *InMemorySessionRepository {
	return &InMemorySessionRepository{
		sessions: make(map[string]*domain.UserSession),
	}
}

func (r *InMemorySessionRepository) SaveSession(ctx context.Context, s *domain.UserSession) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.sessions[s.ID] = s
	return nil
}

func (r *InMemorySessionRepository) GetSessionByID(ctx context.Context, sessionID string) (*domain.UserSession, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	s, exists := r.sessions[sessionID]
	if !exists {
		return nil, ErrSessionNotFound
	}
	return s, nil
}

func (r *InMemorySessionRepository) ListUserSessions(ctx context.Context, userID string) ([]*domain.UserSession, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var results []*domain.UserSession
	for _, s := range r.sessions {
		if s.UserID == userID {
			results = append(results, s)
		}
	}
	return results, nil
}

func (r *InMemorySessionRepository) RevokeSession(ctx context.Context, sessionID, reason string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	s, exists := r.sessions[sessionID]
	if !exists {
		return ErrSessionNotFound
	}
	s.IsRevoked = true
	s.RevocationReason = reason
	return nil
}

func (r *InMemorySessionRepository) RevokeAllUserSessionsExcept(ctx context.Context, userID, currentSessionID, reason string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, s := range r.sessions {
		if s.UserID == userID && s.ID != currentSessionID && !s.IsRevoked {
			s.IsRevoked = true
			s.RevocationReason = reason
		}
	}
	return nil
}

func (r *InMemorySessionRepository) UpdateLastActive(ctx context.Context, sessionID string, t time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	s, exists := r.sessions[sessionID]
	if !exists {
		return ErrSessionNotFound
	}
	s.LastActiveAt = t
	return nil
}

