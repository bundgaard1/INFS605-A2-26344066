package repository

import (
	"context"

	"gorm.io/gorm"
	"osbourne.local/assignment-service/internal/domain"
)

type GORMAssignmentRepository struct {
	db *gorm.DB
}

func NewGORMAssignmentRepository(db *gorm.DB) *GORMAssignmentRepository {
	return &GORMAssignmentRepository{db: db}
}

func (r *GORMAssignmentRepository) CreateAssignment(ctx context.Context, assignment *domain.Assignment) error {
	return r.db.WithContext(ctx).Create(assignment).Error
}

func (r *GORMAssignmentRepository) GetAssignment(ctx context.Context, assignmentID string) (*domain.Assignment, error) {
	var assignment domain.Assignment
	if err := r.db.WithContext(ctx).First(&assignment, "id = ?", assignmentID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &assignment, nil
}

func (r *GORMAssignmentRepository) GetAssignmentsByCourse(ctx context.Context, courseID string) ([]*domain.Assignment, error) {
	var assignments []*domain.Assignment
	if err := r.db.WithContext(ctx).Where("course_id = ?", courseID).Find(&assignments).Error; err != nil {
		return nil, err
	}
	return assignments, nil
}

func (r *GORMAssignmentRepository) CreateSubmission(ctx context.Context, submission *domain.Submission) error {
	return r.db.WithContext(ctx).Create(submission).Error
}

func (r *GORMAssignmentRepository) ListSubmissionsByAssignment(ctx context.Context, assignmentID string) ([]*domain.Submission, error) {
	var submissions []*domain.Submission
	if err := r.db.WithContext(ctx).
		Where("assignment_id = ?", assignmentID).
		Preload("Grade").
		Find(&submissions).Error; err != nil {
		return nil, err
	}
	return submissions, nil
}

func (r *GORMAssignmentRepository) CreateGrade(ctx context.Context, grade *domain.Grade) error {
	return r.db.WithContext(ctx).Create(grade).Error
}

func (r *GORMAssignmentRepository) GetGradeBySubmission(ctx context.Context, submissionID string) (*domain.Grade, error) {
	var grade domain.Grade
	if err := r.db.WithContext(ctx).First(&grade, "submission_id = ?", submissionID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &grade, nil
}
