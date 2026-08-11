package service

import (
	"context"

	"osbourne.local/profile-service/internal/domain"
	"osbourne.local/profile-service/internal/repository"
)

type ProfileService struct {
	repo repository.ProfileRepository
}

// Injects the repository
func NewProfileService(repo repository.ProfileRepository) *ProfileService {
	return &ProfileService{repo: repo}
}

func (s *ProfileService) GetProfile(ctx context.Context, id string) (*domain.UserProfile, error) {
	// Call the database layer
	return s.repo.GetByID(ctx, id)
}
