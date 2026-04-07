# Gin 框架进阶系列（十一）：静态文件服务与文件下载

---

前面十篇的重心是 JSON API，请求进来、处理逻辑、返回 JSON。但实际项目中，几乎一定会遇到需要提供文件的场景：用户上传了头像需要访问、后台生成了 Excel 报表需要下载、前端打包后的静态资源需要托管。Gin 提供了两套机制分别处理这两类需求——路由级别的静态文件服务和处理函数级别的文件响应。

---

## 一、静态文件服务：整个目录对外开放

当你有一整个目录的文件需要通过 HTTP 访问时，用路由级别的 API。

### 三个方法，各有用途

```go
package main

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()

    // 方法一：Static —— 将 URL 路径映射到本地目录
    // 访问 /assets/css/style.css → 读取 ./public/css/style.css
    r.Static("/assets", "./public")

    // 方法二：StaticFS —— 接受 http.FileSystem 接口，控制力更强
    // 可以对接 embed.FS、自定义存储、或限制目录列表
    r.StaticFS("/uploads", http.Dir("/var/www/uploads"))

    // 方法三：StaticFile —— 单个文件映射
    // 只为一个固定路径提供一个固定文件
    r.StaticFile("/favicon.ico", "./resources/favicon.ico")
    r.StaticFile("/robots.txt", "./resources/robots.txt")

    r.Run(":8080")
}
```

`Static` 和 `StaticFS` 的区别在于参数类型。`Static` 直接接收一个目录路径字符串，内部帮你包装成 `http.Dir`。`StaticFS` 接收 `http.FileSystem` 接口，意味着你可以传入任何实现了这个接口的东西——Go 1.16 引入的 `embed.FS`、内存文件系统、甚至自己写的对象存储适配器。

`StaticFile` 是单文件版本。

### 用 embed.FS 打包静态资源

Go 的 `embed` 包可以在编译时把文件嵌入二进制文件里，部署时不需要额外带一堆静态文件：

```go
package main

import (
    "embed"
    "io/fs"
    "net/http"

    "github.com/gin-gonic/gin"
)

//go:embed static/*
var staticFiles embed.FS

func main() {
    r := gin.Default()

    // embed.FS 的根路径包含 "static" 这一层目录
    // 用 fs.Sub 去掉前缀，这样访问 /assets/style.css 就对应 static/style.css
    subFS, _ := fs.Sub(staticFiles, "static")
    r.StaticFS("/assets", http.FS(subFS))

    r.Run(":8080")
}
```

这个技巧在容器化部署时特别有用。Docker 镜像里只需要一个二进制文件，静态资源全部编译进去了，不用担心路径问题、不用 `COPY` 额外的文件夹。

### 禁用目录列表

`http.Dir` 默认允许目录列表——如果用户访问 `/assets/`（末尾有斜杠，且目录下没有 `index.html`），Nginx 会列出目录里所有文件。这在生产环境是个安全隐患。自定义一个 `http.FileSystem` 来禁用它：

```go
package main

import (
    "net/http"
    "os"

    "github.com/gin-gonic/gin"
)

// noListingFS 包装 http.Dir，禁止目录列表
type noListingFS struct {
    fs http.FileSystem
}

func (nfs noListingFS) Open(name string) (http.File, error) {
    f, err := nfs.fs.Open(name)
    if err != nil {
        return nil, err
    }

    // 检查是不是目录
    info, err := f.Stat()
    if err != nil {
        f.Close()
        return nil, err
    }

    // 如果是目录，返回 404 而不是文件列表
    if info.IsDir() {
        // 尝试打开目录下的 index.html
        index := name + "/index.html"
        if _, err := nfs.fs.Open(index); err != nil {
            f.Close()
            return nil, os.ErrNotExist
        }
    }

    return f, nil
}

func main() {
    r := gin.Default()

    // 使用自定义的 FileSystem
    r.StaticFS("/assets", noListingFS{http.Dir("./public")})

    r.Run(":8080")
}
```

这样访问 `/assets/` 会返回 404，而不是暴露目录结构。

---

## 二、文件响应：在处理函数中返回文件

用 `gin.Context` 上的文件方法。

### c.File —— 直接返回文件内容

```go
r.GET("/avatar/:userID", func(c *gin.Context) {
    userID := c.Param("userID")

    // 业务逻辑：查数据库获取头像路径
    avatarPath, err := getUserAvatar(userID)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "avatar not found"})
        return
    }

    // 返回文件，Content-Type 自动检测
    // 图片会直接在浏览器中显示
    c.File(avatarPath)
})
```

`c.File` 会自动根据文件扩展名设置 `Content-Type`。`.png` 文件返回 `image/png`，`.pdf` 返回 `application/pdf`。浏览器收到后会尝试内联显示（图片直接显示、PDF 在浏览器内打开）。

### c.FileAttachment —— 触发浏览器下载

```go
r.GET("/reports/:id/download", func(c *gin.Context) {
    reportID := c.Param("id")

    // 业务逻辑：查找报表文件
    report, err := getReport(reportID)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "report not found"})
        return
    }

    // 服务器上的文件可能叫 "a1b2c3d4.xlsx"
    // 但用户下载时看到的文件名是 "月度销售报表.xlsx"
    c.FileAttachment(report.FilePath, report.DisplayName+".xlsx")
})
```

`FileAttachment` 和 `File` 的区别只有一个 HTTP 头：`Content-Disposition: attachment; filename="月度销售报表.xlsx"`。这个头告诉浏览器"不要尝试打开，直接下载保存"，并且用指定的文件名。

服务器磁盘上的文件名通常是 UUID 或哈希值（避免冲突和安全问题），而用户看到的下载名称是业务层面有意义的名字。`FileAttachment` 帮你做了这个映射。

### c.FileFromFS —— 安全地从受限目录提供文件

```go
r.GET("/docs/:name", func(c *gin.Context) {
    name := c.Param("name")
    // http.Dir 会把访问限制在 /var/www/documents 目录内
    // 即使 name 是 "../../etc/passwd"，也无法逃逸出这个目录
    c.FileFromFS(name, http.Dir("/var/www/documents"))
})
```

这是三个方法中最安全的一个，因为 `http.Dir` 内部会处理路径清理，防止目录遍历攻击。

---

## 三、路径遍历：最危险的错误

这是文件服务中最常见的安全漏洞，必须单独强调。

```go
// ❌ 极度危险 —— 永远不要这样写
r.GET("/files/:name", func(c *gin.Context) {
    name := c.Param("name")
    c.File(name)  // 攻击者控制了文件路径！
})
```

攻击者只要发送 `GET /files/..%2F..%2F..%2Fetc%2Fpasswd`（`../../../etc/passwd` 的 URL 编码），就能读取服务器上的任意文件——密码文件、环境变量、私钥、数据库配置，什么都能拿到。

正确的做法有三种，按安全程度排序：

```go
// ✅ 方法一（推荐）：用 FileFromFS 限制在特定目录
safeFS := http.Dir("/var/www/public")
r.GET("/files/:name", func(c *gin.Context) {
    c.FileFromFS(c.Param("name"), safeFS)
})

// ✅ 方法二：白名单校验
allowedFiles := map[string]string{
    "report":  "/var/www/files/report.pdf",
    "manual":  "/var/www/files/manual.pdf",
}

r.GET("/files/:key", func(c *gin.Context) {
    path, ok := allowedFiles[c.Param("key")]
    if !ok {
        c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
        return
    }
    c.File(path)
})

// ✅ 方法三：清理路径 + 校验前缀
r.GET("/files/:name", func(c *gin.Context) {
    name := filepath.Base(c.Param("name")) // 只取文件名，去掉所有目录部分
    fullPath := filepath.Join("/var/www/public", name)

    // 再次确认路径没有逃逸
    if !strings.HasPrefix(fullPath, "/var/www/public") {
        c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
        return
    }

    c.File(fullPath)
})
```

`filepath.Base` 会把 `../../etc/passwd` 变成 `passwd`，`filepath.Join` 之后得到 `/var/www/public/passwd`，攻击者无法逃逸出目录。但最简单最可靠的还是直接用 `FileFromFS`，把安全校验交给标准库。

---

## 四、实际业务场景：文件上传 + 下载完整流程

把上传和下载串起来，看一个完整的业务场景：

```go
package main

import (
    "crypto/sha256"
    "fmt"
    "io"
    "net/http"
    "os"
    "path/filepath"
    "strings"
    "sync"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
)

// 文件元信息（实际项目中存数据库）
type FileMeta struct {
    ID           string    `json:"id"`
    OriginalName string    `json:"original_name"`
    StoragePath  string    `json:"-"` // 不暴露给前端
    Size         int64     `json:"size"`
    ContentType  string    `json:"content_type"`
    SHA256       string    `json:"sha256"`
    UploadedAt   time.Time `json:"uploaded_at"`
}

var (
    fileStore   = make(map[string]*FileMeta) // 模拟数据库
    fileStoreMu sync.RWMutex
    uploadDir   = "/var/www/uploads" // 上传文件存储目录
)

func main() {
    os.MkdirAll(uploadDir, 0755)

    r := gin.Default()

    files := r.Group("/api/files")
    {
        files.POST("/upload", uploadFile)
        files.GET("/:id", getFileMeta)
        files.GET("/:id/preview", previewFile)   // 浏览器内预览
        files.GET("/:id/download", downloadFile)  // 触发下载
    }

    r.Run(":8080")
}

func uploadFile(c *gin.Context) {
    file, header, err := c.Request.FormFile("file")
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "no file uploaded"})
        return
    }
    defer file.Close()

    // 限制文件大小（10MB）
    if header.Size > 10<<20 {
        c.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": "file too large"})
        return
    }

    // 校验文件类型（白名单）
    allowedTypes := map[string]bool{
        "image/jpeg":      true,
        "image/png":       true,
        "application/pdf": true,
    }
    contentType := header.Header.Get("Content-Type")
    if !allowedTypes[contentType] {
        c.JSON(http.StatusBadRequest, gin.H{"error": "file type not allowed"})
        return
    }

    // 生成安全的存储文件名（UUID，不用原始文件名）
    fileID := uuid.New().String()
    ext := filepath.Ext(header.Filename)
    storageName := fileID + ext
    storagePath := filepath.Join(uploadDir, storageName)

    // 写入磁盘，同时计算哈希
    dst, err := os.Create(storagePath)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file"})
        return
    }
    defer dst.Close()

    hasher := sha256.New()
    written, err := io.Copy(dst, io.TeeReader(file, hasher))
    if err != nil {
        os.Remove(storagePath) // 清理半成品
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file"})
        return
    }

    meta := &FileMeta{
        ID:           fileID,
        OriginalName: header.Filename,
        StoragePath:  storagePath,
        Size:         written,
        ContentType:  contentType,
        SHA256:       fmt.Sprintf("%x", hasher.Sum(nil)),
        UploadedAt:   time.Now(),
    }

    fileStoreMu.Lock()
    fileStore[fileID] = meta
    fileStoreMu.Unlock()

    c.JSON(http.StatusCreated, meta)
}

func getFileMeta(c *gin.Context) {
    meta := findFile(c)
    if meta == nil {
        return
    }
    c.JSON(http.StatusOK, meta)
}

func previewFile(c *gin.Context) {
    meta := findFile(c)
    if meta == nil {
        return
    }
    // c.File 会内联显示：图片直接展示，PDF 在浏览器中打开
    c.File(meta.StoragePath)
}

func downloadFile(c *gin.Context) {
    meta := findFile(c)
    if meta == nil {
        return
    }
    // c.FileAttachment 触发浏览器下载，使用原始文件名
    c.FileAttachment(meta.StoragePath, meta.OriginalName)
}

func findFile(c *gin.Context) *FileMeta {
    id := c.Param("id")
    fileStoreMu.RLock()
    meta, ok := fileStore[id]
    fileStoreMu.RUnlock()

    if !ok {
        c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
        return nil
    }
    return meta
}
```

这段代码体现了几个关键实践。文件存储用 UUID 命名而不是原始文件名，避免了文件名冲突和路径注入。文件类型用白名单校验，不是靠扩展名判断（扩展名可以伪造）。上传时同步计算 SHA256 哈希，可以用于去重和完整性校验。`StoragePath` 字段标记了 `json:"-"`，不会暴露给前端——前端只需要用文件 ID 来访问。

预览和下载是同一个文件的两种不同返回方式：`/preview` 用 `c.File` 内联显示，`/download` 用 `c.FileAttachment` 触发下载。

---

## 五、大文件下载方案

大文件下载用Gin的话，关键是**流式传输 + 分块读取**，避免把整个文件加载到内存。

### 方案1：基础流式传输（推荐）

```go
package main

import (
    "io"
    "net/http"
    "os"

    "github.com/gin-gonic/gin"
)

func downloadFile(c *gin.Context) {
    filePath := "/path/to/large/file.zip"
    
    file, err := os.Open(filePath)
    if err != nil {
        c.String(http.StatusInternalServerError, "打开文件失败")
        return
    }
    defer file.Close()

    stat, err := file.Stat()
    if err != nil {
        c.String(http.StatusInternalServerError, "获取文件信息失败")
        return
    }

    // 设置响应头
    c.Header("Content-Type", "application/octet-stream")
    c.Header("Content-Disposition", "attachment; filename=\"largefile.zip\"")
    c.Header("Content-Length", fmt.Sprintf("%d", stat.Size()))

    // 流式拷贝，内存占用极小
    io.Copy(c.Writer, file)
}
```

### 方案2：支持断点续传（Range请求）

```go
func downloadWithRange(c *gin.Context) {
    filePath := "/path/to/large/file.zip"
    
    file, err := os.Open(filePath)
    if err != nil {
        c.Status(http.StatusNotFound)
        return
    }
    defer file.Close()

    stat, _ := file.Stat()
    
    // 使用http.ServeContent自动处理Range、Last-Modified等
    http.ServeContent(c.Writer, c.Request, stat.Name(), stat.ModTime(), file)
}
```

`http.ServeContent` 是标准库神器，自动支持：
- 断点续传（Range请求）
- 缓存控制（304 Not Modified）
- MIME类型检测

### 方案3：分片下载（超大文件推荐）

如果文件几十G，还可以考虑**分片生成 + 边生成边传输**：

```go
func streamLargeFile(c *gin.Context) {
    c.Header("Content-Type", "application/octet-stream")
    c.Header("Content-Disposition", "attachment; filename=\"hugefile.bin\"")
    c.Header("Transfer-Encoding", "chunked") // 分块传输

    // 使用bufio减少系统调用
    buf := make([]byte, 64*1024) // 64KB缓冲区
    
    file, _ := os.Open("/path/to/file")
    defer file.Close()

    writer := bufio.NewWriterSize(c.Writer, 256*1024)
    defer writer.Flush()

    io.CopyBuffer(writer, file, buf)
}
```

### 关键要点

| 要点 | 说明 |
|------|------|
| **不要**用 `ioutil.ReadAll` | 会爆内存 |
| **不要**用 `c.File()` 传大文件 | 内部可能全读内存 |
| **要**用 `io.Copy` 流式传输 | 固定内存占用 |
| **建议**加缓冲区 | 减少系统调用次数 |
| **生产环境**用 `http.ServeContent` | 自带断点续传 |

几十G的文件，方案2的 `http.ServeContent` 最省事，方案1最直观可控。


## 五、和 Nginx 配合的最佳实践

在部署篇中我们已经配置了 Nginx 反向代理。对于文件服务，有一个性能优化技巧值得了解：**X-Accel-Redirect**。

思路是这样的：让 Go 应用只做权限校验，校验通过后告诉 Nginx "去哪里拿文件"，实际的文件传输由 Nginx 完成。Nginx 做文件 I/O 的效率远高于 Go（sendfile 系统调用、零拷贝）。

Nginx 配置：

```nginx
# 内部路径，不对外直接暴露
location /internal-files/ {
    internal;                          # 只接受内部重定向，外部直接访问返回 404
    alias /var/www/uploads/;           # 实际文件目录
}
```

Go 处理函数：

```go
r.GET("/api/files/:id/download", func(c *gin.Context) {
    meta := findFile(c)
    if meta == nil {
        return
    }

    // 权限校验（示例：检查当前用户是否有权下载）
    // ...

    // 不自己发文件，让 Nginx 去发
    filename := filepath.Base(meta.StoragePath)
    c.Header("Content-Disposition",
        fmt.Sprintf(`attachment; filename="%s"`, meta.OriginalName))
    c.Header("X-Accel-Redirect", "/internal-files/"+filename)
    c.Status(http.StatusOK)
})
```

Go 应用返回的响应里带了 `X-Accel-Redirect` 头，Nginx 看到这个头后不会把 Go 的响应体转发给客户端，而是去 `/internal-files/` 对应的本地目录读取文件，直接发给客户端。整个过程中文件内容不经过 Go 应用的内存，大文件下载时性能差距非常明显。

---

## 六、静态文件要不要走 Go 应用

这取决于你的架构。有三种常见方案。

第一种，Nginx 直接处理静态文件，Go 只处理 API。这是最常见也是性能最好的方案。在 Nginx 配置中用 `location /static/` 直接指向文件目录，请求根本不经过 Go 应用。部署篇中的 Nginx 配置已经这样做了。

第二种，Go 应用用 `embed.FS` 内嵌静态资源。适合前后端一体的小项目，部署最简单——只有一个二进制文件。但每次前端改动都需要重新编译 Go 应用。

第三种，CDN 托管静态文件。前端打包后的产物上传到 CDN（阿里云 OSS、腾讯云 COS、Cloudflare R2 等），HTML 中引用 CDN 地址。这是大型项目的标准做法，Go 应用完全不参与静态文件服务。

---

## 小结

Gin 的文件服务 API 分两层：路由级别的 `Static`/`StaticFS`/`StaticFile` 用于整个目录或固定文件的映射；处理函数级别的 `c.File`/`c.FileFromFS`/`c.FileAttachment` 用于需要业务逻辑介入的场景。记住永远不要把用户输入直接拼进文件路径。用 `FileFromFS` 或白名单，把路径控制权握在自己手里。

大文件下载，使用 `io.Copy` 和 `http.ServeContent` 可以轻松实现**流式传输 + 分块读取**，避免内存溢出。