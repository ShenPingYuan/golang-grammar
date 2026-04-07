# Gin 框架进阶系列（二）：路由详解

---

## 基本路由

Gin 支持所有标准 HTTP 方法：

```go
r.GET("/user", getUser)
r.POST("/user", createUser)
r.PUT("/user", updateUser)
r.DELETE("/user", deleteUser)
r.PATCH("/user", patchUser)
r.OPTIONS("/user", optionsUser)
r.HEAD("/user", headUser)

// 匹配所有方法
r.Any("/all", handler)
```

---

## 路径参数

```go
// :name 必选参数
r.GET("/user/:id", func(c *gin.Context) {
    id := c.Param("id")
    c.JSON(200, gin.H{"id": id})
})

// *action 通配参数，匹配后续所有路径
r.GET("/file/*filepath", func(c *gin.Context) {
    path := c.Param("filepath")
    c.JSON(200, gin.H{"path": path})
})
```

```bash
curl http://localhost:8080/user/42
# {"id":"42"}

curl http://localhost:8080/file/static/css/main.css
# {"path":"/static/css/main.css"}
```

注意 `:id` 返回的是 **string**，需要自己转类型。通配参数 `*filepath` 的值**带前导斜杠**。

---

## 查询参数

```go
// GET /search?q=gin&page=1
r.GET("/search", func(c *gin.Context) {
    q := c.Query("q")               // 无值返回 ""
    page := c.DefaultQuery("page", "1") // 无值返回默认值
    c.JSON(200, gin.H{"q": q, "page": page})
})
```

`Query` 和 `DefaultQuery` 只取 URL 上的 query string。表单参数用 `PostForm` / `DefaultPostForm`，后续篇章会讲。

---

## 路由分组

项目一大，路由散落各处就是灾难。用 `Group` 分组：

```go
func main() {
    r := gin.Default()

    // v1 版本
    v1 := r.Group("/api/v1")
    {
        v1.GET("/users", listUsers)
        v1.POST("/users", createUser)
        v1.GET("/users/:id", getUser)
    }

    // v2 版本
    v2 := r.Group("/api/v2")
    {
        v2.GET("/users", listUsersV2)
    }

    r.Run(":8080")
}
```

分组支持**嵌套**，也支持**给分组单独加中间件**：

```go
authorized := r.Group("/admin")
authorized.Use(AuthMiddleware())
{
    authorized.GET("/dashboard", dashboardHandler) // /admin/dashboard
    authorized.GET("/profile", profileHandler)   // /admin/profile
    authorized.POST("/settings", settingsHandler) // /admin/settings
}
```

花括号 `{}` 纯粹是代码风格，没有语法意义，只是让分组内的路由视觉上更清晰。

---

## 路由拆分

实际项目中不要把路由全堆在 `main.go`。推荐按模块拆文件：

```
router/
├── router.go       // 初始化引擎，注册各模块路由
├── user.go         // 用户相关路由
└── order.go        // 订单相关路由
```

```go
// router/router.go
package router

import "github.com/gin-gonic/gin"

func Setup() *gin.Engine {
    r := gin.New()
    r.Use(gin.Logger(), gin.Recovery())

    RegisterUserRoutes(r)
    RegisterOrderRoutes(r)

    return r
}
```

```go
// router/user.go
package router

import "github.com/gin-gonic/gin"

func RegisterUserRoutes(r *gin.Engine) {
    g := r.Group("/api/users")
    {
        g.GET("", listUsers)
        g.GET("/:id", getUser)
        g.POST("", createUser)
    }
}
```

```go
// main.go
package main

import "gin-blog/router"

func main() {
    r := router.Setup()
    r.Run(":8080")
}
```

---

## 重定向

```go
// HTTP 重定向
r.GET("/old", func(c *gin.Context) {
    c.Redirect(http.StatusMovedPermanently, "/new")
})

// 路由内部重定向（不改变 URL）
r.GET("/internal", func(c *gin.Context) {
    c.Request.URL.Path = "/new"
    r.HandleContext(c)
})
```

---

## 404 与 405 处理

```go
// 自定义 404
r.NoRoute(func(c *gin.Context) {
    c.JSON(404, gin.H{"code": 404, "msg": "page not found"})
})

// 自定义 405（方法不允许）
r.NoMethod(func(c *gin.Context) {
    c.JSON(405, gin.H{"code": 405, "msg": "method not allowed"})
})
```

`NoMethod` 需要配合 `r.HandleMethodNotAllowed = true` 才生效，默认是关闭的。

---

## 小结

这篇覆盖了路由系统的核心内容：路径参数与查询参数的区别、分组与嵌套、路由拆分的工程实践、重定向。

下一篇进入**请求处理**：参数绑定、校验器（validator）与文件上传。