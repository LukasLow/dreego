package main

import (
	d "github.com/LukasLow/dreego/v0/pkg/dom"

	"github.com/LukasLow/dreego/v0/pkg/dreego"
	"github.com/LukasLow/dreego/v0/pkg/scope"
)

// Dashboard — geschützt: ohne Session zurück zum Login.
var Dashboard = dreego.Page{
	Path:   "/dashboard",
	Layout: Shell,
	Nav:    "dashboard",
	Get:    getDashboard,
}

func getDashboard(c *dreego.Ctx) d.View {
	email := c.SessionVal("autorin_email")
	if email == "" {
		return c.Redirect("/login", 303)
	}

	return c.Box(
		scope.CSS(`
.schreibtisch { margin: 36px auto; max-width: 560px }
.schreibtisch h1 { font-size: 28px }
.schreibtisch .blatt { border: 2px solid #e2d9c8; border-radius: 14px; padding: 22px; background: #fff }
.schreibtisch .blatt b { color: #b4451f }
`),
		d.Div(d.Class("schreibtisch"),
			d.H1(d.Text("Dein Schreibtisch")),
			d.Div(d.Class("blatt"),
				d.P(d.Text("Angemeldet als "), d.B(d.Text(email))),
				d.P(d.Text("Projekt: "), d.B(d.Text("Der Wintervogel")), d.Text(" — Kapitel 7")),
				d.P(d.A(d.Href("/logout"), d.Text("Abmelden"))),
			),
		),
	)
}
