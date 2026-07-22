package com.example.admin.user;

import java.util.Collection;
import java.util.List;
import java.util.Optional;
import java.util.UUID;

import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.data.jpa.repository.EntityGraph;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.query.Param;

public interface UserRepository extends JpaRepository<UserEntity, UUID> {
    boolean existsByEmailIgnoreCase(String email);

    boolean existsByEmailIgnoreCaseAndIdNot(String email, UUID id);

    @Query(value = """
            select u from UserEntity u
            where :query = ''
               or lower(u.email) like concat('%', lower(:query), '%')
               or lower(u.displayName) like concat('%', lower(:query), '%')
            """, countQuery = """
            select count(u) from UserEntity u
            where :query = ''
               or lower(u.email) like concat('%', lower(:query), '%')
               or lower(u.displayName) like concat('%', lower(:query), '%')
            """)
    Page<UserEntity> search(@Param("query") String query, Pageable pageable);

    @EntityGraph(attributePaths = "roles")
    Optional<UserEntity> findDetailById(UUID id);

    @Query("""
            select u.id as userId, r.code as roleCode
            from UserEntity u join u.roles r
            where u.id in :userIds
            """)
    List<UserRoleCodeRow> findRoleCodes(@Param("userIds") Collection<UUID> userIds);
}
