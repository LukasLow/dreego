package web

import (
	d "github.com/LukasLow/dreego/v0/pkg/dom"

	"github.com/LukasLow/dreego/v0/pkg/dreego"
	"github.com/LukasLow/dreego/v0/pkg/scope"
)

// bindingPath is the runtime JavaScript binding served by Wails. It mirrors
// GreeterService.Greet in app/greeter.go (module wails-dreego).
const bindingPath = "/bindings/wails-dreego/app/greeterservice.js"

// greeterJS calls the generated Wails binding. It is scoped to the component
// root (as `root`) and inlined with a CSP nonce.
const greeterJS = `
const bindingPath = "` + bindingPath + `";
const button = root.querySelector("#greet");
const input = root.querySelector("#name");
const result = root.querySelector("#result");

async function greet() {
  if (!button || !input || !result) { return; }
  button.disabled = true;
  try {
    const service = await import(bindingPath);
    result.textContent = await service.Greet(input.value);
  } catch (error) {
    result.textContent = "Greeter service unavailable.";
    console.error(error);
  } finally {
    button.disabled = false;
  }
}

button?.addEventListener("click", greet);
input?.addEventListener("keydown", (event) => {
  if (event.key === "Enter") { greet(); }
});
`

// Home is the start page.
var Home = dreego.Page{
	Path:   "/",
	Layout: Shell,
	Head:   func(c *dreego.Ctx) []d.View { return []d.View{d.TitleEl(d.Text("Wails + dreego"))} },
	Get:    getHome,
}

func getHome(c *dreego.Ctx) d.View {
	c.AddCritical(shellCSS)

	greet := c.Box(
		scope.CSS(`
.greet { display:grid; gap:.5rem; }
.greet label { font-size:.78rem; font-weight:650; color:#9aa3b8; }
.greet .row { display:flex; gap:.6rem; }
.greet input { flex:1; min-width:0; padding:.65rem .8rem; border:1px solid rgba(255,255,255,.16);
  border-radius:.7rem; background:rgba(0,0,0,.25); color:#e7ebf5; font:inherit; }
.greet button { padding:.65rem 1.1rem; border:0; border-radius:.7rem; background:#9fffd7;
  color:#07130f; font:inherit; font-weight:750; cursor:pointer; }
.greet .result { min-height:1.4rem; margin:.3rem 0 0; font-size:.95rem; font-weight:650; color:#9fffd7; }
`),
		scope.JS(greeterJS),
		d.Div(d.Class("greet"),
			d.Label(d.Attr("for", "name"), d.Text("Your name")),
			d.Div(d.Class("row"),
				d.Input(d.Attr("id", "name"), d.Name("name"), d.Type("text"),
					d.Attr("placeholder", "world"), d.Attr("autocomplete", "off")),
				d.Button(d.Attr("id", "greet"), d.Type("button"), d.Text("Greet")),
			),
			d.P(d.Attr("id", "result"), d.Class("result"),
				d.Attr("role", "status"), d.Attr("aria-live", "polite")),
		),
	)

	return d.Group([]d.View{
		d.A(d.Href("#main"), d.Class("skip-link"), d.Text("skip to content")),
		d.Main(d.Class("shell"), d.Attr("id", "main"), d.Attr("tabindex", "-1"),
			d.Header(d.Class("topbar"),
				d.Span(d.Class("brand"),
					d.Span(d.Class("brand-mark"), d.Attr("aria-hidden", "true"), d.Text("D")),
					d.Text(" dreego"),
				),
				d.Span(d.Class("native-badge"), d.Text("Native · Wails v3")),
			),

			d.Section(d.Class("card"), d.Attr("aria-labelledby", "page-title"),
				d.P(d.Class("eyebrow"), d.Text("In-process render")),
				d.H1(d.Attr("id", "page-title"), d.Text("Go backend, dreego frontend.")),
				d.P(d.Class("intro"),
					d.Text("No npm, no local server. This page is a Go function, served by the dreego Wails adapter."),
				),
				greet,
			),

			d.Footer(
				d.Span(d.Text("Private by design · no local server")),
				d.A(d.Href("/about"), d.Text("How it works →")),
			),
		),
	})
}
