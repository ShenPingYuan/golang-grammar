package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"myproject/internal/auth"
	"myproject/internal/dto"
	"myproject/internal/event"
	"myproject/internal/handler"
	"myproject/internal/repository"
	"myproject/internal/router"
	"myproject/internal/service"
)

func setupTestRouter() http.Handler {
	bus := event.NewBus()
	userRepo := repository.NewMemoryUserRepository()
	orderRepo := repository.NewMemoryOrderRepository()
	jwtMgr := auth.NewJWTManager("test-secret", 3600_000_000_000) // 1h

	userSvc := service.NewUserService(userRepo, jwtMgr, bus)
	orderSvc := service.NewOrderService(orderRepo, bus)

	uh := handler.NewUserHandler(userSvc)
	oh := handler.NewOrderHandler(orderSvc)
	hh := handler.NewHealthHandler()

	return router.New(uh, oh, hh, jwtMgr)
}

func TestHealthEndpoint(t *testing.T) {
	r := setupTestRouter()
	req := httptest.NewRequest("GET", "/healthz", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestRegisterAndLogin(t *testing.T) {
	r := setupTestRouter()

	// 注册
	body, _ := json.Marshal(dto.CreateUserRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	})
	req := httptest.NewRequest("POST", "/api/v1/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("register: expected 201, got %d, body: %s", w.Code, w.Body.String())
	}

	// 登录
	loginBody, _ := json.Marshal(dto.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	})
	req = httptest.NewRequest("POST", "/api/v1/login", bytes.NewReader(loginBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("login: expected 200, got %d, body: %s", w.Code, w.Body.String())
	}

	var loginResp dto.LoginResponse
	_ = json.NewDecoder(w.Body).Decode(&loginResp)
	if loginResp.Token == "" {
		t.Fatal("login: expected non-empty token")
	}

	// 用 token 访问受保护接口
	req = httptest.NewRequest("GET", "/api/v1/users/me", nil)
	req.Header.Set("Authorization", "Bearer "+loginResp.Token)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("me: expected 200, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestProtectedRouteWithoutToken(t *testing.T) {
	r := setupTestRouter()
	req := httptest.NewRequest("GET", "/api/v1/users", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}