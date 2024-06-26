package response

import (
	"log/slog"
	"net/http"
)

type Builder struct {
	w          http.ResponseWriter
	r          *http.Request
	statusCode int
	body       []byte
	headers    http.Header
}

func New(w http.ResponseWriter, r *http.Request) *Builder {
	return &Builder{
		w:          w,
		r:          r,
		statusCode: http.StatusOK,     // set the default HTTP response code to 200
		headers:    make(http.Header), // create empty header map
	}
}

func (b *Builder) WithStatus(statusCode int) *Builder {
	b.statusCode = statusCode

	return b
}

func (b *Builder) WithHeader(key, value string) *Builder {
	b.headers.Add(key, value)

	return b
}

func (b *Builder) WithBody(body []byte) *Builder {
	b.body = body

	return b
}

func (b *Builder) Write() {
	b.writeHeaders()

	if b.body == nil {
		return
	}

	if _, err := b.w.Write(b.body); err != nil {
		slog.Error("error writing body", slog.Any("error", err))
	}
}

func (b *Builder) writeHeaders() {
	for key, value := range b.headers {
		b.w.Header()[key] = value
	}

	b.w.WriteHeader(b.statusCode)
}
