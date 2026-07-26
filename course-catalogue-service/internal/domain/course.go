package domain

import "time"

type Course struct {
	ID          string    `gorm:"primaryKey" json:"id"`             // f.eks. "261605" (real DB ID)
	Code        string    `gorm:"uniqueIndex;not null" json:"code"` // f.eks. "INFS605"
	Title       string    `gorm:"not null" json:"title"`
	Description string    `json:"description"`
	Credits     int       `gorm:"not null" json:"credits"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
