package api

import (
	"github.com/go-chi/chi/v5"
	chiMW "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
	"github.com/pythoninja/go-redirect/internal/config"
	mw "github.com/pythoninja/go-redirect/internal/server/middleware"
	"github.com/pythoninja/go-redirect/internal/server/route"
	"github.com/pythoninja/go-redirect/internal/storage"
	"net/http"
	"time"
)

type handler struct {
	config *config.Application
	store  *storage.Storage
}

func Router(cfg *config.Application, store *storage.Storage) http.Handler {
	handler := &handler{config: cfg, store: store}
	routes := route.New()

	router := chi.NewRouter()
	router.Use(mw.LogRequests)
	router.Use(httprate.LimitByRealIP(100, 1*time.Minute))
	router.Use(chiMW.RedirectSlashes)

	router.NotFound(handler.notFoundHandler)
	router.MethodNotAllowed(handler.methodNotAllowedHandler)

	// Main router for /
	router.Get(routes.Redirect, handler.linkRedirectHandler)
	//router.New("/panic", func(http.ResponseWriter, *http.Request) { panic("foo") })

	// API router for /v1
	apiRouter := chi.NewRouter()
	apiRouter.Get(routes.ApiHealtcheck, handler.healthcheckHandler)
	apiRouter.Get(routes.ApiListLinks, handler.listLinksHandler)
	apiRouter.Get(routes.ApiShowLink, handler.showLinkHandler)

	// Mount API router to the main router
	router.Mount(routes.ApiPath, apiRouter)

	return router
}
