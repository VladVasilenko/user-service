package service

import (
	"context"
	"errors"
	"fmt"

	"user-service/internal/model"
	"user-service/internal/repository"
)

var (
	ErrInvalidID    = errors.New("invalid user id")
	ErrUserNotFound = errors.New("user not found")
)

type UserRepository interface {
	GetByID(ctx context.Context, id int64) (*model.User, error)
}

type UserService struct {
	repo UserRepository
}

func NewUserService(r UserRepository) *UserService {
	return &UserService{repo: r}
}

func (s *UserService) GetUserProfile(ctx context.Context, id int64) (*model.User, error) {
	if id <= 0 {
		return nil, ErrInvalidID
	}

	u, err := s.repo.GetByID(ctx, id)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get user %d: %w", id, err)
	}

	return u, nil
}
