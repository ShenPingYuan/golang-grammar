package repository

import (
	"context"
	"sync"
	"time"

	apperrors "myproject/internal/errors"
	"myproject/internal/model"
)

// UserRepository 用户数据访问接口
type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	GetByID(ctx context.Context, id string) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	List(ctx context.Context, offset, limit int) ([]*model.User, int, error)
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id string) error
}

// --- 内存实现 ---

type memoryUserRepo struct {
	mu    sync.RWMutex
	users map[string]*model.User
}

func NewMemoryUserRepository() UserRepository {
	return &memoryUserRepo{users: make(map[string]*model.User)}
}

func (r *memoryUserRepo) Create(_ context.Context, user *model.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, u := range r.users {
		if u.Email == user.Email {
			return apperrors.ErrAlreadyExists
		}
	}
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	r.users[user.ID] = user
	return nil
}

func (r *memoryUserRepo) GetByID(_ context.Context, id string) (*model.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	user, ok := r.users[id]
	if !ok {
		return nil, apperrors.ErrNotFound
	}
	return user, nil
}

func (r *memoryUserRepo) GetByEmail(_ context.Context, email string) (*model.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, u := range r.users {
		if u.Email == email {
			return u, nil
		}
	}
	return nil, apperrors.ErrNotFound
}

func (r *memoryUserRepo) List(_ context.Context, offset, limit int) ([]*model.User, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	all := make([]*model.User, 0, len(r.users))
	for _, u := range r.users {
		all = append(all, u)
	}
	total := len(all)
	if offset >= total {
		return []*model.User{}, total, nil
	}
	end := offset + limit
	if end > total {
		end = total
	}
	return all[offset:end], total, nil
}

func (r *memoryUserRepo) Update(_ context.Context, user *model.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.users[user.ID]; !ok {
		return apperrors.ErrNotFound
	}
	user.UpdatedAt = time.Now()
	r.users[user.ID] = user
	return nil
}

func (r *memoryUserRepo) Delete(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.users[id]; !ok {
		return apperrors.ErrNotFound
	}
	delete(r.users, id)
	return nil
}