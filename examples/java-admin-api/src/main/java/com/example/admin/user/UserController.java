package com.example.admin.user;

import java.net.URI;
import java.util.UUID;

import jakarta.servlet.http.HttpServletRequest;
import jakarta.validation.Valid;
import jakarta.validation.constraints.Max;
import jakarta.validation.constraints.Min;
import jakarta.validation.constraints.Size;
import org.springframework.http.ResponseEntity;
import org.springframework.validation.annotation.Validated;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PatchMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.PutMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;

import com.example.admin.common.api.ApiResponse;
import com.example.admin.common.api.PageData;
import com.example.admin.common.web.RequestIdFilter;

@Validated
@RestController
@RequestMapping("/api/users")
public class UserController {
    private final UserService userService;

    public UserController(UserService userService) {
        this.userService = userService;
    }

    @GetMapping
    ApiResponse<PageData<UserView>> list(
            @RequestParam(defaultValue = "") @Size(max = 100) String q,
            @RequestParam(defaultValue = "0") @Min(0) int page,
            @RequestParam(defaultValue = "20") @Min(1) @Max(100) int pageSize,
            HttpServletRequest request
    ) {
        return ApiResponse.success(userService.list(q, page, pageSize), RequestIdFilter.get(request));
    }

    @GetMapping("/{id}")
    ApiResponse<UserView> get(@PathVariable UUID id, HttpServletRequest request) {
        return ApiResponse.success(userService.get(id), RequestIdFilter.get(request));
    }

    @PostMapping
    ResponseEntity<ApiResponse<UserView>> create(
            @Valid @RequestBody UserRequests.Create body,
            HttpServletRequest request
    ) {
        UserView user = userService.create(body);
        return ResponseEntity.created(URI.create("/api/users/" + user.id()))
                .body(ApiResponse.success(user, RequestIdFilter.get(request)));
    }

    @PutMapping("/{id}")
    ApiResponse<UserView> update(
            @PathVariable UUID id,
            @Valid @RequestBody UserRequests.Update body,
            HttpServletRequest request
    ) {
        return ApiResponse.success(userService.update(id, body), RequestIdFilter.get(request));
    }

    @PatchMapping("/{id}/status")
    ApiResponse<UserView> changeStatus(
            @PathVariable UUID id,
            @Valid @RequestBody UserRequests.ChangeStatus body,
            HttpServletRequest request
    ) {
        return ApiResponse.success(userService.changeStatus(id, body), RequestIdFilter.get(request));
    }

    @PutMapping("/{id}/roles")
    ApiResponse<UserView> replaceRoles(
            @PathVariable UUID id,
            @Valid @RequestBody UserRequests.ReplaceRoles body,
            HttpServletRequest request
    ) {
        return ApiResponse.success(userService.replaceRoles(id, body), RequestIdFilter.get(request));
    }
}
