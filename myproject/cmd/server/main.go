package main

import (
	"log"
	"net/http"

	"myproject/internal/auth"
	"myproject/internal/config"
	"myproject/internal/event"
	eventhandler "myproject/internal/event/handler"
	"myproject/internal/handler"
	"myproject/internal/notify"
	"myproject/internal/repository"
	"myproject/internal/router"
	"myproject/internal/service"
	"myproject/pkg/logger"
)

func main() {
	cfg := config.Load("configs/config.yaml")
	l := logger.New(cfg.Log.Level)

	// 事件总线
	bus := event.NewBus()
	notifier := notify.NewConsole()
	eventhandler.Register(bus, notifier, l)

	// 仓储层（内存实现，开箱即用）
	userRepo := repository.NewMemoryUserRepository()
	orderRepo := repository.NewMemoryOrderRepository()

	// 认证
	jwtMgr := auth.NewJWTManager(cfg.JWT.Secret, cfg.JWT.ExpiryDuration())

	// 业务层
	userSvc := service.NewUserService(userRepo, jwtMgr, bus)
	orderSvc := service.NewOrderService(orderRepo, bus)

	// HTTP 处理器
	uh := handler.NewUserHandler(userSvc)
	oh := handler.NewOrderHandler(orderSvc)
	hh := handler.NewHealthHandler()

	// 路由
	r := router.New(uh, oh, hh, jwtMgr)

	l.Info("HTTP server starting on :" + cfg.Server.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Server.Port, r))
}