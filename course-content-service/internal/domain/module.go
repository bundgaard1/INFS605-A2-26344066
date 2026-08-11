package domain

import (
	"time"
)

type Attachment struct {
	ID        string    `json:"id" clover:"id"`
	Name      string    `json:"name" clover:"name"`
	Type      string    `json:"type" clover:"type"` // f.eks. "pdf", "image"
	Path      string    `json:"path" clover:"path"` // Sti til filen
	Size      int64     `json:"size" clover:"size"` // Størrelse i bytes
	CreatedAt time.Time `json:"created_at" clover:"created_at"`
}

type Module struct {
	ID          string       `json:"id" clover:"id"`             // Unik identifikator for modulet
	CourseID    string       `json:"courseId" clover:"courseId"` // ref Course Catalogue
	Title       string       `json:"title" clover:"title"`
	Text        string       `json:"text" clover:"text"`
	Attachments []Attachment `json:"attachments" clover:"attachments"`
	UpdatedAt   time.Time    `json:"updatedAt" clover:"updatedAt"`
}
