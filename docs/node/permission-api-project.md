# Node 权限 API 从零到项目

## 适合谁看

适合已经理解 JavaScript、TypeScript、HTTP、SQL 和基本鉴权概念，但还没有独立完成过可运行后端项目的人。本章从空目录开始，最终得到一个可以迁移、初始化、启动、测试、构建和容器化的用户角色权限 API。

## 这个项目解决什么

你将实现：

- 存活与就绪检查。
- 登录、退出和当前用户。
- 用户分页、创建和启用/停用。
- 角色列表和用户角色替换。
- 动作级权限码校验。
- PostgreSQL 迁移、回滚和种子数据。
- 统一错误响应、结构化日志和 request id。
- 事务、审计日志、会话撤销和优雅停机。
- 单元测试与 Fastify `inject()` API 测试。

第一版不实现菜单树、刷新 Token、Redis 缓存和微服务。这些能力只有在权限 API 本身正确后才值得增加。

## 完整请求链路

```mermaid
flowchart TD
  A["HTTP 请求"] --> B["Fastify 路由匹配与 request id"]
  B --> C["请求体解析与 JSON Schema 校验"]
  C --> D["读取 Bearer 会话"]
  D --> E["加载用户、角色和权限"]
  E --> F["Route handler"]
  F --> G["Service：业务规则与权限"]
  G --> H["Repository：参数化 SQL"]
  H --> I["PostgreSQL"]
  G --> J["事务内审计日志"]
  F --> K["响应序列化"]
  K --> L["状态码、request id、结构化日志"]
```

Fastify 使用 plugin、hook、decorator 和 route lifecycle，不套用 Express 的 `req, res, next` middleware 模型。本章也不额外创建只转发参数的 Controller；route handler 已经是传输层适配器，Service 和 Repository 才是必须稳定的边界。

## 阶段地图

```mermaid
flowchart LR
  A["1. 脚手架"] --> B["2. API 契约"]
  B --> C["3. 数据模型与迁移"]
  C --> D["4. 配置和数据库"]
  D --> E["5. 密码与会话"]
  E --> F["6. Repository"]
  F --> G["7. Service 与事务"]
  G --> H["8. Fastify 路由"]
  H --> I["9. 测试"]
  I --> J["10. 部署与故障注入"]
```

每完成一个阶段就执行对应验证，不要等全部文件写完才第一次启动。

## 技术基线

| 能力 | 选择 | 原因 |
| --- | --- | --- |
| 运行时 | Node.js 24 LTS | 生产使用受支持 LTS；本地、CI、镜像保持一致 |
| Web 框架 | Fastify 5 | 原生 Promise 错误传播、Schema、日志、inject 测试和插件封装 |
| 语言 | TypeScript 7 | 严格类型检查；仍由 `tsc` 构建，不把运行时类型擦除当检查器 |
| 数据库 | PostgreSQL 18 | 约束、事务、锁、JSONB 和成熟诊断能力 |
| 数据访问 | `pg` / node-postgres | 直接看见连接池、参数化 SQL 和同一连接事务 |
| 会话 | 随机不透明 Bearer Token | 数据库可撤销，能够真正实现退出和禁用用户会话失效 |
| 密码 | Node `crypto.scrypt` | 不保存明文；示例不引入本地二进制依赖 |
| 测试 | `node:test` + `app.inject()` | Node 内置测试运行器；不监听真实端口即可测试 HTTP 契约 |

依赖版本以 2026-07-21 验证结果为基线。创建项目时提交 `package-lock.json`，不要让 CI 每次解析出不同依赖树。

## 阶段 1：创建项目

### 1.1 环境检查

```bash
node --version
npm --version
docker --version
```

Node 主版本应为 `v24`。前半程可只用 Docker 启动本地 PostgreSQL，容器化阶段再用 Compose 运行完整 API；也可以连接自己管理的开发数据库。

### 1.2 初始化和安装

```bash
mkdir node-permission-api
cd node-permission-api
npm init -y
npm install fastify@^5.10.0 @fastify/rate-limit@^11.1.0 pg@^8.22.0
npm install -D typescript@^7.0.2 tsx@^4.23.1 @types/node@^24.13.3 @types/pg@^8.20.0
```

将 `package.json` 调整为：

```json
{
  "name": "node-permission-api",
  "version": "1.0.0",
  "private": true,
  "type": "module",
  "engines": {
    "node": ">=24 <25"
  },
  "scripts": {
    "dev": "node --env-file=.env --import=tsx --watch src/server.ts",
    "typecheck": "tsc --noEmit",
    "build": "tsc -p tsconfig.json",
    "start": "node --env-file=.env dist/src/server.js",
    "test": "node --import=tsx --test \"tests/**/*.test.ts\"",
    "db:migrate": "node --env-file=.env --import=tsx scripts/migrate.ts up",
    "db:rollback": "node --env-file=.env --import=tsx scripts/migrate.ts down",
    "db:seed": "node --env-file=.env --import=tsx scripts/seed.ts",
    "db:migrate:prod": "node --env-file=.env dist/scripts/migrate.js up"
  },
  "dependencies": {
    "@fastify/rate-limit": "^11.1.0",
    "fastify": "^5.10.0",
    "pg": "^8.22.0"
  },
  "devDependencies": {
    "@types/node": "^24.13.3",
    "@types/pg": "^8.20.0",
    "tsx": "^4.23.1",
    "typescript": "^7.0.2"
  }
}
```

### 1.3 TypeScript 配置

`tsconfig.json`：

```json
{
  "compilerOptions": {
    "target": "ES2023",
    "module": "NodeNext",
    "moduleResolution": "NodeNext",
    "rootDir": ".",
    "outDir": "dist",
    "strict": true,
    "noUncheckedIndexedAccess": true,
    "exactOptionalPropertyTypes": true,
    "verbatimModuleSyntax": true,
    "erasableSyntaxOnly": true,
    "forceConsistentCasingInFileNames": true,
    "skipLibCheck": true,
    "sourceMap": true
  },
  "include": ["src/**/*.ts", "scripts/**/*.ts"]
}
```

源码中的相对导入写 `.js` 后缀，例如 `import { loadConfig } from './config.js'`。TypeScript 在开发时会解析对应 `.ts`，构建后 Node 可以直接加载生成的 `.js`。

### 1.4 最终目录

```text
node-permission-api/
├─ sql/
│  ├─ 001_init.up.sql
│  └─ 001_init.down.sql
├─ scripts/
│  ├─ migrate.ts
│  └─ seed.ts
├─ src/
│  ├─ app.ts
│  ├─ config.ts
│  ├─ db.ts
│  ├─ errors.ts
│  ├─ password.ts
│  ├─ repositories.ts
│  ├─ services.ts
│  ├─ server.ts
│  └─ types.ts
├─ tests/
│  ├─ app.test.ts
│  └─ password.test.ts
├─ .env.example
├─ compose.yaml
├─ Dockerfile
├─ package.json
├─ package-lock.json
└─ tsconfig.json
```

## 阶段 2：先固定 API 契约

### 2.1 路由表

| 方法与路径 | 是否公开 | 权限 | 成功状态 | 用途 |
| --- | --- | --- | ---: | --- |
| `GET /health/live` | 是 | - | 200 | 进程存活 |
| `GET /health/ready` | 是 | - | 200/503 | 是否愿意接流量 |
| `POST /api/auth/login` | 是 | - | 200 | 创建短期会话 |
| `POST /api/auth/logout` | 否 | 已认证 | 200 | 撤销当前会话 |
| `GET /api/me` | 否 | 已认证 | 200 | 当前用户与权限 |
| `GET /api/users` | 否 | `user:read` | 200 | 分页查询用户 |
| `POST /api/users` | 否 | `user:create` | 201 | 创建用户并分配角色 |
| `PATCH /api/users/:id/status` | 否 | `user:update` | 200 | 启用或停用用户 |
| `GET /api/roles` | 否 | `user:read` | 200 | 表单使用的角色列表 |
| `PUT /api/users/:id/roles` | 否 | `role:grant` | 200 | 原子替换用户角色 |

### 2.2 成功与失败响应

```ts
export type ApiSuccess<T> = {
  ok: true
  data: T
  requestId: string
}

export type ApiFailure = {
  ok: false
  error: {
    code: string
    message: string
  }
  requestId: string
}
```

状态码语义：

| 状态码 | 业务含义 | 客户端下一步 |
| ---: | --- | --- |
| 400 | 请求结构或字段非法 | 修正输入 |
| 401 | 没有有效会话 | 重新登录 |
| 403 | 身份有效但动作无权限 | 显示无权限，不重复登录 |
| 404 | 资源不存在 | 返回列表或刷新数据 |
| 409 | 唯一性、状态或并发冲突 | 保留输入并解释冲突 |
| 429 | 登录尝试过多 | 等待后重试 |
| 500 | 未预期内部错误 | 用 request id 联系支持 |
| 503 | 应用暂时未就绪 | 平台摘流并稍后重试 |

### 2.3 核心类型

`src/types.ts`：

```ts
export type UserStatus = 'enabled' | 'disabled'

export type PermissionCode =
  | 'user:read'
  | 'user:create'
  | 'user:update'
  | 'role:grant'

export type AuthContext = {
  sessionId: string
  userId: string
  email: string
  displayName: string
  permissionCodes: string[]
}

export type UserDTO = {
  id: string
  email: string
  displayName: string
  status: UserStatus
  roleCodes: string[]
  createdAt: string
}

export type RoleDTO = {
  id: string
  code: string
  name: string
}

export type CreateUserInput = {
  email: string
  displayName: string
  password: string
  roleIds: string[]
}
```

数据库 `uuid` 在 JSON 中使用 string，不转换为 JavaScript `number`，避免安全整数和跨语言精度问题。

## 阶段 3：PostgreSQL 数据模型

### 3.1 关系图

```mermaid
erDiagram
  USERS ||--o{ USER_ROLES : assigned
  ROLES ||--o{ USER_ROLES : contains
  ROLES ||--o{ ROLE_PERMISSIONS : grants
  PERMISSIONS ||--o{ ROLE_PERMISSIONS : included
  USERS ||--o{ SESSIONS : owns
  USERS ||--o{ AUDIT_LOGS : acts

  USERS {
    uuid id PK
    text email
    text password_hash
    text status
    timestamptz created_at
  }
  SESSIONS {
    uuid id PK
    char token_hash UK
    timestamptz expires_at
    timestamptz revoked_at
  }
```

### 3.2 Up migration

`sql/001_init.up.sql`：

```sql
CREATE TABLE users (
  id uuid PRIMARY KEY,
  email text NOT NULL,
  display_name text NOT NULL,
  password_hash text NOT NULL,
  status text NOT NULL DEFAULT 'enabled'
    CHECK (status IN ('enabled', 'disabled')),
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX users_email_lower_uk ON users ((lower(email)));
CREATE INDEX users_status_created_idx ON users (status, created_at DESC);

CREATE TABLE roles (
  id uuid PRIMARY KEY,
  code text NOT NULL UNIQUE,
  name text NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE permissions (
  code text PRIMARY KEY,
  name text NOT NULL
);

CREATE TABLE user_roles (
  user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  role_id uuid NOT NULL REFERENCES roles(id) ON DELETE RESTRICT,
  PRIMARY KEY (user_id, role_id)
);

CREATE INDEX user_roles_role_idx ON user_roles (role_id, user_id);

CREATE TABLE role_permissions (
  role_id uuid NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
  permission_code text NOT NULL
    REFERENCES permissions(code) ON DELETE RESTRICT,
  PRIMARY KEY (role_id, permission_code)
);

CREATE TABLE sessions (
  id uuid PRIMARY KEY,
  user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  token_hash char(64) NOT NULL UNIQUE,
  expires_at timestamptz NOT NULL,
  revoked_at timestamptz,
  created_at timestamptz NOT NULL DEFAULT now(),
  CHECK (expires_at > created_at)
);

CREATE INDEX sessions_active_user_idx
  ON sessions (user_id, expires_at DESC)
  WHERE revoked_at IS NULL;

CREATE TABLE audit_logs (
  id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  actor_user_id uuid REFERENCES users(id) ON DELETE SET NULL,
  action text NOT NULL,
  target_type text NOT NULL,
  target_id text NOT NULL,
  detail jsonb NOT NULL DEFAULT '{}'::jsonb,
  request_id text NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX audit_logs_target_idx
  ON audit_logs (target_type, target_id, created_at DESC);

INSERT INTO permissions (code, name) VALUES
  ('user:read', '查看用户'),
  ('user:create', '创建用户'),
  ('user:update', '更新用户'),
  ('role:grant', '分配角色');

COMMENT ON TABLE users IS '后台登录用户，status 控制是否允许认证';
COMMENT ON COLUMN users.email IS '登录邮箱，使用 lower(email) 唯一索引实现大小写不敏感唯一性';
COMMENT ON COLUMN users.password_hash IS 'scrypt 参数、salt 和派生结果，不存明文密码';
COMMENT ON TABLE sessions IS '不透明访问会话，只保存 token SHA-256，不保存原始 token';
COMMENT ON COLUMN sessions.revoked_at IS '退出或禁用用户时写入；非空会话不可继续使用';
COMMENT ON TABLE audit_logs IS '关键权限写操作的不可变审计记录';
```

这里使用 PostgreSQL 的 `uuid`、`timestamptz`、`jsonb`、表达式索引和 `COMMENT ON`，不要混入 MySQL 的 `AUTO_INCREMENT`、`tinyint` 或列内 `COMMENT` 语法。

### 3.3 Down migration

`sql/001_init.down.sql`：

```sql
DROP TABLE IF EXISTS audit_logs;
DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS role_permissions;
DROP TABLE IF EXISTS user_roles;
DROP TABLE IF EXISTS permissions;
DROP TABLE IF EXISTS roles;
DROP TABLE IF EXISTS users;
```

Down 会删除数据，只用于本地或经过备份、停写和影响评估的受控回滚。生产优先采用向前修复迁移，不要为了回滚应用版本就自动删除新表或新列。

### 3.4 为什么这些约束必须在数据库

- 两个并发请求都可能先查到“邮箱不存在”，只有唯一索引能最终阻止重复数据。
- 外键防止出现指向不存在用户或角色的关系。
- 复合主键防止同一角色被重复分配。
- `ON DELETE RESTRICT` 防止仍被引用的权限或角色静默消失。
- 审计日志与业务写入位于同一事务，审计失败时业务也回滚。

## 阶段 4：配置、连接池和事务

### 4.1 环境变量

`.env.example`：

```env
NODE_ENV=development
HOST=0.0.0.0
PORT=3000
DATABASE_URL=postgres://app:app@127.0.0.1:5432/permission_api
DB_POOL_MAX=10
SESSION_TTL_MINUTES=30
ADMIN_EMAIL=admin@example.com
ADMIN_PASSWORD=ChangeThisAdmin123!
```

`.env` 不提交。真实密钥由部署平台注入；种子管理员密码只用于本地开发，并在首次真实登录后替换。

`src/config.ts`：

```ts
export type AppConfig = {
  nodeEnv: 'development' | 'test' | 'production'
  host: string
  port: number
  databaseUrl: string
  dbPoolMax: number
  sessionTtlMinutes: number
}

function required(name: string): string {
  const value = process.env[name]?.trim()
  if (!value) throw new Error(`Missing environment variable: ${name}`)
  return value
}

function positiveInteger(name: string, fallback: number, maximum: number): number {
  const raw = process.env[name]
  if (!raw) return fallback
  const value = Number(raw)
  if (!Number.isInteger(value) || value <= 0 || value > maximum) {
    throw new Error(`${name} must be an integer between 1 and ${maximum}`)
  }
  return value
}

export function loadConfig(): AppConfig {
  const nodeEnv = process.env.NODE_ENV ?? 'development'
  if (!['development', 'test', 'production'].includes(nodeEnv)) {
    throw new Error('NODE_ENV must be development, test or production')
  }

  return {
    nodeEnv: nodeEnv as AppConfig['nodeEnv'],
    host: process.env.HOST?.trim() || '0.0.0.0',
    port: positiveInteger('PORT', 3000, 65535),
    databaseUrl: required('DATABASE_URL'),
    dbPoolMax: positiveInteger('DB_POOL_MAX', 10, 100),
    sessionTtlMinutes: positiveInteger('SESSION_TTL_MINUTES', 30, 1440)
  }
}
```

### 4.2 数据库和事务 helper

`src/db.ts`：

```ts
import { Pool, type PoolClient } from 'pg'
import type { AppConfig } from './config.js'

export type DbExecutor = Pick<PoolClient, 'query'>

export function createPool(config: AppConfig): Pool {
  const pool = new Pool({
    connectionString: config.databaseUrl,
    max: config.dbPoolMax,
    connectionTimeoutMillis: 3_000,
    idleTimeoutMillis: 30_000,
    statement_timeout: 5_000,
    application_name: 'node-permission-api'
  })

  pool.on('error', (error) => {
    console.error('Unexpected idle PostgreSQL client error', error)
  })

  return pool
}

export async function withTransaction<T>(
  pool: Pool,
  work: (client: PoolClient) => Promise<T>
): Promise<T> {
  const client = await pool.connect()
  try {
    await client.query('BEGIN')
    const result = await work(client)
    await client.query('COMMIT')
    return result
  } catch (error) {
    try {
      await client.query('ROLLBACK')
    } catch (rollbackError) {
      console.error('Transaction rollback failed', rollbackError)
    }
    throw error
  } finally {
    client.release()
  }
}
```

事务中的 `BEGIN`、业务 SQL、审计 SQL、`COMMIT/ROLLBACK` 必须使用同一个 `PoolClient`。`pool.query()` 适合单条查询，不适合跨多条语句维持事务。

## 阶段 5：错误、密码和会话

### 5.1 统一应用错误

`src/errors.ts`：

```ts
export class AppError extends Error {
  readonly statusCode: number
  readonly code: string

  constructor(statusCode: number, code: string, message: string) {
    super(message)
    this.name = 'AppError'
    this.statusCode = statusCode
    this.code = code
  }
}

export function mapDatabaseError(error: unknown): never {
  const code = typeof error === 'object' && error !== null && 'code' in error
    ? String(error.code)
    : ''

  if (code === '23505') {
    throw new AppError(409, 'RESOURCE_CONFLICT', '数据已存在或发生唯一性冲突')
  }
  if (code === '23503') {
    throw new AppError(409, 'REFERENCE_CONFLICT', '关联数据不存在或仍在使用')
  }
  throw error
}
```

PostgreSQL `23505` 是唯一约束冲突，`23503` 是外键冲突。不要把数据库原始错误文本返回给客户端，因为它可能包含表名、约束名和 SQL 细节。

### 5.2 密码哈希

`src/password.ts`：

```ts
import { createHash, randomBytes, scrypt, timingSafeEqual } from 'node:crypto'

const KEY_LENGTH = 64
const N = 131_072
const R = 8
const P = 1
const MAX_MEMORY = 256 * 1024 * 1024

function deriveKey(
  password: string,
  salt: Buffer,
  cost: { n: number; r: number; p: number }
): Promise<Buffer> {
  return new Promise((resolve, reject) => {
    scrypt(password, salt, KEY_LENGTH, {
      N: cost.n,
      r: cost.r,
      p: cost.p,
      maxmem: MAX_MEMORY
    }, (error, key) => {
      if (error) reject(error)
      else resolve(key)
    })
  })
}

export async function hashPassword(password: string): Promise<string> {
  const salt = randomBytes(16)
  const key = await deriveKey(password, salt, { n: N, r: R, p: P })
  return ['scrypt', N, R, P, salt.toString('base64url'), key.toString('base64url')].join('$')
}

export async function verifyPassword(password: string, stored: string): Promise<boolean> {
  const [algorithm, nText, rText, pText, saltText, keyText] = stored.split('$')
  if (algorithm !== 'scrypt' || !nText || !rText || !pText || !saltText || !keyText) {
    return false
  }

  const n = Number(nText)
  const r = Number(rText)
  const p = Number(pText)
  if (!Number.isInteger(n) || !Number.isInteger(r) || !Number.isInteger(p)) return false
  if (n < 16_384 || n > N || r !== R || p < 1 || p > 4) return false

  const expected = Buffer.from(keyText, 'base64url')
  if (expected.length !== KEY_LENGTH) return false
  const actual = await deriveKey(password, Buffer.from(saltText, 'base64url'), { n, r, p })
  return timingSafeEqual(actual, expected)
}

export function createSessionToken(): { raw: string; hash: string } {
  const raw = randomBytes(32).toString('base64url')
  const hash = createHash('sha256').update(raw).digest('hex')
  return { raw, hash }
}

export function hashSessionToken(raw: string): string {
  return createHash('sha256').update(raw).digest('hex')
}
```

数据库只保存会话 token 的 SHA-256。原始 token 只在登录成功响应中出现一次；数据库泄露者不能直接拿哈希当作 Bearer token 使用。

密码哈希是昂贵操作，登录路由还必须限流。生产环境应根据机器性能压测参数，并评估 Argon2id、集中身份系统、MFA、弱密码检测和凭据泄露检测；不要机械复制一个成本参数到所有环境。

### 5.3 会话状态机

```mermaid
stateDiagram-v2
  [*] --> Active: "登录成功"
  Active --> Revoked: "主动退出"
  Active --> Revoked: "用户被停用"
  Active --> Expired: "超过 expires_at"
  Active --> Active: "认证并加载最新权限"
  Revoked --> [*]
  Expired --> [*]
```

本项目每个受保护请求都从数据库加载当前用户和权限，因此角色变更立即生效。代价是每次请求多一次数据库查询；后续可以引入短 TTL 缓存，但必须同时设计权限变更失效。

## 阶段 6：Repository

Repository 只负责数据库读写，不决定当前用户是否有权限。所有 SQL 使用占位符，不能拼接用户输入。

`src/repositories.ts`：

```ts
import { randomUUID } from 'node:crypto'
import type { Pool, QueryResultRow } from 'pg'
import type { DbExecutor } from './db.js'
import type {
  AuthContext,
  CreateUserInput,
  RoleDTO,
  UserDTO,
  UserStatus
} from './types.js'

type LoginRow = QueryResultRow & {
  id: string
  email: string
  display_name: string
  password_hash: string
  status: UserStatus
}

type AuthRow = QueryResultRow & {
  session_id: string
  user_id: string
  email: string
  display_name: string
  permission_codes: string[]
}

type UserRow = QueryResultRow & {
  id: string
  email: string
  display_name: string
  status: UserStatus
  role_codes: string[]
  created_at: Date
}

type RoleRow = QueryResultRow & {
  id: string
  code: string
  name: string
}

function toUser(row: UserRow): UserDTO {
  return {
    id: row.id,
    email: row.email,
    displayName: row.display_name,
    status: row.status,
    roleCodes: row.role_codes,
    createdAt: row.created_at.toISOString()
  }
}

export class AuthRepository {
  private readonly pool: Pool

  constructor(pool: Pool) {
    this.pool = pool
  }

  async findUserForLogin(email: string): Promise<LoginRow | null> {
    const result = await this.pool.query<LoginRow>(`
      SELECT id, email, display_name, password_hash, status
      FROM users
      WHERE lower(email) = lower($1)
    `, [email])
    return result.rows[0] ?? null
  }

  async createSession(userId: string, tokenHash: string, expiresAt: Date): Promise<void> {
    await this.pool.query(`
      INSERT INTO sessions (id, user_id, token_hash, expires_at)
      VALUES ($1, $2, $3, $4)
    `, [randomUUID(), userId, tokenHash, expiresAt])
  }

  async findContext(tokenHash: string): Promise<AuthContext | null> {
    const result = await this.pool.query<AuthRow>(`
      SELECT
        s.id AS session_id,
        u.id AS user_id,
        u.email,
        u.display_name,
        COALESCE(
          array_agg(DISTINCT rp.permission_code)
            FILTER (WHERE rp.permission_code IS NOT NULL),
          '{}'::text[]
        ) AS permission_codes
      FROM sessions s
      JOIN users u ON u.id = s.user_id
      LEFT JOIN user_roles ur ON ur.user_id = u.id
      LEFT JOIN role_permissions rp ON rp.role_id = ur.role_id
      WHERE s.token_hash = $1
        AND s.revoked_at IS NULL
        AND s.expires_at > now()
        AND u.status = 'enabled'
      GROUP BY s.id, u.id
    `, [tokenHash])

    const row = result.rows[0]
    return row ? {
      sessionId: row.session_id,
      userId: row.user_id,
      email: row.email,
      displayName: row.display_name,
      permissionCodes: row.permission_codes
    } : null
  }

  async revokeSession(sessionId: string): Promise<void> {
    await this.pool.query(`
      UPDATE sessions SET revoked_at = COALESCE(revoked_at, now()) WHERE id = $1
    `, [sessionId])
  }

  async revokeUserSessions(executor: DbExecutor, userId: string): Promise<void> {
    await executor.query(`
      UPDATE sessions
      SET revoked_at = COALESCE(revoked_at, now())
      WHERE user_id = $1 AND revoked_at IS NULL
    `, [userId])
  }
}

export class UserRepository {
  private readonly pool: Pool

  constructor(pool: Pool) {
    this.pool = pool
  }

  async list(query: { q: string; page: number; pageSize: number }): Promise<{
    items: UserDTO[]
    total: number
  }> {
    const search = `%${query.q}%`
    const countResult = await this.pool.query<{ total: string }>(`
      SELECT count(*)::text AS total
      FROM users
      WHERE $1 = '%%' OR email ILIKE $1 OR display_name ILIKE $1
    `, [search])

    const result = await this.pool.query<UserRow>(`
      SELECT
        u.id,
        u.email,
        u.display_name,
        u.status,
        u.created_at,
        COALESCE(
          array_agg(DISTINCT r.code) FILTER (WHERE r.code IS NOT NULL),
          '{}'::text[]
        ) AS role_codes
      FROM users u
      LEFT JOIN user_roles ur ON ur.user_id = u.id
      LEFT JOIN roles r ON r.id = ur.role_id
      WHERE $1 = '%%' OR u.email ILIKE $1 OR u.display_name ILIKE $1
      GROUP BY u.id
      ORDER BY u.created_at DESC, u.id
      LIMIT $2 OFFSET $3
    `, [search, query.pageSize, (query.page - 1) * query.pageSize])

    return {
      items: result.rows.map(toUser),
      total: Number(countResult.rows[0]?.total ?? 0)
    }
  }

  async create(
    executor: DbExecutor,
    input: CreateUserInput,
    passwordHash: string
  ): Promise<UserDTO> {
    const id = randomUUID()
    const result = await executor.query<UserRow>(`
      INSERT INTO users (id, email, display_name, password_hash)
      VALUES ($1, lower($2), $3, $4)
      RETURNING id, email, display_name, status, created_at, '{}'::text[] AS role_codes
    `, [id, input.email, input.displayName, passwordHash])

    const row = result.rows[0]
    if (!row) throw new Error('INSERT users returned no row')
    return toUser(row)
  }

  async lockById(executor: DbExecutor, userId: string): Promise<UserStatus | null> {
    const result = await executor.query<{ status: UserStatus }>(`
      SELECT status FROM users WHERE id = $1 FOR UPDATE
    `, [userId])
    return result.rows[0]?.status ?? null
  }

  async updateStatus(
    executor: DbExecutor,
    userId: string,
    status: UserStatus
  ): Promise<void> {
    await executor.query(`
      UPDATE users SET status = $2, updated_at = now() WHERE id = $1
    `, [userId, status])
  }

  async replaceRoles(executor: DbExecutor, userId: string, roleIds: string[]): Promise<void> {
    await executor.query('DELETE FROM user_roles WHERE user_id = $1', [userId])
    if (roleIds.length > 0) {
      await executor.query(`
        INSERT INTO user_roles (user_id, role_id)
        SELECT $1, unnest($2::uuid[])
      `, [userId, roleIds])
    }
  }

  async findRoleCodes(executor: DbExecutor, roleIds: string[]): Promise<string[]> {
    const result = await executor.query<{ code: string }>(`
      SELECT code FROM roles WHERE id = ANY($1::uuid[]) ORDER BY code
    `, [roleIds])
    return result.rows.map((row) => row.code)
  }

  async listRoles(): Promise<RoleDTO[]> {
    const result = await this.pool.query<RoleRow>(`
      SELECT id, code, name FROM roles ORDER BY name, id
    `)
    return result.rows.map((row) => ({ id: row.id, code: row.code, name: row.name }))
  }
}

export class AuditRepository {
  async create(executor: DbExecutor, input: {
    actorUserId: string
    action: string
    targetType: string
    targetId: string
    detail: Record<string, unknown>
    requestId: string
  }): Promise<void> {
    await executor.query(`
      INSERT INTO audit_logs (
        actor_user_id, action, target_type, target_id, detail, request_id
      ) VALUES ($1, $2, $3, $4, $5::jsonb, $6)
    `, [
      input.actorUserId,
      input.action,
      input.targetType,
      input.targetId,
      JSON.stringify(input.detail),
      input.requestId
    ])
  }
}
```

分页查询使用稳定的 `created_at, id` 排序。规模增大后应评估游标分页；本项目保留页码分页，便于先理解 URL、`LIMIT/OFFSET` 和总数契约。

## 阶段 7：Service、权限和事务

### 7.1 为什么权限放在 Service

只在路由 hook 检查权限，意味着后台任务、CLI 或另一个入口调用 Service 时可能绕过授权。Service 是用例边界，应再次要求明确的 `AuthContext`。

`src/services.ts`：

```ts
import type { Pool } from 'pg'
import type { AppConfig } from './config.js'
import { mapDatabaseError, AppError } from './errors.js'
import { withTransaction } from './db.js'
import {
  createSessionToken,
  hashPassword,
  hashSessionToken,
  verifyPassword
} from './password.js'
import {
  AuditRepository,
  AuthRepository,
  UserRepository
} from './repositories.js'
import type {
  AuthContext,
  CreateUserInput,
  PermissionCode,
  RoleDTO,
  UserDTO,
  UserStatus
} from './types.js'

export function assertPermission(
  actor: AuthContext,
  permission: PermissionCode
): void {
  if (!actor.permissionCodes.includes(permission)) {
    throw new AppError(403, 'FORBIDDEN', `缺少权限：${permission}`)
  }
}

export interface AuthServicePort {
  login(email: string, password: string): Promise<{
    accessToken: string
    expiresAt: string
  }>
  authenticate(rawToken: string): Promise<AuthContext>
  logout(actor: AuthContext): Promise<void>
}

export interface UserServicePort {
  list(actor: AuthContext, query: {
    q: string
    page: number
    pageSize: number
  }): Promise<{ items: UserDTO[]; total: number }>
  listRoles(actor: AuthContext): Promise<RoleDTO[]>
  create(
    actor: AuthContext,
    input: CreateUserInput,
    requestId: string
  ): Promise<UserDTO>
  setStatus(
    actor: AuthContext,
    userId: string,
    status: UserStatus,
    requestId: string
  ): Promise<void>
  replaceRoles(
    actor: AuthContext,
    userId: string,
    roleIds: string[],
    requestId: string
  ): Promise<void>
}

export class AuthService implements AuthServicePort {
  private readonly repository: AuthRepository
  private readonly config: AppConfig

  constructor(repository: AuthRepository, config: AppConfig) {
    this.repository = repository
    this.config = config
  }

  async login(email: string, password: string): Promise<{
    accessToken: string
    expiresAt: string
  }> {
    const user = await this.repository.findUserForLogin(email)
    const passwordValid = user
      ? await verifyPassword(password, user.password_hash)
      : await verifyPassword(
          password,
          'scrypt$131072$8$1$MDAwMDAwMDAwMDAwMDAwMA$AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA'
        )

    if (!user || !passwordValid || user.status !== 'enabled') {
      throw new AppError(401, 'INVALID_CREDENTIALS', '邮箱或密码不正确')
    }

    const token = createSessionToken()
    const expiresAt = new Date(Date.now() + this.config.sessionTtlMinutes * 60_000)
    await this.repository.createSession(user.id, token.hash, expiresAt)
    return { accessToken: token.raw, expiresAt: expiresAt.toISOString() }
  }

  async authenticate(rawToken: string): Promise<AuthContext> {
    if (rawToken.length < 32 || rawToken.length > 200) {
      throw new AppError(401, 'UNAUTHORIZED', '登录已失效，请重新登录')
    }
    const context = await this.repository.findContext(hashSessionToken(rawToken))
    if (!context) {
      throw new AppError(401, 'UNAUTHORIZED', '登录已失效，请重新登录')
    }
    return context
  }

  async logout(actor: AuthContext): Promise<void> {
    await this.repository.revokeSession(actor.sessionId)
  }
}

export class UserService implements UserServicePort {
  private readonly pool: Pool
  private readonly users: UserRepository
  private readonly auth: AuthRepository
  private readonly audit: AuditRepository

  constructor(
    pool: Pool,
    users: UserRepository,
    auth: AuthRepository,
    audit: AuditRepository
  ) {
    this.pool = pool
    this.users = users
    this.auth = auth
    this.audit = audit
  }

  async list(actor: AuthContext, query: {
    q: string
    page: number
    pageSize: number
  }): Promise<{ items: UserDTO[]; total: number }> {
    assertPermission(actor, 'user:read')
    return this.users.list(query)
  }

  async listRoles(actor: AuthContext): Promise<RoleDTO[]> {
    assertPermission(actor, 'user:read')
    return this.users.listRoles()
  }

  async create(
    actor: AuthContext,
    input: CreateUserInput,
    requestId: string
  ): Promise<UserDTO> {
    assertPermission(actor, 'user:create')
    const roleIds = [...new Set(input.roleIds)]
    const passwordHash = await hashPassword(input.password)

    try {
      return await withTransaction(this.pool, async (client) => {
        const roleCodes = await this.users.findRoleCodes(client, roleIds)
        if (roleCodes.length !== roleIds.length) {
          throw new AppError(400, 'INVALID_ROLE', '包含不存在的角色')
        }

        const user = await this.users.create(client, { ...input, roleIds }, passwordHash)
        await this.users.replaceRoles(client, user.id, roleIds)
        await this.audit.create(client, {
          actorUserId: actor.userId,
          action: 'user.create',
          targetType: 'user',
          targetId: user.id,
          detail: { email: user.email, roleIds },
          requestId
        })
        return { ...user, roleCodes }
      })
    } catch (error) {
      return mapDatabaseError(error)
    }
  }

  async setStatus(
    actor: AuthContext,
    userId: string,
    status: UserStatus,
    requestId: string
  ): Promise<void> {
    assertPermission(actor, 'user:update')
    if (actor.userId === userId && status === 'disabled') {
      throw new AppError(409, 'CANNOT_DISABLE_SELF', '不能停用当前登录用户')
    }

    await withTransaction(this.pool, async (client) => {
      const currentStatus = await this.users.lockById(client, userId)
      if (!currentStatus) throw new AppError(404, 'USER_NOT_FOUND', '用户不存在')
      if (currentStatus === status) return

      await this.users.updateStatus(client, userId, status)
      if (status === 'disabled') {
        await this.auth.revokeUserSessions(client, userId)
      }
      await this.audit.create(client, {
        actorUserId: actor.userId,
        action: 'user.status.update',
        targetType: 'user',
        targetId: userId,
        detail: { before: currentStatus, after: status },
        requestId
      })
    })
  }

  async replaceRoles(
    actor: AuthContext,
    userId: string,
    roleIdsInput: string[],
    requestId: string
  ): Promise<void> {
    assertPermission(actor, 'role:grant')
    const roleIds = [...new Set(roleIdsInput)]

    await withTransaction(this.pool, async (client) => {
      if (!await this.users.lockById(client, userId)) {
        throw new AppError(404, 'USER_NOT_FOUND', '用户不存在')
      }
      if ((await this.users.findRoleCodes(client, roleIds)).length !== roleIds.length) {
        throw new AppError(400, 'INVALID_ROLE', '包含不存在的角色')
      }

      await this.users.replaceRoles(client, userId, roleIds)
      await this.audit.create(client, {
        actorUserId: actor.userId,
        action: 'user.roles.replace',
        targetType: 'user',
        targetId: userId,
        detail: { roleIds },
        requestId
      })
    })
  }
}
```

登录失败时，无论邮箱不存在、密码错误还是用户停用，对外都返回同一错误，降低账户枚举风险。示例对不存在用户也执行一次低成本兼容哈希；高安全系统应采用经过审查的固定假哈希并通过基准测试控制时序差异。

角色替换先用 `SELECT ... FOR UPDATE` 锁定目标用户行。对同一个用户的并发授权会串行执行，避免两个“先删后插”事务无序覆盖。更复杂的协作界面还可以增加 `version` 并使用乐观锁返回 409。

## 阶段 8：Fastify 应用

### 8.1 应用依赖和 Schema

`buildApp()` 不监听端口，便于测试。`server.ts` 只负责装配真实依赖和进程生命周期。

`src/app.ts`：

```ts
import { randomUUID } from 'node:crypto'
import rateLimit from '@fastify/rate-limit'
import Fastify, {
  type FastifyError,
  type FastifyInstance,
  type FastifyRequest
} from 'fastify'
import type { AppConfig } from './config.js'
import { AppError } from './errors.js'
import type { AuthServicePort, UserServicePort } from './services.js'
import type { AuthContext, CreateUserInput, UserStatus } from './types.js'

type AppDependencies = {
  config: AppConfig
  auth: AuthServicePort
  users: UserServicePort
  readiness: {
    isClosing: () => boolean
    check: () => Promise<void>
  }
  logger?: boolean
}

type LoginBody = { email: string; password: string }
type UserListQuery = { q?: string; page?: number; pageSize?: number }
type StatusBody = { status: UserStatus }
type RoleBody = { roleIds: string[] }

const uuidPattern = '^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[1-8][0-9a-fA-F]{3}-[89abAB][0-9a-fA-F]{3}-[0-9a-fA-F]{12}$'
const emailPattern = '^[^\\s@]+@[^\\s@]+\\.[^\\s@]+$'

const failureSchema = {
  type: 'object',
  additionalProperties: false,
  required: ['ok', 'error', 'requestId'],
  properties: {
    ok: { const: false },
    error: {
      type: 'object',
      additionalProperties: false,
      required: ['code', 'message'],
      properties: {
        code: { type: 'string' },
        message: { type: 'string' }
      }
    },
    requestId: { type: 'string' }
  }
} as const

function successSchema(data: Record<string, unknown>) {
  return {
    type: 'object',
    additionalProperties: false,
    required: ['ok', 'data', 'requestId'],
    properties: {
      ok: { const: true },
      data,
      requestId: { type: 'string' }
    }
  } as const
}

function ok<T>(request: FastifyRequest, data: T) {
  return { ok: true as const, data, requestId: request.id }
}

function readBearerToken(request: FastifyRequest): string {
  const authorization = request.headers.authorization
  if (!authorization) throw new AppError(401, 'UNAUTHORIZED', '请先登录')
  const [scheme, token, extra] = authorization.trim().split(/\s+/)
  if (scheme?.toLowerCase() !== 'bearer' || !token || extra) {
    throw new AppError(401, 'UNAUTHORIZED', '认证信息格式不正确')
  }
  return token
}

export async function buildApp(deps: AppDependencies): Promise<FastifyInstance> {
  const app = Fastify({
    logger: deps.logger === false ? false : {
      level: deps.config.nodeEnv === 'development' ? 'debug' : 'info',
      redact: ['req.headers.authorization', 'req.headers.cookie', 'password']
    },
    genReqId: () => randomUUID(),
    bodyLimit: 1_048_576,
    requestTimeout: 15_000
  })

  await app.register(rateLimit, { global: false })

  async function authenticate(request: FastifyRequest): Promise<AuthContext> {
    return deps.auth.authenticate(readBearerToken(request))
  }

  app.setNotFoundHandler((request, reply) => {
    return reply.status(404).send({
      ok: false,
      error: { code: 'ROUTE_NOT_FOUND', message: '接口不存在' },
      requestId: request.id
    })
  })

  app.setErrorHandler((error: FastifyError, request, reply) => {
    if (error.validation) {
      return reply.status(400).send({
        ok: false,
        error: { code: 'VALIDATION_ERROR', message: '请求参数不正确' },
        requestId: request.id
      })
    }

    if (error instanceof AppError) {
      return reply.status(error.statusCode).send({
        ok: false,
        error: { code: error.code, message: error.message },
        requestId: request.id
      })
    }

    if (error.statusCode === 429) {
      return reply.status(429).send({
        ok: false,
        error: { code: 'RATE_LIMITED', message: '请求过于频繁，请稍后再试' },
        requestId: request.id
      })
    }

    request.log.error({ err: error }, 'Unhandled request error')
    return reply.status(500).send({
      ok: false,
      error: { code: 'INTERNAL_ERROR', message: '服务暂时不可用' },
      requestId: request.id
    })
  })

  app.get('/health/live', {
    schema: {
      response: {
        200: successSchema({
          type: 'object',
          additionalProperties: false,
          required: ['status'],
          properties: { status: { const: 'alive' } }
        })
      }
    }
  }, async (request) => ok(request, { status: 'alive' }))

  app.get('/health/ready', async (request, reply) => {
    if (deps.readiness.isClosing()) {
      throw new AppError(503, 'NOT_READY', '服务正在关闭')
    }
    try {
      await deps.readiness.check()
      return ok(request, { status: 'ready' })
    } catch (error) {
      request.log.warn({ err: error }, 'Readiness dependency check failed')
      return reply.status(503).send({
        ok: false,
        error: { code: 'NOT_READY', message: '关键依赖暂不可用' },
        requestId: request.id
      })
    }
  })

  app.post<{ Body: LoginBody }>('/api/auth/login', {
    config: { rateLimit: { max: 5, timeWindow: '1 minute' } },
    schema: {
      body: {
        type: 'object',
        additionalProperties: false,
        required: ['email', 'password'],
        properties: {
          email: { type: 'string', minLength: 3, maxLength: 254, pattern: emailPattern },
          password: { type: 'string', minLength: 12, maxLength: 128 }
        }
      },
      response: {
        200: successSchema({
          type: 'object',
          additionalProperties: false,
          required: ['accessToken', 'expiresAt'],
          properties: {
            accessToken: { type: 'string' },
            expiresAt: { type: 'string' }
          }
        }),
        '4xx': failureSchema
      }
    }
  }, async (request) => ok(
    request,
    await deps.auth.login(request.body.email, request.body.password)
  ))

  app.post('/api/auth/logout', async (request) => {
    const actor = await authenticate(request)
    await deps.auth.logout(actor)
    return ok(request, null)
  })

  app.get('/api/me', async (request) => {
    const actor = await authenticate(request)
    return ok(request, actor)
  })

  app.get<{ Querystring: UserListQuery }>('/api/users', {
    schema: {
      querystring: {
        type: 'object',
        additionalProperties: false,
        properties: {
          q: { type: 'string', maxLength: 100, default: '' },
          page: { type: 'integer', minimum: 1, maximum: 100_000, default: 1 },
          pageSize: { type: 'integer', minimum: 1, maximum: 100, default: 20 }
        }
      }
    }
  }, async (request) => {
    const actor = await authenticate(request)
    const query = {
      q: request.query.q?.trim() ?? '',
      page: request.query.page ?? 1,
      pageSize: request.query.pageSize ?? 20
    }
    return ok(request, { ...await deps.users.list(actor, query), ...query })
  })

  app.get('/api/roles', async (request) => {
    const actor = await authenticate(request)
    return ok(request, await deps.users.listRoles(actor))
  })

  app.post<{ Body: CreateUserInput }>('/api/users', {
    schema: {
      body: {
        type: 'object',
        additionalProperties: false,
        required: ['email', 'displayName', 'password', 'roleIds'],
        properties: {
          email: { type: 'string', minLength: 3, maxLength: 254, pattern: emailPattern },
          displayName: { type: 'string', minLength: 2, maxLength: 80 },
          password: { type: 'string', minLength: 12, maxLength: 128 },
          roleIds: {
            type: 'array',
            minItems: 1,
            maxItems: 20,
            uniqueItems: true,
            items: { type: 'string', pattern: uuidPattern }
          }
        }
      }
    }
  }, async (request, reply) => {
    const actor = await authenticate(request)
    const user = await deps.users.create(actor, request.body, request.id)
    return reply.status(201).send(ok(request, user))
  })

  app.patch<{ Params: { id: string }; Body: StatusBody }>(
    '/api/users/:id/status',
    {
      schema: {
        params: {
          type: 'object',
          required: ['id'],
          properties: { id: { type: 'string', pattern: uuidPattern } }
        },
        body: {
          type: 'object',
          additionalProperties: false,
          required: ['status'],
          properties: { status: { type: 'string', enum: ['enabled', 'disabled'] } }
        }
      }
    },
    async (request) => {
      const actor = await authenticate(request)
      await deps.users.setStatus(actor, request.params.id, request.body.status, request.id)
      return ok(request, null)
    }
  )

  app.put<{ Params: { id: string }; Body: RoleBody }>(
    '/api/users/:id/roles',
    {
      schema: {
        params: {
          type: 'object',
          required: ['id'],
          properties: { id: { type: 'string', pattern: uuidPattern } }
        },
        body: {
          type: 'object',
          additionalProperties: false,
          required: ['roleIds'],
          properties: {
            roleIds: {
              type: 'array',
              minItems: 1,
              maxItems: 20,
              uniqueItems: true,
              items: { type: 'string', pattern: uuidPattern }
            }
          }
        }
      }
    },
    async (request) => {
      const actor = await authenticate(request)
      await deps.users.replaceRoles(actor, request.params.id, request.body.roleIds, request.id)
      return ok(request, null)
    }
  )

  return app
}
```

Fastify Schema 先拒绝多余字段、非法 UUID、过长字符串和错误枚举。需要数据库的校验仍在 Service 内执行。示例给登录和存活接口展示了响应 Schema；用户接口返回显式 `UserDTO`，生产项目还应继续为用户、角色和错误响应补齐 Schema，让序列化层形成第二道字段白名单。

`genReqId` 始终生成服务端 UUID，不直接信任客户端传入的 `x-request-id`。真实网关链路如果要透传 trace id，应校验格式、长度并区分 trace id 与本服务 request id。

## 阶段 9：迁移和种子数据

### 9.1 可校验的迁移脚本

`scripts/migrate.ts`：

```ts
import { createHash } from 'node:crypto'
import { readdir, readFile } from 'node:fs/promises'
import path from 'node:path'
import { fileURLToPath } from 'node:url'
import { createPool } from '../src/db.js'
import { loadConfig } from '../src/config.js'

const root = path.resolve(path.dirname(fileURLToPath(import.meta.url)), '..')
const sqlDirectory = path.join(root, 'sql')
const direction = process.argv[2] ?? 'up'

if (direction !== 'up' && direction !== 'down') {
  throw new Error('Usage: migrate.ts up|down')
}

const pool = createPool(loadConfig())
const client = await pool.connect()

try {
  await client.query(`
    CREATE TABLE IF NOT EXISTS schema_migrations (
      name text PRIMARY KEY,
      checksum char(64) NOT NULL,
      applied_at timestamptz NOT NULL DEFAULT now()
    )
  `)

  if (direction === 'up') {
    const files = (await readdir(sqlDirectory))
      .filter((file) => file.endsWith('.up.sql'))
      .sort()

    for (const file of files) {
      const sql = await readFile(path.join(sqlDirectory, file), 'utf8')
      const checksum = createHash('sha256').update(sql).digest('hex')
      const existing = await client.query<{ checksum: string }>(`
        SELECT checksum FROM schema_migrations WHERE name = $1
      `, [file])

      if (existing.rows[0]) {
        if (existing.rows[0].checksum !== checksum) {
          throw new Error(`Applied migration changed: ${file}`)
        }
        console.log(`skip ${file}`)
        continue
      }

      await client.query('BEGIN')
      try {
        await client.query(sql)
        await client.query(`
          INSERT INTO schema_migrations (name, checksum) VALUES ($1, $2)
        `, [file, checksum])
        await client.query('COMMIT')
        console.log(`up ${file}`)
      } catch (error) {
        await client.query('ROLLBACK')
        throw error
      }
    }
  } else {
    const latest = await client.query<{ name: string }>(`
      SELECT name FROM schema_migrations ORDER BY name DESC LIMIT 1
    `)
    const upFile = latest.rows[0]?.name
    if (!upFile) {
      console.log('nothing to rollback')
    } else {
      const downFile = upFile.replace(/\.up\.sql$/, '.down.sql')
      const sql = await readFile(path.join(sqlDirectory, downFile), 'utf8')
      await client.query('BEGIN')
      try {
        await client.query(sql)
        await client.query('DELETE FROM schema_migrations WHERE name = $1', [upFile])
        await client.query('COMMIT')
        console.log(`down ${downFile}`)
      } catch (error) {
        await client.query('ROLLBACK')
        throw error
      }
    }
  }
} finally {
  client.release()
  await pool.end()
}
```

已执行迁移的 checksum 发生变化时直接失败。不要修改已经进入共享环境的旧迁移；新建下一号迁移描述变更和回滚策略。

### 9.2 幂等种子脚本

`scripts/seed.ts`：

```ts
import { randomUUID } from 'node:crypto'
import { loadConfig } from '../src/config.js'
import { createPool, withTransaction } from '../src/db.js'
import { hashPassword } from '../src/password.js'

const email = process.env.ADMIN_EMAIL?.trim().toLowerCase()
const password = process.env.ADMIN_PASSWORD

if (!email || !password || password.length < 12) {
  throw new Error('ADMIN_EMAIL and ADMIN_PASSWORD (12+ chars) are required')
}

const pool = createPool(loadConfig())
const passwordHash = await hashPassword(password)

try {
  await withTransaction(pool, async (client) => {
    const roleResult = await client.query<{ id: string }>(`
      INSERT INTO roles (id, code, name)
      VALUES ($1, 'admin', '系统管理员')
      ON CONFLICT (code) DO UPDATE SET name = EXCLUDED.name
      RETURNING id
    `, [randomUUID()])
    const roleId = roleResult.rows[0]?.id
    if (!roleId) throw new Error('Unable to create admin role')

    let userResult = await client.query<{ id: string }>(`
      SELECT id FROM users WHERE lower(email) = lower($1)
    `, [email])

    let userId = userResult.rows[0]?.id
    if (userId) {
      await client.query(`
        UPDATE users
        SET password_hash = $2, status = 'enabled', updated_at = now()
        WHERE id = $1
      `, [userId, passwordHash])
    } else {
      userId = randomUUID()
      userResult = await client.query<{ id: string }>(`
        INSERT INTO users (id, email, display_name, password_hash)
        VALUES ($1, $2, '本地管理员', $3)
        RETURNING id
      `, [userId, email, passwordHash])
    }

    await client.query(`
      INSERT INTO user_roles (user_id, role_id)
      VALUES ($1, $2)
      ON CONFLICT DO NOTHING
    `, [userId, roleId])

    await client.query(`
      INSERT INTO role_permissions (role_id, permission_code)
      SELECT $1, code FROM permissions
      ON CONFLICT DO NOTHING
    `, [roleId])
  })
  console.log(`seeded admin ${email}`)
} finally {
  await pool.end()
}
```

脚本重复执行会更新本地管理员密码并补齐角色权限，不会创建多个管理员。生产环境不要保留固定默认密码或在应用每次启动时自动 seed。

## 阶段 10：进程装配与优雅停机

`src/server.ts`：

```ts
import process from 'node:process'
import { buildApp } from './app.js'
import { loadConfig } from './config.js'
import { createPool } from './db.js'
import {
  AuditRepository,
  AuthRepository,
  UserRepository
} from './repositories.js'
import { AuthService, UserService } from './services.js'

const config = loadConfig()
const pool = createPool(config)
const authRepository = new AuthRepository(pool)
const userRepository = new UserRepository(pool)
const authService = new AuthService(authRepository, config)
const userService = new UserService(
  pool,
  userRepository,
  authRepository,
  new AuditRepository()
)

let closing = false

const app = await buildApp({
  config,
  auth: authService,
  users: userService,
  readiness: {
    isClosing: () => closing,
    check: async () => {
      await pool.query('SELECT 1')
    }
  }
})

app.addHook('onClose', async () => {
  await pool.end()
})

await pool.query('SELECT 1')
await app.listen({ host: config.host, port: config.port })

async function shutdown(signal: string, exitCode = 0): Promise<void> {
  if (closing) return
  closing = true
  app.log.info({ signal }, 'Graceful shutdown started')

  const forceExit = setTimeout(() => {
    app.log.fatal('Graceful shutdown timed out')
    process.exit(1)
  }, 10_000)
  forceExit.unref()

  try {
    await app.close()
    clearTimeout(forceExit)
    process.exitCode = exitCode
  } catch (error) {
    app.log.fatal({ err: error }, 'Graceful shutdown failed')
    process.exitCode = 1
  }
}

process.once('SIGTERM', () => void shutdown('SIGTERM'))
process.once('SIGINT', () => void shutdown('SIGINT'))
process.on('uncaughtExceptionMonitor', (error, origin) => {
  app.log.fatal({ err: error, origin }, 'Uncaught exception')
})
process.on('unhandledRejection', (reason) => {
  app.log.fatal({ err: reason }, 'Unhandled promise rejection')
  void shutdown('unhandledRejection', 1)
})
```

关闭顺序是：先让 readiness 失败，再由 `app.close()` 停止接收新请求并等待在途请求，最后触发 `onClose` 释放数据库池。总超时防止发布永久卡住。

`uncaughtExceptionMonitor` 只记录证据，不阻止 Node 按默认策略退出。未捕获异常后继续提供服务是不安全的；可靠重启交给容器平台或外部进程管理器。

## 阶段 11：测试关键边界

### 11.1 密码测试

`tests/password.test.ts`：

```ts
import assert from 'node:assert/strict'
import { describe, it } from 'node:test'
import { hashPassword, verifyPassword } from '../src/password.js'

describe('password', () => {
  it('accepts the correct password and rejects a wrong password', async () => {
    const hash = await hashPassword('CorrectPassword123!')
    assert.equal(await verifyPassword('CorrectPassword123!', hash), true)
    assert.equal(await verifyPassword('WrongPassword123!', hash), false)
  })

  it('rejects an unsupported hash format', async () => {
    assert.equal(await verifyPassword('anything', 'plain-text'), false)
  })
})
```

### 11.2 API 契约测试

`tests/app.test.ts`：

```ts
import assert from 'node:assert/strict'
import { after, before, describe, it } from 'node:test'
import type { FastifyInstance } from 'fastify'
import { buildApp } from '../src/app.js'
import { AppError } from '../src/errors.js'
import { assertPermission } from '../src/services.js'
import type {
  AuthServicePort,
  UserServicePort
} from '../src/services.js'
import type { AuthContext } from '../src/types.js'

const admin: AuthContext = {
  sessionId: 'session-admin',
  userId: '00000000-0000-4000-8000-000000000001',
  email: 'admin@example.com',
  displayName: 'Admin',
  permissionCodes: ['user:read', 'user:create', 'user:update', 'role:grant']
}

const viewer: AuthContext = {
  ...admin,
  sessionId: 'session-viewer',
  userId: '00000000-0000-4000-8000-000000000002',
  email: 'viewer@example.com',
  displayName: 'Viewer',
  permissionCodes: ['user:read']
}

const auth: AuthServicePort = {
  async login(email, password) {
    if (email !== 'admin@example.com' || password !== 'CorrectPassword123!') {
      throw new AppError(401, 'INVALID_CREDENTIALS', '邮箱或密码不正确')
    }
    return { accessToken: 'admin-token-value-long-enough-123456', expiresAt: new Date().toISOString() }
  },
  async authenticate(token) {
    if (token === 'admin-token-value-long-enough-123456') return admin
    if (token === 'viewer-token-value-long-enough-12345') return viewer
    throw new AppError(401, 'UNAUTHORIZED', '登录已失效，请重新登录')
  },
  async logout() {}
}

const users: UserServicePort = {
  async list(actor, query) {
    assertPermission(actor, 'user:read')
    return { items: [], total: query.q ? 0 : 2 }
  },
  async listRoles(actor) {
    assertPermission(actor, 'user:read')
    return []
  },
  async create(actor) {
    assertPermission(actor, 'user:create')
    throw new Error('not needed by this test')
  },
  async setStatus(actor) {
    assertPermission(actor, 'user:update')
  },
  async replaceRoles(actor) {
    assertPermission(actor, 'role:grant')
  }
}

describe('permission API contract', () => {
  let app: FastifyInstance

  before(async () => {
    app = await buildApp({
      config: {
        nodeEnv: 'test',
        host: '127.0.0.1',
        port: 3000,
        databaseUrl: 'postgres://unused',
        dbPoolMax: 1,
        sessionTtlMinutes: 30
      },
      auth,
      users,
      readiness: { isClosing: () => false, check: async () => {} },
      logger: false
    })
  })

  after(async () => {
    await app.close()
  })

  it('returns 401 without a bearer token', async () => {
    const response = await app.inject({ method: 'GET', url: '/api/users' })
    assert.equal(response.statusCode, 401)
    assert.equal(response.json().error.code, 'UNAUTHORIZED')
  })

  it('allows a viewer to list users', async () => {
    const response = await app.inject({
      method: 'GET',
      url: '/api/users?page=1&pageSize=20',
      headers: { authorization: 'Bearer viewer-token-value-long-enough-12345' }
    })
    assert.equal(response.statusCode, 200)
    assert.equal(response.json().data.total, 2)
  })

  it('returns 403 when a viewer creates a user', async () => {
    const response = await app.inject({
      method: 'POST',
      url: '/api/users',
      headers: { authorization: 'Bearer viewer-token-value-long-enough-12345' },
      payload: {
        email: 'new@example.com',
        displayName: 'New User',
        password: 'NewPassword123!',
        roleIds: ['00000000-0000-4000-8000-000000000010']
      }
    })
    assert.equal(response.statusCode, 403)
    assert.equal(response.json().error.code, 'FORBIDDEN')
  })

  it('returns 400 for an invalid payload before the service runs', async () => {
    const response = await app.inject({
      method: 'POST',
      url: '/api/users',
      headers: { authorization: 'Bearer admin-token-value-long-enough-123456' },
      payload: { email: 'wrong' }
    })
    assert.equal(response.statusCode, 400)
    assert.equal(response.json().error.code, 'VALIDATION_ERROR')
  })

  it('returns a stable 404 response', async () => {
    const response = await app.inject({ method: 'GET', url: '/missing' })
    assert.equal(response.statusCode, 404)
    assert.equal(response.json().error.code, 'ROUTE_NOT_FOUND')
  })
})
```

测试通过依赖注入替换数据库 Service，不监听真实端口。数据库集成测试还应启动独立测试库，执行迁移，在每个测试后清理数据，并覆盖唯一约束、外键、事务回滚和并发角色替换；不要复用开发数据库。

执行：

```bash
npm run typecheck
npm test
npm run build
```

## 阶段 12：本地数据库和首次启动

### 12.1 Docker Compose

`compose.yaml`：

```yaml
services:
  postgres:
    image: postgres:18-alpine
    environment:
      POSTGRES_DB: permission_api
      POSTGRES_USER: app
      POSTGRES_PASSWORD: app
    ports:
      - "5432:5432"
    volumes:
      - postgres-data:/var/lib/postgresql
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U app -d permission_api"]
      interval: 5s
      timeout: 3s
      retries: 20

  api:
    build: .
    environment:
      NODE_ENV: production
      HOST: 0.0.0.0
      PORT: 3000
      DATABASE_URL: postgres://app:app@postgres:5432/permission_api
      DB_POOL_MAX: 10
      SESSION_TTL_MINUTES: 30
    ports:
      - "3000:3000"
    depends_on:
      postgres:
        condition: service_healthy

volumes:
  postgres-data:
```

启动并初始化：

```bash
cp .env.example .env
docker compose up -d postgres
npm run db:migrate
npm run db:seed
npm run dev
```

### 12.2 手工调用

健康检查：

```bash
curl -i http://127.0.0.1:3000/health/live
curl -i http://127.0.0.1:3000/health/ready
```

登录并保存 token：

```bash
curl -s http://127.0.0.1:3000/api/auth/login \
  -H 'Content-Type: application/json' \
  --data '{"email":"admin@example.com","password":"ChangeThisAdmin123!"}'
```

将响应中的 `accessToken` 只放在当前终端变量，不写入仓库：

```bash
export ACCESS_TOKEN='替换为登录响应中的值'
curl -s http://127.0.0.1:3000/api/me \
  -H "Authorization: Bearer $ACCESS_TOKEN"
curl -s 'http://127.0.0.1:3000/api/users?page=1&pageSize=20' \
  -H "Authorization: Bearer $ACCESS_TOKEN"
```

主动退出后，同一个 token 应立即得到 401：

```bash
curl -s -X POST http://127.0.0.1:3000/api/auth/logout \
  -H "Authorization: Bearer $ACCESS_TOKEN"
curl -i http://127.0.0.1:3000/api/me \
  -H "Authorization: Bearer $ACCESS_TOKEN"
```

## 阶段 13：构建和容器化

### 13.1 生产镜像

`Dockerfile`：

```dockerfile
FROM node:24-bookworm-slim AS build
WORKDIR /app
COPY package.json package-lock.json ./
RUN npm ci
COPY tsconfig.json ./
COPY src ./src
COPY scripts ./scripts
RUN npm run build

FROM node:24-bookworm-slim AS runtime
ENV NODE_ENV=production
WORKDIR /app
COPY package.json package-lock.json ./
RUN npm ci --omit=dev && npm cache clean --force
COPY --from=build /app/dist ./dist
COPY sql ./dist/sql
USER node
EXPOSE 3000
CMD ["node", "dist/src/server.js"]
```

`.dockerignore`：

```text
node_modules
dist
.env
.git
coverage
npm-debug.log
```

构建镜像、迁移、初始化并运行完整服务：

```bash
docker compose build api
docker compose run --rm api node dist/scripts/migrate.js up
docker compose run --rm \
  -e ADMIN_EMAIL=admin@example.com \
  -e ADMIN_PASSWORD='ChangeThisAdmin123!' \
  api node dist/scripts/seed.js
docker compose up -d api
curl -i http://127.0.0.1:3000/health/ready
```

Compose 中的数据库主机名必须是服务名 `postgres`，不能写 `127.0.0.1`，因为容器内的回环地址只指向 API 容器自身。API 监听 `0.0.0.0`，否则宿主机端口映射无法访问。镜像不复制 `.env`，运行时通过 Compose 或部署平台注入配置。

### 13.2 发布顺序

```mermaid
flowchart LR
  A["构建不可变镜像"] --> B["备份与迁移兼容检查"]
  B --> C["运行 dist/scripts/migrate.js up"]
  C --> D["发布新实例"]
  D --> E["readiness 通过"]
  E --> F["逐步切流"]
  F --> G["观察错误率与延迟"]
  G --> H["完成或回滚应用"]
```

迁移必须与新旧应用版本兼容。推荐顺序是：先增加可空字段或新表，再发布同时兼容新旧结构的代码，完成回填后再收紧约束。不要在同一次发布中先删除旧列再期望旧实例正常运行。

生产迁移任务可使用同一镜像执行：

```bash
node dist/scripts/migrate.js up
```

应用回滚不等于数据库 Down。先回滚到仍兼容新结构的旧镜像；只有迁移确实可逆、数据已备份且没有新版本写入不可恢复数据时，才执行受控 Down。

### 13.3 反向代理前缀

如果应用路由本身包含 `/api`，Nginx `proxy_pass` 不要意外剥掉该前缀：

```nginx
location /api/ {
  proxy_pass http://node_api:3000;
  proxy_http_version 1.1;
  proxy_set_header Host $host;
  proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
  proxy_set_header X-Forwarded-Proto $scheme;
}
```

这里 `proxy_pass` 后没有额外 URI 部分，因此 `/api/users` 仍转发为 `/api/users`。在启用 `trust proxy` 或根据转发头判断客户端 IP、HTTPS 和 Secure Cookie 前，必须只信任受控代理；客户端可自行伪造未经网关清洗的头。

## 阶段 14：故障注入

不要只验证管理员成功登录。按表主动制造失败：

| 编号 | 注入方式 | 预期行为 | 关键证据 |
| --- | --- | --- | --- |
| N01 | 删除 `DATABASE_URL` | 进程启动失败，不监听端口 | 启动日志、退出码 |
| N02 | 请求体多一个 `isAdmin` | 400 `VALIDATION_ERROR` | 响应 + Schema |
| N03 | 不带 Bearer token | 401 `UNAUTHORIZED` | 状态码、request id |
| N04 | 只读用户调用创建接口 | 403 `FORBIDDEN` | 后端 Service 权限判断 |
| N05 | 并发创建同一邮箱 | 只有一个成功，另一个 409 | 唯一索引 `23505` |
| N06 | 角色 ID 不存在 | 400，事务无用户和审计残留 | 数据库查询 |
| N07 | 审计写入故意失败 | 整个写操作回滚 | 事务日志、表数据 |
| N08 | 并发替换同一用户角色 | 用户行锁串行化，不产生半状态 | 锁等待和最终角色 |
| N09 | 停用目标用户 | 目标全部活动会话立即失效 | sessions + 后续 401 |
| N10 | 连续错误登录 6 次 | 第 6 次 429 | 限流响应和日志 |
| N11 | 停止 PostgreSQL | liveness 200，readiness 503 | 两个健康端点 |
| N12 | 修改已应用 migration | 迁移因 checksum 不一致失败 | 迁移日志 |
| N13 | 发送 `SIGTERM` | readiness 失败、在途请求排空、池关闭 | 发布日志 |
| N14 | 查询参数 `page=0` | 400，不进入 Repository | Fastify validation |
| N15 | 制造未知异常 | 客户端只见 500，日志含完整 err 与 request id | 响应/结构化日志 |

### 14.1 验证数据库约束

不要只从 UI 点击。可以并发发送两个同邮箱创建请求，然后查询：

```sql
SELECT lower(email), count(*)
FROM users
GROUP BY lower(email)
HAVING count(*) > 1;
```

结果必须为空。即使 Service 先查询邮箱，最终正确性仍由唯一索引保证。

### 14.2 验证角色替换回滚

在测试数据库临时给 `audit_logs.action` 增加不可能满足的约束，使审计插入失败，再调用角色替换。事务结束后 `user_roles` 必须仍是旧值。验证完成后还原测试数据库，不要在共享环境修改约束。

### 14.3 验证优雅停机

1. 创建一个测试路由或数据库查询，让请求稳定执行 3 秒。
2. 请求开始后向进程发送 `SIGTERM`。
3. 立刻检查 readiness 应变为 503。
4. 原请求应在关闭期限内完成，新请求不再进入。
5. 日志应出现开始关闭和资源释放信息，进程退出码为 0。

## 阶段 15：测试与交付矩阵

| 层级 | 主要覆盖 | 不应该依赖 |
| --- | --- | --- |
| 纯函数单元测试 | 密码格式、权限判断、配置解析 | 网络和共享数据库 |
| Service 测试 | 自我停用、角色校验、事务调用 | Fastify 请求对象 |
| API inject 测试 | Schema、状态码、认证、统一响应 | 真实监听端口 |
| 数据库集成测试 | SQL、约束、迁移、回滚、锁 | 开发数据库现有数据 |
| 容器冒烟测试 | 镜像、配置、监听地址、信号 | 本地 node_modules |
| 故障注入 | 依赖断连、限流、并发、关闭 | 只跑成功路径 |

CI 最小顺序：

```text
npm ci
  -> npm run typecheck
  -> npm test
  -> 启动临时 PostgreSQL
  -> npm run db:migrate
  -> 数据库集成测试
  -> npm run build
  -> docker build
  -> 镜像健康检查
```

测试数据库使用独立数据库名和凭据。每次 CI 从空数据库执行迁移；测试结束销毁实例，而不是依赖执行顺序或共享脏数据。

## 完成标准

### 功能

- [ ] 管理员可以登录、查询用户、创建用户、停用用户和替换角色。
- [ ] 退出、过期和停用用户都会让原会话返回 401。
- [ ] 无动作权限返回 403，不被误报为 401。
- [ ] 搜索和分页有边界校验与稳定排序。
- [ ] 同邮箱并发创建由数据库唯一索引守住。

### 数据与安全

- [ ] 密码和原始会话 token 都没有以明文写入数据库或日志。
- [ ] 所有用户输入 SQL 都使用参数占位符。
- [ ] 用户角色替换、会话撤销和审计在同一事务边界内完成。
- [ ] DDL 包含字段含义、唯一性、外键、删除策略、索引目的和回滚说明。
- [ ] 登录限流、请求体限制和错误脱敏已验证。

### 质量与交付

- [ ] `npm ci`、typecheck、test 和 build 从干净目录通过。
- [ ] 迁移能在空库 up，受控 down 后还能再次 up。
- [ ] liveness 与 readiness 在数据库断连时表现不同。
- [ ] SIGTERM 能排空请求并关闭数据库池。
- [ ] 390px 前端调用、curl 和生产反向代理都使用同一 API 前缀。
- [ ] README 记录环境变量、迁移、seed、启动、测试、构建、部署和回滚步骤。

## 项目交付材料

```text
README.md                 安装、启动、命令和本地账号
ARCHITECTURE.md           请求链路、分层和依赖方向
API_CONTRACT.md           路由、Schema、状态码和权限码
DATABASE.md               表、约束、索引、迁移和回滚
SECURITY.md               会话、密码、限流、日志脱敏和威胁边界
TEST_NOTES.md             自动测试与故障注入证据
DEPLOYMENT.md             镜像、健康检查、迁移发布和回滚
INCIDENT_NOTES.md         至少三次故障复盘
```

## 实际项目常见问题

### 角色修改后权限仍是旧值

先确认权限事实来自数据库、进程内缓存还是 Redis。本项目每次认证重新查询，因此下一条请求立即生效；引入缓存后必须在事务提交后删除受影响用户的权限 key，并验证多实例都能看到失效。

### 接口偶发等待数秒

按 request id 拆分：事件循环延迟、连接池等待、SQL、锁等待、密码哈希和外部依赖。不要只增加池大小；总连接数等于实例数乘每实例池大小，盲目增大可能把等待转移到数据库。

### 登录失败被当成数据库故障

只有凭据错误、会话不存在、过期、撤销或用户停用属于 401。数据库超时和断连属于 500/503，应记录内部错误并让 readiness 反映依赖状态，不能用大范围 `catch` 把所有异常吞成“请重新登录”。

### 返回了 500 而不是 400

Fastify Schema 错误也会进入自定义 error handler。必须先判断 `error.validation`，再处理 `AppError`、429 和未知错误。API 测试要对非法 body 断言 400，避免错误处理器把框架错误全部覆盖成 500。

### 生产构建能通过但启动找不到模块

NodeNext 相对导入需要 `.js` 后缀；`rootDir` 和 `outDir` 会影响最终入口。本章构建结果是 `dist/src/server.js`，Docker 和 `start` 脚本必须指向同一路径。更多运行时问题进入 [Node.js 真实项目问题库](/projects/issues-node)。

## 生产化时仍要补什么

| 教学项目 | 生产系统通常还需要 |
| --- | --- |
| 单数据库认证查询 | 会话缓存、失效广播和数据库降级策略 |
| Bearer token | 浏览器端安全存储决策、CSRF/XSS 威胁模型或 HttpOnly Cookie |
| 单实例内存限流 | Redis/网关分布式限流、账户维度与 IP 维度组合 |
| 角色权限 | 租户、部门、数据范围、超级管理员保护和审批 |
| 基础日志 | trace、指标、告警、日志保留和隐私合规 |
| 单迁移 | expand/contract、在线索引、回填、备份恢复演练 |
| 单区域数据库 | 高可用、连接代理、灾备和容量规划 |

不要把本章示例直接暴露到公网。先完成威胁建模、依赖审计、数据库备份恢复、密钥管理和实际负载测试。

## 参考资料

- [Node.js Releases](https://nodejs.org/en/about/previous-releases)
- [Node.js TypeScript](https://nodejs.org/api/typescript.html)
- [Node.js Test Runner](https://nodejs.org/api/test.html)
- [Node.js Crypto](https://nodejs.org/api/crypto.html)
- [Node.js Process](https://nodejs.org/api/process.html)
- [Fastify Getting Started](https://fastify.dev/docs/latest/Guides/Getting-Started/)
- [Fastify Lifecycle](https://fastify.dev/docs/latest/Reference/Lifecycle/)
- [Fastify Validation and Serialization](https://fastify.dev/docs/latest/Reference/Validation-and-Serialization/)
- [Fastify Testing](https://fastify.dev/docs/latest/Guides/Testing/)
- [node-postgres Pooling](https://node-postgres.com/features/pooling)
- [node-postgres Transactions](https://node-postgres.com/features/transactions)
- [PostgreSQL Constraints](https://www.postgresql.org/docs/current/ddl-constraints.html)
- [PostgreSQL Explicit Locking](https://www.postgresql.org/docs/current/explicit-locking.html)

## 下一步学习

完成项目后进入 [Node.js 专项练习](/roadmap/node-practice)，注入事件循环阻塞、Stream 背压、open handles、SIGTERM 和多实例问题。通用接口与事务问题查看 [后端接口与服务问题](/projects/issues-backend)，Node 运行时问题查看 [Node.js 真实项目问题库](/projects/issues-node)。
