package com.example.admin.user;

import java.time.Instant;
import java.util.LinkedHashSet;
import java.util.Set;
import java.util.UUID;

import jakarta.persistence.Column;
import jakarta.persistence.Entity;
import jakarta.persistence.EnumType;
import jakarta.persistence.Enumerated;
import jakarta.persistence.FetchType;
import jakarta.persistence.Id;
import jakarta.persistence.JoinColumn;
import jakarta.persistence.JoinTable;
import jakarta.persistence.ManyToMany;
import jakarta.persistence.PrePersist;
import jakarta.persistence.PreUpdate;
import jakarta.persistence.Table;
import jakarta.persistence.Version;

import com.example.admin.role.RoleEntity;

@Entity
@Table(name = "users")
public class UserEntity {
    @Id
    private UUID id;

    @Column(nullable = false, length = 254)
    private String email;

    @Column(name = "display_name", nullable = false, length = 80)
    private String displayName;

    @Enumerated(EnumType.STRING)
    @Column(nullable = false, length = 20)
    private UserStatus status;

    @Version
    @Column(nullable = false)
    private long version;

    @ManyToMany(fetch = FetchType.LAZY)
    @JoinTable(
            name = "user_roles",
            joinColumns = @JoinColumn(name = "user_id"),
            inverseJoinColumns = @JoinColumn(name = "role_id")
    )
    private Set<RoleEntity> roles = new LinkedHashSet<>();

    @Column(name = "created_at", nullable = false)
    private Instant createdAt;

    @Column(name = "updated_at", nullable = false)
    private Instant updatedAt;

    protected UserEntity() {
    }

    private UserEntity(UUID id, String email, String displayName) {
        this.id = id;
        this.email = email;
        this.displayName = displayName;
        this.status = UserStatus.ACTIVE;
    }

    public static UserEntity create(String email, String displayName) {
        return new UserEntity(UUID.randomUUID(), email, displayName);
    }

    public void updateProfile(String email, String displayName) {
        this.email = email;
        this.displayName = displayName;
    }

    public void changeStatus(UserStatus status) {
        this.status = status;
    }

    public void replaceRoles(Set<RoleEntity> roles) {
        this.roles.clear();
        this.roles.addAll(roles);
        this.updatedAt = Instant.now();
    }

    @PrePersist
    void beforeInsert() {
        Instant now = Instant.now();
        this.createdAt = now;
        this.updatedAt = now;
    }

    @PreUpdate
    void beforeUpdate() {
        this.updatedAt = Instant.now();
    }

    public UUID getId() {
        return id;
    }

    public String getEmail() {
        return email;
    }

    public String getDisplayName() {
        return displayName;
    }

    public UserStatus getStatus() {
        return status;
    }

    public long getVersion() {
        return version;
    }

    public Set<RoleEntity> getRoles() {
        return Set.copyOf(roles);
    }

    public Instant getCreatedAt() {
        return createdAt;
    }

    public Instant getUpdatedAt() {
        return updatedAt;
    }
}
