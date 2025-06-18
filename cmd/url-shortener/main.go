package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"url_shortener/internal/config"
	"url_shortener/internal/http-server/handlers/delete"
	"url_shortener/internal/http-server/handlers/redirect"
	"url_shortener/internal/http-server/handlers/url/save"
	"url_shortener/internal/lib/logger/sl"
	"url_shortener/internal/storage/sqlite"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

const (
	envProd  = "prod"
	envLocal = "local"
	envDev   = "dev"
)

func main() {

	cfg := config.MustLoad()
	ctx, cancel := context.WithCancel(context.Background())
	log := setupLogger(cfg.Env)

	log.Info("Starting app", slog.String("env", cfg.Env))
	log.Debug("Debug messages are enabled")

	// ssoClient, err := ssogrpc.New(context.Background(), log, cfg.Clients.SSO.Address,
	// 	cfg.Clients.SSO.Timeout, cfg.Clients.SSO.RetriesCount)
	// if err != nil {
	// 	log.Error("failed to init sso client", sl.Err(err))
	// 	os.Exit(1)
	// }

	// isAdm, err := ssoClient.IsAdmin(context.Background(), 1)
	// logg := log.With("IsAdmin?", isAdm)
	// logg.Info("isAdmin: 1?")

	// fmt.Println(isAdm)

	// connStr := fmt.Sprintf("user=%s dbname=%s password=%s host=%s port=%s sslmode=%s", cfg.PostgresDB.User, cfg.Dbname, cfg.PostgresDB.Password, cfg.Host, cfg.Port, cfg.Sslmode)
	// storage, err := psql.New(connStr)

	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("Error create postgresql: %s\n", sl.Err(err))
		os.Exit(1)
	}

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)

	router.Route("/url", func(r chi.Router) {
		r.Use(middleware.BasicAuth("url-shortener", map[string]string{
			cfg.HTTPServer.User: cfg.HTTPServer.Password,
		}))

		r.Post("/", save.New(ctx, log, storage))
		r.Delete("/{alias}", delete.New(ctx, log, storage))
		r.Get("/{alias}", redirect.New(ctx, log, storage))
		// TODO: DELETE /{id}

	})

	router.Get("/{alias}", redirect.New(ctx, log, storage))

	//middleware (при обработке каждого запроса - выполняется цепочка handler-ов, например авторизация)
	log.Info("starting server", slog.String("addres", cfg.Address))

	srv := &http.Server{
		Addr:    cfg.Address,
		Handler: router,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}
	cancel()

	log.Error("server stopped")

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
