package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

func Load() *Config {
	_ = godotenv.Load()

	cfg := &Config{
		TelegramToken: os.Getenv("TELEGRAM_BOT_TOKEN"),
		APIURL:        os.Getenv("API_URL"),
	}

	if cfg.TelegramToken == "" || cfg.APIURL == "" {
		log.Fatal("Missing TELEGRAM_BOT_TOKEN or API_URL")
	}

	return cfg
}
