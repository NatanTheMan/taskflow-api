package handler

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/logan-dev-x/taskflow-api/internal/middleware"
)

func NewRouter(db Pinger, logger *slog.Logger) http.Handler {
	r := chi.NewRouter()
	r.Use(chimw.RequestID)
	r.Use(middleware.Logger(logger))
	r.Use(chimw.Recoverer)

	health := NewHealthHandler(db)
	r.Get("/health", health.Live)
	r.Get("/ready", health.Ready)

	r.Route("/api/v1", func(r chi.Router) {
	})

	return r
}
