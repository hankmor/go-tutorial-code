package main

import (
	"context"

	"gorm.io/gorm"
)

// PostRepository 定义文章数据访问能力，Service 不需要知道具体数据库实现。
type PostRepository interface {
	// List 返回按 ID 升序排列的文章。
	List(ctx context.Context) ([]Post, error)
	// Create 持久化一篇文章，并回填数据库生成的 ID。
	Create(ctx context.Context, post *Post) error
}

type gormPostRepository struct{ db *gorm.DB }

func newGORMPostRepository(db *gorm.DB) *gormPostRepository {
	return &gormPostRepository{db: db}
}

func (r *gormPostRepository) List(ctx context.Context) ([]Post, error) {
	var posts []Post
	err := r.db.WithContext(ctx).Order("id ASC").Find(&posts).Error
	return posts, err
}

func (r *gormPostRepository) Create(ctx context.Context, post *Post) error {
	return r.db.WithContext(ctx).Create(post).Error
}
