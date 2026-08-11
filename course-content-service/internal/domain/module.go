package domain

import (
	"time"
)

// type Attachment struct {
// 	ID        string    `json:"id" clover:"id"`
// 	Name      string    `json:"name" clover:"name"`
// 	Type      string    `json:"type" clover:"type"`
// 	Path      string    `json:"path" clover:"path"`
// 	Size      int64     `json:"size" clover:"size"`
// 	CreatedAt time.Time `json:"created_at" clover:"created_at"`
// }

type Module struct {
	ID       string `json:"id" clover:"id"`
	CourseID string `json:"courseId" clover:"courseId"`
	Title    string `json:"title" clover:"title"`
	Text     string `json:"text" clover:"text"`
	// Attachments []Attachment `json:"attachments" clover:"attachments"`
	UpdatedAt time.Time `json:"updatedAt" clover:"updatedAt"`
}
