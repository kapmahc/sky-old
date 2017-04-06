package sky

import (
	"context"
	"math"
	"net/http"
	"reflect"
	"runtime"

	log "github.com/Sirupsen/logrus"
	"github.com/unrolled/render"
)

const (
	abortIndex = math.MaxUint8
)

// K key type
type K string

// H hash
type H map[string]interface{}

// Context context
type Context struct {
	Request  *http.Request
	Writer   http.ResponseWriter
	render   *render.Render
	handlers []Handler
	index    uint8
}

// Next should be used only inside middleware.
func (p *Context) Next() error {

	idx := p.index
	p.index++

	if idx < uint8(len(p.handlers)) {
		hnd := p.handlers[idx]
		log.Debugf(runtime.FuncForPC(reflect.ValueOf(hnd).Pointer()).Name())
		if err := hnd((p)); err != nil {
			return err
		}
	}
	return nil
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
func (p *Context) Abort(status int, err error) {
	log.Warn(err)
	p.index = abortIndex
	p.HTML(
		status,
		"errors",
		H{"message": err.Error()},
	)
}

// HTML render html
func (p *Context) HTML(status int, name string, binding interface{}, htmlOpt ...render.HTMLOptions) {
	p.render.HTML(p.Writer, status, name, binding, htmlOpt...)
}
