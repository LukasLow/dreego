package frontend

import (
	d "github.com/LukasLow/dreego/v0/pkg/dom"

	"github.com/LukasLow/dreego/v0/pkg/dreego"
)

// Inline SVG icons, taken verbatim from the Wails vanilla template. Inline
// (rather than files) because the page is rendered in Go.
const (
	iconUser  = `<svg class="input-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/><circle cx="12" cy="7" r="4"/></svg>`
	iconArrow = `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><line x1="5" y1="12" x2="19" y2="12"/><polyline points="12 5 19 12 12 19"/></svg>`
	iconClock = `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/></svg>`
	iconDoc   = `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><line x1="7" y1="17" x2="17" y2="7"/><polyline points="7 7 17 7 17 17"/></svg>`
)

// WailsVersion is shown in the footer (the vanilla template injected it too).
const WailsVersion = "v3.0.0-beta.27"

// Index is the template's single page.
var Index = dreego.Page{
	Path:   "/",
	Layout: Shell,
	Head: func(c *dreego.Ctx) []d.View {
		return []d.View{d.TitleEl(d.Text("Wails App"))}
	},
	Get: getIndex,
}

// getIndex renders the page: background, brand row, heading, the input+button
// row, the footer (with the live clock), the toast and the client script.
func getIndex(c *dreego.Ctx) d.View {
	return d.Group([]d.View{
		// Fixed full-bleed background layer (see style.css).
		d.Div(d.Class("bg"), d.Attr("aria-hidden", "true")),

		d.Main(d.Class("container"),
			// Brand row: Wails logo + JavaScript badge.
			d.Header(d.Class("brand"),
				d.A(d.Class("brand-mark"), d.Attr("data-wml-openURL", "https://v3.wails.io"), d.Attr("aria-label", "Wails website"),
					d.Img(d.Src("/static/wails.png"), d.Class("brand-logo"), d.Alt("Wails logo")),
				),
				d.A(d.Class("brand-badge"), d.Attr("data-wml-openURL", "https://developer.mozilla.org/en-US/docs/Web/JavaScript"), d.Attr("aria-label", "JavaScript"),
					d.Img(d.Src("/static/javascript.svg"), d.Alt("JavaScript logo")),
				),
			),

			d.H1(d.Class("title"),
				d.Span(d.Class("title-accent"), d.Text("Wails +")),
				d.Text(" "),
				d.Span(d.Class("title-name"), d.Span(d.Class("title-name-text"), d.Text("dreego"))),
			),
			d.P(d.Class("subtitle"), d.Text("Build beautiful cross-platform apps with Go and dreego.")),

			// Greet input + button.
			d.Div(d.Class("greet"),
				d.Div(d.Class("input-box"), d.ID("input"),
					d.Raw(iconUser),
					d.Input(d.Class("input"), d.ID("name"), d.Type("text"), d.Placeholder("Your name"), d.AutoComplete("off"), d.Attr("aria-label", "input")),
					d.Button(d.Class("btn"), d.ID("greet"), d.Attr("aria-label", "greet-btn"),
						d.Text("Greet "), d.Raw(iconArrow),
					),
				),
			),
		),

		d.Hr(d.Class("footer-divider")),
		d.Footer(d.Class("footer"),
			d.Span(d.Class("footer-version"), d.Span(d.ID("version"), d.Text(WailsVersion))),
			d.Span(d.Class("footer-time"), d.Raw(iconClock),
				d.Span(d.ID("time"), d.Text("Listening for Time event...")),
			),
			d.A(d.Class("footer-docs"), d.Attr("data-wml-openURL", "https://v3.wails.io"), d.Attr("aria-label", "Wails documentation"),
				d.Text("Docs "), d.Raw(iconDoc),
			),
		),

		// Toast: shows the greeting returned by the Go backend. A fixed overlay,
		// so it never reflows the layout.
		d.Div(d.Class("toast"), d.ID("toast"), d.Attr("role", "status"), d.Attr("aria-live", "polite"),
			d.Span(d.Class("toast-label"), d.Text("From Go")),
			d.Span(d.Class("toast-msg"), d.ID("result"), d.Attr("aria-label", "result")),
		),

		// The client module: plain JavaScript (the port of the vanilla main.ts),
		// emitted as a real <script type="module"> with the per-request CSP
		// nonce so it can use ES imports of the runtime and the generated binding.
		d.Script(d.Type("module"), d.Attr("nonce", c.Nonce()), d.Raw(clientJS)),
	})
}
