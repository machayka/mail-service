// Package form zajmują się operacjami na bazie danych
package form

import (
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

var (
	ErrFormNotFound      = errors.New("not found")
	ErrFormAlreadyExists = errors.New("uuid already in database")
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
		"SELECT id, email FROM forms WHERE id = $1",
		id,
	).Scan(&form.ID, &form.Email)

	if err == sql.ErrNoRows {
		return nil, ErrFormNotFound
	}
	return &form, err
}

func (r *Repository) CreateNewForm(form *Form) error {
	err := r.db.QueryRow("INSERT INTO forms (id, email) VALUES ($1, $2) RETURNING id, email", form.ID, form.Email).Scan(&form.ID, &form.Email)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == "23505" {
				return ErrFormAlreadyExists
			}
		}
	}
	return err
}
