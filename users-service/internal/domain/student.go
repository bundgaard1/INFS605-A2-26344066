package domain

type Student struct {
	ID      string   `gorm:"primaryKey"`
	Name    string   `gorm:"not null"`
	Role    string   `gorm:"not null"`
	Courses []Course `gorm:"many2many:student_courses;"` // M2M relation
}

type Course struct {
	ID    string `gorm:"primaryKey"`
	Title string `gorm:"not null"`
}
