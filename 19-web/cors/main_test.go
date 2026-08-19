package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCORSAllowsConfiguredOrigin(t *testing.T) {
	r := NewRouter("http://localhost:3000")
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", resp.Code, http.StatusOK)
	}
	if got := resp.Header().Get("Access-Control-Allow-Origin"); got != "http://localhost:3000" {
		t.Fatalf("allow origin = %q, want configured origin", got)
	}
	if got := resp.Header().Get("Vary"); got != "Origin" {
		t.Fatalf("vary = %q, want Origin", got)
	}
}

func TestCORSRejectsUnconfiguredOrigin(t *testing.T) {
	r := NewRouter("http://localhost:3000")
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	req.Header.Set("Origin", "http://localhost:4000")
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	if got := resp.Header().Get("Access-Control-Allow-Origin"); got != "" {
		t.Fatalf("allow origin = %q, want empty for rejected origin", got)
	}
}

func TestCORSHandlesPreflight(t *testing.T) {
	r := NewRouter("http://localhost:3000")
	req := httptest.NewRequest(http.MethodOptions, "/healthz", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	req.Header.Set("Access-Control-Request-Method", http.MethodPost)
	req.Header.Set("Access-Control-Request-Headers", "Authorization")
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", resp.Code, http.StatusNoContent)
	}
	if got := resp.Header().Get("Access-Control-Allow-Methods"); got != "GET,POST,OPTIONS" {
		t.Fatalf("allow methods = %q, want configured methods", got)
	}
	if got := resp.Header().Get("Access-Control-Allow-Headers"); got != "Content-Type,Authorization" {
		t.Fatalf("allow headers = %q, want configured headers", got)
	}
}

// TestCORSRejectsUnsupportedPreflight 验证预检方法和请求头白名单。
func TestCORSRejectsUnsupportedPreflight(t *testing.T) {
	r := NewRouter("http://localhost:3000")
	request := httptest.NewRequest(http.MethodOptions, "/healthz", nil)
	request.Header.Set("Origin", "http://localhost:3000")
	request.Header.Set("Access-Control-Request-Method", http.MethodDelete)
	request.Header.Set("Access-Control-Request-Headers", "X-Internal-Token")
	recorder := httptest.NewRecorder()

	r.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusForbidden)
	}
}
