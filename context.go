package sky

import (
	"net/http"

	"github.com/go-playground/form"
	"github.com/unrolled/render"
	validator "gopkg.in/go-playground/validator.v9"
)

// Context It allows us to pass variables between middleware, manage the flow, validate the JSON of a request and render a JSON response for example.
type Context struct {
	Request  *http.Request
	Writer   http.ResponseWriter
	validate *validator.Validate
	decoder  *form.Decoder
	render   *render.Render
}

// Abort Prevents pending handlers from being called. Note that this will not stop the current handler.
func (p *Context) Abort(code int, err error) {
	p.render.Text(p.Writer, code, err.Error())
}
