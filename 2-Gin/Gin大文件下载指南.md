# Gin 大文件下载指南

大文件下载用Gin的话，关键是**流式传输 + 分块读取**，避免把整个文件加载到内存。给你几个方案：

## 方案1：基础流式传输（推荐）

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

## 方案2：支持断点续传（Range请求）

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

## 方案3：分片下载（超大文件推荐）

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

## 关键要点

| 要点 | 说明 |
|------|------|
| **不要**用 `ioutil.ReadAll` | 会爆内存 |
| **不要**用 `c.File()` 传大文件 | 内部可能全读内存 |
| **要**用 `io.Copy` 流式传输 | 固定内存占用 |
| **建议**加缓冲区 | 减少系统调用次数 |
| **生产环境**用 `http.ServeContent` | 自带断点续传 |

几十G的文件，方案2的 `http.ServeContent` 最省事，方案1最直观可控。
