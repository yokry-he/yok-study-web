package com.example.admin;

import java.util.UUID;

import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.context.annotation.Import;
import org.springframework.dao.DataIntegrityViolationException;
import org.springframework.dao.OptimisticLockingFailureException;
import org.springframework.jdbc.core.JdbcTemplate;

import com.example.admin.user.UserEntity;
import com.example.admin.user.UserRepository;

import static org.junit.jupiter.api.Assertions.assertThrows;

@Import(TestcontainersConfiguration.class)
@SpringBootTest
class UserPersistenceIntegrationTest {
    @Autowired
    private JdbcTemplate jdbcTemplate;

    @Autowired
    private UserRepository userRepository;

    @Test
    void databaseEnforcesCaseInsensitiveEmailUniqueness() {
        insertUser("CaseSensitive@example.com");

        assertThrows(
                DataIntegrityViolationException.class,
                () -> insertUser("casesensitive@example.com")
        );
    }

    @Test
    void jpaVersionRejectsASecondDetachedUpdate() {
        UserEntity saved = userRepository.saveAndFlush(
                UserEntity.create("optimistic-lock@example.com", "Initial Name")
        );
        UserEntity firstSnapshot = userRepository.findById(saved.getId()).orElseThrow();
        UserEntity secondSnapshot = userRepository.findById(saved.getId()).orElseThrow();

        firstSnapshot.updateProfile(firstSnapshot.getEmail(), "First Writer");
        userRepository.saveAndFlush(firstSnapshot);

        secondSnapshot.updateProfile(secondSnapshot.getEmail(), "Second Writer");
        assertThrows(
                OptimisticLockingFailureException.class,
                () -> userRepository.saveAndFlush(secondSnapshot)
        );
    }

    private void insertUser(String email) {
        jdbcTemplate.update("""
                insert into users (id, email, display_name, status, version, created_at, updated_at)
                values (?, ?, ?, 'ACTIVE', 0, now(), now())
                """, UUID.randomUUID(), email, "Constraint Check");
    }
}
