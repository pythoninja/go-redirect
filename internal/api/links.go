package api

import (
	"errors"
	"github.com/pythoninja/go-redirect/internal/server/response/json"
	"github.com/pythoninja/go-redirect/internal/storage"
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

func (h *handler) showLinkHandler(w http.ResponseWriter, r *http.Request) {
	id, err := readIdParam(r)
	if err != nil {
		json.NotFound(w, r)
		return
	}

	// Get info about link from database
	res, err := h.store.GetLinkById(int64(id))
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrRecordNotFound):
			json.NotFound(w, r)
		default:
			json.ServerError(w, r, err)
		}
		return
	}

	json.Ok(w, r, res)
}
