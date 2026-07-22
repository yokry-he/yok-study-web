package com.example.admin.user;

import java.util.ArrayList;
import java.util.Collection;
import java.util.LinkedHashMap;
import java.util.LinkedHashSet;
import java.util.List;
import java.util.Locale;
import java.util.Map;
import java.util.Set;
import java.util.UUID;

import org.springframework.data.domain.Page;
import org.springframework.data.domain.PageRequest;
import org.springframework.data.domain.Sort;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import com.example.admin.common.api.PageData;
import com.example.admin.common.error.AppException;
import com.example.admin.role.RoleEntity;
import com.example.admin.role.RoleRepository;

@Service
public class UserService {
    private final UserRepository userRepository;
    private final RoleRepository roleRepository;

    public UserService(UserRepository userRepository, RoleRepository roleRepository) {
        this.userRepository = userRepository;
        this.roleRepository = roleRepository;
    }

    @Transactional(readOnly = true)
    public PageData<UserView> list(String query, int page, int pageSize) {
        String normalizedQuery = query == null ? "" : query.trim();
        PageRequest pageable = PageRequest.of(
                page,
                pageSize,
                Sort.by(Sort.Order.desc("createdAt"), Sort.Order.asc("id"))
        );
        Page<UserEntity> result = userRepository.search(normalizedQuery, pageable);
        Map<UUID, List<String>> roleCodes = loadRoleCodes(result.getContent());
        List<UserView> items = result.getContent().stream()
                .map(user -> UserView.from(user, roleCodes.getOrDefault(user.getId(), List.of())))
                .toList();
        return new PageData<>(items, page, pageSize, result.getTotalElements());
    }

    @Transactional(readOnly = true)
    public UserView get(UUID id) {
        UserEntity user = findDetail(id);
        List<String> roleCodes = user.getRoles().stream().map(RoleEntity::getCode).toList();
        return UserView.from(user, roleCodes);
    }

    @Transactional
    public UserView create(UserRequests.Create request) {
        String email = normalizeEmail(request.email());
        if (userRepository.existsByEmailIgnoreCase(email)) {
            throw AppException.conflict("USER_EMAIL_EXISTS", "邮箱已存在");
        }
        Set<RoleEntity> roles = loadActiveRoles(request.roleIds());
        UserEntity user = UserEntity.create(email, request.displayName().trim());
        user.replaceRoles(roles);
        // 手工 UUID 会让 Spring Data 使用 merge；必须保留返回的托管实例，
        // 才能把 flush 后的版本号和时间戳准确返回给客户端。
        user = userRepository.saveAndFlush(user);
        return UserView.from(user, roles.stream().map(RoleEntity::getCode).toList());
    }

    @Transactional
    public UserView update(UUID id, UserRequests.Update request) {
        UserEntity user = findDetail(id);
        requireVersion(user, request.expectedVersion());
        String email = normalizeEmail(request.email());
        if (userRepository.existsByEmailIgnoreCaseAndIdNot(email, id)) {
            throw AppException.conflict("USER_EMAIL_EXISTS", "邮箱已存在");
        }
        user.updateProfile(email, request.displayName().trim());
        userRepository.flush();
        return UserView.from(user, user.getRoles().stream().map(RoleEntity::getCode).toList());
    }

    @Transactional
    public UserView changeStatus(UUID id, UserRequests.ChangeStatus request) {
        UserEntity user = findDetail(id);
        requireVersion(user, request.expectedVersion());
        user.changeStatus(request.status());
        userRepository.flush();
        return UserView.from(user, user.getRoles().stream().map(RoleEntity::getCode).toList());
    }

    @Transactional
    public UserView replaceRoles(UUID id, UserRequests.ReplaceRoles request) {
        UserEntity user = findDetail(id);
        requireVersion(user, request.expectedVersion());
        Set<RoleEntity> roles = loadActiveRoles(request.roleIds());
        user.replaceRoles(roles);
        userRepository.flush();
        return UserView.from(user, roles.stream().map(RoleEntity::getCode).toList());
    }

    private UserEntity findDetail(UUID id) {
        return userRepository.findDetailById(id)
                .orElseThrow(() -> AppException.notFound("USER_NOT_FOUND", "用户不存在"));
    }

    private Set<RoleEntity> loadActiveRoles(Set<UUID> roleIds) {
        List<RoleEntity> roles = roleRepository.findAllByIdInAndStatus(roleIds, "ACTIVE");
        if (roles.size() != roleIds.size()) {
            throw AppException.badRequest("INVALID_ROLE", "包含不存在或已停用的角色");
        }
        return new LinkedHashSet<>(roles);
    }

    private Map<UUID, List<String>> loadRoleCodes(Collection<UserEntity> users) {
        if (users.isEmpty()) {
            return Map.of();
        }
        List<UUID> ids = users.stream().map(UserEntity::getId).toList();
        Map<UUID, List<String>> result = new LinkedHashMap<>();
        for (UserRoleCodeRow row : userRepository.findRoleCodes(ids)) {
            result.computeIfAbsent(row.getUserId(), ignored -> new ArrayList<>()).add(row.getRoleCode());
        }
        return result;
    }

    private void requireVersion(UserEntity user, long expectedVersion) {
        if (user.getVersion() != expectedVersion) {
            throw AppException.conflict("STALE_VERSION", "数据已被其他请求修改，请刷新后重试");
        }
    }

    private String normalizeEmail(String email) {
        return email.trim().toLowerCase(Locale.ROOT);
    }
}
