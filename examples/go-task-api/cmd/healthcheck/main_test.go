package main

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRunAcceptsOnly2xx(t *testing.T) {
	for _, status := range []int{http.StatusOK, http.StatusNoContent} {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(status)
		}))
		err := run(context.Background(), server.URL, server.Client())
		server.Close()
		if err != nil {
			t.Fatalf("status %d error = %v", status, err)
		}
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer server.Close()
	if err := run(context.Background(), server.URL, server.Client()); err == nil {
		t.Fatal("503 error = nil")
	}
}

func TestRunValidatesURLAndDeadline(t *testing.T) {
	client := &http.Client{Timeout: time.Second}
	if err := run(context.Background(), "", client); err == nil {
		t.Fatal("empty URL error = nil")
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := run(ctx, "http://127.0.0.1:1/health", client); !errors.Is(err, context.Canceled) {
		t.Fatalf("canceled error = %v", err)
	}
}
