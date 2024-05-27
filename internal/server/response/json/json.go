package json

import (
	"encoding/json"
	"fmt"
	"github.com/pythoninja/go-redirect/internal/server/response"
	"log/slog"
	"net/http"
	"strings"
)

var contentTypeHeader = "application/json"

type responseWrapper map[string]any

func Ok(w http.ResponseWriter, r *http.Request, body any) {
	resp := response.New(w, r)
	resp.WithStatus(http.StatusOK)
	resp.WithBody(toJson(body))
	resp.WithHeader("Content-Type", contentTypeHeader)
	resp.Write()
}

func NotFound(w http.ResponseWriter, r *http.Request) {
	logWarning(r, http.StatusNotFound)

	message := "resource not found"
	errorResponse(w, r, http.StatusNotFound, message)
}

func MethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	logWarning(r, http.StatusMethodNotAllowed)

	message := fmt.Sprintf("method %s is not allowed for this resource", r.Method)
	errorResponse(w, r, http.StatusMethodNotAllowed, message)
}

func errorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	resp := response.New(w, r)
	resp.WithStatus(status)
	resp.WithBody(toJsonError(message))
	resp.WithHeader("Content-Type", contentTypeHeader)
	resp.Write()
}

func toJson(body any) []byte {
	var message any

	switch m := body.(type) {
	case error:
		message = responseWrapper{"error": m.Error()}
	default:
		message = m
	}

	js, err := json.MarshalIndent(message, "", strings.Repeat(" ", 2))
	if err != nil {
		slog.Error("Unable to marshal JSON data", slog.Any("error", err))
		return []byte("") // todo: return error here, otherwise it will be replaced with ""
	}

	js = append(js, '\n')
	return js
}

func toJsonError(message any) []byte {
	return toJson(responseWrapper{"error": message})
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
