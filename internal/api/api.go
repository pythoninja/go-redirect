package api

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/pythoninja/go-redirect/internal/config"
	"github.com/pythoninja/go-redirect/internal/storage"
	"log/slog"
	"net/http"
)

type handler struct {
	config *config.Application
	store  *storage.Storage
}

var (
	apiVersion       = 1
	basePath         = fmt.Sprintf("/v%d", apiVersion)
	healthcheckRoute = fmt.Sprintf("%s/healthcheck", basePath)
	//listLinksRoute    = fmt.Sprintf("%s/links", basePath)
	//linkInfoRoute     = fmt.Sprintf("%s/link/:id", basePath)
	//linkRedirectRoute = fmt.Sprintf("%s/link/:id", basePath)
)

func Routes(cfg *config.Application, store *storage.Storage) http.Handler {
	handler := &handler{config: cfg, store: store}

	router := httprouter.New()
	router.NotFound = http.HandlerFunc(handler.notFoundHandler)
	router.MethodNotAllowed = http.HandlerFunc(handler.methodNotAllowedHandler)

	slog.Debug("initialize route", "route", healthcheckRoute)
	//slog.Debug("Initialize route", "route", listLinksRoute)
	//slog.Debug("Initialize route", "route", linkInfoRoute)
	//slog.Debug("Initialize route", "route", linkRedirectRoute)

	router.HandlerFunc(http.MethodGet, healthcheckRoute, handler.healthcheckHandler)
	//router.HandlerFunc(http.MethodGet, listLinksRoute, listLinksHandler)
	//router.HandlerFunc(http.MethodGet, linkInfoRoute, linkInfoHandler)
	//router.HandlerFunc(http.MethodGet, linkRedirectRoute, linkRedirectHandler)

	return handler.logRequests(router)
}
