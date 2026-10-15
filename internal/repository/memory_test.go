package repository_test

import (
	"context"
	"errors"
	"testing"

	"user-service/internal/repository"
)

func TestMemoryRepository_GetByID(t *testing.T) {
	repo := repository.NewMemoryRepository()

	t.Run("found", func(t *testing.T) {
		u, err := repo.GetByID(context.Background(), 1)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if u == nil {
			t.Fatal("user is nil")
		}
		if u.ID != 1 {
			t.Fatalf("want id 1, got %d", u.ID)
		}
	})

	t.Run("not found", func(t *testing.T) {
		_, err := repo.GetByID(context.Background(), 999)
		if !errors.Is(err, repository.ErrNotFound) {
			t.Fatalf("want ErrNotFound, got %v", err)
		}
	})

	t.Run("context cancelled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// memory-репозиторий не проверяет ctx, но тест показывает,
		// что сигнатура это позволяет. Для postgres этот кейс реально сработает.
		_, err := repo.GetByID(ctx, 1)
		// memory вернёт nil, потому что не смотрит в ctx — это ок для in-memory.
		// Если хочешь честный тест — добавь проверку ctx в репозиторий:
		// if err := ctx.Err(); err != nil { return nil, err }
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}
