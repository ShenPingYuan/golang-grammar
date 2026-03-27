package main

import (
	"log"
	"net"

	"myproject/internal/auth"
	"myproject/internal/config"
	"myproject/internal/event"
	eventhandler "myproject/internal/event/handler"
	internalgrpc "myproject/internal/grpc"
	"myproject/internal/notify"
	"myproject/internal/repository"
	"myproject/internal/service"
	"myproject/pkg/logger"
)

func main() {
	cfg := config.Load("configs/config.yaml")
	l := logger.New(cfg.Log.Level)

	bus := event.NewBus()
	notifier := notify.NewConsole()
	eventhandler.Register(bus, notifier, l)

	userRepo := repository.NewMemoryUserRepository()
	orderRepo := repository.NewMemoryOrderRepository()

	jwtMgr := auth.NewJWTManager(cfg.JWT.Secret, cfg.JWT.ExpiryDuration())

	userSvc := service.NewUserService(userRepo, jwtMgr, bus)
	orderSvc := service.NewOrderService(orderRepo, bus)

	grpcServer := internalgrpc.NewServer(userSvc, orderSvc, jwtMgr, l)

	lis, err := net.Listen("tcp", ":"+cfg.GRPC.Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	l.Info("gRPC server starting on :" + cfg.GRPC.Port)
	log.Fatal(grpcServer.Serve(lis))
}