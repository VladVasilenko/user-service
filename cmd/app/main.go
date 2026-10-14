package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"user-service/internal/repository"
	userHTTP "user-service/internal/transport/http"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func RecoveryMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Паника перехвачена: %v", err)
				http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
			}
		}()
		next(w, r)
	}
}

func main() {
	// Строка подключения (DSN): указываем localhost, так как порт базы проброшен на наш ПК
	dsn := "postgres://postgres:mysecretpassword@localhost:5432/users_db?sslmode=disable"

	// 1. Открываем пул соединений с базой данных
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		panic(err)
	}
	defer db.Close() // Гарантируем закрытие пула при остановке сервиса

	// 2. Честно проверяем, доступен ли наш Docker-контейнер
	err = db.Ping()
	if err != nil {
		log.Fatalf("База данных в Docker недоступна: %v", err)
	}

	// 3. Инициализируем репозиторий Postgres вместо старого MemoryRepository
	repo := repository.NewPostgresRepository(db)

	// 4. Передаем его в хендлер (интерфейс хендлера остался прежним!)
	handler := userHTTP.NewUserHandler(repo)

	mux := http.NewServeMux()
	mux.HandleFunc("/user", RecoveryMiddleware(handler.GetPlayer))

	fmt.Println("Микросервис запущен на :8080 и подключен к Postgres в Docker!")

	if err := http.ListenAndServe(":8080", mux); err != nil {
		panic(err)
	}
}
