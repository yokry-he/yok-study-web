package task

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/yokry-he/yok-study-web/examples/go-task-api/internal/platform/httpx"
)

const maxTaskRequestBytes int64 = 1 << 20

var errNilTaskService = errors.New("task handler: service is nil")

type TaskService interface {
	Create(context.Context, CreateInput) (Task, error)
	Get(context.Context, int64) (Task, error)
	List(context.Context, ListInput) (Page, error)
	Update(context.Context, int64, UpdateInput) (Task, error)
	ChangeStatus(context.Context, int64, ChangeStatusInput) (Task, error)
	Delete(context.Context, int64, DeleteInput) error
}

type Handler struct {
	service TaskService
}

func NewHandler(service TaskService) *Handler {
	return &Handler{service: service}
}

type createTaskRequest struct {
	OwnerID     int64      `json:"ownerId"`
	Title       string     `json:"title"`
	Description *string    `json:"description"`
	DueAt       *time.Time `json:"dueAt"`
}

type updateTaskRequest struct {
	Title           string     `json:"title"`
	Description     *string    `json:"description"`
	DueAt           *time.Time `json:"dueAt"`
	ExpectedVersion *int64     `json:"expectedVersion"`
}

type changeTaskStatusRequest struct {
	Status          Status `json:"status"`
	ExpectedVersion *int64 `json:"expectedVersion"`
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var request createTaskRequest
	if err := httpx.DecodeJSON(w, r, &request, maxTaskRequestBytes); err != nil {
		h.writeError(w, r, err)
		return
	}
	service, err := h.taskService()
	if err != nil {
		h.writeError(w, r, err)
		return
	}
	created, err := service.Create(r.Context(), CreateInput{
		OwnerID:     request.OwnerID,
		Title:       request.Title,
		Description: request.Description,
		DueAt:       request.DueAt,
	})
	if err != nil {
		h.writeError(w, r, mapTaskError(err))
		return
	}
	w.Header().Set("Location", fmt.Sprintf("/api/tasks/%d", created.ID))
	_ = httpx.WriteData(w, http.StatusCreated, created, taskRequestID(r))
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.ParsePositiveID(r.PathValue("id"))
	if err != nil {
		h.writeError(w, r, err)
		return
	}
	service, err := h.taskService()
	if err != nil {
		h.writeError(w, r, err)
		return
	}
	found, err := service.Get(r.Context(), id)
	if err != nil {
		h.writeError(w, r, mapTaskError(err))
		return
	}
	_ = httpx.WriteData(w, http.StatusOK, found, taskRequestID(r))
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	page, err := parseOptionalTaskIntQuery(r, "page")
	if err != nil {
		h.writeError(w, r, err)
		return
	}
	pageSize, err := parseOptionalTaskIntQuery(r, "pageSize")
	if err != nil {
		h.writeError(w, r, err)
		return
	}
	ownerID, err := parseOptionalPositiveTaskID(r, "ownerId")
	if err != nil {
		h.writeError(w, r, err)
		return
	}
	service, err := h.taskService()
	if err != nil {
		h.writeError(w, r, err)
		return
	}
	result, err := service.List(r.Context(), ListInput{
		Page:     page,
		PageSize: pageSize,
		OwnerID:  ownerID,
		Status:   Status(r.URL.Query().Get("status")),
	})
	if err != nil {
		h.writeError(w, r, mapTaskError(err))
		return
	}
	_ = httpx.WriteData(w, http.StatusOK, result, taskRequestID(r))
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.ParsePositiveID(r.PathValue("id"))
	if err != nil {
		h.writeError(w, r, err)
		return
	}
	var request updateTaskRequest
	if err := httpx.DecodeJSON(w, r, &request, maxTaskRequestBytes); err != nil {
		h.writeError(w, r, err)
		return
	}
	if request.ExpectedVersion == nil {
		h.writeError(w, r, taskValidationAPIError(FieldErrors{"expectedVersion": "必须提供预期版本"}))
		return
	}
	service, err := h.taskService()
	if err != nil {
		h.writeError(w, r, err)
		return
	}
	updated, err := service.Update(r.Context(), id, UpdateInput{
		Title:           request.Title,
		Description:     request.Description,
		DueAt:           request.DueAt,
		ExpectedVersion: *request.ExpectedVersion,
	})
	if err != nil {
		h.writeError(w, r, mapTaskError(err))
		return
	}
	_ = httpx.WriteData(w, http.StatusOK, updated, taskRequestID(r))
}

func (h *Handler) ChangeStatus(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.ParsePositiveID(r.PathValue("id"))
	if err != nil {
		h.writeError(w, r, err)
		return
	}
	var request changeTaskStatusRequest
	if err := httpx.DecodeJSON(w, r, &request, maxTaskRequestBytes); err != nil {
		h.writeError(w, r, err)
		return
	}
	if request.ExpectedVersion == nil {
		h.writeError(w, r, taskValidationAPIError(FieldErrors{"expectedVersion": "必须提供预期版本"}))
		return
	}
	service, err := h.taskService()
	if err != nil {
		h.writeError(w, r, err)
		return
	}
	updated, err := service.ChangeStatus(r.Context(), id, ChangeStatusInput{
		Status:          request.Status,
		ExpectedVersion: *request.ExpectedVersion,
	})
	if err != nil {
		h.writeError(w, r, mapTaskError(err))
		return
	}
	_ = httpx.WriteData(w, http.StatusOK, updated, taskRequestID(r))
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.ParsePositiveID(r.PathValue("id"))
	if err != nil {
		h.writeError(w, r, err)
		return
	}
	expectedVersion, err := parseExpectedVersionQuery(r)
	if err != nil {
		h.writeError(w, r, err)
		return
	}
	service, err := h.taskService()
	if err != nil {
		h.writeError(w, r, err)
		return
	}
	if err := service.Delete(r.Context(), id, DeleteInput{ExpectedVersion: expectedVersion}); err != nil {
		h.writeError(w, r, mapTaskError(err))
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) taskService() (TaskService, error) {
	if h == nil || h.service == nil {
		return nil, errNilTaskService
	}
	return h.service, nil
}

func (h *Handler) writeError(w http.ResponseWriter, r *http.Request, err error) {
	_ = httpx.WriteError(w, err, taskRequestID(r))
}

func mapTaskError(err error) error {
	if err == nil {
		return nil
	}
	var validation *ValidationError
	if errors.As(err, &validation) && validation != nil {
		return httpx.WrapAPIError(
			http.StatusUnprocessableEntity,
			httpx.CodeValidationFailed,
			"任务参数校验失败",
			taskHTTPFields(validation.Fields),
			err,
		)
	}
	switch {
	case errors.Is(err, ErrNotFound):
		return httpx.WrapAPIError(http.StatusNotFound, "TASK_NOT_FOUND", "任务不存在", nil, err)
	case errors.Is(err, ErrOwnerNotFound):
		return httpx.WrapAPIError(http.StatusNotFound, "TASK_OWNER_NOT_FOUND", "任务负责人不存在", nil, err)
	case errors.Is(err, ErrOwnerDisabled):
		return httpx.WrapAPIError(http.StatusUnprocessableEntity, "TASK_OWNER_DISABLED", "任务负责人已被禁用", httpx.FieldErrors{"ownerId": "负责人必须处于启用状态"}, err)
	case errors.Is(err, ErrVersionConflict):
		return httpx.WrapAPIError(http.StatusConflict, "TASK_VERSION_CONFLICT", "任务已被其他请求修改，请刷新后重试", nil, err)
	case errors.Is(err, ErrInvalidTransition):
		return httpx.WrapAPIError(http.StatusConflict, "TASK_INVALID_TRANSITION", "任务状态不允许这样变更", httpx.FieldErrors{"status": "当前状态不能变更到目标状态"}, err)
	default:
		return err
	}
}

func taskValidationAPIError(fields FieldErrors) error {
	return mapTaskError(&ValidationError{Fields: fields})
}

func taskHTTPFields(fields FieldErrors) httpx.FieldErrors {
	converted := make(httpx.FieldErrors, len(fields))
	for field, message := range fields {
		converted[field] = message
	}
	return converted
}

func parseOptionalTaskIntQuery(r *http.Request, name string) (int, error) {
	raw := r.URL.Query().Get(name)
	if raw == "" {
		return 0, nil
	}
	value, err := strconv.Atoi(raw)
	if err != nil {
		return 0, queryArgumentError(name, "必须是整数")
	}
	return value, nil
}

func parseOptionalPositiveTaskID(r *http.Request, name string) (int64, error) {
	raw := r.URL.Query().Get(name)
	if raw == "" {
		return 0, nil
	}
	value, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || value <= 0 {
		return 0, queryArgumentError(name, "必须是正整数")
	}
	return value, nil
}

func parseExpectedVersionQuery(r *http.Request) (int64, error) {
	raw := r.URL.Query().Get("expectedVersion")
	if raw == "" {
		return 0, taskValidationAPIError(FieldErrors{"expectedVersion": "必须提供预期版本"})
	}
	value, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		return 0, queryArgumentError("expectedVersion", "必须是整数")
	}
	return value, nil
}

func queryArgumentError(field, message string) error {
	return httpx.NewAPIError(
		http.StatusBadRequest,
		httpx.CodeInvalidArgument,
		fmt.Sprintf("查询参数 %s %s", field, message),
		httpx.FieldErrors{field: message},
	)
}

func taskRequestID(r *http.Request) string {
	if r == nil {
		return ""
	}
	return httpx.RequestIDFromContext(r.Context())
}
