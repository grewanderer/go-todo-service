package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	"go-todo-service/internal/config"
	"go-todo-service/internal/handlers"
	"go-todo-service/internal/repository/postgres"
	authsvc "go-todo-service/internal/service/auth"
	tasksrv "go-todo-service/internal/service/task"
	"go-todo-service/pkg/logger"
)

func main() {
	log := logger.New(os.Stdout)

	cfg, err := config.Load()
	if err != nil {
		log.Error("failed to load config", map[string]any{"error": err.Error()})
		os.Exit(1)
	}

	db, err := setupDatabase(cfg)
	if err != nil {
		log.Error("database connection failed", map[string]any{"error": err.Error()})
		os.Exit(1)
	}
	defer db.Close()

	userRepo := postgres.NewUserRepository(db)
	taskRepo := postgres.NewTaskRepository(db)

	authService := authsvc.New(userRepo, cfg.JWTSecret, cfg.JWTTTL)
	taskService := tasksrv.New(taskRepo)

	authHandler := handlers.NewAuthHandler(authService, log)
	taskHandler := handlers.NewTaskHandler(taskService, log)
	authMiddleware := handlers.NewAuthMiddleware(cfg.JWTSecret, log)

	router := handlers.NewRouter(authHandler, taskHandler, authMiddleware, log)

	server := &http.Server{
		Addr:         ":" + cfg.ServerPort,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Error("server shutdown error", map[string]any{"error": err.Error()})
		}
	}()

	log.Info("server starting", map[string]any{"port": cfg.ServerPort})
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Error("server failed", map[string]any{"error": err.Error()})
		os.Exit(1)
	}

	log.Info("server stopped", map[string]any{})
}

func setupDatabase(cfg config.Config) (*sql.DB, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
		cfg.DBSSLMode,
	)
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Hour)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}
