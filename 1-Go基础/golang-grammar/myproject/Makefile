.PHONY: run build test lint proto clean docker seed migrate-up migrate-down

# --- 开发 ---
run:
	go run ./cmd/server

run-grpc:
	go run ./cmd/grpcserver

run-worker:
	go run ./cmd/worker

run-scheduler:
	go run ./cmd/scheduler

# --- 构建 ---
build:
	@mkdir -p bin
	go build -o bin/server    ./cmd/server
	go build -o bin/worker    ./cmd/worker
	go build -o bin/scheduler ./cmd/scheduler
	go build -o bin/migrate   ./cmd/migrate
	go build -o bin/cli       ./cmd/cli

# --- 测试 ---
test:
	go test ./... -race -cover -count=1

test-integration:
	go test ./test/... -race -v

coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out

# --- 代码质量 ---
lint:
	golangci-lint run ./...

fmt:
	gofmt -w .
	goimports -w -local myproject .

vet:
	go vet ./...

# --- Protobuf ---
proto:
	bash scripts/protogen.sh

proto-buf:
	buf generate

# --- 数据库迁移 ---
migrate-up:
	go run ./cmd/migrate -direction=up

migrate-down:
	go run ./cmd/migrate -direction=down

# --- 数据填充 ---
seed:
	go run scripts/seed.go

# --- Docker ---
docker:
	docker-compose -f deployments/docker-compose.yaml up --build

docker-down:
	docker-compose -f deployments/docker-compose.yaml down

# --- 清理 ---
clean:
	rm -rf bin/ coverage.out gen/proto/