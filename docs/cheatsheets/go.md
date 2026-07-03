# Go 速查

## 基础命令

| 任务 | 命令 |
| --- | --- |
| 查看版本 | `go version` |
| 查看环境 | `go env` |
| 初始化模块 | `go mod init example.com/app` |
| 整理依赖 | `go mod tidy` |
| 运行项目 | `go run ./cmd/server` |
| 构建二进制 | `go build ./cmd/server` |
| 运行测试 | `go test ./...` |
| 竞态检测 | `go test -race ./...` |
| Benchmark | `go test -bench=. ./...` |
| Fuzzing | `go test -fuzz=FuzzName` |

## 常用类型零值

| 类型 | 零值 |
| --- | --- |
| `string` | `""` |
| `int` | `0` |
| `bool` | `false` |
| `pointer` | `nil` |
| `slice` | `nil` |
| `map` | `nil` |
| `channel` | `nil` |

## 函数和错误

```go
func FindUser(ctx context.Context, id int64) (*User, error) {
    if id <= 0 {
        return nil, ErrInvalidID
    }
    return repo.FindByID(ctx, id)
}
```

错误包装：

```go
return fmt.Errorf("find user id=%d: %w", id, err)
```

判断：

```go
if errors.Is(err, sql.ErrNoRows) {
    return nil, ErrUserNotFound
}
```

## 并发

```go
var wg sync.WaitGroup

for _, item := range items {
    item := item
    wg.Add(1)
    go func() {
        defer wg.Done()
        process(item)
    }()
}

wg.Wait()
```

监听取消：

```go
select {
case result := <-resultCh:
    return result, nil
case <-ctx.Done():
    return nil, ctx.Err()
}
```

## HTTP

```go
func handler(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    _ = ctx
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
}
```

## database/sql

```go
row := db.QueryRowContext(ctx, query, id)

var user User
if err := row.Scan(&user.ID, &user.Name); err != nil {
    return nil, fmt.Errorf("scan user: %w", err)
}
```

事务：

```go
tx, err := db.BeginTx(ctx, nil)
if err != nil {
    return err
}
defer tx.Rollback()

// exec sql

return tx.Commit()
```

## 资源关闭清单

- `rows.Close()`
- `resp.Body.Close()`
- `file.Close()`
- goroutine 退出条件
- context cancel

## 常见排查

| 问题 | 先看 |
| --- | --- |
| goroutine 泄漏 | goroutine profile、channel 阻塞 |
| 内存高 | heap profile、大 slice、缓存 |
| 接口慢 | context、SQL、外部接口、连接池 |
| nil panic | 初始化、指针、map |
| CI 失败 | Go 版本、go.sum、私有依赖 |

## 继续学习

- [Go 学习导览](/go/introduction)
- [并发：goroutine、channel、select](/go/concurrency)
- [Go 常见问题](/go/troubleshooting)
