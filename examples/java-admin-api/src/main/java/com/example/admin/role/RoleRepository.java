package com.example.admin.role;

import java.util.Collection;
import java.util.List;
import java.util.UUID;

import org.springframework.data.jpa.repository.JpaRepository;

public interface RoleRepository extends JpaRepository<RoleEntity, UUID> {
    List<RoleEntity> findAllByStatusOrderByNameAsc(String status);

    List<RoleEntity> findAllByIdInAndStatus(Collection<UUID> ids, String status);
}
