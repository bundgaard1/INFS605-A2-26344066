package repository

import (
	"context"

	"osbourne.local/course-catalogue-service/internal/domain"
)

// CourseCatalogueRepository definerer alle database-operationer for kurser
type CourseCatalogueRepository interface {
	GetCourse(ctx context.Context, courseID string) (*domain.Course, error)
	ListCourses(ctx context.Context, page int32, pageSize int32) ([]*domain.Course, int32, error)
	EnrollStudent(ctx context.Context, courseID string, userID string) error
	GetEnrolledCoursesByUserID(ctx context.Context, userID string) ([]*domain.Course, error)
}
