package database

import (
	"fmt"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"osbourne.local/assignment-service/internal/domain"
)

func NewGORMDB(dbPath string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("Could not initialize GORM: %w", err)
	}

	// Enable WAL via a direct PRAGMA instead of a query string
	db.Exec("PRAGMA journal_mode=WAL;")

	for _, model := range []any{
		&domain.Assignment{},
		&domain.Submission{},
		&domain.Grade{},
	} {
		if err := db.AutoMigrate(model); err != nil {
			return nil, fmt.Errorf("Could not migrate database: %w", err)
		}
	}

	return db, nil
}
