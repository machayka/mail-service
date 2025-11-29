package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
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

	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{Views: engine})
	//	app.Use(recover.New())

	app.Post("/submit/:id", fHandler.FormHandler)
	app.Get("/add/:id", fHandler.NewForm)
	app.Post("/add", fHandler.AddForm)

	log.Fatal(app.Listen(cfg.Server.Port))
}
