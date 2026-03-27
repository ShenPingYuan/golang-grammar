package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"myproject/internal/dto"
	apperrors "myproject/internal/errors"
	"myproject/internal/middleware"
	"myproject/internal/service"
)

type UserHandler struct {
	svc service.UserService
}

func NewUserHandler(svc service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Code: 400, Message: "invalid request body"})
		return
	}
	resp, err := h.svc.Register(r.Context(), req)
	if err != nil {
		handleError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, resp)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Code: 400, Message: "invalid request body"})
		return
	}
	resp, err := h.svc.Login(r.Context(), req)
	if err != nil {
		handleError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	resp, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		handleError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	size, _ := strconv.Atoi(r.URL.Query().Get("size"))
	items, total, err := h.svc.List(r.Context(), page, size)
	if err != nil {
		handleError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, dto.ListResponse{Items: items, Total: total, Page: page, Size: size})
}

// --- 下面是 UserHandler 用不到但演示获取当前用户 ---

func (h *UserHandler) Me(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	resp, err := h.svc.GetByID(r.Context(), userID)
	if err != nil {
		handleError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

// --- 公共辅助函数 ---

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func handleError(w http.ResponseWriter, err error) {
	switch {
	case apperrors.Is(err, apperrors.ErrNotFound):
		writeJSON(w, http.StatusNotFound, dto.ErrorResponse{Code: 404, Message: "not found"})
	case apperrors.Is(err, apperrors.ErrAlreadyExists):
		writeJSON(w, http.StatusConflict, dto.ErrorResponse{Code: 409, Message: "already exists"})
	case apperrors.Is(err, apperrors.ErrUnauthorized):
		writeJSON(w, http.StatusUnauthorized, dto.ErrorResponse{Code: 401, Message: "unauthorized"})
	case apperrors.Is(err, apperrors.ErrForbidden):
		writeJSON(w, http.StatusForbidden, dto.ErrorResponse{Code: 403, Message: "forbidden"})
	case apperrors.Is(err, apperrors.ErrInvalidInput):
		writeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Code: 400, Message: "invalid input"})
	default:
		writeJSON(w, http.StatusInternalServerError, dto.ErrorResponse{Code: 500, Message: "internal error"})
	}
}