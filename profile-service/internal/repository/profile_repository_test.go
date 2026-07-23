package repository_test

import (
	"context"
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"osbourne.local/profile-service/internal/domain"
	"osbourne.local/profile-service/internal/repository"
)

// setupTestDB opretter en helt ren in-memory SQLite database til hver test
func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	// "file::memory:?cache=shared" eller blot ":memory:" sikrer at DB kun lever i RAM under testen
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // Skjul SQL-logs under testkørsel
	})
	if err != nil {
		t.Fatalf("Kunne ikke oprette in-memory database: %v", err)
	}

	// Kør auto-migrering på in-memory DB'en
	err = db.AutoMigrate(&domain.UserProfile{})
	if err != nil {
		t.Fatalf("Fejl ved AutoMigrate i test: %v", err)
	}

	return db
}

func TestGORMProfileRepository_GetByID(t *testing.T) {
	// 1. Arrange: Klargør test-data i in-memory databasen
	db := setupTestDB(t)
	repo := repository.NewGORMProfileRepository(db)
	ctx := context.Background()

	testStudent := domain.UserProfile{
		ID:   "student-123",
		Name: "student",
		Role: domain.RoleStudent,
	}

	if err := db.Create(&testStudent).Error; err != nil {
		t.Fatalf("Kunne ikke indsætte test-data: %v", err)
	}

	// 2. Act: Kald metoden på repositoriet
	result, err := repo.GetByID(ctx, "student-123")

	// 3. Assert: Kontrollér at resultatet er som forventet
	if err != nil {
		t.Fatalf("Forventede ingen fejl, men fik: %v", err)
	}

	if result.Id != testStudent.ID {
		t.Errorf("Forventede ID %s, men fik %s", testStudent.ID, result.Id)
	}

	if result.Name != testStudent.Name {
		t.Errorf("Forventede Fornavn %s, men fik %s", testStudent.Name, result.Name)
	}

}

func TestGORMProfileRepository_GetByID_NotFound(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	repo := repository.NewGORMProfileRepository(db)
	ctx := context.Background()

	// Act: Søg efter et ID der ikke findes
	_, err := repo.GetByID(ctx, "non-existing-id")

	// Assert: Forvent en fejl
	if err == nil {
		t.Error("Forventede en fejl for et ugyldigt ID, men fik nil")
	}
}
