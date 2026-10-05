package dreego

import (
	"io"

	g "github.com/LukasLow/dreego/v0/pkg/dom"

	"github.com/LukasLow/dreego/v0/pkg/scope"
)

// redirect is a Node that represents an HTTP redirect. The framework detects it
// and writes the status and Location header instead of rendering a body.
type redirect struct {
	url  string
	code int
}

// Render implements dom.View. A redirect has no body.
func (red redirect) Render(w io.Writer) error { return nil }

// Redirect returns a Node that redirects the client. Typical use is the
// Post/Redirect/Get pattern:
//
//	return dreego.Redirect("/auth/login", 303)
//
// The code must be a 3xx status; anything else is clamped to 303, so a bad
// value can never panic the server.
func Redirect(url string, code int) g.View {
	if code < 300 || code > 399 {
		code = 303
	}
	return redirect{url: url, code: code}
}

// Document renders a full page with the collector's deduplicated assets. It is
// the usual way for a layout to build its output.
func (c *Ctx) Document(lang string, head g.View, body g.View) g.View {
	return scope.Document(c.collector, lang, head, body)
}
