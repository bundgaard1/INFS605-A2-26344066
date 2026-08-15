package domain

import "time"

type User struct {
	ID   string
	Name string
	Role string
}

type Notification struct {
	ID        string
	Title     string
	Message   string
	Read      bool
	Link      string
	Timestamp time.Time
}

type Module struct {
	ID       string
	CourseID string
	Title    string
	Text     string
}

type Course struct {
	ID          string
	Code        string
	Title       string
	Description string
	Credits     int
	Modules     []Module
}
