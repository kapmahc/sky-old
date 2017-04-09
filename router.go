package sky

import (
	"net/http"
	"reflect"
	"runtime"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/go-playground/form"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
	validator "gopkg.in/go-playground/validator.v9"
)

// Handler http handler func
type Handler func(*Context) error

// New Returns a new blank Router instance without any middleware attached.
func New() *Router {
	return &Router{
		router:   mux.NewRouter(),
		validate: validator.New(),
		decoder:  form.NewDecoder(),
		render:   render.New(),
		handlers: make([]Handler, 0),
	}
}

// Router It contains the muxer, middleware and configuration settings.
type Router struct {
	router   *mux.Router
	validate *validator.Validate
	decoder  *form.Decoder
	render   *render.Render
	handlers []Handler
}

// Run Attaches the router to a http.Server and starts listening and serving HTTP requests.
func (p *Router) Run(fn func(http.Handler) error) error {
	return fn(p.router)
}

// Use Adds a Handler onto the middleware stack. Handlers are invoked in the order they are added to a Negroni.
func (p *Router) Use(handlers ...Handler) {
	p.handlers = append(p.handlers, handlers...)
}

// Group Creates a new router group. You should add all the routes that have common middlwares or the same path prefix.
func (p *Router) Group(path string, fn func(*Router)) {
	handlers := make([]Handler, len(p.handlers))
	copy(handlers, p.handlers)

	fn(&Router{
		router:   p.router.PathPrefix(path).Subrouter(),
		validate: p.validate,
		render:   p.render,
		decoder:  p.decoder,
		handlers: handlers,
	})
}

// GET http get
func (p *Router) GET(name, path string, handlers ...Handler) {
	p.handle(http.MethodGet, name, path, handlers...)
}

// POST http post
func (p *Router) POST(name, path string, handlers ...Handler) {
	p.handle(http.MethodPost, name, path, handlers...)
}

// PUT http put
func (p *Router) PUT(name, path string, handlers ...Handler) {
	p.handle(http.MethodPut, name, path, handlers...)
}

// PATCH http patch
func (p *Router) PATCH(name, path string, handlers ...Handler) {
	p.handle(http.MethodPatch, name, path, handlers...)
}

// DELETE http delete
func (p *Router) DELETE(name, path string, handlers ...Handler) {
	p.handle(http.MethodDelete, name, path, handlers...)
}

func (p *Router) handle(method, name, path string, handlers ...Handler) {
	items := make([]Handler, len(p.handlers))
	copy(items, p.handlers)
	items = append(items, handlers...)

	p.router.HandleFunc(path, func(wrt http.ResponseWriter, req *http.Request) {
		begin := time.Now()
		log.Infof(
			"%s %s %s%s",
			req.Proto,
			req.Method,
			req.Host,
			req.RequestURI,
		)

		ctx := Context{
			Request:  req,
			Writer:   wrt,
			validate: p.validate,
			decoder:  p.decoder,
			render:   p.render,
		}

		for _, h := range items {
			log.Debug(runtime.FuncForPC(reflect.ValueOf(h).Pointer()).Name())
			if err := h(&ctx); err != nil {
				ctx.Abort(http.StatusInternalServerError, err)
				break
			}
		}
		log.Infof("spend time %s", time.Now().Sub(begin))
	}).Name(name).Methods(method)
}
