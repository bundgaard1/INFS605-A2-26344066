package service

import (
	"context"

	"osbourne.local/course-catalogue-service/internal/domain"
	"osbourne.local/course-catalogue-service/internal/repository"
)

type CourseService struct {
	repo repository.CourseCatalogueRepository
}

func NewCourseService(repo repository.CourseCatalogueRepository) *CourseService {
	return &CourseService{repo: repo}
}

func (s *CourseService) GetCourse(ctx context.Context, courseID string) (*domain.Course, error) {
	return s.repo.GetCourse(ctx, courseID)
}

func (s *CourseService) ListCourses(ctx context.Context, page int32, pageSize int32) ([]*domain.Course, int32, error) {
	return s.repo.ListCourses(ctx, page, pageSize)
}

func (s *CourseService) EnrollStudent(ctx context.Context, courseID string, studentID string) error {
	return s.repo.EnrollStudent(ctx, courseID, studentID)
}

func (s *CourseService) GetEnrolledCoursesByUserID(ctx context.Context, userID string) ([]*domain.Course, error) {
	return s.repo.GetEnrolledCoursesByUserID(ctx, userID)
}
