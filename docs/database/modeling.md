# 数据建模与表设计

## 适合谁看

适合准备从“接口字段”走向“业务数据模型”的学习者。数据库建模解决的是：业务对象如何稳定、清晰、可约束地保存下来。

建模不是画几张表。它需要回答：

- 业务对象是什么。
- 对象之间是什么关系。
- 哪些字段是核心事实。
- 哪些字段只是展示缓存。
- 哪些值必须唯一。
- 哪些关系必须存在。
- 哪些数据允许删除。

## 从业务对象开始

以后台权限系统为例：

| 业务对象 | 含义 |
| --- | --- |
| 用户 | 登录系统的人 |
| 角色 | 一组权限的集合 |
| 权限 | 可访问菜单、路由、按钮或接口的能力 |
| 部门 | 组织结构 |
| 用户角色 | 用户和角色的多对多关系 |
| 角色权限 | 角色和权限的多对多关系 |

先列业务对象，再设计表。不要从页面表单字段直接反推数据库结构。

## 关系类型

| 关系 | 示例 | 建模方式 |
| --- | --- | --- |
| 一对一 | 用户和用户资料 | 共享主键或唯一外键 |
| 一对多 | 部门和用户 | 多的一方保存外键 |
| 多对多 | 用户和角色 | 中间表 |
| 树结构 | 部门上下级 | parent_id、闭包表或路径 |

## 示例：角色权限模型

```sql
create table role (
  id bigint generated always as identity primary key,

  -- 角色编码用于权限判断和审计，要求稳定且唯一。
  code text not null,

  -- 角色名称用于后台展示，可修改，但不应影响权限判断。
  name text not null,

  -- 状态控制角色是否可被授权或生效。
  status text not null default 'active'
    check (status in ('active', 'disabled')),

  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),

  constraint role_code_unique unique (code)
);

comment on table role is '角色表：定义一组权限的业务身份，如管理员、运营、财务';
comment on column role.code is '角色编码，权限判断和审计使用，创建后应保持稳定';
comment on column role.status is '角色状态：active=可用，disabled=停用；停用后不应继续授权';

create table permission (
  id bigint generated always as identity primary key,
  code text not null,
  name text not null,
  type text not null check (type in ('menu', 'route', 'button', 'api')),
  created_at timestamptz not null default now(),

  constraint permission_code_unique unique (code)
);

comment on table permission is '权限表：保存菜单、路由、按钮和接口等可授权资源';
comment on column permission.type is '权限类型：menu=菜单，route=路由，button=按钮，api=接口';

create table role_permission (
  role_id bigint not null references role(id) on delete cascade,
  permission_id bigint not null references permission(id) on delete cascade,
  granted_at timestamptz not null default now(),

  -- 复合主键保证同一个角色不能重复绑定同一个权限。
  primary key (role_id, permission_id)
);

comment on table role_permission is '角色权限关联表：保存角色和权限的多对多授权关系';
comment on column role_permission.granted_at is '授权时间，用于审计角色权限变更';

-- 外键列用于反查某个权限被哪些角色使用，需补充索引。
create index role_permission_permission_id_idx
on role_permission (permission_id);
```

## 字段设计原则

### 字段名表达业务

推荐：

```text
created_at
paid_at
deleted_at
approved_by
```

不推荐：

```text
time1
flag
type2
data
```

模糊字段会让后续维护者不知道怎么用。

### 空值要有含义

字段允许 NULL 前要说明含义。

例如：

```sql
deleted_at timestamptz null
```

含义：为空表示未删除，有值表示软删除时间。

如果没有明确含义，就不要随便允许 NULL。

### 状态字段要有限定

状态字段应该有枚举说明和约束。

```sql
status text not null default 'pending'
  check (status in ('pending', 'paid', 'cancelled', 'refunded'))
```

约束的业务原因：订单只能处于这四类状态，避免后端 bug 或脚本写入未知状态。

## 软删除还是物理删除

适合软删除：

- 用户。
- 订单。
- 权限。
- 财务记录。
- 审计相关数据。

可以物理删除：

- 临时草稿。
- 过期验证码。
- 可重建缓存。
- 无业务价值的临时任务记录。

软删除字段：

```sql
deleted_at timestamptz null
```

不要只用 `is_deleted boolean`。删除时间更有审计价值。

## 实际项目问题

### 问题：权限码散落在前端和后端，改名后大面积出错

**原因**

权限码没有作为稳定业务标识维护。

**解决方案**

权限码进入数据库和代码常量，并规定创建后不随意修改。展示名称可以改，code 要谨慎迁移。

### 问题：一个字段被塞进各种 JSON 配置

**原因**

为了开发快，把核心字段放进 `settings`。

**风险**

- 无法加约束。
- 查询困难。
- 索引复杂。
- 文档不清楚。

**解决方案**

如果字段进入核心查询、权限判断或报表统计，应迁移成正式列。

## 最佳实践

- 先建业务对象模型，再建表。
- 每个表写 table comment，每个非显然字段写 column comment。
- 唯一业务规则用唯一约束。
- 多对多关系用中间表，不要用逗号字符串保存 ID。
- 核心状态字段加 check 约束或枚举治理。
- 软删除使用 `deleted_at`，并说明查询默认排除规则。

## 下一步学习

继续学习 [索引与查询优化](/database/indexes)。
