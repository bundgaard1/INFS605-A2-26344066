package repository

import (
	"context"

	"osbourne.local/profile-service/internal/domain"
)

// ProfileRepository definerer alle database-operationer for profiler
type ProfileRepository interface {
	GetByID(ctx context.Context, id string) (*domain.UserProfile, error)
}
