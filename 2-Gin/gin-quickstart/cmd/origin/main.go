package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func main() {
	// 启动测试服务器
	go startServer()

	// 等待服务器启动
	select {}

	// 运行客户端测试（实际使用时放开）
	// runClientTests()
}

// ============ 服务器端 ============

func startServer() {
	http.HandleFunc("/query", handleQuery)
	http.HandleFunc("/form-urlencoded", handleFormURLEncoded)
	http.HandleFunc("/json", handleJSON)
	http.HandleFunc("/multipart", handleMultipart)
	http.HandleFunc("/upload", handleFileUpload)

	fmt.Println("Server starting on :8080")
	http.ListenAndServe(":8080", nil)
}

// 1. Query 参数处理
func handleQuery(w http.ResponseWriter, r *http.Request) {
	// 从 URL 解析
	name := r.URL.Query().Get("name")
	age := r.URL.Query().Get("age")

	fmt.Fprintf(w, "Query参数: name=%s, age=%s\n", name, age)
}

// 2. form-urlencoded 处理
func handleFormURLEncoded(w http.ResponseWriter, r *http.Request) {
	// 解析表单数据 (Content-Type: application/x-www-form-urlencoded)
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	name := r.FormValue("name")
	age := r.FormValue("age")

	fmt.Fprintf(w, "Form数据: name=%s, age=%s\n", name, age)
}

// 3. JSON 处理
func handleJSON(w http.ResponseWriter, r *http.Request) {
	// 读取 body
	body, _ := io.ReadAll(r.Body)
	defer r.Body.Close()

	// 解析 JSON
	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	fmt.Fprintf(w, "JSON数据: %+v\n", data)
}

// 4. multipart/form-data 处理（纯表单，无文件）
func handleMultipart(w http.ResponseWriter, r *http.Request) {
	// 解析 multipart (内存限制 32MB)
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	name := r.FormValue("name")
	desc := r.FormValue("description")

	fmt.Fprintf(w, "Multipart表单: name=%s, description=%s\n", name, desc)
}

// 5. 文件上传处理
func handleFileUpload(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	// 获取普通字段
	user := r.FormValue("user")

	// 获取文件
	file, header, err := r.FormFile("avatar")
	if err != nil {
		http.Error(w, "获取文件失败: "+err.Error(), 400)
		return
	}
	defer file.Close()

	// 读取文件信息
	fileSize := header.Size
	fileName := header.Filename
	contentType := header.Header.Get("Content-Type")

	// 保存文件（示例）
	// dst, _ := os.Create("./uploads/" + fileName)
	// io.Copy(dst, file)

	fmt.Fprintf(w, "文件上传成功!\n用户: %s\n文件名: %s\n大小: %d bytes\n类型: %s\n",
		user, fileName, fileSize, contentType)
}

// ============ 客户端示例 ============

func runClientTests() {
	baseURL := "http://localhost:8080"

	// 1. Query 参数
	testQueryParams(baseURL)

	// 2. form-urlencoded
	testFormURLEncoded(baseURL)

	// 3. JSON
	testJSON(baseURL)

	// 4. multipart 表单
	testMultipartForm(baseURL)

	// 5. 文件上传
	testFileUpload(baseURL)
}

// 1. Query 参数请求
func testQueryParams(baseURL string) {
	// 方式一：手动拼接
	resp, _ := http.Get(baseURL + "/query?name=张三&age=25")
	body, _ := io.ReadAll(resp.Body)
	fmt.Println("Query:", string(body))

	// 方式二：用 url.Values 安全编码
	params := url.Values{}
	params.Set("name", "张三") // 自动处理中文编码
	params.Set("age", "25")
	resp2, _ := http.Get(baseURL + "/query?" + params.Encode())
	body2, _ := io.ReadAll(resp2.Body)
	fmt.Println("Query(编码):", string(body2))
}

// 2. form-urlencoded 请求
func testFormURLEncoded(baseURL string) {
	// 构造表单数据
	data := url.Values{}
	data.Set("name", "李四")
	data.Set("age", "30")

	// 发送 POST
	resp, _ := http.Post(
		baseURL+"/form-urlencoded",
		"application/x-www-form-urlencoded",
		strings.NewReader(data.Encode()),
	)
	body, _ := io.ReadAll(resp.Body)
	fmt.Println("Form:", string(body))
}

// 3. JSON 请求
func testJSON(baseURL string) {
	payload := map[string]interface{}{
		"name": "王五",
		"age":  28,
		"tags": []string{"go", "backend"},
	}
	jsonData, _ := json.Marshal(payload)

	resp, _ := http.Post(
		baseURL+"/json",
		"application/json",
		bytes.NewReader(jsonData),
	)
	body, _ := io.ReadAll(resp.Body)
	fmt.Println("JSON:", string(body))
}

// 4. multipart 表单（无文件）
func testMultipartForm(baseURL string) {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// 添加字段
	writer.WriteField("name", "赵六")
	writer.WriteField("description", "这是一个描述")

	writer.Close() // 必须关闭以写入 boundary

	resp, _ := http.Post(
		baseURL+"/multipart",
		writer.FormDataContentType(), // 自动包含 boundary
		&buf,
	)
	body, _ := io.ReadAll(resp.Body)
	fmt.Println("Multipart:", string(body))
}

// 5. 文件上传 ⭐ 重点
func testFileUpload(baseURL string) {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// 添加普通字段
	writer.WriteField("user", "孙七")

	// 添加文件
	fileWriter, _ := writer.CreateFormFile("avatar", "test.jpg")

	// 模拟文件内容（实际用 os.Open 打开真实文件）
	fileContent := []byte("这里是图片的二进制数据...")
	fileWriter.Write(fileContent)

	// 或者打开真实文件：
	// file, _ := os.Open("photo.jpg")
	// io.Copy(fileWriter, file)

	writer.Close()

	resp, _ := http.Post(
		baseURL+"/upload",
		writer.FormDataContentType(),
		&buf,
	)
	body, _ := io.ReadAll(resp.Body)
	fmt.Println("文件上传:", string(body))
}

// 真实文件上传示例
func uploadRealFile(baseURL, filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// 创建文件字段
	fileWriter, _ := writer.CreateFormFile("avatar", filePath)

	// 复制文件内容
	io.Copy(fileWriter, file)

	// 添加其他字段
	writer.WriteField("user", "current_user")

	writer.Close()

	resp, _ := http.Post(
		baseURL+"/upload",
		writer.FormDataContentType(),
		&buf,
	)
	body, _ := io.ReadAll(resp.Body)
	fmt.Println(string(body))
}