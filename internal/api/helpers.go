package api

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
)

// readIDParam extracts and parses the "id" parameter from the request URL.
// It returns the parsed "id" as int64 if it is a positive integer, otherwise
// an error with message "invalid id parameter".
func readIDParam(r *http.Request) (int64, error) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}

	return id, nil
}

// readAliasParam extracts and returns the value of the "alias" parameter from the request URL.
func readAliasParam(r *http.Request) string {
	return r.PathValue("alias")
}

func normalizeAlias(alias string) string {
	return strings.ToLower(alias)
}
