# Gin 框架进阶系列（四）：中间件机制深入

---

## 中间件本质

Gin 的中间件就是一个 `gin.HandlerFunc`，和普通路由处理函数签名完全一样：

```go
type HandlerFunc func(*gin.Context)
```

没有任何特殊接口，**能当 handler 的就能当中间件**。区别仅在于：中间件通常调用 `c.Next()` 把控制权交给下一个处理函数，处理完再回来执行后续逻辑。

---

## 执行顺序与洋葱模型

```go
func A(c *gin.Context) {
    fmt.Println("A - 前")
    c.Next()
    fmt.Println("A - 后")
}

func B(c *gin.Context) {
    fmt.Println("B - 前")
    c.Next()
    fmt.Println("B - 后")
}

func handler(c *gin.Context) {
    fmt.Println("handler")
    c.JSON(200, gin.H{"msg": "ok"})
}

r.Use(A, B)
r.GET("/test", handler)
```

输出：

```
A - 前
B - 前
handler
B - 后
A - 后
```

像洋葱一样，请求从外层进、从外层出。`c.Next()` 之前是**请求阶段**，之后是**响应阶段**。这是理解所有中间件行为的基础。

---

## c.Next() 与 c.Abort() 原理

Gin 内部维护了一个 handler 链（slice）和一个索引 `index`：

```go
// 简化版源码逻辑
func (c *Context) Next() {
    c.index++
    for c.index < int8(len(c.handlers)) {
        c.handlers[c.index](c)
        c.index++
    }
}

func (c *Context) Abort() {
    c.index = abortIndex // 设为极大值，后续 handler 不再执行
}
```

`c.Next()` 递增索引并依次执行后续 handler。`c.Abort()` 将索引设为极大值，循环条件不满足，链路直接终止。

关键点：`c.Abort()` **不会终止当前函数的执行**，它只是阻止后续 handler 运行。如果你要立刻返回，必须 `return`：

```go
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        if token == "" {
            c.AbortWithStatusJSON(401, gin.H{"msg": "未授权"})
            return // 不 return 的话，下面的代码还会继续执行
        }
        c.Next()
    }
}
```

`c.AbortWithStatusJSON` = `c.Abort()` + `c.JSON()`，一步到位。

---

## 中间件注册方式

```go
// 全局中间件 —— 对所有路由生效
r.Use(Logger(), Recovery())

// 分组中间件 —— 仅对该分组生效
admin := r.Group("/admin")
admin.Use(AuthMiddleware())

// 单路由中间件 —— 仅对该路由生效
r.GET("/debug", DebugOnly(), debugHandler)
```

执行顺序：**全局 → 分组 → 单路由**，按注册顺序排列。

---

## 中间件间传值：c.Set / c.Get

中间件解析出的数据需要传给后续 handler，用 `c.Set` 和 `c.Get`：

```go
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        // 假设解析出 userID
        userID, err := parseToken(token)
        if err != nil {
            c.AbortWithStatusJSON(401, gin.H{"msg": "token 无效"})
            return
        }
        c.Set("userID", userID)
        c.Next()
    }
}

// 后续 handler 取值
r.GET("/profile", AuthMiddleware(), func(c *gin.Context) {
    userID, exists := c.Get("userID")
    if !exists {
        c.JSON(500, gin.H{"msg": "userID not found"})
        return
    }
    c.JSON(200, gin.H{"userID": userID})
})
```

`c.Get` 返回 `any` 类型，需要自己断言。可以封装一个辅助函数避免重复代码：

```go
func GetUserID(c *gin.Context) (int64, bool) {
    v, exists := c.Get("userID")
    if !exists {
        return 0, false
    }
    id, ok := v.(int64)
    return id, ok
}
```

---

## 实战：耗时统计中间件

最经典的中间件范例，利用洋葱模型在请求前后各记一次时间：

```go
func Timer() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()

        c.Next() // 等所有后续 handler 执行完

        duration := time.Since(start)
        status := c.Writer.Status()
        log.Printf("[%d] %s %s - %v", status, c.Request.Method, c.Request.URL.Path, duration)
    }
}
```

`c.Next()` 之后拿到的 `status` 和 `duration` 是最终结果，因为所有 handler 已经跑完了。

---

## 实战：CORS 跨域中间件

```go
func CORS() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Header("Access-Control-Allow-Origin", "*")
        c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
        c.Header("Access-Control-Allow-Headers", "Content-Type,Authorization")
        c.Header("Access-Control-Max-Age", "86400")

        // 预检请求直接返回
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }

        c.Next()
    }
}
```

生产环境不要用 `*`，应该配置为具体的前端域名。

---

## 实战：Recovery 中间件手写版

理解 Gin 内置 `Recovery` 的原理：

```go
func Recovery() gin.HandlerFunc {
    return func(c *gin.Context) {
        defer func() {
            if err := recover(); err != nil {
                log.Printf("[PANIC] %v\n%s", err, debug.Stack())
                c.AbortWithStatusJSON(500, gin.H{"msg": "服务器内部错误"})
            }
        }()
        c.Next()
    }
}
```

`defer` + `recover` 捕获 panic，防止单个请求崩掉整个进程。`debug.Stack()` 打印完整调用栈方便排查。

---

## 中间件踩坑点

**不调用 `c.Next()` 会怎样？** 后续 handler 仍然会执行。因为 Gin 的 for 循环会自动递增索引往下走。`c.Next()` 的作用是让你能在后续 handler 执行完**之后**做事。如果你不需要"后置逻辑"，不调用 `c.Next()` 也没问题。

**`c.Abort()` 之后还能写响应吗？** 可以。`Abort` 只阻止后续 handler，不影响当前函数写响应。但要注意不要多次写响应导致 `http: superfluous response.WriteHeader call` 警告。

**异步 goroutine 中不能用原始 `c`。** 必须用 `c.Copy()`，因为请求结束后，Gin 出于性能考取 `c` 会被复用，直接用会数据竞争：

```go
func AsyncMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        cp := c.Copy() // 拷贝一份
        go func() {
            time.Sleep(3 * time.Second)
            log.Printf("异步处理: %s", cp.Request.URL.Path)
        }()
        c.Next()
    }
}
```

---

## 小结

中间件的全部核心就三件事：洋葱模型决定执行顺序，`c.Next()` 分割前置/后置逻辑，`c.Abort()` 截断链路。掌握了这三点，无论是鉴权、限流、日志、Recovery，都是同一个套路的不同变体。

下一篇进入**Gin + GORM：连接数据库实现 CRUD**。
