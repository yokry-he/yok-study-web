package com.example.admin.role;

import java.util.List;

import jakarta.servlet.http.HttpServletRequest;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import com.example.admin.common.api.ApiResponse;
import com.example.admin.common.web.RequestIdFilter;

@RestController
@RequestMapping("/api/roles")
public class RoleController {
    private final RoleService roleService;

    public RoleController(RoleService roleService) {
        this.roleService = roleService;
    }

    @GetMapping
    ApiResponse<List<RoleView>> list(HttpServletRequest request) {
        return ApiResponse.success(roleService.listActiveRoles(), RequestIdFilter.get(request));
    }
}
