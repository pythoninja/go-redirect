package api

import (
	"github.com/pythoninja/go-redirect/internal/server/response/json"
	"net/http"
)

func (h *handler) notFoundHandler(w http.ResponseWriter, r *http.Request) {
	json.NotFound(w, r)
}

func (h *handler) methodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	json.MethodNotAllowed(w, r)
}
