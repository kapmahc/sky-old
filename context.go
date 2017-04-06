package sky

import (
	"context"
	"net/http"
)

// K key type
type K string

// H hash
type H map[string]interface{}

// Context context
type Context struct {
	Request *http.Request
	Writer  http.ResponseWriter
	written bool
}

// Get get
func (p *Context) Get(key string) interface{} {
	return p.Request.Context().Value(K(key))
}

// Set set
func (p *Context) Set(key string, value interface{}) {
	ctx := context.WithValue(p.Request.Context(), K(key), value)
	p.Request = p.Request.WithContext(ctx)
}

// Abort abort
func (p *Context) Abort(code int, err error) {
	// TODO render templat by different code
	p.written = true
}
