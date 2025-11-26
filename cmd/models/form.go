// Package models jest odpowiedzialny za definicje structów
package models

import (
	"time"

	"github.com/google/uuid"
)

type Form struct {
	ID        uuid.UUID `db:"id"`
	Email     string    `db:"email"`
	CreatedAt time.Time `db:"created_at"`
}

type FormData struct {
	Email   string `json:"email" form:"email"`
	Message string `json:"message" form:"message"`
}
