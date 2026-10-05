package dreego

import (
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"

	"github.com/LukasLow/dreego/v0/pkg/scope"
)

// Layout is the page shell: it wraps a page's head and body and produces the
// full document. It is the Go replacement for a Dreego `layouts/*.dreego` file.
//
// The layout renders the body FIRST, so components can register their CSS and
// JS with the collector; then it puts the collected styles in <head> and the
// scripts at the end of <body>.
type Layout func(c *Ctx, head g.Node, body g.Node) g.Node

// defaultLayout is used when a page declares none. It is intentionally bare.
// The head and body arguments are their *content* — scope.Document writes the
// <head> and <body> elements itself.
func defaultLayout(c *Ctx, head g.Node, body g.Node) g.Node {
	kopf := g.Group([]g.Node{
		h.Meta(h.Charset("utf-8")),
		h.Meta(h.Name("viewport"), h.Content("width=device-width, initial-scale=1")),
		head,
	})

	return scope.Document(c.collector, "de", kopf, body)
}
