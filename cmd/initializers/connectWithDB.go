package initializers

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

var DB *sql.DB

// Database settings

// TODO: pobrać zmienne z .env
const (
	host     = "192.168.1.28"
	port     = 5432 // Default port
	user     = "postgres"
	password = "poqfis-6deqbo-qoqNyr"
	dbname   = "fiber_demo"
)

func Connect() error {
	var err error
	DB, err = sql.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname))
	if err != nil {
		return err
	}
	if err = DB.Ping(); err != nil {
		return err
	}
	return nil
}

func GetDB() *sql.DB {
	return DB
}
