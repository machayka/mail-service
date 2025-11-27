// Package form zajmują się operacjami na bazie danych
package form

import (
	"database/sql"

	"github.com/google/uuid"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetByID(id uuid.UUID) (*Form, error) {
	var form Form
	err := r.db.QueryRow(
		"SELECT id, email, created_at FROM forms WHERE id = $1",
		id,
	).Scan(&form.ID, &form.Email, &form.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &form, err
}
