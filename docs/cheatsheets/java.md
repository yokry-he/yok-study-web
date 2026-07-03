# Java 速查

## 基础命令

| 任务 | 命令 |
| --- | --- |
| 查看 Java 版本 | `java --version` |
| 查看编译器版本 | `javac --version` |
| 编译单文件 | `javac Hello.java` |
| 运行类 | `java Hello` |
| Maven 测试 | `mvn test` |
| Maven 打包 | `mvn package` |
| Spring Boot 启动 | `mvn spring-boot:run` |
| Gradle 测试 | `./gradlew test` |
| Gradle 构建 | `./gradlew build` |

## 常用类型

| 类型 | 用法 |
| --- | --- |
| `String` | 字符串 |
| `BigDecimal` | 金额和精确小数 |
| `LocalDate` | 日期 |
| `LocalDateTime` | 日期时间 |
| `Optional<T>` | 可能不存在的返回值 |
| `List<T>` | 有序列表 |
| `Set<T>` | 去重集合 |
| `Map<K,V>` | key-value 映射 |

## 集合选择

| 场景 | 推荐 |
| --- | --- |
| 普通列表 | `ArrayList` |
| 去重 | `HashSet` |
| 保留插入顺序去重 | `LinkedHashSet` |
| 按 key 查询 | `HashMap` |
| 保留 key 顺序 | `LinkedHashMap` |
| 并发 map | `ConcurrentHashMap` |

## Stream 常用写法

```java
List<Long> ids = users.stream()
    .map(User::id)
    .toList();

Map<Long, User> userMap = users.stream()
    .collect(Collectors.toMap(User::id, item -> item));

Map<Long, List<Order>> ordersByUser = orders.stream()
    .collect(Collectors.groupingBy(Order::userId));
```

## 异常处理

```java
try {
    service.doWork();
} catch (BusinessException e) {
    log.warn("business failed, code={}", e.code(), e);
    throw e;
} catch (Exception e) {
    log.error("system failed", e);
    throw new SystemException("SYSTEM_ERROR", "系统异常");
}
```

## Spring Boot 常用注解

| 注解 | 用法 |
| --- | --- |
| `@SpringBootApplication` | 启动类 |
| `@RestController` | REST Controller |
| `@RequestMapping` | 路由前缀 |
| `@GetMapping` / `@PostMapping` | HTTP 方法 |
| `@Service` | 业务服务 |
| `@Repository` | 数据访问 |
| `@Transactional` | 事务边界 |
| `@ConfigurationProperties` | 配置绑定 |
| `@Valid` | 参数校验 |

## 事务排查

事务不生效时检查：

- 方法是否 `public`。
- 是否同类内部调用。
- 异常是否被吞掉。
- 是否抛出会触发回滚的异常。
- 注解是否放在 Service 层。
- 是否进入了异步线程。

## JVM 排查

| 问题 | 看什么 |
| --- | --- |
| CPU 高 | 线程 dump、高 CPU 线程 |
| 内存高 | heap dump、对象引用链 |
| GC 频繁 | GC 日志、堆使用曲线 |
| 请求卡住 | 线程 dump、锁等待、连接池 |
| 类冲突 | 依赖树、classpath |

## 继续学习

- [Java 学习导览](/java/introduction)
- [JVM 内存、GC 与诊断](/java/jvm-memory-gc)
- [Java 常见问题](/java/troubleshooting)
