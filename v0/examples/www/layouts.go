package main

import (
	d "github.com/LukasLow/dreego/v0/pkg/dom"

	"github.com/LukasLow/dreego/v0/pkg/dreego"
	"github.com/LukasLow/dreego/v0/pkg/scope"
)

// Marketing — die Hülle der öffentlichen Seiten.
// Der Kopf ist Inhalt (meta/title), NICHT das <head>-Element: scope.Document
// baut <head> selbst und setzt die gesammelten Styles hinein.
func Marketing(c *dreego.Ctx, head d.View, body d.View) d.View {
	kopf := d.Group([]d.View{
		d.Meta(d.Charset("utf-8")),
		d.Meta(d.Name("viewport"), d.Content("width=device-width, initial-scale=1")),
		head,
	})

	return c.Document("de", kopf,
		d.Group([]d.View{
			header(c),
			d.Main(d.Attr("id", "main"), body),
			footer(),
		}),
	)
}

// Auth — die schlanke Hülle der Anmeldeseiten.
func Auth(c *dreego.Ctx, head d.View, body d.View) d.View {
	kopf := d.Group([]d.View{
		d.Meta(d.Charset("utf-8")),
		d.Meta(d.Name("viewport"), d.Content("width=device-width, initial-scale=1")),
		head,
	})

	return c.Document("de", kopf, body)
}

func header(c *dreego.Ctx) d.View {
	aktiv := c.Nav()

	link := func(ziel string, beschriftung string, marke string) d.View {
		klasse := "navlink"
		if aktiv == marke {
			klasse += " on"
		}
		return d.A(d.Href(ziel), d.Class(klasse), d.Text(beschriftung))
	}

	return c.Box(
		scope.CSS(`
.mhead { display: flex; align-items: center; justify-content: space-between;
         gap: 16px; padding: 14px 26px; background: #fff;
         border-bottom: 2px solid rgba(0,0,0,.12) }
.mhead .brand { font-weight: 800; font-size: 18px; text-decoration: none; color: #111 }
.mhead .nav { display: flex; gap: 18px }
.mhead .navlink { color: #000; text-decoration: none; font-weight: 700 }
.mhead .navlink.on { color: #003870; text-decoration: underline }
`),
		d.Header(d.Class("mhead"),
			d.A(d.Href("/"), d.Class("brand"), d.Text("dreego")),
			d.Nav(d.Class("nav"),
				link("/", "Start", "home"),
				link("/login", "Anmelden", "login"),
			),
		),
	)
}

func footer() d.View {
	return d.Footer(d.Class("mfoot"), d.P(d.Text("dreego — Go-Webframework, ohne Compiler.")))
}
