package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/pythoninja/go-redirect/internal/config"
	"github.com/pythoninja/go-redirect/internal/server/middleware"
	"github.com/pythoninja/go-redirect/internal/server/route"
	"github.com/pythoninja/go-redirect/internal/storage"
	"log/slog"
	"net/http"
)

type handler struct {
	config *config.Application
	store  *storage.Storage
}

func Router(cfg *config.Application, store *storage.Storage) http.Handler {
	h := &handler{config: cfg, store: store}
	r := route.Configure()
	mw := middleware.Configure()

	router := chi.NewRouter()
	router.Use(mw.LogRequests)
	router.Use(mw.GlobalRateLimiter)
	router.Use(mw.RedirectSlashes)

	router.NotFound(h.notFoundHandler)
	router.MethodNotAllowed(h.methodNotAllowedHandler)

	// Main router for /
	router.Get(r.Redirect, h.linkRedirectHandler)
	slog.Info("registered new route", slog.Any("path", r.Redirect), slog.Any("method", "GET"))

	//router.Configure("/panic", func(http.ResponseWriter, *http.Request) { panic("foo") })

	// API router for /v1
	apiRouter := chi.NewRouter()
	apiRouter.Get(r.ApiHealtcheck, h.healthcheckHandler)
	slog.Info("registered new route", slog.Any("path", r.ApiHealtcheck), slog.Any("method", "GET"))

	apiRouter.Get(r.ApiListLinks, h.listLinksHandler)
	slog.Info("registered new route", slog.Any("path", r.ApiListLinks), slog.Any("method", "GET"))

	apiRouter.Get(r.ApiShowLink, h.showLinkHandler)
	slog.Info("registered new route", slog.Any("path", r.ApiShowLink), slog.Any("method", "GET"))

	// Mount API router to the main router
	router.Mount(r.ApiPath, apiRouter)

	return router
}
