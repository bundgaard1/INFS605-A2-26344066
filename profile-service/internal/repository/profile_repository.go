package repository

import (
	"context"

	"osbourne.local/profile-service/gen/profile"
)

// ProfileRepository definerer alle database-operationer for profiler
type ProfileRepository interface {
	GetByID(ctx context.Context, id string) (*profile.ProfileResponse, error)
}
