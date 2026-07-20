package main

import (
	"context"
	"log"
	"time"

	pb "osbourne.local/shared/proto/student"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main1() {
	// 1. Establish a network connection to the server (Insecure for local Dev/Docker mesh)
	// In production docker compose, "localhost:50051" would be "student-service:50051"
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// 2. Create the type-safe client stub
	client := pb.NewStudentServiceClient(conn)

	// 3. Execute the network call with a timeout context
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := client.GetStudentProfile(ctx, &pb.StudentRequest{StudentId: "12345"})
	if err != nil {
		log.Fatalf("could not fetch profile: %v", err)
	}

	log.Printf("gRPC Response Received: %s %s (%s)", res.GetFirstName(), res.GetLastName(), res.GetRole())
}
