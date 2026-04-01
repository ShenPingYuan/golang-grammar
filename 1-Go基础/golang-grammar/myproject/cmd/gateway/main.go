package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"myproject/internal/config"
	"myproject/pkg/logger"
)

// 简易 HTTP→gRPC 网关示例
// 生产环境建议使用 grpc-gateway (github.com/grpc-ecosystem/grpc-gateway)
func main() {
	cfg := config.Load("configs/config.yaml")
	l := logger.New(cfg.Log.Level)

	target, _ := url.Parse("http://localhost:" + cfg.Server.Port)
	proxy := httputil.NewSingleHostReverseProxy(target)

	l.Info("Gateway starting on :8081 → forwarding to :" + cfg.Server.Port)
	log.Fatal(http.ListenAndServe(":8081", proxy))
}