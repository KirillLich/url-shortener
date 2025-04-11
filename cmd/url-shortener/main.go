package main

import (
	"log"
	"log/slog"
	"os"
	"url-shortener/internal/config"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/storage/sqlite"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)
	log.Info("starting url-shortener", slog.String("env", cfg.Env))

	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}

	// _, err = storage.SaveURL("https://google.com", "google")
	// if err != nil {
	// 	log.Error("failed to save url", sl.Err(err))
	// 	os.Exit(1)
	// }

	myURL, err := storage.GetURL("google")
	if err != nil {
		log.Error("failed to get url", sl.Err(err))
		os.Exit(1)
	}

	log.Debug("url taken", slog.String("url", myURL))

	_ = storage

	//TODO: init router

	//TODO: run server
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = setupConsoleLogger(&slog.HandlerOptions{Level: slog.LevelDebug})
	case envDev:
		log = setupJSONLogger("../../logs/app.json", &slog.HandlerOptions{Level: slog.LevelDebug})
	case envProd:
		log = setupJSONLogger("../../logs/app.json", &slog.HandlerOptions{Level: slog.LevelInfo})
	}

	return log
}

func setupJSONLogger(logFilePath string, opts *slog.HandlerOptions) *slog.Logger {
	// Открываем файл для записи логов (создаем или перезаписываем)
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	// Создаем JSON-обработчик, пишущий в файл
	handler := slog.NewJSONHandler(logFile, opts)

	return slog.New(handler)
}

func setupConsoleLogger(opts *slog.HandlerOptions) *slog.Logger {
	handler := slog.NewTextHandler(os.Stdout, opts)

	return slog.New(handler)
}
