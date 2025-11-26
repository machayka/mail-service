package form

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/machayka/mail-service/cmd/models"
)

func (r *FormRepository) GetByID(id uuid.UUID) (*models.Form, error) {
	var form models.Form
	err := r.db.QueryRow(
		"SELECT id, email, created_at FROM forms WHERE id = $1",
		id,
	).Scan(&form.ID, &form.Email, &form.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &form, err
}
