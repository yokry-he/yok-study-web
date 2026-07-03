# ORM 实战

## 适合谁看

适合已经理解 SQL 基础，但在项目里开始使用 Prisma、TypeORM、Drizzle、Sequelize 或其他 ORM 的学习者：

- 不知道 ORM 和 SQL 应该怎么配合。
- 以为用了 ORM 就不需要懂索引和事务。
- 查询接口很方便，但线上 SQL 很慢。
- 关联查询经常查出太多字段。
- 事务写法看起来对，实际却没有包住全部操作。

ORM 的价值是提高工程效率、统一模型、减少重复 SQL，但它不能替代数据库设计。真正可靠的项目会同时关注模型、查询、事务、迁移、日志和性能。

## ORM 解决什么

ORM 把数据库表映射成代码里的模型。

它通常提供：

| 能力 | 说明 |
| --- | --- |
| 模型定义 | 用 schema 或实体类描述表、字段、关系 |
| 类型提示 | 查询结果能被 TypeScript 推导 |
| CRUD API | 用函数调用代替手写常规 SQL |
| 关联查询 | 查询用户时带出角色、部门、订单等数据 |
| 迁移工具 | 根据模型变化生成数据库迁移 |
| 事务 API | 把多个数据库操作放到同一事务里 |

但 ORM 不会自动解决：

- 错误的数据模型。
- 缺失的索引。
- 不合理的关联查询。
- 错误的事务边界。
- 线上迁移风险。
- 数据权限和审计。

## 推荐分层

不要让 controller 直接调用 ORM。

```text
controller
↓
service
↓
repository
↓
orm client
↓
database
```

职责：

| 层 | 负责 |
| --- | --- |
| controller | HTTP 参数、响应状态、错误转换 |
| service | 业务规则、权限、事务编排 |
| repository | ORM 查询、字段选择、数据映射 |
| orm client | 连接、查询 API、事务 API |
| database | 数据约束、索引、事务隔离 |

这样做的好处是：以后从 Prisma 换成 Drizzle，或者把某个复杂查询改成原生 SQL，不会影响 controller 和大部分业务代码。

## 模型设计不要只跟着页面走

错误做法：

```text
页面上有 name、phone、roleName
↓
数据库就建 name、phone、roleName
```

更稳定的做法：

```text
user
├─ id
├─ username
├─ phone
├─ status

role
├─ id
├─ code
├─ name

user_role
├─ user_id
└─ role_id
```

页面展示字段可以组合出来，但数据库模型要表达业务关系。

## 查询只取需要的字段

很多 ORM 默认会返回整行数据。后台列表页通常不需要密码哈希、备注、大文本、内部配置等字段。

示例：

```ts
const users = await prisma.user.findMany({
  select: {
    id: true,
    username: true,
    displayName: true,
    status: true,
    createdAt: true
  }
})
```

这样做有三个好处：

- 减少网络传输。
- 降低敏感字段误返回的风险。
- 让接口字段更稳定。

如果确实需要关联数据，也要明确选择字段：

```ts
const users = await prisma.user.findMany({
  select: {
    id: true,
    username: true,
    roles: {
      select: {
        role: {
          select: {
            code: true,
            name: true
          }
        }
      }
    }
  }
})
```

不要为了方便直接把深层关联全部 include 出来。

## 关联查询要警惕 N+1

典型问题：

```ts
const users = await userRepository.findUsers()

for (const user of users) {
  user.roles = await roleRepository.findRolesByUserId(user.id)
}
```

如果有 100 个用户，就可能产生 101 次查询。

常见解决方案：

1. 使用 ORM 的关联查询能力。
2. 用 `where in` 一次查出所有关联数据。
3. 对复杂列表写专门的查询方法。
4. 用日志记录 SQL 数量和耗时。

列表页最容易被 N+1 拖慢。新增列表接口时，一定要观察实际 SQL。

## 事务必须由 service 编排

事务代表一个业务动作的原子性。

例如创建用户：

```text
创建用户
↓
绑定角色
↓
写审计日志
```

这些操作应该同成功、同失败。

示例：

```ts
await prisma.$transaction(async (tx) => {
  const user = await tx.user.create({
    data: {
      username,
      displayName
    }
  })

  await tx.userRole.createMany({
    data: roleIds.map((roleId) => ({
      userId: user.id,
      roleId
    }))
  })

  await tx.auditLog.create({
    data: {
      action: 'user.create',
      targetId: user.id,
      operatorId: currentUser.id
    }
  })
})
```

注意：事务里的查询要使用事务对象 `tx`，不要混用全局 client。

## 事务里不要做慢操作

不要在事务里做这些事情：

- 调用第三方 HTTP 接口。
- 上传文件。
- 发送短信或邮件。
- 执行很慢的统计查询。
- 等待用户交互。

事务持有连接和锁。事务越长，越容易导致锁等待、死锁和连接池耗尽。

推荐流程：

```text
参数校验
↓
外部前置校验
↓
开启短事务
↓
写数据库
↓
提交
↓
提交后发送异步通知
```

## 什么时候用原生 SQL

ORM 适合常规 CRUD，但不是所有查询都适合 ORM。

可以考虑原生 SQL 的场景：

- 复杂报表。
- 多表聚合。
- 窗口函数。
- 性能要求很高的列表。
- 数据库特有能力，例如 PostgreSQL JSONB、全文搜索、GIS。

但原生 SQL 也要遵守：

- 使用参数绑定，禁止字符串拼接。
- 写清楚索引依赖。
- 有 EXPLAIN 证据。
- 返回字段明确。
- 放在 repository 层。

## migration 不等于模型自动同步

开发环境可以快速生成迁移，生产环境必须审查迁移。

迁移说明至少写清楚：

- 为什么新增或修改字段。
- 字段是否允许为空。
- 默认值怎么来的。
- 对旧数据有什么影响。
- 是否需要回填。
- 是否会锁大表。
- 如何验证和回滚。

不要在生产环境使用“自动同步表结构”的模式。它可能在你没意识到时修改或删除结构。

## 实际项目问题

### 1. ORM 查询很慢

**现象**

接口代码只有一行 ORM 查询，但页面加载很慢。

**原因**

- include 了太多关联。
- 没有索引。
- 出现 N+1 查询。
- 分页字段不稳定。
- ORM 生成的 SQL 没被查看过。

**解决方案**

1. 打开 SQL 日志。
2. 复制慢 SQL 到数据库执行 `EXPLAIN`。
3. 减少 select 字段。
4. 给 WHERE、JOIN、ORDER BY 对应字段补索引。
5. 必要时改成手写 SQL。

### 2. 删除数据失败

**现象**

删除用户时报外键约束错误。

**原因**

用户被角色、订单、审计日志等记录引用。

**解决方案**

- 核心业务优先软删除。
- 明确哪些关系允许级联删除。
- 不要随意对重要数据开启 cascade delete。
- 删除前给前端返回可理解的业务提示。

### 3. 字段类型和 TypeScript 类型对不上

**现象**

数据库里是 decimal，代码里当 number 直接计算，金额出现精度问题。

**解决方案**

- 金额优先用整数分保存，或使用 decimal 类型并配合 Decimal 库。
- 在 repository 层做数据映射。
- 不要让前端和后端各自猜字段语义。

### 4. 事务没有回滚

**现象**

用户创建失败，但角色绑定已经写入。

**原因**

部分操作没有使用事务对象，而是混用了全局 ORM client。

**解决方案**

- 事务内统一使用 `tx`。
- 把事务编排放到 service。
- repository 支持接收事务 client。
- 给失败路径写集成测试。

## 最佳实践

- ORM 负责提高效率，数据库约束负责兜底。
- controller 不直接调用 ORM。
- 查询默认只选需要字段。
- 列表接口必须关注 N+1。
- 事务由 service 编排，且尽量短。
- 复杂查询允许使用原生 SQL，但必须参数化。
- 迁移进入代码评审，不能无审查自动同步生产库。
- 权限、金额、订单、审计等核心数据要有测试覆盖。

## 参考资料

- [Prisma relation queries](https://www.prisma.io/docs/orm/prisma-client/queries/relation-queries)
- [Prisma transactions and batch queries](https://www.prisma.io/docs/orm/prisma-client/queries/transactions)

## 下一步学习

继续学习 [备份与恢复](/database/backup-recovery)，理解数据出问题后如何恢复。
