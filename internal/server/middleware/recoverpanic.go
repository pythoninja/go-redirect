package middleware

import (
	"fmt"
	"github.com/pythoninja/go-redirect/internal/server/response/json"
	"net/http"
)

func recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				json.ServerError(w, r, fmt.Errorf("%s", err)) //nolint: err113
			}
		}()

		next.ServeHTTP(w, r)
	})
}
