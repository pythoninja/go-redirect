package api

import (
	"log/slog"
	"net/http"
)

func (h *handler) logRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			ip        = r.RemoteAddr
			method    = r.Method
			userAgent = r.UserAgent()
			uri       = r.URL.RequestURI()
		)

		slog.Info("received request",
			slog.String("ip", ip),
			slog.String("method", method),
			slog.String("user-agent", userAgent),
			slog.String("uri", uri))

		next.ServeHTTP(w, r)
	})
}
