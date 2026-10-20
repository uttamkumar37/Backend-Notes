package main

import (
	"context"
	"io"
	"log/slog"
	"net"
	"net/http"
	"os"
	"sync"
	"testing"
	"time"
)

func TestAddressConfiguration(t *testing.T) {
	t.Setenv("INVESTOR_ADDR", "placeholder")
	if err := os.Unsetenv("INVESTOR_ADDR"); err != nil {
		t.Fatal(err)
	}
	if got, err := address(); err != nil || got != "127.0.0.1:8080" {
		t.Fatalf("default = %q, %v", got, err)
	}
	t.Setenv("INVESTOR_ADDR", "127.0.0.1:9090")
	if got, err := address(); err != nil || got != "127.0.0.1:9090" {
		t.Fatalf("override = %q, %v", got, err)
	}
	t.Setenv("INVESTOR_ADDR", "")
	if _, err := address(); err == nil {
		t.Fatal("empty address accepted")
	}
}

func TestRunRejectsOccupiedPort(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()
	t.Setenv("INVESTOR_ADDR", listener.Addr().String())
	if err := run(slog.New(slog.NewTextHandler(io.Discard, nil))); err == nil {
		t.Fatal("run accepted occupied address")
	}
}

func TestShutdownAllowsActiveRequestToFinish(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	entered, release := make(chan struct{}), make(chan struct{})
	var once sync.Once
	unblock := func() { once.Do(func() { close(release) }) }
	defer unblock()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		close(entered)
		<-release
		w.WriteHeader(http.StatusNoContent)
	})
	result := make(chan error, 1)
	go func() { result <- serve(ctx, listener, handler, slog.New(slog.NewTextHandler(io.Discard, nil))) }()
	response := make(chan error, 1)
	go func() {
		client := &http.Client{Timeout: 3 * time.Second}
		r, err := client.Get("http://" + listener.Addr().String())
		if err == nil {
			r.Body.Close()
			if r.StatusCode != http.StatusNoContent {
				err = io.ErrUnexpectedEOF
			}
		}
		response <- err
	}()
	select {
	case <-entered:
	case <-time.After(3 * time.Second):
		t.Fatal("request did not reach handler")
	}
	cancel()
	unblock()
	select {
	case err := <-response:
		if err != nil {
			t.Fatalf("active request interrupted: %v", err)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("request did not finish")
	}
	select {
	case err := <-result:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("server did not stop")
	}
}
