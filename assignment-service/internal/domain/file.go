package domain

// SubmissionFile describes a file that is part of a submission. The file bytes
// live in the FileStorage backend; only the metadata is persisted in SQL.
type SubmissionFile struct {
	RelativePath string // Path inside the file store
	FileName     string // Original filename
	Size         int64  // Size in bytes
}
