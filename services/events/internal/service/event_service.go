package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/dnjstjqppk-dev/facebook/pkg/domain"
	"github.com/dnjstjqppk-dev/facebook/services/events/internal/repository"
)

var (
	ErrInvalidEventTime = errors.New("event end time must be after start time")
	ErrInvalidRSVP      = errors.New("invalid RSVP status")
)

type EventService struct {
	repo repository.EventRepository
}

func NewEventService(repo repository.EventRepository) *EventService {
	return &EventService{repo: repo}
}

func (s *EventService) CreateEvent(ctx context.Context, hostID, title, description, location string, isOnline bool, startsAt, endsAt time.Time) (*domain.Event, error) {
	if endsAt.Before(startsAt) {
		return nil, ErrInvalidEventTime
	}

	event := &domain.Event{
		ID:          generateUUID(),
		HostID:      hostID,
		Title:       title,
		Description: description,
		Location:    location,
		IsOnline:    isOnline,
		StartsAt:    startsAt,
		EndsAt:      endsAt,
		RSVPs:       map[string]domain.RSVPStatus{hostID: domain.RSVPGoing},
		CreatedAt:   time.Now().UTC(),
	}

	if err := s.repo.SaveEvent(ctx, event); err != nil {
		return nil, err
	}
	return event, nil
}

func (s *EventService) RSVPEvent(ctx context.Context, userID, eventID string, status domain.RSVPStatus) error {
	switch status {
	case domain.RSVPGoing, domain.RSVPInterested, domain.RSVPDeclined:
	default:
		return ErrInvalidRSVP
	}

	event, err := s.repo.GetEvent(ctx, eventID)
	if err != nil {
		return err
	}

	if event.RSVPs == nil {
		event.RSVPs = make(map[string]domain.RSVPStatus)
	}
	event.RSVPs[userID] = status

	return s.repo.SaveEvent(ctx, event)
}

func (s *EventService) GetUserCalendarEvents(ctx context.Context, userID string) ([]*domain.Event, error) {
	return s.repo.ListEventsByUser(ctx, userID)
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
