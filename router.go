package sky

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"syscall"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/fvbock/endless"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
)

// Handler handler
type Handler func(*Context) error

// NewRouter new router
func NewRouter(rt *mux.Router, theme string, options ...render.Options) *Router {
	options = append(
		options,
		render.Options{
			Directory:  path.Join("themes", theme, "views"),
			Layout:     LayoutApplication,
			Extensions: []string{".html"},
		},
	)
	rt.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir(path.Join("themes", theme, "assets")))))
	return &Router{
		root: rt,
		render: render.New(
			options...,
		),
		handlers: make([]Handler, 0),
	}
}

// Router router
type Router struct {
	root     *mux.Router
	render   *render.Render
	handlers []Handler
}

// Start start
func (p *Router) Start(port int) error {
	addr := fmt.Sprintf(":%d", port)
	if IsProduction() {
		srv := endless.NewServer(addr, p.root)
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
	return http.ListenAndServe(addr, p.root)
}

// Use use
func (p *Router) Use(handlers ...Handler) {
	p.handlers = append(p.handlers, handlers...)
}

// LoadHTML load html
func (p *Router) LoadHTML(string) {
	// TODO
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
		render:   p.render,
		handlers: items,
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
			render:   p.render,
			handlers: items,
		}

		if err := ctx.Next(); err != nil {
			ctx.render.Text(ctx.Writer, http.StatusInternalServerError, err.Error())
			log.Error(err)
		}

		log.Infof("spend time %s", time.Now().Sub(begin))
	}).Name(name).Methods(method)
}
