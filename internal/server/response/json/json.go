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

// Ok handles generating a successful response for an HTTP request.
// It marshals the given body to JSON format using the toJson function
// and returns the result to the client with a 200 OK status code.
// If there is an error during the marshaling process, it calls the ServerError function.
func Ok(w http.ResponseWriter, r *http.Request, body any) {
	bodyJson, err := toJson(body)
	if err != nil {
		ServerError(w, r, err)
		return
	}

	resp := response.New(w, r)
	resp.WithStatus(http.StatusOK)
	resp.WithBody(bodyJson)
	resp.WithHeader("Content-Type", contentTypeHeader)
	resp.Write()
}

// errorResponse handles generating an error response for an HTTP request.
// It creates a responseWrapper map with the "error" key set to the given message.
// It then marshals the responseWrapper to JSON format using the toJson function
// and returns the result to the client.
// If there is an error during the marshaling process, it logs the error.
func errorResponse(w http.ResponseWriter, r *http.Request, status int, message any, headers http.Header) {
	wrapper := responseWrapper{"errors": message}
	bodyJson, err := toJson(wrapper)
	if err != nil {
		slog.Error(err.Error())
	}

	resp := response.New(w, r)

	if headers != nil {
		for key, values := range headers {
			for _, h := range values {
				resp.WithHeader(key, h)
			}
		}
	}

	resp.WithStatus(status)
	resp.WithBody(bodyJson)
	resp.WithHeader("Content-Type", contentTypeHeader)
	resp.Write()
}

// toJson marshals the given body to JSON format with indentation and returns it as a byte slice.
// If the body is of type error, it wraps the error message in a responseWrapper map.
// It returns an error if there was a problem marshaling the JSON data.
func toJson(body any) ([]byte, error) {
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
		return nil, fmt.Errorf("unable to marshal json: %s", err)
	}

	js = append(js, '\n')
	return js, nil
}
