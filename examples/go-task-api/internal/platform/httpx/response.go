package httpx

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type ErrorBody struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Fields  FieldErrors `json:"fields,omitempty"`
}

type Envelope struct {
	Success   bool       `json:"success"`
	Data      any        `json:"data,omitempty"`
	Error     *ErrorBody `json:"error,omitempty"`
	RequestID string     `json:"requestId"`
}

func WriteData(w http.ResponseWriter, status int, data any, requestID string) error {
	if err := validateResponseStatus(status); err != nil {
		return err
	}

	payload, err := json.Marshal(Envelope{
		Success:   true,
		Data:      data,
		RequestID: requestID,
	})
	if err != nil {
		return fmt.Errorf("marshal data response: %w", err)
	}

	return writeJSON(w, status, payload)
}

func WriteError(w http.ResponseWriter, err error, requestID string) error {
	if err == nil {
		return errors.New("write error response: err must not be nil")
	}

	apiErr := MapError(err)
	if apiErr == nil {
		return nil
	}
	if statusErr := validateResponseStatus(apiErr.Status); statusErr != nil {
		return statusErr
	}

	// 只序列化公开字段，Cause 留给上层日志，不进入响应结构。
	payload, marshalErr := json.Marshal(Envelope{
		Success: false,
		Error: &ErrorBody{
			Code:    apiErr.Code,
			Message: apiErr.Message,
			Fields:  apiErr.Fields,
		},
		RequestID: requestID,
	})
	if marshalErr != nil {
		return fmt.Errorf("marshal error response: %w", marshalErr)
	}

	return writeJSON(w, apiErr.Status, payload)
}

func validateResponseStatus(status int) error {
	if status < 200 || status > 599 {
		return fmt.Errorf("invalid HTTP response status: %d", status)
	}
	return nil
}

func writeJSON(w http.ResponseWriter, status int, payload []byte) error {
	// Header 和状态码必须在完整 envelope 序列化成功后才能提交。
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if _, err := w.Write(payload); err != nil {
		return fmt.Errorf("write JSON response: %w", err)
	}
	return nil
}
