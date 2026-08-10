Her er den samlede Go markdown-oversigt over **domænestructs (GORM/Entities)** for alle 6 mikrotjenester.

Strukturerne er designet således, at de matcher **Bounded Context** princippet, hvor services kun refererer til hinandens entiteter via simple streng-IDs (`UserID`, `CourseID`, osv.) i stedet for direkte database-foreign-keys.

---

## 1. Authentication Service (`auth-service`)

```go
// UserAccount repræsenterer brugerens login-identitet
type UserAccount struct {
	ID           string     `gorm:"primaryKey" json:"id"`                  
	Email        string     `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash string     `gorm:"not null" json:"-"`
	IsActive     bool       `gorm:"default:true" json:"is_active"`
	LastLoginAt  *time.Time `json:"last_login_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

```
---

## 2. Student Profile Service (`profile-service`)

See the `domain` folder in `profile-service`

---

## 3. Course Catalogue Service (`course-catalogue-service`)

See the `domain`folder in `course-catalogue-service`

---

## 4. Course Content Service (`course-content-service`)

See the `domain` folder in `course-content-service`



---

## 5. Assignment & Grading Service (`assignment-service`)

```go

// Assignment repræsenterer en opgave med deadline
type Assignment struct {
	ID          string       `gorm:"primaryKey" json:"id"`
	CourseID    string       `gorm:"index;not null" json:"course_id"` // Refererer til Course.ID
	Title       string       `gorm:"not null" json:"title"`
	Description string       `json:"description"`
	DueDate     time.Time    `gorm:"not null" json:"due_date"`
	Submissions []Submission `gorm:"foreignKey:AssignmentID" json:"submissions,omitempty"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

// Submission er en studerendes aflevering af en opgave
type Submission struct {
	ID           string    `gorm:"primaryKey" json:"id"`
	AssignmentID string    `gorm:"index;not null" json:"assignment_id"`
	StudentID    string    `gorm:"index;not null" json:"student_id"` // Refererer til UserProfile.ID
	FileURL      string    `gorm:"not null" json:"file_url"`        // Link til den afleverede fil
	SubmittedAt  time.Time `json:"submitted_at"`
	Grade        *Grade    `gorm:"foreignKey:SubmissionID" json:"grade,omitempty"`
}

// Grade indeholder den tildelte bedømmelse og feedback fra underviseren
type Grade struct {
	ID           string    `gorm:"primaryKey" json:"id"`
	SubmissionID string    `gorm:"uniqueIndex;not null" json:"submission_id"`
	GraderID     string    `gorm:"not null" json:"grader_id"` // Underviserens/TA'ens UserID
	Score        int       `gorm:"not null" json:"score"`     // f.eks. 85.5 ud af 100
	LetterGrade  string    `json:"letter_grade,omitempty"`    // f.eks. "A", "12", "Pass"
	Feedback     string    `json:"feedback,omitempty"`
	GradedAt     time.Time `json:"graded_at"`
}

```

---

## 6. Notification Service (`notification-service`)

```go

type Notification struct {
	ID        string    `gorm:"primaryKey" json:"id"`
	UserID    string    `gorm:"index;not null" json:"user_id"` // Hvem notifikationen tilhører
	Title     string    `gorm:"not null" json:"title"`         // f.eks. "Ny karakter modtaget"
	Message   string    `gorm:"not null" json:"message"`       // f.eks. "Du har fået 12 i INFS-605"
	LinkURL   string    `json:"link_url,omitempty"`            // f.eks. "/assignments/asg_9921" (så brugeren kan klikke)
	IsRead    bool      `gorm:"default:false;index" json:"is_read"`
	CreatedAt time.Time `gorm:"index" json:"created_at"`
}

```