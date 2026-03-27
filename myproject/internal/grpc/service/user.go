package grpcservice

import (
	"context"
	"encoding/json"

	"myproject/internal/auth"
	"myproject/internal/dto"
	"myproject/internal/service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UserGRPCService 通过 gRPC 暴露 UserService
// 完整版应实现 gen/proto/user 中的 UserServiceServer 接口
type UserGRPCService struct {
	svc    service.UserService
	jwtMgr *auth.JWTManager
}

func RegisterUserService(s *grpc.Server, svc service.UserService, jwtMgr *auth.JWTManager) {
	// 完整版: pb.RegisterUserServiceServer(s, &UserGRPCService{svc: svc, jwtMgr: jwtMgr})
	// 当前为演示，运行 make proto 后替换为生成的注册函数
	_ = &UserGRPCService{svc: svc, jwtMgr: jwtMgr}
}

func (s *UserGRPCService) CreateUser(ctx context.Context, req json.RawMessage) (interface{}, error) {
	var r dto.CreateUserRequest
	if err := json.Unmarshal(req, &r); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	resp, err := s.svc.Register(ctx, r)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return resp, nil
}