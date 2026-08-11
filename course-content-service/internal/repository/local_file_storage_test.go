package repository_test

import (
	"context"
	"os"
	"strings"
	"testing"

	repo "osbourne.local/course-content-service/internal/repository"
)

func setupLocalFileStorage(t *testing.T) (*repo.LocalFileStorage, func()) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	// Initialize the LocalFileStorage with the temporary directory
	fileStorage, err := repo.NewLocalFileStorage(tempDir)
	if err != nil {
		t.Fatalf("failed to initialize LocalFileStorage: %v", err)
	}

	// Return a cleanup function to remove the temporary directory after the test
	cleanup := func() {
		os.RemoveAll(tempDir)
	}

	return fileStorage, cleanup
}

func TestLocalFileStorage_Save(t *testing.T) {
	ctx := context.Background()
	// Initialize the LocalFileStorage with the temporary directory
	fs, cleanup := setupLocalFileStorage(t)
	defer cleanup()

	// Define test data
	relativePath := "test_dir/test_file.txt"
	content := "Hello, World!"

	// Call the Save method
	relativePath, err := fs.Save(ctx, relativePath, strings.NewReader(content))
	if err != nil {
		t.Fatalf("Save method failed: %v", err)
	}

	reader, err := fs.Get(ctx, relativePath)
	if err != nil {
		t.Fatalf("Get method failed: %v", err)
	}
	defer reader.Close()

	// Read the content back from the file
	readContent := make([]byte, len(content))
	_, err = reader.Read(readContent)
	if err != nil {
		t.Fatalf("Failed to read content from file: %v", err)
	}

	if string(readContent) != content {
		t.Errorf("Content mismatch: expected %q, got %q", content, string(readContent))
	}
}

func TestLocalFileStorage_Delete(t *testing.T) {
	ctx := context.Background()
	// Initialize the LocalFileStorage with the temporary directory
	fs, cleanup := setupLocalFileStorage(t)
	defer cleanup()

	// Define test data
	relativePath := "test_dir/test_file_to_delete.txt"
	content := "This file will be deleted."

	// Save the file first
	_, err := fs.Save(ctx, relativePath, strings.NewReader(content))
	if err != nil {
		t.Fatalf("Save method failed: %v", err)
	}

	// Get the file to ensure it exists
	reader, err := fs.Get(ctx, relativePath)
	if err != nil {
		t.Fatalf("Get method failed: %v", err)
	}
	reader.Close()

	// Now delete the file
	err = fs.Delete(ctx, relativePath)
	if err != nil {
		t.Fatalf("Delete method failed: %v", err)
	}

	// Try to get the file after deletion, it should return an error
	_, err = fs.Get(ctx, relativePath)
	if err == nil {
		t.Errorf("Expected error when getting deleted file, but got none")
	}
}
