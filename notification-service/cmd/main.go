package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/wagslane/go-rabbitmq"
	"google.golang.org/grpc"
	"osbourne.local/notification-service/gen/notification"
	"osbourne.local/notification-service/internal/consumer"
	"osbourne.local/notification-service/internal/database"
	"osbourne.local/notification-service/internal/repository"
	"osbourne.local/notification-service/internal/server"
	"osbourne.local/notification-service/internal/service"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "50052"
	}

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Kunne ikke lytte på port :%s: %v", port, err)
	}

	db, err := database.NewGORMDB("./notifications.db")
	if err != nil {
		log.Fatalf("DB fejl: %v", err)
	}

	database.SeedData(db)

	notificationRepo := repository.NewGORMNotificationRepository(db)
	notificationSvc := service.NewNotificationService(notificationRepo)
	notificationGrpcServer := server.NewNotificationServer(notificationSvc)

	// Opret RabbitMQ forbindelse
	amqpURL := os.Getenv("RABBITMQ_URL")
	if amqpURL == "" {
		amqpURL = "amqp://guest:guest@rabbitmq:5672/"
	}
	conn, err := rabbitmq.NewConn(amqpURL)
	if err != nil {
		log.Fatalf("Fejl under oprettelse af RabbitMQ forbindelse: %v", err)
	}
	defer conn.Close()

	rmqConsumer, err := consumer.NewNotificationConsumer(conn, notificationSvc)
	if err != nil {
		log.Fatalf("Fejl under oprettelse af RabbitMQ Consumer: %v", err)
	}
	defer rmqConsumer.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		log.Println("Starter RabbitMQ Consumer worker...")
		if err := rmqConsumer.Start(ctx); err != nil {
			log.Printf("RabbitMQ Consumer stoppede med fejl: %v", err)
		}
	}()

	// Start gRPC server
	grpcServer := grpc.NewServer()
	notification.RegisterNotificationServiceServer(
		grpcServer,
		notificationGrpcServer)

	go func() {
		log.Printf("notification-service (gRPC) kører på port :%s...", port)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Fejl under afvikling af gRPC server: %v", err)
		}
	}()

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
