package api

import (
	"github.com/pythoninja/go-redirect/internal/server/response/json"
	"net/http"
)

func (h *handler) listLinksHandler(w http.ResponseWriter, r *http.Request) {
	res, err := h.store.GetAllLinks()
	if err != nil {
		json.ServerError(w, r, err)
		return
	}

	json.Ok(w, r, res)
}
