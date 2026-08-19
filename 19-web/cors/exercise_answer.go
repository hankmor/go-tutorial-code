package main

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// CORSOptions 定义多来源和凭据模式下的 CORS 配置。
type CORSOptions struct {
	// AllowedOrigins 是允许读取响应的浏览器来源白名单。
	AllowedOrigins []string
	// AllowCredentials 表示是否允许浏览器携带 Cookie 等凭据。
	AllowCredentials bool
}

// CORSWithOptions 返回支持多来源白名单和凭据配置的 CORS 中间件。
func CORSWithOptions(options CORSOptions) gin.HandlerFunc {
	allowedOrigins := make(map[string]struct{}, len(options.AllowedOrigins))
	for _, origin := range options.AllowedOrigins {
		allowedOrigins[origin] = struct{}{}
	}

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		originAllowed := false
		if _, allowed := allowedOrigins[origin]; allowed {
			originAllowed = true
			// 返回匹配到的明确来源，不能在凭据模式下使用通配符。
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Vary", "Origin")
			c.Header("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Content-Type,Authorization")
			if options.AllowCredentials {
				c.Header("Access-Control-Allow-Credentials", "true")
			}
		}

		if c.Request.Method == http.MethodOptions {
			if !originAllowed || !isAllowedExercisePreflight(c) {
				c.AbortWithStatus(http.StatusForbidden)
				return
			}
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

func isAllowedExercisePreflight(c *gin.Context) bool {
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
