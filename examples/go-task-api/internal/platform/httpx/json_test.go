package httpx

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type decodeFixture struct {
	Name string `json:"name"`
}

type failingResponseWriter struct {
	*httptest.ResponseRecorder
	writeErr error
}

func (w *failingResponseWriter) Write([]byte) (int, error) {
	return 0, w.writeErr
}

type trackingResponseWriter struct {
	*httptest.ResponseRecorder
	wroteHeader bool
	wroteBody   bool
}

func newTrackingResponseWriter() *trackingResponseWriter {
	return &trackingResponseWriter{ResponseRecorder: httptest.NewRecorder()}
}

func (w *trackingResponseWriter) WriteHeader(status int) {
	w.wroteHeader = true
	w.ResponseRecorder.WriteHeader(status)
}

func (w *trackingResponseWriter) Write(body []byte) (int, error) {
	w.wroteBody = true
	return w.ResponseRecorder.Write(body)
}

func TestAPIErrorErrorReturnsPublicMessageWithoutCause(t *testing.T) {
	cause := errors.New("database password=secret")
	apiErr := WrapAPIError(http.StatusConflict, "CONFLICT", "资源状态冲突", nil, cause)

	if got := apiErr.Error(); got != "资源状态冲突" {
		t.Fatalf("APIError.Error() = %q, want %q", got, "资源状态冲突")
	}
	if strings.Contains(apiErr.Error(), cause.Error()) {
		t.Fatal("APIError.Error() exposed its private cause")
	}
}

func TestAPIErrorUnwrapSupportsErrorsIs(t *testing.T) {
	cause := errors.New("database unavailable")
	apiErr := WrapAPIError(http.StatusConflict, "CONFLICT", "资源状态冲突", nil, cause)

	if apiErr.Unwrap() != cause {
		t.Fatalf("APIError.Unwrap() = %v, want original cause", apiErr.Unwrap())
	}
	if !errors.Is(apiErr, cause) {
		t.Fatal("errors.Is() did not find the APIError cause")
	}
}

func TestNewAPIErrorPreservesAllFields(t *testing.T) {
	fields := FieldErrors{"name": "姓名不能为空"}
	apiErr := NewAPIError(http.StatusUnprocessableEntity, "VALIDATION_FAILED", "请求参数校验失败", fields)

	assertAPIError(t, apiErr, http.StatusUnprocessableEntity, "VALIDATION_FAILED", "请求参数校验失败")
	if apiErr.Fields["name"] != "姓名不能为空" {
		t.Fatalf("NewAPIError() fields = %#v", apiErr.Fields)
	}
	if apiErr.Cause != nil {
		t.Fatalf("NewAPIError() cause = %v, want nil", apiErr.Cause)
	}
}

func TestWrapAPIErrorPreservesAllFieldsAndCause(t *testing.T) {
	fields := FieldErrors{"state": "状态不允许变更"}
	cause := errors.New("optimistic lock conflict")
	apiErr := WrapAPIError(http.StatusConflict, "CONFLICT", "资源状态冲突", fields, cause)

	assertAPIError(t, apiErr, http.StatusConflict, "CONFLICT", "资源状态冲突")
	if apiErr.Fields["state"] != "状态不允许变更" {
		t.Fatalf("WrapAPIError() fields = %#v", apiErr.Fields)
	}
	if apiErr.Cause != cause {
		t.Fatalf("WrapAPIError() cause = %v, want original cause", apiErr.Cause)
	}
}

func TestMapErrorReturnsNilForNil(t *testing.T) {
	if got := MapError(nil); got != nil {
		t.Fatalf("MapError(nil) = %#v, want nil", got)
	}
}

func TestMapErrorPreservesAPIErrorIdentity(t *testing.T) {
	apiErr := NewAPIError(http.StatusNotFound, "NOT_FOUND", "资源不存在", nil)

	if got := MapError(apiErr); got != apiErr {
		t.Fatalf("MapError() = %p, want original APIError %p", got, apiErr)
	}
}

func TestMapErrorFindsWrappedAPIError(t *testing.T) {
	apiErr := NewAPIError(http.StatusNotFound, "NOT_FOUND", "资源不存在", nil)
	wrapped := fmt.Errorf("service failed: %w", apiErr)

	if got := MapError(wrapped); got != apiErr {
		t.Fatalf("MapError() = %p, want wrapped APIError %p", got, apiErr)
	}
}

func TestMapErrorReturnsNilForCanceledContext(t *testing.T) {
	if got := MapError(context.Canceled); got != nil {
		t.Fatalf("MapError(context.Canceled) = %#v, want nil", got)
	}
}

func TestMapErrorMapsDeadlineExceeded(t *testing.T) {
	apiErr := MapError(context.DeadlineExceeded)

	assertAPIError(
		t,
		apiErr,
		http.StatusGatewayTimeout,
		"DEADLINE_EXCEEDED",
		"请求处理超时，请稍后重试",
	)
	if !errors.Is(apiErr, context.DeadlineExceeded) {
		t.Fatal("deadline API error does not preserve context.DeadlineExceeded")
	}
}

func TestMapErrorMapsUnknownAndPreservesCause(t *testing.T) {
	cause := errors.New("sql query failed: password=secret")
	apiErr := MapError(cause)

	assertAPIError(t, apiErr, http.StatusInternalServerError, "INTERNAL_ERROR", "服务器内部错误")
	if !errors.Is(apiErr, cause) {
		t.Fatal("unknown API error does not preserve its cause")
	}
	if strings.Contains(apiErr.Error(), "password") || strings.Contains(apiErr.Error(), "secret") {
		t.Fatal("unknown API error exposed its private cause")
	}
}

func TestDecodeJSONAcceptsOneObjectAndTrailingWhitespace(t *testing.T) {
	tests := []struct {
		name string
		body string
	}{
		{name: "object", body: `{"name":"demo"}`},
		{name: "trailing whitespace", body: "{\"name\":\"demo\"} \n\t\r"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var dst decodeFixture
			err := decodeBody(tc.body, &dst, 1024)
			if err != nil {
				t.Fatalf("DecodeJSON() returned an unexpected error: %v", err)
			}
			if dst.Name != "demo" {
				t.Fatalf("DecodeJSON() name = %q, want demo", dst.Name)
			}
		})
	}
}

func TestDecodeJSONRejectsEmptyBody(t *testing.T) {
	for _, body := range []string{"", " \n\t\r"} {
		err := decodeBody(body, &decodeFixture{}, 1024)
		assertDecodeError(t, err, ErrEmptyJSONBody, http.StatusBadRequest, "INVALID_JSON")
		assertDoesNotEcho(t, err, body)
	}
}

func TestDecodeJSONRejectsMalformedJSON(t *testing.T) {
	const body = `{"name":"private-value"`
	err := decodeBody(body, &decodeFixture{}, 1024)

	assertDecodeError(t, err, ErrMalformedJSON, http.StatusBadRequest, "INVALID_JSON")
	assertDoesNotEcho(t, err, "private-value")
}

func TestDecodeJSONRejectsUnknownField(t *testing.T) {
	const body = `{"name":"demo","privateValue":"do-not-echo"}`
	err := decodeBody(body, &decodeFixture{}, 1024)

	assertDecodeError(t, err, ErrUnknownJSONField, http.StatusBadRequest, "UNKNOWN_FIELD")
	apiErr := requireAPIError(t, err)
	if !strings.Contains(apiErr.Message, "privateValue") {
		t.Fatalf("DecodeJSON() message %q does not identify the unknown field", apiErr.Message)
	}
	assertDoesNotEcho(t, err, "do-not-echo")
}

func TestDecodeJSONRejectsInvalidFieldType(t *testing.T) {
	const body = `{"name":123456789}`
	err := decodeBody(body, &decodeFixture{}, 1024)

	assertDecodeError(t, err, ErrInvalidJSONType, http.StatusBadRequest, "INVALID_ARGUMENT")
	apiErr := requireAPIError(t, err)
	if apiErr.Fields["name"] == "" {
		t.Fatalf("DecodeJSON() fields = %#v, want a name field error", apiErr.Fields)
	}
	assertDoesNotEcho(t, err, "123456789")
}

func TestDecodeJSONRejectsSecondValue(t *testing.T) {
	const body = `{"name":"demo"} {"name":"private-second-value"}`
	err := decodeBody(body, &decodeFixture{}, 1024)

	assertDecodeError(t, err, ErrMultipleJSONValues, http.StatusBadRequest, "INVALID_JSON")
	assertDoesNotEcho(t, err, "private-second-value")
}

func TestDecodeJSONRejectsBodyOverLimit(t *testing.T) {
	const body = `{"name":"private-value-that-is-too-long"}`
	err := decodeBody(body, &decodeFixture{}, 16)

	assertDecodeError(t, err, ErrBodyTooLarge, http.StatusBadRequest, "BODY_TOO_LARGE")
	assertDoesNotEcho(t, err, "private-value-that-is-too-long")
}

func TestDecodeJSONRejectsNonPositiveMaxBytesWithoutWriting(t *testing.T) {
	for _, maxBytes := range []int64{0, -1} {
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"name":"demo"}`))
		rec := httptest.NewRecorder()
		var dst decodeFixture

		err := DecodeJSON(rec, req, &dst, maxBytes)
		assertAPIError(t, requireAPIError(t, err), http.StatusInternalServerError, "INTERNAL_ERROR", "服务器内部错误")
		assertRecorderUntouched(t, rec)
	}
}

func FuzzDecodeJSON(f *testing.F) {
	f.Add(`{"name":"demo"}`)
	f.Fuzz(func(t *testing.T, body string) {
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
		rec := httptest.NewRecorder()
		var dst decodeFixture
		_ = DecodeJSON(rec, req, &dst, 1024)
	})
}

func TestParsePositiveID(t *testing.T) {
	id, err := ParsePositiveID("9223372036854775807")
	if err != nil {
		t.Fatalf("ParsePositiveID() returned an unexpected error: %v", err)
	}
	if id != 9223372036854775807 {
		t.Fatalf("ParsePositiveID() = %d, want MaxInt64", id)
	}

	for _, raw := range []string{"", "0", "-1", "abc", "9223372036854775808"} {
		t.Run(raw, func(t *testing.T) {
			got, err := ParsePositiveID(raw)
			if got != 0 {
				t.Fatalf("ParsePositiveID(%q) = %d, want 0", raw, got)
			}
			assertAPIError(t, requireAPIError(t, err), http.StatusBadRequest, "INVALID_ARGUMENT", "资源 ID 必须是正整数")
		})
	}
}

func TestWriteDataWritesSuccessEnvelope(t *testing.T) {
	rec := httptest.NewRecorder()
	data := map[string]any{"id": float64(42), "name": "demo"}

	if err := WriteData(rec, http.StatusOK, data, "request-123"); err != nil {
		t.Fatalf("WriteData() returned an unexpected error: %v", err)
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("WriteData() status = %d, want %d", rec.Code, http.StatusOK)
	}
	if got := rec.Header().Get("Content-Type"); got != "application/json; charset=utf-8" {
		t.Fatalf("WriteData() Content-Type = %q", got)
	}

	var envelope Envelope
	decodeResponse(t, rec, &envelope)
	if !envelope.Success || envelope.Error != nil {
		t.Fatalf("WriteData() envelope = %#v", envelope)
	}
	if envelope.RequestID != "request-123" {
		t.Fatalf("WriteData() requestId = %q", envelope.RequestID)
	}
	gotData, ok := envelope.Data.(map[string]any)
	if !ok || gotData["id"] != float64(42) || gotData["name"] != "demo" {
		t.Fatalf("WriteData() data = %#v", envelope.Data)
	}
}

func TestWriteDataMarshalFailureDoesNotCommitResponse(t *testing.T) {
	rec := newTrackingResponseWriter()
	err := WriteData(rec, http.StatusOK, make(chan int), "request-123")
	if err == nil {
		t.Fatal("WriteData() unexpectedly encoded an unsupported value")
	}
	assertTrackingRecorderUntouched(t, rec)
}

func TestWriteDataRejectsInvalidStatusWithoutWriting(t *testing.T) {
	for _, status := range []int{199, 600} {
		rec := newTrackingResponseWriter()
		if err := WriteData(rec, status, map[string]string{"name": "demo"}, "request-123"); err == nil {
			t.Fatalf("WriteData() accepted invalid status %d", status)
		}
		assertTrackingRecorderUntouched(t, rec)
	}
}

func TestWriteDataReturnsWrappedWriteFailure(t *testing.T) {
	writeFailure := errors.New("client connection closed")
	rec := &failingResponseWriter{
		ResponseRecorder: httptest.NewRecorder(),
		writeErr:         writeFailure,
	}

	err := WriteData(rec, http.StatusOK, map[string]string{"name": "demo"}, "request-123")
	if !errors.Is(err, writeFailure) {
		t.Fatalf("WriteData() error = %v, want wrapped write failure", err)
	}
}

func TestWriteErrorWritesPublic422EnvelopeWithFields(t *testing.T) {
	rec := httptest.NewRecorder()
	cause := errors.New("sql: connection postgres://user:secret@database")
	err := WrapAPIError(
		http.StatusUnprocessableEntity,
		"VALIDATION_FAILED",
		"请求参数校验失败",
		FieldErrors{"name": "姓名不能为空"},
		cause,
	)

	if writeErr := WriteError(rec, err, "request-456"); writeErr != nil {
		t.Fatalf("WriteError() returned an unexpected error: %v", writeErr)
	}
	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("WriteError() status = %d, want %d", rec.Code, http.StatusUnprocessableEntity)
	}
	if got := rec.Header().Get("Content-Type"); got != "application/json; charset=utf-8" {
		t.Fatalf("WriteError() Content-Type = %q", got)
	}

	var envelope Envelope
	decodeResponse(t, rec, &envelope)
	if envelope.Success || envelope.Data != nil || envelope.Error == nil {
		t.Fatalf("WriteError() envelope = %#v", envelope)
	}
	if envelope.RequestID != "request-456" {
		t.Fatalf("WriteError() requestId = %q", envelope.RequestID)
	}
	if envelope.Error.Code != "VALIDATION_FAILED" || envelope.Error.Message != "请求参数校验失败" {
		t.Fatalf("WriteError() public error = %#v", envelope.Error)
	}
	if envelope.Error.Fields["name"] != "姓名不能为空" {
		t.Fatalf("WriteError() fields = %#v", envelope.Error.Fields)
	}
	assertBodyDoesNotContain(t, rec, "sql", "postgres", "secret", "connection")
}

func TestWriteErrorMapsDeadline(t *testing.T) {
	rec := httptest.NewRecorder()
	if err := WriteError(rec, context.DeadlineExceeded, "deadline-request"); err != nil {
		t.Fatalf("WriteError() returned an unexpected error: %v", err)
	}

	var envelope Envelope
	decodeResponse(t, rec, &envelope)
	if rec.Code != http.StatusGatewayTimeout || envelope.Error == nil {
		t.Fatalf("WriteError() deadline response: status=%d envelope=%#v", rec.Code, envelope)
	}
	if envelope.Error.Code != "DEADLINE_EXCEEDED" || envelope.Error.Message != "请求处理超时，请稍后重试" {
		t.Fatalf("WriteError() public deadline error = %#v", envelope.Error)
	}
}

func TestWriteErrorMapsUnknownCauseWithoutExposingIt(t *testing.T) {
	rec := httptest.NewRecorder()
	cause := errors.New("SELECT password FROM users; postgres://admin:secret@db; stack trace")
	if err := WriteError(rec, cause, "unknown-request"); err != nil {
		t.Fatalf("WriteError() returned an unexpected error: %v", err)
	}

	var envelope Envelope
	decodeResponse(t, rec, &envelope)
	if rec.Code != http.StatusInternalServerError || envelope.Error == nil {
		t.Fatalf("WriteError() unknown response: status=%d envelope=%#v", rec.Code, envelope)
	}
	if envelope.Error.Code != "INTERNAL_ERROR" || envelope.Error.Message != "服务器内部错误" {
		t.Fatalf("WriteError() public error = %#v", envelope.Error)
	}
	assertBodyDoesNotContain(t, rec, "SELECT", "password", "postgres", "secret", "stack trace")
}

func TestWriteErrorCanceledRequestIsNoOp(t *testing.T) {
	rec := newTrackingResponseWriter()
	if err := WriteError(rec, context.Canceled, "canceled-request"); err != nil {
		t.Fatalf("WriteError() canceled request error = %v, want nil", err)
	}
	assertTrackingRecorderUntouched(t, rec)
}

func TestWriteErrorRejectsNilWithoutWriting(t *testing.T) {
	rec := newTrackingResponseWriter()

	if err := WriteError(rec, nil, "request-123"); err == nil {
		t.Fatal("WriteError() accepted a nil error")
	}
	assertTrackingRecorderUntouched(t, rec)
}

func TestWriteErrorRejectsInvalidAPIErrorStatusWithoutWriting(t *testing.T) {
	for _, status := range []int{199, 600} {
		t.Run(fmt.Sprintf("status_%d", status), func(t *testing.T) {
			rec := newTrackingResponseWriter()
			err := NewAPIError(status, "INVALID_STATUS", "invalid", nil)
			if writeErr := WriteError(rec, err, "request-123"); writeErr == nil {
				t.Fatalf("WriteError() accepted invalid status %d", status)
			}
			assertTrackingRecorderUntouched(t, rec)
		})
	}
}

func TestWriteErrorReturnsWrappedWriteFailure(t *testing.T) {
	writeFailure := errors.New("client connection closed")
	rec := &failingResponseWriter{
		ResponseRecorder: httptest.NewRecorder(),
		writeErr:         writeFailure,
	}

	err := WriteError(rec, NewAPIError(http.StatusBadRequest, "INVALID_ARGUMENT", "参数错误", nil), "request-123")
	if !errors.Is(err, writeFailure) {
		t.Fatalf("WriteError() error = %v, want wrapped write failure", err)
	}
}

func decodeBody(body string, dst any, maxBytes int64) error {
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	rec := httptest.NewRecorder()
	return DecodeJSON(rec, req, dst, maxBytes)
}

func assertDecodeError(t *testing.T, err, sentinel error, status int, code string) {
	t.Helper()
	if !errors.Is(err, sentinel) {
		t.Fatalf("DecodeJSON() error = %v, want errors.Is(..., %v)", err, sentinel)
	}
	apiErr := requireAPIError(t, err)
	if apiErr.Status != status || apiErr.Code != code {
		t.Fatalf("DecodeJSON() API error = %#v, want status %d and code %s", apiErr, status, code)
	}
}

func requireAPIError(t *testing.T, err error) *APIError {
	t.Helper()
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("error type = %T, want *APIError", err)
	}
	return apiErr
}

func assertAPIError(t *testing.T, apiErr *APIError, status int, code, message string) {
	t.Helper()
	if apiErr == nil {
		t.Fatal("APIError is nil")
	}
	if apiErr.Status != status || apiErr.Code != code || apiErr.Message != message {
		t.Fatalf("APIError = %#v, want status=%d code=%s message=%q", apiErr, status, code, message)
	}
}

func assertDoesNotEcho(t *testing.T, err error, privateValue string) {
	t.Helper()
	if privateValue != "" && strings.Contains(err.Error(), privateValue) {
		t.Fatalf("error exposed request content %q", privateValue)
	}
}

func assertRecorderUntouched(t *testing.T, rec *httptest.ResponseRecorder) {
	t.Helper()
	if rec.Body.Len() != 0 {
		t.Fatalf("response body was written before success: %q", rec.Body.String())
	}
	if len(rec.Header()) != 0 {
		t.Fatalf("response headers were written before success: %#v", rec.Header())
	}
}

func assertTrackingRecorderUntouched(t *testing.T, rec *trackingResponseWriter) {
	t.Helper()
	assertRecorderUntouched(t, rec.ResponseRecorder)
	if rec.wroteHeader {
		t.Fatal("response status was committed before success")
	}
	if rec.wroteBody {
		t.Fatal("response body write was attempted before success")
	}
}

func decodeResponse(t *testing.T, rec *httptest.ResponseRecorder, dst any) {
	t.Helper()
	if err := json.Unmarshal(rec.Body.Bytes(), dst); err != nil {
		t.Fatalf("response is not valid JSON: %v; body=%q", err, rec.Body.String())
	}
}

func assertBodyDoesNotContain(t *testing.T, rec *httptest.ResponseRecorder, privateValues ...string) {
	t.Helper()
	body := rec.Body.String()
	for _, privateValue := range privateValues {
		if strings.Contains(body, privateValue) {
			t.Fatalf("response body exposed private value %q: %s", privateValue, body)
		}
	}
}
