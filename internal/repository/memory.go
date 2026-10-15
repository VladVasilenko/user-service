package repository

import (
	"context"
	"sync"

	"user-service/internal/model"
)

type MemoryRepository struct {
	mu   sync.RWMutex
	data map[int64]*model.User
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		data: map[int64]*model.User{
			1: {ID: 1, Name: "Алексей", Email: "alex@example.com"},
			2: {ID: 2, Name: "Иван", Email: "ivan@example.com"},
			3: {ID: 3, Name: "Мария", Email: "maria@example.com"},
		},
	}
}

func (r *MemoryRepository) GetByID(ctx context.Context, id int64) (*model.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	u, ok := r.data[id]
	if !ok {
		return nil, ErrNotFound
	}
	return u, nil
}
