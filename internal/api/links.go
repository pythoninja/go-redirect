package api

import (
	"errors"
	"fmt"
	"github.com/pythoninja/go-redirect/internal/model"
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
		message := map[string]string{"link": fmt.Sprintf("must be a positive integer")}
		json.LinkNotFoundResponse(w, r, message)
		return
	}

	// Get info about link from database
	res, err := h.store.Links.GetLinkById(id)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrRecordNotFound):
			message := map[string]string{"link": fmt.Sprintf("id '%d' not found", id)}
			json.LinkNotFoundResponse(w, r, message)
		default:
			json.ServerError(w, r, err)
		}
		return
	}

	json.Ok(w, r, res)
}

func (h *handler) linkRedirectHandler(w http.ResponseWriter, r *http.Request) {
	alias := readAliasParam(r)

	rawURL, err := h.store.Links.GetUrlByAlias(alias)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrRecordNotFound):
			message := map[string]string{"alias": fmt.Sprintf("'%s' not found", alias)}
			json.LinkNotFoundResponse(w, r, message)
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
	// Send an HTTP 500 error and log the details if returned URL is not valid.
	if validator.ValidateURL(v, parsedURL); !v.Valid() {
		json.ServerErrorWithDetails(w, r, v.Errors)
		return
	}

	err = h.store.Links.UpdateClicksByAlias(alias)
	if err != nil {
		json.ServerError(w, r, err)
	}

	http.Redirect(w, r, rawURL, http.StatusFound)
}

func (h *handler) linkAddHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Url   string `json:"url"`
		Alias string `json:"alias"`
	}

	err := json.ReadBody(w, r, &input)
	if err != nil {
		json.BadRequestResponse(w, r, err)
		return
	}

	link := model.Link{
		Url:   input.Url,
		Alias: input.Alias,
	}

	v := validator.New()

	parsedURL, err := url.ParseRequestURI(input.Url)
	if err != nil {
		json.BadRequestResponse(w, r, errors.New("missed url key or invalid URL provided"))
		return
	}

	validator.ValidateURL(v, parsedURL)
	validator.ValidateAlias(v, input.Alias)

	if !v.Valid() {
		json.FailedValidationResponse(w, r, v.Errors)
		return
	}

	err = h.store.Links.InsertLink(&link)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrDuplicateAlias):
			v.AddError("alias", fmt.Sprintf("'%s' is already exists in database", link.Alias))
			json.FailedValidationResponse(w, r, v.Errors)
		default:
			json.ServerError(w, r, err)
		}
		return
	}

	json.Created(w, r, link)
}
