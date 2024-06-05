package json

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/pythoninja/go-redirect/internal/server/response"
	"io"
	"log/slog"
	"net/http"
	"strings"
)

var contentTypeHeader = "application/json"

type responseWrapper map[string]any

// Ok handles generating a successful response for an HTTP request.
// It marshals the given body to JSON format using the toJSON function
// and returns the result to the client with a 200 OK status code.
// If there is an error during the marshaling process, it calls the ServerError function.
func Ok(w http.ResponseWriter, r *http.Request, body any) {
	bodyJSON, err := toJSON(body)
	if err != nil {
		ServerError(w, r, err)
		return
	}

	resp := response.New(w, r)
	resp.WithStatus(http.StatusOK)
	resp.WithBody(bodyJSON)
	resp.WithHeader("Content-Type", contentTypeHeader)
	resp.Write()
}

func Created(w http.ResponseWriter, r *http.Request, body any) {
	bodyJSON, err := toJSON(body)
	if err != nil {
		ServerError(w, r, err)
		return
	}

	resp := response.New(w, r)
	resp.WithStatus(http.StatusCreated)
	resp.WithBody(bodyJSON)
	resp.WithHeader("Content-Type", contentTypeHeader)
	resp.Write()
}

// errorResponse handles generating an error response for an HTTP request.
// It creates a responseWrapper map with the "error" key set to the given message.
// It then marshals the responseWrapper to JSON format using the toJSON function
// and returns the result to the client.
// If there is an error during the marshaling process, it logs the error.
func errorResponse(w http.ResponseWriter, r *http.Request, status int, message any, headers http.Header) {
	wrapper := responseWrapper{"errors": message}
	bodyJSON, err := toJSON(wrapper)
	if err != nil {
		slog.Error(err.Error())
	}

	resp := response.New(w, r)

	for key, values := range headers {
		for _, h := range values {
			resp.WithHeader(key, h)
		}
	}

	resp.WithStatus(status)
	resp.WithBody(bodyJSON)
	resp.WithHeader("Content-Type", contentTypeHeader)
	resp.Write()
}

// toJSON marshals the given body to JSON format with indentation and returns it as a byte slice.
// If the body is of type error, it wraps the error message in a responseWrapper map.
// It returns an error if there was a problem marshaling the JSON data.
func toJSON(body any) ([]byte, error) {
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

func ReadBody(w http.ResponseWriter, r *http.Request, dst any) error {
	maxBytes := int64(1_048_576)
	r.Body = http.MaxBytesReader(w, r.Body, maxBytes)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError
		var maxBytesError *http.MaxBytesError

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)

		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")

		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)

		case errors.Is(err, io.EOF):
			return fmt.Errorf("body must not be empty")

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body contains unknown key: %s", fieldName)

		case errors.As(err, &maxBytesError):
			return fmt.Errorf("body must not be larger that %d bytes", maxBytesError.Limit)

		case errors.As(err, &invalidUnmarshalError):
			panic(err)

		default:
			return err
		}
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only contain a single JSON value")
	}

	return nil
}
