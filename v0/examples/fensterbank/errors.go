package main

import (
	d "github.com/LukasLow/dreego/v0/pkg/dom"

	"github.com/LukasLow/dreego/v0/pkg/dreego"
	"github.com/LukasLow/dreego/v0/pkg/scope"
)

// RegistriereFehlerseiten setzt eigene 404- und 500-Seiten im Fensterbank-Look.
func RegistriereFehlerseiten(app *dreego.App) {
	app.SetErrorPage(404, func(c *dreego.Ctx) d.Node {
		return fehlerSeite(c, 404, c.T("fehler.404.titel"), c.T("fehler.404.text"))
	})
	app.SetErrorPage(500, func(c *dreego.Ctx) d.Node {
		return fehlerSeite(c, 500, "Etwas ging schief",
			"Wir konnten die Seite nicht laden. Bitte später noch einmal versuchen.")
	})
}

// fehlerSeite rendert eine Fehlerseite IN der Site-Hülle (Shell), damit sie
// aussieht wie die normale Website — kein nackter Text, kein Layoutbruch.
func fehlerSeite(c *dreego.Ctx, status int, titel string, text string) d.Node {
	inhalt := c.Box(
		scope.CSS(`
.fehler { max-width: 560px; margin: 48px auto; text-align: center }
.fehler .zahl { font-size: 64px; font-weight: 800; color: #0b6b4a; margin: 0 }
.fehler h1 { font-size: 26px; margin: 8px 0 10px }
.fehler p { color: #3d5c4e }
.fehler .weg { margin-top: 18px; display: flex; gap: 12px; justify-content: center; flex-wrap: wrap }
.fehler a { font-weight: 700 }
`),
		d.Section(d.Class("fehler"),
			d.P(d.Class("zahl"), d.Text(itoa(status))),
			d.H1(d.Text(titel)),
			d.P(d.Text(text)),
			d.Div(d.Class("weg"),
				d.A(d.Href("/"), d.Text("Zur Startseite")),
				d.A(d.Href("/sorten"), d.Text("Zu den Sorten")),
			),
		),
	)

	return Shell(c, d.TitleEl(d.Text(itoa(status)+" — Fensterbank")), inhalt)
}

// itoa: kleine Zahl-zu-Text-Hilfe (Statuscodes sind dreistellig).
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
