package middleware

import (
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
	"net/http"
	"time"
)

type Middlewares struct {
	LogRequests       func(next http.Handler) http.Handler
	RecoverPanic      func(next http.Handler) http.Handler
	RedirectSlashes   func(next http.Handler) http.Handler
	GlobalRateLimiter func(next http.Handler) http.Handler
	Authorize         func(apiKey string) func(next http.Handler) http.Handler
}

func Configure() Middlewares {
	return Middlewares{
		LogRequests:       logRequests,
		RecoverPanic:      recoverPanic,
		RedirectSlashes:   middleware.RedirectSlashes,
		GlobalRateLimiter: httprate.LimitByRealIP(300, 1*time.Minute),
		Authorize:         keyAuth,
	}
}
