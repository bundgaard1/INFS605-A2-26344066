package repository

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
	"osbourne.local/users-service/gen/student"
	"osbourne.local/users-service/internal/domain"
)

type GORMStudentRepository struct {
	db *gorm.DB
}

func NewGORMStudentRepository(db *gorm.DB) *GORMStudentRepository {
	return &GORMStudentRepository{db: db}
}

func (r *GORMStudentRepository) GetByID(ctx context.Context, id string) (*student.StudentResponse, error) {
	var studentEntity domain.Student

	err := r.db.WithContext(ctx).
		Preload("Courses").
		First(&studentEntity, "id = ?", id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("studerende med id %s blev ikke fundet", id)
		}
		return nil, fmt.Errorf("databasefejl: %w", err)
	}

	// Map GORM domænemodel til Protobuf gRPC response
	var protoCourses []*student.Course
	for _, c := range studentEntity.Courses {
		protoCourses = append(protoCourses, &student.Course{
			Id:    c.ID,
			Title: c.Title,
		})
	}

	return &student.StudentResponse{
		Id:      studentEntity.ID,
		Name:    studentEntity.Name,
		Role:    studentEntity.Role,
		Courses: protoCourses,
	}, nil
}
