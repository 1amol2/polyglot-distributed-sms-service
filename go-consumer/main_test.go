package main

import (
	"context"
	"encoding/json"

	"go.mongodb.org/mongo-driver/v2/bson"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMessagesHandlerRejectsPost(t *testing.T) {
	handler := messagesHandler(func(
		ctx context.Context,
		userID string,
	) ([]bson.M, error) {
		return []bson.M{}, nil
	})

	req := httptest.NewRequest(
		http.MethodPost,
		"/v1/user/user123/messages",
		nil,
	)
	recorder := httptest.NewRecorder()

	handler(recorder, req)

	if recorder.Code != http.StatusMethodNotAllowed {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusMethodNotAllowed,
			recorder.Code,
		)
	}

	if got := recorder.Header().Get("Allow"); got != http.MethodGet {
		t.Errorf("expected Allow: GET, got %q", got)
	}
}

func TestMessagesHandlerRejectsInvalidUserID(t *testing.T) {
	handler := messagesHandler(func(
		ctx context.Context,
		userID string,
	) ([]bson.M, error) {
		return []bson.M{}, nil
	})

	tests := []struct {
		name string
		path string
	}{
		{name: "empty user ID", path: "/v1/user//messages"},
		{name: "nested user ID", path: "/v1/user/a/b/messages"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			recorder := httptest.NewRecorder()

			handler(recorder, req)

			if recorder.Code != http.StatusBadRequest {
				t.Errorf(
					"expected status %d, got %d",
					http.StatusBadRequest,
					recorder.Code,
				)
			}
		})
	}
}

func TestMessagesHandlerRejectsUnknownRoute(t *testing.T) {
	handler := messagesHandler(func(
		ctx context.Context,
		userID string,
	) ([]bson.M, error) {
		return []bson.M{}, nil
	})

	req := httptest.NewRequest(
		http.MethodGet,
		"/v1/user/user123/other",
		nil,
	)
	recorder := httptest.NewRecorder()

	handler(recorder, req)

	if recorder.Code != http.StatusNotFound {
		t.Errorf(
			"expected status %d, got %d",
			http.StatusNotFound,
			recorder.Code,
		)
	}

}

func TestMessagesHandlerReturnsUserMessages(t *testing.T) {
	expectedMessage := bson.M{
		"messageId":   "test-message-1",
		"userId":      "user123",
		"phoneNumber": "911111111111",
		"message":     "Hello from test",
		"status":      "SUCCESS",
	}

	handler := messagesHandler(func(
		ctx context.Context,
		userID string,
	) ([]bson.M, error) {
		if userID != "user123" {
			t.Errorf("expected userID user123, got %q", userID)
		}

		return []bson.M{expectedMessage}, nil
	})

	req := httptest.NewRequest(
		http.MethodGet,
		"/v1/user/user123/messages",
		nil,
	)
	recorder := httptest.NewRecorder()

	handler(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", recorder.Code)
	}

	if got := recorder.Header().Get("Content-Type"); got != "application/json" {
		t.Errorf("expected application/json, got %q", got)
	}

	var messages []bson.M
	if err := json.Unmarshal(recorder.Body.Bytes(), &messages); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(messages) != 1 {
		t.Fatalf("expected 1 message, got %d", len(messages))
	}

	if messages[0]["messageId"] != "test-message-1" {
		t.Errorf("unexpected messageId: %v", messages[0]["messageId"])
	}

	if messages[0]["userId"] != "user123" {
		t.Errorf("unexpected userId: %v", messages[0]["userId"])
	}
}
