package main

import (
	d "github.com/LukasLow/dreego/v0/pkg/dom"

	"github.com/LukasLow/dreego/v0/pkg/dreego"
	"github.com/LukasLow/dreego/v0/pkg/scope"
)

// LoginForm ist die Eingabe — das `in`. Die Tags steuern Binding und Prüfung.
type LoginForm struct {
	Email string `form:"email" validate:"required,email" label:"E-Mail"`
	Team  string `form:"team" validate:"required,min=2" label:"Teamname"`
}

// Login — GET zeigt das Formular (mit CSRF-Feld), POST verarbeitet es.
var Login = dreego.Page{
	Path:   "/login",
	Layout: Auth,
	Nav:    "login",
	Get:    getLogin,
	Post:   postLogin,
}

func getLogin(c *dreego.Ctx) d.View {
	// Flash aus dem vorherigen POST (PRG) — genau einmal sichtbar.
	return loginForm(c, nil, c.FlashGet("login_error"))
}

func postLogin(c *dreego.Ctx) d.View {
	in, form, err_bind := dreego.Bind[LoginForm](c)
	if err_bind != nil {
		c.Flash("login_error", "Bitte Eingaben prüfen.")
		return c.Redirect("/login", 303)
	}

	if form.HasErrors() {
		// Ungültig: PRG zurück zum Formular. Die Meldung reist per Flash.
		c.Flash("login_error", "Bitte Eingaben prüfen.")
		return c.Redirect("/login", 303)
	}

	// Gültig: Session setzen, dann Redirect aufs Dashboard.
	c.SetSessionVal("user_email", in.Email)
	c.SetSessionVal("user_team", in.Team)
	return c.Redirect("/dashboard", 303)
}

func loginForm(c *dreego.Ctx, form *dreego.Form, meldung string) d.View {
	return c.Box(
		scope.CSS(`
.authcard { max-width: 440px; margin: 40px auto; border: 2px solid #0080ff;
            border-radius: 16px; padding: 26px; background: #fff }
.authcard h1 { font-size: 22px; margin: 0 0 16px }
.authcard label { display: block; font-weight: 700; margin: 14px 0 6px }
.authcard input { width: 100%; font: inherit; padding: 10px 12px; box-sizing: border-box;
                  border: 2px solid #0080ff; border-radius: 10px }
.authcard .err { color: #C62828; font-weight: 700; font-size: 14px; margin: 6px 0 0 }
.authcard button { margin-top: 18px; width: 100%; font: inherit; font-weight: 700;
                   padding: 11px 16px; border: none; border-radius: 10px;
                   background: #0080ff; color: #000; cursor: pointer }
`),
		d.Div(d.Class("authcard"),
			d.H1(d.Text("Beitreten")),
			flashBox(meldung),
			d.Form(d.Attr("method", "post"), d.Attr("action", "/login"),
				c.CSRFInput(),
				formField(form, "team", "text", "Teamname"),
				formField(form, "email", "email", "E-Mail"),
				d.Button(d.Attr("type", "submit"), d.Text("Absenden")),
			),
		),
	)
}

// formField rendert ein Feld mit Old-Wert und Fehlermeldungen.
func formField(form *dreego.Form, name string, typ string, label string) d.View {
	alterWert := ""
	var fehler []string
	if form != nil {
		alterWert = form.Old(name)
		fehler = form.Errors(name)
	}

	return d.Group([]d.View{
		d.Label(d.Attr("for", name), d.Text(label)),
		d.Input(d.Attr("id", name), d.Attr("name", name), d.Attr("type", typ),
			d.Attr("value", alterWert)),
		d.Map(fehler, func(meldung string) d.View {
			return d.P(d.Class("err"), d.Attr("role", "alert"), d.Text(meldung))
		}),
	})
}

func flashBox(meldung string) d.View {
	if meldung == "" {
		return d.Group(nil)
	}
	return d.P(d.Class("err"), d.Attr("role", "alert"), d.Text(meldung))
}
