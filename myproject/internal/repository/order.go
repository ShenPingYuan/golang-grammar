package repository

import (
	"context"
	"sync"
	"time"

	apperrors "myproject/internal/errors"
	"myproject/internal/model"
)

type OrderRepository interface {
	Create(ctx context.Context, order *model.Order) error
	GetByID(ctx context.Context, id string) (*model.Order, error)
	ListByUserID(ctx context.Context, userID string, offset, limit int) ([]*model.Order, int, error)
	Update(ctx context.Context, order *model.Order) error
}

// --- 内存实现 ---

type memoryOrderRepo struct {
	mu     sync.RWMutex
	orders map[string]*model.Order
}

func NewMemoryOrderRepository() OrderRepository {
	return &memoryOrderRepo{orders: make(map[string]*model.Order)}
}

func (r *memoryOrderRepo) Create(_ context.Context, order *model.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()
	r.orders[order.ID] = order
	return nil
}

func (r *memoryOrderRepo) GetByID(_ context.Context, id string) (*model.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	order, ok := r.orders[id]
	if !ok {
		return nil, apperrors.ErrNotFound
	}
	return order, nil
}

func (r *memoryOrderRepo) ListByUserID(_ context.Context, userID string, offset, limit int) ([]*model.Order, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*model.Order
	for _, o := range r.orders {
		if o.UserID == userID {
			result = append(result, o)
		}
	}
	total := len(result)
	if offset >= total {
		return []*model.Order{}, total, nil
	}
	end := offset + limit
	if end > total {
		end = total
	}
	return result[offset:end], total, nil
}

func (r *memoryOrderRepo) Update(_ context.Context, order *model.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.orders[order.ID]; !ok {
		return apperrors.ErrNotFound
	}
	order.UpdatedAt = time.Now()
	r.orders[order.ID] = order
	return nil
}