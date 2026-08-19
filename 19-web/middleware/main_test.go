package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

// TestRegister validates successful and failed JSON binding.
func TestRegister(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		wantStatus int
		wantBody   string
	}{
		{
			name:       "valid request",
			body:       "{\"username\":\"ada\",\"password\":\"long-password\",\"email\":\"ada@example.com\",\"age\":20}",
			wantStatus: http.StatusOK,
			wantBody:   "registered",
		},
		{
			name:       "invalid request",
			body:       "{\"username\":\"ada\",\"password\":\"short\",\"email\":\"bad\",\"age\":16}",
			wantStatus: http.StatusBadRequest,
			wantBody:   "invalid request",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(tt.body))
			request.Header.Set("Content-Type", "application/json")
			recorder := httptest.NewRecorder()

			NewRouter().ServeHTTP(recorder, request)

			if recorder.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", recorder.Code, tt.wantStatus)
			}
			if !strings.Contains(recorder.Body.String(), tt.wantBody) {
				t.Fatalf("body = %q, want %q", recorder.Body.String(), tt.wantBody)
			}
			if recorder.Header().Get("X-Request-ID") == "" {
				t.Fatal("X-Request-ID is empty")
			}
		})
	}
}

// TestRequestIDPreservesIncomingValue verifies proxy-provided request IDs are preserved.
func TestRequestIDPreservesIncomingValue(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/admin/users", nil)
	request.Header.Set("X-Request-ID", "trace-123")
	request.Header.Set("Authorization", "Bearer admin-token")
	recorder := httptest.NewRecorder()

	NewRouter().ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
	if got := recorder.Header().Get("X-Request-ID"); got != "trace-123" {
		t.Fatalf("X-Request-ID = %q, want trace-123", got)
	}
}

// TestAuthMiddlewareAborts verifies that missing authentication cannot reach the handler.
func TestAuthMiddlewareAborts(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/admin/users", nil)
	recorder := httptest.NewRecorder()

	NewRouter().ServeHTTP(recorder, request)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusUnauthorized)
	}
	var response ErrorResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if response.Error != "unauthorized" {
		t.Fatalf("error = %q, want unauthorized", response.Error)
	}
}

// TestMiddlewareOrder 验证全局中间件、路由组中间件和 Handler 的洋葱模型顺序。
func TestMiddlewareOrder(t *testing.T) {
	var events []string
	record := func(name string) gin.HandlerFunc {
		return func(c *gin.Context) {
			events = append(events, name+" before")
			c.Next()
			events = append(events, name+" after")
		}
	}

	router := gin.New()
	router.Use(record("global"))
	group := router.Group("/api")
	group.Use(record("group"))
	group.GET("/users", func(c *gin.Context) {
		events = append(events, "handler")
		c.Status(http.StatusOK)
	})

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/api/users", nil))

	want := []string{
		"global before", "group before", "handler", "group after", "global after",
	}
	if !slicesEqual(events, want) {
		t.Fatalf("events = %#v, want %#v", events, want)
	}
}

// TestAbortPreventsHandler 验证 Abort 后不会进入后续 Handler，但当前中间件仍需 return。
func TestAbortPreventsHandler(t *testing.T) {
	handlerCalled := false
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	})
	router.GET("/private", func(c *gin.Context) {
		handlerCalled = true
		c.Status(http.StatusOK)
	})

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/private", nil))

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusUnauthorized)
	}
	if handlerCalled {
		t.Fatal("handler was called after Abort")
	}
}

func slicesEqual(left, right []string) bool {
	if len(left) != len(right) {
		return false
	}
	for index := range left {
		if left[index] != right[index] {
			return false
		}
	}
	return true
}
