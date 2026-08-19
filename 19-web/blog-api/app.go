package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Blog 是博客服务的依赖容器，保存数据库和 HTTP 用例所需的组件。
type Blog struct {
	// DB 是博客服务使用的 GORM 数据库连接。
	DB *gorm.DB
}

// NewBlog 创建使用内存 SQLite 的演示博客服务。
func NewBlog() (*Blog, error) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if err := migrate(db); err != nil {
		return nil, err
	}
	return &Blog{DB: db}, nil
}

// Router 组装 Repository、Service 和 Handler，并返回可测试的路由。
func (b *Blog) Router() *gin.Engine {
	repository := newGORMPostRepository(b.DB)
	service := NewPostService(repository)
	handler := newPostHandler(service)
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())
	router.GET("/healthz", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status": "ok"}) })
	router.GET("/posts", handler.list)
	router.POST("/posts", handler.create)
	return router
}
