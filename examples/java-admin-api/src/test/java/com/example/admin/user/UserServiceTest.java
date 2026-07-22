package com.example.admin.user;

import java.util.Set;
import java.util.UUID;

import org.junit.jupiter.api.Test;

import com.example.admin.common.error.AppException;
import com.example.admin.role.RoleRepository;

import static org.assertj.core.api.Assertions.assertThatThrownBy;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.verifyNoInteractions;
import static org.mockito.Mockito.when;

class UserServiceTest {
    private final UserRepository users = mock(UserRepository.class);
    private final RoleRepository roles = mock(RoleRepository.class);
    private final UserService service = new UserService(users, roles);

    @Test
    void rejectsDuplicateEmailBeforeWriting() {
        UserRequests.Create request = new UserRequests.Create(
                "Ada@Example.com",
                "Ada",
                Set.of(UUID.randomUUID())
        );
        when(users.existsByEmailIgnoreCase("ada@example.com")).thenReturn(true);

        assertThatThrownBy(() -> service.create(request))
                .isInstanceOf(AppException.class)
                .extracting("code")
                .isEqualTo("USER_EMAIL_EXISTS");
        verifyNoInteractions(roles);
    }

    @Test
    void rejectsUnknownOrDisabledRoles() {
        UUID roleId = UUID.randomUUID();
        UserRequests.Create request = new UserRequests.Create(
                "ada@example.com",
                "Ada",
                Set.of(roleId)
        );
        when(users.existsByEmailIgnoreCase("ada@example.com")).thenReturn(false);
        when(roles.findAllByIdInAndStatus(Set.of(roleId), "ACTIVE")).thenReturn(java.util.List.of());

        assertThatThrownBy(() -> service.create(request))
                .isInstanceOf(AppException.class)
                .extracting("code")
                .isEqualTo("INVALID_ROLE");
    }
}
