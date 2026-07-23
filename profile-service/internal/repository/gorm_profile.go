package repository

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
	"osbourne.local/profile-service/gen/profile"
	"osbourne.local/profile-service/internal/domain"
)

type GORMProfileRepository struct {
	db *gorm.DB
}

func NewGORMProfileRepository(db *gorm.DB) *GORMProfileRepository {
	return &GORMProfileRepository{db: db}
}

func (r *GORMProfileRepository) GetByID(ctx context.Context, id string) (*profile.ProfileResponse, error) {
	var profileEntity domain.UserProfile

	err := r.db.WithContext(ctx).
		Preload("Courses").
		First(&profileEntity, "id = ?", id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("studerende med id %s blev ikke fundet", id)
		}
		return nil, fmt.Errorf("databasefejl: %w", err)
	}

	// Map GORM domænemodel til Protobuf gRPC response

	return &profile.ProfileResponse{
		Id:   profileEntity.ID,
		Name: profileEntity.Name,
		Role: string(profileEntity.Role),
	}, nil
}
