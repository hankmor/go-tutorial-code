package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestCORSExerciseMultipleOriginsAndCredentials(t *testing.T) {
	r := gin.New()
	r.Use(CORSWithOptions(CORSOptions{
		AllowedOrigins:   []string{"https://app.example.com", "https://admin.example.com"},
		AllowCredentials: true,
	}))
	r.GET("/data", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	for _, origin := range []string{"https://app.example.com", "https://admin.example.com"} {
		req := httptest.NewRequest(http.MethodGet, "/data", nil)
		req.Header.Set("Origin", origin)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		if got := resp.Header().Get("Access-Control-Allow-Origin"); got != origin {
			t.Fatalf("origin %q: allow origin = %q", origin, got)
		}
		if got := resp.Header().Get("Access-Control-Allow-Credentials"); got != "true" {
			t.Fatalf("origin %q: allow credentials = %q", origin, got)
		}
	}
}

// TestCORSExerciseRejectsUnsupportedPreflight 验证练习实现的预检白名单。
func TestCORSExerciseRejectsUnsupportedPreflight(t *testing.T) {
	r := gin.New()
	r.Use(CORSWithOptions(CORSOptions{AllowedOrigins: []string{"https://app.example.com"}}))

	request := httptest.NewRequest(http.MethodOptions, "/data", nil)
	request.Header.Set("Origin", "https://app.example.com")
	request.Header.Set("Access-Control-Request-Method", http.MethodDelete)
	request.Header.Set("Access-Control-Request-Headers", "X-Internal-Token")
	recorder := httptest.NewRecorder()
	r.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusForbidden)
	}
}
