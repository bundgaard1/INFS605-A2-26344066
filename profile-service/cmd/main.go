package main

import (
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
		log.Fatalf("Kunne ikke lytte på port :%s: %v", port, err)
	}

	db, err := database.NewGORMDB("./users.db")
	if err != nil {
		log.Fatalf("DB fejl: %v", err)
	}

	database.SeedData(db)

	profileRepo := repository.NewGORMProfileRepository(db)
	profileSvc := service.NewProfileService(profileRepo)
	profileGrpcServer := server.NewProfileServer(profileSvc)

	grpcServer := grpc.NewServer()
	profile.RegisterProfileServiceServer(grpcServer, profileGrpcServer)

	go func() {
		log.Printf("profile-service (gRPC) kører på port :%s...", port)
		if err := grpcServer.Serve(lis); err != nil && err != grpc.ErrServerStopped {
			log.Fatalf("Fejl under afvikling af gRPC server: %v", err)
		}
	}()

	// 4. Graceful Shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("Modtog nedlukningssignal. Lukker pænt...")

	done := make(chan struct{})
	go func() {
		grpcServer.GracefulStop()
		close(done)
	}()

	select {
	case <-done:
		log.Println("gRPC server lukket pænt.")
	case <-time.After(5 * time.Second):
		log.Println("Tidsfrist overskredet - tvinger nedlukning.")
		grpcServer.Stop()
	}
}
