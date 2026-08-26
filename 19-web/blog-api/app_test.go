package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type fakePostRepository struct {
	posts       []Post
	listErr     error
	createErr   error
	createCalls int
}

func (f *fakePostRepository) List(context.Context) ([]Post, error) {
	if f.listErr != nil {
		return nil, f.listErr
	}
	return append([]Post(nil), f.posts...), nil
}

func (f *fakePostRepository) Create(_ context.Context, post *Post) error {
	f.createCalls++
	if f.createErr != nil {
		return f.createErr
	}
	post.ID = uint(len(f.posts) + 1)
	f.posts = append(f.posts, *post)
	return nil
}

func TestPostServiceCreatePost(t *testing.T) {
	repositoryFailure := errors.New("repository unavailable")
	tests := []struct {
		name        string
		request     CreatePostRequest
		createErr   error
		wantTitle   string
		wantBody    string
		wantErr     error
		wantCreates int
	}{
		{
			name:        "normalizes and creates valid post",
			request:     CreatePostRequest{Title: "  Go testing  ", Body: " reliable tests "},
			wantTitle:   "Go testing",
			wantBody:    "reliable tests",
			wantCreates: 1,
		},
		{
			name:        "accepts title at rune limit",
			request:     CreatePostRequest{Title: strings.Repeat("界", 120), Body: "body"},
			wantTitle:   strings.Repeat("界", 120),
			wantBody:    "body",
			wantCreates: 1,
		},
		{
			name:    "rejects blank title",
			request: CreatePostRequest{Title: " ", Body: "body"},
			wantErr: ErrInvalidPost,
		},
		{
			name:    "rejects blank body",
			request: CreatePostRequest{Title: "title", Body: " \t"},
			wantErr: ErrInvalidPost,
		},
		{
			name:    "rejects title longer than rune limit",
			request: CreatePostRequest{Title: strings.Repeat("界", 121), Body: "body"},
			wantErr: ErrInvalidPost,
		},
		{
			name:        "propagates repository failure",
			request:     CreatePostRequest{Title: "title", Body: "body"},
			createErr:   repositoryFailure,
			wantErr:     repositoryFailure,
			wantCreates: 1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repository := &fakePostRepository{createErr: test.createErr}
			service := NewPostService(repository)

			post, err := service.CreatePost(context.Background(), test.request)
			if !errors.Is(err, test.wantErr) {
				t.Fatalf("CreatePost() error = %v, want %v", err, test.wantErr)
			}
			if repository.createCalls != test.wantCreates {
				t.Fatalf("Create() calls = %d, want %d", repository.createCalls, test.wantCreates)
			}
			if test.wantErr == nil && (post.Title != test.wantTitle || post.Body != test.wantBody || post.ID != 1) {
				t.Fatalf("CreatePost() = %+v, want title=%q body=%q id=1", post, test.wantTitle, test.wantBody)
			}
		})
	}
}

func TestPostServiceListPostsPropagatesRepositoryFailure(t *testing.T) {
	wantErr := errors.New("query unavailable")
	service := NewPostService(&fakePostRepository{listErr: wantErr})

	if _, err := service.ListPosts(context.Background()); !errors.Is(err, wantErr) {
		t.Fatalf("ListPosts() error = %v, want %v", err, wantErr)
	}
}

func TestGORMPostRepositoryCreatesAndListsPosts(t *testing.T) {
	blog := newTestBlog(t)
	repository := newGORMPostRepository(blog.DB)
	post := Post{Title: "Repository test", Body: "persisted by GORM"}

	if err := repository.Create(context.Background(), &post); err != nil {
		t.Fatal(err)
	}
	posts, err := repository.List(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if post.ID == 0 || len(posts) != 1 || posts[0] != post {
		t.Fatalf("persisted post = %+v, listed posts = %+v", post, posts)
	}
}

func TestRouterCreatesAndListsPosts(t *testing.T) {
	router := newTestBlog(t).Router()

	createRequest := httptest.NewRequest(http.MethodPost, "/posts", strings.NewReader(`{"title":"Go architecture","body":"layers"}`))
	createRequest.Header.Set("Content-Type", "application/json")
	createResponse := httptest.NewRecorder()
	router.ServeHTTP(createResponse, createRequest)
	if createResponse.Code != http.StatusCreated {
		t.Fatalf("POST /posts status = %d, want %d; body=%s", createResponse.Code, http.StatusCreated, createResponse.Body.String())
	}
	var created PostResponse
	if err := json.Unmarshal(createResponse.Body.Bytes(), &created); err != nil {
		t.Fatalf("decode create response: %v", err)
	}
	if created.ID == 0 || created.Title != "Go architecture" || created.Body != "layers" {
		t.Fatalf("POST /posts response = %+v", created)
	}

	listResponse := httptest.NewRecorder()
	router.ServeHTTP(listResponse, httptest.NewRequest(http.MethodGet, "/posts", nil))
	if listResponse.Code != http.StatusOK {
		t.Fatalf("GET /posts status = %d, want %d; body=%s", listResponse.Code, http.StatusOK, listResponse.Body.String())
	}
	var posts []PostResponse
	if err := json.Unmarshal(listResponse.Body.Bytes(), &posts); err != nil {
		t.Fatalf("decode list response: %v", err)
	}
	if len(posts) != 1 || posts[0] != created {
		t.Fatalf("GET /posts response = %+v, want %+v", posts, []PostResponse{created})
	}
}

func TestRouterRejectsInvalidCreateRequests(t *testing.T) {
	tests := []struct {
		name string
		body string
	}{
		{name: "malformed JSON", body: `{"title":`},
		{name: "missing body", body: `{"title":"Go"}`},
		{name: "blank title", body: `{"title":" ","body":"content"}`},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			router := newTestBlog(t).Router()
			request := httptest.NewRequest(http.MethodPost, "/posts", strings.NewReader(test.body))
			request.Header.Set("Content-Type", "application/json")
			response := httptest.NewRecorder()

			router.ServeHTTP(response, request)
			if response.Code != http.StatusBadRequest {
				t.Fatalf("POST /posts status = %d, want %d; body=%s", response.Code, http.StatusBadRequest, response.Body.String())
			}
			var payload map[string]string
			if err := json.Unmarshal(response.Body.Bytes(), &payload); err != nil {
				t.Fatalf("decode error response: %v", err)
			}
			if payload["error"] == "" {
				t.Fatalf("POST /posts error response = %v", payload)
			}
		})
	}
}

func TestRouterReturnsInternalServerErrorWhenDatabaseIsUnavailable(t *testing.T) {
	blog := newTestBlog(t)
	sqlDB, err := blog.DB.DB()
	if err != nil {
		t.Fatal(err)
	}
	if err := sqlDB.Close(); err != nil {
		t.Fatal(err)
	}
	response := httptest.NewRecorder()

	blog.Router().ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/posts", nil))
	if response.Code != http.StatusInternalServerError {
		t.Fatalf("GET /posts status = %d, want %d; body=%s", response.Code, http.StatusInternalServerError, response.Body.String())
	}
}

func newTestBlog(t *testing.T) *Blog {
	t.Helper()
	gin.SetMode(gin.TestMode)
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatal(err)
	}
	// SQLite 的 :memory: 数据库按连接隔离，限制为一个连接可避免测试偶发读到空库。
	sqlDB.SetMaxOpenConns(1)
	t.Cleanup(func() { _ = sqlDB.Close() })
	if err := migrate(db); err != nil {
		t.Fatal(err)
	}
	return &Blog{DB: db}
}
