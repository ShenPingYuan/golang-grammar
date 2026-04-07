# Gin 框架进阶系列（九）：优雅关闭（Graceful Shutdown）

## 为什么需要优雅关闭

在生产环境中，部署新版本、扩缩容、服务器维护，都需要重启服务。如果直接 kill 进程，正在处理中的请求会被突然中断——用户看到 502 错误，更糟的是数据库写到一半的事务变成脏数据。

优雅关闭实现的是：收到关闭信号后，停止接受新请求，等待正在处理的请求全部完成，然后再退出。

> 参考：https://gin-gonic.com/zh-cn/docs/server-config/graceful-restart-or-stop/

## 完整实现

```go
// cmd/server/main.go
package main

import (
    "context"
    "errors"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "your-project/internal/database"
    "your-project/internal/router"
)

func main() {
    // 初始化数据库
    db, err := database.Init(os.Getenv("DB_DSN"))
    if err != nil {
        log.Fatalf("failed to connect database: %v", err)
    }

    // 初始化路由
    r := router.Setup(db)

    // 不要用 r.Run()，自己创建 http.Server 来获得更多控制权
    srv := &http.Server{
        Addr:         ":8080",
        Handler:      r,
        ReadTimeout:  10 * time.Second,
        WriteTimeout: 30 * time.Second,
        IdleTimeout:  60 * time.Second,
    }

    // 在 goroutine 中启动服务器
    go func() {
        log.Printf("server starting on %s", srv.Addr)
        if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
            log.Fatalf("server failed: %v", err)
        }
    }()

    // 主 goroutine 等待关闭信号
    quit := make(chan os.Signal, 1)
    // SIGINT: Ctrl+C
    // SIGTERM: docker stop, k8s pod termination
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    sig := <-quit
    log.Printf("received signal: %v, shutting down gracefully...", sig)

    // 给正在处理的请求最多 30 秒时间完成
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    // Shutdown 做了三件事：
    // 1. 停止接受新连接
    // 2. 等待已有连接上的请求处理完毕
    // 3. 超时后强制关闭
    if err := srv.Shutdown(ctx); err != nil {
        log.Printf("server forced shutdown: %v", err)
    }

    // 关闭数据库连接池
    sqlDB, err := db.DB()
    if err == nil {
        sqlDB.Close()
    }

    log.Println("server exited")
}
```

## 逐行解析

**为什么不用 `r.Run()`？** `gin.Engine.Run()` 内部也是创建 `http.Server` 然后调 `ListenAndServe()`，但它没有返回 `*http.Server` 实例，你拿不到它，也就没法调 `Shutdown()`。自己创建 `http.Server` 就拥有了完整的控制权。

**三个 Timeout 的含义。** `ReadTimeout` 是从连接被接受到请求 body 完全读取完毕的最大时间，防止慢速客户端拖住连接。`WriteTimeout` 是从请求 body 读取完毕到响应写完的最大时间，防止 handler 卡死。`IdleTimeout` 是 keep-alive 连接在两次请求之间的最大空闲时间。

```
客户端建立连接 → [ReadTimeout] → 请求读完 → handler 处理 → [WriteTimeout] → 响应发完
                                                                              ↓
                                                                    [IdleTimeout] → 等待下一个请求或关闭
```

**`signal.Notify` 为什么用带缓冲的 channel？** `quit := make(chan os.Signal, 1)` 的缓冲大小为 1。这是因为 `signal.Notify` 不会阻塞等你接收——如果 channel 满了，信号就丢了。缓冲为 1 保证在你还没执行到 `<-quit` 时收到的信号不会丢失。

**`Shutdown` 的 30 秒超时。** 正常情况下 `Shutdown` 会等所有请求处理完后返回。但如果有请求卡住了（比如一个下载大文件的请求），不能无限等下去。30 秒超时到了就强制关闭，这时那些还没完成的请求会收到一个 context canceled 的错误。

**关闭顺序很重要。** 先关 HTTP 服务器（停止接受新请求、等待旧请求完成），再关数据库连接。反过来的话，正在处理中的请求试图访问数据库就会报错。

## 与 Docker/Kubernetes 的配合

Docker 执行 `docker stop` 时会先发 SIGTERM，等 10 秒（默认）后发 SIGKILL。Kubernetes 的 Pod 终止流程也是先 SIGTERM 再 SIGKILL，默认宽限期（`terminationGracePeriodSeconds`）是 30 秒。

你的 Shutdown 超时应该比容器的宽限期短一些，留出清理的余量：

```yaml
# Kubernetes Deployment 配置
spec:
  terminationGracePeriodSeconds: 60  # 给 60 秒宽限
  containers:
    - name: app
      # ...
      lifecycle:
        preStop:
          exec:
            command: ["sh", "-c", "sleep 5"]
            # preStop hook 等 5 秒，让服务先从负载均衡摘除
```

```go
// 应用的 Shutdown 超时设为 45 秒
// 60（宽限期）- 5（preStop）- 10（余量）= 45
ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
```

时间线是这样的：Kubernetes 发 SIGTERM → preStop 等 5 秒（期间 Service 把 Pod 从 endpoints 摘除，新请求不再路由过来）→ 应用开始 Shutdown，最多 45 秒处理完剩余请求 → 应用退出。如果 50 秒后应用还没退出，第 60 秒 Kubernetes 会发 SIGKILL 强制杀掉。

## 健康检查端点

优雅关闭还需要配合健康检查。当应用开始 Shutdown 后，健康检查应该返回"不健康"，让负载均衡器知道不要再发新请求过来：

```go
// 用一个原子变量标记服务器状态
var isShuttingDown atomic.Bool

func Health(c *gin.Context) {
    if isShuttingDown.Load() {
        c.JSON(http.StatusServiceUnavailable, gin.H{
            "status": "shutting_down",
        })
        return
    }
    c.JSON(http.StatusOK, gin.H{
        "status": "ok",
    })
}

// 在 main 函数中，收到信号后立即标记
sig := <-quit
isShuttingDown.Store(true)
// 然后再 Shutdown...
```

---