# Gin 框架进阶系列（十）：项目部署——Docker 容器化 + Nginx 反向代理

---

## 为什么容器化？

"在我电脑上是好的"这句话之所以成为程序员经典名言，就是因为开发环境和生产环境之间存在无数差异：操作系统版本、依赖库版本、环境变量、文件路径、网络配置……任何一项不一致都可能让程序出问题。

Docker 解决的核心问题就是这个——把应用和它的运行环境打包成一个标准化的镜像，在哪里运行都一样。Nginx 解决的则是另一个问题——你的 Go 应用不应该直接暴露给公网，它前面需要一层专业的反向代理来处理 TLS 终止、静态文件服务、负载均衡、请求缓冲等工作。

这篇文章从零开始，把一个 Gin 项目从本机搬到生产服务器上，完整走一遍。

---

## 一、部署前的项目结构

简化版项目结构大致这样：

```
my-gin-project/
├── cmd/
│   └── server/
│       └── main.go              # 入口文件
├── internal/
│   ├── handler/                 # 请求处理器
│   ├── middleware/               # 中间件
│   ├── model/                   # 数据模型
│   ├── service/                 # 业务逻辑
│   ├── repository/              # 数据访问层
│   ├── database/                # 数据库初始化
│   ├── router/                  # 路由注册
│   └── ...          
├── pkg/                         # 公共工具包
│   ├── cache/
│   ├── pool/
│   └── response/
├── configs/                     # 配置文件
│   ├── config.yaml
│   ├── config.example.yaml
│   ├── config.dev.yaml
│   └── config.prod.yaml
├── deployments/                 # 部署相关文件
│   ├── docker/
│   │   ├── Dockerfile
│   │   └── .dockerignore
│   ├── nginx/
│   │   ├── nginx.conf
│   │   └── conf.d/
│   │       └── app.conf
│   └── docker-compose.yaml
├── scripts/                     # 脚本
│   └── wait-for-it.sh
├── go.mod
├── go.sum
├── Makefile
└── .env.example
```

`deployments` 目录集中放所有部署相关的文件。有些人喜欢把 Dockerfile 放在项目根目录，也没问题，但当部署文件越来越多（docker-compose、nginx 配置、k8s manifest），集中放更整洁。

---

## 二、Dockerfile：多阶段构建

### 为什么需要多阶段构建

Go 的编译工具链大约 1GB，加上项目依赖的源码，构建时的镜像可能有 1.5GB 以上。但编译出来的二进制文件通常只有 10-30MB。如果把整个编译环境都打包进最终镜像，就是白白浪费 1.4GB 的空间——不仅拉取慢、推送慢，攻击面也更大。

多阶段构建的思路很简单：第一个阶段用完整的 Go 环境编译代码，第二个阶段只把编译出来的二进制文件和必要的运行时文件复制到一个极小的基础镜像里。

### 完整的 Dockerfile

```dockerfile
# deployments/docker/Dockerfile

# ============================================
# 阶段一：编译
# ============================================
FROM golang:1.26-alpine AS builder

# 安装编译时可能需要的系统依赖
# gcc 和 musl-dev 是 CGO 需要的（如果你用了 SQLite 或某些需要 CGO 的库）
# 如果完全不需要 CGO，可以去掉这行
RUN apk add --no-cache gcc musl-dev

WORKDIR /build

# 先复制 go.mod 和 go.sum，单独下载依赖
# 这一步利用了 Docker 的层缓存机制：
# 只要 go.mod 和 go.sum 没变，这一层就不会重新执行
# 即使你改了业务代码，依赖下载这一步也会命中缓存
COPY go.mod go.sum ./
RUN go mod download && go mod verify 

# 再复制全部源码
COPY . .

# 编译
# CGO_ENABLED=0：禁用 CGO，编译出纯静态链接的二进制文件
#                不依赖任何系统动态库，可以在 scratch/distroless 上运行
# -trimpath：去掉二进制文件中的本地路径信息（安全考虑）
# -ldflags：
#   -s 去掉符号表
#   -w 去掉 DWARF 调试信息
#   -X 注入版本信息（构建时传入）
ARG APP_VERSION=dev
RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags="-s -w -X main.version=${APP_VERSION}" \
    -o /build/server \
    ./cmd/server/

# ============================================
# 阶段二：运行
# ============================================
FROM alpine:3.20

# 时区和证书
# ca-certificates：如果你的应用需要调用外部 HTTPS 服务
# tzdata：时区数据，很多日志和定时任务依赖正确的时区
RUN apk add --no-cache ca-certificates tzdata

# 创建非 root 用户
# 永远不要用 root 运行生产服务
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

# 从构建阶段复制编译好的二进制文件
COPY --from=builder /build/server .

# 复制配置文件（如果需要）
COPY configs/config.yaml ./configs/

# 切换到非 root 用户
USER appuser

# 暴露端口（文档性质，不实际映射端口）
EXPOSE 8080

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=10s --retries=3 \
    CMD wget -qO- http://localhost:8080/health || exit 1

# 启动
ENTRYPOINT ["./server"]
```

### 逐段解释

**依赖下载和源码复制为什么要分开？** Docker 构建镜像时，每一条指令都会生成一个"层"（layer）。如果某一层的输入没有变化，Docker 就直接使用缓存，不重新执行。`go.mod` 和 `go.sum` 只有在依赖发生变化时才会改动，平时改业务代码不会碰它们。把 `COPY go.mod go.sum` 和 `RUN go mod download` 放在 `COPY . .` 之前，意味着日常开发改业务代码重新构建时，依赖下载那一层命中缓存，直接跳过。`go mod download` 在依赖多的项目里可能花好几分钟，缓存住了以后每次构建省掉这几分钟。

**为什么用 alpine 而不是 scratch 或 distroless？** `scratch` 是空镜像，连 shell 都没有，最小但调试困难——你没法 `docker exec` 进去看日志、查文件。`distroless` 比 scratch 多了证书和时区，但同样没有 shell。`alpine` 只有 5MB 大小，但自带了 shell 和包管理器，出问题时可以进容器排查。对于大多数情况来说，alpine 是安全和便利的最佳平衡点。

**`-ldflags="-s -w"` 能缩小多少？** 通常能减少 25%-30% 的体积。一个 30MB 的二进制文件可以缩到 20MB 左右。`-s` 去掉符号表，`-w` 去掉调试信息，生产环境不需要这些。

**HEALTHCHECK 的作用。** Docker 会按照设定的间隔执行健康检查命令。如果连续 3 次失败（`--retries=3`），容器状态会变为 `unhealthy`。Docker Compose 的 `depends_on` 配合 `condition: service_healthy` 可以用这个来控制启动顺序。Kubernetes 则有自己的健康检查机制（liveness/readiness probe），不使用 Docker 的 HEALTHCHECK。

### .dockerignore

和 `.gitignore` 类似，`.dockerignore` 告诉 Docker 构建时忽略哪些文件，减少构建上下文的大小，加快构建速度：

```
# deployments/docker/.dockerignore

.git
.github
.vscode
.idea

*.md
LICENSE
Makefile

# 测试文件不需要打进镜像
*_test.go
testdata/

# 本地编译的二进制文件
/tmp
/bin

# 环境变量文件（敏感信息不能打进镜像）
.env
.env.*
!.env.example

# 前端构建产物（如果有）
node_modules/

# 部署目录本身
deployments/
```

重点是 `.env` 文件绝不能打进镜像。数据库密码、JWT 密钥这些敏感信息应该通过运行时的环境变量或 Secret 管理工具注入，不是烧在镜像里。

### 构建镜像

```bash
# 在项目根目录执行
# -f 指定 Dockerfile 路径
# --build-arg 传入构建参数
# 最后的 . 是构建上下文（项目根目录）
docker build \
    -f deployments/docker/Dockerfile \
    --build-arg APP_VERSION=$(git describe --tags --always) \
    -t my-gin-project:$(git describe --tags --always) \
    -t my-gin-project:latest \
    .
```

`git describe --tags --always` 会输出类似 `v1.2.3` 或 `abc1234` 的版本标识。同时打两个 tag：一个带版本号（方便回滚到特定版本），一个 `latest`（方便开发时快速拉取）。

---

## 三、用 Makefile 简化操作

构建、测试、运行的命令越来越长，用 Makefile 封装起来：

```makefile
# Makefile

APP_NAME := my-gin-project
VERSION  := $(shell git describe --tags --always --dirty)
DOCKER_IMAGE := $(APP_NAME):$(VERSION)

.PHONY: build run test lint docker-build docker-push clean

# 本地编译
build:
	CGO_ENABLED=0 go build -trimpath \
		-ldflags="-s -w -X main.version=$(VERSION)" \
		-o bin/server ./cmd/server/

# 本地运行
run:
	go run ./cmd/server/

# 运行测试
test:
	go test -race -count=1 ./...

# 代码检查
lint:
	golangci-lint run ./...

# 构建 Docker 镜像
docker-build:
	docker build \
		-f deployments/docker/Dockerfile \
		--build-arg APP_VERSION=$(VERSION) \
		-t $(DOCKER_IMAGE) \
		-t $(APP_NAME):latest \
		.

# 推送到镜像仓库
docker-push: docker-build
	docker tag $(APP_NAME):latest registry.example.com/$(DOCKER_IMAGE)
	docker push registry.example.com/$(DOCKER_IMAGE)

# 本地启动全部服务（数据库 + 应用 + Nginx）
up:
	cd deployments && docker compose up -d --build

# 停止全部服务
down:
	cd deployments && docker compose down

# 查看日志
logs:
	cd deployments && docker compose logs -f app

# 清理
clean:
	rm -rf bin/
	docker rmi $(APP_NAME):latest 2>/dev/null || true
```

以后日常操作就简单了：

```bash
make build          # 本地编译
make test           # 跑测试
make docker-build   # 构建镜像
make up             # 本地启动全套环境
make logs           # 看日志
make down           # 停掉
```

---

## 四、Docker Compose：本地完整环境

开发和测试时需要一套完整的环境：数据库、Redis（如果用了）、应用本身、Nginx。Docker Compose 可以一条命令把所有东西拉起来。

```yaml
# deployments/docker-compose.yaml

services:
  # ==========================================
  # MySQL 数据库
  # ==========================================
  mysql:
    image: mysql:8.0
    container_name: gin-mysql
    restart: unless-stopped
    environment:
      MYSQL_ROOT_PASSWORD: ${DB_ROOT_PASSWORD:-rootpassword}
      MYSQL_DATABASE: ${DB_NAME:-gin_app}
      MYSQL_USER: ${DB_USER:-ginuser}
      MYSQL_PASSWORD: ${DB_PASSWORD:-ginpassword}
    ports:
      - "3306:3306"    # 映射到宿主机，方便本地用 GUI 工具连接
    volumes:
      - mysql_data:/var/lib/mysql
      # 初始化 SQL，容器首次启动时自动执行
      # - ./init.sql:/docker-entrypoint-initdb.d/init.sql
    healthcheck:
      test: ["CMD", "mysqladmin", "ping", "-h", "localhost"]
      interval: 10s
      timeout: 5s
      retries: 5
      start_period: 30s  # MySQL 启动需要时间，给 30 秒缓冲
    networks:
      - backend

  # ==========================================
  # Redis（如果需要）
  # ==========================================
  redis:
    image: redis:7-alpine
    container_name: gin-redis
    restart: unless-stopped
    command: redis-server --maxmemory 128mb --maxmemory-policy allkeys-lru
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 10s
      timeout: 3s
      retries: 3
    networks:
      - backend

  # ==========================================
  # Go 应用
  # ==========================================
  app:
    build:
      context: ..                          # 构建上下文是项目根目录
      dockerfile: deployments/docker/Dockerfile
      args:
        APP_VERSION: ${APP_VERSION:-dev}
    container_name: gin-app
    restart: unless-stopped
    environment:
      - GIN_MODE=release
      - DB_DSN=${DB_USER:-ginuser}:${DB_PASSWORD:-ginpassword}@tcp(mysql:3306)/${DB_NAME:-gin_app}?charset=utf8mb4&parseTime=True&loc=Local
      - REDIS_ADDR=redis:6379
      - JWT_SECRET=${JWT_SECRET:-your-secret-key-change-in-production}
      - PORT=8080
    ports:
      - "8080:8080"    # 调试时直接访问应用
    depends_on:
      mysql:
        condition: service_healthy
      redis:
        condition: service_healthy
    networks:
      - backend

  # ==========================================
  # Nginx 反向代理
  # ==========================================
  nginx:
    image: nginx:1.27-alpine
    container_name: gin-nginx
    restart: unless-stopped
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx/nginx.conf:/etc/nginx/nginx.conf:ro
      - ./nginx/conf.d:/etc/nginx/conf.d:ro
      - ./nginx/ssl:/etc/nginx/ssl:ro           # TLS 证书
      - nginx_logs:/var/log/nginx
    depends_on:
      - app
    networks:
      - backend

volumes:
  mysql_data:
  redis_data:
  nginx_logs:

networks:
  backend:
    driver: bridge
```

### 几个关键细节

**`depends_on` 配合 `condition: service_healthy`。** 普通的 `depends_on` 只保证容器启动了，不保证服务就绪。MySQL 容器启动后还需要几十秒初始化数据库，这段时间内应用连不上数据库就会报错退出。加了 `condition: service_healthy` 后，Docker Compose 会等 MySQL 的健康检查通过后再启动应用。

**环境变量的 `${VAR:-default}` 语法。** 这是 shell 的默认值语法——如果环境变量 `VAR` 存在就用它的值，不存在就用 `default`。这样不需要 `.env` 文件也能启动，但你可以创建一个 `.env` 文件来覆盖默认值。

**DB_DSN 中的 `mysql:3306`。** 在 Docker Compose 的同一个网络中，容器之间用服务名互相访问。这里的 `mysql` 就是 `services` 下面定义的 MySQL 服务名，Docker 会自动做 DNS 解析。

**volume 的 `:ro` 后缀。** `ro` 表示 read-only，容器内只能读不能写。Nginx 配置文件没有理由被容器修改，设为只读更安全。

### .env 文件

在 `deployments` 目录下创建一个 `.env` 文件给本地开发用：

```bash
# deployments/.env
# 这个文件不要提交到 Git

DB_ROOT_PASSWORD=rootpassword
DB_NAME=gin_app
DB_USER=ginuser
DB_PASSWORD=ginpassword
JWT_SECRET=local-dev-secret-key
APP_VERSION=dev
```

同时提供一个 `.env.example` 作为模板提交到仓库：

```bash
# deployments/.env.example
DB_ROOT_PASSWORD=
DB_NAME=gin_app
DB_USER=
DB_PASSWORD=
JWT_SECRET=
APP_VERSION=
```

---

## 五、Nginx 配置

### 主配置文件

```nginx
# deployments/nginx/nginx.conf

# 工作进程数，auto 会根据 CPU 核心数自动设置
worker_processes auto;

# 每个工作进程的最大连接数
events {
    worker_connections 1024;
    # 使用 epoll（Linux 高性能 I/O 模型）
    use epoll;
    # 一个进程同时接受多个新连接
    multi_accept on;
}

http {
    # 基础设置
    include       /etc/nginx/mime.types;
    default_type  application/octet-stream;

    # 启用 sendfile，减少文件发送时的上下文切换
    sendfile    on;
    tcp_nopush  on;
    tcp_nodelay on;

    # 连接超时
    keepalive_timeout 65;

    # Gzip 压缩
    # 在 Nginx 层做压缩，Go 应用就不需要自己压缩了
    gzip on;
    gzip_vary on;
    gzip_proxied any;
    gzip_comp_level 4;        # 压缩级别 1-9，4 是性能和压缩率的平衡点
    gzip_min_length 1024;     # 小于 1KB 的不压缩
    gzip_types
        text/plain
        text/css
        text/xml
        text/javascript
        application/json
        application/javascript
        application/xml
        application/xml+rss;

    # 日志格式
    log_format main '$remote_addr - $remote_user [$time_local] '
                    '"$request" $status $body_bytes_sent '
                    '"$http_referer" "$http_user_agent" '
                    '$request_time $upstream_response_time';

    access_log /var/log/nginx/access.log main;
    error_log  /var/log/nginx/error.log warn;

    # 限制请求体大小（防止上传超大文件）
    client_max_body_size 10m;

    # 请求缓冲区
    client_body_buffer_size 16k;
    client_header_buffer_size 1k;

    # 隐藏 Nginx 版本号（安全考虑）
    server_tokens off;

    # 加载虚拟主机配置
    include /etc/nginx/conf.d/*.conf;
}
```

`sendfile on` 让 Nginx 在发送文件时使用内核的 `sendfile()` 系统调用，数据直接从文件描述符传到 socket，不经过用户空间，效率更高。`tcp_nopush` 配合 `sendfile` 使用，让响应头和文件内容在一个 TCP 包里发出去，减少网络包数量。`tcp_nodelay` 在 keep-alive 连接上禁用 Nagle 算法，减少小包的发送延迟。

`gzip_comp_level 4` 这个值很关键。压缩级别越高，压缩率越好但 CPU 消耗越大。level 4 和 level 9 的压缩率差距只有 2%-3%，但 CPU 消耗差了一倍多。

### 虚拟主机配置（HTTP）

先从最简单的 HTTP 开始，后面再加 HTTPS：

```nginx
# deployments/nginx/conf.d/app.conf

# 上游服务定义
# 如果有多个应用实例，在这里列出来实现负载均衡
upstream gin_app {
    # 默认使用轮询策略
    server app:8080;

    # 如果有多个实例：
    # server app1:8080 weight=3;   # 权重更高，分配更多请求
    # server app2:8080 weight=1;
    # server app3:8080 backup;     # 备用，其他都挂了才启用

    # 保持长连接到上游，减少 TCP 握手开销
    keepalive 32;
}

server {
    listen 80;
    server_name yourdomain.com;   # 换成你的域名，本地测试用 localhost

    # 安全响应头
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;
    add_header Referrer-Policy "strict-origin-when-cross-origin" always;

    # API 请求转发到 Go 应用
    location /api/ {
        proxy_pass http://gin_app;

        # 传递客户端真实信息给后端
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;

        # 使用 HTTP/1.1 和 keep-alive 连接到上游
        proxy_http_version 1.1;
        proxy_set_header Connection "";

        # 超时设置
        proxy_connect_timeout 5s;     # 连接上游的超时
        proxy_send_timeout    10s;    # 发送请求到上游的超时
        proxy_read_timeout    30s;    # 等待上游响应的超时

        # 缓冲设置
        # 开启缓冲后，Nginx 会先把上游的完整响应读到缓冲区
        # 然后再慢慢发给客户端
        # 这样即使客户端网速慢，也不会长时间占用上游连接
        proxy_buffering on;
        proxy_buffer_size 4k;
        proxy_buffers 8 8k;
    }

    # 健康检查端点
    location /health {
        proxy_pass http://gin_app;
        proxy_set_header Host $host;

        # 健康检查不需要记录访问日志
        access_log off;
    }

    # 静态文件（如果有前端资源）
    location /static/ {
        alias /var/www/static/;
        expires 7d;                   # 浏览器缓存 7 天
        add_header Cache-Control "public, immutable";

        # 静态文件不需要走代理
        access_log off;
    }

    # 禁止访问隐藏文件
    location ~ /\. {
        deny all;
        access_log off;
        log_not_found off;
    }

    # 默认返回 404
    location / {
        return 404 '{"error": "not found"}';
        add_header Content-Type application/json;
    }
}
```

### 核心概念：反向代理做了什么

客户端发请求到 Nginx，Nginx 根据 `location` 规则决定把请求转发给谁。对于 `/api/` 开头的请求，Nginx 转发给 `gin_app`（也就是你的 Go 应用）。Go 应用处理完把响应返回给 Nginx，Nginx 再返回给客户端。

这个过程中 Nginx 做了很多额外的工作。它通过 `proxy_set_header X-Real-IP` 把客户端的真实 IP 传给后端，否则 Go 应用看到的 `RemoteAddr` 永远是 Nginx 容器的内网 IP。它通过 `proxy_buffering` 把上游响应缓冲住，即使客户端网速很慢（比如手机 3G 网络），也不会长时间占用 Go 应用的连接和 goroutine。

在 Go 应用中，你需要信任 Nginx 传过来的 `X-Real-IP` 头：

```go
// 在 Gin 中配置信任代理
r := gin.New()
// 只信任 Docker 内网的代理
r.SetTrustedProxies([]string{"172.16.0.0/12", "192.168.0.0/16"})
// 这样 c.ClientIP() 就会返回 X-Real-IP 中的真实客户端 IP
```

### 加上 HTTPS

生产环境必须上 HTTPS。这里用 Let's Encrypt 的免费证书（通过 certbot 获取），也可以用你自己购买的证书：

```nginx
# deployments/nginx/conf.d/app.conf

# HTTP -> HTTPS 重定向
server {
    listen 80;
    server_name yourdomain.com;

    # Let's Encrypt 证书验证路径，certbot 需要访问这个
    location /.well-known/acme-challenge/ {
        root /var/www/certbot;
    }

    # 其他所有请求重定向到 HTTPS
    location / {
        return 301 https://$host$request_uri;
    }
}

# HTTPS 服务
server {
    listen 443 ssl http2;
    server_name yourdomain.com;

    # TLS 证书
    ssl_certificate     /etc/nginx/ssl/fullchain.pem;
    ssl_certificate_key /etc/nginx/ssl/privkey.pem;

    # TLS 安全配置
    # 只允许 TLS 1.2 和 1.3，禁用老旧的 SSL 和 TLS 1.0/1.1
    ssl_protocols TLSv1.2 TLSv1.3;

    # 密码套件，优先使用服务器端的选择
    ssl_prefer_server_ciphers on;
    ssl_ciphers 'ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES128-GCM-SHA256:ECDHE-ECDSA-AES256-GCM-SHA384:ECDHE-RSA-AES256-GCM-SHA384';

    # SSL 会话缓存，避免每次连接都完整握手
    ssl_session_cache shared:SSL:10m;
    ssl_session_timeout 10m;
    ssl_session_tickets off;

    # HSTS：告诉浏览器以后只用 HTTPS 访问
    add_header Strict-Transport-Security "max-age=63072000; includeSubDomains" always;

    # 安全响应头
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;

    # API 转发（和 HTTP 版本一样）
    location /api/ {
        proxy_pass http://gin_app;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_http_version 1.1;
        proxy_set_header Connection "";
        proxy_connect_timeout 5s;
        proxy_send_timeout    10s;
        proxy_read_timeout    30s;
        proxy_buffering on;
        proxy_buffer_size 4k;
        proxy_buffers 8 8k;
    }

    location /health {
        proxy_pass http://gin_app;
        proxy_set_header Host $host;
        access_log off;
    }

    location /static/ {
        alias /var/www/static/;
        expires 7d;
        add_header Cache-Control "public, immutable";
        access_log off;
    }

    location ~ /\. {
        deny all;
    }

    location / {
        return 404 '{"error": "not found"}';
        add_header Content-Type application/json;
    }
}
```

TLS 终止在 Nginx 层完成，Nginx 和后端 Go 应用之间走的是 HTTP 明文（因为在同一个 Docker 网络里，是安全的）。这样 Go 应用不需要处理证书，配置更简单，性能也更好——Nginx 处理 TLS 的效率比 Go 高。

---

## 六、Certbot 自动申请证书

把 certbot 也容器化，实现证书的自动申请和续期：

```yaml
# 在 docker-compose.yaml 的 services 中添加

  certbot:
    image: certbot/certbot
    container_name: gin-certbot
    volumes:
      - ./nginx/ssl:/etc/letsencrypt
      - certbot_www:/var/www/certbot
    entrypoint: "/bin/sh -c 'trap exit TERM; while :; do certbot renew; sleep 12h & wait $${!}; done;'"
```

首次申请证书：

```bash
# 确保域名已经解析到服务器 IP，且 80 端口可访问
docker compose run --rm certbot certonly \
    --webroot \
    --webroot-path /var/www/certbot \
    -d yourdomain.com \
    --email your@email.com \
    --agree-tos \
    --no-eff-email
```

证书会保存在 `./nginx/ssl/live/yourdomain.com/` 目录下。certbot 容器会每 12 小时检查一次是否需要续期（Let's Encrypt 证书有效期 90 天，certbot 在到期前 30 天自动续期）。

> 参考：[使用 Let’s Encrypt 免费申请泛域名 SSL 证书，并实现自动续期](https://www.cnblogs.com/michaelshen/p/18538178)

---

## 七、生产服务器部署流程

本地开发和测试验证完之后，部署到生产服务器的步骤。

### 方式一：直接在服务器上构建

适合小团队、单服务器的场景：

```bash
# 1. SSH 到服务器
ssh user@your-server

# 2. 拉取最新代码
cd /opt/my-gin-project
git pull origin main

# 3. 构建并重启
cd deployments
docker compose down
docker compose up -d --build

# 4. 检查状态
docker compose ps
docker compose logs -f app
```

### 方式二：推送镜像到仓库

适合多服务器、CI/CD 的场景：

```bash
# === 在 CI 服务器或本地 ===

# 1. 构建镜像
docker build \
    -f deployments/docker/Dockerfile \
    --build-arg APP_VERSION=$(git describe --tags --always) \
    -t registry.example.com/my-gin-project:$(git describe --tags --always) \
    -t registry.example.com/my-gin-project:latest \
    .

# 2. 推送到镜像仓库
docker push registry.example.com/my-gin-project:$(git describe --tags --always)
docker push registry.example.com/my-gin-project:latest

# === 在生产服务器 ===

# 3. 拉取新镜像
docker pull registry.example.com/my-gin-project:latest

# 4. 重启服务
cd /opt/deployments
docker compose up -d
```

生产环境的 `docker-compose.yaml` 不需要 `build` 配置，直接用镜像：

```yaml
# 生产环境 docker-compose.prod.yaml

services:
  app:
    image: registry.example.com/my-gin-project:latest
    container_name: gin-app
    restart: unless-stopped
    environment:
      - GIN_MODE=release
      - DB_DSN=${DB_DSN}
      - JWT_SECRET=${JWT_SECRET}
      - PORT=8080
    depends_on:
      mysql:
        condition: service_healthy
    networks:
      - backend
```

### 零停机部署

如果要求部署时不中断服务，可以利用 `docker compose up -d` 的滚动更新能力，或者更实际的方案——跑两个应用实例：

```yaml
# docker-compose.prod.yaml 中配置多实例

  app:
    image: registry.example.com/my-gin-project:latest
    deploy:
      replicas: 2          # 两个实例
      update_config:
        parallelism: 1     # 每次更新一个
        delay: 10s          # 两次更新之间间隔 10 秒
        order: start-first  # 先启动新的，再停旧的
    # ...
```

```nginx
# Nginx upstream 中配置两个实例
upstream gin_app {
    server app:8080;
    keepalive 32;
}
```

Docker Compose 的 `deploy` 配置在 `docker compose up` 时部分生效（完整功能需要 Docker Swarm）。对于简单的单机多实例场景已经够用。如果需要更完善的零停机部署，建议上 Kubernetes。

---

## 八、GitHub Actions CI/CD

自动化整个流程——代码推送后自动测试、构建、部署：

```yaml
# .github/workflows/deploy.yaml

name: Build and Deploy

on:
  push:
    branches: [main]
    tags: ['v*']

env:
  REGISTRY: ghcr.io
  IMAGE_NAME: ${{ github.repository }}

jobs:
  # ============================
  # 阶段一：测试
  # ============================
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-go@v5
        with:
          go-version: '1.26'

      - name: Run tests
        run: go test -race -count=1 ./...

      - name: Run linter
        uses: golangci/golangci-lint-action@v6
        with:
          version: latest

  # ============================
  # 阶段二：构建并推送镜像
  # ============================
  build:
    needs: test
    runs-on: ubuntu-latest
    permissions:
      contents: read
      packages: write

    steps:
      - uses: actions/checkout@v4

      - name: Login to Container Registry
        uses: docker/login-action@v3
        with:
          registry: ${{ env.REGISTRY }}
          username: ${{ github.actor }}
          password: ${{ secrets.GITHUB_TOKEN }}

      - name: Extract metadata
        id: meta
        uses: docker/metadata-action@v5
        with:
          images: ${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}
          tags: |
            type=ref,event=branch
            type=semver,pattern={{version}}
            type=sha,prefix=

      - name: Build and push
        uses: docker/build-push-action@v5
        with:
          context: .
          file: deployments/docker/Dockerfile
          push: true
          tags: ${{ steps.meta.outputs.tags }}
          build-args: |
            APP_VERSION=${{ github.sha }}

  # ============================
  # 阶段三：部署到服务器
  # ============================
  deploy:
    needs: build
    runs-on: ubuntu-latest
    if: github.ref == 'refs/heads/main'

    steps:
      - name: Deploy to server
        uses: appleboy/ssh-action@v1
        with:
          host: ${{ secrets.SERVER_HOST }}
          username: ${{ secrets.SERVER_USER }}
          key: ${{ secrets.SERVER_SSH_KEY }}
          script: |
            cd /opt/my-gin-project/deployments
            docker compose pull app
            docker compose up -d app
            # 等待健康检查通过
            sleep 10
            curl -f http://localhost:8080/health || exit 1
            echo "Deploy successful"
```

在 GitHub 仓库的 Settings → Secrets 中配置 `SERVER_HOST`、`SERVER_USER`、`SERVER_SSH_KEY`。每次推送到 main 分支，会自动跑测试、构建镜像、部署到服务器。打 tag 时只构建镜像不自动部署（生产环境的 tag 发布通常需要手动确认）。

---

## 九、日志管理

容器化后日志管理方式和传统部署不同。容器里的日志应该输出到 stdout/stderr，而不是写文件。

### 应用日志输出到 stdout

```go
// 在 main.go 中初始化 zap logger
import "go.uber.org/zap"

func initLogger() *zap.Logger {
    config := zap.NewProductionConfig()

    // 输出到 stdout，不写文件
    // Docker 会自动收集 stdout 的内容
    config.OutputPaths = []string{"stdout"}
    config.ErrorOutputPaths = []string{"stderr"}

    // 生产环境用 JSON 格式，方便日志系统解析
    config.Encoding = "json"

    logger, _ := config.Build()
    return logger
}
```

### Docker 日志驱动配置

```yaml
# docker-compose.yaml 中配置日志
services:
  app:
    # ...
    logging:
      driver: json-file
      options:
        max-size: "50m"      # 单个日志文件最大 50MB
        max-file: "5"        # 最多保留 5 个文件
        # 总共最多 250MB 日志，自动轮转
```

不配置 `max-size` 的话，Docker 的 `json-file` 日志驱动会无限增长，迟早把磁盘撑满。这是生产环境最常见的事故之一。

查看日志：

```bash
# 实时查看应用日志
docker compose logs -f app

# 查看最近 100 行
docker compose logs --tail 100 app

# 查看某个时间段
docker compose logs --since "2024-01-01T10:00:00" app
```

---

## 十、监控与告警

部署完不是终点，还需要知道线上服务的运行状态。

### 基础监控端点

```go
// internal/handler/monitor.go
package handler

import (
    "runtime"
    "time"

    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

var startTime = time.Now()

func RegisterMonitorRoutes(r *gin.Engine, db *gorm.DB) {
    monitor := r.Group("/internal")
    {
        monitor.GET("/health", healthCheck(db))
        monitor.GET("/metrics", metrics(db))
    }
}

func healthCheck(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        // 检查数据库连接
        sqlDB, err := db.DB()
        if err != nil {
            c.JSON(503, gin.H{"status": "error", "db": err.Error()})
            return
        }
        if err := sqlDB.Ping(); err != nil {
            c.JSON(503, gin.H{"status": "error", "db": err.Error()})
            return
        }

        c.JSON(200, gin.H{"status": "ok"})
    }
}

func metrics(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        var mem runtime.MemStats
        runtime.ReadMemStats(&mem)

        sqlDB, _ := db.DB()
        dbStats := sqlDB.Stats()

        c.JSON(200, gin.H{
            "uptime":           time.Since(startTime).String(),
            "goroutines":       runtime.NumGoroutine(),
            "memory": gin.H{
                "alloc_mb":       mem.Alloc / 1024 / 1024,
                "sys_mb":         mem.Sys / 1024 / 1024,
                "gc_cycles":      mem.NumGC,
                "gc_pause_total": time.Duration(mem.PauseTotalNs).String(),
            },
            "db": gin.H{
                "open_connections": dbStats.OpenConnections,
                "in_use":          dbStats.InUse,
                "idle":            dbStats.Idle,
                "wait_count":      dbStats.WaitCount,
                "wait_duration":   dbStats.WaitDuration.String(),
            },
        })
    }
}
```

Nginx 层要确保 `/internal/*` 路径不对外暴露：

```nginx
# 内部监控端点，只允许内网访问
location /internal/ {
    # 只允许特定 IP 访问
    allow 10.0.0.0/8;
    allow 172.16.0.0/12;
    allow 192.168.0.0/16;
    allow 127.0.0.1;
    deny all;

    proxy_pass http://gin_app;
    proxy_set_header Host $host;
}
```

---

## 十一、安全加固清单

部署到生产环境，安全是不能忽视的：

**应用层面。** 容器用非 root 用户运行（Dockerfile 中已经配了 `USER appuser`）。敏感配置通过环境变量注入，不烧在镜像里。`.env` 文件不提交到 Git。JWT 密钥、数据库密码在生产环境用足够强的随机值。

**Nginx 层面。** 隐藏 Nginx 版本号（`server_tokens off`）。配置安全响应头（X-Frame-Options、X-Content-Type-Options 等）。TLS 只启用 1.2 和 1.3。启用 HSTS。限制请求体大小（`client_max_body_size`）。内部端点限制访问 IP。

**Docker 层面。** 基础镜像定期更新（Alpine 有安全补丁时及时升级）。不在镜像里安装不需要的软件包。Volume 挂载用 `:ro` 尽可能设为只读。不使用 `--privileged` 运行容器。

**网络层面。** 只暴露必要的端口（80 和 443 给 Nginx，其他端口不映射到宿主机）。数据库端口（3306）在生产环境不对外暴露，只在 Docker 内部网络中可访问。开启服务器防火墙（UFW 或 iptables），只放行 22（SSH）、80、443。

---

## 十二、排错指南

部署后最常遇到的问题和解决方法。

**应用启动后立即退出。** 执行 `docker compose logs app` 查看错误日志。最常见的原因是数据库连不上——检查 DSN 配置、确认 MySQL 容器已经就绪。

**Nginx 返回 502 Bad Gateway。** 说明 Nginx 连不上后端应用。检查应用是否在运行（`docker compose ps`），检查 upstream 配置的服务名和端口是否正确，确认两个容器在同一个 Docker 网络中。

**Nginx 返回 504 Gateway Timeout。** 后端应用处理超时。检查 Nginx 的 `proxy_read_timeout` 是否够长，同时排查 Go 应用中的慢查询或外部调用超时。

**容器健康检查一直 unhealthy。** 进入容器手动执行健康检查命令看报什么错：`docker exec gin-app wget -qO- http://localhost:8080/health`。

**磁盘空间不足。** Docker 的镜像、容器、卷、构建缓存会逐渐占满磁盘。定期清理：`docker system prune -a --volumes`（注意这会删除所有未使用的资源，包括未运行的容器和未被引用的卷）。

---

## 完整架构回顾

```
                     互联网
                       │
                       ▼
              ┌────────────────┐
              │   云服务商防火墙  │  只放行 80, 443, 22
              └────────┬───────┘
                       │
                       ▼
              ┌────────────────┐
              │     Nginx      │  TLS 终止、Gzip、安全头、限流
              │   (443/80)     │  静态文件服务、请求缓冲
              └────────┬───────┘
                       │ HTTP (Docker 内网)
                       ▼
              ┌────────────────┐
              │   Gin 应用      │  业务逻辑、JSON API
              │   (8080)       │  优雅关闭、健康检查
              └───┬────────┬───┘
                  │        │
            ┌─────┘        └─────┐
            ▼                    ▼
   ┌────────────────┐   ┌────────────────┐
   │     MySQL      │   │     Redis      │
   │   (3306)       │   │   (6379)       │
   └────────────────┘   └────────────────┘

   所有组件运行在同一个 Docker 网络中
   只有 Nginx 的 80/443 端口暴露到宿主机
```

---