package domain

import "time"

type File struct {
	ID   string `bson:"id" json:"id"` // UUID
	Name string `bson:"name" json:"name"`
	Type string `bson:"type" json:"type"` // f.eks. "pdf", "image"
	URL  string `bson:"url" json:"url"`   // URL til filen
	Size int64  `bson:"size" json:"size"` // Størrelse i bytes
}

// Module er hoveddokumentet i MongoDB/NoSQL samlingen "modules"
type Module struct {
	ID       string `bson:"_id" json:"id"`              // UUID
	CourseID string `bson:"course_id" json:"course_id"` // ref Course Catalogue
	Title    string `bson:"title" json:"title"`
	Text     string `bson:"text" json:"text"` // Beskrivelse af modulet

	Files []File `bson:"files" json:"files"`

	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}
