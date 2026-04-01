package observability

import "log/slog"

// InitTracer 初始化链路追踪
// 生产环境请使用 go.opentelemetry.io/otel
func InitTracer(serviceName, endpoint string) {
	slog.Info("tracing configured (stub)", "service", serviceName, "endpoint", endpoint)
	// TODO:
	// tp := sdktrace.NewTracerProvider(...)
	// otel.SetTracerProvider(tp)
}