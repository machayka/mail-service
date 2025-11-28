package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/machayka/mail-service/config"
	"github.com/machayka/mail-service/internal/form"
	"github.com/machayka/mail-service/internal/initializers"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	db, err := initializers.Connect(cfg)
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

	log.Fatal(app.Listen(cfg.Server.Port))
}
