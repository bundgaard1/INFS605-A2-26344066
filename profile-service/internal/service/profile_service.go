package service

import (
	"context"

	"osbourne.local/profile-service/gen/profile"
	"osbourne.local/profile-service/internal/repository"
)

type ProfileService struct {
	repo repository.ProfileRepository
}

// Injecter repositoriet
func NewProfileService(repo repository.ProfileRepository) *ProfileService {
	return &ProfileService{repo: repo}
}

func (s *ProfileService) GetProfile(ctx context.Context, id string) (*profile.ProfileResponse, error) {
	// Kald databaselaget
	return s.repo.GetByID(ctx, id)
}
