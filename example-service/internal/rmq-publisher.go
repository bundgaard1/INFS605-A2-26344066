package main

import (
	"context"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main4() {
	// 1. Dial the RabbitMQ TCP network endpoint
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	// 2. Open a virtual channel inside the TCP connection
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	// 3. Declare the queue to ensure it exists before publishing
	q, err := ch.QueueDeclare(
		"enrollment_queue", // Queue name
		false,              // Durable (survives broker restart)
		false,              // Delete when unused
		false,              // Exclusive
		false,              // No-wait
		nil,                // Arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	// 4. Publish the raw byte payload asynchronously
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body := []byte(`{"student_id":"12345", "course_id":"CS-401", "action":"enrolled"}`)
	err = ch.PublishWithContext(ctx,
		"",     // Exchange name (using default direct exchange)
		q.Name, // Routing key (matches queue name)
		false,  // Mandatory
		false,  // Immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		log.Fatalf("Failed to publish a message: %v", err)
	}
	log.Println("Successfully emitted Asynchronous Event to RabbitMQ")
}
