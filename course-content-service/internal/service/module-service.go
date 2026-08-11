package service

import "osbourne.local/course-content-service/internal/domain"

type ModuleService struct {
	repo      domain.ContentRepository
	fileStore domain.ObjectStore
}

func NewModuleService(repo domain.ContentRepository, fileStore domain.ObjectStore) *ModuleService {
	return &ModuleService{
		repo:      repo,
		fileStore: fileStore,
	}
}
