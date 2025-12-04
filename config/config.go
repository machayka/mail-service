// Package config zajmuje się zabezpieczeniem zmiennych z pliku .env
package config

import (
	"fmt"
	"log"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	DB struct {
		Host     string `env:"DB_HOST,required"`
		Port     int    `env:"DB_PORT,required"`
		User     string `env:"DB_USER,required"`
		Password string `env:"DB_PASSWORD,required"`
		Name     string `env:"DB_NAME,required"`
	}

	Server struct {
		Port string `env:"PORT" envDefault:":420"`
	}

	SMTP struct {
		User string `env:"SMTP_USER,required"`
		Pass string `env:"SMTP_PASS,required"`
	}

	Stripe struct {
		Key           string `env:"STRIPE_KEY,required"`
		PriceID       string `env:"PRICE_ID,required"`
		Domain        string `env:"DOMAIN,required"`
		WebhookSecret string `env:"WEBHOOK_SECRET,required"`
		PortalLink    string `env:"PORTAL_LINK,required"`
	}
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	cfg := &Config{}

	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	if len(cfg.Server.Port) > 0 && cfg.Server.Port[0] != ':' {
		cfg.Server.Port = ":" + cfg.Server.Port
	}

	return cfg, nil
}

func (c *Config) ConnectionString() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		c.DB.Host,
		c.DB.Port,
		c.DB.User,
		c.DB.Password,
		c.DB.Name,
	)
}
