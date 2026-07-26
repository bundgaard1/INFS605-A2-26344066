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
