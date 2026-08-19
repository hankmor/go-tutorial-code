package main

import (
	"fmt"
	"net/http"
)

// AppError 表示可以安全转换为 HTTP 响应的业务错误。
type AppError struct {
	Status  int
	Code    string
	Message string
	Err     error
}

// Error 返回面向日志和错误链的摘要信息。
func (e *AppError) Error() string {
	return e.Message
}

// Unwrap 返回底层错误，供 errors.Is 和 errors.As 继续判断。
func (e *AppError) Unwrap() error {
	return e.Err
}

// ValidationError 将绑定或字段校验错误转换为稳定的客户端错误。
func ValidationError(err error) error {
	return &AppError{
		Status:  http.StatusBadRequest,
		Code:    "INVALID_REQUEST",
		Message: "invalid request",
		Err:     fmt.Errorf("validation failed: %w", err),
	}
}
