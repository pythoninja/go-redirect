package json

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
)

func NotFound(w http.ResponseWriter, r *http.Request) {
	message := "resource not found"
	errorResponse(w, r, http.StatusNotFound, message)
}

func MethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("method %s is not allowed for this resource", r.Method)
	errorResponse(w, r, http.StatusMethodNotAllowed, message)
}

func ServerError(w http.ResponseWriter, r *http.Request, err error) {
	logError(r, http.StatusInternalServerError, err)

	message := "the server encountered a problem and could not process your request"
	errorResponse(w, r, http.StatusInternalServerError, message)
}

func ServerErrorWithDetails(w http.ResponseWriter, r *http.Request, message map[string]string) {
	var validationErrors string
	for key, value := range message {
		validationErrors += fmt.Sprintf("%s: %s; ", key, value)
	}

	logError(r, http.StatusInternalServerError, errors.New(strings.TrimSuffix(validationErrors, "; ")))
	errorResponse(w, r, http.StatusInternalServerError, message)
}

func logError(r *http.Request, status int, err error) {
	logEntry := slog.With(
		slog.String("ip", r.RemoteAddr),
		slog.String("method", r.Method),
		slog.String("uri", r.RequestURI),
		slog.String("user_agent", r.UserAgent()),
		slog.Int("http_code", status))

	if err != nil {
		logEntry.Error("error logged",
			slog.String("details", err.Error()))

	} else {
		logEntry.Warn("warning logged")
	}
}
