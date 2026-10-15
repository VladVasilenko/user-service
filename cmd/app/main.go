package main

import (
	"log"

	"user-service/internal/app"
)

func main() {
	application, err := app.NewApp()
	if err != nil {
		log.Fatalf("Критическая ошибка инициализации приложения: %v", err)
	}
	defer application.Close()

	if err := application.Run(); err != nil {
		log.Fatalf("Микросервис аварийно остановлен: %v", err)
	}
}
