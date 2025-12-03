package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html/v2"
	"github.com/machayka/mail-service/config"
	"github.com/machayka/mail-service/internal/form"
	"github.com/machayka/mail-service/internal/initializers"
	"github.com/machayka/mail-service/internal/mail"
	"github.com/machayka/mail-service/internal/payments"
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

	mailSender := mail.NewService(cfg)
	paymentClient := payments.NewPaymentClient(cfg)

	fRepo := form.NewRepository(db)
	fService := form.NewService(fRepo, mailSender, paymentClient)
	fHandler := form.NewHandler(fService)

	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{Views: engine})
	app.Use(recover.New())

	app.Post("/submit/:id", fHandler.FormSubmit)
	app.Get("/add/:id", fHandler.NewForm)
	app.Post("/add", fHandler.AddForm)

	// TODO: to ":id" jest dla No Code customer portal: https://docs.stripe.com/customer-management/activate-no-code-customer-portal
	app.Get("/success/:id", fHandler.PaymentSuccess)
	app.Post("/webhook", fHandler.HandleWebhook(cfg))

	log.Fatal(app.Listen(cfg.Server.Port))
}
