package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"osbourne.local/profile-service/gen/profile"
	"osbourne.local/profile-service/internal/database"
	"osbourne.local/profile-service/internal/repository"
	"osbourne.local/profile-service/internal/server"
	"osbourne.local/profile-service/internal/service"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "50051"
	}

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Could not listen on port :%s: %v", port, err)
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./profiles.db"
	}

	fmt.Println("Database path:", dbPath) // Debugging line to check the DB_PATH value

	db, err := database.NewGORMDB(dbPath)
	if err != nil {
		log.Fatalf("Database error: %v", err)
	}

	database.SeedData(db)

	profileRepo := repository.NewGORMProfileRepository(db)
	profileSvc := service.NewProfileService(profileRepo)
	profileGrpcServer := server.NewProfileServer(profileSvc)

	grpcServer := grpc.NewServer()
	profile.RegisterProfileServiceServer(grpcServer, profileGrpcServer)

	go func() {
		log.Printf("profile-service (gRPC) running on port :%s...", port)
		if err := grpcServer.Serve(lis); err != nil && err != grpc.ErrServerStopped {
			log.Fatalf("Error while running gRPC server: %v", err)
		}
	}()

	// 4. Graceful Shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("Received shutdown signal. Shutting down gracefully...")

	done := make(chan struct{})
	go func() {
		grpcServer.GracefulStop()
		close(done)
	}()

	select {
	case <-done:
		log.Println("gRPC server shut down gracefully.")
	case <-time.After(5 * time.Second):
		log.Println("Timeout exceeded - forcing shutdown.")
		grpcServer.Stop()
	}
}
