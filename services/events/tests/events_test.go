package tests

import (
	"context"
	"testing"
	"time"

	"github.com/dnjstjqppk-dev/facebook/pkg/domain"
	"github.com/dnjstjqppk-dev/facebook/services/events/internal/repository"
	"github.com/dnjstjqppk-dev/facebook/services/events/internal/service"
)

func setupTestEventService() (*service.EventService, repository.EventRepository) {
	repo := repository.NewInMemoryEventRepository()
	svc := service.NewEventService(repo)
	return svc, repo
}

func TestCreateEvent_OnlineAndInPerson(t *testing.T) {
	svc, _ := setupTestEventService()
	ctx := context.Background()

	now := time.Now().UTC()
	start := now.Add(24 * time.Hour)
	end := start.Add(2 * time.Hour)

	event, err := svc.CreateEvent(ctx, "organizer-1", "Tech Meetup Seoul", "Developer talks and networking", "Gangnam, Seoul", false, start, end)
	if err != nil {
		t.Fatalf("failed to create event: %v", err)
	}

	if event.Title != "Tech Meetup Seoul" || event.RSVPs["organizer-1"] != domain.RSVPGoing {
		t.Errorf("event fields mismatch: %v", event)
	}

	// Invalid end time
	_, err = svc.CreateEvent(ctx, "organizer-1", "Bad Event", "", "", true, end, start)
	if err != service.ErrInvalidEventTime {
		t.Errorf("expected ErrInvalidEventTime, got: %v", err)
	}
}

func TestEventRSVP_GoingInterestedDeclined(t *testing.T) {
	svc, repo := setupTestEventService()
	ctx := context.Background()

	now := time.Now().UTC()
	event, _ := svc.CreateEvent(ctx, "organizer-1", "Concert in Park", "Live music", "Central Park", false, now.Add(48*time.Hour), now.Add(52*time.Hour))

	// User 1 says GOING
	err := svc.RSVPEvent(ctx, "guest-1", event.ID, domain.RSVPGoing)
	if err != nil {
		t.Fatalf("RSVP going failed: %v", err)
	}

	// User 2 says INTERESTED
	_ = svc.RSVPEvent(ctx, "guest-2", event.ID, domain.RSVPInterested)

	e, _ := repo.GetEvent(ctx, event.ID)
	if e.RSVPs["guest-1"] != domain.RSVPGoing || e.RSVPs["guest-2"] != domain.RSVPInterested {
		t.Errorf("RSVP statuses mismatch")
	}

	// User 1 changes to DECLINED
	_ = svc.RSVPEvent(ctx, "guest-1", event.ID, domain.RSVPDeclined)
	e, _ = repo.GetEvent(ctx, event.ID)
	if e.RSVPs["guest-1"] != domain.RSVPDeclined {
		t.Errorf("expected DECLINED status")
	}
}

func TestEvent_FilterByUserCalendar(t *testing.T) {
	svc, _ := setupTestEventService()
	ctx := context.Background()

	now := time.Now().UTC()
	e1, _ := svc.CreateEvent(ctx, "host-a", "Hackathon", "Build cool apps", "Online", true, now.Add(10*time.Hour), now.Add(20*time.Hour))
	e2, _ := svc.CreateEvent(ctx, "host-b", "Private Dinner", "Dinner", "Restaurant", false, now.Add(15*time.Hour), now.Add(17*time.Hour))

	// User joins Hackathon (GOING) but declines dinner
	_ = svc.RSVPEvent(ctx, "user-developer", e1.ID, domain.RSVPGoing)
	_ = svc.RSVPEvent(ctx, "user-developer", e2.ID, domain.RSVPDeclined)

	calendar, err := svc.GetUserCalendarEvents(ctx, "user-developer")
	if err != nil || len(calendar) != 1 {
		t.Fatalf("expected 1 event in calendar, got %d", len(calendar))
	}
	if calendar[0].ID != e1.ID {
		t.Errorf("expected Hackathon in calendar")
	}
}

