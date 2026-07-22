package httpx

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
	"time"
)

type failingReader struct {
	err error
}

func (r failingReader) Read([]byte) (int, error) {
	return 0, r.err
}

type closeTrackingBody struct {
	reader io.Reader
	closed bool
}

func (b *closeTrackingBody) Read(payload []byte) (int, error) {
	return b.reader.Read(payload)
}

func (b *closeTrackingBody) Close() error {
	b.closed = true
	return nil
}

type flushTrackingWriter struct {
	header  http.Header
	flushed bool
}

func (w *flushTrackingWriter) Header() http.Header {
	if w.header == nil {
		w.header = make(http.Header)
	}
	return w.header
}

func (*flushTrackingWriter) Write(payload []byte) (int, error) {
	return len(payload), nil
}

func (*flushTrackingWriter) WriteHeader(int) {}

func (w *flushTrackingWriter) Flush() {
	w.flushed = true
}

type headerSequenceWriter struct {
	header   http.Header
	statuses []int
	body     bytes.Buffer
}

func (w *headerSequenceWriter) Header() http.Header {
	if w.header == nil {
		w.header = make(http.Header)
	}
	return w.header
}

func (w *headerSequenceWriter) Write(payload []byte) (int, error) {
	return w.body.Write(payload)
}

func (w *headerSequenceWriter) WriteHeader(status int) {
	w.statuses = append(w.statuses, status)
}

func TestRequestIDKeepsValidCallerValue(t *testing.T) {
	values := []string{
		"a",
		"request-123_ABC.example:test",
		strings.Repeat("x", 128),
	}

	for _, value := range values {
		t.Run(value[:min(len(value), 20)], func(t *testing.T) {
			originalReader := requestIDReader
			requestIDReader = failingReader{err: errors.New("random source must not be used")}
			t.Cleanup(func() { requestIDReader = originalReader })

			var contextValue string
			handler := RequestID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				contextValue = RequestIDFromContext(r.Context())
				w.WriteHeader(http.StatusNoContent)
			}))

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodGet, "/request-id", nil)
			request.Header.Set("X-Request-ID", value)
			handler.ServeHTTP(recorder, request)

			if recorder.Code != http.StatusNoContent {
				t.Fatalf("status = %d, want %d", recorder.Code, http.StatusNoContent)
			}
			if got := recorder.Header().Get("X-Request-ID"); got != value {
				t.Fatalf("response request id = %q, want %q", got, value)
			}
			if contextValue != value {
				t.Fatalf("context request id = %q, want %q", contextValue, value)
			}
		})
	}
}

func TestRequestIDReplacesMissingAndInvalidValues(t *testing.T) {
	values := []string{
		"",
		"contains space",
		"contains/slash",
		"中文",
		strings.Repeat("x", 129),
	}
	pattern := regexp.MustCompile(`^req_[0-9a-f]{32}$`)

	for _, value := range values {
		name := value
		if name == "" {
			name = "missing"
		}
		t.Run(name[:min(len(name), 20)], func(t *testing.T) {
			originalReader := requestIDReader
			requestIDReader = bytes.NewReader(bytes.Repeat([]byte{0xab}, 16))
			t.Cleanup(func() { requestIDReader = originalReader })

			var contextValue string
			handler := RequestID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				contextValue = RequestIDFromContext(r.Context())
				w.WriteHeader(http.StatusNoContent)
			}))

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodGet, "/request-id", nil)
			if value != "" {
				request.Header.Set("X-Request-ID", value)
			}
			handler.ServeHTTP(recorder, request)

			responseValue := recorder.Header().Get("X-Request-ID")
			if !pattern.MatchString(responseValue) {
				t.Fatalf("generated request id = %q, want req_ plus 32 lowercase hex digits", responseValue)
			}
			if contextValue != responseValue {
				t.Fatalf("context request id = %q, response request id = %q", contextValue, responseValue)
			}
		})
	}
}

func TestRequestIDHandlesRandomSourceFailure(t *testing.T) {
	originalReader := requestIDReader
	requestIDReader = failingReader{err: errors.New("random source secret")}
	t.Cleanup(func() { requestIDReader = originalReader })

	called := false
	handler := RequestID(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		called = true
	}))
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/request-id", nil))

	if called {
		t.Fatal("downstream handler ran after request id generation failed")
	}
	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusInternalServerError)
	}
	assertErrorEnvelope(t, recorder, CodeInternalError, "")
	if strings.Contains(recorder.Body.String(), "random source secret") {
		t.Fatal("random source error leaked to response")
	}
}

func TestResponseRecorderTracksImplicitStatusAndBytes(t *testing.T) {
	base := httptest.NewRecorder()
	recorder := &responseRecorder{ResponseWriter: base}

	written, err := recorder.Write([]byte("hello"))
	if err != nil {
		t.Fatalf("write: %v", err)
	}
	if written != 5 || recorder.bytes != 5 {
		t.Fatalf("written = %d, tracked bytes = %d, want 5 and 5", written, recorder.bytes)
	}
	if recorder.statusCode() != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.statusCode(), http.StatusOK)
	}
	if _, ok := any(recorder).(http.Flusher); ok {
		t.Fatal("response recorder must not directly advertise http.Flusher")
	}
}

func TestResponseRecorderUnwrapLetsResponseControllerReachFlusher(t *testing.T) {
	for _, layers := range []int{1, 2} {
		t.Run(fmt.Sprintf("layers_%d", layers), func(t *testing.T) {
			base := &flushTrackingWriter{}
			var wrapped http.ResponseWriter = base
			for range layers {
				wrapped = &responseRecorder{ResponseWriter: wrapped}
			}

			if err := http.NewResponseController(wrapped).Flush(); err != nil {
				t.Fatalf("flush through %d recorder layers: %v", layers, err)
			}
			if !base.flushed {
				t.Fatalf("flush did not reach underlying writer through %d layers", layers)
			}
		})
	}
}

func TestResponseRecorderForwardsInformationalHeadersBeforeFinalStatus(t *testing.T) {
	base := &headerSequenceWriter{}
	recorder := &responseRecorder{ResponseWriter: base}

	recorder.WriteHeader(http.StatusContinue)
	recorder.WriteHeader(http.StatusEarlyHints)
	if recorder.committed() {
		t.Fatal("informational response unexpectedly committed final status")
	}
	recorder.WriteHeader(http.StatusCreated)
	recorder.WriteHeader(http.StatusNoContent)

	if got := fmt.Sprint(base.statuses); got != "[100 103 201]" {
		t.Fatalf("forwarded statuses = %s, want [100 103 201]", got)
	}
	if recorder.statusCode() != http.StatusCreated {
		t.Fatalf("final status = %d, want %d", recorder.statusCode(), http.StatusCreated)
	}
	if !recorder.committed() {
		t.Fatal("final response status was not marked committed")
	}
}

func TestRecoverCanWriteFinalErrorAfterEarlyHints(t *testing.T) {
	base := &headerSequenceWriter{}
	handler := Recover(nil)(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusEarlyHints)
		panic("secret after early hints")
	}))

	var recovered any
	func() {
		defer func() { recovered = recover() }()
		handler.ServeHTTP(base, httptest.NewRequest(http.MethodGet, "/early-hints", nil))
	}()

	if recovered != nil {
		t.Fatalf("recover propagated %#v after informational response", recovered)
	}
	if got := fmt.Sprint(base.statuses); got != "[103 500]" {
		t.Fatalf("forwarded statuses = %s, want [103 500]", got)
	}
	var envelope Envelope
	if err := json.Unmarshal(base.body.Bytes(), &envelope); err != nil {
		t.Fatalf("decode final error response: %v; body = %s", err, base.body.String())
	}
	if envelope.Success || envelope.Error == nil || envelope.Error.Code != CodeInternalError {
		t.Fatalf("unexpected final error envelope: %#v", envelope)
	}
}

func TestChainUsesDeclaredOuterToInnerOrder(t *testing.T) {
	var order []string
	middleware := func(name string) func(http.Handler) http.Handler {
		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				order = append(order, name+":before")
				next.ServeHTTP(w, r)
				order = append(order, name+":after")
			})
		}
	}
	final := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		order = append(order, "final")
		w.WriteHeader(http.StatusNoContent)
	})

	Chain(final, middleware("first"), middleware("second")).ServeHTTP(
		httptest.NewRecorder(),
		httptest.NewRequest(http.MethodGet, "/chain", nil),
	)

	want := "first:before,second:before,final,second:after,first:after"
	if got := strings.Join(order, ","); got != want {
		t.Fatalf("order = %q, want %q", got, want)
	}

	called := false
	Chain(http.HandlerFunc(func(http.ResponseWriter, *http.Request) { called = true })).ServeHTTP(
		httptest.NewRecorder(),
		httptest.NewRequest(http.MethodGet, "/empty-chain", nil),
	)
	if !called {
		t.Fatal("empty middleware chain did not call final handler")
	}
}

func TestAccessLogRecordsOneStructuredEntryWithoutSensitiveData(t *testing.T) {
	var logs bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&logs, nil))
	handler := Chain(
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte("ok"))
		}),
		RequestID,
		AccessLog(logger),
	)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/tasks?token=query-secret", strings.NewReader("body-secret"))
	request.Header.Set("X-Request-ID", "caller-request-1")
	request.Header.Set("Authorization", "Bearer authorization-secret")
	request.Header.Set("Cookie", "session=cookie-secret")
	handler.ServeHTTP(recorder, request)

	entries := decodeLogEntries(t, logs.String())
	if len(entries) != 1 {
		t.Fatalf("log entries = %d, want 1: %s", len(entries), logs.String())
	}
	entry := entries[0]
	assertLogValue(t, entry, "method", http.MethodPost)
	assertLogValue(t, entry, "path", "/tasks")
	assertLogValue(t, entry, "request_id", "caller-request-1")
	assertLogNumber(t, entry, "status", http.StatusCreated)
	assertLogNumber(t, entry, "bytes", 2)
	if _, ok := entry["duration"]; !ok {
		t.Fatal("access log missing duration")
	}
	for _, secret := range []string{"query-secret", "body-secret", "authorization-secret", "cookie-secret"} {
		if strings.Contains(logs.String(), secret) {
			t.Fatalf("access log leaked %q: %s", secret, logs.String())
		}
	}
}

func TestAccessLogAllowsNilLogger(t *testing.T) {
	called := false
	handler := AccessLog(nil)(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		called = true
		w.WriteHeader(http.StatusNoContent)
	}))
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/nil-logger", nil))

	if !called || recorder.Code != http.StatusNoContent {
		t.Fatalf("called = %t, status = %d", called, recorder.Code)
	}
}

func TestRecoverHidesPanicAndKeepsRequestID(t *testing.T) {
	var logs bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&logs, nil))
	handler := Chain(
		http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
			panic("database password=secret")
		}),
		RequestID,
		AccessLog(logger),
		Recover(logger),
	)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/panic?token=query-secret", nil)
	request.Header.Set("X-Request-ID", "panic-request-1")
	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusInternalServerError)
	}
	assertErrorEnvelope(t, recorder, CodeInternalError, "panic-request-1")
	for _, secret := range []string{"password=secret", "query-secret"} {
		if strings.Contains(recorder.Body.String(), secret) || strings.Contains(logs.String(), secret) {
			t.Fatalf("panic path leaked %q", secret)
		}
	}

	entries := decodeLogEntries(t, logs.String())
	accessEntries := 0
	panicEntries := 0
	for _, entry := range entries {
		switch entry["msg"] {
		case "http request":
			accessEntries++
			assertLogNumber(t, entry, "status", http.StatusInternalServerError)
		case "panic recovered":
			panicEntries++
			assertLogValue(t, entry, "request_id", "panic-request-1")
			if _, ok := entry["stack"]; !ok {
				t.Fatal("panic log missing stack")
			}
		}
	}
	if accessEntries != 1 || panicEntries != 1 {
		t.Fatalf("access entries = %d, panic entries = %d, logs = %s", accessEntries, panicEntries, logs.String())
	}
}

func TestRecoverDoesNotRewriteCommittedResponse(t *testing.T) {
	var logs bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&logs, nil))
	handler := Recover(logger)(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusAccepted)
		_, _ = w.Write([]byte("partial"))
		panic("sensitive panic value")
	}))

	recorder := httptest.NewRecorder()
	var recovered any
	func() {
		defer func() { recovered = recover() }()
		handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/committed", nil))
	}()

	if recovered != http.ErrAbortHandler {
		t.Fatalf("recovered value = %#v, want http.ErrAbortHandler", recovered)
	}
	if recorder.Code != http.StatusAccepted {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusAccepted)
	}
	if got := recorder.Body.String(); got != "partial" {
		t.Fatalf("body = %q, want partial", got)
	}
	if strings.Contains(logs.String(), "sensitive panic value") {
		t.Fatal("panic value leaked to log")
	}
}

func TestRecoverPreservesErrAbortHandlerSemantics(t *testing.T) {
	var logs bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&logs, nil))
	handler := Recover(logger)(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		panic(http.ErrAbortHandler)
	}))

	var recovered any
	func() {
		defer func() { recovered = recover() }()
		handler.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/abort", nil))
	}()

	if recovered != http.ErrAbortHandler {
		t.Fatalf("recovered value = %#v, want http.ErrAbortHandler", recovered)
	}
	if logs.Len() != 0 {
		t.Fatalf("ErrAbortHandler must not be logged: %s", logs.String())
	}
}

func TestDeadlinePassesDeadlineExceededToDownstream(t *testing.T) {
	observed := make(chan error, 1)
	handler := Deadline(5 * time.Millisecond)(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		<-r.Context().Done()
		observed <- r.Context().Err()
	}))

	handler.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/deadline", nil))
	select {
	case err := <-observed:
		if !errors.Is(err, context.DeadlineExceeded) {
			t.Fatalf("context error = %v, want deadline exceeded", err)
		}
	case <-time.After(time.Second):
		t.Fatal("downstream did not observe deadline")
	}
}

func TestDeadlineDoesNotDetachHandlerAndNonPositiveTimeoutPassesThrough(t *testing.T) {
	started := time.Now()
	Deadline(time.Millisecond)(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		time.Sleep(20 * time.Millisecond)
	})).ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/sync", nil))
	if elapsed := time.Since(started); elapsed < 18*time.Millisecond {
		t.Fatalf("middleware returned before downstream handler: %s", elapsed)
	}

	for _, timeout := range []time.Duration{0, -time.Second} {
		called := false
		Deadline(timeout)(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
			called = true
			if _, ok := r.Context().Deadline(); ok {
				t.Fatal("non-positive timeout unexpectedly added a deadline")
			}
		})).ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/passthrough", nil))
		if !called {
			t.Fatalf("timeout %s did not call downstream", timeout)
		}
	}
}

func TestRequireJSONSkipsReadOnlyMethods(t *testing.T) {
	for _, method := range []string{http.MethodGet, http.MethodHead, http.MethodDelete} {
		t.Run(method, func(t *testing.T) {
			called := false
			handler := RequireJSON(1024)(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				called = true
				w.WriteHeader(http.StatusNoContent)
			}))
			recorder := httptest.NewRecorder()
			handler.ServeHTTP(recorder, httptest.NewRequest(method, "/tasks", nil))
			if !called || recorder.Code != http.StatusNoContent {
				t.Fatalf("called = %t, status = %d", called, recorder.Code)
			}
		})
	}
}

func TestRequireJSONRejectsMissingMalformedAndUnsupportedMediaTypes(t *testing.T) {
	tests := []struct {
		name        string
		contentType string
	}{
		{name: "missing"},
		{name: "malformed", contentType: `application/json; charset="unterminated`},
		{name: "other media type", contentType: "text/plain"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			called := false
			handler := Chain(
				http.HandlerFunc(func(http.ResponseWriter, *http.Request) { called = true }),
				RequestID,
				RequireJSON(1024),
			)
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodPost, "/tasks", strings.NewReader(`{"title":"test"}`))
			request.Header.Set("X-Request-ID", "media-request-1")
			if test.contentType != "" {
				request.Header.Set("Content-Type", test.contentType)
			}
			handler.ServeHTTP(recorder, request)

			if called {
				t.Fatal("downstream ran for unsupported media type")
			}
			if recorder.Code != http.StatusUnsupportedMediaType {
				t.Fatalf("status = %d, want %d", recorder.Code, http.StatusUnsupportedMediaType)
			}
			assertErrorEnvelope(t, recorder, CodeUnsupportedMedia, "media-request-1")
		})
	}
}

func TestRequireJSONAcceptsParametersWithoutConsumingBody(t *testing.T) {
	body := `{"title":"learn Go"}`
	originalBody := &closeTrackingBody{reader: strings.NewReader(body)}
	var downstreamBody string
	handler := RequireJSON(1024)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		payload, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("read downstream body: %v", err)
		}
		downstreamBody = string(payload)
		w.WriteHeader(http.StatusNoContent)
	}))

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/tasks", nil)
	request.Body = originalBody
	request.ContentLength = int64(len(body))
	request.Header.Set("Content-Type", "application/json; charset=utf-8")
	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusNoContent)
	}
	if downstreamBody != body {
		t.Fatalf("downstream body = %q, want %q", downstreamBody, body)
	}
	if !originalBody.closed {
		t.Fatal("original request body was not closed after being restored for downstream")
	}
}

func TestRequireJSONRejectsOversizedBodyBeforeHandler(t *testing.T) {
	const limit = int64(1 << 20)
	originalBody := &closeTrackingBody{reader: strings.NewReader(strings.Repeat("x", int(limit)+1))}
	called := false
	handler := Chain(
		http.HandlerFunc(func(http.ResponseWriter, *http.Request) { called = true }),
		RequestID,
		RequireJSON(limit),
	)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/tasks", nil)
	request.Body = originalBody
	request.ContentLength = limit + 1
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Request-ID", "large-body-request-1")
	handler.ServeHTTP(recorder, request)

	if called {
		t.Fatal("downstream ran for oversized body")
	}
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusBadRequest)
	}
	assertErrorEnvelope(t, recorder, CodeBodyTooLarge, "large-body-request-1")
	if !originalBody.closed {
		t.Fatal("known-length oversized request body was not closed")
	}
}

func TestRequireJSONRejectsUnknownLengthOversizedBodyAndClosesIt(t *testing.T) {
	const limit = int64(1024)
	originalBody := &closeTrackingBody{reader: strings.NewReader(strings.Repeat("x", int(limit)+1))}
	called := false
	handler := Chain(
		http.HandlerFunc(func(http.ResponseWriter, *http.Request) { called = true }),
		RequestID,
		RequireJSON(limit),
	)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/tasks", nil)
	request.Body = originalBody
	request.ContentLength = -1
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Request-ID", "unknown-length-request")
	handler.ServeHTTP(recorder, request)

	if called {
		t.Fatal("downstream ran for unknown-length oversized body")
	}
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusBadRequest)
	}
	assertErrorEnvelope(t, recorder, CodeBodyTooLarge, "unknown-length-request")
	if !originalBody.closed {
		t.Fatal("unknown-length oversized request body was not closed")
	}
}

func TestRequireJSONReturnsStableErrorAndClosesBodyWhenReadFails(t *testing.T) {
	originalBody := &closeTrackingBody{reader: failingReader{err: errors.New("transport password=secret")}}
	called := false
	handler := Chain(
		http.HandlerFunc(func(http.ResponseWriter, *http.Request) { called = true }),
		RequestID,
		RequireJSON(1024),
	)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPatch, "/tasks/1", nil)
	request.Body = originalBody
	request.ContentLength = -1
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Request-ID", "read-failure-request")
	handler.ServeHTTP(recorder, request)

	if called {
		t.Fatal("downstream ran after request body read failed")
	}
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusBadRequest)
	}
	assertErrorEnvelope(t, recorder, CodeInvalidJSON, "read-failure-request")
	if strings.Contains(recorder.Body.String(), "transport password=secret") {
		t.Fatal("request body read error leaked to response")
	}
	if !originalBody.closed {
		t.Fatal("request body was not closed after read failure")
	}
}

func TestRequireJSONHandlesNonPositiveLimit(t *testing.T) {
	for _, limit := range []int64{0, -1} {
		called := false
		handler := Chain(
			http.HandlerFunc(func(http.ResponseWriter, *http.Request) { called = true }),
			RequestID,
			RequireJSON(limit),
		)
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodPut, "/tasks/1", strings.NewReader(`{"title":"test"}`))
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("X-Request-ID", "invalid-limit-request")
		handler.ServeHTTP(recorder, request)

		if called {
			t.Fatalf("downstream ran for limit %d", limit)
		}
		if recorder.Code != http.StatusInternalServerError {
			t.Fatalf("limit %d status = %d, want %d", limit, recorder.Code, http.StatusInternalServerError)
		}
		assertErrorEnvelope(t, recorder, CodeInternalError, "invalid-limit-request")
	}
}

func TestFixedMiddlewareComposition(t *testing.T) {
	var logs bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&logs, nil))
	body := `{"title":"ship middleware"}`
	var gotBody string
	router := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		payload, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("read body: %v", err)
		}
		gotBody = string(payload)
		w.WriteHeader(http.StatusCreated)
	})
	handler := Chain(
		router,
		RequestID,
		AccessLog(logger),
		Recover(logger),
		Deadline(100*time.Millisecond),
		RequireJSON(1<<20),
	)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/tasks", strings.NewReader(body))
	request.Header.Set("Content-Type", "application/json; charset=utf-8")
	request.Header.Set("X-Request-ID", "fixed-chain-request")
	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusCreated)
	}
	if gotBody != body {
		t.Fatalf("body = %q, want %q", gotBody, body)
	}
	if got := recorder.Header().Get("X-Request-ID"); got != "fixed-chain-request" {
		t.Fatalf("response request id = %q", got)
	}
	entries := decodeLogEntries(t, logs.String())
	if len(entries) != 1 || entries[0]["msg"] != "http request" {
		t.Fatalf("unexpected access logs: %s", logs.String())
	}
}

func assertErrorEnvelope(t *testing.T, recorder *httptest.ResponseRecorder, code, requestID string) {
	t.Helper()
	var envelope Envelope
	if err := json.Unmarshal(recorder.Body.Bytes(), &envelope); err != nil {
		t.Fatalf("decode response envelope: %v; body = %s", err, recorder.Body.String())
	}
	if envelope.Success {
		t.Fatal("error envelope marked successful")
	}
	if envelope.Error == nil || envelope.Error.Code != code {
		t.Fatalf("error = %#v, want code %q", envelope.Error, code)
	}
	if envelope.RequestID != requestID {
		t.Fatalf("request id = %q, want %q", envelope.RequestID, requestID)
	}
}

func decodeLogEntries(t *testing.T, raw string) []map[string]any {
	t.Helper()
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return nil
	}
	lines := strings.Split(trimmed, "\n")
	entries := make([]map[string]any, 0, len(lines))
	for _, line := range lines {
		var entry map[string]any
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			t.Fatalf("decode log entry: %v; line = %s", err, line)
		}
		entries = append(entries, entry)
	}
	return entries
}

func assertLogValue(t *testing.T, entry map[string]any, key, want string) {
	t.Helper()
	if got, _ := entry[key].(string); got != want {
		t.Fatalf("log %s = %q, want %q; entry = %#v", key, got, want, entry)
	}
}

func assertLogNumber(t *testing.T, entry map[string]any, key string, want int) {
	t.Helper()
	got, ok := entry[key].(float64)
	if !ok || int(got) != want {
		t.Fatalf("log %s = %#v, want %d; entry = %#v", key, entry[key], want, entry)
	}
}
