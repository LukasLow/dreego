// Package components sammelt wiederverwendbare Bausteine der www-Site.
// Jede Komponente ist eine Go-Funktion, die einen Komponenten-Scope über den
// Request-Context anmeldet (c.Box) — HTML, CSS und JS in einer Funktion.
package components

import (
	d "github.com/LukasLow/dreego/v0/pkg/dom"

	"github.com/LukasLow/dreego/v0/pkg/dreego"
	"github.com/LukasLow/dreego/v0/pkg/scope"
)

// Hero — Startseiten-Kopf mit CSS und einem kleinen Client-Skript in EINER
// Funktion. Das Skript bekommt den Komponenten-Root als `root`.
func Hero(c *dreego.Ctx) d.Node {
	return c.Box(
		scope.CSS(`
.hero { text-align: center; padding: 48px 20px }
.hero h1 { font-size: 40px; margin: 0 0 12px; color: #111 }
.hero p { font-size: 18px; color: #333; max-width: 560px; margin: 0 auto }
.hero button { margin-top: 20px; font: inherit; font-weight: 700; padding: 12px 22px;
               border: none; border-radius: 10px; background: #0080ff; color: #000; cursor: pointer }
@media (max-width: 640px) { .hero h1 { font-size: 30px } }
`),
		scope.JS(`
var knopf = root.querySelector("button");
var zaehler = root.querySelector(".klicks");
knopf.addEventListener("click", function () {
    zaehler.textContent = String(Number(zaehler.textContent) + 1);
});
`),
		d.Section(d.Class("hero"),
			d.H1(d.Text("Ist alles bereit?")),
			d.P(d.Text("Eine Seite in reinem Go — HTML, CSS und JS in einer Funktion.")),
			d.Button(d.Text("Klick mich")),
			d.P(d.Text("Klicks: "), d.Span(d.Class("klicks"), d.Text("0"))),
		),
	)
}
