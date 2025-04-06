package config

import "github.com/joho/godotenv"

type Config struct {
	Secret string
	Port   string // Default value is ":8080"
}

// New returns the new Config instance.
func New(secret string, port string) *Config {
	c := &Config{
		Secret: secret,
		Port:   port,
	}

	return c
}

func (c *Config) Load() error {
	err := godotenv.Load()
	if err != nil {
		return err
	}

	return nil
}
