package main

import (
	"fmt"
	"log/slog"
	"os"
	"url_shortener/internal/config"
	"url_shortener/internal/lib/logger/sl"
	"url_shortener/internal/storage/psql"
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

	connStr := fmt.Sprintf("user=%s dbname=%s password=%s host=%s port=%s sslmode=%s", cfg.User, cfg.Dbname, cfg.Password, cfg.Host, cfg.Port, cfg.Sslmode)
	stor, err := psql.New(connStr)
	if err != nil {
		log.Error("Error create postgresql: %s\n", sl.Err(err))
		os.Exit(1)
	}
	_ = stor

	id, err := stor.SaveURL("yandex.ru", "r")
	if err != nil {
		log.Error("Error create postgresql: %s\n", sl.Err(err))
		os.Exit(1)
	}

	fmt.Println(id)

	//TODO: init router - chi, "chi render"

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
