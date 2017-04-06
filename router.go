package sky

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"syscall"

	log "github.com/Sirupsen/logrus"
	"github.com/fvbock/endless"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
)

// Handler handler
type Handler func(*Context) error

// NewRouter new router
func NewRouter(opt render.Options) *Router {
	return &Router{
		root:     mux.NewRouter(),
		render:   render.New(opt),
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
	}
	http.Handle("/", p.root)
	return nil
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
func (p *Router) Walk(func(name, method, path string)) {
	// TODO
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
		ctx := Context{Request: req, Writer: wrt}
		for _, hnd := range items {
			if err := hnd(&ctx); err != nil {
				ctx.Abort(http.StatusInternalServerError, err)
				return
			}
			if ctx.written {
				return
			}
		}
	}).Methods(method)
}
