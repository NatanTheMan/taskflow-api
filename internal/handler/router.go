package handler

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(db Pinger) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)

	health := NewHealthHandler(db)
	r.Get("/health", health.Live)
	r.Get("/ready", health.Ready)

	return r
}
