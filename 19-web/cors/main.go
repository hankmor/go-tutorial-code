package main

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// CORS 只允许配置中的前端来源，生产环境不要使用任意来源。
func CORS(allowedOrigin string) gin.HandlerFunc {
	return func(c *gin.Context) {
		originAllowed := c.GetHeader("Origin") == allowedOrigin
		if originAllowed {
			c.Header("Access-Control-Allow-Origin", allowedOrigin)
			c.Header("Vary", "Origin")
			c.Header("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Content-Type,Authorization")
		}
		if c.Request.Method == http.MethodOptions {
			if !originAllowed || !isAllowedPreflight(c) {
				c.AbortWithStatus(http.StatusForbidden)
				return
			}
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

func isAllowedPreflight(c *gin.Context) bool {
	if method := c.GetHeader("Access-Control-Request-Method"); method != "" && method != http.MethodGet && method != http.MethodPost && method != http.MethodOptions {
		return false
	}
	for _, header := range strings.Split(c.GetHeader("Access-Control-Request-Headers"), ",") {
		switch strings.ToLower(strings.TrimSpace(header)) {
		case "", "content-type", "authorization":
		default:
			return false
		}
	}
	return true
}

// NewRouter 创建带有 CORS 白名单和健康检查路由的 Gin 引擎。
func NewRouter(allowedOrigin string) *gin.Engine {
	r := gin.New()
	r.Use(CORS(allowedOrigin))
	r.GET("/healthz", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	return r
}

func main() {
	r := NewRouter("http://localhost:3000")
	_ = r.Run(":8084")
}
