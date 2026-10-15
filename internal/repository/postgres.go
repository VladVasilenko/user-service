package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"user-service/internal/model"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) GetByID(ctx context.Context, id int64) (*model.User, error) {
	var u model.User
	err := r.db.QueryRowContext(ctx,
		"SELECT id, email, name FROM users WHERE id = $1", id,
	).Scan(&u.ID, &u.Email, &u.Name)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("query user %d: %w", id, err)
	}
	return &u, nil
}
