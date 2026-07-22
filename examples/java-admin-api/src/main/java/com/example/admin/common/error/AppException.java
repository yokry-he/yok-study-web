package com.example.admin.common.error;

import org.springframework.http.HttpStatus;

public final class AppException extends RuntimeException {
    private final HttpStatus status;
    private final String code;

    private AppException(HttpStatus status, String code, String message) {
        super(message);
        this.status = status;
        this.code = code;
    }

    public static AppException badRequest(String code, String message) {
        return new AppException(HttpStatus.BAD_REQUEST, code, message);
    }

    public static AppException notFound(String code, String message) {
        return new AppException(HttpStatus.NOT_FOUND, code, message);
    }

    public static AppException conflict(String code, String message) {
        return new AppException(HttpStatus.CONFLICT, code, message);
    }

    public HttpStatus status() {
        return status;
    }

    public String code() {
        return code;
    }
}
