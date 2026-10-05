package main

import (
	d "github.com/LukasLow/dreego/v0/pkg/dom"

	"github.com/LukasLow/dreego/v0/pkg/dreego"
)

// basisCSS wird inline ausgeliefert — kein weisser Blitz.
var basisCSS = mustRead("public/basis.css")

func mustRead(pfad string) string {
	daten, err_lesen := publicFS.ReadFile(pfad)
	if err_lesen != nil {
		panic("portal: " + pfad + " fehlt im Embed: " + err_lesen.Error())
	}
	return string(daten)
}

// Shell ist die Hülle aller Portal-Seiten.
func Shell(c *dreego.Ctx, head d.Node, body d.Node) d.Node {
	c.AddCritical(basisCSS)

	kopf := d.Group([]d.Node{
		d.Meta(d.Charset("utf-8")),
		d.Meta(d.Name("viewport"), d.Content("width=device-width, initial-scale=1")),
		d.Meta(d.Name("theme-color"), d.Content("#eef1f6")),
		head,
	})

	return c.Document("de", kopf, d.Group([]d.Node{
		nav(c),
		d.Main(body),
		fuss(),
	}))
}

func nav(c *dreego.Ctx) d.Node {
	// Cross-Site-Link über den urls-Helfer: absolute Adresse zur öffentlichen
	// Site (anderer Port), nicht relativ.
	return d.Header(d.Class("kopf"),
		d.A(d.Href(markt.Abs("/")), d.Class("marke"), d.Text("Fensterbank "), d.Em(d.Text("Portal"))),
		d.Nav(
			d.A(d.Href(public.Abs("/")), d.Text("Zur Website")),
			d.A(d.Href(markt.Abs("/login")), d.Text("Anmelden")),
			d.A(d.Href(markt.Abs("/konto")), d.Text("Konto")),
		),
	)
}

func fuss() d.Node {
	return d.Footer(d.Class("fuss"),
		d.P(d.Text("Portal — zweite Site der Demo (eigener Port), gebaut mit dreego.")))
}
