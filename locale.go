package sky

import (
	"net/http"
	"time"

	"golang.org/x/text/language"
)

const (
	// LOCALE locale key
	LOCALE = "locale"
)

// LocaleMiddleware new locale-middleware
func LocaleMiddleware(matcher language.Matcher) Handler {
	return func(c *Context) error {
		write := false
		// 1. Check URL arguments.
		lang := c.Request.URL.Query().Get(LOCALE)

		// 2. Get language information from cookies.
		if len(lang) == 0 {
			if ck, er := c.Request.Cookie(LOCALE); er == nil {
				lang = ck.Value
			}
		} else {
			write = true
		}
		// 3. Get language information from 'Accept-Language'.
		if len(lang) == 0 {
			al := c.Request.Header.Get("Accept-Language")
			if len(al) > 4 {
				lang = al[:5] // Only compare first 5 letters.
			}
		}

		tag, _, _ := matcher.Match(language.Make(lang))
		ts := tag.String()
		if ts != lang {
			write = true
		}

		if write {
			http.SetCookie(c.Writer, &http.Cookie{
				Name:    LOCALE,
				Value:   ts,
				Expires: time.Now().AddDate(10, 0, 0),
				Path:    "/",
			})
		}

		c.Set(LOCALE, ts)
		return c.Next()
	}
}
