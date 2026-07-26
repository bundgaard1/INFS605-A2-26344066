package database

import (
	"fmt"

	// Brug den officielle GORM dialekt-adapter fra glebarez:

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"osbourne.local/notification-service/internal/domain"
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

	err = db.AutoMigrate(&domain.Notification{})
	if err != nil {
		return nil, fmt.Errorf("fejl ved AutoMigrate: %w", err)
	}

	return db, nil
}

func SeedData(db *gorm.DB) {
	var count int64
	db.Model(&domain.Notification{}).Count(&count)
	if count > 0 {
		return // Data findes allerede
	}

	notification := domain.Notification{
		ID:      "1",
		UserID:  "12345",
		Title:   "Ny karakter modtaget",
		Message: "Du har fået 12 i INFS-605, link tager dig til root",
		LinkURL: "/",
		IsRead:  false,
	}

	notification2 := domain.Notification{
		ID:      "2",
		UserID:  "12345",
		Title:   "Velkommen til platformen",
		Message: "Tak for at tilmelde dig vores platform!",
		LinkURL: "/profile",
		IsRead:  true,
	}

	db.Create(&notification)
	db.Create(&notification2)
}
