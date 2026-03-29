package router

import (
	"github.com/gorilla/mux"

	"myproject/internal/auth"
	"myproject/internal/handler"
	"myproject/internal/middleware"
)

func New(
	uh *handler.UserHandler,
	oh *handler.OrderHandler,
	hh *handler.HealthHandler,
	jwtMgr *auth.JWTManager,
) *mux.Router {
	r := mux.NewRouter()

	// 全局中间件
	r.Use(middleware.Recovery)
	r.Use(middleware.Logging)
	r.Use(middleware.CORS)

	// 健康检查（无需认证）
	r.HandleFunc("/healthz", hh.Health).Methods("GET")

	// API v1
	api := r.PathPrefix("/api/v1").Subrouter()

	// 公开路由
	api.HandleFunc("/register", uh.Register).Methods("POST")
	api.HandleFunc("/login", uh.Login).Methods("POST")

	// 需认证的路由
	protected := api.PathPrefix("").Subrouter()
	protected.Use(middleware.Auth(jwtMgr))
	protected.Use(middleware.RateLimit(100))

	// 用户
	protected.HandleFunc("/users/me", uh.Me).Methods("GET")
	protected.HandleFunc("/users", uh.List).Methods("GET")
	protected.HandleFunc("/users/{id}", uh.GetByID).Methods("GET")

	// 订单
	protected.HandleFunc("/orders", oh.Create).Methods("POST")
	protected.HandleFunc("/orders", oh.List).Methods("GET")
	protected.HandleFunc("/orders/{id}", oh.GetByID).Methods("GET")
	protected.HandleFunc("/orders/{id}/pay", oh.Pay).Methods("POST")

	return r
}