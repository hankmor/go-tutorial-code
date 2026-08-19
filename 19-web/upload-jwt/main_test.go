package main

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// TestTokenAndProtectedRoute verifies token generation and authentication.
func TestTokenAndProtectedRoute(t *testing.T) {
	secret := []byte(strings.Repeat("s", 32))
	router := NewRouter(secret, t.TempDir())

	tokenResponse := httptest.NewRecorder()
	router.ServeHTTP(tokenResponse, httptest.NewRequest(http.MethodGet, "/token", nil))
	if tokenResponse.Code != http.StatusOK {
		t.Fatalf("token status = %d, want %d", tokenResponse.Code, http.StatusOK)
	}

	request := httptest.NewRequest(http.MethodGet, "/api/profile", nil)
	request.Header.Set("Authorization", "Bearer "+tokenResponse.Body.String())
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("profile status = %d, want %d", response.Code, http.StatusOK)
	}
}

// TestExpiredTokenIsRejected verifies expiration is enforced.
func TestExpiredTokenIsRejected(t *testing.T) {
	secret := []byte(strings.Repeat("s", 32))
	claims := Claims{
		UserID: 1,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "go-tutorial",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Minute)),
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secret)
	if err != nil {
		t.Fatal(err)
	}

	request := httptest.NewRequest(http.MethodGet, "/api/profile", nil)
	request.Header.Set("Authorization", "Bearer "+token)
	response := httptest.NewRecorder()
	NewRouter(secret, t.TempDir()).ServeHTTP(response, request)
	if response.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusUnauthorized)
	}
}

// TestWrongIssuerIsRejected verifies that a valid signature from another issuer is not accepted.
func TestWrongIssuerIsRejected(t *testing.T) {
	secret := []byte(strings.Repeat("s", 32))
	claims := Claims{
		UserID: 1,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "another-service",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute)),
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secret)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := ParseToken(token, secret); err == nil {
		t.Fatal("ParseToken() error = nil, want wrong-issuer error")
	}
}

// TestWrongAlgorithmIsRejected verifies that the parser does not accept another HMAC algorithm.
func TestWrongAlgorithmIsRejected(t *testing.T) {
	secret := []byte(strings.Repeat("s", 32))
	claims := Claims{
		UserID: 1,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "go-tutorial",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute)),
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS512, claims).SignedString(secret)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := ParseToken(token, secret); err == nil {
		t.Fatal("ParseToken() error = nil, want wrong-algorithm error")
	}
}

// TestShortSecretIsRejected verifies the minimum signing-key requirement.
func TestShortSecretIsRejected(t *testing.T) {
	if _, err := GenerateToken(1, []byte("too-short")); err == nil {
		t.Fatal("GenerateToken() error = nil, want short-secret error")
	}
}

// TestUploadRejectsUnsupportedType verifies extension and file-content checks.
func TestUploadRejectsUnsupportedType(t *testing.T) {
	secret := []byte(strings.Repeat("s", 32))
	token, err := GenerateToken(1, secret)
	if err != nil {
		t.Fatal(err)
	}

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("file", "payload.txt")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = part.Write([]byte("not an image"))
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}

	request := httptest.NewRequest(http.MethodPost, "/api/upload", &body)
	request.Header.Set("Authorization", "Bearer "+token)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	response := httptest.NewRecorder()
	NewRouter(secret, t.TempDir()).ServeHTTP(response, request)
	if response.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusBadRequest)
	}
}

// TestUploadRejectsOversizedRequest verifies the application request-body limit.
func TestUploadRejectsOversizedRequest(t *testing.T) {
	secret := []byte(strings.Repeat("s", 32))
	token, err := GenerateToken(1, secret)
	if err != nil {
		t.Fatal(err)
	}
	request := httptest.NewRequest(http.MethodPost, "/api/upload", strings.NewReader(strings.Repeat("x", int(maxUploadSize)+1)))
	request.Header.Set("Authorization", "Bearer "+token)
	response := httptest.NewRecorder()
	NewRouter(secret, t.TempDir()).ServeHTTP(response, request)
	if response.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusRequestEntityTooLarge)
	}
}

// TestUploadSavesGeneratedName verifies a valid PNG is stored without the original path.
func TestUploadSavesGeneratedName(t *testing.T) {
	secret := []byte(strings.Repeat("s", 32))
	token, err := GenerateToken(1, secret)
	if err != nil {
		t.Fatal(err)
	}
	uploadDir := t.TempDir()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("file", "../../avatar.png")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = part.Write([]byte("\x89PNG\r\n\x1a\n"))
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}

	request := httptest.NewRequest(http.MethodPost, "/api/upload", &body)
	request.Header.Set("Authorization", "Bearer "+token)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	response := httptest.NewRecorder()
	NewRouter(secret, uploadDir).ServeHTTP(response, request)
	if response.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d, body = %s", response.Code, http.StatusCreated, response.Body.String())
	}

	var result struct {
		Name string `json:"name"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	if strings.Contains(result.Name, "avatar") || strings.Contains(result.Name, "..") {
		t.Fatalf("generated name = %q", result.Name)
	}
	matches, err := filepath.Glob(filepath.Join(uploadDir, "*"))
	if err != nil || len(matches) != 1 {
		t.Fatalf("saved files = %v, err = %v", matches, err)
	}
}
