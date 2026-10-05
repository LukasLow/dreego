package main

import (
	d "github.com/LukasLow/dreego/v0/pkg/dom"

	"github.com/LukasLow/dreego/v0/pkg/dreego"
)

// basisCSS ist das komplette Grund-Stylesheet, zur Laufzeit aus dem Embed
// gelesen. Es wird inline in <head> gesetzt, damit vor dem ersten Paint
// KEINE render-blockierende CSS-Anfrage offen ist — das verhindert den
// weissen Blitz beim Refresh.
//
// Die @font-face-Verweise darin zeigen weiter auf /public/fonts/*.woff2; die
// Schriften laden asynchron (font-display: swap), der Text ist aber sofort in
// der Fallback-Schrift sichtbar — nie weiss.
var basisCSS = mustReadPublic("public/style.css")

func mustReadPublic(pfad string) string {
	daten, err_lesen := publicFS.ReadFile(pfad)
	if err_lesen != nil {
		// Kein stiller Default: fehlt diese Datei, ist das Styling kaputt.
		panic("federkiel: " + pfad + " nicht im Embed: " + err_lesen.Error())
	}
	return string(daten)
}

// Shell ist die Hülle aller Federkiel-Seiten.
func Shell(c *dreego.Ctx, head d.View, body d.View) d.View {
	// Das komplette Basis-CSS inline -> kein blockierender Request -> kein
	// weisser Blitz. (scope.Collector dedupliziert es pro Seite.)
	c.AddCritical(basisCSS)

	kopf := d.Group([]d.View{
		d.Meta(d.Charset("utf-8")),
		d.Meta(d.Name("viewport"), d.Content("width=device-width, initial-scale=1")),
		// theme-color faerbt zusaetzlich die Browser-/Systemleiste (mobil).
		d.Meta(d.Name("theme-color"), d.Content("#fffdf7")),
		d.Link(d.Rel("icon"), d.Type("image/svg+xml"), d.Href("/public/favicon.svg")),
		// Fonts vorladen: die Anfrage startet parallel zum HTML-Parsen,
		// statt erst nach dem CSS — einer der groessten Latenz-Hebel.
		d.Link(d.Rel("preload"), d.As("font"), d.Type("font/woff2"),
			d.Href("/public/fonts/atkinson-400.woff2"), d.Attr("crossorigin", "")),
		d.Link(d.Rel("preload"), d.As("font"), d.Type("font/woff2"),
			d.Href("/public/fonts/atkinson-700.woff2"), d.Attr("crossorigin", "")),
		d.Script(d.Src("/public/app.js"), d.Defer()),
		// Bewusst KEIN <link rel="stylesheet"> mehr: das Basis-CSS ist inline.
		head,
	})

	return c.Document("de", kopf, d.Group([]d.View{
		nav(c),
		d.Main(body),
		footer(),
	}))
}

func nav(c *dreego.Ctx) d.View {
	aktiv := c.Nav()

	link := func(ziel string, text string, marke string) d.View {
		eigenschaften := []d.View{d.Href(ziel), d.Text(text)}
		if aktiv == marke {
			eigenschaften = append(eigenschaften, d.Class("on"))
		}
		return d.A(eigenschaften...)
	}

	return d.Header(d.Class("sitenav"),
		d.A(d.Href("/"), d.Class("brand"), d.Text("Feder"), d.Em(d.Text("kiel"))),
		d.Nav(
			link("/features", "Funktionen", "features"),
			link("/pricing", "Preise", "pricing"),
			link("/impressum", "Impressum", "impressum"),
			d.A(d.Href("/login"), d.Class("btn ghost"), d.Text("Anmelden")),
		),
	)
}

func footer() d.View {
	return d.Footer(
		d.P(d.Text("Federkiel — eine erfundene Demo-App, gebaut mit dreego.")),
	)
}
