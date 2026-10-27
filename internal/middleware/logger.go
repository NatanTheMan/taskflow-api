// Package middleware
package middleware

import (
	"log/slog"
	"net/http"
	"time"

	chimw "github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func Logger(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			wrapper := chimw.NewWrapResponseWriter(w, r.ProtoMajor)

			start := time.Now()
			next.ServeHTTP(wrapper, r)
			latency := time.Since(start).Milliseconds()

			route := chi.RouteContext(r.Context()).RoutePattern()
			args := []any{
				"method", r.Method,
				"status", wrapper.Status(),
				"route", route,
				"latency", latency,
				"request_id", chimw.GetReqID(r.Context()),
			}

			msg := "HTTP request"
			if route == "/health" || route == "/ready" {
				logger.Debug(msg, args...)
				return
			}

			switch wrapper.Status() / 100 {
			case 5:
				logger.Error(msg, args...)
			case 4:
				logger.Warn(msg, args...)
			default:
				logger.Info(msg, args...)
			}
		})
	}
}
