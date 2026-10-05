package ui

import (
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"

	"github.com/LukasLow/dreego/v0/pkg/dreego"
	"github.com/LukasLow/dreego/v0/pkg/scope"
)

// Alert renders a message box. tone is Info, Success, Warning or Danger. The
// role is "alert" for Warning/Danger (assertive), "status" otherwise — so screen
// readers announce it correctly.
//
//	ui.Alert(c, "Bitte Eingaben prüfen.", ui.Danger)
func Alert(c *dreego.Ctx, message string, tone Tone) g.Node {
	klasse := "u-alert u-alert-info"
	role := "status"
	switch tone {
	case Success:
		klasse = "u-alert u-alert-success"
	case Warning:
		klasse = "u-alert u-alert-warning"
		role = "alert"
	case Danger:
		klasse = "u-alert u-alert-danger"
		role = "alert"
	}

	return c.Box(
		scope.CSS(`
.u-alert { border: 1.6px solid; border-radius: 10px; padding: 11px 14px; font-weight: 700 }
.u-alert-info { background: #eef4ff; border-color: #0080ff; color: #003870 }
.u-alert-success { background: #e7f7ee; border-color: #0a7d2e; color: #0a6b2e }
.u-alert-warning { background: #fff4e0; border-color: #E8820C; color: #8a5000 }
.u-alert-danger { background: #fdecec; border-color: #C62828; color: #9c1c1c }
`),
		h.P(h.Class(klasse), g.Attr("role", role), g.Text(message)),
	)
}
