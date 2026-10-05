package web

import (
	d "github.com/LukasLow/dreego/v0/pkg/dom"

	"github.com/LukasLow/dreego/v0/pkg/dreego"
	"github.com/LukasLow/dreego/v0/pkg/scope"
)

// About explains the design.
var About = dreego.Page{
	Path:   "/about",
	Layout: Shell,
	Head:   func(c *dreego.Ctx) []d.View { return []d.View{d.TitleEl(d.Text("About — Wails + dreego"))} },
	Get:    getAbout,
}

func getAbout(c *dreego.Ctx) d.View {
	c.AddCritical(shellCSS)

	zeile := func(bezeichnung string, wert string) d.View {
		return d.Div(
			d.Dt(d.Text(bezeichnung)),
			d.Dd(d.Text(wert)),
		)
	}

	return d.Group([]d.View{
		d.A(d.Href("#main"), d.Class("skip-link"), d.Text("skip to content")),
		d.Main(d.Class("shell"), d.Attr("id", "main"), d.Attr("tabindex", "-1"),
			c.Box(
				scope.CSS(`
.about .back { display:inline-flex; align-items:center; gap:.45rem; color:#9fffd7; font-size:.78rem; font-weight:750; text-decoration:none; }
.about article { max-width:32rem; margin:clamp(2rem,8vh,4rem) auto 0; }
.about dl { display:grid; grid-template-columns:1fr 1fr; gap:.6rem; margin:0; }
.about dl div { padding:.85rem 1rem; border:1px solid rgba(255,255,255,.1); border-radius:.9rem; background:rgba(255,255,255,.045); }
.about dt { margin-bottom:.3rem; color:#747d97; font-size:.64rem; font-weight:700; letter-spacing:.1em; text-transform:uppercase; }
.about dd { margin:0; font-size:.85rem; font-weight:700; }
.about .note { display:flex; gap:.7rem; margin:1.2rem 0 0; padding:1rem; border-left:2px solid #b3a6ff;
  background:rgba(179,166,255,.07); color:#a4abc0; font-size:.78rem; line-height:1.55; }
`),
				d.Div(d.Class("about"),
					d.A(d.Class("back"), d.Href("/"), d.Text("← Back")),
					d.Article(
						d.P(d.Class("eyebrow"), d.Text("Under the hood")),
						d.H1(d.Text("Small surface. Clear ownership.")),
						d.P(d.Class("intro"),
							d.Text("A native desktop window powered by Wails, with dreego handling only what it should: rendering pages and serving assets in-process.")),
						d.Dl(
							zeile("Transport", "In-process handler"),
							zeile("Local ports", "None"),
							zeile("UI source", "Go (dreego)"),
							zeile("State owner", "Explicit Go service"),
						),
						d.P(d.Class("note"),
							d.Text("The application owns its Wails window, services, bindings and lifecycle. dreego is the adapter — not the wrapper.")),
					),
				),
			),
		),
	})
}
