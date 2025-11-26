// Package form zajmują się operacjami na bazie danych
package form

import (
	"database/sql"
)

type FormRepository struct {
	db *sql.DB
}

func NewFormRepository(db *sql.DB) *FormRepository {
	return &FormRepository{db: db}
}
