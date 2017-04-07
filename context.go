package sky

import (
	"context"
	"math"
	"net"
	"net/http"
	"reflect"
	"runtime"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/go-playground/form"
	"github.com/gorilla/csrf"
	"github.com/unrolled/render"
	validator "gopkg.in/go-playground/validator.v9"
)

const (
	abortIndex = math.MaxUint8
	// DATA data key
	DATA = "data"
)

// K key type
type K string

// H hash
type H map[string]interface{}

// Context context
type Context struct {
	Request  *http.Request
	Writer   http.ResponseWriter
	params   map[string]string
	render   *render.Render
	decoder  *form.Decoder
	validate *validator.Validate
	handlers []Handler
	index    uint8
}

// Query get query
func (p *Context) Query(k string) string {
	return p.Request.URL.Query().Get(k)
}

// Param get url param
func (p *Context) Param(k string) string {
	return p.params[k]
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

// Bind bind to form
func (p *Context) Bind(fm interface{}) error {
	if err := p.Request.ParseForm(); err != nil {
		return err
	}

	if err := p.decoder.Decode(fm, p.Request.Form); err != nil {
		return err
	}
	return p.validate.Struct(fm)
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

// ClientIP implements a best effort algorithm to return the real client IP, it parses
// X-Real-IP and X-Forwarded-For in order to work properly with reverse-proxies such us: nginx or haproxy.
func (p *Context) ClientIP() string {
	// if p.engine.ForwardedByClientIP {
	// 	clientIP := strings.TrimSpace(c.requestHeader("X-Real-Ip"))
	// 	if len(clientIP) > 0 {
	// 		return clientIP
	// 	}
	// 	clientIP = c.requestHeader("X-Forwarded-For")
	// 	if index := strings.IndexByte(clientIP, ','); index >= 0 {
	// 		clientIP = clientIP[0:index]
	// 	}
	// 	clientIP = strings.TrimSpace(clientIP)
	// 	if len(clientIP) > 0 {
	// 		return clientIP
	// 	}
	// }
	if ip, _, err := net.SplitHostPort(strings.TrimSpace(p.Request.RemoteAddr)); err == nil {
		return ip
	}
	return ""
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

// JSON render json
func (p *Context) JSON(status int, value interface{}) {
	p.Writer.Header().Set("X-CSRF-Token", csrf.Token(p.Request))
	p.render.JSON(p.Writer, status, value)
}

// Redirect redirect to
func (p *Context) Redirect(status int, url string) {
	http.Redirect(p.Writer, p.Request, url, status)
}
