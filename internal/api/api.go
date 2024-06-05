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
	logAPIEntry := slog.With(slog.String("root", r.API.Path))

	// Main rootRouter for /
	rootRouter.Get(r.Root.Redirect, h.redirectLinkHandler)
	logRootEntry.Info("registered new route", slog.Any("path", r.Root.Redirect), slog.Any("method", "GET"))

	// API rootRouter for /v1
	apiRouter := chi.NewRouter()

	apiRouter.Get(r.API.Healtcheck, h.healthcheckHandler)
	logAPIEntry.Info("registered new route", slog.Any("path", r.API.Healtcheck), slog.Any("method", "GET"))

	apiRouter.With(mw.Authorize(cfg.Config.APISecretKey)).Get(r.API.ListLinks, h.listLinksHandler)
	logAPIEntry.Info("registered new route", slog.Any("path", r.API.ListLinks), slog.Any("method", "GET"))

	apiRouter.With(mw.Authorize(cfg.Config.APISecretKey)).Get(r.API.ShowLink, h.showLinkHandler)
	logAPIEntry.Info("registered new route", slog.Any("path", r.API.ShowLink), slog.Any("method", "GET"))

	apiRouter.With(mw.Authorize(cfg.Config.APISecretKey)).Post(r.API.AddLink, h.addLinkHandler)
	logAPIEntry.Info("registered new route", slog.Any("path", r.API.AddLink), slog.Any("method", "POST"))

	apiRouter.With(mw.Authorize(cfg.Config.APISecretKey)).Put(r.API.UpdateLink, h.updateLinkHandler)
	logAPIEntry.Info("registered new route", slog.Any("path", r.API.UpdateLink), slog.Any("method", "PUT"))

	apiRouter.With(mw.Authorize(cfg.Config.APISecretKey)).Delete(r.API.DeleteLink, h.deleteLinkHandler)
	logAPIEntry.Info("registered new route", slog.Any("path", r.API.DeleteLink), slog.Any("method", "DELETE"))

	// Mount API rootRouter to the main rootRouter
	rootRouter.Mount(r.API.Path, apiRouter)

	return rootRouter
}
