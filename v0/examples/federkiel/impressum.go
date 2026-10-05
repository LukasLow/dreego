package main

import (
	d "github.com/LukasLow/dreego/v0/pkg/dom"

	"github.com/LukasLow/dreego/v0/pkg/dreego"
	"github.com/LukasLow/dreego/v0/pkg/markdown"
	"github.com/LukasLow/dreego/v0/pkg/scope"
)

// impressumMd ist die Rechtsseite als Markdown (Dreego: `lang="md"`).
const impressumMd = `
# Impressum

Angaben gemäß § 5 DDG. Federkiel ist eine **erfundene Demo-Anwendung**.

## Anbieter

Federkiel Verlag (Fiktion)
Beispielgasse 1
12345 Musterstadt

## Kontakt

E-Mail: [hallo@federkiel.example](mailto:hallo@federkiel.example)

## Umsatzsteuer

USt-IdNr.: DE000000000 (Platzhalter)

## Verantwortlich für den Inhalt

F. Kiel, Adresse wie oben

---

Diese Seite wird aus **Markdown** gerendert — mit Überschriften, Listen und
Links, ganz ohne Roh-HTML.
`

// Impressum — Markdown-Seite.
var Impressum = dreego.Page{
	Path:   "/impressum",
	Layout: Shell,
	Nav:    "impressum",
	Head: func(c *dreego.Ctx) []d.Node {
		return []d.Node{d.TitleEl(d.Text("Impressum — Federkiel"))}
	},
	Get: getImpressum,
}

func getImpressum(c *dreego.Ctx) d.Node {
	return c.Box(
		scope.CSS(`
.recht { max-width: 720px; margin: 34px auto }
.recht h1 { font-size: 32px }
.recht h2 { font-size: 20px; margin-top: 26px }
`),
		d.Div(d.Class("recht"),
			d.Group(markdown.ToNodes(impressumMd)),
		),
	)
}
