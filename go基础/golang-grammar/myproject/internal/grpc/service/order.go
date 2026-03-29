package grpcservice

import (
	"myproject/internal/service"

	"google.golang.org/grpc"
)

type OrderGRPCService struct {
	svc service.OrderService
}

func RegisterOrderService(s *grpc.Server, svc service.OrderService) {
	// 完整版: pb.RegisterOrderServiceServer(s, &OrderGRPCService{svc: svc})
	_ = &OrderGRPCService{svc: svc}
}