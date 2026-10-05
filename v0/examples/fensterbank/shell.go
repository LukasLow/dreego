package main

import (
	d "github.com/LukasLow/dreego/v0/pkg/dom"

	"github.com/LukasLow/dreego/v0/pkg/dreego"
	"github.com/LukasLow/dreego/v0/pkg/urls"
)

// markt ist die Adresse der zweiten Site (Portal, eigener Port). Cross-Site-Links
// laufen über den urls-Helfer, nicht relativ.
var markt = urls.FromEnv("PORTAL_URL", "http://localhost:4001")

// basisCSS wird zur Laufzeit aus dem Embed gelesen und inline gesetzt — kein
// blockierender Stylesheet-Link, kein weisser Blitz.
var basisCSS = mustRead("public/basis.css")

func mustRead(pfad string) string {
	daten, err_lesen := publicFS.ReadFile(pfad)
	if err_lesen != nil {
		panic("fensterbank: " + pfad + " fehlt im Embed: " + err_lesen.Error())
	}
	return string(daten)
}

// Shell ist die Hülle aller Seiten.
func Shell(c *dreego.Ctx, head d.Node, body d.Node) d.Node {
	c.AddCritical(basisCSS)

	kopf := d.Group([]d.Node{
		d.Meta(d.Charset("utf-8")),
		d.Meta(d.Name("viewport"), d.Content("width=device-width, initial-scale=1")),
		d.Meta(d.Name("theme-color"), d.Content("#f2f7f4")),
		head,
	})

	return c.Document("de", kopf, d.Group([]d.Node{
		nav(c),
		d.Main(body),
		fuss(c),
	}))
}

func nav(c *dreego.Ctx) d.Node {
	aktiv := c.Nav()
	link := func(ziel, text, marke string) d.Node {
		eigenschaften := []d.Node{d.Href(ziel), d.Text(text)}
		if aktiv == marke {
			eigenschaften = append(eigenschaften, d.Class("on"))
		}
		return d.A(eigenschaften...)
	}
	return d.Header(d.Class("kopf"),
		d.A(d.Href("/"), d.Class("marke"), d.Text("Fenster"), d.Em(d.Text("bank"))),
		d.Nav(
			link("/", c.T("nav.start"), "start"),
			link("/sorten", c.T("nav.sorten"), "sorten"),
			link("/pflege", c.T("nav.pflege"), "pflege"),
			link("/kontakt", c.T("nav.kontakt"), "kontakt"),
			link("/shop", c.T("nav.shop"), "shop"),
			d.A(d.Href(markt.Abs("/login")), d.Text(c.T("nav.anmelden"))),
		),
	)
}

func fuss(c *dreego.Ctx) d.Node {
	return d.Footer(d.Class("fuss"), d.P(d.Text(c.T("fuss"))))
}
