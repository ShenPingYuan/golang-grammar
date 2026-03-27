package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"

	"myproject/internal/auth"
	"myproject/internal/dto"
	apperrors "myproject/internal/errors"
	"myproject/internal/event"
	"myproject/internal/model"
	"myproject/internal/repository"
	"myproject/internal/util"
)

type UserService interface {
	Register(ctx context.Context, req dto.CreateUserRequest) (*dto.UserResponse, error)
	Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error)
	GetByID(ctx context.Context, id string) (*dto.UserResponse, error)
	List(ctx context.Context, page, size int) ([]*dto.UserResponse, int, error)
}

type userService struct {
	repo   repository.UserRepository
	jwt    *auth.JWTManager
	events event.Bus
}

func NewUserService(repo repository.UserRepository, jwt *auth.JWTManager, events event.Bus) UserService {
	return &userService{repo: repo, jwt: jwt, events: events}
}

func (s *userService) Register(ctx context.Context, req dto.CreateUserRequest) (*dto.UserResponse, error) {
	if !util.IsValidEmail(req.Email) {
		return nil, apperrors.ErrInvalidInput
	}
	if !util.IsMinLength(req.Password, 6) {
		return nil, apperrors.ErrInvalidInput
	}

	hash, err := util.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		ID:       generateID(),
		Username: req.Username,
		Email:    req.Email,
		Password: hash,
		Role:     auth.RoleUser,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	// 发布用户创建事件
	_ = s.events.Publish(ctx, event.UserCreatedEvent{
		UserID:   user.ID,
		Username: user.Username,
		Email:    user.Email,
	})

	return toUserResponse(user), nil
}

func (s *userService) Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := s.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, apperrors.ErrUnauthorized
	}
	if !util.CheckPassword(user.Password, req.Password) {
		return nil, apperrors.ErrUnauthorized
	}

	token, err := s.jwt.Generate(user.ID, user.Role)
	if err != nil {
		return nil, err
	}
	return &dto.LoginResponse{Token: token}, nil
}

func (s *userService) GetByID(ctx context.Context, id string) (*dto.UserResponse, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return toUserResponse(user), nil
}

func (s *userService) List(ctx context.Context, page, size int) ([]*dto.UserResponse, int, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}
	offset := (page - 1) * size
	users, total, err := s.repo.List(ctx, offset, size)
	if err != nil {
		return nil, 0, err
	}
	resp := make([]*dto.UserResponse, len(users))
	for i, u := range users {
		resp[i] = toUserResponse(u)
	}
	return resp, total, nil
}

func toUserResponse(u *model.User) *dto.UserResponse {
	return &dto.UserResponse{
		ID:       u.ID,
		Username: u.Username,
		Email:    u.Email,
		Role:     u.Role,
	}
}

func generateID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}