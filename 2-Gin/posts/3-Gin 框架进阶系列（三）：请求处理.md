# Gin 框架进阶系列（三）：请求处理

---

## 参数绑定

Gin 用 `Bind` 系列方法将请求数据映射到 struct。核心分两类：

| 方法              | 绑定失败时                 |
| ----------------- | -------------------------- |
| `ShouldBind` 系列 | 返回 error，**由自己处理** |
| `Bind` 系列       | 自动返回 400，**不推荐**   |

实际开发统一用 `Should` 系列，掌握主动权。

---

## JSON 绑定

最常见的场景：前端传 JSON body。

```go
type CreateUserReq struct {
    Name  string `json:"name"  binding:"required"`
    Email string `json:"email" binding:"required,email"`
    Age   int    `json:"age"   binding:"required,gte=1,lte=150"`
}

r.POST("/user", func(c *gin.Context) {
    var req CreateUserReq
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"code": 400, "msg": err.Error()})
        return
    }
    c.JSON(200, gin.H{"code": 0, "data": req})
})
```

```bash
# 正常请求
curl -X POST http://localhost:8080/user \
  -H "Content-Type: application/json" \
  -d '{"name":"tom","email":"tom@example.com","age":25}'

# 缺少字段
curl -X POST http://localhost:8080/user \
  -H "Content-Type: application/json" \
  -d '{"name":"tom"}'
# 返回校验错误信息：
# {"code":400,"msg":"Key: 'CreateUserReq.Email' Error:Field 'Email' is required"}
```

---

## 表单绑定

```go
type LoginReq struct {
    Username string `form:"username" binding:"required"`
    Password string `form:"password" binding:"required,min=6"`
}

r.POST("/login", func(c *gin.Context) {
    var req LoginReq
    if err := c.ShouldBind(&req); err != nil {
        c.JSON(400, gin.H{"msg": err.Error()})
        return
    }
    c.JSON(200, gin.H{"user": req.Username})
})
```

`ShouldBind` 会根据 `Content-Type` 自动选择绑定器：`application/json` 走 JSON，`application/x-www-form-urlencoded` 走表单，`multipart/form-data` 走多部分表单。

---

## 查询参数绑定

```go
type SearchReq struct {
    Keyword string `form:"q"    binding:"required"`
    Page    int    `form:"page" binding:"gte=1"`
    Size    int    `form:"size" binding:"gte=1,lte=100"`
}

// GET /search?q=gin&page=1&size=10
r.GET("/search", func(c *gin.Context) {
    var req SearchReq
    if err := c.ShouldBindQuery(&req); err != nil {
        c.JSON(400, gin.H{"msg": err.Error()})
        return
    }
    c.JSON(200, gin.H{"data": req})
})
```

注意 struct tag 用的是 `form` 而不是 `json`，query 参数和表单参数共用 `form` tag。

---

## URI 参数绑定

```go
type UserURI struct {
    ID int `uri:"id" binding:"required,gte=1"`
}

// GET /user/42
r.GET("/user/:id", func(c *gin.Context) {
    var uri UserURI
    if err := c.ShouldBindUri(&uri); err != nil {
        c.JSON(400, gin.H{"msg": err.Error()})
        return
    }
    c.JSON(200, gin.H{"id": uri.ID})
})
```

直接拿到 `int` 类型，省去手动 `strconv`。

---

## 常用校验规则速查

Gin 底层用的是 [go-playground/validator](https://github.com/go-playground/validator)，常用标签：

```go
binding:"required"          // 必填
binding:"email"             // 邮箱格式
binding:"url"               // URL 格式
binding:"min=6"             // 字符串最小长度 6 / 数字最小值 6
binding:"max=20"            // 字符串最大长度 20 / 数字最大值 20
binding:"gte=1,lte=100"     // 数字范围 [1, 100]
binding:"len=11"            // 长度恰好为 11
binding:"oneof=male female" // 枚举值
binding:"eqfield=Password"  // 与另一个字段相等（确认密码场景）
```

多个规则用逗号隔开，是 **AND** 关系。

---

## 自定义校验器

内置规则不够用时，注册自定义校验：

```go
import (
    "github.com/gin-gonic/gin/binding"
    "github.com/go-playground/validator/v10"
)

// 校验手机号（简易版）
var validMobile validator.Func = func(fl validator.FieldLevel) bool {
    mobile := fl.Field().String()
    return len(mobile) == 11 && mobile[0] == '1'
}

func main() {
    r := gin.Default()

    // 注册
    if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
        v.RegisterValidation("mobile", validMobile)
    }

    type Req struct {
        Phone string `json:"phone" binding:"required,mobile"`
    }

    r.POST("/sms", func(c *gin.Context) {
        var req Req
        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(400, gin.H{"msg": err.Error()})
            return
        }
        c.JSON(200, gin.H{"phone": req.Phone})
    })

    r.Run(":8080")
}
```

---

## 校验错误友好化

默认错误信息对用户不友好，翻译一下：

```go
import (
    "github.com/go-playground/validator/v10"
)

func formatValidationErrors(err error) map[string]string {
    errs := make(map[string]string)
    if ve, ok := err.(validator.ValidationErrors); ok {
        for _, fe := range ve {
            switch fe.Tag() {
            case "required":
                errs[fe.Field()] = "不能为空"
            case "email":
                errs[fe.Field()] = "邮箱格式不正确"
            case "min":
                errs[fe.Field()] = "长度不能小于 " + fe.Param()
            case "max":
                errs[fe.Field()] = "长度不能大于 " + fe.Param()
            default:
                errs[fe.Field()] = "校验失败: " + fe.Tag()
            }
        }
    }
    return errs
}

// 使用
if err := c.ShouldBindJSON(&req); err != nil {
    c.JSON(400, gin.H{"code": 400, "errors": formatValidationErrors(err)})
    return
}
```

```json
{
    "code": 400,
    "errors": {
        "Email": "邮箱格式不正确",
        "Age": "校验失败: gte"
    }
}
```

---

## 单文件上传

```go
r.POST("/upload", func(c *gin.Context) {
    file, err := c.FormFile("file")
    if err != nil {
        c.JSON(400, gin.H{"msg": "文件获取失败"})
        return
    }

    // 保存到本地
    dst := "./uploads/" + file.Filename
    if err := c.SaveUploadedFile(file, dst); err != nil {
        c.JSON(500, gin.H{"msg": "保存失败"})
        return
    }

    c.JSON(200, gin.H{
        "msg":  "上传成功",
        "name": file.Filename,
        "size": file.Size,
    })
})
```

```bash
curl -X POST http://localhost:8080/upload \
  -F "file=@./test.png"
```

---

## 多文件上传

```go
r.POST("/uploads", func(c *gin.Context) {
    form, _ := c.MultipartForm()
    files := form.File["files"] // 字段名 files

    for _, f := range files {
        dst := "./uploads/" + f.Filename
        c.SaveUploadedFile(f, dst)
    }

    c.JSON(200, gin.H{
        "msg":   "上传成功",
        "count": len(files),
    })
})
```

---

## 限制上传大小

**MaxMultipartMemory**

```go
// 内存缓冲区限制为 8MB，超出部分写入临时文件
r.MaxMultipartMemory = 8 << 20 // 8 MiB
```

这个值控制的是**内存缓冲区大小**，超出部分会写入临时文件，并非硬性拒绝。要做真正的大小限制，需要在中间件里检查 `Content-Length` 或读取后判断 `file.Size`。

**ContentLength**

通过 `c.Request.ContentLength` 可以获取请求体大小，但需要注意的是, Content-Length 可以被伪造或缺失，更稳妥的做法是在拿到文件后检查 file.Size，或者第三种方式

```go
func LimitUploadSize(maxBytes int64) gin.HandlerFunc {
    return func(c *gin.Context) {
        if c.Request.ContentLength > maxBytes {
            c.AbortWithStatusJSON(413, gin.H{
                "msg": "文件大小超出限制",
            })
            return
        }
        c.Next()
    }
}

// 使用
r.POST("/upload", LimitUploadSize(8<<20), func(c *gin.Context) {
    // ...上传逻辑
})
```

**MaxBytesReader**

```go
const MaxUploadSize = 1 << 20 // 1 MB

r.POST("/upload", func(c *gin.Context) {
    // 包装请求体，硬性限制读取字节数
    c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, MaxUploadSize)

    // 解析 multipart form，触发实际读取
    if err := c.Request.ParseMultipartForm(MaxUploadSize); err != nil {
        // 判断是否为超出大小的错误
        if _, ok := err.(*http.MaxBytesError); ok {
            c.JSON(http.StatusRequestEntityTooLarge, gin.H{
                "error": fmt.Sprintf("文件过大，最大允许 %d 字节", MaxUploadSize),
            })
            return
        }
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    file, header, err := c.Request.FormFile("file")
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "缺少文件字段"})
        return
    }
    defer file.Close()

    // 保存文件...
    c.JSON(http.StatusOK, gin.H{
        "msg":  "上传成功",
        "name": header.Filename,
    })
})
```
```bash
# 小文件 -- 成功
curl -X POST http://localhost:8080/upload -F "file=@small.txt"
# {"msg":"上传成功","name":"small.txt"}

# 大文件 -- 拒绝
curl -X POST http://localhost:8080/upload -F "file=@large.zip"
# {"error":"文件过大，最大允许 1048576 字节"}
```

三个关键点：

**`http.MaxBytesReader`** 是核心。它包装了 `c.Request.Body`，读取超过指定字节数后立刻返回错误，不会把整个大文件读进内存或写入磁盘，从根源上防御了大文件耗尽资源的 DoS 攻击。

**`ParseMultipartForm`** 触发实际读取。调用它时才真正开始消费请求体，此时 `MaxBytesReader` 的限制才生效。

**`*http.MaxBytesError`** 类型断言用于区分"文件过大"和其他解析错误，前者返回 413，后者返回 400，语义更清晰。

之前提到的检查 `Content-Length` 和 `file.Size` 的方式有明显缺陷：`Content-Length` 可以伪造或缺失，`file.Size` 则意味着文件已经被完整读取了，为时已晚。`http.MaxBytesReader` 是**读取阶段就拦截**，是官方推荐的正确姿势。

官方文档：[限制上传大小](https://gin-gonic.com/zh-cn/docs/routing/upload-file/limit-bytes/)


---

## 小结

这篇覆盖了 Gin 请求处理的完整链路：JSON / 表单 / Query / URI 四种绑定方式、validator 内置规则与自定义校验、错误信息友好化、单文件与多文件上传。

下一篇进入**中间件机制深入**：执行顺序、`c.Next()` / `c.Abort()` 原理与常见中间件实现。