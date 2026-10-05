package main

import (
	d "github.com/LukasLow/dreego/v0/pkg/dom"

	"github.com/LukasLow/dreego/v0/pkg/dreego"
	"github.com/LukasLow/dreego/v0/pkg/scope"
)

// errorPage zeigt einen technischen Fehler als Seite. In echt wuerde das ueber
// die Fehlermeldung laufen (infopech) — hier nur ein Minimalbeispiel.
func errorPage(c *dreego.Ctx, ursache error) d.Node {
	return c.Box(
		scope.CSS(`.errbox { border: 2px solid #C62828; border-radius: 12px; padding: 18px;
                    max-width: 560px; margin: 40px auto; background: #fff }
                    .errbox h1 { color: #C62828; font-size: 20px; margin: 0 0 8px }`),
		d.Div(d.Class("errbox"),
			d.H1(d.Text("Etwas ging schief")),
			d.P(d.Text(ursache.Error())),
		),
	)
}
