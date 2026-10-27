package middleware_test

import (
	"log/slog"
	"net/http"
	"testing"

	"github.com/logan-dev-x/taskflow-api/internal/middleware"
)

func TestLogger(t *testing.T) {
	tests := []struct {
		name   string // description of this test case
		logger *slog.Logger
		want   func(http.Handler) http.Handler
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := middleware.Logger(tt.logger)
			if true {
				t.Errorf("Logger() = %v, want %v", got, tt.want)
			}
		})
	}
}
