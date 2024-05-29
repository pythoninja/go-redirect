package middleware

import (
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"time"
)

func LogRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		logEntry := slog.With(
			slog.String("ip", r.RemoteAddr),
			slog.String("method", r.Method),
			slog.String("user-agent", r.UserAgent()),
			slog.String("uri", r.URL.RequestURI()),
		)

		start := time.Now()

		next.ServeHTTP(ww, r)

		logEntry.Info("request completed",
			slog.Int("status", ww.Status()),
			slog.Int("bytes", ww.BytesWritten()),
			slog.String("duration", time.Since(start).String()),
		)
	})
}
