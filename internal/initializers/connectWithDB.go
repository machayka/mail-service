// Package initializers jest odpowiedzialne za funkcje przy starcie funkcji main()
package initializers

import (
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/machayka/mail-service/config"
)

func Connect(cfg *config.Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.ConnectionString())
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
