package com.example.admin.common.error;

import java.util.LinkedHashMap;
import java.util.Map;

import jakarta.servlet.http.HttpServletRequest;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.dao.DataIntegrityViolationException;
import org.springframework.dao.OptimisticLockingFailureException;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.http.converter.HttpMessageNotReadableException;
import org.springframework.web.HttpMediaTypeNotSupportedException;
import org.springframework.web.HttpRequestMethodNotSupportedException;
import org.springframework.web.bind.MethodArgumentNotValidException;
import org.springframework.web.bind.annotation.ExceptionHandler;
import org.springframework.web.bind.annotation.RestControllerAdvice;
import org.springframework.web.method.annotation.HandlerMethodValidationException;
import org.springframework.web.method.annotation.MethodArgumentTypeMismatchException;
import org.springframework.web.servlet.resource.NoResourceFoundException;

import com.example.admin.common.api.ApiError;
import com.example.admin.common.api.ApiResponse;
import com.example.admin.common.web.RequestIdFilter;

@RestControllerAdvice
public final class GlobalExceptionHandler {
    private static final Logger log = LoggerFactory.getLogger(GlobalExceptionHandler.class);

    @ExceptionHandler(AppException.class)
    ResponseEntity<ApiResponse<Void>> handleApp(
            AppException exception,
            HttpServletRequest request
    ) {
        return failure(exception.status(), exception.code(), exception.getMessage(), Map.of(), request);
    }

    @ExceptionHandler(MethodArgumentNotValidException.class)
    ResponseEntity<ApiResponse<Void>> handleBodyValidation(
            MethodArgumentNotValidException exception,
            HttpServletRequest request
    ) {
        Map<String, String> fields = new LinkedHashMap<>();
        exception.getBindingResult().getFieldErrors().forEach(error ->
                fields.putIfAbsent(error.getField(), error.getDefaultMessage()));
        return failure(HttpStatus.BAD_REQUEST, "VALIDATION_ERROR", "请求参数不正确", fields, request);
    }

    @ExceptionHandler(HandlerMethodValidationException.class)
    ResponseEntity<ApiResponse<Void>> handleMethodValidation(
            HandlerMethodValidationException exception,
            HttpServletRequest request
    ) {
        return failure(HttpStatus.BAD_REQUEST, "VALIDATION_ERROR", "路径或查询参数不正确", Map.of(), request);
    }

    @ExceptionHandler(MethodArgumentTypeMismatchException.class)
    ResponseEntity<ApiResponse<Void>> handleTypeMismatch(
            MethodArgumentTypeMismatchException exception,
            HttpServletRequest request
    ) {
        return failure(HttpStatus.BAD_REQUEST, "INVALID_ARGUMENT", "路径或查询参数类型不正确", Map.of(), request);
    }

    @ExceptionHandler(HttpMessageNotReadableException.class)
    ResponseEntity<ApiResponse<Void>> handleUnreadable(
            HttpMessageNotReadableException exception,
            HttpServletRequest request
    ) {
        return failure(HttpStatus.BAD_REQUEST, "INVALID_JSON", "请求 JSON 无法解析", Map.of(), request);
    }

    @ExceptionHandler(HttpRequestMethodNotSupportedException.class)
    ResponseEntity<ApiResponse<Void>> handleMethodNotAllowed(
            HttpRequestMethodNotSupportedException exception,
            HttpServletRequest request
    ) {
        return failure(HttpStatus.METHOD_NOT_ALLOWED, "METHOD_NOT_ALLOWED", "请求方法不受支持", Map.of(), request);
    }

    @ExceptionHandler(HttpMediaTypeNotSupportedException.class)
    ResponseEntity<ApiResponse<Void>> handleUnsupportedMediaType(
            HttpMediaTypeNotSupportedException exception,
            HttpServletRequest request
    ) {
        return failure(HttpStatus.UNSUPPORTED_MEDIA_TYPE, "UNSUPPORTED_MEDIA_TYPE", "请求 Content-Type 不受支持", Map.of(), request);
    }

    @ExceptionHandler(OptimisticLockingFailureException.class)
    ResponseEntity<ApiResponse<Void>> handleOptimisticLock(
            OptimisticLockingFailureException exception,
            HttpServletRequest request
    ) {
        return failure(HttpStatus.CONFLICT, "STALE_VERSION", "数据已被其他请求修改，请刷新后重试", Map.of(), request);
    }

    @ExceptionHandler(DataIntegrityViolationException.class)
    ResponseEntity<ApiResponse<Void>> handleConstraint(
            DataIntegrityViolationException exception,
            HttpServletRequest request
    ) {
        return failure(HttpStatus.CONFLICT, "DATA_CONFLICT", "数据与现有记录冲突", Map.of(), request);
    }

    @ExceptionHandler(NoResourceFoundException.class)
    ResponseEntity<ApiResponse<Void>> handleNotFound(
            NoResourceFoundException exception,
            HttpServletRequest request
    ) {
        return failure(HttpStatus.NOT_FOUND, "ROUTE_NOT_FOUND", "请求地址不存在", Map.of(), request);
    }

    @ExceptionHandler(Exception.class)
    ResponseEntity<ApiResponse<Void>> handleUnexpected(
            Exception exception,
            HttpServletRequest request
    ) {
        log.error("Unhandled request failure", exception);
        return failure(HttpStatus.INTERNAL_SERVER_ERROR, "INTERNAL_ERROR", "服务暂时不可用", Map.of(), request);
    }

    private ResponseEntity<ApiResponse<Void>> failure(
            HttpStatus status,
            String code,
            String message,
            Map<String, String> fields,
            HttpServletRequest request
    ) {
        String requestId = RequestIdFilter.get(request);
        return ResponseEntity.status(status)
                .body(ApiResponse.failure(new ApiError(code, message, fields), requestId));
    }
}
