package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestNewRouter 验证基础路由、参数校验和静态路由优先级。
func TestNewRouter(t *testing.T) {
	tests := []struct {
		name       string
		target     string
		wantStatus int
		wantBody   string
	}{
		{name: "ping", target: "/ping", wantStatus: http.StatusOK, wantBody: "pong"},
		{name: "static route wins", target: "/user/create", wantStatus: http.StatusOK, wantBody: "static"},
		{name: "dynamic route", target: "/user/ada", wantStatus: http.StatusOK, wantBody: "ada"},
		{name: "invalid page", target: "/search?page=abc", wantStatus: http.StatusBadRequest, wantBody: "positive integer"},
		{name: "search name", target: "/search?name=Ada", wantStatus: http.StatusOK, wantBody: "Ada"},
		{name: "wildcard path", target: "/files/docs/readme.md", wantStatus: http.StatusOK, wantBody: "docs/readme.md"},
		{name: "invalid user id", target: "/api/v1/users/not-a-number", wantStatus: http.StatusBadRequest, wantBody: "positive integer"},
		{name: "users", target: "/api/v1/users", wantStatus: http.StatusOK, wantBody: "Alice"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, tt.target, nil)
			recorder := httptest.NewRecorder()

			NewRouter().ServeHTTP(recorder, request)

			if recorder.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", recorder.Code, tt.wantStatus)
			}
			if !strings.Contains(recorder.Body.String(), tt.wantBody) {
				t.Fatalf("body = %q, want substring %q", recorder.Body.String(), tt.wantBody)
			}
		})
	}
}
