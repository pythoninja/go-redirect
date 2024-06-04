package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/pythoninja/go-redirect/internal/config"
	"github.com/pythoninja/go-redirect/internal/server/middleware"
	"github.com/pythoninja/go-redirect/internal/server/response/json"
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

	rootRouter := chi.NewRouter()
	rootRouter.Use(mw.LogRequests)
	rootRouter.Use(mw.RecoverPanic)

	if cfg.Config.EnableRateLimiter {
		rootRouter.Use(mw.GlobalRateLimiter)
	}

	rootRouter.Use(mw.RedirectSlashes)

	rootRouter.NotFound(json.NotFound)
	rootRouter.MethodNotAllowed(json.MethodNotAllowed)

	logRootEntry := slog.With(slog.String("root", r.Root.Path))
	logApiEntry := slog.With(slog.String("root", r.Api.Path))

	// Main rootRouter for /
	rootRouter.Get(r.Root.Redirect, h.redirectLinkHandler)
	logRootEntry.Info("registered new route", slog.Any("path", r.Root.Redirect), slog.Any("method", "GET"))

	// API rootRouter for /v1
	apiRouter := chi.NewRouter()

	apiRouter.Get(r.Api.Healtcheck, h.healthcheckHandler)
	logApiEntry.Info("registered new route", slog.Any("path", r.Api.Healtcheck), slog.Any("method", "GET"))

	apiRouter.With(mw.Authorize(cfg.Config.APISecretKey)).Get(r.Api.ListLinks, h.listLinksHandler)
	logApiEntry.Info("registered new route", slog.Any("path", r.Api.ListLinks), slog.Any("method", "GET"))

	apiRouter.With(mw.Authorize(cfg.Config.APISecretKey)).Get(r.Api.ShowLink, h.showLinkHandler)
	logApiEntry.Info("registered new route", slog.Any("path", r.Api.ShowLink), slog.Any("method", "GET"))

	apiRouter.With(mw.Authorize(cfg.Config.APISecretKey)).Post(r.Api.AddLink, h.addLinkHandler)
	logApiEntry.Info("registered new route", slog.Any("path", r.Api.AddLink), slog.Any("method", "POST"))

	apiRouter.With(mw.Authorize(cfg.Config.APISecretKey)).Put(r.Api.UpdateLink, h.updateLinkHandler)
	logApiEntry.Info("registered new route", slog.Any("path", r.Api.UpdateLink), slog.Any("method", "PUT"))

	apiRouter.With(mw.Authorize(cfg.Config.APISecretKey)).Delete(r.Api.DeleteLink, h.deleteLinkHandler)
	logApiEntry.Info("registered new route", slog.Any("path", r.Api.DeleteLink), slog.Any("method", "DELETE"))

	// Mount API rootRouter to the main rootRouter
	rootRouter.Mount(r.Api.Path, apiRouter)

	return rootRouter
}
