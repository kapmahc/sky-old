package sky

import (
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
)

// URLFor url-for
type URLFor struct {
	Router *mux.Router `inject:""`
}

// URL builds a url for the route.
func (p *URLFor) URL(name string, pairs ...interface{}) string {

	rt := p.Router.Get(name)
	if rt == nil {
		return name
	}
	var params []string
	for _, v := range pairs {
		switch t := v.(type) {
		case string:
			params = append(params, v.(string))
		default:
			log.Warn("unknown type", t)
		}
	}
	url, err := rt.URL(params...)
	if err != nil {
		log.Error(err)
		return name
	}
	return url.String()
}
