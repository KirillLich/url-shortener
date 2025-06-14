package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"
	"url-shortener/internal/config"
	"url-shortener/internal/http-server/handlers/redirect"
	"url-shortener/internal/http-server/handlers/url/save"
	mwLogger "url-shortener/internal/http-server/middleware/logger"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/storage/sqlite"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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

	router := chi.NewRouter()

	//middleware
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger) //built-in logger middleware
	router.Use(mwLogger.New(log)) //custome logger
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat) //specific for chi

	//authorization
	router.Route("/url", func(r chi.Router) {
		r.Use(middleware.BasicAuth("url-shortener", map[string]string{
			cfg.HTTPServer.User: cfg.HTTPServer.Password,
		}))
		r.Post("/", save.New(log, storage))
		//TODO: delete
	})

	router.Get("/{alias}", redirect.New(log, storage))

	log.Info("starting server", slog.String("adress", cfg.Address))

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to run server")
	}

	log.Error("server stopped")
	//TODO: run server
}

// TODO: add pretty logger
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
