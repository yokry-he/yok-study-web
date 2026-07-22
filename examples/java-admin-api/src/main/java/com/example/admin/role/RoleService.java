package com.example.admin.role;

import java.util.List;

import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

@Service
public class RoleService {
    private final RoleRepository roleRepository;

    public RoleService(RoleRepository roleRepository) {
        this.roleRepository = roleRepository;
    }

    @Transactional(readOnly = true)
    public List<RoleView> listActiveRoles() {
        return roleRepository.findAllByStatusOrderByNameAsc("ACTIVE")
                .stream()
                .map(RoleView::from)
                .toList();
    }
}
