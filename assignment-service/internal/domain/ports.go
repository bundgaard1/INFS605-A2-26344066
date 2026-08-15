package domain

import (
	"context"
	"io"
)

type AssignmentRepository interface {
	CreateAssignment(ctx context.Context, assignment *Assignment) error
	GetAssignment(ctx context.Context, assignmentID string) (*Assignment, error)
	GetAssignmentsByCourse(ctx context.Context, courseID string) ([]*Assignment, error)
	CreateSubmission(ctx context.Context, submission *Submission) error
	ListSubmissionsByAssignment(ctx context.Context, assignmentID string) ([]*Submission, error)
	CreateGrade(ctx context.Context, grade *Grade) error
	GetGradeBySubmission(ctx context.Context, submissionID string) (*Grade, error)
}

// FileStorage persists the raw bytes of submission files on the local
// filesystem. Paths are relative to the storage root and are built from
// server-side identifiers only.
type FileStorage interface {
	Save(ctx context.Context, relativePath string, src io.Reader) (string, error)
	Get(ctx context.Context, relativePath string) (io.ReadCloser, error)
	Delete(ctx context.Context, relativePath string) error
}
