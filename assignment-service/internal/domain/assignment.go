package domain

import "time"

type Assignment struct {
	ID          string       `gorm:"primaryKey" json:"id"`
	CourseID    string       `gorm:"index;not null" json:"course_id"` // Refers to Course.ID
	Title       string       `gorm:"not null" json:"title"`
	Description string       `json:"description"`
	DueDate     time.Time    `gorm:"not null" json:"due_date"`
	Submissions []Submission `gorm:"foreignKey:AssignmentID" json:"submissions,omitempty"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

type Submission struct {
	ID           string    `gorm:"primaryKey" json:"id"`
	AssignmentID string    `gorm:"index;not null" json:"assignment_id"`
	StudentID    string    `gorm:"index;not null" json:"student_id"`
	FileURL      string    `gorm:"not null" json:"file_url"` // Relative path inside the file store
	FileName     string    `json:"file_name"`                // Original filename, kept for metadata
	FileSize     int64     `json:"file_size"`                // Size in bytes
	SubmittedAt  time.Time `json:"submitted_at"`
	Grade        *Grade    `gorm:"foreignKey:SubmissionID" json:"grade,omitempty"`
}

type Grade struct {
	ID           string    `gorm:"primaryKey" json:"id"`
	SubmissionID string    `gorm:"uniqueIndex;not null" json:"submission_id"`
	GraderID     string    `gorm:"not null" json:"grader_id"`
	Score        int       `gorm:"not null" json:"score"`
	LetterGrade  string    `json:"letter_grade,omitempty"`
	Feedback     string    `json:"feedback,omitempty"`
	GradedAt     time.Time `json:"graded_at"`
}
