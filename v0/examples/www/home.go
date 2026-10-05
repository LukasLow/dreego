package main

import (
	d "github.com/LukasLow/dreego/v0/pkg/dom"

	"github.com/LukasLow/dreego/v0/examples/www/components"
	"github.com/LukasLow/dreego/v0/pkg/dreego"
)

// Home — die Startseite. Route oben, Handler darunter.
var Home = dreego.Page{
	Path:   "/",
	Layout: Marketing,
	Nav:    "home",
	Head:   headHome,
	Get:    getHome,
}

func headHome(c *dreego.Ctx) []d.Node {
	return []d.Node{
		d.TitleEl(d.Text("dreego — Startseite")),
		d.Meta(d.Name("description"), d.Content("Eine Seite, gebaut mit dreego.")),
	}
}

func getHome(c *dreego.Ctx) d.Node {
	return d.Group([]d.Node{
		components.Hero(c),
		d.Section(d.Class("wrap"),
			d.H2(d.Text("Sicherheit")),
			d.P(d.Text("Session (verschlüsselt), CSRF, Flash — alles im Kern.")),
			d.P(d.A(d.Href("/login"), d.Text("zum Demo-Login"))),
		),
	})
}
