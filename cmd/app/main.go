package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/machayka/mail-service/cmd/initializers"
)

func init() {
	initializers.LoadEnvVariables()
}

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	err := app.Listen(os.Getenv("PORT"))
	if err != nil {
		log.Fatalf("Can't start the app: ", err)
	}
}
