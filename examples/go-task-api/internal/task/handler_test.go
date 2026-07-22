package task

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/yokry-he/yok-study-web/examples/go-task-api/internal/platform/httpx"
)

type stubTaskService struct {
	create       func(context.Context, CreateInput) (Task, error)
	get          func(context.Context, int64) (Task, error)
	list         func(context.Context, ListInput) (Page, error)
	update       func(context.Context, int64, UpdateInput) (Task, error)
	changeStatus func(context.Context, int64, ChangeStatusInput) (Task, error)
	delete       func(context.Context, int64, DeleteInput) error
}

func (s *stubTaskService) Create(ctx context.Context, input CreateInput) (Task, error) {
	return s.create(ctx, input)
}
func (s *stubTaskService) Get(ctx context.Context, id int64) (Task, error) {
	return s.get(ctx, id)
}
func (s *stubTaskService) List(ctx context.Context, input ListInput) (Page, error) {
	return s.list(ctx, input)
}
func (s *stubTaskService) Update(ctx context.Context, id int64, input UpdateInput) (Task, error) {
	return s.update(ctx, id, input)
}
func (s *stubTaskService) ChangeStatus(ctx context.Context, id int64, input ChangeStatusInput) (Task, error) {
	return s.changeStatus(ctx, id, input)
}
func (s *stubTaskService) Delete(ctx context.Context, id int64, input DeleteInput) error {
	return s.delete(ctx, id, input)
}

func TestTaskHandlerCreateReturns201AndLocation(t *testing.T) {
	dueAt := time.Date(2032, 1, 2, 3, 4, 5, 0, time.UTC)
	service := &stubTaskService{create: func(_ context.Context, input CreateInput) (Task, error) {
		if input.OwnerID != 7 || input.Title != "学习 Go" || input.Description == nil || *input.Description != "详细计划" || input.DueAt == nil || !input.DueAt.Equal(dueAt) {
			t.Fatalf("Create input = %+v", input)
		}
		return Task{ID: 11, OwnerID: 7, Title: input.Title, Status: StatusTodo}, nil
	}}
	handler := NewHandler(service)
	recorder := serveTaskHandler(t, handler.Create, http.MethodPost, "/api/tasks", `{"ownerId":7,"title":"学习 Go","description":"详细计划","dueAt":"2032-01-02T03:04:05Z"}`, nil)
	if recorder.Code != http.StatusCreated || recorder.Header().Get("Location") != "/api/tasks/11" {
		t.Fatalf("status=%d Location=%q body=%s", recorder.Code, recorder.Header().Get("Location"), recorder.Body.String())
	}
}

func TestTaskHandlerListParsesFilters(t *testing.T) {
	service := &stubTaskService{list: func(_ context.Context, input ListInput) (Page, error) {
		if input.Page != 3 || input.PageSize != 15 || input.OwnerID != 7 || input.Status != StatusDoing {
			t.Fatalf("List input = %+v", input)
		}
		return Page{Items: []Task{{ID: 1}}, Page: 3, PageSize: 15, Total: 1}, nil
	}}
	handler := NewHandler(service)
	recorder := serveTaskHandler(t, handler.List, http.MethodGet, "/api/tasks?page=3&pageSize=15&ownerId=7&status=DOING", "", nil)
	if recorder.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", recorder.Code, recorder.Body.String())
	}
}

func TestTaskHandlerGetUpdateStatusAndDelete(t *testing.T) {
	service := &stubTaskService{
		get: func(_ context.Context, id int64) (Task, error) {
			return Task{ID: id}, nil
		},
		update: func(_ context.Context, id int64, input UpdateInput) (Task, error) {
			if id != 5 || input.Title != "更新" || input.ExpectedVersion != 0 {
				t.Fatalf("Update(%d, %+v)", id, input)
			}
			return Task{ID: id, Title: input.Title, Version: 1}, nil
		},
		changeStatus: func(_ context.Context, id int64, input ChangeStatusInput) (Task, error) {
			if id != 5 || input.Status != StatusDoing || input.ExpectedVersion != 1 {
				t.Fatalf("ChangeStatus(%d, %+v)", id, input)
			}
			return Task{ID: id, Status: StatusDoing, Version: 2}, nil
		},
		delete: func(_ context.Context, id int64, input DeleteInput) error {
			if id != 5 || input.ExpectedVersion != 2 {
				t.Fatalf("Delete(%d, %+v)", id, input)
			}
			return nil
		},
	}
	handler := NewHandler(service)
	pathValues := map[string]string{"id": "5"}

	if recorder := serveTaskHandler(t, handler.Get, http.MethodGet, "/api/tasks/5", "", pathValues); recorder.Code != 200 {
		t.Fatalf("Get status=%d", recorder.Code)
	}
	if recorder := serveTaskHandler(t, handler.Update, http.MethodPut, "/api/tasks/5", `{"title":"更新","description":null,"dueAt":null,"expectedVersion":0}`, pathValues); recorder.Code != 200 {
		t.Fatalf("Update status=%d body=%s", recorder.Code, recorder.Body.String())
	}
	if recorder := serveTaskHandler(t, handler.ChangeStatus, http.MethodPatch, "/api/tasks/5/status", `{"status":"DOING","expectedVersion":1}`, pathValues); recorder.Code != 200 {
		t.Fatalf("ChangeStatus status=%d body=%s", recorder.Code, recorder.Body.String())
	}
	recorder := serveTaskHandler(t, handler.Delete, http.MethodDelete, "/api/tasks/5?expectedVersion=2", "", pathValues)
	if recorder.Code != http.StatusNoContent || recorder.Body.Len() != 0 {
		t.Fatalf("Delete status=%d body=%q", recorder.Code, recorder.Body.String())
	}
}

func TestTaskHandlerRejectsInvalidPathAndQueryValues(t *testing.T) {
	service := &stubTaskService{
		get:  func(context.Context, int64) (Task, error) { t.Fatal("service ran"); return Task{}, nil },
		list: func(context.Context, ListInput) (Page, error) { t.Fatal("service ran"); return Page{}, nil },
	}
	handler := NewHandler(service)
	recorder := serveTaskHandler(t, handler.Get, http.MethodGet, "/api/tasks/0", "", map[string]string{"id": "0"})
	assertTaskError(t, recorder, http.StatusBadRequest, httpx.CodeInvalidArgument)
	recorder = serveTaskHandler(t, handler.List, http.MethodGet, "/api/tasks?page=abc", "", nil)
	assertTaskError(t, recorder, http.StatusBadRequest, httpx.CodeInvalidArgument)
	recorder = serveTaskHandler(t, handler.List, http.MethodGet, "/api/tasks?ownerId=9223372036854775808", "", nil)
	assertTaskError(t, recorder, http.StatusBadRequest, httpx.CodeInvalidArgument)
}

func TestTaskHandlerRequiresExpectedVersionForWrites(t *testing.T) {
	service := &stubTaskService{
		update: func(context.Context, int64, UpdateInput) (Task, error) { t.Fatal("Update ran"); return Task{}, nil },
		changeStatus: func(context.Context, int64, ChangeStatusInput) (Task, error) {
			t.Fatal("ChangeStatus ran")
			return Task{}, nil
		},
		delete: func(context.Context, int64, DeleteInput) error { t.Fatal("Delete ran"); return nil },
	}
	handler := NewHandler(service)
	pathValues := map[string]string{"id": "1"}
	tests := []struct {
		name    string
		method  string
		path    string
		body    string
		handler http.HandlerFunc
	}{
		{name: "update", method: http.MethodPut, path: "/api/tasks/1", body: `{"title":"更新"}`, handler: handler.Update},
		{name: "status", method: http.MethodPatch, path: "/api/tasks/1/status", body: `{"status":"DOING"}`, handler: handler.ChangeStatus},
		{name: "delete", method: http.MethodDelete, path: "/api/tasks/1", handler: handler.Delete},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			recorder := serveTaskHandler(t, test.handler, test.method, test.path, test.body, pathValues)
			errorBody := assertTaskError(t, recorder, http.StatusUnprocessableEntity, httpx.CodeValidationFailed)
			if errorBody.Fields["expectedVersion"] == "" {
				t.Fatalf("fields=%+v", errorBody.Fields)
			}
		})
	}
}

func TestTaskHandlerMapsDomainErrors(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		wantStatus int
		wantCode   string
	}{
		{name: "validation", err: &ValidationError{Fields: FieldErrors{"title": "标题无效"}}, wantStatus: 422, wantCode: httpx.CodeValidationFailed},
		{name: "not found", err: ErrNotFound, wantStatus: 404, wantCode: "TASK_NOT_FOUND"},
		{name: "owner not found", err: ErrOwnerNotFound, wantStatus: 404, wantCode: "TASK_OWNER_NOT_FOUND"},
		{name: "owner disabled", err: ErrOwnerDisabled, wantStatus: 422, wantCode: "TASK_OWNER_DISABLED"},
		{name: "version", err: ErrVersionConflict, wantStatus: 409, wantCode: "TASK_VERSION_CONFLICT"},
		{name: "transition", err: ErrInvalidTransition, wantStatus: 409, wantCode: "TASK_INVALID_TRANSITION"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			service := &stubTaskService{create: func(context.Context, CreateInput) (Task, error) {
				return Task{}, test.err
			}}
			handler := NewHandler(service)
			recorder := serveTaskHandler(t, handler.Create, http.MethodPost, "/api/tasks", `{"ownerId":1,"title":"任务"}`, nil)
			assertTaskError(t, recorder, test.wantStatus, test.wantCode)
		})
	}
}

func TestTaskHandlerRejectsUnknownJSONField(t *testing.T) {
	service := &stubTaskService{create: func(context.Context, CreateInput) (Task, error) {
		t.Fatal("service ran")
		return Task{}, nil
	}}
	handler := NewHandler(service)
	recorder := serveTaskHandler(t, handler.Create, http.MethodPost, "/api/tasks", `{"ownerId":1,"title":"任务","priority":1}`, nil)
	assertTaskError(t, recorder, http.StatusBadRequest, httpx.CodeUnknownField)
}

func TestTaskHandlerMiddlewareRejectsMediaTypeAndOversizedBody(t *testing.T) {
	service := &stubTaskService{create: func(context.Context, CreateInput) (Task, error) {
		t.Fatal("service ran")
		return Task{}, nil
	}}
	handler := NewHandler(service)
	final := http.HandlerFunc(handler.Create)
	chain := httpx.Chain(final, httpx.RequestID, httpx.RequireJSON(1<<20))

	request := httptest.NewRequest(http.MethodPost, "/api/tasks", strings.NewReader(`{"ownerId":1,"title":"任务"}`))
	request.Header.Set("Content-Type", "text/plain")
	recorder := httptest.NewRecorder()
	chain.ServeHTTP(recorder, request)
	assertTaskErrorWithoutFixedID(t, recorder, http.StatusUnsupportedMediaType, httpx.CodeUnsupportedMedia)

	request = httptest.NewRequest(http.MethodPost, "/api/tasks", strings.NewReader(strings.Repeat("x", (1<<20)+1)))
	request.Header.Set("Content-Type", "application/json")
	recorder = httptest.NewRecorder()
	chain.ServeHTTP(recorder, request)
	assertTaskErrorWithoutFixedID(t, recorder, http.StatusBadRequest, httpx.CodeBodyTooLarge)
}

func TestTaskHandlerMapsDeadlineAndCancellation(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		wantStatus int
		wantBody   bool
	}{
		{name: "deadline", err: context.DeadlineExceeded, wantStatus: http.StatusGatewayTimeout, wantBody: true},
		{name: "canceled", err: context.Canceled, wantStatus: http.StatusOK, wantBody: false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			service := &stubTaskService{get: func(context.Context, int64) (Task, error) { return Task{}, test.err }}
			handler := NewHandler(service)
			recorder := serveTaskHandler(t, handler.Get, http.MethodGet, "/api/tasks/1", "", map[string]string{"id": "1"})
			if recorder.Code != test.wantStatus {
				t.Fatalf("status=%d want=%d", recorder.Code, test.wantStatus)
			}
			if test.wantBody {
				assertTaskError(t, recorder, http.StatusGatewayTimeout, httpx.CodeDeadlineExceeded)
			} else if recorder.Body.Len() != 0 {
				t.Fatalf("canceled body=%q", recorder.Body.String())
			}
		})
	}
}

type taskEnvelope[T any] struct {
	Success   bool             `json:"success"`
	Data      T                `json:"data"`
	Error     *httpx.ErrorBody `json:"error"`
	RequestID string           `json:"requestId"`
}

func serveTaskHandler(
	t *testing.T,
	handler http.HandlerFunc,
	method string,
	path string,
	body string,
	pathValues map[string]string,
) *httptest.ResponseRecorder {
	t.Helper()
	request := httptest.NewRequest(method, path, strings.NewReader(body))
	request.Header.Set("X-Request-ID", "handler-request")
	for key, value := range pathValues {
		request.SetPathValue(key, value)
	}
	recorder := httptest.NewRecorder()
	httpx.RequestID(handler).ServeHTTP(recorder, request)
	return recorder
}

func assertTaskError(t *testing.T, recorder *httptest.ResponseRecorder, status int, code string) *httpx.ErrorBody {
	t.Helper()
	return assertTaskErrorWithRequestID(t, recorder, status, code, "handler-request")
}

func assertTaskErrorWithoutFixedID(t *testing.T, recorder *httptest.ResponseRecorder, status int, code string) *httpx.ErrorBody {
	t.Helper()
	return assertTaskErrorWithRequestID(t, recorder, status, code, "")
}

func assertTaskErrorWithRequestID(t *testing.T, recorder *httptest.ResponseRecorder, status int, code, requestID string) *httpx.ErrorBody {
	t.Helper()
	if recorder.Code != status {
		t.Fatalf("status=%d want=%d body=%s", recorder.Code, status, recorder.Body.String())
	}
	var envelope taskEnvelope[json.RawMessage]
	if err := json.Unmarshal(recorder.Body.Bytes(), &envelope); err != nil {
		t.Fatalf("decode response %q: %v", recorder.Body.String(), err)
	}
	if envelope.Success || envelope.Error == nil || envelope.Error.Code != code {
		t.Fatalf("envelope=%+v", envelope)
	}
	if requestID != "" && envelope.RequestID != requestID {
		t.Fatalf("requestId=%q want=%q", envelope.RequestID, requestID)
	}
	if envelope.RequestID == "" {
		t.Fatal("missing requestId")
	}
	return envelope.Error
}

func TestTaskHandlerDoesNotExposeUnknownErrors(t *testing.T) {
	service := &stubTaskService{get: func(context.Context, int64) (Task, error) {
		return Task{}, errors.New("sql password=secret")
	}}
	handler := NewHandler(service)
	recorder := serveTaskHandler(t, handler.Get, http.MethodGet, "/api/tasks/1", "", map[string]string{"id": "1"})
	assertTaskError(t, recorder, http.StatusInternalServerError, httpx.CodeInternalError)
	if strings.Contains(recorder.Body.String(), "secret") {
		t.Fatal("internal detail leaked")
	}
}
