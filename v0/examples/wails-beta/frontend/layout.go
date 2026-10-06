package frontend

import (
	d "github.com/LukasLow/dreego/v0/pkg/dom"

	"github.com/LukasLow/dreego/v0/pkg/dreego"
)

// Shell is the page layout: it wires the vanilla template's <head> (charset,
// viewport, icon, stylesheet) around the page body. The stylesheet and images
// are served by dreego from frontend/static at /static/….
func Shell(c *dreego.Ctx, head d.View, body d.View) d.View {
	kopf := d.Group([]d.View{
		d.Meta(d.Charset("utf-8")),
		d.Meta(d.Name("viewport"), d.Content("width=device-width, initial-scale=1.0, viewport-fit=cover")),
		d.Link(d.Rel("icon"), d.Type("image/svg+xml"), d.Href("/static/wails.png")),
		d.Link(d.Rel("stylesheet"), d.Href("/static/style.css")),
		head,
	})

	return c.Document("en", kopf, body)
}
