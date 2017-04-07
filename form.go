package sky

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/csrf"
)

// NewForm new form
func NewForm(c *Context) Form {
	return Form{
		"id":             "fm-" + uuid.New().String(),
		"lang":           c.Get(LOCALE),
		"method":         http.MethodPost,
		"action":         c.Request.URL.RequestURI(),
		"fields":         make([]interface{}, 0),
		csrf.TemplateTag: csrf.TemplateField(c.Request),
	}
}

// Form form
type Form map[string]interface{}

// Title set title
func (p Form) Title(t string) {
	p["title"] = t
}

// Email email field
type Email map[string]interface{}
