package main

import (
	"fmt"
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/kapmahc/sky"
)

func handler(id int) sky.Handler {
	return func(c *sky.Context) error {
		fmt.Fprintf(c.Writer, "handler %d\n", id)
		return nil
	}
}

func main() {
	log.SetLevel(log.DebugLevel)
	rt := sky.New()
	rt.Use(
		handler(1),
		handler(2),
		handler(3),
	)

	rt.GET("hello", "/hello", handler(4), handler(5))
	rt.GET("hi", "/hi", handler(6), handler(7))

	rt.Group(
		"/api",
		func(r *sky.Router) {
			r.Use(
				handler(11),
				handler(12),
			)
			r.GET("api.hi", "/hi", handler(13), handler(14))
		},
	)

	rt.Run(func(h http.Handler) error {
		http.ListenAndServe(":3000", h)
		return nil
	})

	time.Sleep(20)
}
