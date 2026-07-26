package domain

import "time"

type Enrollment struct {
	ID         string    `gorm:"primaryKey" json:"id"`
	CourseID   string    `gorm:"index;not null;compositeIndex:idx_student_course,unique" json:"course_id"`
	UserID     string    `gorm:"index;not null;compositeIndex:idx_student_course,unique" json:"user_id"` // Refererer til
	EnrolledAt time.Time `json:"enrolled_at"`
}
