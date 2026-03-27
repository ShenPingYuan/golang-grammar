package grpcinterceptor

import (
	"context"
	"time"

	"myproject/pkg/logger"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func Logging(l *logger.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()
		resp, err := handler(ctx, req)
		st, _ := status.FromError(err)
		l.Info("gRPC request",
			"method", info.FullMethod,
			"code", st.Code().String(),
			"duration", time.Since(start).String(),
		)
		return resp, err
	}
}