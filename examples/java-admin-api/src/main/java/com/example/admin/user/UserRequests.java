package com.example.admin.user;

import java.util.Set;
import java.util.UUID;

import jakarta.validation.constraints.Email;
import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.NotEmpty;
import jakarta.validation.constraints.NotNull;
import jakarta.validation.constraints.PositiveOrZero;
import jakarta.validation.constraints.Size;

public final class UserRequests {
    private UserRequests() {
    }

    public record Create(
            @NotBlank @Email @Size(max = 254) String email,
            @NotBlank @Size(min = 2, max = 80) String displayName,
            @NotEmpty @Size(max = 20) Set<UUID> roleIds
    ) {
    }

    public record Update(
            @NotBlank @Email @Size(max = 254) String email,
            @NotBlank @Size(min = 2, max = 80) String displayName,
            @NotNull @PositiveOrZero Long expectedVersion
    ) {
    }

    public record ChangeStatus(
            @NotNull UserStatus status,
            @NotNull @PositiveOrZero Long expectedVersion
    ) {
    }

    public record ReplaceRoles(
            @NotEmpty @Size(max = 20) Set<UUID> roleIds,
            @NotNull @PositiveOrZero Long expectedVersion
    ) {
    }
}
