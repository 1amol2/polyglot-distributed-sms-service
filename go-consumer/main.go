package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/segmentio/kafka-go"
)

type SmsEvent struct {
	MessageID   string          `json:"messageId"`
	PhoneNumber string          `json:"phoneNumber"`
	Message     string          `json:"message"`
	Status      string          `json:"status"`
	CreatedAt   json.RawMessage `json:"createdAt"`
}

func main() {
	broker := os.Getenv("KAFKA_BROKERS")
	if broker == "" {
		broker = "localhost:9092"
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{broker},
		Topic:       "sms-events",
		GroupID:     "polyglot-go-consumer",
		StartOffset: kafka.FirstOffset,
	})

	defer func() {
		if err := reader.Close(); err != nil {
			log.Printf("Error closing Kafka reader: %v", err)
		}
	}()

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	log.Printf("Go consumer listening to sms-events at %s", broker)

	for {
		message, err := reader.FetchMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				log.Println("Shutting down consumer")
				return
			}

			log.Printf("Failed to fetch Kafka message: %v", err)
			continue
		}

		var event SmsEvent
		if err := json.Unmarshal(message.Value, &event); err != nil {
			log.Printf("Invalid SMS event at offset %d: %v",
				message.Offset, err)
			continue
		}

		log.Printf(
			"SMS event received: ID=%s Phone=%s Status=%s Message=%q",
			event.MessageID,
			event.PhoneNumber,
			event.Status,
			event.Message,
		)

		if err := reader.CommitMessages(ctx, message); err != nil {
			log.Printf("Failed to commit Kafka offset: %v", err)
			continue
		}
	}
}
