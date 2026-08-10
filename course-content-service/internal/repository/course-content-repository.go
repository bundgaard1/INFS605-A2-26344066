package repository

import (
	"context"

	"osbourne.local/course-content-service/internal/domain"
)

type CourseContentRepository interface {
	ListModules(ctx context.Context, courseID string) ([]*domain.Module, int32, error)

	GetModule(ctx context.Context, moduleID string) (*domain.Module, error)
	CreateModule(ctx context.Context, module *domain.Module) error
	UpdateModule(ctx context.Context, module *domain.Module) error
	DeleteModule(ctx context.Context, moduleID string) error

	AddFileToModule(ctx context.Context, moduleID string, file *domain.File) error
	RemoveFileFromModule(ctx context.Context, moduleID string, fileID string) error
}
