// Package initializers wywoływane są w main.go na początku kodu
package initializers

import (
	"errors"

	"github.com/joho/godotenv"
)

func LoadEnvVariables() error {
	err := godotenv.Load()
	if err != nil {
		return errors.New("error loading .env file")
	}
	return nil
}
