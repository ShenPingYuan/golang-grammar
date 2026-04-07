`context.Context` 是 Go 里用来**传递截止时间、取消信号和请求范围值**的接口。

---

## 核心用途

| 场景         | 说明                           |
| ------------ | ------------------------------ |
| **超时控制** | 设置请求最大执行时间           |
| **取消信号** | 主动取消正在进行的操作         |
| **传递数据** | 在函数调用链中共享请求级别的值 |

---

## 基本用法

### 1. 创建 Context

```go
// 背景 context，所有 context 的根
ctx := context.Background()

// 带取消
ctx, cancel := context.WithCancel(parentCtx)
defer cancel()  // 记得调用，释放资源

// 带超时（常用）
ctx, cancel := context.WithTimeout(parentCtx, 5*time.Second)
defer cancel()

// 带截止时间
ctx, cancel := context.WithDeadline(parentCtx, time.Now().Add(5*time.Second))
```

### 2. 检查取消/超时

```go
select {
case <-ctx.Done():
    // context 被取消或超时了
    return ctx.Err()  // context.Canceled 或 context.DeadlineExceeded
default:
    // 继续执行
}
```

### 3. 传递数据（不推荐滥用）

```go
// 存
ctx := context.WithValue(parentCtx, "userID", 123)

// 取
userID := ctx.Value("userID").(int)
```

**注意：** 不要放业务关键数据，只放请求范围元数据（如 traceID、userID 等）。

---

## 实际例子：HTTP 请求超时

```go
func fetchData(ctx context.Context) error {
    req, err := http.NewRequestWithContext(ctx, "GET", "https://api.example.com", nil)
    if err != nil {
        return err
    }
    
    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return err  // 如果 ctx 超时，这里会返回错误
    }
    defer resp.Body.Close()
    
    // 处理响应...
    return nil
}

func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()
    
    if err := fetchData(ctx); err != nil {
        log.Println("失败:", err)
    }
}
```

---

## 关键规则

1. **Context 是只读的** — 派生新的 context 来添加功能
2. **不要存 nil context** — 不确定时用 `context.Background()`
3. **及时调用 cancel()** — 防止 goroutine 泄漏
4. **函数第一个参数** — 标准写法是 `func DoSomething(ctx context.Context, ...)`

---

## Gin 里的 Context

Gin 的 `*gin.Context` 和 Go 的 `context.Context` 是两回事：

```go
// Gin context 包装了标准 context
func handler(c *gin.Context) {
    // 获取标准 context
    ctx := c.Request.Context()
    
    // 传给需要标准 context 的函数
    data, err := db.Query(ctx, "SELECT * FROM users")
}
```

Gin 的 `c.Set()` / `c.Get()` 是框架层面的，和标准 `context.WithValue` 无关。

---

哪里还想深入了解？超时控制、取消传播、还是实际项目中的用法？