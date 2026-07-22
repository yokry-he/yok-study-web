package httpx

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"mime"
	"net/http"
	"runtime/debug"
	"strings"
	"time"
)

const (
	requestIDHeader      = "X-Request-ID"
	requestIDRandomBytes = 16
	requestIDMaxLength   = 128
)

type requestIDKey struct{}

// requestIDReader 默认使用密码学安全随机源，包内变量仅用于覆盖不可用随机源的失败路径。
var requestIDReader io.Reader = rand.Reader

func RequestIDFromContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	value, _ := ctx.Value(requestIDKey{}).(string)
	return value
}

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get(requestIDHeader)
		if !validRequestID(requestID) {
			generated, err := generateRequestID()
			if err != nil {
				_ = WriteError(w, WrapAPIError(
					http.StatusInternalServerError,
					CodeInternalError,
					"服务器内部错误",
					nil,
					fmt.Errorf("generate request id: %w", err),
				), "")
				return
			}
			requestID = generated
		}

		w.Header().Set(requestIDHeader, requestID)
		ctx := context.WithValue(r.Context(), requestIDKey{}, requestID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func validRequestID(value string) bool {
	if len(value) == 0 || len(value) > requestIDMaxLength {
		return false
	}
	for i := 0; i < len(value); i++ {
		char := value[i]
		if (char >= 'a' && char <= 'z') ||
			(char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9') {
			continue
		}
		switch char {
		case '.', '_', ':', '-':
			continue
		default:
			return false
		}
	}
	return true
}

func generateRequestID() (string, error) {
	random := make([]byte, requestIDRandomBytes)
	if _, err := io.ReadFull(requestIDReader, random); err != nil {
		return "", err
	}
	return "req_" + hex.EncodeToString(random), nil
}

type responseRecorder struct {
	http.ResponseWriter
	status int
	bytes  int
}

// Unwrap 供 http.ResponseController 查找底层真实支持的可选能力。
func (r *responseRecorder) Unwrap() http.ResponseWriter {
	return r.ResponseWriter
}

func (r *responseRecorder) WriteHeader(status int) {
	if r.status != 0 {
		return
	}
	if status >= 100 && status <= 199 {
		r.ResponseWriter.WriteHeader(status)
		return
	}
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

func (r *responseRecorder) Write(payload []byte) (int, error) {
	if r.status == 0 {
		r.WriteHeader(http.StatusOK)
	}
	written, err := r.ResponseWriter.Write(payload)
	r.bytes += written
	return written, err
}

func (r *responseRecorder) statusCode() int {
	if r.status == 0 {
		return http.StatusOK
	}
	return r.status
}

func (r *responseRecorder) committed() bool {
	return r.status != 0
}

func Chain(final http.Handler, middleware ...func(http.Handler) http.Handler) http.Handler {
	for i := len(middleware) - 1; i >= 0; i-- {
		final = middleware[i](final)
	}
	return final
}

func AccessLog(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		if logger == nil {
			return next
		}
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			startedAt := time.Now()
			recorder := &responseRecorder{ResponseWriter: w}
			defer func() {
				logger.InfoContext(
					r.Context(),
					"http request",
					"method", r.Method,
					"path", r.URL.Path,
					"status", recorder.statusCode(),
					"bytes", recorder.bytes,
					"duration", time.Since(startedAt),
					"request_id", RequestIDFromContext(r.Context()),
				)
			}()
			next.ServeHTTP(recorder, r)
		})
	}
}

func Recover(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			recorder := &responseRecorder{ResponseWriter: w}
			defer func() {
				recovered := recover()
				if recovered == nil {
					return
				}

				if panicErr, ok := recovered.(error); ok && errors.Is(panicErr, http.ErrAbortHandler) {
					panic(http.ErrAbortHandler)
				}

				requestID := RequestIDFromContext(r.Context())
				if logger != nil {
					// 不记录 panic 值；其字符串可能包含密码、SQL 或用户输入。
					logger.ErrorContext(
						r.Context(),
						"panic recovered",
						"method", r.Method,
						"path", r.URL.Path,
						"request_id", requestID,
						"panic_type", fmt.Sprintf("%T", recovered),
						"stack", string(debug.Stack()),
					)
				}

				if recorder.committed() {
					panic(http.ErrAbortHandler)
				}
				_ = WriteError(recorder, NewAPIError(
					http.StatusInternalServerError,
					CodeInternalError,
					"服务器内部错误",
					nil,
				), requestID)
			}()

			next.ServeHTTP(recorder, r)
		})
	}
}

func Deadline(timeout time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		if timeout <= 0 {
			return next
		}
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), timeout)
			defer cancel()
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RequireJSON(maxBytes int64) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !methodHasJSONBody(r.Method) {
				next.ServeHTTP(w, r)
				return
			}

			requestID := RequestIDFromContext(r.Context())
			if maxBytes <= 0 {
				_ = WriteError(w, WrapAPIError(
					http.StatusInternalServerError,
					CodeInternalError,
					"服务器内部错误",
					nil,
					fmt.Errorf("RequireJSON maxBytes must be positive: %d", maxBytes),
				), requestID)
				return
			}

			mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
			if err != nil || !strings.EqualFold(mediaType, "application/json") {
				_ = WriteError(w, NewAPIError(
					http.StatusUnsupportedMediaType,
					CodeUnsupportedMedia,
					"Content-Type 必须是 application/json",
					nil,
				), requestID)
				return
			}

			if r.ContentLength > maxBytes {
				_ = r.Body.Close()
				writeBodyTooLarge(w, requestID)
				return
			}

			limitedBody := http.MaxBytesReader(w, r.Body, maxBytes)
			payload, readErr := io.ReadAll(limitedBody)
			closeErr := limitedBody.Close()
			if readErr == nil {
				readErr = closeErr
			}
			if readErr != nil {
				var maxBytesErr *http.MaxBytesError
				if errors.As(readErr, &maxBytesErr) {
					writeBodyTooLarge(w, requestID)
					return
				}
				_ = WriteError(w, WrapAPIError(
					http.StatusBadRequest,
					CodeInvalidJSON,
					"无法读取请求体",
					nil,
					readErr,
				), requestID)
				return
			}

			r.Body = io.NopCloser(bytes.NewReader(payload))
			next.ServeHTTP(w, r)
		})
	}
}

func methodHasJSONBody(method string) bool {
	switch method {
	case http.MethodPost, http.MethodPut, http.MethodPatch:
		return true
	default:
		return false
	}
}

func writeBodyTooLarge(w http.ResponseWriter, requestID string) {
	_ = WriteError(w, WrapAPIError(
		http.StatusBadRequest,
		CodeBodyTooLarge,
		"请求体超过大小限制",
		nil,
		ErrBodyTooLarge,
	), requestID)
}
