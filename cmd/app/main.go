package main

import (
	"log"

	"user-service/internal/app"
)

func main() {
	// 1. Собираем приложение (внутри отработает весь наш DI конвейер)
	application, err := app.NewApp()
	if err != nil {
		log.Fatalf("Критическая ошибка инициализации приложения: %v", err)
	}
	// Гарантируем, что когда функция main будет завершаться, база данных закроется
	defer application.Close()

	// 2. Включаем зажигание — запускаем HTTP-сервер микросервиса
	if err := application.Run(); err != nil {
		log.Fatalf("Микросервис аварийно остановлен: %v", err)
	}
}
