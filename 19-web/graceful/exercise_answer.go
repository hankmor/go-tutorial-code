package main

import (
	"context"
	"errors"
	"net/http"
	"time"
)

// ShutdownWithCleanup 在关闭预算内停止 HTTP 服务，并在随后执行资源清理函数。
func ShutdownWithCleanup(
	parent context.Context,
	srv *http.Server,
	timeout time.Duration,
	cleanups ...func() error,
) error {
	shutdownCtx, cancel := context.WithTimeout(parent, timeout)
	defer cancel()

	shutdownErr := srv.Shutdown(shutdownCtx)
	cleanupErrors := make([]error, 0, len(cleanups)+1)
	if shutdownErr != nil {
		cleanupErrors = append(cleanupErrors, shutdownErr)
	}
	for _, cleanup := range cleanups {
		if cleanup != nil {
			cleanupErrors = append(cleanupErrors, cleanup())
		}
	}
	return errors.Join(cleanupErrors...)
}
