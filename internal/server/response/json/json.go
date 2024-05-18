package json

import (
	"encoding/json"
	"errors"
	"github.com/pythoninja/go-redirect/internal/server/response"
	"log/slog"
	"net/http"
	"strings"
)

const contentTypeHeader = "application/json"

func OK(w http.ResponseWriter, r *http.Request, body any) {
	resp := response.New(w, r)
	resp.WithStatus(http.StatusOK)
	resp.WithBody(toJson(body))
	resp.WithHeader("Content-Type", contentTypeHeader)
	resp.Write()
}

func NotFound(w http.ResponseWriter, r *http.Request) {
	logWarning(r, http.StatusNotFound)

	err := errors.New("not found")
	errorResponse(w, r, http.StatusNotFound, err)
}

func MethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	logWarning(r, http.StatusMethodNotAllowed)

	err := errors.New("method not allowed")
	errorResponse(w, r, http.StatusMethodNotAllowed, err)
}

func errorResponse(w http.ResponseWriter, r *http.Request, status int, err error) {
	resp := response.New(w, r)
	resp.WithStatus(status)
	resp.WithBody(toJsonError(err))
	resp.WithHeader("Content-Type", contentTypeHeader)
	resp.Write()
}

func toJson(body any) []byte {
	js, err := json.MarshalIndent(body, "", strings.Repeat(" ", 2))
	if err != nil {
		slog.Error("Unable to marshal JSON data", slog.Any("error", err))
		return []byte("")
	}

	js = append(js, '\n')
	return js
}

func toJsonError(err error) []byte {
	type jsonWrapper struct {
		Message string `json:"error"`
	}

	message := jsonWrapper{Message: err.Error()}
	return toJson(message)
}

func logWarning(r *http.Request, status int) {
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
