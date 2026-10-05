package ui

import (
	g "github.com/LukasLow/dreego/v0/pkg/dom"

	"github.com/LukasLow/dreego/v0/pkg/dreego"
	"github.com/LukasLow/dreego/v0/pkg/scope"
)

// Field renders a labelled form input with its validation messages.
//
//	ui.Field(c, form, "email", "E-Mail", "email")
//
// value and errors come from the Form (Old/Errors); pass nil form for a fresh
// form. The label is associated with the input via for/id (accessibility).
func Field(c *dreego.Ctx, form *dreego.Form, name string, label string, typ string) g.View {
	alterWert := ""
	var fehler []string
	if form != nil {
		alterWert = form.Old(name)
		fehler = form.Errors(name)
	}

	return c.Box(
		scope.CSS(`
.u-field label { display: block; font-weight: 700; font-size: 14px; margin: 14px 0 6px }
.u-field input { width: 100%; box-sizing: border-box; font: inherit; font-size: 15px;
                 padding: 11px 13px; border: 1.6px solid #0080ff; border-radius: 11px;
                 background: #fff; color: #000 }
.u-field input:focus { outline: none; box-shadow: 0 0 0 3px rgba(0,128,255,.20) }
.u-field .err { color: #C62828; font-weight: 700; font-size: 13.5px; margin: 6px 0 0 }
`),
		g.Div(g.Class("u-field"),
			g.Label(g.Attr("for", name), g.Text(label)),
			g.Input(
				g.Attr("id", name), g.Attr("name", name), g.Attr("type", typ),
				g.Attr("value", alterWert),
			),
			g.Map(fehler, func(meldung string) g.View {
				return g.P(g.Class("err"), g.Attr("role", "alert"), g.Text(meldung))
			}),
		),
	)
}
