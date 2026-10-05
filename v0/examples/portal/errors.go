package main

import (
	d "github.com/LukasLow/dreego/v0/pkg/dom"

	"github.com/LukasLow/dreego/v0/pkg/dreego"
	"github.com/LukasLow/dreego/v0/pkg/scope"
)

// RegistriereFehlerseiten setzt eigene 404- und 500-Seiten im Portal-Look.
func RegistriereFehlerseiten(app *dreego.App) {
	app.SetErrorPage(404, func(c *dreego.Ctx) d.View {
		return fehlerSeite(c, 404, "Nicht gefunden",
			"Diese Adresse gehört nicht zum Portal.")
	})
	app.SetErrorPage(500, func(c *dreego.Ctx) d.View {
		return fehlerSeite(c, 500, "Etwas ging schief",
			"Bitte später noch einmal versuchen.")
	})
}

func fehlerSeite(c *dreego.Ctx, status int, titel string, text string) d.View {
	inhalt := c.Box(
		scope.CSS(`
.fehler { max-width: 520px; margin: 48px auto; text-align: center }
.fehler .zahl { font-size: 60px; font-weight: 800; color: #274bdb; margin: 0 }
.fehler h1 { font-size: 24px; margin: 8px 0 10px }
.fehler a { font-weight: 700 }
`),
		d.Section(d.Class("fehler"),
			d.P(d.Class("zahl"), d.Text(itoa(status))),
			d.H1(d.Text(titel)),
			d.P(d.Text(text)),
			d.P(d.A(d.Href(markt.Abs("/")), d.Text("Zurück zum Portal"))),
		),
	)
	return Shell(c, d.TitleEl(d.Text(itoa(status)+" — Portal")), inhalt)
}

func itoa(zahl int) string {
	if zahl == 0 {
		return "0"
	}
	var ziffern []byte
	for zahl > 0 {
		ziffern = append([]byte{byte('0' + zahl%10)}, ziffern...)
		zahl /= 10
	}
	return string(ziffern)
}
