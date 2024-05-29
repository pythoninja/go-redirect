package api

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	chiMW "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
	"github.com/pythoninja/go-redirect/internal/config"
	mw "github.com/pythoninja/go-redirect/internal/server/middleware"
	"github.com/pythoninja/go-redirect/internal/storage"
	"net/http"
	"time"
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
	showLinkRoute     = fmt.Sprintf("%s/link/{id}", basePath)
	linkRedirectRoute = "/{alias}"
)

func Routes(cfg *config.Application, store *storage.Storage) http.Handler {
	handler := &handler{config: cfg, store: store}

	router := chi.NewRouter()
	router.Use(mw.LogRequests)
	router.Use(httprate.LimitByRealIP(100, 1*time.Minute))
	router.Use(chiMW.RedirectSlashes)

	router.NotFound(handler.notFoundHandler)
	router.MethodNotAllowed(handler.methodNotAllowedHandler)

	router.Get(healthcheckRoute, handler.healthcheckHandler)
	router.Get(linkRedirectRoute, handler.linkRedirectHandler)
	router.Get(listLinksRoute, handler.listLinksHandler)
	router.Get(showLinkRoute, handler.showLinkHandler)
	//router.Get("/panic", func(http.ResponseWriter, *http.Request) { panic("foo") })

	return router
}
