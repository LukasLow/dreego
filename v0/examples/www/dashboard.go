package main

import (
	d "github.com/LukasLow/dreego/v0/pkg/dom"

	"github.com/LukasLow/dreego/v0/pkg/dreego"
	"github.com/LukasLow/dreego/v0/pkg/scope"
)

// Dashboard — eine geschützte Seite. Ohne Session leitet sie zum Login.
var Dashboard = dreego.Page{
	Path:   "/dashboard",
	Layout: Marketing,
	Nav:    "dashboard",
	Get:    getDashboard,
}

func getDashboard(c *dreego.Ctx) d.View {
	email := c.SessionVal("user_email")
	if email == "" {
		// Kein Login: zurück zum Formular (der Guard ist hier die Seite selbst,
		// weil dreego noch keine Middleware pro Route kennt).
		return c.Redirect("/login", 303)
	}

	return c.Box(
		scope.CSS(`
.dash { max-width: 560px; margin: 40px auto }
.dash h1 { font-size: 24px }
.dash .karte { border: 2px solid #0080ff; border-radius: 14px; padding: 20px; background: #fff }
.dash a { color: #003870; font-weight: 700 }
`),
		d.Div(d.Class("dash"),
			d.H1(d.Text("Dashboard")),
			d.Div(d.Class("karte"),
				d.P(d.Text("Angemeldet als "), d.B(d.Text(email))),
				d.P(d.A(d.Href("/logout"), d.Text("Abmelden"))),
			),
		),
	)
}
