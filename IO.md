## 12. 文件与 I/O 操作

Go 的 I/O 体系围绕两个核心接口构建：`io.Reader` 和 `io.Writer`。理解这一点后，所有 I/O 操作都是一脉相通的。

```go
type Reader interface {
    Read(p []byte) (n int, err error)
}
type Writer interface {
    Write(p []byte) (n int, err error)
}
```

> 文件、网络连接、压缩流、加密流……都实现了这两个接口，可以自由组合。

---

### 12.1 打开与关闭文件

```go
// os.Open —— 只读打开
file, err := os.Open("data.txt")
if err != nil {
    log.Fatal(err)
}
defer file.Close()

// os.Create —— 创建/截断文件（可写）
file, err := os.Create("output.txt")

// os.OpenFile —— 完全控制模式和权限
file, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
```

**常用标志位组合**：

| 场景               | 标志位                                       |
|-------------------|---------------------------------------------|
| 只读               | `os.O_RDONLY`                               |
| 只写（覆盖）        | `os.O_WRONLY \| os.O_CREATE \| os.O_TRUNC`  |
| 追加写入            | `os.O_APPEND \| os.O_CREATE \| os.O_WRONLY` |
| 读写               | `os.O_RDWR \| os.O_CREATE`                  |
| 文件必须不存在才创建  | `os.O_CREATE \| os.O_EXCL \| os.O_WRONLY`   |

> ⚠️ **必须用 `defer file.Close()`**，否则文件描述符泄漏。在循环中打开文件时尤其注意，不要在循环里只 defer（要封装成函数或手动 close）。

---

### 12.2 读取文件

#### 一次性读取整个文件（小文件）

```go
// 最简方式：os.ReadFile（Go 1.16+）
data, err := os.ReadFile("config.json")
if err != nil {
    log.Fatal(err)
}
fmt.Println(string(data))

// 等价的手动方式
file, _ := os.Open("config.json")
defer file.Close()
data, _ := io.ReadAll(file)
```

> 适合小文件（< 几十 MB）。整个内容加载到内存，大文件会 OOM。

#### 按行读取（最常用）

```go
func readLines(path string) ([]string, error) {
    file, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    var lines []string
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        lines = append(lines, scanner.Text())
    }
    return lines, scanner.Err()
}
```

> ⚠️ `bufio.Scanner` 默认最大行长为 **64KB**。超长行需要调大缓冲：

```go
scanner := bufio.NewScanner(file)
scanner.Buffer(make([]byte, 0), 10*1024*1024) // 最大 10MB 一行
```

#### 按固定大小块读取（大文件）

```go
func readInChunks(path string) error {
    file, err := os.Open(path)
    if err != nil {
        return err
    }
    defer file.Close()

    buf := make([]byte, 32*1024) // 32KB 缓冲区
    for {
        n, err := file.Read(buf)
        if n > 0 {
            // 处理 buf[:n]
            fmt.Printf("read %d bytes\n", n)
        }
        if err == io.EOF {
            break
        }
        if err != nil {
            return err
        }
    }
    return nil
}
```

#### 使用 bufio.Reader 灵活读取

```go
file, _ := os.Open("data.txt")
defer file.Close()
reader := bufio.NewReader(file)

// 按分隔符读取
for {
    // ReadString 读到分隔符为止（包含分隔符）
    line, err := reader.ReadString('\n')
    if len(line) > 0 {
        line = strings.TrimRight(line, "\r\n")
        fmt.Println(line)
    }
    if err == io.EOF {
        break
    }
    if err != nil {
        log.Fatal(err)
    }
}

// 按字节读取
b, err := reader.ReadByte()

// 窥探但不消费
peeked, err := reader.Peek(5) // 查看前 5 个字节
```

#### 从指定位置读取（随机读取）

```go
file, _ := os.Open("data.bin")
defer file.Close()

// Seek 到指定位置
file.Seek(100, io.SeekStart)   // 从头偏移 100 字节
file.Seek(-50, io.SeekEnd)     // 从末尾倒退 50 字节
file.Seek(10, io.SeekCurrent)  // 从当前位置前进 10 字节

buf := make([]byte, 64)
n, _ := file.Read(buf)

// ReadAt 不改变文件偏移量（并发安全）
n, err := file.ReadAt(buf, 1024) // 从偏移 1024 处读取
```

#### io.SectionReader —— 读取文件的一个片段

```go
file, _ := os.Open("large.bin")
defer file.Close()

// 从偏移 1000 处读取 500 字节
section := io.NewSectionReader(file, 1000, 500)
data, _ := io.ReadAll(section)
```

---

### 12.3 写入文件

#### 一次性写入（小内容）

```go
// os.WriteFile（Go 1.16+）
err := os.WriteFile("output.txt", []byte("Hello, Go!\n"), 0644)
if err != nil {
    log.Fatal(err)
}
```

#### 使用 bufio.Writer 高效写入

裸 `file.Write()` 每次都触发系统调用。`bufio.Writer` 在内存中累积数据，批量写入。

```go
func writeLines(path string, lines []string) error {
    file, err := os.Create(path)
    if err != nil {
        return err
    }
    defer file.Close()

    writer := bufio.NewWriter(file)
    for _, line := range lines {
        // WriteString 比 Write([]byte(...)) 更高效
        _, err := writer.WriteString(line + "\n")
        if err != nil {
            return err
        }
    }
    return writer.Flush() // 关键！将缓冲区剩余数据刷到磁盘
}
```

> ⚠️ **忘记 `Flush()` 是常见 Bug**。数据还在缓冲区里，文件可能为空或不完整。

#### 使用 fmt.Fprintf 格式化写入

```go
file, _ := os.Create("report.txt")
defer file.Close()

writer := bufio.NewWriter(file)
defer writer.Flush()

fmt.Fprintf(writer, "User: %s\n", "Alice")
fmt.Fprintf(writer, "Score: %d\n", 95)
fmt.Fprintf(writer, "Rate: %.2f%%\n", 87.5)
```

#### 追加写入

```go
func appendToFile(path, content string) error {
    file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }
    defer file.Close()

    _, err = file.WriteString(content)
    return err
}

appendToFile("app.log", "[2024-03-15 10:30:00] server started\n")
```

#### 在指定位置写入

```go
file, _ := os.OpenFile("data.bin", os.O_RDWR, 0644)
defer file.Close()

file.Seek(100, io.SeekStart)
file.Write([]byte{0xFF, 0xFE})

// WriteAt 不改变偏移量
file.WriteAt([]byte("hello"), 200)
```

---

### 12.4 大文件处理实战

#### 示例1：大文件复制（流式，内存恒定）

```go
func copyFile(src, dst string) (int64, error) {
    sourceFile, err := os.Open(src)
    if err != nil {
        return 0, err
    }
    defer sourceFile.Close()

    destFile, err := os.Create(dst)
    if err != nil {
        return 0, err
    }
    defer destFile.Close()

    // io.Copy 内部使用 32KB 缓冲区
    // 无论文件多大，内存占用恒定
    nBytes, err := io.Copy(destFile, sourceFile)
    if err != nil {
        return 0, err
    }

    // 确保数据落盘
    err = destFile.Sync()
    return nBytes, err
}

// 指定缓冲区大小
func copyFileBuffered(src, dst string, bufSize int) (int64, error) {
    sourceFile, _ := os.Open(src)
    defer sourceFile.Close()
    destFile, _ := os.Create(dst)
    defer destFile.Close()

    buf := make([]byte, bufSize)
    return io.CopyBuffer(destFile, sourceFile, buf)
}

// 使用：1MB 缓冲区复制大文件
copyFileBuffered("10GB.dat", "backup.dat", 1024*1024)
```

#### 示例2：大文件逐行处理（如日志分析）

处理一个 10GB 的 Nginx 日志文件，统计每个 IP 的请求次数：

```go
func analyzeLog(path string) (map[string]int, error) {
    file, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    ipCount := make(map[string]int)
    scanner := bufio.NewScanner(file)
    // Nginx 日志一行通常不超过 4KB，默认 64KB 足够
    lineNum := 0

    for scanner.Scan() {
        lineNum++
        line := scanner.Text()

        // 提取第一个字段（IP 地址）
        // 日志格式：192.168.1.1 - - [15/Mar/2024:10:30:00 +0800] "GET / HTTP/1.1" 200 612
        if idx := strings.IndexByte(line, ' '); idx > 0 {
            ip := line[:idx]
            ipCount[ip]++
        }

        // 每百万行打印一次进度
        if lineNum%1_000_000 == 0 {
            fmt.Printf("processed %d million lines...\n", lineNum/1_000_000)
        }
    }

    if err := scanner.Err(); err != nil {
        return nil, fmt.Errorf("scan error at line %d: %w", lineNum, err)
    }

    return ipCount, nil
}
```

#### 示例3：大文件分块写入（生成测试数据）

```go
func generateLargeFile(path string, sizeGB int) error {
    file, err := os.Create(path)
    if err != nil {
        return err
    }
    defer file.Close()

    writer := bufio.NewWriterSize(file, 4*1024*1024) // 4MB 缓冲
    defer writer.Flush()

    targetBytes := int64(sizeGB) * 1024 * 1024 * 1024
    var written int64

    for written < targetBytes {
        line := fmt.Sprintf("%s [INFO] request_id=%d user_agent=Mozilla/5.0 path=/api/v1/data\n",
            time.Now().Format(time.RFC3339), rand.Intn(1_000_000))
        n, err := writer.WriteString(line)
        if err != nil {
            return err
        }
        written += int64(n)
    }

    fmt.Printf("wrote %d bytes (%.2f GB)\n", written, float64(written)/(1024*1024*1024))
    return nil
}
```

#### 示例4：带进度条的大文件复制

```go
// 自定义 Reader 包装器，跟踪读取进度
type ProgressReader struct {
    reader     io.Reader
    total      int64
    current    int64
    lastReport float64
}

func (pr *ProgressReader) Read(p []byte) (int, error) {
    n, err := pr.reader.Read(p)
    pr.current += int64(n)

    if pr.total > 0 {
        percent := float64(pr.current) / float64(pr.total) * 100
        // 每增长 1% 才打印，避免刷屏
        if percent-pr.lastReport >= 1 {
            fmt.Printf("\r  progress: %.1f%% (%s / %s)",
                percent, formatBytes(pr.current), formatBytes(pr.total))
            pr.lastReport = percent
        }
    }
    return n, err
}

func formatBytes(b int64) string {
    const unit = 1024
    if b < unit {
        return fmt.Sprintf("%d B", b)
    }
    div, exp := int64(unit), 0
    for n := b / unit; n >= unit; n /= unit {
        div *= unit
        exp++
    }
    return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}

func copyWithProgress(src, dst string) error {
    sourceFile, err := os.Open(src)
    if err != nil {
        return err
    }
    defer sourceFile.Close()

    info, _ := sourceFile.Stat()

    destFile, err := os.Create(dst)
    if err != nil {
        return err
    }
    defer destFile.Close()

    pr := &ProgressReader{
        reader: sourceFile,
        total:  info.Size(),
    }

    _, err = io.Copy(destFile, pr)
    fmt.Println() // 换行
    return err
}
```

---

### 12.5 临时文件与临时目录

```go
// 创建临时文件
tmpFile, err := os.CreateTemp("", "myapp-*.txt")
// 第一个参数为空表示使用系统默认临时目录
// 第二个参数中 * 会被替换为随机字符串
if err != nil {
    log.Fatal(err)
}
defer os.Remove(tmpFile.Name()) // 用完删除
defer tmpFile.Close()

fmt.Println(tmpFile.Name()) // 如 /tmp/myapp-384729481.txt
tmpFile.WriteString("temporary data")

// 创建临时目录
tmpDir, err := os.MkdirTemp("", "myapp-*")
if err != nil {
    log.Fatal(err)
}
defer os.RemoveAll(tmpDir) // 递归删除

fmt.Println(tmpDir) // 如 /tmp/myapp-291837465
```

---

### 12.6 文件信息与权限

```go
// 获取文件信息
info, err := os.Stat("data.txt")
if err != nil {
    if os.IsNotExist(err) {
        fmt.Println("file does not exist")
    }
    log.Fatal(err)
}

fmt.Println(info.Name())    // "data.txt"
fmt.Println(info.Size())    // 字节数
fmt.Println(info.Mode())    // "-rw-r--r--"
fmt.Println(info.ModTime()) // 修改时间
fmt.Println(info.IsDir())   // 是否是目录

// Lstat 不跟踪符号链接（Stat 会跟踪）
info, err = os.Lstat("symlink.txt")

// 修改权限
os.Chmod("script.sh", 0755)

// 修改所有者（需要 root）
os.Chown("data.txt", uid, gid)

// 修改时间戳
os.Chtimes("data.txt", time.Now(), time.Now())

// 判断文件类型
if info.Mode().IsRegular() { /* 普通文件 */ }
if info.Mode().IsDir()     { /* 目录 */ }
if info.Mode()&os.ModeSymlink != 0 { /* 符号链接 */ }
```

---

### 12.7 目录操作

```go
// 创建目录
os.Mkdir("logs", 0755)          // 单层
os.MkdirAll("a/b/c/d", 0755)   // 递归创建

// 读取目录内容
entries, err := os.ReadDir(".")
if err != nil {
    log.Fatal(err)
}
for _, entry := range entries {
    info, _ := entry.Info()
    fmt.Printf("%-30s %10d %s\n", entry.Name(), info.Size(), info.Mode())
}

// 递归遍历目录树
err := filepath.Walk("/var/log", func(path string, info os.FileInfo, err error) error {
    if err != nil {
        return err
    }
    if info.IsDir() && info.Name() == ".git" {
        return filepath.SkipDir // 跳过整个目录
    }
    fmt.Println(path)
    return nil
})

// Go 1.16+ filepath.WalkDir 更高效（不调用 Stat）
err := filepath.WalkDir("/var/log", func(path string, d fs.DirEntry, err error) error {
    if err != nil {
        return err
    }
    fmt.Println(path, d.IsDir())
    return nil
})

// 实际场景：查找所有 .go 文件
func findGoFiles(root string) ([]string, error) {
    var files []string
    err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
        if err != nil {
            return err
        }
        if !d.IsDir() && filepath.Ext(path) == ".go" {
            files = append(files, path)
        }
        return nil
    })
    return files, err
}

// Glob 模式匹配
matches, _ := filepath.Glob("/var/log/*.log")
for _, m := range matches {
    fmt.Println(m)
}
```

---

### 12.8 filepath 包详解

```go
import "path/filepath"

// 拼接路径（自动处理分隔符）
p := filepath.Join("home", "user", "docs", "file.txt")
// Linux: "home/user/docs/file.txt"
// Windows: "home\\user\\docs\\file.txt"

// 提取各部分
filepath.Base("/a/b/c.txt")   // "c.txt"
filepath.Dir("/a/b/c.txt")    // "/a/b"
filepath.Ext("/a/b/c.tar.gz") // ".gz"

// 去掉扩展名
name := strings.TrimSuffix("report.csv", filepath.Ext("report.csv")) // "report"

// 绝对路径
abs, _ := filepath.Abs("./data.txt") // "/home/user/project/data.txt"

// 相对路径
rel, _ := filepath.Rel("/home/user", "/home/user/project/main.go")
// "project/main.go"

// 清理路径
filepath.Clean("/a/b/../c/./d") // "/a/c/d"

// 分割路径和文件名
dir, file := filepath.Split("/home/user/data.txt")
// dir = "/home/user/", file = "data.txt"

// 匹配模式
matched, _ := filepath.Match("*.txt", "readme.txt") // true
matched, _ = filepath.Match("log_202?-*", "log_2024-03") // true
```

---

### 12.9 io 包核心工具

```go
import "io"

// io.Copy：流式复制（已在 12.4 演示）
io.Copy(dst, src)             // 使用默认 32KB 缓冲
io.CopyN(dst, src, 1024)     // 只复制 1024 字节
io.CopyBuffer(dst, src, buf) // 使用自定义缓冲

// io.ReadAll：读取全部内容
data, err := io.ReadAll(reader)

// io.ReadFull：确保读满 buf
buf := make([]byte, 100)
_, err := io.ReadFull(reader, buf) // 不够 100 字节会报错

// io.LimitReader：限制最大读取量
limited := io.LimitReader(file, 1024*1024) // 最多读 1MB
data, _ := io.ReadAll(limited)

// io.MultiReader：串联多个 Reader
header := strings.NewReader("HEADER\n")
body, _ := os.Open("body.txt")
footer := strings.NewReader("\nFOOTER")
combined := io.MultiReader(header, body, footer)
io.Copy(os.Stdout, combined) // 依次输出三部分

// io.MultiWriter：同时写入多个目标
logFile, _ := os.Create("app.log")
multi := io.MultiWriter(os.Stdout, logFile)
fmt.Fprintln(multi, "this goes to both console and file")

// io.TeeReader：读取的同时复制一份
var buf bytes.Buffer
tee := io.TeeReader(resp.Body, &buf)
io.ReadAll(tee)
// 此时 buf 中也有一份 resp.Body 的完整内容

// io.Pipe：同步的内存管道
pr, pw := io.Pipe()
go func() {
    defer pw.Close()
    json.NewEncoder(pw).Encode(data) // 写入端
}()
json.NewDecoder(pr).Decode(&result)  // 读取端

// io.NopCloser：给 Reader 加一个空的 Close 方法
rc := io.NopCloser(strings.NewReader("hello"))
// rc 实现了 io.ReadCloser
```

---

### 12.10 bufio 高级用法

```go
import "bufio"

// 自定义缓冲区大小
reader := bufio.NewReaderSize(file, 1*1024*1024) // 1MB 读缓冲
writer := bufio.NewWriterSize(file, 1*1024*1024) // 1MB 写缓冲

// bufio.Scanner 自定义分割函数
scanner := bufio.NewScanner(file)

// 内置分割函数
scanner.Split(bufio.ScanLines)  // 默认：按行
scanner.Split(bufio.ScanWords)  // 按单词
scanner.Split(bufio.ScanBytes)  // 按字节
scanner.Split(bufio.ScanRunes)  // 按 UTF-8 字符

// 自定义分割：按 \0 分割（处理 null-terminated 数据）
scanner.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
    if atEOF && len(data) == 0 {
        return 0, nil, nil
    }
    if i := bytes.IndexByte(data, 0); i >= 0 {
        return i + 1, data[:i], nil
    }
    if atEOF {
        return len(data), data, nil
    }
    return 0, nil, nil // 需要更多数据
})

// ReadWriter：组合 Reader 和 Writer
rw := bufio.NewReadWriter(
    bufio.NewReader(conn),
    bufio.NewWriter(conn),
)
```

---

### 12.11 实战：CSV 文件处理

```go
import "encoding/csv"

// 读取 CSV
func readCSV(path string) ([][]string, error) {
    file, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    reader := csv.NewReader(file)
    reader.Comma = ','        // 分隔符，默认逗号
    reader.Comment = '#'      // 注释行前缀
    reader.TrimLeadingSpace = true

    // 一次读取全部
    records, err := reader.ReadAll()
    return records, err
}

// 逐行读取大 CSV（省内存）
func processLargeCSV(path string) error {
    file, err := os.Open(path)
    if err != nil {
        return err
    }
    defer file.Close()

    reader := csv.NewReader(file)

    // 读取表头
    header, err := reader.Read()
    if err != nil {
        return err
    }
    fmt.Println("Columns:", header)

    lineNum := 1
    for {
        record, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return fmt.Errorf("line %d: %w", lineNum, err)
        }
        lineNum++

        // 处理每一行
        // record[0], record[1], ...
        _ = record
    }
    return nil
}

// 写入 CSV
func writeCSV(path string, data [][]string) error {
    file, err := os.Create(path)
    if err != nil {
        return err
    }
    defer file.Close()

    // 写入 UTF-8 BOM（让 Excel 正确识别中文）
    file.Write([]byte{0xEF, 0xBB, 0xBF})

    writer := csv.NewWriter(file)
    defer writer.Flush()

    // 写表头
    writer.Write([]string{"姓名", "年龄", "邮箱"})

    // 写数据
    for _, row := range data {
        if err := writer.Write(row); err != nil {
            return err
        }
    }

    return writer.Error()
}
```

---

### 12.12 实战：JSON 文件流式处理

当 JSON 文件非常大（如包含百万条记录的数组），不适合 `json.Unmarshal` 全部加载：

```go
// 流式解析大 JSON 数组
// 文件内容：[{"id":1,"name":"Alice"},{"id":2,"name":"Bob"},...]
func processLargeJSON(path string) error {
    file, err := os.Open(path)
    if err != nil {
        return err
    }
    defer file.Close()

    decoder := json.NewDecoder(file)

    // 读取 '[' 开始符号
    token, err := decoder.Token()
    if err != nil {
        return err
    }
    if delim, ok := token.(json.Delim); !ok || delim != '[' {
        return fmt.Errorf("expected '[', got %v", token)
    }

    type User struct {
        ID   int    `json:"id"`
        Name string `json:"name"`
    }

    count := 0
    for decoder.More() {
        var user User
        if err := decoder.Decode(&user); err != nil {
            return fmt.Errorf("record %d: %w", count, err)
        }
        count++
        // 处理每条记录...
    }

    fmt.Printf("processed %d records\n", count)
    return nil
}

// 流式写入大 JSON 数组
func writeLargeJSON(path string, users <-chan User) error {
    file, err := os.Create(path)
    if err != nil {
        return err
    }
    defer file.Close()

    writer := bufio.NewWriter(file)
    defer writer.Flush()

    encoder := json.NewEncoder(writer)
    encoder.SetIndent("", "  ")

    writer.WriteString("[\n")
    first := true
    for user := range users {
        if !first {
            writer.WriteString(",\n")
        }
        first = false
        if err := encoder.Encode(user); err != nil {
            return err
        }
    }
    writer.WriteString("]\n")
    return nil
}
```

---

### 12.13 实战：安全写入（原子写入）

直接写入文件，如果中途断电或程序崩溃，可能导致文件损坏。安全做法是**先写临时文件，再原子重命名**：

```go
func atomicWriteFile(path string, data []byte, perm os.FileMode) error {
    dir := filepath.Dir(path)

    // 在同一目录创建临时文件（确保在同一文件系统，rename 才是原子的）
    tmpFile, err := os.CreateTemp(dir, ".tmp-*")
    if err != nil {
        return err
    }
    tmpPath := tmpFile.Name()

    // 出错时清理临时文件
    defer func() {
        if err != nil {
            os.Remove(tmpPath)
        }
    }()

    // 写入数据
    if _, err = tmpFile.Write(data); err != nil {
        tmpFile.Close()
        return err
    }

    // 确保数据落盘
    if err = tmpFile.Sync(); err != nil {
        tmpFile.Close()
        return err
    }

    if err = tmpFile.Close(); err != nil {
        return err
    }

    // 设置权限
    if err = os.Chmod(tmpPath, perm); err != nil {
        return err
    }

    // 原子重命名
    return os.Rename(tmpPath, path)
}

// 使用
config := []byte(`{"port": 8080, "debug": false}`)
atomicWriteFile("config.json", config, 0644)
```

---

### 12.14 实战：监控文件变化

简单的轮询方案（生产环境推荐使用 `fsnotify` 库）：

```go
func watchFile(path string, interval time.Duration, callback func()) error {
    initialStat, err := os.Stat(path)
    if err != nil {
        return err
    }

    for {
        time.Sleep(interval)

        stat, err := os.Stat(path)
        if err != nil {
            return err
        }

        if stat.Size() != initialStat.Size() || stat.ModTime() != initialStat.ModTime() {
            callback()
            initialStat = stat
        }
    }
}

// 使用
go watchFile("config.yaml", 2*time.Second, func() {
    fmt.Println("config file changed! reloading...")
    // 重新加载配置...
})
```

---

### 12.15 内存中的 I/O

有时候不需要真正的文件，可以在内存中使用相同的 Reader/Writer 接口：

```go
import (
    "bytes"
    "strings"
)

// strings.Reader —— 从字符串读取
reader := strings.NewReader("Hello, Go!")
data, _ := io.ReadAll(reader) // []byte("Hello, Go!")

// bytes.Buffer —— 可读可写的内存缓冲区
var buf bytes.Buffer
buf.WriteString("hello ")
buf.WriteString("world")
buf.WriteByte('!')
fmt.Println(buf.String()) // "hello world!"
fmt.Println(buf.Len())    // 12

// 从 Buffer 读取
line, _ := buf.ReadString(' ') // "hello "

// bytes.Reader —— 只读，支持 Seek
r := bytes.NewReader([]byte("Hello"))
r.Seek(2, io.SeekStart)
b, _ := r.ReadByte() // 'l'

// 常见场景：构造 HTTP 请求体
body := bytes.NewBufferString(`{"name":"Alice"}`)
req, _ := http.NewRequest("POST", url, body)
```

---

### 12.16 总结：选择合适的读写方式

| 场景 | 推荐方式 | 内存占用 |
|-----|---------|---------|
| 小文件（< 10MB）整体读取 | `os.ReadFile` | 文件大小 |
| 小文件整体写入 | `os.WriteFile` | 文件大小 |
| 按行处理 | `bufio.Scanner` | 单行大小 |
| 大文件流式复制 | `io.Copy` | 32KB（固定） |
| 大文件块处理 | `file.Read` + 固定 buf | buf 大小 |
| 高频写入（日志等） | `bufio.Writer` | 缓冲大小 |
| 大 JSON 流式解析 | `json.NewDecoder` | 单条记录 |
| 大 CSV 逐行处理 | `csv.Reader.Read` | 单行大小 |
| 构造内存数据 | `bytes.Buffer` | 数据大小 |
| 需要原子写入 | 临时文件 + `os.Rename` | 文件大小 |


这个扩展版本覆盖了文件 I/O 的各个方面，从基础操作到大文件处理、流式 JSON/CSV、原子写入、进度追踪等真实场景。你可以直接替换原文中的第 12 节。需要我对其他章节也做类似的深化吗？