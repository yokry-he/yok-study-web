//go:build integration

package tests_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/yokry-he/yok-study-web/examples/go-task-api/internal/app"
	"github.com/yokry-he/yok-study-web/examples/go-task-api/internal/config"
	"github.com/yokry-he/yok-study-web/examples/go-task-api/internal/platform/httpx"
	taskdomain "github.com/yokry-he/yok-study-web/examples/go-task-api/internal/task"
	userdomain "github.com/yokry-he/yok-study-web/examples/go-task-api/internal/user"
)

type integrationEnvelope[T any] struct {
	Success   bool             `json:"success"`
	Data      T                `json:"data"`
	Error     *httpx.ErrorBody `json:"error"`
	RequestID string           `json:"requestId"`
}

type integrationResponse[T any] struct {
	Status   int
	Header   http.Header
	Envelope integrationEnvelope[T]
	Body     []byte
}

func TestTaskLifecycle(t *testing.T) {
	server := startAPI(t)

	createdUser := postJSON[userdomain.User](t, server.URL+"/api/users", map[string]any{
		"name":  "张三",
		"email": "user@example.com",
	})
	if createdUser.ID <= 0 || createdUser.Status != userdomain.StatusActive || createdUser.Version != 0 {
		t.Fatalf("创建用户结果不符合契约: %+v", createdUser)
	}

	duplicate := requestJSON[userdomain.User](t, http.MethodPost, server.URL+"/api/users", map[string]any{
		"name":  "重复用户",
		"email": "USER@example.com",
	}, http.StatusConflict)
	assertIntegrationError(t, duplicate, "EMAIL_CONFLICT")

	createdTask := postJSON[taskdomain.Task](t, server.URL+"/api/tasks", map[string]any{
		"ownerId":     createdUser.ID,
		"title":       "完成 Go 文档",
		"description": "编写可运行的 API 示例",
	})
	if createdTask.OwnerID != createdUser.ID || createdTask.Status != taskdomain.StatusTodo || createdTask.Version != 0 {
		t.Fatalf("创建任务结果不符合契约: %+v", createdTask)
	}

	listed := requestJSON[taskdomain.Page](
		t,
		http.MethodGet,
		server.URL+"/api/tasks?ownerId="+formatInt64(createdUser.ID)+"&status=TODO&page=1&pageSize=10",
		nil,
		http.StatusOK,
	)
	if listed.Envelope.Data.Total != 1 || len(listed.Envelope.Data.Items) != 1 || listed.Envelope.Data.Items[0].ID != createdTask.ID {
		t.Fatalf("筛选后的任务列表不符合契约: %+v", listed.Envelope.Data)
	}

	found := requestJSON[taskdomain.Task](t, http.MethodGet, server.URL+"/api/tasks/"+formatInt64(createdTask.ID), nil, http.StatusOK)
	if found.Envelope.Data.ID != createdTask.ID || found.Envelope.Data.Description == nil {
		t.Fatalf("任务详情不符合契约: %+v", found.Envelope.Data)
	}

	updated := requestJSON[taskdomain.Task](t, http.MethodPut, server.URL+"/api/tasks/"+formatInt64(createdTask.ID), map[string]any{
		"title":           "完成 Go 全量文档",
		"description":     "补齐 API、部署和排障说明",
		"expectedVersion": createdTask.Version,
	}, http.StatusOK)
	if updated.Envelope.Data.Version != 1 || updated.Envelope.Data.Title != "完成 Go 全量文档" {
		t.Fatalf("更新后的任务不符合契约: %+v", updated.Envelope.Data)
	}

	staleUpdate := requestJSON[taskdomain.Task](t, http.MethodPut, server.URL+"/api/tasks/"+formatInt64(createdTask.ID), map[string]any{
		"title":           "过期写入",
		"expectedVersion": createdTask.Version,
	}, http.StatusConflict)
	assertIntegrationError(t, staleUpdate, "TASK_VERSION_CONFLICT")

	doing := requestJSON[taskdomain.Task](t, http.MethodPatch, server.URL+"/api/tasks/"+formatInt64(createdTask.ID)+"/status", map[string]any{
		"status":          taskdomain.StatusDoing,
		"expectedVersion": updated.Envelope.Data.Version,
	}, http.StatusOK)
	if doing.Envelope.Data.Status != taskdomain.StatusDoing || doing.Envelope.Data.Version != 2 {
		t.Fatalf("任务进入 DOING 后不符合契约: %+v", doing.Envelope.Data)
	}

	done := requestJSON[taskdomain.Task](t, http.MethodPatch, server.URL+"/api/tasks/"+formatInt64(createdTask.ID)+"/status", map[string]any{
		"status":          taskdomain.StatusDone,
		"expectedVersion": doing.Envelope.Data.Version,
	}, http.StatusOK)
	if done.Envelope.Data.Status != taskdomain.StatusDone || done.Envelope.Data.Version != 3 {
		t.Fatalf("任务进入 DONE 后不符合契约: %+v", done.Envelope.Data)
	}

	invalidRollback := requestJSON[taskdomain.Task](t, http.MethodPatch, server.URL+"/api/tasks/"+formatInt64(createdTask.ID)+"/status", map[string]any{
		"status":          taskdomain.StatusDoing,
		"expectedVersion": done.Envelope.Data.Version,
	}, http.StatusConflict)
	assertIntegrationError(t, invalidRollback, "TASK_INVALID_TRANSITION")

	staleDelete := requestJSON[struct{}](t, http.MethodDelete, server.URL+"/api/tasks/"+formatInt64(createdTask.ID)+"?expectedVersion=2", nil, http.StatusConflict)
	assertIntegrationError(t, staleDelete, "TASK_VERSION_CONFLICT")

	deleted := requestJSON[struct{}](t, http.MethodDelete, server.URL+"/api/tasks/"+formatInt64(createdTask.ID)+"?expectedVersion=3", nil, http.StatusNoContent)
	if len(deleted.Body) != 0 {
		t.Fatalf("DELETE 204 响应体必须为空，实际为 %q", deleted.Body)
	}

	missing := requestJSON[taskdomain.Task](t, http.MethodGet, server.URL+"/api/tasks/"+formatInt64(createdTask.ID), nil, http.StatusNotFound)
	assertIntegrationError(t, missing, "TASK_NOT_FOUND")
}

func TestAPIProtocolBoundaries(t *testing.T) {
	server := startAPI(t)

	t.Run("非法和溢出 ID", func(t *testing.T) {
		for _, id := range []string{"0", "-1", "9223372036854775808"} {
			response := requestJSON[struct{}](t, http.MethodGet, server.URL+"/api/tasks/"+id, nil, http.StatusBadRequest)
			assertIntegrationError(t, response, httpx.CodeInvalidArgument)
		}
	})

	t.Run("严格 JSON 和大小限制", func(t *testing.T) {
		unknown := requestRaw[struct{}](t, http.MethodPost, server.URL+"/api/users", `{"name":"张三","email":"strict@example.com","extra":true}`, "application/json", http.StatusBadRequest)
		assertIntegrationError(t, unknown, httpx.CodeUnknownField)

		multiple := requestRaw[struct{}](t, http.MethodPost, server.URL+"/api/users", `{"name":"张三","email":"strict@example.com"} {"second":true}`, "application/json", http.StatusBadRequest)
		assertIntegrationError(t, multiple, httpx.CodeInvalidJSON)

		oversizedBody := `{"name":"` + strings.Repeat("a", (1<<20)+1) + `","email":"large@example.com"}`
		oversized := requestRaw[struct{}](t, http.MethodPost, server.URL+"/api/users", oversizedBody, "application/json", http.StatusBadRequest)
		assertIntegrationError(t, oversized, httpx.CodeBodyTooLarge)
	})

	t.Run("405 和 415", func(t *testing.T) {
		method := requestRaw[struct{}](t, http.MethodDelete, server.URL+"/api/users", "", "", http.StatusMethodNotAllowed)
		assertIntegrationError(t, method, httpx.CodeMethodNotAllowed)
		if got := method.Header.Get("Allow"); got != "GET, POST" {
			t.Fatalf("Allow = %q, want GET, POST", got)
		}

		mediaType := requestRaw[struct{}](t, http.MethodPost, server.URL+"/api/users", `{"name":"张三","email":"media@example.com"}`, "text/plain", http.StatusUnsupportedMediaType)
		assertIntegrationError(t, mediaType, httpx.CodeUnsupportedMedia)
	})

	t.Run("三个写接口都要求 expectedVersion", func(t *testing.T) {
		createdUser := postJSON[userdomain.User](t, server.URL+"/api/users", map[string]any{
			"name":  "李四",
			"email": "version@example.com",
		})
		createdTask := postJSON[taskdomain.Task](t, server.URL+"/api/tasks", map[string]any{
			"ownerId": createdUser.ID,
			"title":   "验证版本参数",
		})

		update := requestJSON[struct{}](t, http.MethodPut, server.URL+"/api/tasks/"+formatInt64(createdTask.ID), map[string]any{
			"title": "缺少版本",
		}, http.StatusUnprocessableEntity)
		assertExpectedVersionError(t, update)

		status := requestJSON[struct{}](t, http.MethodPatch, server.URL+"/api/tasks/"+formatInt64(createdTask.ID)+"/status", map[string]any{
			"status": taskdomain.StatusDoing,
		}, http.StatusUnprocessableEntity)
		assertExpectedVersionError(t, status)

		deleteResponse := requestJSON[struct{}](t, http.MethodDelete, server.URL+"/api/tasks/"+formatInt64(createdTask.ID), nil, http.StatusUnprocessableEntity)
		assertExpectedVersionError(t, deleteResponse)
	})
}

func TestMiddlewareProtocolBoundariesOverHTTP(t *testing.T) {
	t.Run("deadline", func(t *testing.T) {
		handler := httpx.Chain(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				<-r.Context().Done()
				_ = httpx.WriteError(w, r.Context().Err(), httpx.RequestIDFromContext(r.Context()))
			}),
			httpx.RequestID,
			httpx.Recover(nil),
			httpx.Deadline(10*time.Millisecond),
		)
		server := httptest.NewServer(handler)
		defer server.Close()

		response := requestJSON[struct{}](t, http.MethodGet, server.URL, nil, http.StatusGatewayTimeout)
		assertIntegrationError(t, response, httpx.CodeDeadlineExceeded)
	})

	t.Run("panic 响应不泄露内部值", func(t *testing.T) {
		const secret = "database-password=never-expose"
		handler := httpx.Chain(
			http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
				panic(errors.New(secret))
			}),
			httpx.RequestID,
			httpx.Recover(slog.New(slog.NewJSONHandler(io.Discard, nil))),
		)
		server := httptest.NewServer(handler)
		defer server.Close()

		response := requestJSON[struct{}](t, http.MethodGet, server.URL, nil, http.StatusInternalServerError)
		assertIntegrationError(t, response, httpx.CodeInternalError)
		if bytes.Contains(response.Body, []byte(secret)) {
			t.Fatalf("panic 响应泄露了内部值: %s", response.Body)
		}
	})
}

func startAPI(t *testing.T) *httptest.Server {
	t.Helper()
	fixture := startPostgres(t)
	cfg := config.Config{
		Environment: "test",
		LogLevel:    slog.LevelError,
		HTTP: config.HTTPConfig{
			RequestTimeout:  2 * time.Second,
			ShutdownTimeout: 2 * time.Second,
		},
	}
	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))
	application := app.New(cfg, logger, fixture.db)
	server := httptest.NewServer(application.Handler)
	t.Cleanup(server.Close)
	return server
}

func postJSON[T any](t *testing.T, url string, input any) T {
	t.Helper()
	response := requestJSON[T](t, http.MethodPost, url, input, http.StatusCreated)
	if location := response.Header.Get("Location"); location == "" {
		t.Fatal("201 响应缺少 Location Header")
	}
	return response.Envelope.Data
}

func requestJSON[T any](t *testing.T, method, url string, input any, wantStatus int) integrationResponse[T] {
	t.Helper()
	var body string
	contentType := ""
	if input != nil {
		payload, err := json.Marshal(input)
		if err != nil {
			t.Fatalf("序列化请求 JSON: %v", err)
		}
		body = string(payload)
		contentType = "application/json"
	}
	return requestRaw[T](t, method, url, body, contentType, wantStatus)
}

func requestRaw[T any](t *testing.T, method, url, body, contentType string, wantStatus int) integrationResponse[T] {
	t.Helper()
	request, err := http.NewRequestWithContext(context.Background(), method, url, strings.NewReader(body))
	if err != nil {
		t.Fatalf("创建 HTTP 请求: %v", err)
	}
	if contentType != "" {
		request.Header.Set("Content-Type", contentType)
	}
	request.Header.Set("X-Request-ID", "integration-request")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Fatalf("执行 HTTP 请求: %v", err)
	}
	defer response.Body.Close()
	payload, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatalf("读取 HTTP 响应: %v", err)
	}
	result := integrationResponse[T]{
		Status: response.StatusCode,
		Header: response.Header.Clone(),
		Body:   payload,
	}
	if response.StatusCode != wantStatus {
		t.Fatalf("%s %s status = %d, want %d, body=%s", method, request.URL.Path, response.StatusCode, wantStatus, payload)
	}
	if wantStatus == http.StatusNoContent {
		return result
	}
	if err := json.Unmarshal(payload, &result.Envelope); err != nil {
		t.Fatalf("解析 HTTP 响应 envelope: %v, body=%s", err, payload)
	}
	if result.Envelope.RequestID == "" || response.Header.Get("X-Request-ID") == "" {
		t.Fatalf("响应缺少 request id: header=%q body=%q", response.Header.Get("X-Request-ID"), result.Envelope.RequestID)
	}
	if result.Envelope.RequestID != response.Header.Get("X-Request-ID") {
		t.Fatalf("Header 和 envelope 的 request id 不一致: header=%q body=%q", response.Header.Get("X-Request-ID"), result.Envelope.RequestID)
	}
	return result
}

func assertIntegrationError[T any](t *testing.T, response integrationResponse[T], code string) {
	t.Helper()
	if response.Envelope.Success || response.Envelope.Error == nil {
		t.Fatalf("响应不是错误 envelope: %+v", response.Envelope)
	}
	if response.Envelope.Error.Code != code {
		t.Fatalf("错误码 = %q, want %q, body=%s", response.Envelope.Error.Code, code, response.Body)
	}
}

func assertExpectedVersionError[T any](t *testing.T, response integrationResponse[T]) {
	t.Helper()
	assertIntegrationError(t, response, httpx.CodeValidationFailed)
	if response.Envelope.Error.Fields["expectedVersion"] == "" {
		t.Fatalf("缺少 expectedVersion 字段错误: %+v", response.Envelope.Error)
	}
}

func formatInt64(value int64) string {
	return strconv.FormatInt(value, 10)
}
