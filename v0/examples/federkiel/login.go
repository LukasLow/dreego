package main

import (
	d "github.com/LukasLow/dreego/v0/pkg/dom"

	"github.com/LukasLow/dreego/v0/pkg/dreego"
	"github.com/LukasLow/dreego/v0/pkg/scope"
)

// LoginForm ist die Eingabe (das `in`).
type LoginForm struct {
	Email string `form:"email" validate:"required,email" label:"E-Mail"`
	Stift string `form:"stift" validate:"required,min=3" label:"Lieblingsstift"`
}

// Login — GET zeigt das Formular, POST prüft und legt die Session an.
var Login = dreego.Page{
	Path:   "/login",
	Layout: Shell,
	Nav:    "login",
	Get:    getLogin,
	Post:   postLogin,
}

func getLogin(c *dreego.Ctx) d.Node {
	return loginKarte(c, nil, c.FlashGet("login_fehler"))
}

func postLogin(c *dreego.Ctx) d.Node {
	in, form, err_bind := dreego.Bind[LoginForm](c)
	if err_bind != nil {
		c.Flash("login_fehler", "Bitte Eingaben prüfen.")
		return c.Redirect("/login", 303)
	}

	if form.HasErrors() {
		c.Flash("login_fehler", "Bitte Eingaben prüfen.")
		return c.Redirect("/login", 303)
	}

	c.SetSessionVal("autorin_email", in.Email)
	c.SetSessionVal("autorin_stift", in.Stift)
	return c.Redirect("/dashboard", 303)
}

func loginKarte(c *dreego.Ctx, form *dreego.Form, meldung string) d.Node {
	return c.Box(
		scope.CSS(`
.login { max-width: 440px; margin: 44px auto; border: 2px solid #e2d9c8; border-radius: 16px; padding: 28px; }
.login h1 { font-size: 26px; margin: 0 0 6px }
.login .lead { color: #5c5347; margin: 0 0 20px }
.login label { display: block; font-weight: 700; margin: 14px 0 6px }
.login input { width: 100%; font: inherit; padding: 10px 12px; border: 2px solid #b4451f; border-radius: 10px; }
.login .fehler { color: #9c1c1c; font-weight: 700; margin: 4px 0 0 }
.login .flash { background: #fbe9e4; border: 2px solid #b4451f; border-radius: 10px; padding: 10px 14px; margin: 0 0 14px; font-weight: 700 }
.login button { margin-top: 20px; width: 100%; }
`),
		d.Div(d.Class("login"),
			d.H1(d.Text("Anmelden")),
			d.P(d.Class("lead"), d.Text("Federkiel fragt nur nach dem Nötigsten.")),
			flash(meldung),
			d.Form(d.Attr("method", "post"), d.Attr("action", "/login"),
				c.CSRFInput(),
				feld(form, "email", "email", "E-Mail"),
				feld(form, "stift", "text", "Lieblingsstift (Geheimnis)"),
				d.Button(d.Type("submit"), d.Class("btn"), d.Text("Weiter")),
			),
		),
	)
}

func feld(form *dreego.Form, name string, typ string, label string) d.Node {
	alterWert := ""
	var fehler []string
	if form != nil {
		alterWert = form.Old(name)
		fehler = form.Errors(name)
	}
	return d.Group([]d.Node{
		d.Label(d.Attr("for", name), d.Text(label)),
		d.Input(d.Attr("id", name), d.Attr("name", name), d.Attr("type", typ), d.Attr("value", alterWert)),
		d.Map(fehler, func(meldung string) d.Node {
			return d.P(d.Class("fehler"), d.Attr("role", "alert"), d.Text(meldung))
		}),
	})
}

func flash(meldung string) d.Node {
	if meldung == "" {
		return d.Group(nil)
	}
	return d.P(d.Class("flash"), d.Attr("role", "alert"), d.Text(meldung))
}
