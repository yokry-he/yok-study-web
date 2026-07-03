# 图解 Java 核心概念

## 这个页面解决什么

如果你刚开始学 Java，很容易被 JDK、JVM、类、对象、堆、栈、线程、事务、Spring Bean 这些概念分散注意力。

这一页先用图把 Java 后端项目的核心模型串起来。读完后再进入具体章节，会更容易理解每个知识点在项目里的位置。

## 一张图理解 Java 程序从哪里来到哪里去

```mermaid
flowchart LR
  A["User.java<br/>源码文件"] --> B["javac<br/>编译器"]
  B --> C["User.class<br/>字节码"]
  C --> D["ClassLoader<br/>类加载器"]
  D --> E["JVM Runtime<br/>运行时"]
  E --> F["Interpreter<br/>解释执行"]
  E --> G["JIT Compiler<br/>热点代码编译"]
  F --> H["CPU 执行"]
  G --> H
```

这张图要记住三件事：

1. Java 代码不是直接运行源码，而是先编译成 `.class` 字节码。
2. JVM 负责加载字节码、管理内存、执行代码和做 GC。
3. 热点代码会被 JIT 编译优化，所以 Java 不是简单的“解释型语言”。

实际项目中，`java --version`、`javac --version`、Maven 编译版本和服务器 JDK 版本必须一致或兼容，否则就会出现本地能跑、服务器不能跑的问题。

## 一张图理解 JDK、JRE、JVM、Maven 的关系

```mermaid
flowchart TD
  A["JDK<br/>开发工具包"] --> B["javac<br/>编译工具"]
  A --> C["jar / javadoc / jdeps<br/>开发诊断工具"]
  A --> D["JRE<br/>运行环境"]
  D --> E["JVM<br/>执行字节码"]
  D --> F["标准类库<br/>java.util / java.time / java.net"]
  G["Maven / Gradle"] --> H["下载依赖"]
  G --> I["执行测试"]
  G --> J["打包 jar"]
  J --> E
```

初学者可以这样理解：

- 写代码、编译、测试、打包，需要 JDK。
- 运行 Java 程序，核心依赖 JVM。
- Maven 和 Gradle 不替代 JDK，它们是构建工具，会调用 JDK 完成编译和打包。

## 一张图理解对象、引用、堆和栈

```mermaid
flowchart LR
  subgraph Stack["线程栈 Stack"]
    A["局部变量 user"]
    B["局部变量 order"]
  end

  subgraph Heap["堆 Heap"]
    C["User 对象<br/>id=1<br/>name=Ada"]
    D["Order 对象<br/>id=100<br/>userId=1"]
  end

  A -- "引用地址" --> C
  B -- "引用地址" --> D
```

关键理解：

- 局部变量通常在线程栈里。
- `new User()` 创建出来的对象在堆里。
- 变量里保存的不是整个对象，而是指向对象的引用。
- 如果没有任何地方还能引用某个对象，它才可能被 GC 回收。

这也是为什么缓存、静态集合、ThreadLocal 使用不当会造成内存泄漏：业务已经不需要对象了，但某个引用还一直保留着它。

## 一张图理解方法调用栈

```mermaid
sequenceDiagram
  participant C as UserController
  participant S as UserService
  participant R as UserRepository
  participant D as Database

  C->>S: getUser(id)
  S->>R: findById(id)
  R->>D: select * from users where id = ?
  D-->>R: row
  R-->>S: User
  S-->>C: UserView
```

当异常发生时，堆栈通常会从最底层一路打印到入口层。排查时不要只看第一行，要找最关键的 `Caused by`：

```text
Controller
  -> Service
     -> Repository
        -> JDBC Driver
           -> Database error
```

如果错误是 SQL 字段不存在，真正原因通常在 Repository 或数据库迁移，而不是 Controller。

## 一张图理解 Spring Boot 请求链路

```mermaid
flowchart TD
  A["浏览器 / 前端请求"] --> B["Filter<br/>跨域、日志、traceId"]
  B --> C["Interceptor<br/>登录态、权限"]
  C --> D["Controller<br/>路由、参数校验"]
  D --> E["Service<br/>业务编排、事务边界"]
  E --> F["Repository / Mapper<br/>数据访问"]
  F --> G[("Database")]
  E --> H["External Client<br/>调用外部系统"]
  E --> I["Cache<br/>Redis"]
  D --> J["Response<br/>统一 JSON"]
```

这张图能帮助你判断代码应该放在哪里：

| 代码 | 应该放哪里 |
| --- | --- |
| 参数格式校验 | Controller / Request DTO |
| 登录态解析 | Filter / Interceptor |
| 是否允许操作某个订单 | Service |
| SQL 查询 | Repository / Mapper |
| 调用支付系统 | External Client |
| 事务控制 | Service |
| 统一错误格式 | Exception Handler |

最常见的新手问题是 Controller 里写满业务逻辑，最后权限、事务、日志、参数校验全混在一起。

## 一张图理解事务为什么要放在 Service

```mermaid
sequenceDiagram
  participant API as Controller
  participant S as OrderService
  participant TX as Transaction Manager
  participant DB as Database

  API->>S: createOrder(command)
  S->>TX: begin
  S->>DB: insert order
  S->>DB: update stock
  S->>DB: insert payment_record
  alt 全部成功
    S->>TX: commit
    S-->>API: success
  else 任一步失败
    S->>TX: rollback
    S-->>API: error
  end
```

事务不是为了包住一条 SQL，而是为了包住一个完整业务动作。

例如“创建订单”通常包括：

- 写订单表。
- 扣库存。
- 写支付记录。
- 写操作日志或事件。

这些动作要么一起成功，要么一起失败，所以事务边界应该放在 Service 层，而不是 Controller 或 Repository 里随意开启。

## 一张图理解线程池和虚拟线程

```mermaid
flowchart TD
  A["HTTP 请求"] --> B{"执行模型"}

  B --> C["传统线程池"]
  C --> D["固定数量平台线程"]
  D --> E["请求排队"]
  E --> F["线程被数据库 / HTTP I/O 阻塞"]

  B --> G["虚拟线程"]
  G --> H["每个请求一个轻量任务"]
  H --> I["等待 I/O 时让出载体线程"]
  I --> J["更容易支撑大量阻塞式 I/O"]

  F --> K["仍要限制数据库连接和外部接口"]
  J --> K
```

虚拟线程解决的是“阻塞等待时线程成本高”的问题，不解决所有并发问题。

仍然必须控制：

- 数据库连接池。
- 外部接口并发。
- Redis 连接。
- 业务锁竞争。
- 限流和超时。

否则虚拟线程可能让请求更容易并发出去，反而更快打满下游资源。

## 一张图理解 Java 后端排错顺序

```mermaid
flowchart TD
  A["线上问题"] --> B{"表现是什么"}
  B --> C["启动失败"]
  B --> D["接口报错"]
  B --> E["接口变慢"]
  B --> F["内存升高"]
  B --> G["CPU 升高"]

  C --> C1["看最底部 Caused by"]
  C1 --> C2["检查配置、端口、Bean、数据源、依赖版本"]

  D --> D1["根据 traceId 找日志"]
  D1 --> D2["定位 Controller / Service / Repository / 外部系统"]

  E --> E1["拆分耗时"]
  E1 --> E2["SQL、外部 HTTP、Redis、锁、线程池"]

  F --> F1["看 GC 和堆曲线"]
  F1 --> F2["heap dump 分析大对象和引用链"]

  G --> G1["找高 CPU 线程"]
  G1 --> G2["线程 dump 定位循环、锁竞争、序列化"]
```

真实项目里不要凭感觉修改。正确顺序是：

1. 先确认现象。
2. 再收集日志、指标、dump。
3. 再定位层次。
4. 最后修改代码和补测试。

## 建议阅读顺序

如果你是第一次学 Java，建议按这个顺序：

1. 先读本页，建立整体图像。
2. 再读 [环境、JDK 与构建工具](/java/setup-tooling)。
3. 再读 [语法与面向对象](/java/syntax-oop)。
4. 学到 Spring Boot 前，回来看“请求链路”和“事务边界”两张图。
5. 遇到线上问题时，回来看“排错顺序”。

## 下一步学习

继续学习 [环境、JDK 与构建工具](/java/setup-tooling)。
