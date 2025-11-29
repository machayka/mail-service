// Package form jest odpowiedzialny za definicje structów
package form

import (
	"github.com/google/uuid"
)

type Form struct {
	ID    uuid.UUID `db:"id" form:"id"`
	Email string    `db:"email" form:"email"`
}

type FormData struct {
	Email   string `json:"email" form:"email"`
	Message string `json:"message" form:"message"`
}
