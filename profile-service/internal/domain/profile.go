package domain

import "time"

type UserRole string

const (
	RoleStudent UserRole = "Student"
	RoleTeacher UserRole = "Teacher"
)

// UserProfile indeholder brugerens personlige stamdata
type UserProfile struct {
	ID        string    `gorm:"primaryKey" json:"id"` // Samme ID som UserAccount.ID
	Name      string    `gorm:"not null" json:"name"`
	Role      UserRole  `gorm:"type:string;default:'Student';not null" json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
