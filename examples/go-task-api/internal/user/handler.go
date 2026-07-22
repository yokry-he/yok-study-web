package user

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/yokry-he/yok-study-web/examples/go-task-api/internal/platform/httpx"
)

const maxUserRequestBytes int64 = 1 << 20

var errNilUserService = errors.New("user handler: service is nil")

type UserService interface {
	Create(context.Context, CreateInput) (User, error)
	Get(context.Context, int64) (User, error)
	List(context.Context, ListInput) (Page, error)
	ChangeStatus(context.Context, int64, ChangeStatusInput) (User, error)
}

type Handler struct {
	service UserService
}

func NewHandler(service UserService) *Handler {
	return &Handler{service: service}
}

type createUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type changeUserStatusRequest struct {
	Status          Status `json:"status"`
	ExpectedVersion *int64 `json:"expectedVersion"`
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var request createUserRequest
	if err := httpx.DecodeJSON(w, r, &request, maxUserRequestBytes); err != nil {
		h.writeError(w, r, err)
		return
	}
	service, err := h.userService()
	if err != nil {
		h.writeError(w, r, err)
		return
	}
	created, err := service.Create(r.Context(), CreateInput{Name: request.Name, Email: request.Email})
	if err != nil {
		h.writeError(w, r, mapUserError(err))
		return
	}

	w.Header().Set("Location", fmt.Sprintf("/api/users/%d", created.ID))
	_ = httpx.WriteData(w, http.StatusCreated, created, requestID(r))
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.ParsePositiveID(r.PathValue("id"))
	if err != nil {
		h.writeError(w, r, err)
		return
	}
	service, err := h.userService()
	if err != nil {
		h.writeError(w, r, err)
		return
	}
	found, err := service.Get(r.Context(), id)
	if err != nil {
		h.writeError(w, r, mapUserError(err))
		return
	}
	_ = httpx.WriteData(w, http.StatusOK, found, requestID(r))
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	page, err := parseOptionalIntQuery(r, "page")
	if err != nil {
		h.writeError(w, r, err)
		return
	}
	pageSize, err := parseOptionalIntQuery(r, "pageSize")
	if err != nil {
		h.writeError(w, r, err)
		return
	}
	service, err := h.userService()
	if err != nil {
		h.writeError(w, r, err)
		return
	}
	result, err := service.List(r.Context(), ListInput{
		Page:     page,
		PageSize: pageSize,
		Status:   Status(r.URL.Query().Get("status")),
	})
	if err != nil {
		h.writeError(w, r, mapUserError(err))
		return
	}
	_ = httpx.WriteData(w, http.StatusOK, result, requestID(r))
}

func (h *Handler) ChangeStatus(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.ParsePositiveID(r.PathValue("id"))
	if err != nil {
		h.writeError(w, r, err)
		return
	}
	var request changeUserStatusRequest
	if err := httpx.DecodeJSON(w, r, &request, maxUserRequestBytes); err != nil {
		h.writeError(w, r, err)
		return
	}
	if request.ExpectedVersion == nil {
		h.writeError(w, r, userValidationAPIError(FieldErrors{
			"expectedVersion": "必须提供预期版本",
		}))
		return
	}
	service, err := h.userService()
	if err != nil {
		h.writeError(w, r, err)
		return
	}
	updated, err := service.ChangeStatus(r.Context(), id, ChangeStatusInput{
		Status:          request.Status,
		ExpectedVersion: *request.ExpectedVersion,
	})
	if err != nil {
		h.writeError(w, r, mapUserError(err))
		return
	}
	_ = httpx.WriteData(w, http.StatusOK, updated, requestID(r))
}

func (h *Handler) userService() (UserService, error) {
	if h == nil || h.service == nil {
		return nil, errNilUserService
	}
	return h.service, nil
}

func (h *Handler) writeError(w http.ResponseWriter, r *http.Request, err error) {
	_ = httpx.WriteError(w, err, requestID(r))
}

func mapUserError(err error) error {
	if err == nil {
		return nil
	}
	var validation *ValidationError
	if errors.As(err, &validation) && validation != nil {
		return httpx.WrapAPIError(
			http.StatusUnprocessableEntity,
			httpx.CodeValidationFailed,
			"用户参数校验失败",
			userHTTPFields(validation.Fields),
			err,
		)
	}
	switch {
	case errors.Is(err, ErrEmailConflict):
		return httpx.WrapAPIError(http.StatusConflict, "EMAIL_CONFLICT", "邮箱已被使用", nil, err)
	case errors.Is(err, ErrNotFound):
		return httpx.WrapAPIError(http.StatusNotFound, "USER_NOT_FOUND", "用户不存在", nil, err)
	case errors.Is(err, ErrVersionConflict):
		return httpx.WrapAPIError(http.StatusConflict, "USER_VERSION_CONFLICT", "用户已被其他请求修改，请刷新后重试", nil, err)
	default:
		return err
	}
}

func userValidationAPIError(fields FieldErrors) error {
	err := &ValidationError{Fields: fields}
	return mapUserError(err)
}

func userHTTPFields(fields FieldErrors) httpx.FieldErrors {
	converted := make(httpx.FieldErrors, len(fields))
	for field, message := range fields {
		converted[field] = message
	}
	return converted
}

func parseOptionalIntQuery(r *http.Request, name string) (int, error) {
	raw := r.URL.Query().Get(name)
	if raw == "" {
		return 0, nil
	}
	value, err := strconv.Atoi(raw)
	if err != nil {
		return 0, httpx.NewAPIError(
			http.StatusBadRequest,
			httpx.CodeInvalidArgument,
			fmt.Sprintf("查询参数 %s 必须是整数", name),
			httpx.FieldErrors{name: "必须是整数"},
		)
	}
	return value, nil
}

func requestID(r *http.Request) string {
	if r == nil {
		return ""
	}
	return httpx.RequestIDFromContext(r.Context())
}
