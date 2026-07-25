package consumer

import (
	"context"
	"log"

	"github.com/wagslane/go-rabbitmq"
	"google.golang.org/protobuf/proto"
	"osbourne.local/notification-service/gen/events"
	"osbourne.local/notification-service/internal/service"
)

type NotificationConsumer struct {
	rmq *rabbitmq.Consumer
	svc *service.NotificationService
}

func NewNotificationConsumer(conn *rabbitmq.Conn, svc *service.NotificationService) (*NotificationConsumer, error) {
	nc := &NotificationConsumer{
		svc: svc,
	}

	// 1. Opret RabbitMQ Consumer direkte med biblioteket
	consumer, err := rabbitmq.NewConsumer(
		conn,
		"notifications_queue",
	)

	if err != nil {
		return nil, err
	}

	nc.rmq = consumer
	return nc, nil
}

// Start åbner for lyttelsen på RabbitMQ beskeder
func (c *NotificationConsumer) Start(ctx context.Context) error {
	log.Println("[CONSUMER] NotificationConsumer lytter på RabbitMQ...")

	// Run modtager en handler-funktion, der kaldes for hver besked
	err := c.rmq.Run(func(d rabbitmq.Delivery) rabbitmq.Action {
		return c.processDelivery(ctx, d.Body)
	})

	return err
}

// Close lukker forbrugeren pænt ved shutdown
func (c *NotificationConsumer) Close() {
	if c.rmq != nil {
		c.rmq.Close()
	}
}

func (c *NotificationConsumer) processDelivery(ctx context.Context, body []byte) rabbitmq.Action {
	// 1. Unmarshal til fælles Protobuf Envelope
	var envelope events.EventEnvelope
	if err := proto.Unmarshal(body, &envelope); err != nil {
		log.Printf("[CONSUMER] Ugyldig EventEnvelope format: %v", err)
		// Smid beskeden væk (NackWithoutRequeue = Send ikke tilbage på køen)
		return rabbitmq.NackDiscard
	}

	// 2. Rute ud fra envelope.Type
	switch envelope.Type {
	case "student.created":
		var event events.StudentCreatedEvent
		if err := proto.Unmarshal(envelope.Payload, &event); err != nil {
			log.Printf("[CONSUMER] Fejl ved unmarshal af StudentCreatedEvent: %v", err)
			return rabbitmq.NackDiscard
		}

		// Kalder forretningslogik
		err := c.svc.CreateNotification(ctx,
			event.GetStudentId(),
			"Velkommen til platformen!",
			"Jeg er glad for at se dig på Osbourne. Vi håber, du får en god oplevelse her.",
			"/")

		if err != nil {
			log.Printf("[CONSUMER] Fejl ved gemning af notifikation: %v", err)
			return rabbitmq.NackRequeue
		}
	default:
		log.Printf("[CONSUMER] Ignorerer ukendt event-type: %s", envelope.Type)
	}

	// Alt gik godt -> Godkend beskeden
	return rabbitmq.Ack
}
