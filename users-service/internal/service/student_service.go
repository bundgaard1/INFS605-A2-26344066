package service

import (
	"context"

	"osbourne.local/users-service/gen/student"
	"osbourne.local/users-service/internal/repository"
)

type StudentService struct {
	repo repository.StudentRepository
}

// Injecter repositoriet
func NewStudentService(repo repository.StudentRepository) *StudentService {
	return &StudentService{repo: repo}
}

func (s *StudentService) GetProfile(ctx context.Context, id string) (*student.StudentResponse, error) {
	// Kald databaselaget
	return s.repo.GetByID(ctx, id)
}
