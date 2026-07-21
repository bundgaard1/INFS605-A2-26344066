package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	"osbourne.local/users-service/gen/student"
)

// server-structen implementerer det genererede StudentServiceServer interface
type server struct {
	student.UnimplementedStudentServiceServer
}

// GetStudentProfile håndterer gRPC-kaldet og returnerer en sample user
func (s *server) GetStudentProfile(ctx context.Context, req *student.StudentRequest) (*student.StudentResponse, error) {
	log.Printf("Modtog forespørgsel for student_id: %s", req.GetStudentId())

	// Returner en simpel sample user
	return &student.StudentResponse{
		Id:   req.GetStudentId(),
		Name: "Andy Osborne",
		Role: "Student",
		Courses: []*student.Course{
			{Id: "INFS-605", Title: "Microservices Programming Project"},
			{Id: "COMP-901", Title: "Distributed Systems"},
		},
	}, nil
}

func main() {
	// 1. Åbn en TCP-port til gRPC-trafik
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Kunne ikke lytte på port :50051: %v", err)
	}

	// 2. Opret en ny gRPC server
	grpcServer := grpc.NewServer()

	// 3. Registrer vores server-implementation
	student.RegisterStudentServiceServer(grpcServer, &server{})

	log.Println("Users-Service (gRPC) kører på port :50051...")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Fejl under afvikling af gRPC server: %v", err)
	}
}
