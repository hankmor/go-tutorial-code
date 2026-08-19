package main

import (
	"context"
	"errors"
	"strings"
)

var (
	// ErrInvalidPost 表示文章输入不符合业务规则。
	ErrInvalidPost = errors.New("invalid post")
)

// PostService 定义文章业务用例。
type PostService struct{ repository PostRepository }

// NewPostService 创建文章业务服务。
func NewPostService(repository PostRepository) *PostService {
	return &PostService{repository: repository}
}

// ListPosts 查询文章列表。
func (s *PostService) ListPosts(ctx context.Context) ([]Post, error) {
	return s.repository.List(ctx)
}

// CreatePost 校验输入并创建文章。
func (s *PostService) CreatePost(ctx context.Context, req CreatePostRequest) (Post, error) {
	post := Post{Title: strings.TrimSpace(req.Title), Body: strings.TrimSpace(req.Body)}
	if post.Title == "" || post.Body == "" || len([]rune(post.Title)) > 120 {
		return Post{}, ErrInvalidPost
	}
	if err := s.repository.Create(ctx, &post); err != nil {
		return Post{}, err
	}
	return post, nil
}
