package main

import (
	"context"
	"net"
	"net/http"
	"testing"
	"time"
)

func TestServerReportsPortBindingFailure(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := runHTTPServer(ctx, &http.Server{Addr: listener.Addr().String()}); err == nil {
		t.Fatal("startup failure was swallowed")
	}
}

func TestServerStopsOnContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := runHTTPServer(ctx, &http.Server{Addr: "127.0.0.1:0"}); err != nil {
		t.Fatalf("graceful stop: %v", err)
	}
}
