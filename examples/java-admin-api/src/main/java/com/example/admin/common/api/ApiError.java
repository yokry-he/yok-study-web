package com.example.admin.common.api;

import java.util.Map;

public record ApiError(String code, String message, Map<String, String> fields) {
    public ApiError {
        fields = fields == null ? Map.of() : Map.copyOf(fields);
    }
}
