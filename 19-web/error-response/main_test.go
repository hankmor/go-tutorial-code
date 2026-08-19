package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestResponses verifies success, business error, panic recovery and validation responses.
func TestResponses(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		target     string
		body       string
		wantStatus int
		wantCode   string
	}{
		{name: "health", method: http.MethodGet, target: "/healthz", wantStatus: http.StatusOK, wantCode: "OK"},
		{name: "not found", method: http.MethodGet, target: "/users/0", wantStatus: http.StatusNotFound, wantCode: "USER_NOT_FOUND"},
		{name: "panic", method: http.MethodGet, target: "/panic", wantStatus: http.StatusInternalServerError, wantCode: "INTERNAL_ERROR"},
		{
			name:       "invalid registration",
			method:     http.MethodPost,
			target:     "/register",
			body:       "{\"username\":\"ada\",\"password\":\"short\",\"email\":\"bad\"}",
			wantStatus: http.StatusBadRequest,
			wantCode:   "INVALID_REQUEST",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.method, tt.target, strings.NewReader(tt.body))
			if tt.method == http.MethodPost {
				request.Header.Set("Content-Type", "application/json")
			}
			recorder := httptest.NewRecorder()

			NewRouter().ServeHTTP(recorder, request)

			if recorder.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", recorder.Code, tt.wantStatus)
			}
			var response Response
			if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
				t.Fatalf("decode response: %v", err)
			}
			if response.Code != tt.wantCode {
				t.Fatalf("code = %q, want %q", response.Code, tt.wantCode)
			}
			if strings.Contains(recorder.Body.String(), "runtime/debug") {
				t.Fatal("response leaked internal stack details")
			}
		})
	}
}

// TestSuccessfulRegistration verifies that valid input produces a success response.
func TestSuccessfulRegistration(t *testing.T) {
	request := httptest.NewRequest(
		http.MethodPost,
		"/register",
		strings.NewReader("{\"username\":\"ada\",\"password\":\"long-password\",\"email\":\"ada@example.com\"}"),
	)
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	NewRouter().ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
}
