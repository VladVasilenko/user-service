package repository

import (
	"database/sql"
)

type PostgresRepository struct {
	db *sql.DB
}

// Имя функции ДОЛЖНО быть с большой буквы New...
func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{
		db: db,
	}
}

func (r *PostgresRepository) GetByID(id int) (string, bool, error) {
	var name string

	err := r.db.QueryRow("SELECT name FROM users WHERE id = $1", id).Scan(&name)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", false, nil
		}
		return "", false, err
	}

	return name, true, nil
}
