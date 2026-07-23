package repository

import (
	"context"

	"osbourne.local/users-service/gen/student"
)

// StudentRepository definerer alle database-operationer for studerende
type StudentRepository interface {
	GetByID(ctx context.Context, id string) (*student.StudentResponse, error)
}
