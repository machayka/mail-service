package main

import (
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

	app.Listen(os.Getenv("PORT"))
}
