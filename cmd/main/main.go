package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"log"
	"log/slog"
	"net/http"
	"os"
	"todoapp/internal/config"
	"todoapp/internal/http-server/handlers/task/deleteTask"
	"todoapp/internal/http-server/handlers/task/getTask"
	"todoapp/internal/http-server/handlers/task/save"
	"todoapp/internal/http-server/handlers/task/updateTask"
	myLogger "todoapp/internal/http-server/middleware/logger"
	"todoapp/internal/lib/logger/slogLib"
	"todoapp/internal/storage/sqlite"
)

const (
	envLocal = "local"
	envProd  = "prod"
)

func main() {
	initEnv()
	configPath, exists := os.LookupEnv("CONFIG_PATH")
	if !exists {
		log.Fatal("CONFIG_PATH is not set")
	}

	cfg := config.MustLoad(configPath)
	logger := setUpLogger(cfg.Env)
	logger.Debug("Debug messages are enabled")
	logger.Info("Starting app...")

	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		logger.Error("Cant initialize storage", slogLib.Err(err))
		os.Exit(1)
	}

	_ = storage

	router := chi.NewRouter()

	setUpMiddlewares(router, logger)
	setUpHandlers(router, logger, *storage)

	logger.Info("Starting server", slog.String("address", cfg.Address))

	server := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HttpServer.Timeout,
		WriteTimeout: cfg.HttpServer.Timeout,
		IdleTimeout:  cfg.HttpServer.IdleTimeout,
	}
	if err := server.ListenAndServe(); err != nil {
		logger.Error("Failed to start server")
	}
}

func setUpLogger(env string) *slog.Logger {
	var logger *slog.Logger
	switch env {
	case envLocal:
		logger = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		logger = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}
	return logger
}

func setUpHandlers(router *chi.Mux, logger *slog.Logger, storage sqlite.Storage) {
	router.Post("/task", save.New(logger, &storage))
	router.Delete("/task", deleteTask.New(logger, &storage))
	router.Patch("/task", updateTask.New(logger, &storage))
	router.Get("/task", getTask.New(logger, &storage))
}

func setUpMiddlewares(router *chi.Mux, logger *slog.Logger) {
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(myLogger.New(logger))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
}

func initEnv() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Cant find env file", err)
	}
}
