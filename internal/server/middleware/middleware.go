package middleware

import (
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
	"net/http"
	"time"
)

type Middlewares struct {
	LogRequests       func(next http.Handler) http.Handler
	RedirectSlashes   func(next http.Handler) http.Handler
	GlobalRateLimiter func(next http.Handler) http.Handler
}

func Configure() Middlewares {
	return Middlewares{
		LogRequests:       LogRequests,
		RedirectSlashes:   middleware.RedirectSlashes,
		GlobalRateLimiter: httprate.LimitByRealIP(100, 1*time.Minute),
	}
}
