# Java Admin API

本目录是 Java 文档模块的可运行用户角色 API，正文见 [Spring Boot 从零到项目落地](../../docs/java/spring-boot-project-from-zero.md)。

## 基线

- Java 25
- Spring Boot 4.1.0
- Maven 3.9.11
- PostgreSQL 18
- Testcontainers 2

## 测试

```bash
mvn -B -ntp test
```

测试需要可用的 Docker daemon，以便 Testcontainers 启动 PostgreSQL。

## 启动

```bash
docker compose up --build
```

启动后访问：

- API：`http://127.0.0.1:8080/api/roles`
- Liveness：`http://127.0.0.1:8080/actuator/health/liveness`
- Readiness：`http://127.0.0.1:8080/actuator/health/readiness`

本示例中的数据库密码只用于本地 Compose。生产环境必须由部署平台注入密钥。
