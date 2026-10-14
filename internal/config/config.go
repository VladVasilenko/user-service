package config

import (
	"flag"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort string
	DBString   string
	TgToken    string
}

func LoadConfig() (*Config, error) {
	// 1. Объявляем флаг запуска -config, по умолчанию он равен ".env"
	configPath := flag.String("config", ".env", "Путь к файлу конфигурации .env")
	flag.Parse() // Парсим флаги, переданные в терминале

	// 2. Загружаем файл по вытащенному пути
	// Теперь если файла нет — мы честно возвращаем ошибку наружу, а не игнорируем её через "_"
	err := godotenv.Overload(*configPath)
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		ServerPort: os.Getenv("SERVER_PORT"),
		DBString:   os.Getenv("DB_DSN"),
		TgToken:    os.Getenv("TG_BOT_TOKEN"),
	}

	if cfg.ServerPort == "" {
		cfg.ServerPort = "8080"
	}

	return cfg, nil
}
