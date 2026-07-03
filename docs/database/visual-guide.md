# 图解数据库核心概念

## 这个页面解决什么

数据库文档如果只讲 SQL 语法，很难解决真实项目问题。实际开发更需要理解：表怎么设计、索引怎么工作、事务为什么会锁、慢查询怎么定位、缓存和数据库如何配合。

## 适合谁看

适合已经能写增删改查，但对表关系、索引、事务锁、慢查询、缓存一致性和数据库排错还缺少整体模型的人。

## 一张图理解后端如何访问数据库

```mermaid
flowchart TD
  A["Controller / Handler"] --> B["Service<br/>业务和事务边界"]
  B --> C["Repository / Mapper<br/>数据访问"]
  C --> D["连接池"]
  D --> E[("Database")]
  B --> F["Cache<br/>Redis"]
  C --> G["Migration<br/>表结构版本"]
```

关键理解：

- Service 决定业务动作和事务边界。
- Repository/Mapper 负责 SQL 和数据映射。
- 连接池负责复用数据库连接。
- 迁移脚本负责表结构变更可追踪。

## 一张图理解表设计

```mermaid
erDiagram
  departments ||--o{ employees : contains
  roles ||--o{ user_roles : grants
  users ||--o{ user_roles : owns

  departments {
    bigint id PK
    string name
    bigint parent_id
  }

  employees {
    bigint id PK
    bigint department_id FK
    string name
    string status
  }

  users {
    bigint id PK
    string username
    string password_hash
  }

  roles {
    bigint id PK
    string code
    string name
  }

  user_roles {
    bigint user_id FK
    bigint role_id FK
  }
```

设计表时要先回答：

- 这个表代表什么业务对象。
- 主键是什么。
- 哪些字段必须唯一。
- 哪些字段可以为空。
- 和其他表是什么关系。
- 查询最常按什么条件过滤。

## 一张图理解 B+Tree 索引查找

```mermaid
flowchart TD
  A["Root 根节点"] --> B["Internal 节点<br/>10 / 30 / 60"]
  B --> C["Leaf<br/>1-9"]
  B --> D["Leaf<br/>10-29"]
  B --> E["Leaf<br/>30-59"]
  B --> F["Leaf<br/>60+"]
  D --> G["找到 key=18 对应行位置"]
```

索引不是“越多越好”。它能加快查询，但会增加写入成本和存储成本。

适合建索引：

- 高频查询条件。
- 高频排序字段。
- join 字段。
- 唯一约束字段。

不适合盲目建索引：

- 区分度很低的字段。
- 很少查询的字段。
- 频繁更新的大字段。

## 一张图理解联合索引最左前缀

```mermaid
flowchart LR
  A["联合索引<br/>(tenant_id, status, created_at)"] --> B["tenant_id 可用"]
  B --> C["tenant_id + status 可用"]
  C --> D["tenant_id + status + created_at 可用"]
  A --> E["跳过 tenant_id 直接查 status<br/>通常用不好这个索引"]
```

如果索引是 `(tenant_id, status, created_at)`，查询最好从 `tenant_id` 开始。跳过最左列，数据库很可能无法充分利用索引。

## 一张图理解事务隔离和锁

```mermaid
sequenceDiagram
  participant T1 as Transaction A
  participant DB as Database Row
  participant T2 as Transaction B

  T1->>DB: update stock set count=count-1 where id=1
  DB-->>T1: 加行锁
  T2->>DB: update same row
  DB-->>T2: 等待锁释放
  T1->>DB: commit
  DB-->>T2: 获得锁继续执行
```

事务越长，锁持有越久。不要在事务中做：

- 外部 HTTP 调用。
- 大文件处理。
- 人工确认等待。
- 大批量无分页处理。

## 一张图理解慢查询排查

```mermaid
flowchart TD
  A["接口慢"] --> B["看接口耗时拆分"]
  B --> C{"SQL 慢吗"}
  C -- "否" --> D["查外部接口、缓存、锁、网络"]
  C -- "是" --> E["看慢 SQL 日志"]
  E --> F["EXPLAIN 执行计划"]
  F --> G{"是否走索引"}
  G -- "否" --> H["补索引或改 SQL"]
  G -- "是" --> I["检查扫描行数、排序、回表、分页深度"]
```

慢查询不要靠猜。至少要看：

- SQL 文本。
- 参数。
- 执行耗时。
- 扫描行数。
- 执行计划。
- 表数据量。

## 一张图理解缓存和数据库

```mermaid
flowchart TD
  A["请求读取数据"] --> B{"Redis 有缓存吗"}
  B -- "有" --> C["返回缓存"]
  B -- "没有" --> D["查数据库"]
  D --> E["写入缓存"]
  E --> F["返回数据"]
  G["数据更新"] --> H["更新数据库"]
  H --> I["删除或更新缓存"]
```

缓存不能替代数据库一致性。常见策略：

- 读多写少：查询缓存，更新后删除缓存。
- 热点数据：设置过期时间和防击穿。
- 强一致场景：谨慎使用缓存，或引入版本号。

## 一张图理解数据库问题定位

```mermaid
flowchart TD
  A["数据库问题"] --> B{"现象"}
  B --> C["查询慢"]
  B --> D["连接池耗尽"]
  B --> E["死锁 / 锁等待"]
  B --> F["数据不一致"]
  B --> G["迁移失败"]

  C --> C1["慢 SQL + EXPLAIN + 索引"]
  D --> D1["连接池配置、慢 SQL、连接未释放"]
  E --> E1["事务顺序、锁范围、事务时长"]
  F --> F1["事务边界、并发更新、缓存同步"]
  G --> G1["回滚脚本、字段默认值、历史数据"]
```

## 下一步学习

继续学习 [数据建模与表设计](/database/modeling)，或进入 [索引与查询优化](/database/indexes)。
