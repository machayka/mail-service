// Package form zajmują się operacjami na bazie danych
package form

import (
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

var (
	ErrNotFound          = errors.New("not found")
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
		return nil, ErrNotFound
	}
	return &form, err
}

func (r *Repository) CreateNewForm(formID, email, customerID, subscriptionID string) error {
	_, err := r.db.Exec("INSERT INTO forms (id, email, stripe_customer_id, stripe_subscription_id) VALUES ($1, $2, $3, $4)", formID, email, customerID, subscriptionID)
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

func (r *Repository) GetStripeCustomerID(email string) (string, error) {
	var stripeCustomerID string
	err := r.db.QueryRow("SELECT stripe_customer_id FROM forms WHERE email = $1", email).Scan(&stripeCustomerID)
	return stripeCustomerID, err
}

func (r *Repository) DeleteForm(subscriptionID string) error {
	_, err := r.db.Exec("DELETE FROM forms WHERE stripe_subscription_id = $1", subscriptionID)
	return err
}
