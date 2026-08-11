package domain

import (
	"context"
)

// type ObjectStore interface {
// 	Save(ctx context.Context, relativePath string, src io.Reader) (string, error)
// 	Get(ctx context.Context, relativePath string) (io.ReadCloser, error)
// 	Delete(ctx context.Context, relativePath string) error
// }

type ContentRepository interface {
	ListModules(ctx context.Context, courseID string) ([]*Module, error)
	GetModule(ctx context.Context, moduleID string) (*Module, error)
	CreateModule(ctx context.Context, module *Module) error
	UpdateModule(ctx context.Context, module *Module) error
	DeleteModule(ctx context.Context, moduleID string) error

	// AddAttachmentToModule(ctx context.Context, moduleID string, attachment *Attachment) error
	// RemoveAttachmentFromModule(ctx context.Context, moduleID string, attachmentID string) error
}
