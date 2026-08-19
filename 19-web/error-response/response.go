package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 是所有接口返回的统一外层结构。
type Response struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// Success 写入成功响应。
func Success(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Response{
		Code:    "OK",
		Message: "success",
		Data:    data,
	})
}

// Failure 写入失败响应，并保留 HTTP 状态码和业务错误码的独立语义。
func Failure(c *gin.Context, status int, code string, message string) {
	c.JSON(status, Response{
		Code:    code,
		Message: message,
	})
}
