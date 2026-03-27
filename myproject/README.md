# MyProject

Go 全栈项目模板，包含 HTTP API、gRPC、事件总线、定时任务、消息队列。

## 快速开始

```bash
# 克隆后进入项目
cd myproject

# 下载依赖
go mod tidy

# 启动 HTTP 服务（内存存储，开箱即用）
make run

# 另开终端测试
curl http://localhost:8080/healthz

# 注册用户
curl -X POST http://localhost:8080/api/v1/register \
  -H "Content-Type: application/json" \
  -d '{"username":"alice","email":"alice@example.com","password":"password123"}'

# 登录
curl -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"email":"alice@example.com","password":"password123"}'

# 用返回的 token 访问受保护接口
curl http://localhost:8080/api/v1/users/me \
  -H "Authorization: Bearer <your-token>"

# 创建订单
curl -X POST http://localhost:8080/api/v1/orders \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your-token>" \
  -d '{"product":"Go 编程指南","amount":99.00}'
```

## gRPC 服务

```bash
# 1. 安装 protoc + Go 插件
# 2. 生成代码
make proto

# 3. 启动 gRPC 服务
make run-grpc
```

## 其他命令

```bash
make test           # 运行测试
make lint           # 代码检查
make build          # 编译所有二进制
make docker         # Docker Compose 启动
make seed           # 填充测试数据
make run-worker     # 启动后台 Worker
make run-scheduler  # 启动定时任务
```

## 技术栈

- **HTTP**: gorilla/mux
- **gRPC**: google.golang.org/grpc + protobuf
- **认证**: JWT (golang-jwt/jwt)
- **配置**: YAML + 环境变量覆盖
- **日志**: log/slog (标准库)
- **存储**: 内存实现 (可替换为 MySQL/PostgreSQL)
- **缓存**: 内存实现 (可替换为 Redis)
- **消息队列**: 内存实现 (可替换为 Kafka/RabbitMQ/NATS)

---

## 完整调用链路示意

```
用户请求 POST /api/v1/register
    │
    ▼
router.go (gorilla/mux 匹配路由)
    │
    ▼
middleware: Recovery → Logging → CORS
    │
    ▼
handler/user.go  Register()
    │  解析 JSON → dto.CreateUserRequest
    ▼
service/user.go  Register()
    │  校验输入 → 哈希密码 → 构造 model.User
    ▼
repository/user.go  Create()
    │  写入内存 map（生产环境写入 MySQL）
    ▼
event/bus.go  Publish(UserCreatedEvent)
    │
    ▼
event/handler/user_created.go  onUserCreated()
    │
    ▼
notify/webhook.go  Send()  → 控制台输出（生产环境发邮件/短信）
    │
    ▼
返回 201 + dto.UserResponse JSON
```

`go mod tidy` → `make run` 即可启动 HTTP 服务并通过 curl 测试全部流程。