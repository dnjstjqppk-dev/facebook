package service

import (
	"context"
	"errors"

	"github.com/dnjstjqppk-dev/facebook/pkg/domain"
	"github.com/dnjstjqppk-dev/facebook/services/profiles/internal/repository"
)

type UpdateProfileRequest struct {
	FirstName          *string `json:"first_name"`
	LastName           *string `json:"last_name"`
	Bio                *string `json:"bio"`
	AvatarURL          *string `json:"avatar_url"`
	CoverURL           *string `json:"cover_url"`
	PrivacyLevel       *string `json:"privacy_level"` // PUBLIC, FRIENDS, ONLY_ME
	IsProfessionalMode *bool   `json:"is_professional_mode"`
}

type ProfileService struct {
	repo repository.ProfileRepository
}

func NewProfileService(repo repository.ProfileRepository) *ProfileService {
	return &ProfileService{repo: repo}
}

func (s *ProfileService) GetProfile(ctx context.Context, viewerID, username string) (*domain.Profile, error) {
	p, err := s.repo.GetProfileByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	// ABAC Field Privacy Filtering
	// If profile is ONLY_ME and viewer is not the owner, mask private info
	if p.PrivacyLevel == "ONLY_ME" && viewerID != p.UserID {
		masked := *p
		masked.Bio = ""
		masked.DateOfBirth = ""
		return &masked, nil
	}

	return p, nil
}

func (s *ProfileService) UpdateProfile(ctx context.Context, userID string, req *UpdateProfileRequest) (*domain.Profile, error) {
	p, err := s.repo.GetProfileByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if req.FirstName != nil && *req.FirstName != "" {
		p.FirstName = *req.FirstName
	}
	if req.LastName != nil && *req.LastName != "" {
		p.LastName = *req.LastName
	}
	if req.Bio != nil {
		p.Bio = *req.Bio
	}
	if req.AvatarURL != nil {
		p.AvatarURL = *req.AvatarURL
	}
	if req.CoverURL != nil {
		p.CoverURL = *req.CoverURL
	}
	if req.PrivacyLevel != nil {
		if *req.PrivacyLevel == "PUBLIC" || *req.PrivacyLevel == "FRIENDS" || *req.PrivacyLevel == "ONLY_ME" {
			p.PrivacyLevel = *req.PrivacyLevel
		} else {
			return nil, errors.New("invalid privacy level")
		}
	}
	if req.IsProfessionalMode != nil {
		p.IsProfessionalMode = *req.IsProfessionalMode
	}

	if err := s.repo.SaveProfile(ctx, p); err != nil {
		return nil, err
	}

	return p, nil
}

func (s *ProfileService) SeedProfile(ctx context.Context, p *domain.Profile) error {
	return s.repo.SaveProfile(ctx, p)
}
