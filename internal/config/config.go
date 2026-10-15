package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort string
	DBString   string
}

func LoadConfig() (*Config, error) {
	// .env опционален: в проде переменные уже в окружении
	_ = godotenv.Load(".env")

	cfg := &Config{
		ServerPort: getEnv("SERVER_PORT", "8080"),
		DBString:   os.Getenv("DB_DSN"),
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) validate() error {
	if c.DBString == "" {
		return errors.New("DB_DSN is required")
	}

	port, err := strconv.Atoi(c.ServerPort)
	if err != nil {
		return fmt.Errorf("SERVER_PORT must be a number: %w", err)
	}
	if port < 1 || port > 65535 {
		return fmt.Errorf("SERVER_PORT out of range: %d", port)
	}

	return nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
