package httpx

import (
	"context"
	"errors"
	"net/http"
)

// API 错误码由所有 HTTP Handler 共享，避免各业务包重复定义协议字符串。
const (
	CodeInvalidJSON      = "INVALID_JSON"
	CodeUnknownField     = "UNKNOWN_FIELD"
	CodeBodyTooLarge     = "BODY_TOO_LARGE"
	CodeInvalidArgument  = "INVALID_ARGUMENT"
	CodeValidationFailed = "VALIDATION_FAILED"
	CodeNotFound         = "NOT_FOUND"
	CodeConflict         = "CONFLICT"
	CodeMethodNotAllowed = "METHOD_NOT_ALLOWED"
	CodeUnsupportedMedia = "UNSUPPORTED_MEDIA_TYPE"
	CodeDeadlineExceeded = "DEADLINE_EXCEEDED"
	CodeInternalError    = "INTERNAL_ERROR"
)

type FieldErrors map[string]string

type APIError struct {
	Status  int
	Code    string
	Message string
	Fields  FieldErrors
	Cause   error
}

func (e *APIError) Error() string {
	if e == nil {
		return ""
	}
	return e.Message
}

func (e *APIError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Cause
}

func NewAPIError(status int, code, message string, fields FieldErrors) *APIError {
	return &APIError{
		Status:  status,
		Code:    code,
		Message: message,
		Fields:  fields,
	}
}

func WrapAPIError(status int, code, message string, fields FieldErrors, cause error) *APIError {
	return &APIError{
		Status:  status,
		Code:    code,
		Message: message,
		Fields:  fields,
		Cause:   cause,
	}
}

func MapError(err error) *APIError {
	if err == nil {
		return nil
	}

	var apiErr *APIError
	if errors.As(err, &apiErr) && apiErr != nil {
		return apiErr
	}
	if errors.Is(err, context.Canceled) {
		return nil
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return WrapAPIError(
			http.StatusGatewayTimeout,
			CodeDeadlineExceeded,
			"请求处理超时，请稍后重试",
			nil,
			err,
		)
	}

	return WrapAPIError(
		http.StatusInternalServerError,
		CodeInternalError,
		"服务器内部错误",
		nil,
		err,
	)
}
