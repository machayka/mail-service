package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/machayka/mail-service/cmd/initializers"
)

func init() {
	if err := initializers.LoadEnvVariables(); err != nil {
		log.Fatal(err)
	}
	if err := initializers.Connect(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	app := fiber.New()
	db := initializers.GetDB()

	app.Get("/", func(c *fiber.Ctx) error {
		rows, err := db.Query("SELECT id, name, salary, age FROM employees order by id")
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}
		defer rows.Close()
		fmt.Println("Rows: ", rows)
		return c.SendString("Hello, World!")
	})

	// Potrzebuje handlera do sprawdzenia czy id w bazie istieje

	if err := app.Listen(os.Getenv("PORT")); err != nil {
		log.Fatal(err)
	}
}
