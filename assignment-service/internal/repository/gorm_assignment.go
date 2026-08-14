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

func (r *GORMAssignmentRepository) SetGrade(ctx context.Context, grade *domain.Grade) error {
	return r.db.WithContext(ctx).Create(grade).Error
}
