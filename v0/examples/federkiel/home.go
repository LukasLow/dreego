package main

import (
	d "github.com/LukasLow/dreego/v0/pkg/dom"

	"github.com/LukasLow/dreego/v0/pkg/dreego"
	"github.com/LukasLow/dreego/v0/pkg/scope"
)

// Home — Startseite mit Hero (CSS+JS in der Komponente) und einem Zähler, der
// zeigt, dass Scoped-JS nur den eigenen Bereich anfasst.
var Home = dreego.Page{
	Path:   "/",
	Layout: Shell,
	Nav:    "home",
	Head: func(c *dreego.Ctx) []d.View {
		return []d.View{
			d.TitleEl(d.Text("Federkiel — schreiben ohne Ablenkung")),
			d.Meta(d.Name("description"), d.Content("Die ablenkungsfreie Schreib-App für Romanautorinnen.")),
		}
	},
	Get: getHome,
}

func getHome(c *dreego.Ctx) d.View {
	return d.Group([]d.View{
		c.Box(
			scope.CSS(`
.hero { text-align: center; padding: 56px 16px 40px }
.hero h1 { font-size: 46px; margin: 0 0 14px; line-height: 1.1; letter-spacing: -.5px }
.hero h1 em { color: #b4451f; font-style: normal }
.hero p { font-size: 20px; color: #5c5347; max-width: 560px; margin: 0 auto 26px }
.hero .row { display: flex; gap: 14px; justify-content: center; flex-wrap: wrap }
@media (max-width: 640px) { .hero h1 { font-size: 32px } .hero p { font-size: 17px } }
`),
			d.Section(d.Class("hero"),
				d.H1(d.Text("Schreiben. Nur "), d.Em(d.Text("schreiben")), d.Text(".")),
				d.P(d.Text("Federkiel hält alles fern, was dich vom nächsten Kapitel abhält.")),
				d.Div(d.Class("row"),
					d.A(d.Href("/login"), d.Class("btn"), d.Text("Kostenlos beginnen")),
					d.A(d.Href("/features"), d.Class("btn ghost"), d.Text("Was es kann")),
				),
			),
		),

		c.Box(
			scope.CSS(`
.zaehler { text-align: center; border: 2px solid #e2d9c8; border-radius: 14px; padding: 24px; max-width: 460px; margin: 10px auto 40px }
.zaehler .zahl { font-size: 40px; font-weight: 700; color: #b4451f }
.zaehler button { font: inherit; font-weight: 700; cursor: pointer; padding: 9px 18px; border: none; border-radius: 9px; background: #b4451f; color: #fff }
`),
			scope.JS(`
var knopf = root.querySelector("button");
var wert = root.querySelector(".wert");
knopf.addEventListener("click", function () {
    wert.textContent = String(Number(wert.textContent) + 137);
});
`),
			d.Div(d.Class("zaehler"),
				d.P(d.Text("Heutiges Schreibziel")),
				d.P(d.Class("zahl"), d.Span(d.Class("wert"), d.Text("0")), d.Text(" Wörter")),
				d.Button(d.Text("+137 Wörter")),
			),
		),
	})
}
