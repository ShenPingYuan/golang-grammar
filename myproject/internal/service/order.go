package service

import (
	"context"

	"myproject/internal/dto"
	apperrors "myproject/internal/errors"
	"myproject/internal/event"
	"myproject/internal/model"
	"myproject/internal/repository"
)

type OrderService interface {
	Create(ctx context.Context, userID string, req dto.CreateOrderRequest) (*dto.OrderResponse, error)
	GetByID(ctx context.Context, id string) (*dto.OrderResponse, error)
	ListByUser(ctx context.Context, userID string, page, size int) ([]*dto.OrderResponse, int, error)
	Pay(ctx context.Context, id string) (*dto.OrderResponse, error)
}

type orderService struct {
	repo   repository.OrderRepository
	events event.Bus
}

func NewOrderService(repo repository.OrderRepository, events event.Bus) OrderService {
	return &orderService{repo: repo, events: events}
}

func (s *orderService) Create(ctx context.Context, userID string, req dto.CreateOrderRequest) (*dto.OrderResponse, error) {
	if req.Amount <= 0 {
		return nil, apperrors.ErrInvalidInput
	}

	order := &model.Order{
		ID:      generateID(),
		UserID:  userID,
		Product: req.Product,
		Amount:  req.Amount,
		Status:  model.OrderStatusPending,
	}

	if err := s.repo.Create(ctx, order); err != nil {
		return nil, err
	}
	return toOrderResponse(order), nil
}

func (s *orderService) GetByID(ctx context.Context, id string) (*dto.OrderResponse, error) {
	order, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return toOrderResponse(order), nil
}

func (s *orderService) ListByUser(ctx context.Context, userID string, page, size int) ([]*dto.OrderResponse, int, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}
	offset := (page - 1) * size
	orders, total, err := s.repo.ListByUserID(ctx, userID, offset, size)
	if err != nil {
		return nil, 0, err
	}
	resp := make([]*dto.OrderResponse, len(orders))
	for i, o := range orders {
		resp[i] = toOrderResponse(o)
	}
	return resp, total, nil
}

func (s *orderService) Pay(ctx context.Context, id string) (*dto.OrderResponse, error) {
	order, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if order.Status == model.OrderStatusPaid {
		return nil, apperrors.ErrInvalidInput
	}

	order.Status = model.OrderStatusPaid
	if err := s.repo.Update(ctx, order); err != nil {
		return nil, err
	}

	_ = s.events.Publish(ctx, event.OrderPaidEvent{
		OrderID: order.ID,
		UserID:  order.UserID,
		Amount:  order.Amount,
	})

	return toOrderResponse(order), nil
}

func toOrderResponse(o *model.Order) *dto.OrderResponse {
	return &dto.OrderResponse{
		ID:      o.ID,
		UserID:  o.UserID,
		Product: o.Product,
		Amount:  o.Amount,
		Status:  o.Status,
	}
}