package response

import (
	"net/http"
)

type Builder struct {
	w          http.ResponseWriter
	r          *http.Request
	statusCode int
	body       any
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
	b.headers.Set(key, value)
	return b
}

func (b *Builder) WithBody(body any) *Builder {
	b.body = body
	return b
}

func (b *Builder) Write() {
	if b.body == nil {
		b.writeHeaders()
		return
	}

	b.writeHeaders()

	switch v := b.body.(type) {
	case []byte:
		b.w.Write(v)
	case string:
		b.w.Write([]byte(v))
	case error:
		b.w.Write([]byte(v.Error()))
	}
}

func (b *Builder) writeHeaders() {
	// Set additional headers if any
	for key, value := range b.headers {
		b.w.Header()[key] = value
	}

	b.w.WriteHeader(b.statusCode)
}