package database_test

import (
	"os"
	"testing"

	c "github.com/ostafen/clover/v2"
	"osbourne.local/course-content-service/internal/database"
)

func TestSeedCloverData_Idempotent(t *testing.T) {
	dir, err := os.MkdirTemp("", "clover-seed-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	db, err := c.Open(dir)
	if err != nil {
		t.Fatal(err)
	}

	if err := database.SeedCloverData(db); err != nil {
		t.Fatalf("first seed failed: %v", err)
	}
	if err := database.SeedCloverData(db); err != nil {
		t.Fatalf("second seed on same DB must not fail: %v", err)
	}

	if err := db.Close(); err != nil {
		t.Fatalf("close: %v", err)
	}

	db2, err := c.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer db2.Close()

	if err := database.SeedCloverData(db2); err != nil {
		t.Fatalf("seed after reopen must not fail: %v", err)
	}
}
