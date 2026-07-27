package database

import (
	"fmt"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"osbourne.local/course-catalogue-service/internal/domain"
)

func NewGORMDB(dbPath string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("Could not initialize GORM: %w", err)
	}

	// Aktiver WAL via en direkte PRAGMA i stedet for query string
	db.Exec("PRAGMA journal_mode=WAL;")

	err = db.AutoMigrate(&domain.Course{})
	if err != nil {
		return nil, fmt.Errorf("Could not migrate database: %w", err)
	}

	err = db.AutoMigrate(&domain.Enrollment{})
	if err != nil {
		return nil, fmt.Errorf("Could not migrate database: %w", err)
	}

	return db, nil
}

func SeedData(db *gorm.DB) {
	var count int64
	db.Model(&domain.Course{}).Count(&count)
	if count > 0 {
		return
	}

	courses := []domain.Course{
		{
			ID:          "1",
			Code:        "CS101",
			Title:       "Introduction to Computer Science",
			Description: "Learn the basics of computer science.",
			Credits:     10,
		},
		{
			ID:          "2",
			Code:        "CS102",
			Title:       "Data Structures and Algorithms",
			Description: "Learn about data structures and algorithms.",
			Credits:     10,
		},
	}
	db.CreateInBatches(&courses, 10)

	enrollments := []domain.Enrollment{
		{
			ID:       "3",
			CourseID: "1",
			UserID:   "12345",
		},
		{
			ID:       "4",
			CourseID: "2",
			UserID:   "12345",
		},
	}
	db.CreateInBatches(&enrollments, 10)
}
