APP_NAME=user-service
BINARY_PATH=./bin/$(APP_NAME)
MAIN_PATH=./cmd/app/main.go
# Переменная с DSN для мигратора (подставляем те же доступы, что и в .env)
DB_DSN="postgres://postgres:mysecretpassword@localhost:5432/users_db?sslmode=disable"

.PHONY: up down restart build run clean migrate-up migrate-down

up:
	docker compose up -d

down:
	docker compose down

restart:
	docker compose down && docker compose up -d

build:
	go build -o $(BINARY_PATH) $(MAIN_PATH)

run:
	go run $(MAIN_PATH)

clean:
	rm -rf ./bin

# --- СЕКЦИЯ МИГРАЦИЙ ---

# Накат миграций (Аналог php artisan migrate)
migrate-up:
	docker run --rm -v $(shell pwd)/migrations:/migrations --network host migrate/migrate -path=/migrations -database $(DB_DSN) up

# Откат всех миграций (Аналог php artisan migrate:rollback)
migrate-down:
	docker run --rm -v $(shell pwd)/migrations:/migrations --network host migrate/migrate -path=/migrations -database $(DB_DSN) down -all
