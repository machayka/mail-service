// Package form zajmują się operacjami na bazie danych
package form

import (
	"database/sql"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetByID(id string) (*Form, error) {
	var form Form
	err := r.db.QueryRow(
		"SELECT id, email, created_at FROM forms WHERE id = $1",
		id,
	).Scan(&form.ID, &form.Email, &form.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, ErrFormNotFound
	}
	return &form, err
}
