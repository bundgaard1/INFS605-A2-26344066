package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"osbourne.local/users-service/gen/student"
)

type server struct {
	student.UnimplementedStudentServiceServer
}

func (s *server) GetStudentProfile(ctx context.Context, req *student.StudentRequest) (*student.StudentResponse, error) {
	log.Printf("Modtog forespørgsel for student_id: %s", req.GetStudentId())

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
	// Dynamisk port med fallback
	port := os.Getenv("PORT")
	if port == "" {
		port = "50051"
	}

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Kunne ikke lytte på port :%s: %v", port, err)
	}

	grpcServer := grpc.NewServer()
	student.RegisterStudentServiceServer(grpcServer, &server{})

	// 1. Kør serveren i en baggrunds-goroutine
	go func() {
		log.Printf("Users-Service (gRPC) kører på port :%s...", port)
		if err := grpcServer.Serve(lis); err != nil && err != grpc.ErrServerStopped {
			log.Fatalf("Fejl under afvikling af gRPC server: %v", err)
		}
	}()

	// 2. Lyt efter nedluknings-signaler (Ctrl+C, Docker stop m.m.)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Blokér indtil signalet modtages
	<-stop
	log.Println("Modtog nedlukningssignal. Lukker gRPC server pænt...")

	// 3. Graceful shutdown: Vent på at aktive gRPC-kald afsluttes
	done := make(chan struct{})
	go func() {
		grpcServer.GracefulStop()
		close(done)
	}()

	// Tving nedlukning hvis det tager mere end 5 sekunder
	select {
	case <-done:
		log.Println("gRPC server lukket pænt.")
	case <-time.After(5 * time.Second):
		log.Println("Tidsfrist overskredet - tvinger nedlukning.")
		grpcServer.Stop()
	}
}
