package app

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

func NewApp() (*App, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	db, err := sql.Open("pgx", cfg.DBString)
	if err != nil {
		return nil, fmt.Errorf("db open error: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(1 * time.Minute)

	pingCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(pingCtx); err != nil {
		db.Close()
		return nil, fmt.Errorf("db ping error: %w", err)
	}

	repo := repository.NewPostgresRepository(db)
	userService := service.NewUserService(repo)
	handler := userHTTP.NewUserHandler(userService)

	router := userHTTP.InitRoutes(userHTTP.Handlers{User: handler}, db)

	serverAddr := fmt.Sprintf(":%s", cfg.ServerPort)
	httpServer := &http.Server{
		Addr:              serverAddr,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	return &App{
		cfg:        cfg,
		db:         db,
		httpServer: httpServer,
	}, nil
}

func (a *App) Run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	errCh := make(chan error, 1)

	go func() {
		log.Printf("[START] Микросервис запущен на порту %s", a.cfg.ServerPort)
		if err := a.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	select {
	case err := <-errCh:
		return fmt.Errorf("http server error: %w", err)
	case <-ctx.Done():
		log.Println("[SHUTDOWN] Получен сигнал завершения, начинаем graceful shutdown")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := a.httpServer.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("http server shutdown error: %w", err)
	}

	log.Println("[SHUTDOWN] HTTP-сервер остановлен")
	return nil
}

func (a *App) Close() {
	if a.db == nil {
		return
	}
	if err := a.db.Close(); err != nil {
		log.Printf("[SHUTDOWN] Ошибка закрытия БД: %v", err)
		return
	}
	log.Println("[SHUTDOWN] Соединение с базой данных закрыто")
}
