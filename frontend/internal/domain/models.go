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

type Course struct {
	ID          string
	Code        string
	Title       string
	Description string
	Credits     int
}

type PageData struct {
	User User
	Data any
}
