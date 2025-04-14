package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	secret string
	port   string // Default value - ":8080"
}

func New() *Config {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	secret := os.Getenv("SECRET")
	port := os.Getenv("PORT")

	if port == "" {
		port = ":8080"
	}

	if secret == "" {
		log.Fatal("SECRET is not set")
	}

	c := &Config{
		secret: secret,
		port:   port,
	}

	return c
}

func (c *Config) GetSecret() string {
	return c.secret
}

func (c *Config) GetPort() string {
	return c.port
}
