package main

import (
	"fmt"
	"log/slog"
	"os"
	"url_shortener/internal/config"
	"url_shortener/internal/storage/psql"
)

const (
	envProd  = "prod"
	envLocal = "local"
	envDev   = "dev"
)

func main() {

	// os.Setenv("CONFIG_PATH", "./config/local.yaml")
	// fmt.Println("get pathvar: ", os.Getenv("CONFIG_PATH"))

	// TODO: init config - cleanenv
	cfg := config.MustLoad()
	fmt.Println(cfg)

	log := setupLogger(cfg.Env)
	log.Info("Starting app", slog.String("env", cfg.Env))

	log.Debug("Debug messages are enabled")

	//TODO: init logger - slog

	connStr := "user=postgres dbname=postgres password=1233 host=localhost port=5432 sslmode=disable"
	stor, err := psql.New(connStr)
	if err != nil {
		// log.Error("Error create postgresql: %s\n", err)
		fmt.Println(err)
	}
	fmt.Println(stor)

	//TODO: init storage

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
