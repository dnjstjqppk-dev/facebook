package tests

import (
	"context"
	"testing"
	"time"

	"github.com/dnjstjqppk-dev/facebook/pkg/domain"
	"github.com/dnjstjqppk-dev/facebook/services/profiles/internal/repository"
	"github.com/dnjstjqppk-dev/facebook/services/profiles/internal/service"
)

func setupTestProfileService() *service.ProfileService {
	repo := repository.NewInMemoryProfileRepository()
	return service.NewProfileService(repo)
}

func TestGetAndSeedProfile(t *testing.T) {
	svc := setupTestProfileService()
	ctx := context.Background()

	p := &domain.Profile{
		UserID:             "user-100",
		Username:           "johndoe",
		FirstName:          "John",
		LastName:           "Doe",
		Bio:                "Software Engineer & Nova Creator",
		DateOfBirth:        "1995-10-10",
		Gender:             "MALE",
		PrivacyLevel:       "PUBLIC",
		IsProfessionalMode: false,
		CreatedAt:          time.Now(),
	}

	if err := svc.SeedProfile(ctx, p); err != nil {
		t.Fatalf("seeding failed: %v", err)
	}

	fetched, err := svc.GetProfile(ctx, "visitor-1", "johndoe")
	if err != nil {
		t.Fatalf("failed to get profile: %v", err)
	}

	if fetched.FirstName != "John" || fetched.Bio != "Software Engineer & Nova Creator" {
		t.Errorf("unexpected profile data: %+v", fetched)
	}
}

func TestUpdateProfile(t *testing.T) {
	svc := setupTestProfileService()
	ctx := context.Background()

	p := &domain.Profile{
		UserID:       "user-100",
		Username:     "johndoe",
		FirstName:    "John",
		LastName:     "Doe",
		Bio:          "Old bio",
		PrivacyLevel: "PUBLIC",
	}
	_ = svc.SeedProfile(ctx, p)

	newBio := "Updated bio in Nova"
	newFirst := "Jonathan"
	isPro := true

	updated, err := svc.UpdateProfile(ctx, "user-100", &service.UpdateProfileRequest{
		FirstName:          &newFirst,
		Bio:                &newBio,
		IsProfessionalMode: &isPro,
	})
	if err != nil {
		t.Fatalf("failed to update profile: %v", err)
	}

	if updated.FirstName != "Jonathan" || updated.Bio != "Updated bio in Nova" || !updated.IsProfessionalMode {
		t.Errorf("update failed to apply: %+v", updated)
	}
}

func TestPrivacyLevelMasking(t *testing.T) {
	svc := setupTestProfileService()
	ctx := context.Background()

	p := &domain.Profile{
		UserID:       "user-secret",
		Username:     "secretuser",
		FirstName:    "Secret",
		LastName:     "Agent",
		Bio:          "Classified Information",
		DateOfBirth:  "1990-01-01",
		PrivacyLevel: "ONLY_ME",
	}
	_ = svc.SeedProfile(ctx, p)

	// Owner views own profile -> sees bio
	ownerView, err := svc.GetProfile(ctx, "user-secret", "secretuser")
	if err != nil || ownerView.Bio != "Classified Information" {
		t.Errorf("owner should see full bio")
	}

	// Stranger views profile -> bio and DOB masked
	strangerView, err := svc.GetProfile(ctx, "stranger-99", "secretuser")
	if err != nil {
		t.Fatalf("stranger should still find profile: %v", err)
	}
	if strangerView.Bio != "" || strangerView.DateOfBirth != "" {
		t.Errorf("stranger should not see masked bio/DOB on ONLY_ME profile: %+v", strangerView)
	}
}
