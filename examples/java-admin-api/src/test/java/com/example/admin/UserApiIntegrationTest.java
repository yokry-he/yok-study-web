package com.example.admin;

import com.jayway.jsonpath.JsonPath;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.boot.webmvc.test.autoconfigure.AutoConfigureMockMvc;
import org.springframework.context.annotation.Import;
import org.springframework.http.MediaType;
import org.springframework.test.web.servlet.MockMvc;

import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.get;
import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.patch;
import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.post;
import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.put;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.header;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.jsonPath;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.status;

@Import(TestcontainersConfiguration.class)
@SpringBootTest
@AutoConfigureMockMvc
class UserApiIntegrationTest {
    private static final String ADMIN_ROLE_ID = "00000000-0000-0000-0000-000000000001";

    @Autowired
    private MockMvc mvc;

    @Test
    void completesUserLifecycleAndRejectsStaleUpdates() throws Exception {
        mvc.perform(get("/api/roles"))
                .andExpect(status().isOk())
                .andExpect(header().exists("X-Request-Id"))
                .andExpect(jsonPath("$.data.length()").value(2));

        String createBody = """
                {
                  "email": "Ada@Example.com",
                  "displayName": "Ada",
                  "roleIds": ["%s"]
                }
                """.formatted(ADMIN_ROLE_ID);

        var created = mvc.perform(post("/api/users")
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(createBody))
                .andExpect(status().isCreated())
                .andExpect(jsonPath("$.success").value(true))
                .andExpect(jsonPath("$.data.email").value("ada@example.com"))
                .andExpect(jsonPath("$.data.roleCodes[0]").value("ADMIN"))
                .andExpect(jsonPath("$.data.createdAt").isNotEmpty())
                .andReturn();

        String userId = created.getResponse().getHeader("Location").replace("/api/users/", "");
        int currentVersion = JsonPath.read(created.getResponse().getContentAsString(), "$.data.version");

        mvc.perform(get("/api/users").param("q", "ada"))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.data.total").value(1))
                .andExpect(jsonPath("$.data.items[0].roleCodes[0]").value("ADMIN"));

        mvc.perform(patch("/api/users/{id}/status", userId)
                        .contentType(MediaType.APPLICATION_JSON)
                        .content("{\"status\":\"DISABLED\",\"expectedVersion\":%d}"
                                .formatted(currentVersion)))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.data.status").value("DISABLED"))
                .andExpect(jsonPath("$.data.version").value(currentVersion + 1));

        mvc.perform(patch("/api/users/{id}/status", userId)
                        .contentType(MediaType.APPLICATION_JSON)
                        .content("{\"status\":\"ACTIVE\",\"expectedVersion\":%d}"
                                .formatted(currentVersion)))
                .andExpect(status().isConflict())
                .andExpect(jsonPath("$.error.code").value("STALE_VERSION"));

        mvc.perform(post("/api/users")
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(createBody))
                .andExpect(status().isConflict())
                .andExpect(jsonPath("$.error.code").value("USER_EMAIL_EXISTS"));
    }

    @Test
    void returnsFieldErrorsForInvalidInput() throws Exception {
        mvc.perform(post("/api/users")
                        .contentType(MediaType.APPLICATION_JSON)
                        .content("{\"email\":\"bad\",\"displayName\":\"\",\"roleIds\":[]}"))
                .andExpect(status().isBadRequest())
                .andExpect(jsonPath("$.error.code").value("VALIDATION_ERROR"))
                .andExpect(jsonPath("$.error.fields.email").exists())
                .andExpect(jsonPath("$.error.fields.displayName").exists())
                .andExpect(jsonPath("$.error.fields.roleIds").exists());
    }

    @Test
    void returnsTheSameErrorContractForUnknownRoutes() throws Exception {
        mvc.perform(get("/api/not-found"))
                .andExpect(status().isNotFound())
                .andExpect(header().exists("X-Request-Id"))
                .andExpect(jsonPath("$.success").value(false))
                .andExpect(jsonPath("$.error.code").value("ROUTE_NOT_FOUND"))
                .andExpect(jsonPath("$.requestId").isNotEmpty());
    }

    @Test
    void rejectsUpdatesWithoutAnExplicitVersion() throws Exception {
        var created = mvc.perform(post("/api/users")
                        .contentType(MediaType.APPLICATION_JSON)
                        .content("""
                                {
                                  "email": "missing-version@example.com",
                                  "displayName": "Missing Version",
                                  "roleIds": ["%s"]
                                }
                                """.formatted(ADMIN_ROLE_ID)))
                .andExpect(status().isCreated())
                .andReturn();
        String userId = created.getResponse().getHeader("Location").replace("/api/users/", "");

        mvc.perform(patch("/api/users/{id}/status", userId)
                        .contentType(MediaType.APPLICATION_JSON)
                        .content("{\"status\":\"DISABLED\"}"))
                .andExpect(status().isBadRequest())
                .andExpect(jsonPath("$.error.code").value("VALIDATION_ERROR"))
                .andExpect(jsonPath("$.error.fields.expectedVersion").exists());

        mvc.perform(put("/api/users/{id}", userId)
                        .contentType(MediaType.APPLICATION_JSON)
                        .content("""
                                {
                                  "email": "missing-version@example.com",
                                  "displayName": "Updated Name"
                                }
                                """))
                .andExpect(status().isBadRequest())
                .andExpect(jsonPath("$.error.code").value("VALIDATION_ERROR"))
                .andExpect(jsonPath("$.error.fields.expectedVersion").exists());

        mvc.perform(put("/api/users/{id}/roles", userId)
                        .contentType(MediaType.APPLICATION_JSON)
                        .content("{\"roleIds\":[\"%s\"]}".formatted(ADMIN_ROLE_ID)))
                .andExpect(status().isBadRequest())
                .andExpect(jsonPath("$.error.code").value("VALIDATION_ERROR"))
                .andExpect(jsonPath("$.error.fields.expectedVersion").exists());
    }

    @Test
    void mapsCommonHttpInputFailuresToStableClientErrors() throws Exception {
        mvc.perform(get("/api/users/not-a-uuid"))
                .andExpect(status().isBadRequest())
                .andExpect(jsonPath("$.error.code").value("INVALID_ARGUMENT"));

        mvc.perform(get("/api/users").param("page", "not-a-number"))
                .andExpect(status().isBadRequest())
                .andExpect(jsonPath("$.error.code").value("INVALID_ARGUMENT"));

        mvc.perform(put("/api/roles"))
                .andExpect(status().isMethodNotAllowed())
                .andExpect(jsonPath("$.error.code").value("METHOD_NOT_ALLOWED"));

        mvc.perform(post("/api/users")
                        .contentType(MediaType.TEXT_PLAIN)
                        .content("not-json"))
                .andExpect(status().isUnsupportedMediaType())
                .andExpect(jsonPath("$.error.code").value("UNSUPPORTED_MEDIA_TYPE"));
    }
}
