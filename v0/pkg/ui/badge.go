package ui

import (
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"

	"github.com/LukasLow/dreego/v0/pkg/dreego"
	"github.com/LukasLow/dreego/v0/pkg/scope"
)

// Badge renders a small status label. tone is Info, Success, Warning or Danger.
//
//	ui.Badge(c, "Aktiv", ui.Success)
func Badge(c *dreego.Ctx, text string, tone Tone) g.Node {
	klasse := "u-badge u-badge-info"
	switch tone {
	case Success:
		klasse = "u-badge u-badge-success"
	case Warning:
		klasse = "u-badge u-badge-warning"
	case Danger:
		klasse = "u-badge u-badge-danger"
	}

	return c.Box(
		scope.CSS(`
.u-badge { display: inline-block; font-size: 13px; font-weight: 700; padding: 3px 10px;
           border-radius: 999px; border: 1.6px solid }
.u-badge-info { background: #eef4ff; color: #003870; border-color: #0080ff }
.u-badge-success { background: #e7f7ee; color: #0a6b2e; border-color: #0a7d2e }
.u-badge-warning { background: #fff4e0; color: #8a5000; border-color: #E8820C }
.u-badge-danger { background: #fdecec; color: #9c1c1c; border-color: #C62828 }
`),
		h.Span(h.Class(klasse), g.Text(text)),
	)
}
