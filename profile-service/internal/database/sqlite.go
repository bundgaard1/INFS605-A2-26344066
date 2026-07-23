package database

import (
	"fmt"

	// Brug den officielle GORM dialekt-adapter fra glebarez:

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"osbourne.local/profile-service/internal/domain"
)

func NewGORMDB(dbPath string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("kunne ikke forbinde til SQLite via GORM: %w", err)
	}

	// Aktiver WAL via en direkte PRAGMA i stedet for query string
	db.Exec("PRAGMA journal_mode=WAL;")

	err = db.AutoMigrate(&domain.UserProfile{})
	if err != nil {
		return nil, fmt.Errorf("fejl ved AutoMigrate: %w", err)
	}

	return db, nil
}

func SeedData(db *gorm.DB) {
	var count int64
	db.Model(&domain.UserProfile{}).Count(&count)
	if count > 0 {
		return // Data findes allerede
	}

	student := domain.UserProfile{
		ID:   "12345",
		Name: "Andy Osborne",
		Role: "Student",
	}

	db.Create(&student)
}
