package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/machayka/mail-service/internal/form"
	"github.com/machayka/mail-service/internal/initializers"
)

func init() {
	if err := initializers.LoadEnvVariables(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	db, err := initializers.Connect()
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	fRepo := form.NewRepository(db)
	fService := form.NewService(fRepo)
	fHandler := form.NewHandler(fService)

	app := fiber.New()
	//	app.Use(recover.New())

	app.Post("/forms/:id", fHandler.FormHandler)

	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("missing PORT env")
	}
	if string(port[0]) != ":" {
		port = ":" + port
	}
	if err := app.Listen(port); err != nil {
		log.Fatal(err)
	}
}
