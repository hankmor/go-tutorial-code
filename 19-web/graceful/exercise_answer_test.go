package main

import (
	"context"
	"errors"
	"net"
	"net/http"
	"testing"
	"time"
)

func TestShutdownExerciseTimesOutAndStillCleansResources(t *testing.T) {
	requestStarted := make(chan struct{})
	releaseRequest := make(chan struct{})
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		close(requestStarted)
		<-releaseRequest
		w.WriteHeader(http.StatusOK)
	})

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	srv := &http.Server{Handler: handler}
	go func() {
		_ = srv.Serve(listener)
	}()

	requestDone := make(chan struct{})
	go func() {
		defer close(requestDone)
		resp, requestErr := http.Get("http://" + listener.Addr().String())
		if requestErr == nil {
			_ = resp.Body.Close()
		}
	}()
	<-requestStarted

	cleaned := false
	err = ShutdownWithCleanup(context.Background(), srv, 20*time.Millisecond, func() error {
		cleaned = true
		return nil
	})
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("shutdown error = %v, want context deadline exceeded", err)
	}
	if !cleaned {
		t.Fatal("cleanup function was not called")
	}

	close(releaseRequest)
	_ = srv.Close()
	<-requestDone
}
