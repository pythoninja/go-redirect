package json

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
)

func NotFound(w http.ResponseWriter, r *http.Request) {
	logMessage(r, http.StatusNotFound, nil)

	message := "resource not found"
	errorResponse(w, r, http.StatusNotFound, message)
}

func MethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	logMessage(r, http.StatusMethodNotAllowed, nil)

	message := fmt.Sprintf("method %s is not allowed for this resource", r.Method)
	errorResponse(w, r, http.StatusMethodNotAllowed, message)
}

func ServerError(w http.ResponseWriter, r *http.Request, err error) {
	logMessage(r, http.StatusInternalServerError, err)

	message := "the server encountered a problem and could not process your request"
	errorResponse(w, r, http.StatusInternalServerError, message)
}

func ServerErrorWithDetails(w http.ResponseWriter, r *http.Request, message map[string]string) {
	var validationErrors string
	for key, value := range message {
		validationErrors += fmt.Sprintf("%s: %s; ", key, value)
	}

	logMessage(r, http.StatusInternalServerError, errors.New(strings.TrimSuffix(validationErrors, "; ")))
	errorResponse(w, r, http.StatusInternalServerError, message)
}

func logMessage(r *http.Request, status int, err error) {
	if err != nil {
		slog.Error(http.StatusText(status),
			slog.String("ip", r.RemoteAddr),
			slog.Group("request",
				slog.String("method", r.Method),
				slog.String("uri", r.RequestURI),
				slog.String("user_agent", r.UserAgent()),
			),
			slog.Group("response", slog.Int("http_code", status)),
			slog.String("details", err.Error()),
		)
	} else {
		slog.Warn(http.StatusText(status),
			slog.String("ip", r.RemoteAddr),
			slog.Group("request",
				slog.String("method", r.Method),
				slog.String("uri", r.RequestURI),
				slog.String("user_agent", r.UserAgent()),
			),
			slog.Group("response", slog.Int("http_code", status)),
		)
	}
}
