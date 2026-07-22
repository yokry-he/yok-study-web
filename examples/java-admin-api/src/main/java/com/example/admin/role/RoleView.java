package com.example.admin.role;

import java.util.UUID;

public record RoleView(UUID id, String code, String name) {
    static RoleView from(RoleEntity role) {
        return new RoleView(role.getId(), role.getCode(), role.getName());
    }
}
