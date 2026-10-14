package main

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"github.com/logan-dev-x/taskflow-api/internal/config"
)

func main() {
	godotenv.Load(".env")
	_, err := config.Load()
	if err != nil {
		slog.Error("invalid configuration", "error", err)
		os.Exit(1)
	}
}
