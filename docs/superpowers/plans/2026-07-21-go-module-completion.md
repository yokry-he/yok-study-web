# Go 模块全量完善 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 交付一个标准库优先、可构建可测试的 Go 任务 API，并把现有 Go 文档完善为包含 26 张核心图、16 个真实问题和 12 个练习的完整学习闭环。

**Architecture:** 示例使用 `net/http`、`database/sql`、pgx 和 PostgreSQL 18，按 `cmd -> app -> user/task -> repository` 组织依赖。文档通过 VitePress `<<<` 导入经过测试的示例源码；准确流程使用 Mermaid，两个抽象心智模型使用 imagegen 教学图片，所有视觉资产进入统一登记。

**Tech Stack:** Go 1.26.5、`net/http`、`log/slog`、`database/sql`、pgx v5.10.0、golang-migrate v4.19.1、Testcontainers for Go v0.43.0、PostgreSQL 18.4、Docker Compose、VitePress 1.6.4、Mermaid 11.16。

---

## 执行约束

- 仓库路径：`/Users/yokry/Documents/Codex/2026-07-01/yok-study-web`。
- 当前 Java 批次改动未提交，执行时不得还原、覆盖或混入无关格式化。
- 本计划不自动提交、推送或部署；每个任务用测试和 `git diff --check` 作为检查点。
- 先写失败测试，再写最小实现；集成测试用 `//go:build integration` 与普通测试隔离。
- Go 依赖解析显式使用 `GOPROXY=https://proxy.golang.org,direct`，避免当前机器配置的镜像返回过旧 pgx 版本。本机 Go 1.26.0 用于快速测试，最终必须在 `golang:1.26.5-bookworm` 中再跑一次完整测试。
- 每次新增 Mermaid 后都要浏览器确认 SVG，不把构建通过当作渲染通过。

## 文件结构

### 示例工程

```text
examples/go-task-api/
├─ .dockerignore
├─ .env.example
├─ API_CONTRACT.md
├─ Dockerfile
├─ README.md
├─ TROUBLESHOOTING.md
├─ compose.yaml
├─ go.mod
├─ go.sum
├─ cmd/
│  ├─ api/main.go
│  ├─ healthcheck/main.go
│  └─ migrate/main.go
├─ internal/
│  ├─ app/
│  │  ├─ app.go
│  │  ├─ router.go
│  │  ├─ router_test.go
│  │  └─ server.go
│  ├─ config/
│  │  ├─ config.go
│  │  └─ config_test.go
│  ├─ platform/
│  │  ├─ database/
│  │  │  ├─ database.go
│  │  │  └─ migrate.go
│  │  └─ httpx/
│  │     ├─ error.go
│  │     ├─ json.go
│  │     ├─ json_test.go
│  │     ├─ middleware.go
│  │     ├─ middleware_test.go
│  │     └─ response.go
│  ├─ task/
│  │  ├─ handler.go
│  │  ├─ handler_test.go
│  │  ├─ model.go
│  │  ├─ repository.go
│  │  ├─ repository_postgres.go
│  │  ├─ service.go
│  │  └─ service_test.go
│  └─ user/
│     ├─ handler.go
│     ├─ handler_test.go
│     ├─ model.go
│     ├─ repository.go
│     ├─ repository_postgres.go
│     ├─ service.go
│     └─ service_test.go
├─ migrations/
│  ├─ 000001_create_users_tasks.down.sql
│  ├─ 000001_create_users_tasks.up.sql
│  └─ embed.go
└─ tests/
   ├─ api_integration_test.go
   └─ postgres_test.go
```

### 文档和视觉资产

```text
docs/go/*.md
docs/projects/issues-go.md
docs/roadmap/go-practice.md
docs/public/images/go/go-api-request-journey.webp
docs/public/images/go/go-concurrency-workshop.webp
docs/contribute/visual-asset-register.md
docs/.vitepress/config.ts
docs/contribute/module-status.md
docs/technologies/index.md
docs/technologies/expansion-plan.md
docs/projects/real-world-issues.md
docs/projects/issues-backend.md
docs/roadmap/introduction.md
docs/roadmap/practice-labs.md
docs/roadmap/reading-guide.md
README.md
.gitignore
```

## Task 1: 初始化 Go 模块和配置契约

**Files:**
- Create: `examples/go-task-api/go.mod`
- Create: `examples/go-task-api/.env.example`
- Create: `examples/go-task-api/internal/config/config_test.go`
- Create: `examples/go-task-api/internal/config/config.go`

- [ ] **Step 1: 创建固定版本的模块文件**

```go
module github.com/yokry-he/yok-study-web/examples/go-task-api

go 1.26.0

require (
	github.com/golang-migrate/migrate/v4 v4.19.1
	github.com/jackc/pgx/v5 v5.10.0
	github.com/testcontainers/testcontainers-go v0.43.0
	github.com/testcontainers/testcontainers-go/modules/postgres v0.43.0
)
```

`.env.example` 固定包含：

```dotenv
APP_ENV=development
HTTP_ADDR=:8080
HTTP_READ_HEADER_TIMEOUT=5s
HTTP_READ_TIMEOUT=10s
HTTP_WRITE_TIMEOUT=15s
HTTP_IDLE_TIMEOUT=60s
HTTP_REQUEST_TIMEOUT=10s
HTTP_SHUTDOWN_TIMEOUT=15s
DATABASE_URL=postgres://app:app@127.0.0.1:5432/taskdb?sslmode=disable
DB_MAX_OPEN_CONNS=20
DB_MAX_IDLE_CONNS=10
DB_CONN_MAX_LIFETIME=30m
DB_CONN_MAX_IDLE_TIME=5m
LOG_LEVEL=info
```

- [ ] **Step 2: 先写配置失败测试**

```go
func TestLoadRejectsMissingDatabaseURL(t *testing.T) {
	t.Setenv("DATABASE_URL", "")
	_, err := Load()
	if !errors.Is(err, ErrDatabaseURLRequired) {
		t.Fatalf("expected ErrDatabaseURLRequired, got %v", err)
	}
}

func TestLoadParsesDefaults(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://app:secret@localhost:5432/taskdb?sslmode=disable")
	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.HTTP.Addr != ":8080" || cfg.Database.MaxOpenConns != 20 {
		t.Fatalf("unexpected defaults: %+v", cfg)
	}
}
```

- [ ] **Step 3: 运行测试确认因 `Load` 不存在而失败**

Run: `cd examples/go-task-api && go test ./internal/config`

Expected: FAIL，错误包含 `undefined: Load`。

- [ ] **Step 4: 实现强类型配置和脱敏错误**

```go
var ErrDatabaseURLRequired = errors.New("DATABASE_URL is required")

type Config struct {
	Environment string
	LogLevel    slog.Level
	HTTP        HTTPConfig
	Database    DatabaseConfig
}

type HTTPConfig struct {
	Addr              string
	ReadHeaderTimeout time.Duration
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
	RequestTimeout    time.Duration
	ShutdownTimeout   time.Duration
}

type DatabaseConfig struct {
	URL             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}
```

`Load` 使用 `os.LookupEnv`、`strconv.Atoi` 和 `time.ParseDuration`，任何无效字段返回 `invalid DB_MAX_OPEN_CONNS: must be a positive integer` 这类不含连接串的错误；`LOG_LEVEL` 只接受 `debug`、`info`、`warn`、`error`。

- [ ] **Step 5: 覆盖无效数字、持续时间和日志级别并通过测试**

Run: `cd examples/go-task-api && go test ./internal/config -v`

Expected: PASS，至少覆盖默认值、缺失 DSN、非法数字、非法持续时间和非法日志级别。

## Task 2: 建立统一 JSON、ID 和错误协议

**Files:**
- Create: `examples/go-task-api/internal/platform/httpx/error.go`
- Create: `examples/go-task-api/internal/platform/httpx/response.go`
- Create: `examples/go-task-api/internal/platform/httpx/json.go`
- Create: `examples/go-task-api/internal/platform/httpx/json_test.go`

- [ ] **Step 1: 先写严格 JSON 解码测试和 Fuzz target**

```go
type decodeFixture struct {
	Name string `json:"name"`
}

func TestDecodeJSONRejectsUnknownField(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"name":"demo","extra":true}`))
	rec := httptest.NewRecorder()
	var dst decodeFixture
	err := DecodeJSON(rec, req, &dst, 1024)
	if !errors.Is(err, ErrUnknownJSONField) {
		t.Fatalf("expected unknown field error, got %v", err)
	}
}

func FuzzDecodeJSON(f *testing.F) {
	f.Add(`{"name":"demo"}`)
	f.Fuzz(func(t *testing.T, body string) {
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
		rec := httptest.NewRecorder()
		var dst decodeFixture
		_ = DecodeJSON(rec, req, &dst, 1024)
	})
}
```

- [ ] **Step 2: 运行测试确认失败**

Run: `cd examples/go-task-api && go test ./internal/platform/httpx`

Expected: FAIL，错误包含 `undefined: DecodeJSON`。

- [ ] **Step 3: 实现响应类型和错误码**

```go
type FieldErrors map[string]string

type APIError struct {
	Status  int
	Code    string
	Message string
	Fields  FieldErrors
	Cause   error
}

type ErrorBody struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Fields  FieldErrors `json:"fields,omitempty"`
}

type Envelope struct {
	Success   bool       `json:"success"`
	Data      any        `json:"data,omitempty"`
	Error     *ErrorBody `json:"error,omitempty"`
	RequestID string     `json:"requestId"`
}
```

定义 `INVALID_JSON`、`UNKNOWN_FIELD`、`BODY_TOO_LARGE`、`INVALID_ARGUMENT`、`VALIDATION_FAILED`、`NOT_FOUND`、`CONFLICT`、`METHOD_NOT_ALLOWED`、`UNSUPPORTED_MEDIA_TYPE`、`DEADLINE_EXCEEDED`、`INTERNAL_ERROR`。

- [ ] **Step 4: 实现严格解码、单值检查和整数 ID 解析**

```go
func DecodeJSON(w http.ResponseWriter, r *http.Request, dst any, maxBytes int64) error {
	r.Body = http.MaxBytesReader(w, r.Body, maxBytes)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		return classifyDecodeError(err)
	}
	if err := dec.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return ErrMultipleJSONValues
	}
	return nil
}

func ParsePositiveID(raw string) (int64, error) {
	id, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || id <= 0 {
		return 0, NewAPIError(http.StatusBadRequest, "INVALID_ARGUMENT", "资源 ID 必须是正整数", nil)
	}
	return id, nil
}
```

- [ ] **Step 5: 实现 `WriteData`、`WriteError` 和 context 错误映射**

`WriteData` 必须设置 `Content-Type: application/json; charset=utf-8`；`WriteError` 只向客户端返回公开消息，`Cause` 只交给日志层。`context.DeadlineExceeded` 映射 504，`context.Canceled` 在未提交响应时不写业务错误。

- [ ] **Step 6: 运行单元测试和短时 Fuzz**

Run:

```bash
cd examples/go-task-api
go test ./internal/platform/httpx -v
go test -run=^$ -fuzz=FuzzDecodeJSON -fuzztime=10s ./internal/platform/httpx
```

Expected: 两条命令 PASS，无 panic。

## Task 3: 建立 request id、日志、恢复和超时中间件

**Files:**
- Create: `examples/go-task-api/internal/platform/httpx/middleware.go`
- Create: `examples/go-task-api/internal/platform/httpx/middleware_test.go`

- [ ] **Step 1: 写中间件顺序和 panic 回归测试**

```go
func TestRecoverHidesPanicAndKeepsRequestID(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))
	h := RequestID(AccessLog(logger)(Recover(logger)(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		panic("database password=secret")
	}))))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/panic", nil))
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("got %d", rec.Code)
	}
	if strings.Contains(rec.Body.String(), "secret") {
		t.Fatal("panic detail leaked")
	}
	if rec.Header().Get("X-Request-ID") == "" {
		t.Fatal("missing request id")
	}
}
```

- [ ] **Step 2: 运行测试确认中间件尚不存在**

Run: `cd examples/go-task-api && go test ./internal/platform/httpx -run Middleware -v`

Expected: FAIL，出现未定义符号。

- [ ] **Step 3: 实现 request id 上下文和状态记录器**

```go
type requestIDKey struct{}

func RequestIDFromContext(ctx context.Context) string {
	value, _ := ctx.Value(requestIDKey{}).(string)
	return value
}

type responseRecorder struct {
	http.ResponseWriter
	status int
	bytes  int
}
```

合法调用方 request id 限制为 1 到 128 个 `[A-Za-z0-9._:-]` 字符；否则用 `crypto/rand` 生成 `req_` 加 32 位十六进制字符串。

- [ ] **Step 4: 实现固定中间件链**

```go
func Chain(final http.Handler, middleware ...func(http.Handler) http.Handler) http.Handler {
	for i := len(middleware) - 1; i >= 0; i-- {
		final = middleware[i](final)
	}
	return final
}
```

路由最终以 `Chain(router, RequestID, AccessLog(logger), Recover(logger), Deadline(cfg.HTTP.RequestTimeout), RequireJSON(1<<20))` 组装。Deadline 只调用 `context.WithTimeout` 并传递给下游，不额外启动 handler goroutine。`RequireJSON` 只检查 POST、PUT、PATCH，并在进入 Handler 前限制 body 大小。

- [ ] **Step 5: 覆盖调用方 request id、非法 request id、panic、状态码、耗时、deadline 和 JSON 媒体类型**

`RequireJSON` 测试固定覆盖：GET 无 Content-Type 仍进入下游、POST `text/plain` 返回 415、POST `application/json; charset=utf-8` 进入下游、超过 1 MiB 的 body 返回 400 `BODY_TOO_LARGE`。

Run: `cd examples/go-task-api && go test ./internal/platform/httpx -v`

Expected: PASS；panic 内容不出现在响应，超时 context 能被下游观察。

## Task 4: 用测试定义用户领域

**Files:**
- Create: `examples/go-task-api/internal/user/model.go`
- Create: `examples/go-task-api/internal/user/repository.go`
- Create: `examples/go-task-api/internal/user/service_test.go`
- Create: `examples/go-task-api/internal/user/service.go`

- [ ] **Step 1: 定义用户模型和仓储契约**

```go
type Status string

const (
	StatusActive   Status = "ACTIVE"
	StatusDisabled Status = "DISABLED"
)

type User struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Status    Status    `json:"status"`
	Version   int64     `json:"version"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Repository interface {
	Create(context.Context, CreateParams) (User, error)
	Get(context.Context, int64) (User, error)
	List(context.Context, ListFilter) (Page, error)
	UpdateStatus(context.Context, int64, Status, int64) (User, error)
}
```

- [ ] **Step 2: 先写 Service 表格测试**

覆盖：姓名去除首尾空格、邮箱转小写、非法邮箱、非法状态、分页上限、仓储冲突、资源不存在和旧版本冲突。

```go
func TestCreateNormalizesEmail(t *testing.T) {
	repo := &fakeRepository{}
	service := NewService(repo)
	_, err := service.Create(context.Background(), CreateInput{Name: "  张三  ", Email: "USER@EXAMPLE.COM"})
	if err != nil {
		t.Fatal(err)
	}
	if repo.created.Name != "张三" || repo.created.Email != "user@example.com" {
		t.Fatalf("unexpected params: %+v", repo.created)
	}
}
```

- [ ] **Step 3: 运行测试确认 Service 不存在**

Run: `cd examples/go-task-api && go test ./internal/user -v`

Expected: FAIL，错误包含 `undefined: NewService`。

- [ ] **Step 4: 实现用户业务规则**

`Create` 默认状态为 `ACTIVE`；名称长度 2 到 64；邮箱通过 `mail.ParseAddress` 后要求解析地址与原始地址一致；分页默认 1/20，上限 100；`ChangeStatusInput.ExpectedVersion` 必须大于等于 0。

- [ ] **Step 5: 运行用户单元测试**

Run: `cd examples/go-task-api && go test ./internal/user -v`

Expected: PASS。

## Task 5: 用测试定义任务状态机

**Files:**
- Create: `examples/go-task-api/internal/task/model.go`
- Create: `examples/go-task-api/internal/task/repository.go`
- Create: `examples/go-task-api/internal/task/service_test.go`
- Create: `examples/go-task-api/internal/task/service.go`

- [ ] **Step 1: 定义任务模型和允许的状态转换**

```go
type Status string

const (
	StatusTodo      Status = "TODO"
	StatusDoing     Status = "DOING"
	StatusDone      Status = "DONE"
	StatusCancelled Status = "CANCELLED"
)

var allowedTransitions = map[Status]map[Status]bool{
	StatusTodo:  {StatusDoing: true, StatusCancelled: true},
	StatusDoing: {StatusTodo: true, StatusDone: true, StatusCancelled: true},
	StatusDone:  {},
	StatusCancelled: {},
}
```

- [ ] **Step 2: 写失败测试覆盖创建、编辑和状态转换**

```go
func TestChangeStatusRejectsDoneToDoing(t *testing.T) {
	repo := &fakeRepository{current: Task{ID: 1, Status: StatusDone, Version: 3}}
	service := NewService(repo)
	_, err := service.ChangeStatus(context.Background(), 1, ChangeStatusInput{Status: StatusDoing, ExpectedVersion: 3})
	if !errors.Is(err, ErrInvalidTransition) {
		t.Fatalf("expected ErrInvalidTransition, got %v", err)
	}
}
```

测试还要覆盖：负责人不存在或停用、标题空白、标题超过 128、截止时间早于当前时间、pageSize 大于 100、更新与删除版本冲突。

- [ ] **Step 3: 运行测试确认失败**

Run: `cd examples/go-task-api && go test ./internal/task -v`

Expected: FAIL，Service 尚未实现。

- [ ] **Step 4: 实现 Task Service**

Service 依赖 `task.Repository` 和最小 `UserReader`：

```go
type UserReader interface {
	Get(context.Context, int64) (user.User, error)
}

type Service struct {
	repo  Repository
	users UserReader
	now   func() time.Time
}
```

通过注入 `now` 让截止时间测试可重复；创建时只接受 ACTIVE 用户；PUT 更新不直接改变状态；状态变化只走 `ChangeStatus`。

- [ ] **Step 5: 运行任务单元测试**

Run: `cd examples/go-task-api && go test ./internal/task -v`

Expected: PASS。

## Task 6: 先用真实 PostgreSQL 测试定义数据库行为

**Files:**
- Create: `examples/go-task-api/migrations/embed.go`
- Create: `examples/go-task-api/migrations/000001_create_users_tasks.up.sql`
- Create: `examples/go-task-api/migrations/000001_create_users_tasks.down.sql`
- Create: `examples/go-task-api/internal/platform/database/database.go`
- Create: `examples/go-task-api/internal/platform/database/migrate.go`
- Create: `examples/go-task-api/tests/postgres_test.go`

- [ ] **Step 1: 编写带中文注释的向上迁移**

```sql
create table users (
  id bigint generated always as identity primary key,
  name varchar(64) not null,
  email varchar(254) not null,
  status varchar(16) not null default 'ACTIVE',
  version bigint not null default 0,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  constraint ck_users_status check (status in ('ACTIVE', 'DISABLED')),
  constraint ck_users_version check (version >= 0)
);

create unique index uk_users_email_lower on users (lower(email));

create table tasks (
  id bigint generated always as identity primary key,
  owner_id bigint not null references users(id) on delete restrict,
  title varchar(128) not null,
  description text,
  status varchar(16) not null default 'TODO',
  due_at timestamptz,
  version bigint not null default 0,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  constraint ck_tasks_status check (status in ('TODO', 'DOING', 'DONE', 'CANCELLED')),
  constraint ck_tasks_version check (version >= 0)
);

create index idx_tasks_owner_status on tasks(owner_id, status, id);
create index idx_tasks_status_due_at on tasks(status, due_at, id);
```

同一文件继续为两张表、全部字段、两个状态约束、两个版本约束、唯一索引和查询索引添加中文 `comment on`。向下迁移按 `tasks`、`users` 顺序删除。

- [ ] **Step 2: 嵌入迁移文件并实现迁移器**

```go
//go:embed *.sql
var Files embed.FS
```

`database.NewMigrator(db)` 使用 `source/iofs` 和 migrate PostgreSQL driver；`Up` 将 `migrate.ErrNoChange` 视为成功，`DownOne` 调用 `Steps(-1)`，`Version` 返回版本号和 dirty 状态。

- [ ] **Step 3: 创建带 build tag 的 Testcontainers 测试基座**

```go
//go:build integration

func startPostgres(t *testing.T) *sql.DB {
	t.Helper()
	ctx := context.Background()
	container, err := postgres.Run(ctx, "postgres:18.4",
		postgres.WithDatabase("taskdb"),
		postgres.WithUsername("app"),
		postgres.WithPassword("app"),
		postgres.BasicWaitStrategies(),
	)
	if err != nil {
		t.Fatal(err)
	}
	connectionString, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		_ = testcontainers.TerminateContainer(container)
		t.Fatal(err)
	}
	db, err := sql.Open("pgx", connectionString)
	if err != nil {
		_ = testcontainers.TerminateContainer(container)
		t.Fatal(err)
	}
	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		_ = testcontainers.TerminateContainer(container)
		t.Fatal(err)
	}
	migrator, err := database.NewMigrator(db)
	if err != nil {
		t.Fatal(err)
	}
	if err := migrator.Up(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = db.Close()
		_ = testcontainers.TerminateContainer(container)
	})
	return db
}
```

- [ ] **Step 4: 写数据库约束失败测试**

测试必须直接执行 SQL 验证：`USER@example.com` 和 `user@example.com` 触发唯一冲突；非法用户/任务状态被 CHECK 拒绝；有任务的用户无法删除；向下迁移删除两张表。

- [ ] **Step 5: 运行集成测试**

Run: `cd examples/go-task-api && GOPROXY=https://proxy.golang.org,direct go test -tags=integration ./tests -v`

Expected: PASS，实际启动 `postgres:18.4`，不能以跳过代替。

## Task 7: 实现 PostgreSQL Repository 和乐观锁

**Files:**
- Create: `examples/go-task-api/internal/user/repository_postgres.go`
- Create: `examples/go-task-api/internal/task/repository_postgres.go`
- Modify: `examples/go-task-api/tests/postgres_test.go`

- [ ] **Step 1: 先写仓储集成测试**

```go
func TestOptimisticUpdateAllowsOnlyOneWriter(t *testing.T) {
	db := startPostgres(t)
	users := user.NewPostgresRepository(db)
	created, err := users.Create(context.Background(), user.CreateParams{Name: "张三", Email: "user@example.com", Status: user.StatusActive})
	if err != nil { t.Fatal(err) }

	var wg sync.WaitGroup
	results := make(chan error, 2)
	for _, status := range []user.Status{user.StatusDisabled, user.StatusActive} {
		wg.Add(1)
		go func(status user.Status) {
			defer wg.Done()
			_, err := users.UpdateStatus(context.Background(), created.ID, status, created.Version)
			results <- err
		}(status)
	}
	wg.Wait()
	close(results)
	// 断言一个 nil、一个 user.ErrVersionConflict。
}
```

同时测试稳定排序分页、owner/status 筛选、context 取消、任务删除版本冲突。

- [ ] **Step 2: 运行测试确认 Repository 构造器不存在**

Run: `cd examples/go-task-api && go test -tags=integration ./tests -run 'Repository|Optimistic' -v`

Expected: FAIL，出现未定义构造器或方法。

- [ ] **Step 3: 实现用户 Repository**

创建使用 `insert ... returning`；列表使用参数化 SQL 和 `count(*) over()`；状态更新使用：

```sql
update users
set status = $1, version = version + 1, updated_at = now()
where id = $2 and version = $3
returning id, name, email, status, version, created_at, updated_at
```

零行时再查询 ID 是否存在，区分 `ErrNotFound` 与 `ErrVersionConflict`。pgx `PgError` 的 `23505` 映射 `ErrEmailConflict`，不依赖错误字符串。

- [ ] **Step 4: 实现任务 Repository**

所有列表排序固定为 `created_at desc, id desc`；分页使用 `limit` 和 `offset`；PUT、状态更新和 DELETE 都在 `where id = $id and version = $expected` 中校验版本。事务中只使用 `*sql.Tx` 的 `QueryContext/ExecContext`，不混用 `*sql.DB`。

- [ ] **Step 5: 运行单元、集成和 race 测试**

Run:

```bash
cd examples/go-task-api
go test ./...
go test -race ./...
go test -tags=integration ./tests -v
```

Expected: 全部 PASS。

## Task 8: 用 httptest 定义并实现用户和任务 Handler

**Files:**
- Create: `examples/go-task-api/internal/user/handler_test.go`
- Create: `examples/go-task-api/internal/user/handler.go`
- Create: `examples/go-task-api/internal/task/handler_test.go`
- Create: `examples/go-task-api/internal/task/handler.go`

- [ ] **Step 1: 写用户 Handler 失败测试**

覆盖：创建 201、列表 200、详情 200、非法 ID 400、未知 JSON 字段 400、字段校验 422、邮箱冲突 409、未找到 404、版本冲突 409、缺失 expectedVersion 422。

```go
func TestChangeStatusRequiresExpectedVersion(t *testing.T) {
	h := NewHandler(&stubService{})
	req := httptest.NewRequest(http.MethodPatch, "/api/users/1/status", strings.NewReader(`{"status":"DISABLED"}`))
	req.SetPathValue("id", "1")
	rec := httptest.NewRecorder()
	h.ChangeStatus(rec, req)
	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("got %d", rec.Code)
	}
}
```

- [ ] **Step 2: 写任务 Handler 失败测试**

覆盖：创建、列表筛选、详情、PUT、状态 PATCH、DELETE 查询参数版本、非法 ID、非法 page、错误 content type、超大 body 和 context deadline。

- [ ] **Step 3: 运行测试确认 Handler 尚未实现**

Run: `cd examples/go-task-api && go test ./internal/user ./internal/task -run Handler -v`

Expected: FAIL。

- [ ] **Step 4: 实现 DTO 和 HTTP 映射**

请求 DTO 的 `expectedVersion` 使用指针以区分缺失和 0：

```go
type changeStatusRequest struct {
	Status          Status `json:"status"`
	ExpectedVersion *int64 `json:"expectedVersion"`
}
```

创建返回 201 和 `Location`；删除返回 204 且无 JSON body；GET 不要求 Content-Type。POST、PUT、PATCH 的媒体类型和 body 上限由 Task 3 的 `RequireJSON` 统一处理，Handler 只调用严格 `DecodeJSON`。

- [ ] **Step 5: 运行 Handler 测试**

Run: `cd examples/go-task-api && go test ./internal/user ./internal/task -v`

Expected: PASS。

## Task 9: 组装路由、健康检查和服务生命周期

**Files:**
- Create: `examples/go-task-api/internal/app/router_test.go`
- Create: `examples/go-task-api/internal/app/router.go`
- Create: `examples/go-task-api/internal/app/server.go`
- Create: `examples/go-task-api/internal/app/app.go`
- Create: `examples/go-task-api/cmd/api/main.go`
- Create: `examples/go-task-api/cmd/migrate/main.go`
- Create: `examples/go-task-api/cmd/healthcheck/main.go`
- Create: `examples/go-task-api/go.sum`

- [ ] **Step 1: 写 404、405、Allow 和健康状态测试**

```go
func TestMethodNotAllowedIncludesAllow(t *testing.T) {
	h := MethodNotAllowed(http.MethodGet, http.MethodPost)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodDelete, "/api/users", nil))
	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("got %d", rec.Code)
	}
	if rec.Header().Get("Allow") != "GET, POST" {
		t.Fatalf("unexpected Allow: %q", rec.Header().Get("Allow"))
	}
}
```

readiness fake 在正常状态返回 200；开始 shutdown 或 PingContext 失败返回 503；liveness 始终只反映进程可服务。

- [ ] **Step 2: 运行路由测试确认失败**

Run: `cd examples/go-task-api && go test ./internal/app -v`

Expected: FAIL。

- [ ] **Step 3: 用 Go 1.22+ ServeMux pattern 注册路由**

```go
mux.HandleFunc("GET /health/live", health.Live)
mux.HandleFunc("GET /health/ready", health.Ready)
mux.HandleFunc("GET /api/users", users.List)
mux.HandleFunc("POST /api/users", users.Create)
mux.HandleFunc("GET /api/users/{id}", users.Get)
mux.HandleFunc("PATCH /api/users/{id}/status", users.ChangeStatus)
mux.HandleFunc("GET /api/tasks", tasks.List)
mux.HandleFunc("POST /api/tasks", tasks.Create)
mux.HandleFunc("GET /api/tasks/{id}", tasks.Get)
mux.HandleFunc("PUT /api/tasks/{id}", tasks.Update)
mux.HandleFunc("PATCH /api/tasks/{id}/status", tasks.ChangeStatus)
mux.HandleFunc("DELETE /api/tasks/{id}", tasks.Delete)
```

为 404 和 405 增加统一 JSON 包装；405 根据已注册方法生成稳定 `Allow`。

- [ ] **Step 4: 实现显式超时和优雅关闭**

```go
server := &http.Server{
	Addr:              cfg.HTTP.Addr,
	Handler:           handler,
	ReadHeaderTimeout: cfg.HTTP.ReadHeaderTimeout,
	ReadTimeout:       cfg.HTTP.ReadTimeout,
	WriteTimeout:      cfg.HTTP.WriteTimeout,
	IdleTimeout:       cfg.HTTP.IdleTimeout,
}
```

应用对外契约固定为：

```go
type App struct {
	Handler   http.Handler
	server    *http.Server
	readiness *Readiness
	db        *sql.DB
}

func New(cfg config.Config, logger *slog.Logger, db *sql.DB) *App
func (a *App) Run(ctx context.Context, shutdownTimeout time.Duration) error
```

`cmd/api` 使用 `signal.NotifyContext` 创建 ctx 并传给 `Run`；收到信号后 `Run` 先调用 readiness 的 `StartShutdown`，再用 `context.WithTimeout` 调用 `server.Shutdown`，最后关闭 DB。`http.ErrServerClosed` 不作为错误退出。

- [ ] **Step 5: 实现三个命令**

- `cmd/api`：加载配置、创建 slog、打开并 Ping 数据库、组装 App、运行服务。
- `cmd/migrate`：只接受 `up`、`down 1`、`version`，非法参数打印用法并返回非零状态。
- `cmd/healthcheck`：2 秒超时访问 `HEALTHCHECK_URL`，仅 2xx 返回 0。

- [ ] **Step 6: 整理依赖并运行全量 Go 静态验证**

Run:

```bash
cd examples/go-task-api
GOPROXY=https://proxy.golang.org,direct go mod tidy
test -z "$(gofmt -l .)"
go vet ./...
go test ./...
go test -race ./...
```

Expected: 全部成功。

## Task 10: 完成 API 集成测试、契约和排障文档

**Files:**
- Create: `examples/go-task-api/tests/api_integration_test.go`
- Create: `examples/go-task-api/API_CONTRACT.md`
- Create: `examples/go-task-api/TROUBLESHOOTING.md`

- [ ] **Step 1: 先写完整 API 集成流程**

测试顺序固定为：创建用户、重复邮箱冲突、创建任务、列表筛选、详情、合法 PUT、旧版本 PUT 冲突、合法状态推进、非法状态回退、旧版本 DELETE 冲突、合法 DELETE、再次查询 404。

```go
func TestTaskLifecycle(t *testing.T) {
	server := startAPI(t)
	user := postJSON[user.User](t, server.URL+"/api/users", map[string]any{
		"name": "张三", "email": "user@example.com",
	})
	createdTask := postJSON[task.Task](t, server.URL+"/api/tasks", map[string]any{
		"ownerId": user.ID, "title": "完成 Go 文档",
	})
	if createdTask.OwnerID != user.ID || createdTask.Status != task.StatusTodo || createdTask.Version != 0 {
		t.Fatalf("unexpected task: %+v", createdTask)
	}
}
```

同一测试文件定义以下两个完整 helper 契约：

```go
func startAPI(t *testing.T) *httptest.Server {
	t.Helper()
	db := startPostgres(t)
	cfg := config.Config{
		Environment: "test",
		LogLevel:    slog.LevelError,
		HTTP: config.HTTPConfig{
			RequestTimeout:  2 * time.Second,
			ShutdownTimeout: 2 * time.Second,
		},
	}
	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))
	application := app.New(cfg, logger, db)
	server := httptest.NewServer(application.Handler)
	t.Cleanup(server.Close)
	return server
}

func postJSON[T any](t *testing.T, url string, input any) T {
	t.Helper()
	body, err := json.Marshal(input)
	if err != nil {
		t.Fatal(err)
	}
	response, err := http.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", response.StatusCode)
	}
	var envelope struct {
		Data      T      `json:"data"`
		RequestID string `json:"requestId"`
	}
	if err := json.NewDecoder(response.Body).Decode(&envelope); err != nil {
		t.Fatal(err)
	}
	if envelope.RequestID == "" {
		t.Fatal("missing request id")
	}
	return envelope.Data
}
```

- [ ] **Step 2: 验证协议边界**

增加独立测试覆盖非法/溢出 ID、未知字段、两个 JSON 值、超大 body、405 + Allow、415、请求 deadline、panic 脱敏，以及缺失 expectedVersion 的三个写接口。

- [ ] **Step 3: 编写完整 API 契约**

`API_CONTRACT.md` 对 12 条路由逐条列出：方法、路径、Header、查询参数、请求 JSON、成功 JSON、错误码、版本并发语义和 curl 示例。DELETE 明确使用 `?expectedVersion=`，所有列表说明稳定排序。

- [ ] **Step 4: 编写可执行排障手册**

`TROUBLESHOOTING.md` 包含：端口占用、配置解析失败、数据库连不上、迁移 dirty、连接池等待、版本冲突、请求超时、容器不健康和优雅关闭超时。每项给出命令、预期证据和修复后验证。

- [ ] **Step 5: 运行集成测试**

Run: `cd examples/go-task-api && go test -tags=integration ./... -v`

Expected: PASS，Tests run 不为 0，Docker 中无残留测试容器。

## Task 11: 容器化并完成真实 smoke test

**Files:**
- Create: `examples/go-task-api/.dockerignore`
- Create: `examples/go-task-api/Dockerfile`
- Create: `examples/go-task-api/compose.yaml`
- Create: `examples/go-task-api/README.md`
- Modify: `.gitignore`

- [ ] **Step 1: 编写多阶段 Dockerfile**

```dockerfile
FROM golang:1.26.5-bookworm AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/api ./cmd/api \
 && CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/migrate ./cmd/migrate \
 && CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/healthcheck ./cmd/healthcheck

FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /app
COPY --from=build /out/api /out/migrate /out/healthcheck /app/
USER nonroot:nonroot
ENTRYPOINT ["/app/api"]
```

- [ ] **Step 2: 编写只绑定本机的 Compose**

数据库使用 `postgres:18.4`，API 和数据库端口分别绑定 `127.0.0.1:8080:8080`、`127.0.0.1:5432:5432`。`migrate` 使用相同镜像并覆盖 entrypoint 为 `/app/migrate up`；API 依赖数据库健康且迁移成功；API 健康检查执行 `/app/healthcheck`；`stop_grace_period: 20s`。

- [ ] **Step 3: 编写 README 完整流程**

README 包含本地 Go 启动、Docker 启动、迁移、测试、race、Fuzz、构建、12 条接口快速验证、关闭和清理命令，并明确“无认证，只能用于学习和本地环境”。

- [ ] **Step 4: 更新忽略规则**

```gitignore
# Go 示例构建与测试产物
examples/go-task-api/bin/
examples/go-task-api/coverage.out
examples/go-task-api/*.test
examples/go-task-api/.env
```

- [ ] **Step 5: 执行真实 smoke test**

Run:

```bash
cd examples/go-task-api
docker compose up --build -d
curl --fail http://127.0.0.1:8080/health/live
curl --fail http://127.0.0.1:8080/health/ready
docker compose ps
docker compose stop api
docker compose down -v
```

Expected: 两个健康接口为 2xx，`docker compose ps` 显示 API healthy；停止日志包含 shutdown 开始和完成；最终无该项目容器和 volume。

## Task 12: 重写 Go 导览和基础章节

**Files:**
- Modify: `docs/go/introduction.md`
- Modify: `docs/go/setup-modules.md`
- Modify: `docs/go/syntax-types.md`
- Modify: `docs/go/interfaces-composition.md`
- Modify: `docs/go/errors-logging-config.md`
- Modify: `docs/go/concurrency.md`
- Modify: `docs/go/context-http.md`
- Modify: `docs/go/database-transaction.md`
- Modify: `docs/go/testing.md`
- Modify: `docs/go/project-deployment.md`
- Modify: `docs/go/performance.md`
- Modify: `docs/go/troubleshooting.md`
- Modify: `docs/go/grpc-service-communication.md`

- [ ] **Step 1: 固定每篇教学结构**

除导览和排错页外，每篇都包含：

```markdown
## 适合谁看
## 先建立心智模型
## 从最小示例开始
## 放进真实项目
## 常见错误与根因
## 验证清单
## 下一步学习
```

- [ ] **Step 2: 更新版本和选型事实**

导览写明 Go 1.26.5 的核对日期与官方发布链接；项目基线解释 `go 1.26.0` 与补丁工具链的区别；明确标准库主路线及 Gin/Chi/GORM/sqlc 决策表。

- [ ] **Step 3: 补足语言难点**

`syntax-types.md` 增加 slice 扩容与共享、map 并发、指针接收者和零值；`interfaces-composition.md` 增加动态类型/值与 nil interface；`errors-logging-config.md` 增加错误链和 slog 字段契约。

- [ ] **Step 4: 补足并发和服务边界**

`concurrency.md` 增加 channel 所有权、关闭规则、sync 选择、race 和 goroutine 泄漏；`context-http.md` 修复重复 `Timeout` 函数声明并解释 deadline middleware 不额外启动 handler goroutine。

- [ ] **Step 5: 补足数据库、测试和交付**

`database-transaction.md` 对齐示例连接池、乐观锁和迁移；`testing.md` 使用示例真实测试；`project-deployment.md` 对齐 distroless、探针和 shutdown；`performance.md` 建立指标到 profile 的证据链；`troubleshooting.md` 链接问题库。

- [ ] **Step 6: 运行文档检查**

Run: `npm run docs:check`

Expected: 通过，无新增缺失章节或内部路由错误。

## Task 13: 完成 26 张 Go 核心图和两张教学图片

**Prerequisite:** 在执行本任务前完成 `2026-07-21-existing-modules-visual-enrichment.md` 的 Task 1-3，使 `DocFigure`、视觉资产登记和自动检查已经可用。

**Files:**
- Modify: `docs/go/visual-guide.md`
- Create: `docs/public/images/go/go-api-request-journey.webp`
- Create: `docs/public/images/go/go-concurrency-workshop.webp`
- Modify/Create: `docs/contribute/visual-asset-register.md`

- [ ] **Step 1: 按规格重写 26 个编号章节**

每节使用：问题、Mermaid、逐步解释、项目对应位置、误区和自测。运行以下命令确认数量：

```bash
test "$(rg -c '^```mermaid' docs/go/visual-guide.md)" -eq 26
test "$(rg -c '^## [0-9]+\.' docs/go/visual-guide.md)" -eq 26
```

- [ ] **Step 2: 使用 built-in imagegen 生成请求旅程图**

最终 prompt：

```text
Use case: scientific-educational
Asset type: Chinese Go programming tutorial illustration
Primary request: create a precise cutaway educational illustration of one web request traveling through a Go HTTP server, business service, repository, PostgreSQL database, then returning a response; show cancellation flowing backward as a separate red signal
Style/medium: clean flat technical illustration with subtle depth, crisp geometric shapes, professional documentation quality
Composition/framing: wide 16:9, five clearly separated stages from left to right, generous padding
Color palette: mint green, cyan, charcoal, white, small amber and red accents
Constraints: no text, no letters, no logos, no code, no watermark; visually distinguish normal data flow from cancellation; do not imitate a real product screenshot
Avoid: decorative gradients, dark background, fantasy imagery, illegible micro-details
```

将 built-in 输出复制到 `tmp/imagegen/go-api-request-journey.png`，再执行 `cwebp -quiet -q 86 tmp/imagegen/go-api-request-journey.png -o docs/public/images/go/go-api-request-journey.webp`。检查尺寸、清晰度和无错误文字后删除临时 PNG。

- [ ] **Step 3: 使用 built-in imagegen 生成并发工作坊类比图**

最终 prompt：

```text
Use case: scientific-educational
Asset type: Chinese Go concurrency tutorial illustration
Primary request: create an educational visual analogy for goroutines and channels: several independent workers receive bounded tasks from one conveyor channel, one coordinator owns and closes the channel, and a clearly visible stop signal lets every worker exit cleanly
Style/medium: clean flat technical illustration, friendly but professional, crisp edges
Composition/framing: wide 16:9, coordinator on the left, bounded conveyor in the center, workers on the right, visible exit paths
Color palette: fresh green, cyan, navy-gray, white, small amber stop accent
Constraints: no text, no letters, no logos, no code, no watermark; exactly one channel owner; no worker stranded; this is explicitly an analogy, not a runtime architecture diagram
Avoid: crowded scene, decorative characters, dark background, gradients, arbitrary machinery
```

将输出复制到 `tmp/imagegen/go-concurrency-workshop.png`，使用 `cwebp -quiet -q 86` 转为 `docs/public/images/go/go-concurrency-workshop.webp`，验证后删除临时 PNG。图注明确“这是帮助理解所有权和退出条件的类比，准确规则以相邻 Mermaid 和代码为准”。

- [ ] **Step 4: 登记两个生成资产**

登记本地路径、使用页面、中文 alt、caption、完整 prompt、`built-in image_gen`、实际生成日期和人工核对结论。正文使用 `DocFigure`，不使用裸 Markdown 图片。

- [ ] **Step 5: 在桌面和移动端检查图示**

Expected: 26 个非空 Mermaid SVG、2 张 WebP 返回 200，无 `.mermaid-diagram__error`，页面无横向溢出。

## Task 14: 把项目文档连接到真实源码

**Files:**
- Modify: `docs/go/http-api-project-from-zero.md`

- [ ] **Step 1: 修正现有内容错误和技术分叉**

删除 ER 图重复字段；固定 PostgreSQL 18、标准库 ServeMux、database/sql、pgx、migrate 和 slog；不再写“任选框架或迁移工具”。

- [ ] **Step 2: 用真实代码导入替换漂移片段**

使用以下形式导入 `go.mod`、迁移、配置、httpx、Service、Repository、Handler、路由、测试、Dockerfile 和 compose：

```markdown
<<< ../../examples/go-task-api/internal/platform/httpx/json.go{go}
```

导入路径必须通过 VitePress build 验证。

- [ ] **Step 3: 按从零顺序组织项目**

顺序固定为：目标与接口、创建模块、目录、配置、数据库迁移、错误契约、用户领域、任务领域、路由和中间件、测试、Docker、smoke、排障、扩展框架时机。

- [ ] **Step 4: 添加图、真实请求和验收表**

至少包含 6 张项目专用 Mermaid：总体架构、请求链、中间件顺序、数据模型、乐观锁并发、关闭顺序。curl 输出必须来自 Task 11 的真实运行结果。

- [ ] **Step 5: 构建验证代码导入**

Run: `npm run docs:build`

Expected: 所有 `<<<` 导入成功，无文件不存在错误。

## Task 15: 编写 16 个 Go 真实项目问题

**Files:**
- Create: `docs/projects/issues-go.md`

- [ ] **Step 1: 创建统一证据模板和快速定位表**

```markdown
## 证据记录模板

| 项目 | 记录内容 |
| --- | --- |
| 现象 | 时间、接口、错误比例、影响范围 |
| 运行信息 | Go 版本、提交、镜像、配置摘要 |
| 证据 | 指标、日志、profile、trace、SQL、goroutine stack |
| 假设 | 可被证伪的根因判断 |
| 修复 | 最小改动与风险 |
| 回归 | 命令、期望结果、观察窗口 |
```

- [ ] **Step 2: 按规格写满 16 个编号问题**

每个问题必须包含：现象、最小复现、错误方向、证据采集、根因、修复、回归测试、预防清单。使用命令确认：

```bash
test "$(rg -c '^## 问题 [0-9]+：' docs/projects/issues-go.md)" -eq 16
```

- [ ] **Step 3: 添加至少 10 张非重复图**

重点给 nil interface、slice 持有、map race、goroutine leak、channel deadlock、context 传播、连接池等待、事务泄漏、丢失更新和关闭顺序配图。

- [ ] **Step 4: 连接 Go 排错页和后端总问题库**

开头链接 `/go/troubleshooting`；末尾链接 `/go/performance`、`/go/testing`、`/projects/issues-backend` 和 `/roadmap/go-practice`。

- [ ] **Step 5: 检查路由和图数**

Run: `npm run docs:check && test "$(rg -c '^```mermaid' docs/projects/issues-go.md)" -ge 10`

Expected: 通过。

## Task 16: 编写 12 个递进练习

**Files:**
- Create: `docs/roadmap/go-practice.md`

- [ ] **Step 1: 建立统一实验规则**

固定 Go 1.26.5、PostgreSQL 18.4、命令记录、失败证据和清理要求；说明哪些练习需要 Docker。

- [ ] **Step 2: 按规格写满 12 个练习**

每个练习包含目标、起始条件、任务步骤、限制、验证命令、通过标准、常见失败、进阶挑战。确认数量：

```bash
test "$(rg -c '^## 练习 [0-9]+：' docs/roadmap/go-practice.md)" -eq 12
```

- [ ] **Step 3: 添加至少 7 张练习图**

为模块依赖、错误链、channel 退出、context、HTTP 中间件、事务并发和最终交付流水线配图。

- [ ] **Step 4: 最终练习引用真实示例**

练习 12 直接使用 `examples/go-task-api`，验收命令与 README、API_CONTRACT 一致，不复制第二套项目。

- [ ] **Step 5: 运行数量和路由检查**

Run: `npm run docs:check && test "$(rg -c '^```mermaid' docs/roadmap/go-practice.md)" -ge 7`

Expected: 通过。

## Task 17: 更新导航、状态和跨模块入口

**Files:**
- Modify: `docs/.vitepress/config.ts`
- Modify: `docs/contribute/module-status.md`
- Modify: `docs/technologies/index.md`
- Modify: `docs/technologies/expansion-plan.md`
- Modify: `docs/projects/real-world-issues.md`
- Modify: `docs/projects/issues-backend.md`
- Modify: `docs/roadmap/introduction.md`
- Modify: `docs/roadmap/practice-labs.md`
- Modify: `docs/roadmap/reading-guide.md`
- Modify: `README.md`

- [ ] **Step 1: 在 Go 侧边栏加入两个入口**

```ts
const goIssuesRoute = '/projects/issues-go'
const goPracticeRoute = '/roadmap/go-practice'

{ text: 'Go 真实项目问题库', link: goIssuesRoute },
{ text: 'Go 专项练习', link: goPracticeRoute },
```

顺序放在性能诊断之后、常见问题之前。

- [ ] **Step 2: 更新技术库和模块状态**

Go 状态写明 26 组图解、可运行任务 API、16 类问题、12 个练习、PostgreSQL 18、race/Fuzz/Testcontainers 和容器交付；下一步只保留微服务治理、云原生和 CLI 等明确非目标。

- [ ] **Step 3: 更新问题库与学习路线交叉链接**

后端问题总览增加 Go 专项；路线增加“Go 导览 -> 图解 -> HTTP 项目 -> 练习 -> 问题库”的明确顺序。

- [ ] **Step 4: 更新 README 示例结构和验证命令**

增加 `examples/java-admin-api`、`examples/go-task-api` 说明，以及各自独立测试命令；不把 Docker 集成测试描述为默认无依赖测试。

- [ ] **Step 5: 运行文档和差异检查**

Run: `npm run docs:check && git diff --check`

Expected: 通过，配置路由都存在。

## Task 18: 完成最终自动化与浏览器验收

**Files:**
- Verify only; only edit files when validation finds a specific defect.

- [ ] **Step 1: 运行 Go 完整验证**

```bash
cd examples/go-task-api
test -z "$(gofmt -l .)"
go vet ./...
go test ./...
go test -race ./...
go test -run=^$ -fuzz=FuzzDecodeJSON -fuzztime=10s ./internal/platform/httpx
go test -tags=integration ./... -v
docker run --rm -v "$PWD:/src" -w /src golang:1.26.5-bookworm sh -lc 'go test ./... && go test -race ./...'
build_dir="$(mktemp -d)"
trap 'rm -rf "$build_dir"' EXIT
go build -o "$build_dir/api" ./cmd/api
go build -o "$build_dir/migrate" ./cmd/migrate
go build -o "$build_dir/healthcheck" ./cmd/healthcheck
```

Expected: 所有命令成功，integration tests 实际运行。

- [ ] **Step 2: 运行 Compose 业务 smoke**

除健康检查外，使用 curl 创建用户和任务，执行一次合法更新和一次旧版本冲突；断言冲突为 409 且 `error.code` 为 `TASK_VERSION_CONFLICT`。最后发送终止信号并检查关闭日志。

- [ ] **Step 3: 运行文档生产验证**

Run: `cd /Users/yokry/Documents/Codex/2026-07-01/yok-study-web && npm run docs:check && npm run docs:build && git diff --check`

Expected: 通过；只允许记录已经存在的 `env` 高亮 fallback 和大 chunk 警告。

- [ ] **Step 4: 浏览器验证四个核心页面**

在 `1440x900` 和 `390x844` 打开：

- `/go/visual-guide`
- `/go/http-api-project-from-zero`
- `/projects/issues-go`
- `/roadmap/go-practice`

断言：Mermaid 分别为 26、至少 6、至少 10、至少 7；全部 SVG 非空；2 张 WebP 返回 200；没有 `.mermaid-diagram__error`；`document.documentElement.scrollWidth === document.documentElement.clientWidth`；控制台无错误。

- [ ] **Step 5: 检查工作区和残留资源**

Run:

```bash
git status --short
docker ps -a --filter name=go-task-api
docker volume ls --filter name=go-task-api
```

Expected: 只有预期源码和文档改动；无 `bin`、测试二进制、覆盖率文件、target 或测试容器残留。

- [ ] **Step 6: 记录最终结果但不提交**

最终交付说明列出：示例测试数量、26/16/12 内容计数、浏览器视口结果、图片路径与 prompts、已有警告和未执行项。等待用户明确要求后再提交、推送或部署。
