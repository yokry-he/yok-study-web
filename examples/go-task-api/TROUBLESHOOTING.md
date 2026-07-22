# Go Task API 排障手册

这份手册按“现象 -> 证据 -> 根因 -> 修复 -> 验证”组织。不要看到错误就立即重启或删除数据；先保留 request ID、时间点、命令输出和最近日志，才能判断问题发生在配置、网络、数据库、迁移、业务并发还是关闭流程。

> 命令默认在 `examples/go-task-api` 目录执行。示例中的密码只用于本地学习环境，不要把真实连接串粘贴到工单、聊天记录或 Git。

## 0. 先收集最小证据

```bash
date -u
go version
docker version
docker compose version
docker compose ps
docker compose logs --since=10m api migrate postgres
curl -i http://127.0.0.1:8080/health/live
curl -i http://127.0.0.1:8080/health/ready
```

记录以下内容：

- 失败请求的方法、路径、HTTP 状态和 `X-Request-ID`。
- 问题第一次发生的时间与是否稳定复现。
- `/health/live` 与 `/health/ready` 是否同时失败。
- `docker compose ps` 中每个服务的状态和退出码。
- 最近一次迁移版本，但不要记录含密码的 `DATABASE_URL`。

快速判断：

| live | ready | 通常表示 |
|---|---|---|
| 连接失败 | 连接失败 | API 未启动、端口错误或进程已退出 |
| 200 | 503 | API 存活，但数据库不可用、迁移未就绪或正在关闭 |
| 200 | 200 | 基础设施正常，继续检查具体请求和业务错误码 |

## 1. 端口已被占用

### 现象

- 本地启动立即返回 `bind: address already in use`。
- Compose 提示无法绑定 `127.0.0.1:8080` 或 `127.0.0.1:5432`。
- 浏览器访问到的内容与当前代码不一致，实际命中了旧进程。

### 收集证据

```bash
lsof -nP -iTCP:8080 -sTCP:LISTEN
lsof -nP -iTCP:5432 -sTCP:LISTEN
docker ps --format 'table {{.ID}}\t{{.Names}}\t{{.Ports}}'
```

预期证据是占用端口的 PID、进程名或容器名。不要在没有确认归属时直接 `kill -9`。

### 修复

若占用者是本项目旧 Compose：

```bash
docker compose down
```

若占用者是确认可以停止的本地进程，先发送可优雅处理的信号：

```bash
kill -TERM <PID>
```

如果必须并行运行多个项目，可以覆盖 PostgreSQL 的主机侧端口：

```bash
POSTGRES_PORT=55432 docker compose up -d
```

容器内数据库地址仍是 `postgres:5432`。若 8080 被占用，需要先停止本项目旧进程，或在单独的本地覆盖文件中调整 API 主机端口；容器内端口和 `HTTP_ADDR=:8080` 不需要跟着改。

### 验证

```bash
lsof -nP -iTCP:8080 -sTCP:LISTEN
docker compose up -d
curl --fail-with-body http://127.0.0.1:8080/health/live
```

## 2. 配置解析失败

### 现象

启动时出现以下类型错误：

- `DATABASE_URL is required`
- `invalid DB_MAX_OPEN_CONNS: must be a positive integer`
- `invalid HTTP_REQUEST_TIMEOUT`
- `invalid LOG_LEVEL`

### 收集证据

只检查变量是否存在和非敏感值，不要打印完整连接串：

```bash
test -n "$DATABASE_URL" && echo 'DATABASE_URL 已设置' || echo 'DATABASE_URL 未设置'
printf 'HTTP_ADDR=%s\n' "$HTTP_ADDR"
printf 'HTTP_REQUEST_TIMEOUT=%s\n' "$HTTP_REQUEST_TIMEOUT"
printf 'DB_MAX_OPEN_CONNS=%s\n' "$DB_MAX_OPEN_CONNS"
printf 'DB_MAX_IDLE_CONNS=%s\n' "$DB_MAX_IDLE_CONNS"
printf 'LOG_LEVEL=%s\n' "$LOG_LEVEL"
```

数字配置必须是正整数；持续时间使用 Go duration，如 `500ms`、`10s`、`5m`；日志级别只接受 `debug`、`info`、`warn`、`error`。

### 修复

从模板创建仅供本地使用的环境文件，然后按需修改：

```bash
cp .env.example .env
set -a
source .env
set +a
```

`.env` 已加入忽略规则，不应提交。Compose 会自行提供容器内连接串；不要把宿主机的 `127.0.0.1` 数据库地址直接用于 API 容器。

### 验证

```bash
go run ./cmd/migrate version
go run ./cmd/healthcheck
```

第二条命令前需设置 `HEALTHCHECK_URL=http://127.0.0.1:8080/health/ready`。

## 3. 数据库连接失败

### 现象

- API 启动失败，或 `/health/live` 为 200、`/health/ready` 为 503。
- 日志显示数据库连接、Ping、DNS、认证或 TLS 失败。
- 本地可以连接，API 容器却不能连接。

### 收集证据

```bash
docker compose ps postgres
docker compose logs --since=10m postgres
docker compose exec postgres pg_isready -U app -d taskdb
docker compose exec postgres psql -U app -d taskdb -c 'select current_database(), current_user, version();'
```

再确认运行环境使用的主机名：

- 宿主机运行 API：`127.0.0.1:5432`。
- Compose 中 API 连接数据库：`postgres:5432`。
- 容器中的 `127.0.0.1` 指向容器自身，不是 PostgreSQL 服务。

### 修复

1. PostgreSQL 未启动：`docker compose up -d postgres`。
2. 数据库尚未 ready：等待 `pg_isready` 成功，不要只依赖容器处于 running。
3. 用户名、密码或数据库名错误：对照 `compose.yaml` 与 `DATABASE_URL`。
4. 容器中使用了错误主机名：改为 Compose 服务名 `postgres`。
5. 本地连接使用了不匹配的 TLS：学习环境连接串明确使用 `sslmode=disable`。

### 验证

```bash
docker compose exec postgres pg_isready -U app -d taskdb
curl --fail-with-body http://127.0.0.1:8080/health/ready
```

## 4. 迁移处于 dirty 状态

### 现象

- `migrate up` 报告 dirty database version。
- `migrate version` 输出类似 `version=1 dirty=true`。
- 数据库已启动，但 API 因表或字段缺失而返回 500。

### 收集证据

```bash
go run ./cmd/migrate version
docker compose run --rm migrate version
docker compose exec postgres psql -U app -d taskdb \
  -c 'select version, dirty from schema_migrations;'
```

保存失败迁移的日志，并检查数据库结构是否只执行了一部分：

```bash
docker compose exec postgres psql -U app -d taskdb -c '\dt+'
docker compose exec postgres psql -U app -d taskdb -c '\d+ users'
docker compose exec postgres psql -U app -d taskdb -c '\d+ tasks'
```

### 根因

迁移执行中途断电、进程被强杀、SQL 执行失败或人工修改 `schema_migrations` 都可能留下 dirty 标记。dirty 不是“再跑一次就会自动修好”的普通错误。

### 修复

本地一次性学习数据可直接重建，风险最低：

```bash
docker compose down -v
docker compose up --build -d
```

需要保留数据时：

1. 先备份数据库。
2. 对照对应版本的 up/down SQL 检查每个对象是否已执行。
3. 手工修复到一个明确、完整的版本状态。
4. 使用受控的迁移工具修正版本标记。

不要在未核对结构时直接把 `dirty` 改为 `false`，也不要在生产数据上执行 `docker compose down -v`。

### 验证

```bash
go run ./cmd/migrate version
```

预期当前版本为 `version=1 dirty=false`，随后运行：

```bash
go test -tags=integration ./tests -run TestPostgresMigrationContract -v
```

## 5. 数据库连接池等待或请求变慢

### 现象

- 并发上升后请求耗时突然增加。
- `/health/ready` 偶尔超时，但 PostgreSQL 本身仍在运行。
- 日志中多个请求接近 `HTTP_REQUEST_TIMEOUT`。

### 收集证据

先看当前连接与等待事件：

```bash
docker compose exec postgres psql -U app -d taskdb -c "
select pid, state, wait_event_type, wait_event,
       now() - query_start as query_age,
       left(query, 120) as query
from pg_stat_activity
where datname = 'taskdb'
order by query_start;"
```

检查配置关系：

```bash
printf 'DB_MAX_OPEN_CONNS=%s DB_MAX_IDLE_CONNS=%s\n' \
  "$DB_MAX_OPEN_CONNS" "$DB_MAX_IDLE_CONNS"
```

`DB_MAX_IDLE_CONNS` 不能大于 `DB_MAX_OPEN_CONNS`。连接数过小会排队，盲目调大又可能耗尽 PostgreSQL 的 `max_connections`。

### 修复

1. 先定位慢 SQL、锁等待或未及时结束的事务。
2. 确认查询使用迁移中定义的过滤与分页索引。
3. 缩短事务范围，不要在持有事务时等待网络或用户输入。
4. 根据真实并发与数据库容量小幅调整池大小。
5. 保留 `HTTP_REQUEST_TIMEOUT`，避免无限等待占满连接。

不要把“增加连接数”当成第一修复手段，它会掩盖慢查询和锁竞争。

### 验证

重复同一负载，比较访问日志中的 `duration`，并再次检查 `pg_stat_activity`：不应持续存在异常长事务或锁等待。

## 6. 乐观锁版本冲突

### 现象

接口返回：

- `409 USER_VERSION_CONFLICT`
- `409 TASK_VERSION_CONFLICT`

这通常不是服务故障，而是两个客户端基于同一旧版本同时修改资源。

### 收集证据

```bash
curl -sS "$BASE_URL/api/tasks/<TASK_ID>" | jq '.data | {id, status, version, updatedAt}'
```

对照失败请求发送的 `expectedVersion`。如果服务端版本更大，说明已有其他写请求成功。

### 修复

1. 重新 GET 最新资源。
2. 比较本地编辑内容与服务端最新内容。
3. 由用户确认覆盖、合并或放弃。
4. 使用最新 `version` 重新提交一次。

不要对 409 进行无限自动重试。无条件重试会覆盖其他用户刚完成的修改。

### 验证

```bash
CURRENT=$(curl -sS "$BASE_URL/api/tasks/<TASK_ID>" | jq -r '.data.version')
curl --fail-with-body \
  -X PATCH "$BASE_URL/api/tasks/<TASK_ID>/status" \
  -H 'Content-Type: application/json' \
  -d "{\"status\":\"DOING\",\"expectedVersion\":$CURRENT}"
```

目标状态还必须符合状态转换表；版本正确但转换非法会返回 `TASK_INVALID_TRANSITION`。

## 7. 请求超时

### 现象

- API 返回 `504 DEADLINE_EXCEEDED`。
- 客户端先断开，服务端日志中出现取消或写响应失败。
- 请求耗时接近 `HTTP_REQUEST_TIMEOUT`。

### 收集证据

使用固定 request ID 重放，便于日志关联：

```bash
curl -i --max-time 15 \
  -H 'X-Request-ID: timeout-investigation-001' \
  "$BASE_URL/api/tasks?page=1&pageSize=100"
docker compose logs --since=5m api | grep 'timeout-investigation-001'
```

同时检查 PostgreSQL 等待：

```bash
docker compose exec postgres psql -U app -d taskdb -c "
select pid, state, wait_event_type, wait_event, now() - query_start as age,
       left(query, 120) as query
from pg_stat_activity
where datname = 'taskdb' and state <> 'idle';"
```

### 修复

- 查询慢：分析执行计划、过滤条件和索引。
- 锁等待：缩短持锁事务，定位阻塞者。
- 请求工作量过大：使用分页，减小 `pageSize`。
- 下游偶发超时：只对幂等读取采用有上限、带退避的重试。
- 超时配置确实过小：基于观测数据调整，而不是直接设为无限。

### 验证

相同请求应在新的超时预算内完成；访问日志中的状态应为 2xx，数据库中不应残留已取消请求对应的长查询。

## 8. Compose 容器不健康

### 现象

`docker compose ps` 显示 API 或 PostgreSQL 为 `unhealthy`、`restarting` 或已退出。

### 收集证据

```bash
docker compose ps
docker compose logs --since=10m api migrate postgres
docker inspect --format '{{json .State.Health}}' go-task-api-api-1 | jq
```

若容器名不同，先从 `docker compose ps -q api` 获取 ID：

```bash
API_ID=$(docker compose ps -q api)
docker inspect --format '{{json .State.Health}}' "$API_ID" | jq
```

常见证据：

- PostgreSQL healthcheck 失败：数据库还没完成恢复或认证配置不一致。
- migrate 非零退出：迁移失败，API 因 `service_completed_successfully` 不会启动。
- API healthcheck 失败：`HEALTHCHECK_URL` 错误、API 尚未监听或 readiness 为 503。

### 修复

1. 先修最上游的 PostgreSQL。
2. 再单独运行迁移并观察输出：`docker compose run --rm migrate up`。
3. 迁移成功后重建 API：`docker compose up --build -d api`。
4. 不要通过删除 healthcheck 或依赖条件绕过问题。

### 验证

```bash
docker compose ps
curl --fail-with-body http://127.0.0.1:8080/health/live
curl --fail-with-body http://127.0.0.1:8080/health/ready
```

API 最终应显示 `healthy`，两个接口都返回 2xx。

## 9. 优雅关闭超时

### 现象

- `docker compose stop api` 超过预期时间。
- 关闭期间新请求仍进入业务处理。
- 达到 `HTTP_SHUTDOWN_TIMEOUT` 后连接被强制关闭。

### 收集证据

在一个终端持续观察状态和日志：

```bash
docker compose logs -f api
```

另一个终端发送停止命令并记录耗时：

```bash
time docker compose stop api
```

关闭开始后，readiness 应先变为 503；已经进入的请求可以在关闭预算内完成。检查仍在执行的数据库语句：

```bash
docker compose exec postgres psql -U app -d taskdb -c "
select pid, state, now() - query_start as age, left(query, 120) as query
from pg_stat_activity
where datname = 'taskdb' and state <> 'idle';"
```

### 修复

- 找出不响应 request context 取消的 Handler 或数据库调用。
- 所有外部调用使用请求 context，并设置比总关闭预算更短的超时。
- 缩短长事务，避免关闭时等待不可控的后台工作。
- Compose 的 `stop_grace_period` 应大于 `HTTP_SHUTDOWN_TIMEOUT`，给应用留下清理时间。
- 只有超时后才强制关闭连接；不要把 `SIGKILL` 作为常规停止方式。

### 验证

```bash
docker compose up -d api
curl --fail-with-body http://127.0.0.1:8080/health/ready
time docker compose stop api
docker compose ps -a api
```

进程应在关闭预算内退出，随后可以重新启动且数据库连接正常。

## 10. 排障结束后的回归检查

```bash
go test ./...
go test -race ./...
go test -tags=integration ./... -v
docker compose up --build -d
curl --fail-with-body http://127.0.0.1:8080/health/live
curl --fail-with-body http://127.0.0.1:8080/health/ready
docker compose down -v
```

最后确认没有测试或 Compose 残留资源：

```bash
docker ps -a --filter name=go-task-api
docker volume ls --filter name=go-task-api
```
