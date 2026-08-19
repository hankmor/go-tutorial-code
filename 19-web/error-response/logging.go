package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const maxLogBodySize = 4 << 10

// HTTPLoggingMiddleware 记录请求和响应的关键信息，并对 JSON 中的敏感字段脱敏。
func HTTPLoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestBody := &limitedBuffer{max: maxLogBodySize}
		if c.Request.Body != nil {
			c.Request.Body = io.NopCloser(io.TeeReader(c.Request.Body, requestBody))
		}

		writer := &responseCapture{ResponseWriter: c.Writer}
		c.Writer = writer
		start := time.Now()

		c.Next()

		log.Printf(
			"http request method=%s path=%s request_body=%s response_status=%d response_body=%s duration=%s",
			c.Request.Method,
			c.Request.URL.Path,
			formatLogBody(requestBody.Bytes(), c.GetHeader("Content-Type")),
			writer.Status(),
			formatLogBody(writer.body.Bytes(), writer.Header().Get("Content-Type")),
			time.Since(start),
		)
	}
}

// responseCapture 在保留 Gin ResponseWriter 行为的同时复制响应内容用于日志。
type responseCapture struct {
	gin.ResponseWriter
	body limitedBuffer
}

func (w *responseCapture) Write(data []byte) (int, error) {
	_, _ = w.body.Write(data)
	return w.ResponseWriter.Write(data)
}

func (w *responseCapture) WriteString(data string) (int, error) {
	_, _ = w.body.Write([]byte(data))
	return w.ResponseWriter.WriteString(data)
}

// limitedBuffer 限制日志缓存大小，避免异常大的请求或响应占用过多内存。
type limitedBuffer struct {
	bytes.Buffer
	max       int
	truncated bool
}

func (b *limitedBuffer) Write(data []byte) (int, error) {
	remaining := b.max - b.Len()
	if remaining > 0 {
		if len(data) > remaining {
			_, _ = b.Buffer.Write(data[:remaining])
			b.truncated = true
		} else {
			_, _ = b.Buffer.Write(data)
		}
	} else if len(data) > 0 {
		b.truncated = true
	}
	return len(data), nil
}

func formatLogBody(body []byte, contentType string) string {
	if len(body) == 0 {
		return "<empty>"
	}
	if strings.Contains(strings.ToLower(contentType), "json") {
		var value any
		if err := json.Unmarshal(body, &value); err == nil {
			redactLogFields(value)
			encoded, err := json.Marshal(value)
			if err == nil {
				return string(encoded)
			}
		}
	}
	return fmt.Sprintf("%q", string(body))
}

func redactLogFields(value any) {
	object, ok := value.(map[string]any)
	if !ok {
		return
	}
	for key, nested := range object {
		switch strings.ToLower(key) {
		case "password", "token", "authorization", "secret":
			object[key] = "[REDACTED]"
		default:
			redactLogFields(nested)
		}
	}
}

var _ http.ResponseWriter = (*responseCapture)(nil)
