package service_test

import (
	"context"
	"errors"
	"testing"

	"user-service/internal/model"
	"user-service/internal/repository"
	"user-service/internal/service"
)

type mockRepo struct {
	user *model.User
	err  error
}

func (m *mockRepo) GetByID(ctx context.Context, id int64) (*model.User, error) {
	return m.user, m.err
}

func TestUserService_GetUserProfile(t *testing.T) {
	t.Run("invalid id", func(t *testing.T) {
		svc := service.NewUserService(&mockRepo{})

		_, err := svc.GetUserProfile(context.Background(), 0)
		if !errors.Is(err, service.ErrInvalidID) {
			t.Fatalf("want ErrInvalidID, got %v", err)
		}
	})

	t.Run("not found", func(t *testing.T) {
		svc := service.NewUserService(&mockRepo{err: repository.ErrNotFound})

		_, err := svc.GetUserProfile(context.Background(), 1)
		if !errors.Is(err, service.ErrUserNotFound) {
			t.Fatalf("want ErrUserNotFound, got %v", err)
		}
	})

	t.Run("success", func(t *testing.T) {
		want := &model.User{ID: 1, Name: "Алексей"}
		svc := service.NewUserService(&mockRepo{user: want})

		got, err := svc.GetUserProfile(context.Background(), 1)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.ID != want.ID || got.Name != want.Name {
			t.Fatalf("want %+v, got %+v", want, got)
		}
	})
}
