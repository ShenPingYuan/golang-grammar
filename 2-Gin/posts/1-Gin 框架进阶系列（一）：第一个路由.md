# Gin 框架进阶系列（一）：安装与第一个路由

---

## 前置条件

Go 1.18+，已配置好 `GOPATH` 和 `GOPROXY`（国内建议 `https://goproxy.cn,direct`）。

---

## 初始化项目

```bash
mkdir gin-blog && cd gin-blog
go mod init gin-blog
go get -u github.com/gin-gonic/gin
```

---

## 最小可运行示例

```go
package main

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default() // 内置 Logger + Recovery 中间件

    r.GET("/ping", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "message": "pong",
        })
    })

    r.Run(":8080") // 默认 0.0.0.0:8080
}
```

```bash
go run main.go
curl http://localhost:8080/ping
# {"message":"pong"}
```

`gin.Default()` 与 `gin.New()` 的区别：前者自带 Logger 和 Recovery 两个中间件，后者是裸引擎。生产环境建议用 `gin.New()` 自行控制中间件。

---

## gin.H 是什么

```go
// gin 源码
type H map[string]any
```

就是 `map[string]interface{}` 的别名，方便写 JSON 响应。结构化场景建议用 struct 替代：

```go
type Response struct {
    Code int    `json:"code"`
    Msg  string `json:"msg"`
    Data any    `json:"data"`
}

r.GET("/ping", func(c *gin.Context) {
    c.JSON(http.StatusOK, Response{
        Code: 0,
        Msg:  "success",
        Data: nil,
    })
})
```

---

## 运行模式

Gin 有三种模式：`debug`、`release`、`test`。

```go
// 方式一：代码设置
gin.SetMode(gin.ReleaseMode)

// 方式二：环境变量
// export GIN_MODE=release
```

`debug` 模式会打印路由表和调试日志，**部署时务必切到 `release`**，否则性能白白浪费在日志 I/O 上。

---

## 小结

这一篇只做了一件事：用最少代码跑起来一个 Gin 服务，理解 `gin.Default()`、`gin.H`、运行模式三个核心概念。

下一篇进入**路由系统详解**：分组、路径参数、查询参数、重定向与 404 处理。