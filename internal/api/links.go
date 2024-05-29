package api

import (
	"errors"
	"github.com/pythoninja/go-redirect/internal/server/response/json"
	"github.com/pythoninja/go-redirect/internal/storage"
	"github.com/pythoninja/go-redirect/internal/validator"
	"net/http"
	"net/url"
)

func (h *handler) listLinksHandler(w http.ResponseWriter, r *http.Request) {
	res, err := h.store.Links.GetAllLinks()
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
	res, err := h.store.Links.GetLinkById(id)
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

func (h *handler) linkRedirectHandler(w http.ResponseWriter, r *http.Request) {
	alias := readAliasParam(r)
	if alias == "" {
		json.NotFound(w, r)
		return
	}

	rawURL, err := h.store.Links.GetUrlByAlias(alias)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrRecordNotFound):
			json.NotFound(w, r)
		default:
			json.ServerError(w, r, err)
		}
		return
	}

	v := validator.New()
	parsedURL, err := url.ParseRequestURI(rawURL)
	if err != nil {
		json.ServerError(w, r, err)
		return
	}

	// This shouldn't ever happen, but we want to validate the URL before sending it to the client.
	// And send an HTTP 500 error with details if returned URL is not valid.
	if v.ValidateURL(parsedURL); !v.Valid() {
		json.ServerErrorWithDetails(w, r, v.Errors)
		return
	}

	err = h.store.Links.UpdateClicksByAlias(alias)
	if err != nil {
		json.ServerError(w, r, err)
	}

	http.Redirect(w, r, rawURL, http.StatusFound)
}
