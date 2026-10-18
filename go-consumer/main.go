package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/segmentio/kafka-go"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type SmsEvent struct {
	MessageID   string          `json:"messageId" bson:"messageId"`
	PhoneNumber string          `json:"phoneNumber" bson:"phoneNumber"`
	Message     string          `json:"message" bson:"message"`
	Status      string          `json:"status" bson:"status"`
	CreatedAt   json.RawMessage `json:"createdAt" bson:"createdAt"`
	UserID      string          `json:"userId" bson:"userId"`
}
type DeadLetterEvent struct {
	OriginalTopic     string    `json:"originalTopic"`
	OriginalPartition int       `json:"originalPartition"`
	OriginalOffset    int64     `json:"originalOffset"`
	OriginalKey       string    `json:"originalKey"`
	OriginalValue     string    `json:"originalValue"`
	Reason            string    `json:"reason"`
	FailedAt          time.Time `json:"failedAt"`
}

func parseCreatedAt(raw json.RawMessage) (time.Time, error) {
	var seconds float64

	if err := json.Unmarshal(raw, &seconds); err == nil {
		wholeSeconds := int64(seconds)
		nanoseconds := int64((seconds - float64(wholeSeconds)) * 1e9)
		return time.Unix(wholeSeconds, nanoseconds).UTC(), nil
	}

	var timestamp string
	if err := json.Unmarshal(raw, &timestamp); err != nil {
		return time.Time{}, err
	}

	return time.Parse(time.RFC3339Nano, timestamp)
}

type MessagesFinder func(
	ctx context.Context,
	userID string,
) ([]bson.M, error)

func messagesHandler(findMessages MessagesFinder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.Header().Set("Allow", http.MethodGet)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		const prefix = "/v1/user/"
		const suffix = "/messages"

		path := strings.TrimPrefix(r.URL.Path, prefix)
		if !strings.HasPrefix(r.URL.Path, prefix) ||
			!strings.HasSuffix(path, suffix) {
			http.NotFound(w, r)
			return
		}

		userID := strings.TrimSuffix(path, suffix)
		if userID == "" || strings.Contains(userID, "/") {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		queryCtx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		messages, err := findMessages(queryCtx, userID)
		if err != nil {
			log.Printf("Query messages for user %s: %v", userID, err)
			http.Error(w, "Failed to retrieve messages",
				http.StatusInternalServerError)
			return
		}

		if messages == nil {
			messages = []bson.M{}
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(messages); err != nil {
			log.Printf("Write response for user %s: %v", userID, err)
		}
	}
}
func validateEvent(event SmsEvent) (time.Time, error) {
	if event.MessageID == "" {
		return time.Time{}, fmt.Errorf("messageId is required")
	}

	if event.UserID == "" {
		return time.Time{}, fmt.Errorf("userId is required")
	}

	if event.PhoneNumber == "" {
		return time.Time{}, fmt.Errorf("phoneNumber is required")
	}

	if event.Message == "" {
		return time.Time{}, fmt.Errorf("message is required")
	}

	if event.Status != "SUCCESS" && event.Status != "FAILURE" {
		return time.Time{}, fmt.Errorf("invalid status %q", event.Status)
	}

	createdAt, err := parseCreatedAt(event.CreatedAt)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid createdAt: %w", err)
	}

	return createdAt, nil
}
func publishToDeadLetter(
	ctx context.Context,
	writer *kafka.Writer,
	message kafka.Message,
	reason string,
) error {
	event := DeadLetterEvent{
		OriginalTopic:     message.Topic,
		OriginalPartition: message.Partition,
		OriginalOffset:    message.Offset,
		OriginalKey:       string(message.Key),
		OriginalValue:     string(message.Value),
		Reason:            reason,
		FailedAt:          time.Now().UTC(),
	}

	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal dead-letter event: %w", err)
	}

	deadLetterKey := fmt.Sprintf(
		"%s-%d-%d",
		message.Topic,
		message.Partition,
		message.Offset,
	)

	return writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(deadLetterKey),
		Value: payload,
	})
}
func findUserMessages(collection *mongo.Collection) MessagesFinder {
	return func(ctx context.Context, userID string) ([]bson.M, error) {
		cursor, err := collection.Find(ctx, bson.M{"userId": userID})
		if err != nil {
			return nil, err
		}
		defer cursor.Close(ctx)

		var messages []bson.M
		if err := cursor.All(ctx, &messages); err != nil {
			return nil, err
		}

		return messages, nil
	}
}
func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	client, err := mongo.Connect(options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Connect to MongoDB: %v", err)
	}
	defer func() {
		if err := client.Disconnect(context.Background()); err != nil {
			log.Printf("Disconnect MongoDB: %v", err)
		}
	}()

	pingCtx, cancelPing := context.WithTimeout(ctx, 5*time.Second)
	defer cancelPing()

	if err := client.Ping(pingCtx, nil); err != nil {
		log.Fatalf("Ping MongoDB: %v", err)
	}
	log.Println("Connected to MongoDB")

	collection := client.Database("polyglot").Collection("sms_messages")

	http.HandleFunc("/v1/user/", messagesHandler(findUserMessages(collection)))

	go func() {
		log.Println("HTTP API listening on :8081")
		if err := http.ListenAndServe(":8081", nil); err != nil {
			log.Fatalf("HTTP server failed: %v", err)
		}
	}()

	broker := os.Getenv("KAFKA_BROKERS")
	if broker == "" {
		broker = "localhost:9092"
	}
	deadLetterWriter := &kafka.Writer{
		Addr:         kafka.TCP(broker),
		Topic:        "sms-events-dead-letter",
		Balancer:     &kafka.Hash{},
		RequiredAcks: kafka.RequireAll,
	}

	defer func() {
		if err := deadLetterWriter.Close(); err != nil {
			log.Printf("Close dead-letter writer: %v", err)
		}
	}()
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{broker},
		Topic:       "sms-events",
		GroupID:     "polyglot-go-consumer",
		StartOffset: kafka.FirstOffset,
	})
	defer func() {
		if err := reader.Close(); err != nil {
			log.Printf("Close Kafka reader: %v", err)
		}
	}()

	log.Printf("Consuming sms-events from %s", broker)

	for {
		message, err := reader.FetchMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				log.Println("Shutting down consumer")
				return
			}
			log.Printf("Fetch Kafka message: %v", err)
			time.Sleep(time.Second)
			continue
		}

		var event SmsEvent

		if err := json.Unmarshal(message.Value, &event); err != nil {
			reason := fmt.Sprintf("invalid JSON: %v", err)

			if dlqErr := publishToDeadLetter(
				ctx, deadLetterWriter, message, reason,
			); dlqErr != nil {
				log.Printf("Publish invalid JSON to DLQ: %v", dlqErr)
				continue
			}

			if err := reader.CommitMessages(ctx, message); err != nil {
				log.Printf("Commit invalid JSON event: %v", err)
			}
			continue
		}

		createdAt, err := validateEvent(event)
		if err != nil {
			if dlqErr := publishToDeadLetter(
				ctx, deadLetterWriter, message, err.Error(),
			); dlqErr != nil {
				log.Printf("Publish invalid event to DLQ: %v", dlqErr)
				continue
			}

			if err := reader.CommitMessages(ctx, message); err != nil {
				log.Printf("Commit invalid event: %v", err)
			}
			continue
		}
		document := bson.M{
			"_id":         event.MessageID,
			"messageId":   event.MessageID,
			"phoneNumber": event.PhoneNumber,
			"message":     event.Message,
			"status":      event.Status,
			"createdAt":   createdAt,
			"userId":      event.UserID,
		}

		insertCtx, cancelInsert := context.WithTimeout(ctx, 5*time.Second)
		_, err = collection.InsertOne(insertCtx, document)
		cancelInsert()

		if err != nil {
			if mongo.IsDuplicateKeyError(err) {
				log.Printf("Event already stored: %s", event.MessageID)
				if err := reader.CommitMessages(ctx, message); err != nil {
					log.Printf("Commit duplicate event: %v", err)
				}
				continue
			}

			log.Printf("Save event %s to MongoDB: %v", event.MessageID, err)
			continue
		}

		log.Printf("Saved SMS event to MongoDB: %s", event.MessageID)

		if err := reader.CommitMessages(ctx, message); err != nil {
			log.Printf("Commit Kafka offset: %v", err)
		}
	}
}
