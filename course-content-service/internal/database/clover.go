package database

import (
	"fmt"

	"github.com/ostafen/clover/v2"
	c "github.com/ostafen/clover/v2"
	"github.com/ostafen/clover/v2/document"
	"osbourne.local/course-content-service/internal/domain"
)

func NewCloverDB(connectionString string) (*clover.DB, error) {
	db, err := clover.Open(connectionString)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func SeedCloverData(db *c.DB) error {
	if ok, err := db.HasCollection("modules"); ok {
		return nil
	} else if err != nil {
		return fmt.Errorf("Could not check collection: %w", err)
	}

	err := db.CreateCollection("modules")
	if err != nil {
		return fmt.Errorf("Could not create collection: %w", err)
	}

	module := domain.Module{
		ID:          "1",
		CourseID:    "1",
		Title:       "Module 1",
		Text:        "This is the first module.",
		Attachments: []domain.Attachment{},
	}

	doc := document.NewDocumentOf(module)

	if doc == nil {
		return fmt.Errorf("Could not create document")
	}

	docId, err := db.InsertOne("modules", doc)
	if err != nil {
		return fmt.Errorf("Could not insert document: %w", err)
	}
	fmt.Println(docId)

	return nil
}
