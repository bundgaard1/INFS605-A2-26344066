package database

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	// Brug den officielle GORM dialekt-adapter fra glebarez:

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"osbourne.local/profile-service/internal/domain"
)

func NewGORMDB(dbPath string) (*gorm.DB, error) {
	// Ensure folder exists
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		log.Fatalf("Failed to create DB directory: %v", err)
	}

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("kunne ikke forbinde til SQLite via GORM: %w", err)
	}

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
		return
	}

	student := domain.UserProfile{
		ID:   "12345",
		Name: "Andy Osborne",
		Role: "Student",
	}

	db.Create(&student)
}
