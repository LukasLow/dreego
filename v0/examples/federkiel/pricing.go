package main

import (
	d "github.com/LukasLow/dreego/v0/pkg/dom"

	"github.com/LukasLow/dreego/v0/pkg/dreego"
	"github.com/LukasLow/dreego/v0/pkg/scope"
)

// Pricing — Preisseite.
var Pricing = dreego.Page{
	Path:   "/pricing",
	Layout: Shell,
	Nav:    "pricing",
	Head: func(c *dreego.Ctx) []d.Node {
		return []d.Node{d.TitleEl(d.Text("Preise — Federkiel"))}
	},
	Get: getPricing,
}

func getPricing(c *dreego.Ctx) d.Node {
	return c.Box(
		scope.CSS(`
.preise { padding: 34px 0; text-align: center }
.preise h1 { font-size: 36px; margin: 0 0 8px }
.preise .lead { color: #5c5347; margin: 0 0 26px }
.stufen { display: grid; grid-template-columns: 1fr 1fr 1fr; gap: 16px; text-align: left }
.stufe { border: 2px solid #e2d9c8; border-radius: 16px; padding: 22px 20px; display: flex; flex-direction: column }
.stufe.tipp { border-color: #b4451f }
.stufe .preis { font-size: 34px; font-weight: 700; margin: 6px 0 2px }
.stufe ul { padding-left: 20px; margin: 12px 0 18px; color: #5c5347 }
.stufe .btn { margin-top: auto; text-align: center }
@media (max-width: 640px) { .preise h1 { font-size: 28px } .stufen { grid-template-columns: 1fr } }
`),
		d.Section(d.Class("preise"),
			d.H1(d.Text("Preise")),
			d.P(d.Class("lead"), d.Text("Einmal zahlen, für immer schreiben. Keine Abo-Falle.")),
			d.Div(d.Class("stufen"),
				stufe("Skizze", "0 €", []string{"1 Manuskript", "Alle Funktionen", "Kein Konto nötig"}, false),
				stufe("Autorin", "29 €", []string{"Unbegrenzt Manuskripte", "Figuren-Kartei", "Export DOCX/EPUB"}, true),
				stufe("Verlag", "99 €", []string{"Bis zu 10 Autorinnen", "Geteilte Kapitel", "Priorisierter Support"}, false),
			),
		),
	)
}

func stufe(name string, preis string, punkte []string, tipp bool) d.Node {
	klasse := "stufe"
	if tipp {
		klasse += " tipp"
	}
	return d.Div(d.Class(klasse),
		d.H3(d.Text(name)),
		d.P(d.Class("preis"), d.Text(preis)),
		d.Ul(d.Map(punkte, func(punkt string) d.Node { return d.Li(d.Text(punkt)) })),
		d.A(d.Href("/login"), d.Class("btn"), d.Text("Wählen")),
	)
}
