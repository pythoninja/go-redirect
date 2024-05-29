package api

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

// readIdParam extracts and parses the "id" parameter from the request URL.
// It returns the parsed "id" as int64 if it is a positive integer, otherwise
// an error with message "invalid id parameter".
func readIdParam(r *http.Request) (int64, error) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}

	return id, nil
}

// readAliasParam extracts and returns the value of the "alias" parameter from the request URL.
func readAliasParam(r *http.Request) string {
	return chi.URLParam(r, "alias")
}
