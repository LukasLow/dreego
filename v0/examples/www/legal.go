package main

import (
	d "github.com/LukasLow/dreego/v0/pkg/dom"

	"github.com/LukasLow/dreego/v0/pkg/dreego"
	"github.com/LukasLow/dreego/v0/pkg/markdown"
	"github.com/LukasLow/dreego/v0/pkg/scope"
)

// impressumText ist die Rechtsseite als Markdown (wie Dreego `<body lang="md">`).
const impressumText = `
# Impressum

Angaben gemäß § 5 DDG.

## Anbieter

Platzhalter — Betreiberangaben hier eintragen.

## Kontakt

E-Mail: [hallo@example.de](mailto:hallo@example.de)

## Umsatzsteuer

Platzhalter — USt-IdNr. oder Hinweis auf § 19 UStG.

---

Diese Seite wird aus **Markdown** gerendert, sicher und ohne Roh-HTML.
`

// Legal — eine Markdown-Seite.
var Legal = dreego.Page{
	Path:   "/impressum",
	Layout: Marketing,
	Nav:    "legal",
	Get:    getLegal,
}

func getLegal(c *dreego.Ctx) d.Node {
	return c.Box(
		scope.CSS(`
.legal { max-width: 720px; margin: 30px auto; padding: 0 18px }
.legal h1 { font-size: 30px }
.legal h2 { font-size: 20px; margin-top: 24px }
.legal a { color: #003870; font-weight: 700 }
`),
		d.Main(d.Attr("id", "main"), d.Class("legal"),
			d.Group(markdown.ToNodes(impressumText)),
		),
	)
}
