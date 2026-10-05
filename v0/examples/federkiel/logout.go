package main

import (
	d "github.com/LukasLow/dreego/v0/pkg/dom"

	"github.com/LukasLow/dreego/v0/pkg/dreego"
	"github.com/LukasLow/dreego/v0/pkg/scope"
)

// Logout — löscht die Session. Bewusst POST + CSRF, damit kein fremdes Bild
// einen per <img src="/logout"> abmelden kann.
var Logout = dreego.Page{
	Path:   "/logout",
	Layout: Shell,
	Get:    getLogoutForm,
	Post:   postLogout,
}

func getLogoutForm(c *dreego.Ctx) d.Node {
	return c.Box(
		scope.CSS(`
.abmelden { max-width: 420px; margin: 44px auto; text-align: center }
.abmelden button { margin-top: 8px }
`),
		d.Div(d.Class("abmelden"),
			d.H1(d.Text("Abmelden?")),
			d.Form(d.Attr("method", "post"), d.Attr("action", "/logout"),
				c.CSRFInput(),
				d.Button(d.Type("submit"), d.Class("btn"), d.Text("Ja, abmelden")),
			),
		),
	)
}

func postLogout(c *dreego.Ctx) d.Node {
	c.DestroySession()
	return c.Redirect("/", 303)
}
