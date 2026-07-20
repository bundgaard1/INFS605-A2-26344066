package main

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main3() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	// Ensure the queue exists here as well, in case consumer boots before publisher
	q, err := ch.QueueDeclare("enrollment_queue", false, false, false, false, nil)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	// 1. Register a consumer network stream to read from the queue
	msgs, err := ch.Consume(
		q.Name, // Queue name
		"",     // Consumer identifier tags
		true,   // Auto-Acknowledge (Removes message from queue immediately on read)
		false,  // Exclusive
		false,  // No-local
		false,  // No-wait
		nil,    // Args
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	// 2. Read indefinitely from the Go channel stream using a blocking loop
	var forever chan struct{}

	go func() {
		for d := range msgs {
			log.Printf("Async Event Received: %s", d.Body)
			// This is where your code triggers emails, logs, or metrics collections.
		}
	}()

	log.Println(" [*] Notification worker waiting for events. To exit press CTRL+C")
	<-forever // Blocks the main thread permanently so the goroutine keeps listening
}
