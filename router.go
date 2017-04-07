package sky

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path"
	"syscall"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/fvbock/endless"
	"github.com/go-playground/form"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
	validator "gopkg.in/go-playground/validator.v9"
)

// Handler handler
type Handler func(*Context) error

// NewRouter new router
func NewRouter(rt *mux.Router, theme string, funcs template.FuncMap) *Router {

	rt.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir(path.Join("themes", theme, "assets")))))
	return &Router{
		root: rt,
		render: render.New(
			render.Options{
				Directory:  path.Join("themes", theme, "views"),
				Layout:     "application",
				Extensions: []string{".html"},
				Funcs:      []template.FuncMap{funcs},
			},
		),
		handlers: make([]Handler, 0),
		decoder:  form.NewDecoder(),
		validate: validator.New(),
	}
}

// Router router
type Router struct {
	root     *mux.Router
	render   *render.Render
	handlers []Handler
	decoder  *form.Decoder
	validate *validator.Validate
}

// Start start
func (p *Router) Start(port int, key []byte) error {
	addr := fmt.Sprintf(":%d", port)
	hnd := csrf.Protect(
		key,
		csrf.RequestHeader("Authenticity-Token"),
		csrf.FieldName("authenticity_token"),
		csrf.Secure(IsProduction()),
		csrf.CookieName("_csrf_"),
		csrf.Path("/"),
	)(p.root)
	if IsProduction() {
		srv := endless.NewServer(addr, hnd)
		srv.BeforeBegin = func(add string) {
			fd, err := os.OpenFile(path.Join("tmp", "pid"), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
			if err != nil {
				log.Error(err)
				return
			}
			defer fd.Close()
			pid := syscall.Getpid()
			log.Printf("pid is %d", pid)
			fmt.Fprint(fd, pid)
		}
		return srv.ListenAndServe()
	}
	return http.ListenAndServe(addr, hnd)
}

// Use use
func (p *Router) Use(handlers ...Handler) {
	p.handlers = append(p.handlers, handlers...)
}

// Walk walk routes
func (p *Router) Walk(fn func(name, method, path string) error) error {
	return p.root.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		tpl, err := route.GetPathTemplate()
		if err != nil {
			return err
		}
		// TODO methods https://github.com/gorilla/mux/pull/244
		fn(route.GetName(), "", tpl)
		return nil
	})
}

// Group group
func (p *Router) Group(path string, fn func(*Router), handlers ...Handler) {
	rt := p.root.PathPrefix(path).Subrouter()
	items := make([]Handler, len(p.handlers))
	copy(items, p.handlers)
	items = append(items, handlers...)

	fn(&Router{
		root:     rt,
		handlers: items,
		render:   p.render,
		decoder:  p.decoder,
		validate: p.validate,
	})
}

// Get http aget
func (p *Router) Get(name, path string, handlers ...Handler) {
	p.add(http.MethodGet, name, path, handlers...)
}

// Post http post
func (p *Router) Post(name, path string, handlers ...Handler) {
	p.add(http.MethodPost, name, path, handlers...)
}

// Delete http delete
func (p *Router) Delete(name, path string, handlers ...Handler) {
	p.add(http.MethodDelete, name, path, handlers...)
}

func (p *Router) add(method, name, path string, handlers ...Handler) {
	items := make([]Handler, len(p.handlers))
	copy(items, p.handlers)
	items = append(items, handlers...)

	p.root.HandleFunc(path, func(wrt http.ResponseWriter, req *http.Request) {
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
			params:   mux.Vars(req),
			handlers: items,
			render:   p.render,
			decoder:  p.decoder,
			validate: p.validate,
		}

		if err := ctx.Next(); err != nil {
			ctx.render.Text(ctx.Writer, http.StatusInternalServerError, err.Error())
			log.Error(err)
		}

		log.Infof("spend time %s", time.Now().Sub(begin))
	}).Name(name).Methods(method)
}
