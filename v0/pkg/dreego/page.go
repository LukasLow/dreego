package dreego

import (
	"net/http"

	g "github.com/LukasLow/dreego/v0/pkg/dom"
)

// HeadFunc builds the page's <head> contributions (title, meta, links). They
// are merged into the layout head and deduplicated there.
type HeadFunc func(c *Ctx) []g.View

// HandlerFunc renders the page body for one HTTP method.
type HandlerFunc func(c *Ctx) g.View

// APIHandler responds with a non-HTML body (JSON, XML, plain text, …). It is
// the dreego replacement for Dreego's `<server type="json">` routes.
type APIHandler func(c *Ctx) error

// Page is one route. Declared at the top of its file, with the handler
// functions directly below it:
//
//	var Pricing = dreego.Page{
//	    Path:   "/pricing",
//	    Layout: layouts.Marketing,
//	    Nav:    "pricing",
//	    Head:   headPricing,
//	    Get:    getPricing,
//	}
//
// Path is the single source of the URL. Dynamic segments use brackets:
// "/projekte/[id]" — read them with c.Param("id").
//
// Get and Post select the HTTP method for an HTML page. API turns the page into
// a non-HTML endpoint instead (no layout, no render). When Middleware is set, it
// wraps only this page.
type Page struct {
	Path       string
	Layout     Layout
	Nav        string
	Head       HeadFunc
	Get        HandlerFunc
	Post       HandlerFunc
	API        APIHandler
	Methods    []string
	Middleware []func(http.Handler) http.Handler
}

// methods lists the HTTP methods this page answers. An explicit Methods list
// wins; otherwise it is derived from the handlers (API defaults to GET).
func (seite Page) methods() []string {
	if len(seite.Methods) > 0 {
		return seite.Methods
	}

	var methoden []string
	if seite.Get != nil {
		methoden = append(methoden, "GET")
	}
	if seite.Post != nil {
		methoden = append(methoden, "POST")
	}
	if len(methoden) == 0 && seite.API != nil {
		methoden = append(methoden, "GET")
	}
	return methoden
}

// handlerFor returns the handler for the given HTTP method, or nil.
func (seite Page) handlerFor(method string) HandlerFunc {
	switch method {
	case "GET", "HEAD":
		return seite.Get
	case "POST":
		return seite.Post
	default:
		return nil
	}
}

// allowedMethods lists the methods this page answers.
func (seite Page) allowedMethods() string {
	switch {
	case seite.Get != nil && seite.Post != nil:
		return "GET, HEAD, POST"
	case seite.Get != nil:
		return "GET, HEAD"
	case seite.Post != nil:
		return "POST"
	default:
		return ""
	}
}

// headNode collects the page head into one node.
func (seite Page) headNode(c *Ctx) g.View {
	if seite.Head == nil {
		return g.Group(nil)
	}
	return g.Group(seite.Head(c))
}
