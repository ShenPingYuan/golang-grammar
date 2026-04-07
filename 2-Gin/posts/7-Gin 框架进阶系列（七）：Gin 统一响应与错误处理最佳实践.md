# Gin 框架进阶系列（七）：Gin 统一响应与错误处理最佳实践

---

## 为什么需要统一

翻开很多 Gin 项目的代码，你会看到每个 handler 里都在重复同样的事情：手动拼 `gin.H{"code": 0, "msg": "ok", "data": ...}`，出错时有的返回 400、有的返回 200 带错误码，字段名一会儿是 `msg` 一会儿是 `message`。前端对接时苦不堪言，后端自己维护也头疼。

统一响应和错误处理要解决的就是三个问题：所有接口返回格式一致，错误信息分层（用户看到什么、日志记录什么），handler 专注业务逻辑而不是重复拼 JSON。

---

## 统一响应结构

先定义一个全局通用的响应体：

```go
// pkg/response/response.go
package response

import (
    "net/http"

    "github.com/gin-gonic/gin"
)

// Response 统一响应结构
type Response struct {
    Code int         `json:"code"`         // 业务状态码，0 表示成功
    Msg  string      `json:"msg"`          // 提示信息
    Data interface{} `json:"data,omitempty"` // 业务数据，无数据时省略该字段
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
    c.JSON(http.StatusOK, Response{
        Code: 0,
        Msg:  "ok",
        Data: data,
    })
}

// SuccessWithMsg 成功响应，自定义消息
func SuccessWithMsg(c *gin.Context, msg string, data interface{}) {
    c.JSON(http.StatusOK, Response{
        Code: 0,
        Msg:  msg,
        Data: data,
    })
}

// Fail 失败响应
func Fail(c *gin.Context, httpStatus int, code int, msg string) {
    c.JSON(httpStatus, Response{
        Code: code,
        Msg:  msg,
    })
}

// FailWithData 失败响应，携带附加数据（如字段校验详情）
func FailWithData(c *gin.Context, httpStatus int, code int, msg string, data interface{}) {
    c.JSON(httpStatus, Response{
        Code: code,
        Msg:  msg,
        Data: data,
    })
}
```

这里做了一个重要设计：**HTTP 状态码和业务状态码分离**。HTTP 状态码给网关、负载均衡器、监控系统看（200 表示请求本身成功，401 表示未认证，500 表示服务端故障），业务状态码给前端看（10001 表示参数错误，20001 表示余额不足等等）。两者职责不同，不应混用。

所有接口统一返回的 JSON 永远只有 `code`、`msg`、`data` 三个字段，前端只需要一套解析逻辑。

---

## 自定义业务错误码

错误码不应该在代码里散落 magic number，需要集中管理：

```go
// pkg/errcode/errcode.go
package errcode

// AppError 应用级错误
type AppError struct {
    HttpStatus int    `json:"-"`           // HTTP 状态码，不序列化
    Code       int    `json:"code"`        // 业务错误码
    Msg        string `json:"msg"`         // 面向用户的提示
    Internal   string `json:"-"`           // 内部错误信息，仅日志使用
}

// 实现 error 接口
func (e *AppError) Error() string {
    if e.Internal != "" {
        return e.Internal
    }
    return e.Msg
}

// New 创建 AppError
func New(httpStatus, code int, msg string) *AppError {
    return &AppError{
        HttpStatus: httpStatus,
        Code:       code,
        Msg:        msg,
    }
}

// WithInternal 附加内部错误信息
func (e *AppError) WithInternal(err error) *AppError {
    // 返回新对象，不修改原始定义
    newErr := *e
    newErr.Internal = err.Error()
    return &newErr
}
```

`WithInternal` 返回新对象而不是修改原对象，这一点很关键。因为下面我们会把错误码定义成包级变量，如果直接修改原始变量，并发场景下会出问题。

---

## 集中定义错误码

```go
// pkg/errcode/code.go
package errcode

import "net/http"

// 通用错误 10xxx
var (
    ErrBadRequest     = New(http.StatusBadRequest, 10001, "请求参数错误")
    ErrUnauthorized   = New(http.StatusUnauthorized, 10002, "未登录或 Token 已过期")
    ErrForbidden      = New(http.StatusForbidden, 10003, "权限不足")
    ErrNotFound       = New(http.StatusNotFound, 10004, "资源不存在")
    ErrInternalServer = New(http.StatusInternalServerError, 10005, "服务器内部错误")
    ErrTooManyRequest = New(http.StatusTooManyRequests, 10006, "请求过于频繁")
)

// 用户模块 20xxx
var (
    ErrUserExist      = New(http.StatusConflict, 20001, "该邮箱已注册")
    ErrUserNotFound   = New(http.StatusNotFound, 20002, "用户不存在")
    ErrPasswordWrong  = New(http.StatusUnauthorized, 20003, "邮箱或密码错误")
    ErrTokenGenerate  = New(http.StatusInternalServerError, 20004, "Token 生成失败")
)

// 订单模块 30xxx
var (
    ErrOrderNotFound  = New(http.StatusNotFound, 30001, "订单不存在")
    ErrOrderCancelled = New(http.StatusBadRequest, 30002, "订单已取消，无法操作")
)
```

按模块分段编码：10xxx 通用，20xxx 用户，30xxx 订单。这样前端看到错误码就能大致定位问题域。

---

## 全局错误处理中间件

核心思路是让 handler 不再直接调用 `c.JSON`，而是通过 `c.Error()` 把错误挂到上下文，由统一的中间件收口处理：

```go
// middleware/error_handler.go
package middleware

import (
    "errors"
    "log"
    "net/http"

    "your-project/pkg/errcode"
    "your-project/pkg/response"

    "github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Next() // 先执行后续 handler

        // 没有错误，直接返回
        if len(c.Errors) == 0 {
            return
        }

        // 取最后一个错误（通常也是最重要的那个）
        err := c.Errors.Last().Err

        // 判断是否是自定义 AppError
        var appErr *errcode.AppError
        if errors.As(err, &appErr) {
            // 记录内部错误到日志
            if appErr.Internal != "" {
                log.Printf("[BizError] code=%d msg=%s internal=%s path=%s",
                    appErr.Code, appErr.Msg, appErr.Internal, c.Request.URL.Path)
            }
            response.Fail(c, appErr.HttpStatus, appErr.Code, appErr.Msg)
            return
        }

        // 未知错误，统一返回 500，不暴露内部信息
        log.Printf("[UnknownError] err=%v path=%s", err, c.Request.URL.Path)
        response.Fail(c, http.StatusInternalServerError,
            errcode.ErrInternalServer.Code,
            errcode.ErrInternalServer.Msg)
    }
}
```

这个中间件做了两件非常重要的事。第一，它区分了业务错误和未知错误——业务错误返回友好的提示，未知错误一律返回"服务器内部错误"，绝不把堆栈或数据库信息暴露给客户端。第二，内部错误信息只进日志，不进响应，这就是错误分层。

---

## 重构 Handler

有了统一响应和错误中间件后，handler 变得极其干净：

### 重构前（散乱版本）

```go
func Login(c *gin.Context) {
    var req LoginReq
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"code": 400, "message": "参数错误"})
        return
    }

    var user model.User
    if err := config.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
        c.JSON(401, gin.H{"code": 401, "msg": "邮箱或密码错误"})
        return
    }

    if !pkg.CheckPassword(req.Password, user.Password) {
        c.JSON(401, gin.H{"code": 401, "msg": "邮箱或密码不对"})
        return
    }

    token, err := pkg.GenerateToken(user.ID, user.Email)
    if err != nil {
        c.JSON(500, gin.H{"code": 500, "msg": "token 生成失败"})
        return
    }

    c.JSON(200, gin.H{"code": 0, "msg": "ok", "data": gin.H{"token": token}})
}
```

字段名不统一（`message` vs `msg`），HTTP 状态码和业务码混为一谈，错误提示各自为战。

### 重构后（统一版本）

```go
func Login(c *gin.Context) {
    var req LoginReq
    if err := c.ShouldBindJSON(&req); err != nil {
        c.Error(errcode.ErrBadRequest.WithInternal(err))
        return
    }

    var user model.User
    if err := config.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
        c.Error(errcode.ErrPasswordWrong)
        return
    }

    if !pkg.CheckPassword(req.Password, user.Password) {
        c.Error(errcode.ErrPasswordWrong)
        return
    }

    token, err := pkg.GenerateToken(user.ID, user.Email)
    if err != nil {
        c.Error(errcode.ErrTokenGenerate.WithInternal(err))
        return
    }

    response.Success(c, gin.H{"token": token})
}
```

handler 里再也看不到 HTTP 状态码和 JSON 拼接的噪音。每一行都在表达业务意图：绑定失败就是参数错误，查不到就是密码错误，生成失败就是 Token 异常，成功就返回数据。读代码的人三秒就能看懂。

---

## 参数校验错误的友好输出

`ShouldBindJSON` 返回的校验错误信息默认是英文且不够友好。我们可以把校验错误翻译成结构化的字段级提示：

```go
// pkg/response/validator.go
package response

import (
    "github.com/gin-gonic/gin"
    "github.com/go-playground/validator/v10"

    "your-project/pkg/errcode"
)

type FieldError struct {
    Field   string `json:"field"`
    Message string `json:"message"`
}

// HandleValidationError 处理参数校验错误
func HandleValidationError(c *gin.Context, err error) {
    errs, ok := err.(validator.ValidationErrors)
    if !ok {
        // 不是 validator 的错误（比如 JSON 语法错误）
        c.Error(errcode.ErrBadRequest.WithInternal(err))
        return
    }

    fieldErrors := make([]FieldError, 0, len(errs))
    for _, e := range errs {
        fieldErrors = append(fieldErrors, FieldError{
            Field:   e.Field(),
            Message: msgForTag(e),
        })
    }

    // 使用 Abort 直接返回，不再走 ErrorHandler
    c.AbortWithStatusJSON(400, Response{
        Code: errcode.ErrBadRequest.Code,
        Msg:  "请求参数校验失败",
        Data: fieldErrors,
    })
}

func msgForTag(e validator.FieldError) string {
    switch e.Tag() {
    case "required":
        return "不能为空"
    case "email":
        return "邮箱格式不正确"
    case "min":
        return "长度不能小于 " + e.Param()
    case "max":
        return "长度不能大于 " + e.Param()
    case "oneof":
        return "值必须是以下之一: " + e.Param()
    default:
        return "校验失败: " + e.Tag()
    }
}
```

handler 中使用：

```go
if err := c.ShouldBindJSON(&req); err != nil {
    response.HandleValidationError(c, err)
    return
}
```

前端收到的响应：

```json
{
    "code": 10001,
    "msg": "请求参数校验失败",
    "data": [
        {"field": "Email", "message": "邮箱格式不正确"},
        {"field": "Password", "message": "长度不能小于 6"}
    ]
}
```

每个字段具体哪里不合格，一目了然。

---

## Panic Recovery 增强

Gin 自带的 `Recovery()` 中间件会捕获 panic 并返回 500，但它的响应格式不符合我们的统一结构。替换掉它：

```go
// middleware/recovery.go
package middleware

import (
    "log"
    "net/http"
    "runtime/debug"

    "your-project/pkg/errcode"
    "your-project/pkg/response"

    "github.com/gin-gonic/gin"
)

func CustomRecovery() gin.HandlerFunc {
    return func(c *gin.Context) {
        defer func() {
            if r := recover(); r != nil {
                // 打印完整堆栈
                log.Printf("[PANIC] %v\n%s", r, debug.Stack())

                // 返回统一格式
                response.Fail(c,
                    http.StatusInternalServerError,
                    errcode.ErrInternalServer.Code,
                    errcode.ErrInternalServer.Msg,
                )
                c.Abort()
            }
        }()
        c.Next()
    }
}
```

注册时用 `CustomRecovery` 替换默认的 `Recovery`：

```go
r := gin.New() // 不用 gin.Default()，它会自动加 Logger + Recovery
r.Use(gin.Logger())
r.Use(middleware.CustomRecovery())
r.Use(middleware.ErrorHandler())
```

注意中间件的顺序。`CustomRecovery` 必须在最外层，这样它才能兜住所有下游中间件和 handler 中的 panic。`ErrorHandler` 放在 `CustomRecovery` 之后，处理正常的业务错误。

---

## 404 和 405 处理

Gin 对未匹配的路由默认返回 `404 page not found` 纯文本，格式同样不统一：

```go
// main.go
r.NoRoute(func(c *gin.Context) {
    response.Fail(c, http.StatusNotFound, errcode.ErrNotFound.Code, "接口不存在")
})

r.NoMethod(func(c *gin.Context) {
    response.Fail(c, http.StatusMethodNotAllowed, 10007, "请求方法不允许")
})
```

别忘了启用 405 检测，Gin 默认是关闭的：

```go
r.HandleMethodNotAllowed = true
```

---

## 分页响应的封装

列表接口几乎都需要分页，单独封装一个分页响应：

```go
// pkg/response/pagination.go
package response

import "github.com/gin-gonic/gin"

type PageResult struct {
    List     interface{} `json:"list"`
    Total    int64       `json:"total"`
    Page     int         `json:"page"`
    PageSize int         `json:"page_size"`
}

func SuccessWithPage(c *gin.Context, list interface{}, total int64, page, pageSize int) {
    Success(c, PageResult{
        List:     list,
        Total:    total,
        Page:     page,
        PageSize: pageSize,
    })
}
```

handler 中使用：

```go
func ListUsers(c *gin.Context) {
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
    if pageSize > 100 {
        pageSize = 100 // 防止一次拉太多
    }

    var users []model.User
    var total int64

    config.DB.Model(&model.User{}).Count(&total)
    config.DB.Offset((page - 1) * pageSize).Limit(pageSize).Find(&users)

    response.SuccessWithPage(c, users, total, page, pageSize)
}
```

前端拿到的格式永远是：

```json
{
    "code": 0,
    "msg": "ok",
    "data": {
        "list": [...],
        "total": 42,
        "page": 1,
        "page_size": 10
    }
}
```

---

## 请求 ID 追踪

当用户反馈"接口报错了"，你需要快速从海量日志里找到那一条。给每个请求分配一个唯一 ID，响应头和日志都带上它：

```go
// middleware/request_id.go
package middleware

import (
    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
)

func RequestID() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 优先从请求头取（上游网关可能已经生成了）
        rid := c.GetHeader("X-Request-ID")
        if rid == "" {
            rid = uuid.New().String()
        }
        c.Set("request_id", rid)
        c.Header("X-Request-ID", rid)
        c.Next()
    }
}
```

在 `ErrorHandler` 中把 request_id 一起记进日志：

```go
rid, _ := c.Get("request_id")
log.Printf("[BizError] request_id=%v code=%d msg=%s internal=%s path=%s",
    rid, appErr.Code, appErr.Msg, appErr.Internal, c.Request.URL.Path)
```

前端报错时只需要提供 `X-Request-ID`，后端一条命令就能 grep 到。

---

## 完整中间件注册顺序

```go
func main() {
    r := gin.New()

    // 第一层：请求 ID（最早分配，后续所有中间件都能用）
    r.Use(middleware.RequestID())

    // 第二层：Panic 兜底
    r.Use(middleware.CustomRecovery())

    // 第三层：日志
    r.Use(gin.Logger())

    // 第四层：统一错误处理
    r.Use(middleware.ErrorHandler())

    // 404 / 405
    r.HandleMethodNotAllowed = true
    r.NoRoute(func(c *gin.Context) {
        response.Fail(c, 404, errcode.ErrNotFound.Code, "接口不存在")
    })
    r.NoMethod(func(c *gin.Context) {
        response.Fail(c, 405, 10007, "请求方法不允许")
    })

    // 业务路由
    registerRoutes(r)

    r.Run(":8080")
}
```

中间件的顺序就是洋葱模型从外到内的顺序。RequestID 在最外层保证所有后续操作都能拿到 ID；Recovery 在第二层兜住一切 panic；Logger 在第三层记录请求日志；ErrorHandler 在第四层收口业务错误。

---

## 设计原则总结

第一，响应结构固定为 `code` + `msg` + `data`，没有例外。前端只需要一个拦截器就能统一处理所有接口。

第二，错误分两层。面向用户的 `Msg` 要可读、可展示（"邮箱或密码错误"），面向开发者的 `Internal` 只进日志（"record not found: SELECT * FROM users WHERE email = ..."）。永远不要把数据库错误、堆栈信息返回给客户端。

第三，handler 只做三件事——绑定参数、调用业务逻辑、返回成功响应或 `c.Error`。任何格式拼装、状态码映射、日志记录的工作都交给中间件和工具函数。handler 越薄，越不容易出 bug。

第四，错误码集中管理，按模块分段。新增错误只需要在 `errcode/code.go` 里加一行，grep 错误码就能找到所有使用的地方。

---