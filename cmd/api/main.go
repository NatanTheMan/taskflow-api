package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/logan-dev-x/taskflow-api/internal/config"
	"github.com/logan-dev-x/taskflow-api/internal/database"
	"github.com/logan-dev-x/taskflow-api/internal/handler"
)

func main() {
	godotenv.Load(".env")
	cfg, err := config.Load()
	if err != nil {
		slog.Error("invalid configuration", "error", err)
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	pool, err := database.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		println("no db")
	}
	defer pool.Close()

	server := &http.Server{
		Addr:              ":" + cfg.HTTPPort,
		Handler:           handler.NewRouter(pool),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	server.ListenAndServe()
}
