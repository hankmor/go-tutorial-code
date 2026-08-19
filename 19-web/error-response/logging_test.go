package main

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestHTTPLoggingMiddleware 验证日志同时包含请求和响应，并且不会暴露密码。
func TestHTTPLoggingMiddleware(t *testing.T) {
	var logs bytes.Buffer
	originalWriter := log.Writer()
	log.SetOutput(&logs)
	t.Cleanup(func() { log.SetOutput(originalWriter) })

	request := httptest.NewRequest(
		http.MethodPost,
		"/register",
		strings.NewReader(`{"username":"ada","password":"do-not-log","email":"ada@example.com"}`),
	)
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	NewRouter().ServeHTTP(recorder, request)

	output := logs.String()
	if !strings.Contains(output, "request_body=") {
		t.Fatalf("log does not contain request body: %q", output)
	}
	if !strings.Contains(output, "response_body=") {
		t.Fatalf("log does not contain response body: %q", output)
	}
	if strings.Contains(output, "do-not-log") {
		t.Fatalf("log leaked password: %q", output)
	}
}

// TestHTTPLoggingMiddlewareLogsRecoveredPanic 验证恢复后的 Panic 响应仍会进入访问日志。
func TestHTTPLoggingMiddlewareLogsRecoveredPanic(t *testing.T) {
	var logs bytes.Buffer
	originalWriter := log.Writer()
	log.SetOutput(&logs)
	t.Cleanup(func() { log.SetOutput(originalWriter) })

	recorder := httptest.NewRecorder()
	NewRouter().ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/panic", nil))

	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusInternalServerError)
	}
	if !strings.Contains(logs.String(), "response_status=500") {
		t.Fatalf("log does not contain recovered panic status: %q", logs.String())
	}
}
