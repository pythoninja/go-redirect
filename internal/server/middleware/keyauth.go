package middleware

import (
	"github.com/pythoninja/go-redirect/internal/server/response/json"
	"log/slog"
	"net/http"
	"strings"
)

func keyAuth(apiKey string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Vary", "Authorization")

			authorizationHeader := r.Header.Get("Authorization")

			if authorizationHeader == "" {
				slog.Warn("blank authorization header")
				json.AuthenticationRequiredResponse(w, r)
				return
			}

			headerParts := strings.Split(authorizationHeader, " ")
			if len(headerParts) != 2 || headerParts[0] != "Key" {
				slog.Warn("invalid authorization header", slog.Any("header", authorizationHeader))
				json.InvalidAPIKeyResponse(w, r)
				return
			}

			key := headerParts[1]

			if key != apiKey {
				slog.Info("invalid api key provided", slog.Any("key", key))
				json.InvalidAPIKeyResponse(w, r)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
