package main

import (
	"log"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

// ErrorHandler 将未知 Panic 转换为脱敏的统一错误响应。
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if recovered := recover(); recovered != nil {
				log.Printf("panic: %v\n%s", recovered, debug.Stack())
				Failure(c, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
				c.Abort()
			}
		}()
		c.Next()
	}
}
