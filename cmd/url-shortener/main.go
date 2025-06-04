package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"url_shortener/internal/config"
	"url_shortener/internal/http-server/handlers/redirect"
	"url_shortener/internal/http-server/handlers/url/save"
	"url_shortener/internal/lib/logger/sl"
	"url_shortener/internal/storage/psql"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

const (
	envProd  = "prod"
	envLocal = "local"
	envDev   = "dev"
)

func main() {

	// os.Setenv("CONFIG_PATH", "./config/local.yaml")
	// os.Setenv("CGO_ENABLED", "1")
	// fmt.Println("get pathvar: ", os.Getenv("CONFIG_PATH"))

	// TODO: init config - cleanenv
	cfg := config.MustLoad()
	fmt.Println(cfg)

	//TODO: init logger - slog
	log := setupLogger(cfg.Env)
	log.Info("Starting app", slog.String("env", cfg.Env))
	log.Debug("Debug messages are enabled")

	//TODO: init storage
	// connStr := "user=postgres dbname=postgres password=1233 host=localhost port=5432 sslmode=disable"

	connStr := fmt.Sprintf("user=%s dbname=%s password=%s host=%s port=%s sslmode=%s", cfg.PostgresDB.User, cfg.Dbname, cfg.PostgresDB.Password, cfg.Host, cfg.Port, cfg.Sslmode)
	storage, err := psql.New(connStr)
	if err != nil {
		log.Error("Error create postgresql: %s\n", sl.Err(err))
		os.Exit(1)
	}

	str, err := storage.GetURL("g")
	if err != nil {
		log.Error("URL not found: %s\n", sl.Err(err))
	}

	// slog.Info("URL: %s", slog.StringValue(str))
	fmt.Println(str)

	//TODO: init router - chi, "chi render"

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)

	router.Route("/url", func(r chi.Router) {
		r.Use(middleware.BasicAuth("url-shortener", map[string]string{
			cfg.HTTPServer.User: cfg.HTTPServer.Password,
		}))

		r.Post("/", save.New(log, storage))

		// TODO: DELETE /{id}

	})

	router.Get("/{alias}", redirect.New(log, storage))
	//middleware (при обработке каждого запроса - выполняется цепочка handler-ов, например авторизация)
	log.Info("starting server", slog.String("addres", cfg.Address))

	srv := &http.Server{
		Addr:    cfg.Address,
		Handler: router,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}
	log.Error("server stopped")

	//TODO: run server

}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log

}
