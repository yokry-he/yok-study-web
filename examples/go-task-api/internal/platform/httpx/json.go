package httpx

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

// 稳定 sentinel 让调用方通过 errors.Is 判断解码失败类别。
var (
	ErrEmptyJSONBody      = errors.New("empty JSON body")
	ErrMalformedJSON      = errors.New("malformed JSON")
	ErrUnknownJSONField   = errors.New("unknown JSON field")
	ErrInvalidJSONType    = errors.New("invalid JSON type")
	ErrBodyTooLarge       = errors.New("JSON body too large")
	ErrMultipleJSONValues = errors.New("multiple JSON values")
)

func DecodeJSON(w http.ResponseWriter, r *http.Request, dst any, maxBytes int64) error {
	if maxBytes <= 0 {
		return WrapAPIError(
			http.StatusInternalServerError,
			CodeInternalError,
			"服务器内部错误",
			nil,
			fmt.Errorf("DecodeJSON maxBytes must be positive: %d", maxBytes),
		)
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxBytes)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(dst); err != nil {
		return classifyDecodeError(err)
	}

	var extra any
	if err := decoder.Decode(&extra); err != nil {
		if errors.Is(err, io.EOF) {
			return nil
		}
		return classifyDecodeError(err)
	}

	return WrapAPIError(
		http.StatusBadRequest,
		CodeInvalidJSON,
		"请求体只能包含一个 JSON 值",
		nil,
		ErrMultipleJSONValues,
	)
}

func ParsePositiveID(raw string) (int64, error) {
	id, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || id <= 0 {
		return 0, NewAPIError(
			http.StatusBadRequest,
			CodeInvalidArgument,
			"资源 ID 必须是正整数",
			nil,
		)
	}
	return id, nil
}

func classifyDecodeError(err error) error {
	var maxBytesErr *http.MaxBytesError
	if errors.As(err, &maxBytesErr) {
		return WrapAPIError(
			http.StatusBadRequest,
			CodeBodyTooLarge,
			"请求体超过大小限制",
			nil,
			errors.Join(ErrBodyTooLarge, err),
		)
	}

	if errors.Is(err, io.EOF) {
		return WrapAPIError(
			http.StatusBadRequest,
			CodeInvalidJSON,
			"请求体不能为空",
			nil,
			ErrEmptyJSONBody,
		)
	}

	var syntaxErr *json.SyntaxError
	if errors.As(err, &syntaxErr) || errors.Is(err, io.ErrUnexpectedEOF) {
		return WrapAPIError(
			http.StatusBadRequest,
			CodeInvalidJSON,
			"请求体包含无效 JSON",
			nil,
			errors.Join(ErrMalformedJSON, err),
		)
	}

	if field, ok := unknownJSONField(err); ok {
		message := "请求体包含未知字段"
		if field != "" {
			message = fmt.Sprintf("请求体包含未知字段 %q", field)
		}
		return WrapAPIError(
			http.StatusBadRequest,
			CodeUnknownField,
			message,
			nil,
			errors.Join(ErrUnknownJSONField, err),
		)
	}

	var typeErr *json.UnmarshalTypeError
	if errors.As(err, &typeErr) {
		var fields FieldErrors
		if typeErr.Field != "" {
			fields = FieldErrors{typeErr.Field: "字段类型不正确"}
		}
		return WrapAPIError(
			http.StatusBadRequest,
			CodeInvalidArgument,
			"请求参数类型不正确",
			fields,
			errors.Join(ErrInvalidJSONType, err),
		)
	}

	var invalidTarget *json.InvalidUnmarshalError
	if errors.As(err, &invalidTarget) {
		return WrapAPIError(
			http.StatusInternalServerError,
			CodeInternalError,
			"服务器内部错误",
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

func unknownJSONField(err error) (string, bool) {
	const prefix = "json: unknown field "
	quoted, ok := strings.CutPrefix(err.Error(), prefix)
	if !ok {
		return "", false
	}

	field, unquoteErr := strconv.Unquote(quoted)
	if unquoteErr != nil {
		return "", true
	}
	return field, true
}
