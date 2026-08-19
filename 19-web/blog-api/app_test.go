package main

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type fakePostRepository struct{ posts []Post }

func (f *fakePostRepository) List(context.Context) ([]Post, error) {
	return append([]Post(nil), f.posts...), nil
}

func (f *fakePostRepository) Create(_ context.Context, post *Post) error {
	post.ID = uint(len(f.posts) + 1)
	f.posts = append(f.posts, *post)
	return nil
}

func TestPostServiceCreatePost(t *testing.T) {
	service := NewPostService(&fakePostRepository{})
	post, err := service.CreatePost(context.Background(), CreatePostRequest{Title: "  Go  ", Body: " body "})
	if err != nil || post.Title != "Go" || post.Body != "body" || post.ID != 1 {
		t.Fatalf("CreatePost() = %+v, %v", post, err)
	}
}

func TestPostServiceRejectsInvalidPost(t *testing.T) {
	service := NewPostService(&fakePostRepository{})
	_, err := service.CreatePost(context.Background(), CreatePostRequest{Title: " ", Body: "body"})
	if !errors.Is(err, ErrInvalidPost) {
		t.Fatalf("CreatePost() error = %v, want ErrInvalidPost", err)
	}
}

func TestRouterCreatesAndListsPosts(t *testing.T) {
	blog, err := NewBlog()
	if err != nil {
		t.Fatal(err)
	}
	sqlDB, err := blog.DB.DB()
	if err != nil {
		t.Fatal(err)
	}
	defer sqlDB.Close()
	router := blog.Router()

	createRequest := httptest.NewRequest(http.MethodPost, "/posts", strings.NewReader(`{"title":"Go architecture","body":"layers"}`))
	createRequest.Header.Set("Content-Type", "application/json")
	createResponse := httptest.NewRecorder()
	router.ServeHTTP(createResponse, createRequest)
	if createResponse.Code != http.StatusCreated {
		t.Fatalf("POST /posts status = %d, want %d", createResponse.Code, http.StatusCreated)
	}

	listResponse := httptest.NewRecorder()
	router.ServeHTTP(listResponse, httptest.NewRequest(http.MethodGet, "/posts", nil))
	if listResponse.Code != http.StatusOK || !strings.Contains(listResponse.Body.String(), "Go architecture") {
		t.Fatalf("GET /posts response = %d %s", listResponse.Code, listResponse.Body.String())
	}
}
