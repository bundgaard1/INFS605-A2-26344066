Here is the complete Go markdown overview of the **domain structs (GORM/Entities)** for all 6 microservices.

The structs are designed to match the **Bounded Context** principle, where services only reference each other's entities via simple string IDs (`UserID`, `CourseID`, etc.) instead of direct database foreign keys.

---

## 1. Authentication Service (`auth-service`)

```go
// UserAccount represents the user's login identity

```
---

## 2. Student Profile Service (`profile-service`)

See the `domain` folder in `profile-service`

---

## 3. Course Catalogue Service (`course-catalogue-service`)

See the `domain` folder in `course-catalogue-service`

---

## 4. Course Content Service (`course-content-service`)

See the `domain` folder in `course-content-service`

---

## 5. Assignment & Grading Service (`assignment-service`)

```go

// Assignment represents an assignment with a deadline
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

// Submission is a student's submission of an assignment
type Submission struct {
	ID           string    `gorm:"primaryKey" json:"id"`
	AssignmentID string    `gorm:"index;not null" json:"assignment_id"`
	StudentID    string    `gorm:"index;not null" json:"student_id"` // Refers to UserProfile.ID
	FileURL      string    `gorm:"not null" json:"file_url"`        // Link to the submitted file
	SubmittedAt  time.Time `json:"submitted_at"`
	Grade        *Grade    `gorm:"foreignKey:SubmissionID" json:"grade,omitempty"`
}

// Grade contains the awarded assessment and feedback from the instructor
type Grade struct {
	ID           string    `gorm:"primaryKey" json:"id"`
	SubmissionID string    `gorm:"uniqueIndex;not null" json:"submission_id"`
	GraderID     string    `gorm:"not null" json:"grader_id"` // The instructor's/TA's UserID
	Score        int       `gorm:"not null" json:"score"`     // e.g. 85.5 out of 100
	LetterGrade  string    `json:"letter_grade,omitempty"`    // e.g. "A", "12", "Pass"
	Feedback     string    `json:"feedback,omitempty"`
	GradedAt     time.Time `json:"graded_at"`
}

```

---

## 6. Notification Service (`notification-service`)

See the `domain` folder in `notification-service`
