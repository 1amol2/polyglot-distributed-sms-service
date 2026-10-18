
package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMessagesHandlerRejectsPost(t *testing.T) {
	handler := messagesHandler(nil)

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
	handler := messagesHandler(nil)

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
	handler := messagesHandler(nil)

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
