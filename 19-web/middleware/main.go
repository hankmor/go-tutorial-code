package main

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// ErrorResponse 是接口向客户端返回的统一错误结构。
type ErrorResponse struct {
	Error string `json:"error"`
}

// RegisterRequest 是用户注册接口接收的 JSON 请求。
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=20"`
	Password string `json:"password" binding:"required,min=8"`
	Email    string `json:"email" binding:"required,email"`
	Age      int    `json:"age" binding:"required,gte=18"`
}

// RequestID 为请求设置可追踪的请求 ID。
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = newRequestID()
		}
		c.Header("X-Request-ID", requestID)
		c.Next()
	}
}

// TimingMiddleware 记录请求完成所用的时间。
func TimingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		log.Printf("%s %s %s", c.Request.Method, c.Request.URL.Path, time.Since(start))
	}
}

// AuthMiddleware 通过演示 Token 保护管理接口。
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("Authorization") != "Bearer admin-token" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{Error: "unauthorized"})
			return
		}
		c.Next()
	}
}

// NewRouter 创建本章示例使用的 Gin 路由引擎。
func NewRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery(), RequestID(), TimingMiddleware())

	router.POST("/register", registerHandler)

	admin := router.Group("/admin")
	admin.Use(AuthMiddleware())
	admin.GET("/users", listUsersHandler)

	return router
}

func newRequestID() string {
	var raw [16]byte
	if _, err := rand.Read(raw[:]); err != nil {
		return "request-id-unavailable"
	}
	return hex.EncodeToString(raw[:])
}

func registerHandler(c *gin.Context) {
	var request RegisterRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid request"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "registered",
		"user": gin.H{
			"username": request.Username,
			"email":    request.Email,
		},
	})
}

func listUsersHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"users": []string{"Alice", "Bob"},
	})
}

func main() {
	if err := NewRouter().Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
