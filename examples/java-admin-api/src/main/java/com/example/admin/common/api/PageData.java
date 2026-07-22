package com.example.admin.common.api;

import java.util.List;

public record PageData<T>(List<T> items, int page, int pageSize, long total) {
    public PageData {
        items = List.copyOf(items);
    }
}
