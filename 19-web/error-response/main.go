package main

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RegisterRequest 是注册接口接收的 JSON 请求。
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=20"`
	Password string `json:"password" binding:"required,min=8"`
	Email    string `json:"email" binding:"required,email"`
}

// NewRouter 创建本章错误处理示例的 Gin 路由引擎。
func NewRouter() *gin.Engine {
	router := gin.New()
	router.Use(HTTPLoggingMiddleware(), ErrorHandler())

	router.GET("/healthz", healthHandler)
	router.GET("/users/:id", getUserHandler)
	router.GET("/panic", panicHandler)
	router.POST("/register", registerHandler)

	return router
}

// healthHandler 返回成功响应。
func healthHandler(c *gin.Context) {
	Success(c, gin.H{"status": "ok"})
}

// getUserHandler 演示业务错误和资源不存在响应。
func getUserHandler(c *gin.Context) {
	if c.Param("id") == "0" {
		writeError(c, &AppError{
			Status:  http.StatusNotFound,
			Code:    "USER_NOT_FOUND",
			Message: "user not found",
		})
		return
	}
	Success(c, gin.H{"id": c.Param("id"), "name": "Ada"})
}

// panicHandler 演示未知 Panic 的统一恢复。
func panicHandler(c *gin.Context) {
	panic("simulated panic")
}

// registerHandler 演示绑定失败和脱敏错误响应。
func registerHandler(c *gin.Context) {
	var request RegisterRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		writeError(c, ValidationError(err))
		return
	}
	Success(c, gin.H{"username": request.Username, "email": request.Email})
}

func writeError(c *gin.Context, err error) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		Failure(c, appErr.Status, appErr.Code, appErr.Message)
		return
	}
	Failure(c, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
}

func main() {
	if err := NewRouter().Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
