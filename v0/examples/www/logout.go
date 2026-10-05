package main

import (
	d "github.com/LukasLow/dreego/v0/pkg/dom"

	"github.com/LukasLow/dreego/v0/pkg/dreego"
)

// Logout — löscht die Session und leitet zur Startseite.
// Hinweis: ein Logout per GET ist per CSRF auslösbar (<img src="/logout">).
// Akzeptiert als bewusste Vereinfachung; für Produktion POST + CSRF-Token.
var Logout = dreego.Page{
	Path: "/logout",
	Get:  getLogout,
}

func getLogout(c *dreego.Ctx) d.View {
	c.DestroySession()
	return c.Redirect("/", 303)
}
