package grpc

import (
	"myproject/internal/auth"
	grpcinterceptor "myproject/internal/grpc/interceptor"
	grpcservice "myproject/internal/grpc/service"
	"myproject/internal/service"
	"myproject/pkg/logger"

	"google.golang.org/grpc"
)

// NewServer 创建并配置 gRPC 服务器
// 注意: 需先运行 make proto 生成 gen/proto/ 代码
func NewServer(
	userSvc service.UserService,
	orderSvc service.OrderService,
	jwtMgr *auth.JWTManager,
	l *logger.Logger,
) *grpc.Server {
	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpcinterceptor.Recovery(l),
			grpcinterceptor.Logging(l),
			grpcinterceptor.Auth(jwtMgr),
		),
	)

	// 注册 gRPC 服务
	grpcservice.RegisterUserService(s, userSvc, jwtMgr)
	grpcservice.RegisterOrderService(s, orderSvc)

	return s
}