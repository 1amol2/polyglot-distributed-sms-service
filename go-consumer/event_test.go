package main

import (
	"encoding/json"
	"testing"
	"time"
)

func validTestEvent() SmsEvent {
	return SmsEvent{
		MessageID:   "msg-123",
		UserID:      "user-123",
		PhoneNumber: "+919876543210",
		Message:     "Hello",
		Status:      "SUCCESS",
		CreatedAt:   json.RawMessage(`"2026-10-18T10:00:00Z"`),
	}
}

func TestValidateEventAcceptsValidEvent(t *testing.T) {
	event := validTestEvent()

	createdAt, err := validateEvent(event)
	if err != nil {
		t.Fatalf("expected valid event, got error: %v", err)
	}

	expected := time.Date(
		2026, 10, 18, 10, 0, 0, 0, time.UTC,
	)

	if !createdAt.Equal(expected) {
		t.Errorf("expected timestamp %v, got %v", expected, createdAt)
	}
}

func TestValidateEventRejectsInvalidEvents(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(*SmsEvent)
	}{
		{
			name: "missing messageId",
			mutate: func(e *SmsEvent) {
				e.MessageID = ""
			},
		},
		{
			name: "missing userId",
			mutate: func(e *SmsEvent) {
				e.UserID = ""
			},
		},
		{
			name: "missing phone number",
			mutate: func(e *SmsEvent) {
				e.PhoneNumber = ""
			},
		},
		{
			name: "missing message",
			mutate: func(e *SmsEvent) {
				e.Message = ""
			},
		},
		{
			name: "invalid status",
			mutate: func(e *SmsEvent) {
				e.Status = "PENDING"
			},
		},
		{
			name: "invalid timestamp",
			mutate: func(e *SmsEvent) {
				e.CreatedAt = json.RawMessage(`"not-a-timestamp"`)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := validTestEvent()
			tt.mutate(&event)

			_, err := validateEvent(event)
			if err == nil {
				t.Fatal("expected validation error, got nil")
			}
		})
	}
}
