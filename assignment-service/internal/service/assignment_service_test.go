package service

import (
	"context"
	"errors"
	"io"
	"regexp"
	"strings"
	"testing"

	"osbourne.local/assignment-service/internal/domain"
)

type fakeRepository struct {
	assignments map[string]*domain.Assignment
	submissions []*domain.Submission
	failCreate  bool
	lastGrade   *domain.Grade
}

func newFakeRepository() *fakeRepository {
	return &fakeRepository{
		assignments: map[string]*domain.Assignment{},
	}
}

func (f *fakeRepository) CreateAssignment(ctx context.Context, assignment *domain.Assignment) error {
	f.assignments[assignment.ID] = assignment
	return nil
}

func (f *fakeRepository) GetAssignment(ctx context.Context, assignmentID string) (*domain.Assignment, error) {
	return f.assignments[assignmentID], nil
}

func (f *fakeRepository) CreateSubmission(ctx context.Context, submission *domain.Submission) error {
	if f.failCreate {
		return errors.New("db down")
	}
	f.submissions = append(f.submissions, submission)
	return nil
}

func (f *fakeRepository) ListSubmissionsByAssignment(ctx context.Context, assignmentID string) ([]*domain.Submission, error) {
	var out []*domain.Submission
	for _, s := range f.submissions {
		if s.AssignmentID == assignmentID {
			out = append(out, s)
		}
	}
	return out, nil
}

func (f *fakeRepository) SetGrade(ctx context.Context, grade *domain.Grade) error {
	f.lastGrade = grade
	return nil
}

type fakeStorage struct {
	files   map[string][]byte
	deleted []string
}

func newFakeStorage() *fakeStorage {
	return &fakeStorage{files: map[string][]byte{}}
}

func (s *fakeStorage) Save(ctx context.Context, relativePath string, src io.Reader) (string, error) {
	data, err := io.ReadAll(src)
	if err != nil {
		return "", err
	}
	s.files[relativePath] = data
	return relativePath, nil
}

func (s *fakeStorage) Get(ctx context.Context, relativePath string) (io.ReadCloser, error) {
	return nil, nil
}

func (s *fakeStorage) Delete(ctx context.Context, relativePath string) error {
	delete(s.files, relativePath)
	s.deleted = append(s.deleted, relativePath)
	return nil
}

var uuidLike = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)

func setupService(t *testing.T) (*AssignmentService, *fakeRepository, *fakeStorage) {
	repo := newFakeRepository()
	storage := newFakeStorage()
	svc := NewAssignmentService(repo, storage)
	return svc, repo, storage
}

func TestSubmitAssignment_SavesUnderUuidPathAndPersistsMetadata(t *testing.T) {
	svc, repo, storage := setupService(t)

	repo.assignments["assn_1"] = &domain.Assignment{ID: "assn_1", Title: "HW1"}

	submission, err := svc.SubmitAssignment(context.Background(), SubmitAssignmentInput{
		AssignmentID: "assn_1",
		StudentID:    "student_1",
		FileName:     "report.pdf",
		Size:         10,
	}, strings.NewReader("filecontent"))
	if err != nil {
		t.Fatalf("SubmitAssignment failed: %v", err)
	}

	// Path must be built only from server-side IDs + a UUID name
	parts := strings.Split(submission.FileURL, "/")
	if len(parts) != 4 || parts[0] != "submissions" || parts[1] != "assn_1" {
		t.Fatalf("unexpected file url: %q", submission.FileURL)
	}
	if parts[2] != submission.ID {
		t.Errorf("directory should be the submission id, got %q", parts[2])
	}
	name := strings.TrimSuffix(parts[3], ".pdf")
	if !uuidLike.MatchString(name) {
		t.Errorf("stored filename should be a uuid, got %q", parts[3])
	}
	if strings.Contains(parts[3], "report") {
		t.Errorf("original filename leaked into storage path: %q", submission.FileURL)
	}

	if submission.FileName != "report.pdf" || submission.FileSize != 10 {
		t.Errorf("metadata mismatch: %+v", submission)
	}

	if got := string(storage.files[submission.FileURL]); got != "filecontent" {
		t.Errorf("stored content mismatch: %q", got)
	}

	if len(repo.submissions) != 1 {
		t.Fatalf("expected 1 submission in repo, got %d", len(repo.submissions))
	}
}

func TestSubmitAssignment_RemovesFileWhenDbFails(t *testing.T) {
	svc, repo, storage := setupService(t)

	repo.assignments["assn_1"] = &domain.Assignment{ID: "assn_1"}
	repo.failCreate = true

	_, err := svc.SubmitAssignment(context.Background(), SubmitAssignmentInput{
		AssignmentID: "assn_1",
		StudentID:    "student_1",
		FileName:     "report.pdf",
		Size:         10,
	}, strings.NewReader("filecontent"))
	if err == nil {
		t.Fatal("expected error when DB write fails")
	}

	if len(storage.files) != 0 {
		t.Errorf("file should have been rolled back, still stored: %v", storage.files)
	}
	if len(storage.deleted) != 1 {
		t.Errorf("expected exactly 1 delete call, got %d", len(storage.deleted))
	}
}

func TestSubmitAssignment_AssignmentNotFound(t *testing.T) {
	svc, _, storage := setupService(t)

	_, err := svc.SubmitAssignment(context.Background(), SubmitAssignmentInput{
		AssignmentID: "missing",
		StudentID:    "student_1",
		FileName:     "report.pdf",
		Size:         10,
	}, strings.NewReader("filecontent"))
	if !errors.Is(err, ErrAssignmentNotFound) {
		t.Fatalf("expected ErrAssignmentNotFound, got %v", err)
	}

	if len(storage.files) != 0 {
		t.Errorf("no file should be stored for a missing assignment: %v", storage.files)
	}
}

func TestSanitizedExtension(t *testing.T) {
	cases := map[string]string{
		"report.pdf":         ".pdf",
		"report.PDF":         ".pdf",
		"../evil/report.PDF": ".pdf",
		"noext":              "",
		"report.":            "",
		"report.exe;rm":      "",
		"report.tar.gz":      ".gz",
		"a.pdf  ":            "",
		"report.aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa": ".aaaaaaaaaaaaaaa", // capped at 16 chars incl. dot
	}
	for in, want := range cases {
		if got := sanitizedExtension(in); got != want {
			t.Errorf("sanitizedExtension(%q) = %q, want %q", in, got, want)
		}
	}
}
