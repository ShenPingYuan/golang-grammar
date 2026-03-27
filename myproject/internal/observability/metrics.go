package observability

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

// 简易指标收集器（生产环境请使用 github.com/prometheus/client_golang）

var (
	requestCount  atomic.Int64
	requestErrors atomic.Int64
)

func IncRequestCount()  { requestCount.Add(1) }
func IncRequestErrors() { requestErrors.Add(1) }

// MetricsHandler 返回指标端点
func MetricsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("# HELP http_requests_total Total HTTP requests\n"))
		w.Write([]byte("# TYPE http_requests_total counter\n"))
		w.Write([]byte("http_requests_total " + itoa(requestCount.Load()) + "\n"))
		w.Write([]byte("# HELP http_request_errors_total Total HTTP request errors\n"))
		w.Write([]byte("# TYPE http_request_errors_total counter\n"))
		w.Write([]byte("http_request_errors_total " + itoa(requestErrors.Load()) + "\n"))
	}
}

func itoa(n int64) string {
	return fmt.Sprintf("%d", n)
}
