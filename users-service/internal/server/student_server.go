package server

import (
	"context"
	"log"

	"osbourne.local/users-service/gen/student"
	"osbourne.local/users-service/internal/service"
)

type StudentServer struct {
	student.UnimplementedStudentServiceServer
	studentSvc *service.StudentService
}

func NewStudentServer(studentSvc *service.StudentService) *StudentServer {
	return &StudentServer{
		studentSvc: studentSvc,
	}
}

func (s *StudentServer) GetStudentProfile(ctx context.Context, req *student.StudentRequest) (*student.StudentResponse, error) {
	log.Printf("Modtog gRPC-forespørgsel for student_id: %s", req.GetStudentId())

	return s.studentSvc.GetProfile(ctx, req.GetStudentId())
}
