# Go Task API

这是文档站 Go 模块配套的完整示例项目：使用 Go 标准库 `net/http`、`database/sql`、pgx、PostgreSQL 和 Docker Compose 实现用户与任务 API。代码重点展示分层、严格 JSON、统一错误、状态机、乐观锁、数据库迁移、健康检查、优雅关闭和真实集成测试。

> 本项目没有认证、授权、限流、审计和生产密钥管理，只能用于学习与本地实验，不能直接暴露到公网或部署到生产环境。

## 1. 技术与版本

- Go `1.26.5`：容器构建和最终验证版本。
- `net/http`：Go 1.22+ method pattern 路由。
- `database/sql` + pgx v5：PostgreSQL 驱动与连接池。
- PostgreSQL `18.4`：开发和集成测试数据库。
- golang-migrate：嵌入式 SQL 迁移。
- Testcontainers for Go：测试时启动真实 PostgreSQL。
- Docker Compose：本地一键启动、迁移与健康检查。

## 2. 项目结构

```text
.
├── cmd/
│   ├── api/                 # HTTP 服务进程
│   ├── healthcheck/         # 容器健康检查命令
│   └── migrate/             # up、down 1、version
├── internal/
│   ├── app/                 # 路由装配、Server 和生命周期
│   ├── config/              # 强类型环境变量配置
│   ├── platform/
│   │   ├── database/        # 数据库连接和迁移器
│   │   └── httpx/           # JSON、错误与中间件
│   ├── task/                # 任务领域、Service、Repository、Handler
│   └── user/                # 用户领域、Service、Repository、Handler
├── migrations/              # 嵌入二进制的 PostgreSQL SQL
├── tests/                   # PostgreSQL 与 API 集成测试
├── API_CONTRACT.md          # 12 条路由的完整契约
├── TROUBLESHOOTING.md       # 可执行排障手册
├── Dockerfile
└── compose.yaml
```

依赖方向固定为：

```text
cmd -> app -> handler -> service -> repository -> database/sql -> PostgreSQL
```

Handler 只处理 HTTP 协议，Service 负责校验和状态规则，Repository 负责 SQL 与错误映射。`internal` 包保证示例外部不能绕过装配层直接依赖内部实现。

## 3. 前置条件

本地运行需要：

```bash
go version
docker version
docker compose version
```

推荐使用 Go `1.26.5`。普通单元测试不需要 Docker；带 `integration` tag 的测试会通过 Testcontainers 拉起 `postgres:18.4`，因此必须先启动 Docker Desktop 或兼容 Docker daemon。

## 4. 本地 Go 启动

### 4.1 准备环境变量

```bash
cp .env.example .env
set -a
source .env
set +a
```

`.env` 已被 Git 忽略。默认 `DATABASE_URL` 连接宿主机的 `127.0.0.1:5432`。

### 4.2 只启动 PostgreSQL

```bash
docker compose up -d postgres
docker compose ps postgres
docker compose exec postgres pg_isready -U app -d taskdb
```

### 4.3 执行迁移

```bash
go run ./cmd/migrate up
go run ./cmd/migrate version
```

预期版本输出：

```text
version=1 dirty=false
```

迁移命令只接受三个显式操作：

```bash
go run ./cmd/migrate up
go run ./cmd/migrate down 1
go run ./cmd/migrate version
```

`down 1` 会删除示例表，只应在本地明确需要回滚时使用。

### 4.4 启动 API

```bash
go run ./cmd/api
```

另开终端验证：

```bash
curl --fail-with-body http://127.0.0.1:8080/health/live
curl --fail-with-body http://127.0.0.1:8080/health/ready
```

按 `Ctrl+C` 会触发优雅关闭：先关闭 readiness，再等待在途请求，最后关闭数据库连接池。

## 5. Docker Compose 一键启动

### 5.1 构建并启动

```bash
docker compose up --build -d
```

启动顺序由 Compose 明确约束：

1. PostgreSQL 通过 `pg_isready`。
2. `migrate` 容器成功执行所有迁移并退出 0。
3. API 启动，并由镜像内 `/app/healthcheck` 检查 readiness。

查看状态和日志：

```bash
docker compose ps
docker compose logs migrate
docker compose logs -f api
```

预期 `postgres` 与 `api` 为 healthy，`migrate` 为成功退出。

### 5.2 单独执行迁移命令

```bash
docker compose run --rm migrate up
docker compose run --rm migrate version
docker compose run --rm migrate down 1
```

### 5.3 停止与清理

仅停止并保留数据库数据：

```bash
docker compose stop
```

删除容器和网络，但保留命名卷：

```bash
docker compose down
```

删除容器、网络和数据库卷：

```bash
docker compose down -v
```

`down -v` 会永久删除该 Compose 项目的本地数据库数据，执行前先确认没有需要保留的内容。

若宿主机 5432 已被其他项目占用，可以只覆盖宿主机端口：

```bash
POSTGRES_PORT=55432 docker compose up --build -d
```

API 容器仍通过 `postgres:5432` 访问数据库，不受这个宿主机端口覆盖影响。使用本地 Go 进程连接该数据库时，需要把 `.env` 中的端口同步改为 `55432`。

## 6. 配置说明

| 变量 | 默认值 | 作用 |
|---|---|---|
| `APP_ENV` | `development` | 环境名称 |
| `HTTP_ADDR` | `:8080` | HTTP 监听地址 |
| `HTTP_READ_HEADER_TIMEOUT` | `5s` | 读取 Header 超时 |
| `HTTP_READ_TIMEOUT` | `10s` | 读取完整请求超时 |
| `HTTP_WRITE_TIMEOUT` | `15s` | 写响应超时 |
| `HTTP_IDLE_TIMEOUT` | `60s` | Keep-Alive 空闲超时 |
| `HTTP_REQUEST_TIMEOUT` | `10s` | 单个 Handler context 截止时间 |
| `HTTP_SHUTDOWN_TIMEOUT` | `15s` | 优雅关闭预算 |
| `DATABASE_URL` | 见 `.env.example` | PostgreSQL 连接串，必填 |
| `DB_MAX_OPEN_CONNS` | `20` | 最大打开连接数 |
| `DB_MAX_IDLE_CONNS` | `10` | 最大空闲连接数 |
| `DB_CONN_MAX_LIFETIME` | `30m` | 单连接最长寿命 |
| `DB_CONN_MAX_IDLE_TIME` | `5m` | 单连接最长空闲时间 |
| `LOG_LEVEL` | `info` | `debug`、`info`、`warn`、`error` |

配置解析失败会在进程监听端口前返回错误。错误信息不会回显数据库连接串。

## 7. 测试与质量检查

### 7.1 普通测试

```bash
go test ./...
```

### 7.2 竞态检测

```bash
go test -race ./...
```

### 7.3 静态检查与格式

```bash
test -z "$(gofmt -l .)"
go vet ./...
```

### 7.4 真实 PostgreSQL 集成测试

```bash
GOPROXY=https://proxy.golang.org,direct \
  go test -tags=integration ./... -count=1 -v
```

测试覆盖迁移结构、中文数据库注释、约束、索引、分页快照、乐观锁并发、Repository、完整 API 生命周期和 HTTP 协议边界。测试不能因为 Docker 不可用而静默跳过。

### 7.5 集成竞态检测

```bash
GOPROXY=https://proxy.golang.org,direct \
  go test -race -tags=integration ./... -count=1
```

### 7.6 短时 Fuzz

```bash
go test ./internal/platform/httpx \
  -run '^$' \
  -fuzz '^FuzzDecodeJSON$' \
  -fuzztime 10s
```

Fuzz corpus 可能生成在 Go 缓存目录；发现问题后应把最小复现输入加入普通回归测试。

### 7.7 本机构建

```bash
mkdir -p bin
CGO_ENABLED=0 go build -trimpath -o bin/api ./cmd/api
CGO_ENABLED=0 go build -trimpath -o bin/migrate ./cmd/migrate
CGO_ENABLED=0 go build -trimpath -o bin/healthcheck ./cmd/healthcheck
```

### 7.8 固定 Go 镜像验证

```bash
docker run --rm \
  -v "$PWD:/src" \
  -w /src \
  golang:1.26.5-bookworm \
  sh -lc 'go test ./... && go test -race ./...'
```

## 8. 12 条接口快速验证

完整字段、错误码和并发语义见 [API_CONTRACT.md](./API_CONTRACT.md)。先设置：

```bash
export BASE_URL=http://127.0.0.1:8080
```

### 8.1 Liveness

```bash
curl --fail-with-body "$BASE_URL/health/live"
```

### 8.2 Readiness

```bash
curl --fail-with-body "$BASE_URL/health/ready"
```

### 8.3 用户列表

```bash
curl --fail-with-body "$BASE_URL/api/users?page=1&pageSize=20&status=ACTIVE"
```

### 8.4 创建用户

```bash
curl --fail-with-body -i \
  -X POST "$BASE_URL/api/users" \
  -H 'Content-Type: application/json' \
  -d '{"name":"张三","email":"user@example.com"}'
```

### 8.5 用户详情

```bash
curl --fail-with-body "$BASE_URL/api/users/1"
```

### 8.6 修改用户状态

```bash
curl --fail-with-body \
  -X PATCH "$BASE_URL/api/users/1/status" \
  -H 'Content-Type: application/json' \
  -d '{"status":"DISABLED","expectedVersion":0}'
```

创建任务前需要负责人保持 `ACTIVE`；如果刚执行了上一步，请创建另一个用户或把状态改回 `ACTIVE` 并使用最新版本。

### 8.7 任务列表

```bash
curl --fail-with-body "$BASE_URL/api/tasks?page=1&pageSize=20&ownerId=1&status=TODO"
```

### 8.8 创建任务

```bash
curl --fail-with-body -i \
  -X POST "$BASE_URL/api/tasks" \
  -H 'Content-Type: application/json' \
  -d '{"ownerId":1,"title":"完成 Go 文档","description":"补齐示例和排障说明"}'
```

### 8.9 任务详情

```bash
curl --fail-with-body "$BASE_URL/api/tasks/1"
```

### 8.10 更新任务

```bash
curl --fail-with-body \
  -X PUT "$BASE_URL/api/tasks/1" \
  -H 'Content-Type: application/json' \
  -d '{"title":"完成 Go 全量文档","description":null,"dueAt":null,"expectedVersion":0}'
```

### 8.11 修改任务状态

```bash
curl --fail-with-body \
  -X PATCH "$BASE_URL/api/tasks/1/status" \
  -H 'Content-Type: application/json' \
  -d '{"status":"DOING","expectedVersion":1}'
```

### 8.12 删除任务

```bash
curl --fail-with-body -i \
  -X DELETE "$BASE_URL/api/tasks/1?expectedVersion=2"
```

返回 `204` 时响应体为空。

## 9. 常见协议要点

- `POST`、`PUT`、`PATCH` 必须使用 `application/json`。
- 请求 JSON 不能包含未知字段，也不能连续发送两个 JSON 值。
- 请求体最大 1 MiB。
- 列表使用稳定排序，翻页时相同创建时间的数据仍由 ID 决定顺序。
- 所有修改和删除使用 `expectedVersion`，冲突返回 409。
- `DONE` 和 `CANCELLED` 是终态，不能再转换。
- `DELETE` 成功返回 204，不返回 JSON envelope。
- 500 响应不会暴露 SQL、密码、panic 值或内部堆栈。

## 10. 运行与安全边界

这个示例刻意不包含：

- 登录、JWT、Session 或 RBAC。
- HTTPS 终止和可信代理配置。
- 限流、防滥用、审计和数据保留策略。
- 密钥管理、数据库备份和灾难恢复。
- 跨实例分布式追踪与指标系统。
- 生产发布、滚动升级和零停机迁移策略。

因此它适合学习“一个可靠 API 的核心骨架”，不代表完整的生产安全基线。

遇到启动、迁移、连接池、超时、容器健康或关闭问题时，按 [TROUBLESHOOTING.md](./TROUBLESHOOTING.md) 的证据步骤处理。
