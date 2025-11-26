package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	formHandler "github.com/machayka/mail-service/cmd/handlers/form"
	"github.com/machayka/mail-service/cmd/initializers"
	formRepo "github.com/machayka/mail-service/cmd/repositories/form"
	formService "github.com/machayka/mail-service/cmd/services/form"
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
	db := initializers.GetDB()

	defer func() {
		if err := db.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	fRepo := formRepo.NewFormRepository(db)
	fService := formService.NewFormService(fRepo)
	fHandler := formHandler.NewFormHandler(fService)

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
