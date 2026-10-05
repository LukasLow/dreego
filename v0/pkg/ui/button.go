package ui

import (
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"

	"github.com/LukasLow/dreego/v0/pkg/dreego"
	"github.com/LukasLow/dreego/v0/pkg/scope"
)

// Button renders a link styled as a button. variant is Primary or Ghost.
//
//	ui.Button(c, "/pricing", "Preise ansehen", ui.Primary)
func Button(c *dreego.Ctx, href string, label string, variant Variant) g.Node {
	klasse := "u-btn u-btn-primary"
	if variant == Ghost {
		klasse = "u-btn u-btn-ghost"
	}

	return c.Box(
		scope.CSS(`
.u-btn { display: inline-block; font: inherit; font-weight: 700; cursor: pointer;
         padding: 11px 20px; border-radius: 10px; text-decoration: none; text-align: center;
         border: 2px solid #0080ff }
.u-btn-primary { background: #0080ff; color: #000 }
.u-btn-primary:hover { background: #37a1ff }
.u-btn-ghost { background: transparent; color: #003870 }
.u-btn-ghost:hover { background: #eef4ff }
`),
		h.A(h.Href(href), h.Class(klasse), g.Text(label)),
	)
}
