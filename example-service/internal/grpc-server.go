package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	// Replace with your actual generated proto package import path
	pb "github.com/yourusername/shared/proto/student"
)

// Server implements the generated pb.StudentServiceServer interface
type server struct {
	pb.UnimplementedStudentServiceServer
}

func (s *server) GetStudentProfile(ctx context.Context, req *pb.StudentRequest) (*pb.StudentResponse, error) {
	log.Printf("Received gRPC request for Student ID: %s", req.GetStudentId())

	// Mock database lookup
	return &pb.StudentResponse{
		StudentId: req.GetStudentId(),
		FirstName: "John",
		LastName:  "Doe",
		Role:      "student",
	}, nil
}

func main2() {
	// 1. Open a raw TCP network listener
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// 2. Initialize the gRPC network server instance
	grpcServer := grpc.NewServer()

	// 3. Register our implementation code with the gRPC server
	pb.RegisterStudentServiceServer(grpcServer, &server{})

	log.Println("gRPC Server listening on port :50051...")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
