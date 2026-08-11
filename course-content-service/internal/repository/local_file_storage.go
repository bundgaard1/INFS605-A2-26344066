package repository

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type LocalFileStorage struct {
	baseDir string
}

func NewLocalFileStorage(baseDir string) (*LocalFileStorage, error) {
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create upload base directory: %w", err)
	}
	return &LocalFileStorage{baseDir: baseDir}, nil
}

func (s *LocalFileStorage) Save(ctx context.Context, relativePath string, src io.Reader) (string, error) {
	fullPath := filepath.Join(s.baseDir, relativePath)

	// Ensure subdirectories exist (e.g. ./uploads/modules/mod_123/)
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return "", fmt.Errorf("failed to create target directory: %w", err)
	}

	dst, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to create local file: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return "", fmt.Errorf("failed to write file contents: %w", err)
	}

	return relativePath, nil
}

func (s *LocalFileStorage) Get(ctx context.Context, relativePath string) (io.ReadCloser, error) {
	fullPath := filepath.Join(s.baseDir, relativePath)
	file, err := os.Open(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	return file, nil
}

func (s *LocalFileStorage) Delete(ctx context.Context, relativePath string) error {
	fullPath := filepath.Join(s.baseDir, relativePath)
	err := os.Remove(fullPath)
	if os.IsNotExist(err) {
		return nil // File already gone
	}
	return err
}
