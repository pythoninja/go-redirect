package api

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/pythoninja/go-redirect/internal/config"
	"github.com/pythoninja/go-redirect/internal/server/middleware"
	"github.com/pythoninja/go-redirect/internal/storage"
	"log/slog"
	"net/http"
)

type handler struct {
	config *config.Application
	store  *storage.Storage
}

var (
	apiVersion        = 1
	basePath          = fmt.Sprintf("/v%d", apiVersion)
	healthcheckRoute  = fmt.Sprintf("%s/healthcheck", basePath)
	listLinksRoute    = fmt.Sprintf("%s/links", basePath)
	showLinkRoute     = fmt.Sprintf("%s/link/:id", basePath)
	linkRedirectRoute = "/r/:alias"
)

func Routes(cfg *config.Application, store *storage.Storage) http.Handler {
	handler := &handler{config: cfg, store: store}

	router := httprouter.New()
	router.NotFound = http.HandlerFunc(handler.notFoundHandler)
	router.MethodNotAllowed = http.HandlerFunc(handler.methodNotAllowedHandler)

	slog.Debug("initialize route", "route", healthcheckRoute)
	slog.Debug("initialize route", "route", listLinksRoute)
	slog.Debug("initialize route", "route", showLinkRoute)
	slog.Debug("initialize route", "route", linkRedirectRoute)

	router.HandlerFunc(http.MethodGet, healthcheckRoute, handler.healthcheckHandler)
	router.HandlerFunc(http.MethodGet, listLinksRoute, handler.listLinksHandler)
	router.HandlerFunc(http.MethodGet, showLinkRoute, handler.showLinkHandler)
	router.HandlerFunc(http.MethodGet, linkRedirectRoute, handler.linkRedirectHandler)

	return middleware.LogRequests(router)
}
