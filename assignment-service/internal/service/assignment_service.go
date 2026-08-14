package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"osbourne.local/assignment-service/internal/domain"
)

var ErrAssignmentNotFound = errors.New("assignment not found")

type AssignmentService struct {
	repo      domain.AssignmentRepository
	fileStore domain.FileStorage
}

func NewAssignmentService(repo domain.AssignmentRepository, fileStore domain.FileStorage) *AssignmentService {
	return &AssignmentService{
		repo:      repo,
		fileStore: fileStore,
	}
}

func (s *AssignmentService) CreateAssignment(ctx context.Context, assignment *domain.Assignment) error {
	if assignment.ID == "" {
		assignment.ID = uuid.NewString()
	}
	assignment.CreatedAt = time.Now()
	assignment.UpdatedAt = time.Now()

	return s.repo.CreateAssignment(ctx, assignment)
}

func (s *AssignmentService) GetAssignment(ctx context.Context, assignmentID string) (*domain.Assignment, error) {
	return s.repo.GetAssignment(ctx, assignmentID)
}

type SubmitAssignmentInput struct {
	AssignmentID string
	StudentID    string
	FileName     string
	Size         int64
}

func (s *AssignmentService) SubmitAssignment(ctx context.Context, in SubmitAssignmentInput, src io.Reader) (*domain.Submission, error) {
	// 1. Validate that the assignment exists
	assignment, err := s.repo.GetAssignment(ctx, in.AssignmentID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch assignment: %w", err)
	}
	if assignment == nil {
		return nil, ErrAssignmentNotFound
	}

	// 2. Build the storage path from server-side identifiers only, so the
	//    relative path is safe against path traversal.
	submissionID := uuid.NewString()
	fileID := uuid.NewString()
	relativePath := fmt.Sprintf(
		"submissions/%s/%s/%s%s",
		in.AssignmentID,
		submissionID,
		fileID,
		sanitizedExtension(in.FileName),
	)

	// 3. Save the physical file
	savedPath, err := s.fileStore.Save(ctx, relativePath, src)
	if err != nil {
		return nil, fmt.Errorf("failed to save file: %w", err)
	}

	// 4. Persist the submission metadata
	submission := &domain.Submission{
		ID:           submissionID,
		AssignmentID: in.AssignmentID,
		StudentID:    in.StudentID,
		FileURL:      savedPath,
		FileName:     in.FileName,
		FileSize:     in.Size,
		SubmittedAt:  time.Now(),
	}
	if err := s.repo.CreateSubmission(ctx, submission); err != nil {
		// ROLLBACK: if the database write fails, remove the saved file again
		_ = s.fileStore.Delete(ctx, savedPath)
		return nil, fmt.Errorf("failed to persist submission: %w", err)
	}

	return submission, nil
}

func (s *AssignmentService) ListSubmissionsByAssignment(ctx context.Context, assignmentID string) ([]*domain.Submission, error) {
	return s.repo.ListSubmissionsByAssignment(ctx, assignmentID)
}

func (s *AssignmentService) GradeSubmission(ctx context.Context, submissionID string, graderID string, score int, feedback string) (*domain.Grade, error) {
	grade := &domain.Grade{
		ID:           uuid.NewString(),
		SubmissionID: submissionID,
		GraderID:     graderID,
		Score:        score,
		Feedback:     feedback,
		GradedAt:     time.Now(),
	}

	if err := s.repo.SetGrade(ctx, grade); err != nil {
		return nil, fmt.Errorf("failed to store grade: %w", err)
	}

	return grade, nil
}

// sanitizedExtension returns the lowercase file extension from filename,
// allowing only [a-z0-9] characters after the leading dot. It returns an
// empty string when there is no safe extension to keep.
func sanitizedExtension(filename string) string {
	ext := strings.ToLower(filepath.Ext(filepath.Base(filename)))
	if len(ext) < 2 {
		return ""
	}

	for _, r := range ext {
		if r == '.' {
			continue
		}
		if (r < 'a' || r > 'z') && (r < '0' || r > '9') {
			return ""
		}
	}

	if len(ext) > 16 {
		ext = ext[:16]
	}
	return ext
}
