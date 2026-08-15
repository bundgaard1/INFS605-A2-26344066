package repository

import (
	"context"

	"gorm.io/gorm"
	"osbourne.local/course-catalogue-service/internal/domain"
)

type GORMCourseCatalogueRepository struct {
	db *gorm.DB
}

func NewGORMCourseCatalogueRepository(db *gorm.DB) *GORMCourseCatalogueRepository {
	return &GORMCourseCatalogueRepository{db: db}
}

func (r *GORMCourseCatalogueRepository) GetCourse(ctx context.Context, courseID string) (*domain.Course, error) {
	var course domain.Course
	if err := r.db.WithContext(ctx).First(&course, "id = ?", courseID).Error; err != nil {
		return nil, err
	}
	return &course, nil
}

func (r *GORMCourseCatalogueRepository) ListCourses(ctx context.Context, page int32, pageSize int32) ([]*domain.Course, int32, error) {
	var courses []*domain.Course
	var total int64

	offset := (page - 1) * pageSize

	if err := r.db.WithContext(ctx).Model(&domain.Course{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).Limit(int(pageSize)).Offset(int(offset)).Find(&courses).Error; err != nil {
		return nil, 0, err
	}

	return courses, int32(total), nil
}

func (r *GORMCourseCatalogueRepository) CreateEnrollment(ctx context.Context, enrollment *domain.Enrollment) error {
	return r.db.WithContext(ctx).Create(enrollment).Error
}

func (r *GORMCourseCatalogueRepository) GetEnrolledCoursesByUserID(ctx context.Context, userID string) ([]*domain.Course, error) {
	var courses []*domain.Course

	err := r.db.WithContext(ctx).
		Table("courses").
		Joins("JOIN enrollments ON enrollments.course_id = courses.id").
		Where("enrollments.user_id = ?", userID).
		Find(&courses).Error

	if err != nil {
		return nil, err
	}

	return courses, nil
}
