package app

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"user-service/internal/config"
	"user-service/internal/repository"
	"user-service/internal/service"
	userHTTP "user-service/internal/transport/http"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type App struct {
	cfg        *config.Config
	db         *sql.DB
	httpServer *http.Server
}

// NewApp инициализирует все слои приложения, базы данных и воркеров
func NewApp() (*App, error) {
	// 1. Загружаем конфиг
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	// 2. Инициализируем БД
	db, err := sql.Open("pgx", cfg.DBString)
	if err != nil {
		return nil, fmt.Errorf("db open error: %w", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("db ping error: %w", err)
	}

	// 3. Собираем зависимости конвейера (DI)
	repo := repository.NewPostgresRepository(db)
	userService := service.NewUserService(repo)
	handler := userHTTP.NewUserHandler(userService)

	// 4. Собираем роуты
	router := userHTTP.InitRoutes(handler)

	// 5. Настраиваем сам HTTP-сервер
	serverAddr := fmt.Sprintf(":%s", cfg.ServerPort)
	httpServer := &http.Server{
		Addr:    serverAddr,
		Handler: router,
	}

	return &App{
		cfg:        cfg,
		db:         db,
		httpServer: httpServer,
	}, nil
}

// Run запускает скомпилированное приложение
func (a *App) Run() error {
	log.Printf("[START] Микросервис запущен на порту %s", a.cfg.ServerPort)
	return a.httpServer.ListenAndServe()
}

// Close безопасно закрывает критические ресурсы при остановке
func (a *App) Close() {
	if a.db != nil {
		_ = a.db.Close()
		log.Println("[SHUTDOWN] Соединение с базой данных закрыто")
	}
}
