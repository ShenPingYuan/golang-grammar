# Gin 框架进阶系列（六）：Gin 认证与授权——JWT 鉴权实战

---

## 为什么选 JWT

传统 Session 方案把用户状态存在服务端，需要集中式存储（内存、Redis），水平扩展时每个节点都得共享 Session。JWT（JSON Web Token）把用户信息编码进 Token 本身，服务端只需验签不需存储，天然适合无状态的 RESTful API 和微服务架构。

当然 JWT 不是万能的，它也有明确的短板，文末会专门讨论。

---

## JWT 结构

一个 JWT 由三段 Base64URL 编码的字符串用 `.` 拼接而成：

```
eyJhbGciOiJIUzI1NiJ9.eyJ1c2VyX2lkIjoxfQ.SIGNATURE
│       Header       │      Payload      │  Signature │
```


**Header** 声明签名算法（如 HS256）。**Payload** 携带声明（claims），包含自定义数据和标准字段如 `exp`（过期时间）。**Signature** 用密钥对前两段做 HMAC 签名，保证内容不可篡改。

Payload 只是 Base64 编码，**不是加密**，任何人都能解码。所以永远不要在 Payload 里放密码等敏感信息。

> 了解JWT更多信息：[jwt.io](https://jwt.io)

---

## 依赖安装

```bash
go get -u github.com/golang-jwt/jwt/v5
```

`golang-jwt/jwt` 是社区维护最活跃的 JWT 库，也是原 `dgrijalva/jwt-go` 的官方继任者。

---

## 项目结构

```
├── main.go
├── config/
│   └── database.go
├── model/
│   └── user.go
├── handler/
│   ├── auth.go        // 注册、登录
│   └── user.go        // 需要鉴权的接口
├── middleware/
│   └── jwt.go         // JWT 中间件
├── pkg/
│   └── jwt.go         // Token 生成与解析
└── go.mod
```

把 Token 的生成和解析抽到 `pkg/jwt.go`，中间件和 handler 都只调用它，职责清晰，且可复用。

---

## 用户模型

在上一篇的基础上加一个密码字段：

```go
// model/user.go
package model

import "gorm.io/gorm"

type User struct {
    gorm.Model
    Name     string `gorm:"type:varchar(50);not null"       json:"name"`
    Email    string `gorm:"type:varchar(100);uniqueIndex"   json:"email"`
    Password string `gorm:"type:varchar(255);not null"      json:"-"` // json:"-" 永远不序列化
}
```

`json:"-"` 保证无论在哪返回 User，密码都不会被输出到 JSON，真实开发使用 UserDto 来返回数据。

---

## 密码加密

密码必须用 bcrypt 哈希存储，永远不要明文或 MD5（太弱了）：

```go
go get -u golang.org/x/crypto/bcrypt
```

```go
// pkg/password.go
package pkg

import "golang.org/x/crypto/bcrypt"

func HashPassword(raw string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(raw), bcrypt.DefaultCost)
    return string(bytes), err
}

func CheckPassword(raw, hashed string) bool {
    return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(raw)) == nil
}
```

`bcrypt.DefaultCost` 是 10，每次哈希大约 100ms，足够抵抗暴力破解。

---

## Token 生成与解析

```go
// pkg/jwt.go
package pkg

import (
    "errors"
    "time"

    "github.com/golang-jwt/jwt/v5"
)

// 密钥，生产环境从配置或环境变量读取
var jwtSecret = []byte("your-256-bit-secret")

type Claims struct {
    UserID uint   `json:"user_id"`
    Email  string `json:"email"`
    jwt.RegisteredClaims
}

// 生成 Token
func GenerateToken(userID uint, email string) (string, error) {
    claims := Claims{
        UserID: userID,
        Email:  email,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            Issuer:    "gin-demo",
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtSecret)
}

// 解析 Token
func ParseToken(tokenString string) (*Claims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (interface{}, error) {
        // 确保签名算法一致，防止 alg:none 攻击
        if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, errors.New("unexpected signing method")
        }
        return jwtSecret, nil
    })
    if err != nil {
        return nil, err
    }

    claims, ok := token.Claims.(*Claims)
    if !ok || !token.Valid {
        return nil, errors.New("invalid token")
    }

    return claims, nil
}
```

几个关键点。`RegisteredClaims` 包含 `exp`、`iat`、`iss` 等标准字段，`jwt/v5` 会在解析时自动检查 `exp` 是否过期，过期直接返回错误，不需要手动判断。签名算法校验那一步是防御经典的 `alg:none` 攻击——攻击者把 Header 的算法改成 `none` 绕过签名验证。

---

## JWT 中间件

```go
// middleware/jwt.go
package middleware

import (
    "strings"

    "your-project/pkg"

    "github.com/gin-gonic/gin"
)

func JWTAuth() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 从 Authorization 头取 Token
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.AbortWithStatusJSON(401, gin.H{"code": 401, "msg": "缺少 Authorization 头"})
            return
        }

        // 格式: Bearer <token>
        parts := strings.SplitN(authHeader, " ", 2)
        if len(parts) != 2 || parts[0] != "Bearer" {
            c.AbortWithStatusJSON(401, gin.H{"code": 401, "msg": "Authorization 格式错误，应为 Bearer <token>"})
            return
        }

        // 解析
        claims, err := pkg.ParseToken(parts[1])
        if err != nil {
            c.AbortWithStatusJSON(401, gin.H{"code": 401, "msg": "Token 无效或已过期"})
            return
        }

        // 将用户信息写入上下文，后续 handler 可直接取用
        c.Set("userID", claims.UserID)
        c.Set("email", claims.Email)
        c.Next()
    }
}
```

不暴露具体的解析错误给客户端，统一返回"Token 无效或已过期"，避免泄露内部信息。

取值辅助函数：

```go
// middleware/jwt.go
func GetCurrentUserID(c *gin.Context) uint {
    id, _ := c.Get("userID")
    return id.(uint)
}
```

---

## 注册与登录

```go
// handler/auth.go
package handler

import (
    "net/http"

    "your-project/config"
    "your-project/model"
    "your-project/pkg"

    "github.com/gin-gonic/gin"
)

type RegisterReq struct {
    Name     string `json:"name"     binding:"required"`
    Email    string `json:"email"    binding:"required,email"`
    Password string `json:"password" binding:"required,min=6,max=32"`
}

type LoginReq struct {
    Email    string `json:"email"    binding:"required,email"`
    Password string `json:"password" binding:"required"`
}

// 注册
func Register(c *gin.Context) {
    var req RegisterReq
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": err.Error()})
        return
    }

    // 密码加密
    hashed, err := pkg.HashPassword(req.Password)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "密码加密失败"})
        return
    }

    user := model.User{
        Name:     req.Name,
        Email:    req.Email,
        Password: hashed,
    }

    if err := config.DB.Create(&user).Error; err != nil {
        c.JSON(http.StatusConflict, gin.H{"code": 409, "msg": "邮箱已注册"})
        return
    }

    c.JSON(http.StatusCreated, gin.H{"code": 0, "msg": "注册成功"})
}

// 登录
func Login(c *gin.Context) {
    var req LoginReq
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": err.Error()})
        return
    }

    // 查用户
    var user model.User
    if err := config.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "邮箱或密码错误"})
        return
    }

    // 校验密码
    if !pkg.CheckPassword(req.Password, user.Password) {
        c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "邮箱或密码错误"})
        return
    }

    // 生成 Token
    token, err := pkg.GenerateToken(user.ID, user.Email)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "Token 生成失败"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "code": 0,
        "msg":  "登录成功",
        "data": gin.H{
            "token": token,
            "user": gin.H{
                "id":    user.ID,
                "name":  user.Name,
                "email": user.Email,
            },
        },
    })
}
```

登录失败时不要区分"邮箱不存在"和"密码错误"，统一返回"邮箱或密码错误"，防止攻击者枚举有效邮箱。

---

## 受保护路由

```go
// handler/user.go
// 获取用户个人信息
func Profile(c *gin.Context) {
    userID := middleware.GetCurrentUserID(c)

    var user model.User
    if err := config.DB.First(&user, userID).Error; err != nil {
        c.JSON(404, gin.H{"code": 404, "msg": "用户不存在"})
        return
    }

    c.JSON(200, gin.H{
        "code": 0,
        "data": gin.H{
            "id":    user.ID,
            "name":  user.Name,
            "email": user.Email,
        },
    })
}
```

---

## 路由注册

```go
// main.go
func main() {
    config.InitDB()
    config.DB.AutoMigrate(&model.User{})

    r := gin.Default()

    // 公开接口
    public := r.Group("/api/v1")
    {
        public.POST("/register", handler.Register)
        public.POST("/login", handler.Login)
    }

    // 需要鉴权的接口
    auth := r.Group("/api/v1")
    auth.Use(middleware.JWTAuth())
    {
        auth.GET("/profile", handler.Profile)
        auth.GET("/users", handler.ListUsers)
        auth.PUT("/users/:id", handler.UpdateUser)
        auth.DELETE("/users/:id", handler.DeleteUser)
    }

    r.Run(":8080")
}
```

同一个前缀 `/api/v1` 分成两个 Group，一个公开一个鉴权，清晰明了。

---

## 测试流程

```bash
# 1. 注册
curl -X POST http://localhost:8080/api/v1/register \
  -H "Content-Type: application/json" \
  -d '{"name":"tom","email":"tom@example.com","password":"123456"}'
# {"code":0,"msg":"注册成功"}

# 2. 登录，拿到 Token
curl -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"email":"tom@example.com","password":"123456"}'
# {"code":0,"data":{"token":"eyJhbG...","user":{...}},"msg":"登录成功"}

# 3. 不带 Token 访问受保护接口
curl http://localhost:8080/api/v1/profile
# {"code":401,"msg":"缺少 Authorization 头"}

# 4. 带 Token 访问
curl http://localhost:8080/api/v1/profile \
  -H "Authorization: Bearer eyJhbG..."
# {"code":0,"data":{"id":1,"name":"tom","email":"tom@example.com"}}
```

---

## Token 滑动刷新机制

Access Token 设置较短的过期时间（如 2 小时），同时签发一个较长有效期的 Refresh Token（如 7 天）。Access Token 过期后，客户端用 Refresh Token 换取新的 Access Token，无需重新登录：

```go
func GenerateTokenPair(userID uint, email string) (accessToken, refreshToken string, err error) {
    // Access Token: 2 小时
    accessClaims := Claims{
        UserID: userID,
        Email:  email,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(2 * time.Hour)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            Issuer:    "gin-demo",
        },
    }
    at := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
    accessToken, err = at.SignedString(jwtSecret)
    if err != nil {
        return
    }

    // Refresh Token: 7 天
    refreshClaims := Claims{
        UserID: userID,
        Email:  email,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            Issuer:    "gin-demo",
        },
    }
    rt := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
    refreshToken, err = rt.SignedString(jwtSecret)
    return
}
```

```go
// handler/auth.go
func RefreshToken(c *gin.Context) {
    var req struct {
        RefreshToken string `json:"refresh_token" binding:"required"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"code": 400, "msg": err.Error()})
        return
    }

    claims, err := pkg.ParseToken(req.RefreshToken)
    if err != nil {
        c.JSON(401, gin.H{"code": 401, "msg": "Refresh Token 无效或已过期"})
        return
    }

    // 签发新的 Token 对
    accessToken, refreshToken, err := pkg.GenerateTokenPair(claims.UserID, claims.Email)
    if err != nil {
        c.JSON(500, gin.H{"code": 500, "msg": "Token 生成失败"})
        return
    }

    c.JSON(200, gin.H{
        "code": 0,
        "data": gin.H{
            "access_token":  accessToken,
            "refresh_token": refreshToken,
        },
    })
}
```

路由注册在公开分组：

```go
public.POST("/refresh", handler.RefreshToken)
```

---

## 基于角色的权限控制（RBAC 简易版）

给 User 模型加一个 `Role` 字段：

```go
type User struct {
    gorm.Model
    Name     string `gorm:"type:varchar(50);not null"     json:"name"`
    Email    string `gorm:"type:varchar(100);uniqueIndex" json:"email"`
    Password string `gorm:"type:varchar(255);not null"    json:"-"`
    Role     string `gorm:"type:varchar(20);default:user" json:"role"` // user / admin
}
```

Claims 中也携带角色：

```go
type Claims struct {
    UserID uint   `json:"user_id"`
    Email  string `json:"email"`
    Role   string `json:"role"`
    jwt.RegisteredClaims
}
```

角色校验中间件：

```go
// middleware/role.go
func RequireRole(roles ...string) gin.HandlerFunc {
    return func(c *gin.Context) {
        role, exists := c.Get("role")
        if !exists {
            c.AbortWithStatusJSON(401, gin.H{"code": 401, "msg": "未认证"})
            return
        }

        roleStr := role.(string)
        for _, r := range roles {
            if r == roleStr {
                c.Next()
                return
            }
        }

        c.AbortWithStatusJSON(403, gin.H{"code": 403, "msg": "权限不足"})
    }
}
```

使用方式：

```go
auth := r.Group("/api/v1")
auth.Use(middleware.JWTAuth())
{
    auth.GET("/profile", handler.Profile)               // 所有登录用户
    auth.DELETE("/users/:id", middleware.RequireRole("admin"), handler.DeleteUser) // 仅管理员
}
```

401 表示"你是谁我不知道"（未认证），403 表示"我知道你是谁但你没权限"（未授权），两个状态码语义不同，不要混用。

---

## JWT 的局限性与应对

**无法主动踢人。** Token 签发后，在过期前服务端无法让它失效。如果用户修改了密码或被封禁，已签发的 Token 仍然有效。应对方案是维护一个黑名单（Redis 存被作废的 Token ID，查询开销很小），或者将 Token 有效期设得足够短。

**Payload 不加密。** 任何人拿到 Token 都能 Base64 解码看到里面的内容。不要放敏感信息，用户 ID 和角色足够了。

**Token 体积较大。** 相比 Session ID 的几十字节，JWT 动辄几百字节，每次请求都要带上。对大多数 API 场景影响不大，但在极端带宽敏感的场景需要注意。

**密钥管理。** HS256 是对称签名，服务端泄露密钥等于全部 Token 被攻破。生产环境应该用足够长的随机密钥，定期轮换，或者用 RS256 非对称签名（公钥验签、私钥签发）。

---

## 小结

本篇覆盖了 JWT 鉴权的完整链路：密码 bcrypt 存储、Token 签发与解析、中间件拦截校验、双 Token 刷新机制、RBAC 角色控制。记住：密码永远哈希存储，Token 里不放敏感数据，不要信任客户端传来的任何东西。