package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"osbourne.local/course-catalogue-service/internal/domain"
)

type CourseService struct {
	repo domain.CourseCatalogueRepository
}

func NewCourseService(repo domain.CourseCatalogueRepository) *CourseService {
	return &CourseService{repo: repo}
}

func (s *CourseService) GetCourse(ctx context.Context, courseID string) (*domain.Course, error) {
	return s.repo.GetCourse(ctx, courseID)
}

func (s *CourseService) ListCourses(ctx context.Context, page int32, pageSize int32) ([]*domain.Course, int32, error) {
	return s.repo.ListCourses(ctx, page, pageSize)
}

func (s *CourseService) EnrollStudent(ctx context.Context, courseID string, studentID string) error {

	courses, err := s.GetEnrolledCoursesByUserID(ctx, studentID)
	if err != nil {
		return err
	}

	for _, course := range courses {
		if course.ID == courseID {
			return errors.New("student is already enrolled in this course")
		}
	}

	e := &domain.Enrollment{
		ID:         uuid.NewString(),
		CourseID:   courseID,
		UserID:     studentID,
		EnrolledAt: time.Now(),
	}

	return s.repo.CreateEnrollment(ctx, e)
}

func (s *CourseService) GetEnrolledCoursesByUserID(ctx context.Context, userID string) ([]*domain.Course, error) {
	return s.repo.GetEnrolledCoursesByUserID(ctx, userID)
}
