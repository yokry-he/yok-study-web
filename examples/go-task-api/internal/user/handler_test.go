package user

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/yokry-he/yok-study-web/examples/go-task-api/internal/platform/httpx"
)

type stubUserService struct {
	create       func(context.Context, CreateInput) (User, error)
	get          func(context.Context, int64) (User, error)
	list         func(context.Context, ListInput) (Page, error)
	changeStatus func(context.Context, int64, ChangeStatusInput) (User, error)
}

func (s *stubUserService) Create(ctx context.Context, input CreateInput) (User, error) {
	return s.create(ctx, input)
}

func (s *stubUserService) Get(ctx context.Context, id int64) (User, error) {
	return s.get(ctx, id)
}

func (s *stubUserService) List(ctx context.Context, input ListInput) (Page, error) {
	return s.list(ctx, input)
}

func (s *stubUserService) ChangeStatus(ctx context.Context, id int64, input ChangeStatusInput) (User, error) {
	return s.changeStatus(ctx, id, input)
}

func TestUserHandlerCreateReturns201AndLocation(t *testing.T) {
	want := User{ID: 42, Name: "张三", Email: "user@example.com", Status: StatusActive}
	service := &stubUserService{create: func(_ context.Context, input CreateInput) (User, error) {
		if input.Name != " 张三 " || input.Email != "USER@example.com" {
			t.Fatalf("Create input = %+v", input)
		}
		return want, nil
	}}
	handler := NewHandler(service)
	recorder := serveUserHandler(t, handler.Create, http.MethodPost, "/api/users", `{"name":" 张三 ","email":"USER@example.com"}`, nil)

	if recorder.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201", recorder.Code)
	}
	if got := recorder.Header().Get("Location"); got != "/api/users/42" {
		t.Fatalf("Location = %q", got)
	}
	var envelope userEnvelope[User]
	decodeUserEnvelope(t, recorder, &envelope)
	if !envelope.Success || envelope.Data.ID != want.ID || envelope.RequestID != "handler-request" {
		t.Fatalf("envelope = %+v", envelope)
	}
}

func TestUserHandlerListParsesFilters(t *testing.T) {
	service := &stubUserService{list: func(_ context.Context, input ListInput) (Page, error) {
		if input.Page != 2 || input.PageSize != 10 || input.Status != StatusDisabled {
			t.Fatalf("List input = %+v", input)
		}
		return Page{Items: []User{{ID: 7}}, Page: 2, PageSize: 10, Total: 1}, nil
	}}
	handler := NewHandler(service)
	recorder := serveUserHandler(t, handler.List, http.MethodGet, "/api/users?page=2&pageSize=10&status=DISABLED", "", nil)
	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, body=%s", recorder.Code, recorder.Body.String())
	}
}

func TestUserHandlerGetAndInvalidID(t *testing.T) {
	service := &stubUserService{get: func(_ context.Context, id int64) (User, error) {
		if id != 9 {
			t.Fatalf("Get id = %d", id)
		}
		return User{ID: id}, nil
	}}
	handler := NewHandler(service)

	recorder := serveUserHandler(t, handler.Get, http.MethodGet, "/api/users/9", "", map[string]string{"id": "9"})
	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d", recorder.Code)
	}
	recorder = serveUserHandler(t, handler.Get, http.MethodGet, "/api/users/nope", "", map[string]string{"id": "nope"})
	assertUserError(t, recorder, http.StatusBadRequest, httpx.CodeInvalidArgument)
}

func TestUserHandlerCreateRejectsUnknownJSONField(t *testing.T) {
	service := &stubUserService{create: func(context.Context, CreateInput) (User, error) {
		t.Fatal("service must not run")
		return User{}, nil
	}}
	handler := NewHandler(service)
	recorder := serveUserHandler(t, handler.Create, http.MethodPost, "/api/users", `{"name":"张三","email":"a@example.com","role":"admin"}`, nil)
	assertUserError(t, recorder, http.StatusBadRequest, httpx.CodeUnknownField)
}

func TestUserHandlerMapsDomainErrors(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		wantStatus int
		wantCode   string
		wantField  string
	}{
		{name: "validation", err: &ValidationError{Fields: FieldErrors{"name": "姓名无效"}}, wantStatus: 422, wantCode: httpx.CodeValidationFailed, wantField: "name"},
		{name: "email conflict", err: ErrEmailConflict, wantStatus: 409, wantCode: "EMAIL_CONFLICT"},
		{name: "not found", err: ErrNotFound, wantStatus: 404, wantCode: "USER_NOT_FOUND"},
		{name: "version conflict", err: ErrVersionConflict, wantStatus: 409, wantCode: "USER_VERSION_CONFLICT"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			service := &stubUserService{create: func(context.Context, CreateInput) (User, error) {
				return User{}, test.err
			}}
			handler := NewHandler(service)
			recorder := serveUserHandler(t, handler.Create, http.MethodPost, "/api/users", `{"name":"张三","email":"a@example.com"}`, nil)
			errorBody := assertUserError(t, recorder, test.wantStatus, test.wantCode)
			if test.wantField != "" && errorBody.Fields[test.wantField] == "" {
				t.Fatalf("fields = %+v", errorBody.Fields)
			}
		})
	}
}

func TestUserHandlerChangeStatusRequiresExpectedVersion(t *testing.T) {
	service := &stubUserService{changeStatus: func(context.Context, int64, ChangeStatusInput) (User, error) {
		t.Fatal("service must not run")
		return User{}, nil
	}}
	handler := NewHandler(service)
	recorder := serveUserHandler(t, handler.ChangeStatus, http.MethodPatch, "/api/users/1/status", `{"status":"DISABLED"}`, map[string]string{"id": "1"})
	errorBody := assertUserError(t, recorder, http.StatusUnprocessableEntity, httpx.CodeValidationFailed)
	if errorBody.Fields["expectedVersion"] == "" {
		t.Fatalf("fields = %+v", errorBody.Fields)
	}
}

func TestUserHandlerChangeStatusAcceptsVersionZero(t *testing.T) {
	service := &stubUserService{changeStatus: func(_ context.Context, id int64, input ChangeStatusInput) (User, error) {
		if id != 1 || input.Status != StatusDisabled || input.ExpectedVersion != 0 {
			t.Fatalf("ChangeStatus(%d, %+v)", id, input)
		}
		return User{ID: 1, Status: StatusDisabled, Version: 1}, nil
	}}
	handler := NewHandler(service)
	recorder := serveUserHandler(t, handler.ChangeStatus, http.MethodPatch, "/api/users/1/status", `{"status":"DISABLED","expectedVersion":0}`, map[string]string{"id": "1"})
	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, body=%s", recorder.Code, recorder.Body.String())
	}
}

type userEnvelope[T any] struct {
	Success   bool             `json:"success"`
	Data      T                `json:"data"`
	Error     *httpx.ErrorBody `json:"error"`
	RequestID string           `json:"requestId"`
}

func serveUserHandler(
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

func decodeUserEnvelope(t *testing.T, recorder *httptest.ResponseRecorder, destination any) {
	t.Helper()
	if err := json.Unmarshal(recorder.Body.Bytes(), destination); err != nil {
		t.Fatalf("decode response %q: %v", recorder.Body.String(), err)
	}
}

func assertUserError(t *testing.T, recorder *httptest.ResponseRecorder, status int, code string) *httpx.ErrorBody {
	t.Helper()
	if recorder.Code != status {
		t.Fatalf("status = %d, want %d; body=%s", recorder.Code, status, recorder.Body.String())
	}
	var envelope userEnvelope[json.RawMessage]
	decodeUserEnvelope(t, recorder, &envelope)
	if envelope.Success || envelope.Error == nil || envelope.Error.Code != code || envelope.RequestID != "handler-request" {
		t.Fatalf("error envelope = %+v", envelope)
	}
	return envelope.Error
}

func TestUserHandlerDoesNotExposeUnknownErrors(t *testing.T) {
	secretErr := errors.New("database password=secret")
	service := &stubUserService{get: func(context.Context, int64) (User, error) {
		return User{}, secretErr
	}}
	handler := NewHandler(service)
	recorder := serveUserHandler(t, handler.Get, http.MethodGet, "/api/users/1", "", map[string]string{"id": "1"})
	assertUserError(t, recorder, http.StatusInternalServerError, httpx.CodeInternalError)
	if strings.Contains(recorder.Body.String(), "secret") {
		t.Fatal("internal error leaked")
	}
}
