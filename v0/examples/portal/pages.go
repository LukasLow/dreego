package main

import (
	d "github.com/LukasLow/dreego/v0/pkg/dom"

	"github.com/LukasLow/dreego/v0/pkg/dreego"
	"github.com/LukasLow/dreego/v0/pkg/scope"
)

// Register hängt die Portal-Seiten an die App.
func Register(app *dreego.App) {
	app.Page(Start)
	app.Page(Login)
	app.Page(Konto)
	app.Page(Logout)
}

// Start — die Portal-Startseite. Das Logo und "/" dürfen nicht ins Leere
// führen: ist man angemeldet, geht es zum Konto, sonst zur Anmeldung.
var Start = dreego.Page{
	Path:   "/",
	Layout: Shell,
	Head:   func(c *dreego.Ctx) []d.View { return []d.View{d.TitleEl(d.Text("Portal"))} },
	Get: func(c *dreego.Ctx) d.View {
		if c.SessionVal("konto_email") != "" {
			return c.Redirect(markt.Abs("/konto"), 303)
		}
		return c.Redirect(markt.Abs("/login"), 303)
	},
}

// LoginForm ist die Eingabe (das `in`).
type LoginForm struct {
	Email string `form:"email" validate:"required,email" label:"E-Mail"`
}

// Login — GET zeigt das Formular, POST legt die Session an.
var Login = dreego.Page{
	Path:   "/login",
	Layout: Shell,
	Head:   func(c *dreego.Ctx) []d.View { return []d.View{d.TitleEl(d.Text("Anmelden — Portal"))} },
	Get:    getLogin,
	Post:   postLogin,
}

func getLogin(c *dreego.Ctx) d.View {
	return form(c, c.FlashGet("fehler"))
}

func postLogin(c *dreego.Ctx) d.View {
	in, formWerte, err_bind := dreego.Bind[LoginForm](c)
	if err_bind != nil || formWerte.HasErrors() {
		c.Flash("fehler", "Bitte eine gültige E-Mail eingeben.")
		return c.Redirect(markt.Abs("/login"), 303)
	}
	c.SetSessionVal("konto_email", in.Email)
	return c.Redirect(markt.Abs("/konto"), 303)
}

func form(c *dreego.Ctx, meldung string) d.View {
	return c.Box(
		scope.CSS(`
.karte { max-width: 420px; margin: 44px auto; background: #fff; border: 2px solid #d3d9e8;
         border-radius: 16px; padding: 26px }
.karte h1 { margin: 0 0 16px; font-size: 24px }
.karte label { display: block; font-weight: 700; margin: 12px 0 6px }
.karte input { width: 100%; font: inherit; padding: 10px 12px; border: 2px solid #274bdb; border-radius: 10px }
.karte button { margin-top: 18px; width: 100%; font: inherit; font-weight: 700; cursor: pointer;
                padding: 11px 16px; border: none; border-radius: 10px; background: #274bdb; color: #fff }
.fehler { background: #fde8e8; border: 2px solid #c62828; border-radius: 10px; padding: 10px 14px; font-weight: 700 }
`),
		d.Div(d.Class("karte"),
			d.H1(d.Text("Anmelden")),
			flash(meldung),
			d.Form(d.Attr("method", "post"), d.Attr("action", "/login"),
				c.CSRFInput(),
				d.Label(d.Attr("for", "email"), d.Text("E-Mail")),
				d.Input(d.Attr("id", "email"), d.Attr("name", "email"), d.Attr("type", "email")),
				d.Button(d.Attr("type", "submit"), d.Text("Weiter")),
			),
		),
	)
}

func flash(meldung string) d.View {
	if meldung == "" {
		return d.Group(nil)
	}
	return d.P(d.Class("fehler"), d.Attr("role", "alert"), d.Text(meldung))
}

// Konto — geschützte Seite: ohne Session zurück zum Login.
var Konto = dreego.Page{
	Path:   "/konto",
	Layout: Shell,
	Head:   func(c *dreego.Ctx) []d.View { return []d.View{d.TitleEl(d.Text("Konto — Portal"))} },
	Get:    getKonto,
}

func getKonto(c *dreego.Ctx) d.View {
	email := c.SessionVal("konto_email")
	if email == "" {
		return c.Redirect(markt.Abs("/login"), 303)
	}
	return c.Box(
		scope.CSS(`.karte { max-width: 520px; margin: 40px auto; background: #fff; border: 2px solid #d3d9e8; border-radius: 16px; padding: 24px }
                   .karte b { color: #274bdb }`),
		d.Div(d.Class("karte"),
			d.H1(d.Text("Dein Konto")),
			d.P(d.Text("Angemeldet als "), d.B(d.Text(email))),
			d.P(d.A(d.Href(public.Abs("/sorten")), d.Text("Zur öffentlichen Website (anderer Port)"))),
			d.P(d.A(d.Href(markt.Abs("/logout")), d.Text("Abmelden"))),
		),
	)
}

// Logout — POST + CSRF, damit kein fremdes Bild abmelden kann.
var Logout = dreego.Page{
	Path:   "/logout",
	Layout: Shell,
	Get: func(c *dreego.Ctx) d.View {
		return c.Box(
			scope.CSS(`.karte { max-width: 420px; margin: 44px auto; text-align: center }`),
			d.Div(d.Class("karte"),
				d.H1(d.Text("Abmelden?")),
				d.Form(d.Attr("method", "post"), d.Attr("action", "/logout"),
					c.CSRFInput(),
					d.Button(d.Attr("type", "submit"), d.Text("Ja, abmelden")),
				),
			),
		)
	},
	Post: func(c *dreego.Ctx) d.View {
		c.DestroySession()
		return c.Redirect(markt.Abs("/"), 303)
	},
}
