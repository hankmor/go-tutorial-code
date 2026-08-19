package main

import "gorm.io/gorm"

// Post 是数据库中的博客文章模型。
// Body 不直接作为 HTTP 请求模型使用，避免数据库字段和外部 API 契约互相绑死。
type Post struct {
	ID    uint   `json:"id" gorm:"primaryKey"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

// CreatePostRequest 描述创建文章时允许客户端提交的字段。
type CreatePostRequest struct {
	Title string `json:"title" binding:"required"`
	Body  string `json:"body" binding:"required"`
}

// PostResponse 描述返回给 HTTP 客户端的文章数据。
type PostResponse struct {
	ID    uint   `json:"id"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

func toPostResponse(post Post) PostResponse {
	return PostResponse{ID: post.ID, Title: post.Title, Body: post.Body}
}

func migrate(db *gorm.DB) error {
	return db.AutoMigrate(&Post{})
}
