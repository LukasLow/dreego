package ui

import (
	g "github.com/LukasLow/dreego/v0/pkg/dom"

	"github.com/LukasLow/dreego/v0/pkg/dreego"
	"github.com/LukasLow/dreego/v0/pkg/scope"
)

// Card renders a titled box with arbitrary body content.
//
//	ui.Card(c, "Projekte", P(g.Text("Alles an einem Ort.")))
func Card(c *dreego.Ctx, titel string, body g.View) g.View {
	return c.Box(
		scope.CSS(`
.u-card { height: 100%; box-sizing: border-box; background: #fff; border: 1.6px solid #0080ff;
          border-radius: 16px; padding: 20px 22px }
.u-card h3 { margin: 0 0 8px; font-size: 18px }
.u-card p { margin: 0; color: #333 }
`),
		g.Div(g.Class("u-card"),
			g.H3(g.Text(titel)),
			body,
		),
	)
}
