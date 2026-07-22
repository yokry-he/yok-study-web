package app

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/yokry-he/yok-study-web/examples/go-task-api/internal/config"
	"github.com/yokry-he/yok-study-web/examples/go-task-api/internal/platform/httpx"
	taskdomain "github.com/yokry-he/yok-study-web/examples/go-task-api/internal/task"
	userdomain "github.com/yokry-he/yok-study-web/examples/go-task-api/internal/user"
)

type fakePinger struct {
	err error
}

func (p *fakePinger) PingContext(context.Context) error {
	return p.err
}

func TestMethodNotAllowedIncludesStableAllow(t *testing.T) {
	handler := httpx.RequestID(MethodNotAllowed(http.MethodPost, http.MethodGet, http.MethodPost))
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodDelete, "/api/users", nil)
	request.Header.Set("X-Request-ID", "method-request")
	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d", recorder.Code)
	}
	if got := recorder.Header().Get("Allow"); got != "GET, POST" {
		t.Fatalf("Allow = %q", got)
	}
	assertAppError(t, recorder, httpx.CodeMethodNotAllowed, "method-request")
}

func TestRouterReturnsJSON404And405(t *testing.T) {
	router := newTestRouter(&fakePinger{})

	recorder := serveApp(t, router, http.MethodGet, "/missing", "")
	if recorder.Code != http.StatusNotFound {
		t.Fatalf("404 status = %d body=%s", recorder.Code, recorder.Body.String())
	}
	assertAppError(t, recorder, httpx.CodeNotFound, "router-request")

	recorder = serveApp(t, router, http.MethodDelete, "/api/users", "")
	if recorder.Code != http.StatusMethodNotAllowed {
		t.Fatalf("405 status = %d body=%s", recorder.Code, recorder.Body.String())
	}
	if got := recorder.Header().Get("Allow"); got != "GET, POST" {
		t.Fatalf("Allow = %q", got)
	}
	assertAppError(t, recorder, httpx.CodeMethodNotAllowed, "router-request")
}

func TestRouterDispatchesRegisteredPattern(t *testing.T) {
	router := newTestRouter(&fakePinger{})
	recorder := serveApp(t, router, http.MethodGet, "/api/users/123", "")
	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("registered route status = %d body=%s", recorder.Code, recorder.Body.String())
	}
	if strings.Contains(recorder.Body.String(), httpx.CodeNotFound) {
		t.Fatal("registered route was treated as 404")
	}
}

func TestHealthHandlersTrackReadinessAndShutdown(t *testing.T) {
	pinger := &fakePinger{}
	readiness := NewReadiness(pinger)
	router := NewRouter(
		slog.New(slog.NewJSONHandler(io.Discard, nil)),
		readiness,
		userdomain.NewHandler(nil),
		taskdomain.NewHandler(nil),
		time.Second,
	)

	if recorder := serveApp(t, router, http.MethodGet, "/health/live", ""); recorder.Code != http.StatusOK {
		t.Fatalf("live status = %d", recorder.Code)
	}
	if recorder := serveApp(t, router, http.MethodGet, "/health/ready", ""); recorder.Code != http.StatusOK {
		t.Fatalf("ready status = %d body=%s", recorder.Code, recorder.Body.String())
	}

	pinger.err = errors.New("database unavailable")
	if recorder := serveApp(t, router, http.MethodGet, "/health/ready", ""); recorder.Code != http.StatusServiceUnavailable {
		t.Fatalf("failed ready status = %d", recorder.Code)
	}
	if recorder := serveApp(t, router, http.MethodGet, "/health/live", ""); recorder.Code != http.StatusOK {
		t.Fatalf("live depends on database: %d", recorder.Code)
	}

	pinger.err = nil
	readiness.StartShutdown()
	readiness.StartShutdown()
	if recorder := serveApp(t, router, http.MethodGet, "/health/ready", ""); recorder.Code != http.StatusServiceUnavailable {
		t.Fatalf("shutdown ready status = %d", recorder.Code)
	}
}

func TestNewHTTPServerCopiesTimeouts(t *testing.T) {
	cfg := config.HTTPConfig{
		Addr:              "127.0.0.1:9876",
		ReadHeaderTimeout: time.Second,
		ReadTimeout:       2 * time.Second,
		WriteTimeout:      3 * time.Second,
		IdleTimeout:       4 * time.Second,
	}
	called := false
	handler := http.HandlerFunc(func(http.ResponseWriter, *http.Request) { called = true })
	server := NewHTTPServer(cfg, handler)
	if server.Addr != cfg.Addr || server.Handler == nil ||
		server.ReadHeaderTimeout != cfg.ReadHeaderTimeout || server.ReadTimeout != cfg.ReadTimeout ||
		server.WriteTimeout != cfg.WriteTimeout || server.IdleTimeout != cfg.IdleTimeout {
		t.Fatalf("server = %+v", server)
	}
	server.Handler.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil))
	if !called {
		t.Fatal("server did not retain handler")
	}
}

func TestRouterRequiresJSONOnlyForBodyMethods(t *testing.T) {
	router := newTestRouter(&fakePinger{})
	getRecorder := serveApp(t, router, http.MethodGet, "/api/users", "")
	if getRecorder.Code == http.StatusUnsupportedMediaType {
		t.Fatal("GET unexpectedly required Content-Type")
	}

	recorder := serveApp(t, router, http.MethodPost, "/api/users", `{"name":"张三","email":"a@example.com"}`)
	if recorder.Code != http.StatusUnsupportedMediaType {
		t.Fatalf("POST without Content-Type status = %d", recorder.Code)
	}
}

func TestAppRunShutsDownReadinessAndDatabase(t *testing.T) {
	addr := reserveAddress(t)
	db, err := sql.Open("pgx", "postgres://app:secret@127.0.0.1:1/taskdb?sslmode=disable")
	if err != nil {
		t.Fatalf("sql.Open: %v", err)
	}
	cfg := config.Config{HTTP: config.HTTPConfig{
		Addr:              addr,
		ReadHeaderTimeout: time.Second,
		ReadTimeout:       time.Second,
		WriteTimeout:      time.Second,
		IdleTimeout:       time.Second,
		RequestTimeout:    time.Second,
	}}
	var logs bytes.Buffer
	application := New(cfg, slog.New(slog.NewJSONHandler(&logs, nil)), db)
	ctx, cancel := context.WithCancel(context.Background())
	result := make(chan error, 1)
	go func() { result <- application.Run(ctx, time.Second) }()

	waitForLiveServer(t, "http://"+addr+"/health/live")
	cancel()
	select {
	case err := <-result:
		if err != nil {
			t.Fatalf("Run() error = %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Run did not finish after cancellation")
	}
	if !application.readiness.shuttingDown.Load() {
		t.Fatal("readiness was not disabled before shutdown")
	}
	if err := db.PingContext(context.Background()); err == nil {
		t.Fatal("database remained open after Run returned")
	}
	for _, message := range []string{"HTTP 服务开始监听", "HTTP 服务开始关闭", "HTTP 服务关闭完成"} {
		if !strings.Contains(logs.String(), message) {
			t.Fatalf("lifecycle logs missing %q: %s", message, logs.String())
		}
	}
}

func TestAppRunValidatesInputs(t *testing.T) {
	var nilApp *App
	if err := nilApp.Run(context.Background(), time.Second); !errors.Is(err, ErrNilApp) {
		t.Fatalf("nil Run error = %v", err)
	}
	application := New(config.Config{}, nil, nil)
	if err := application.Run(nil, time.Second); !errors.Is(err, ErrNilRunContext) {
		t.Fatalf("nil context error = %v", err)
	}
	if err := application.Run(context.Background(), 0); !errors.Is(err, ErrInvalidShutdownTTL) {
		t.Fatalf("zero timeout error = %v", err)
	}
}

func newTestRouter(pinger Pinger) http.Handler {
	return NewRouter(
		slog.New(slog.NewJSONHandler(io.Discard, nil)),
		NewReadiness(pinger),
		userdomain.NewHandler(nil),
		taskdomain.NewHandler(nil),
		50*time.Millisecond,
	)
}

func reserveAddress(t *testing.T) string {
	t.Helper()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("reserve address: %v", err)
	}
	addr := listener.Addr().String()
	if err := listener.Close(); err != nil {
		t.Fatalf("release reserved address: %v", err)
	}
	return addr
}

func waitForLiveServer(t *testing.T, url string) {
	t.Helper()
	client := &http.Client{Timeout: 200 * time.Millisecond}
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		response, err := client.Get(url)
		if err == nil {
			_ = response.Body.Close()
			if response.StatusCode == http.StatusOK {
				return
			}
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatalf("server %s did not become live", url)
}

func serveApp(t *testing.T, handler http.Handler, method, path, body string) *httptest.ResponseRecorder {
	t.Helper()
	request := httptest.NewRequest(method, path, strings.NewReader(body))
	request.Header.Set("X-Request-ID", "router-request")
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)
	return recorder
}

func assertAppError(t *testing.T, recorder *httptest.ResponseRecorder, code, requestID string) {
	t.Helper()
	var envelope struct {
		Success   bool             `json:"success"`
		Error     *httpx.ErrorBody `json:"error"`
		RequestID string           `json:"requestId"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &envelope); err != nil {
		t.Fatalf("decode %q: %v", recorder.Body.String(), err)
	}
	if envelope.Success || envelope.Error == nil || envelope.Error.Code != code || envelope.RequestID != requestID {
		t.Fatalf("envelope = %+v", envelope)
	}
}
