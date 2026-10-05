package main

import (
	d "github.com/LukasLow/dreego/v0/pkg/dom"

	"github.com/LukasLow/dreego/v0/pkg/dreego"
	"github.com/LukasLow/dreego/v0/pkg/scope"
)

// Features — Leistungsseite mit Kartenliste.
var Features = dreego.Page{
	Path:   "/features",
	Layout: Shell,
	Nav:    "features",
	Head: func(c *dreego.Ctx) []d.Node {
		return []d.Node{d.TitleEl(d.Text("Funktionen — Federkiel"))}
	},
	Get: getFeatures,
}

type funktion struct {
	titel string
	text  string
}

func getFeatures(c *dreego.Ctx) d.Node {
	posten := []funktion{
		{"Fokus-Modus", "Keine Benachrichtigungen, kein Zähler, nur die Seite."},
		{"Kapitel-Ordner", "Ordne Szenen, Figuren und Notizen nach Kapiteln."},
		{"Wortziel", "Ein ruhiges Tagesziel statt einer Deadline."},
		{"Figuren-Kartei", "Wer war noch gleich Maras Bruder?"},
		{"Lesemodus", "Der ganze Text in Serifenschrift, ohne Rahmen."},
		{"Export", "Als Markdown, DOCX oder EPUB herausgeben."},
	}

	return c.Box(
		scope.CSS(`
.features { padding: 34px 0 }
.features h1 { font-size: 36px; margin: 0 0 8px }
.features .lead { color: #5c5347; margin: 0 0 26px }
.karten { display: grid; grid-template-columns: 1fr 1fr; gap: 16px }
.karte { border: 2px solid #e2d9c8; border-radius: 14px; padding: 18px 20px; height: 100% }
.karte h3 { margin: 0 0 6px; font-size: 20px }
.karte p { margin: 0; color: #5c5347; }
@media (max-width: 640px) { .features h1 { font-size: 28px } .karten { grid-template-columns: 1fr } }
`),
		d.Section(d.Class("features"),
			d.H1(d.Text("Funktionen")),
			d.P(d.Class("lead"), d.Text("Sechs Dinge, die Federkiel bewusst einfach hält.")),
			d.Div(d.Class("karten"),
				d.Map(posten, func(posten funktion) d.Node {
					return d.Div(d.Class("karte"),
						d.H3(d.Text(posten.titel)),
						d.P(d.Text(posten.text)),
					)
				}),
			),
		),
	)
}
