package repository

import (
	"context"
	"errors"
	"sync"

	"github.com/dnjstjqppk-dev/facebook/services/identity/internal/domain"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user with this email already exists")
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *domain.User, profile *domain.Profile) error
	GetUserByEmail(ctx context.Context, email string) (*domain.User, *domain.Profile, error)
	GetUserByID(ctx context.Context, id string) (*domain.User, *domain.Profile, error)
}

type InMemoryUserRepository struct {
	mu       sync.RWMutex
	users    map[string]*domain.User    // key: id
	byEmail  map[string]string          // email -> id
	profiles map[string]*domain.Profile // user_id -> profile
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		users:    make(map[string]*domain.User),
		byEmail:  make(map[string]string),
		profiles: make(map[string]*domain.Profile),
	}
}

func (r *InMemoryUserRepository) CreateUser(ctx context.Context, user *domain.User, profile *domain.Profile) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.byEmail[user.Email]; exists {
		return ErrUserAlreadyExists
	}

	r.users[user.ID] = user
	r.byEmail[user.Email] = user.ID
	r.profiles[user.ID] = profile

	return nil
}

func (r *InMemoryUserRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, *domain.Profile, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	id, exists := r.byEmail[email]
	if !exists {
		return nil, nil, ErrUserNotFound
	}

	user := r.users[id]
	profile := r.profiles[id]
	return user, profile, nil
}

func (r *InMemoryUserRepository) GetUserByID(ctx context.Context, id string) (*domain.User, *domain.Profile, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.users[id]
	if !exists {
		return nil, nil, ErrUserNotFound
	}
	profile := r.profiles[id]
	return user, profile, nil
}
