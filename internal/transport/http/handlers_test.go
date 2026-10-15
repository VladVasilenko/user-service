package http_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"user-service/internal/model"
	"user-service/internal/service"
	userHTTP "user-service/internal/transport/http"
)

type mockService struct {
	user *model.User
	err  error
}

func (m *mockService) GetUserProfile(ctx context.Context, id int64) (*model.User, error) {
	return m.user, m.err
}

func TestUserHandler_GetUser(t *testing.T) {
	t.Run("not found", func(t *testing.T) {
		h := userHTTP.NewUserHandler(&mockService{err: service.ErrUserNotFound})

		req := httptest.NewRequest(http.MethodGet, "/users/999", nil)
		req.SetPathValue("id", "999")
		rec := httptest.NewRecorder()

		h.GetUser(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Fatalf("want 404, got %d", rec.Code)
		}
		if !strings.Contains(rec.Body.String(), "user not found") {
			t.Fatalf("unexpected body: %s", rec.Body.String())
		}
	})

	t.Run("success", func(t *testing.T) {
		h := userHTTP.NewUserHandler(&mockService{
			user: &model.User{ID: 1, Name: "Алексей"},
		})

		req := httptest.NewRequest(http.MethodGet, "/users/1", nil)
		req.SetPathValue("id", "1")
		rec := httptest.NewRecorder()

		h.GetUser(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("want 200, got %d", rec.Code)
		}
		if !strings.Contains(rec.Body.String(), "Алексей") {
			t.Fatalf("unexpected body: %s", rec.Body.String())
		}
	})
}
