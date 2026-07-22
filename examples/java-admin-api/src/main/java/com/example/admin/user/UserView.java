package com.example.admin.user;

import java.time.Instant;
import java.util.List;
import java.util.UUID;

public record UserView(
        UUID id,
        String email,
        String displayName,
        UserStatus status,
        List<String> roleCodes,
        long version,
        Instant createdAt,
        Instant updatedAt
) {
    static UserView from(UserEntity user, List<String> roleCodes) {
        return new UserView(
                user.getId(),
                user.getEmail(),
                user.getDisplayName(),
                user.getStatus(),
                roleCodes.stream().sorted().toList(),
                user.getVersion(),
                user.getCreatedAt(),
                user.getUpdatedAt()
        );
    }
}
