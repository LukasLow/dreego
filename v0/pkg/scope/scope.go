// Package scope is dreego's first addon: component-scoped CSS and JavaScript
// for the dom layer.
//
// A component writes its CSS and JS right next to its HTML, in the same Go
// function:
//
//	func Zaehler(c *scope.Collector) g.View {
//	    return c.Box(
//	        scope.CSS(`.zaehler { display: flex; gap: 12px }`),
//	        scope.JS(`
//	            var knopf = root.querySelector("button");
//	            knopf.addEventListener("click", function () {
//	                var wert = root.querySelector(".wert");
//	                wert.textContent = String(Number(wert.textContent) + 1);
//	            });
//	        `),
//	        Div(Class("zaehler"),
//	            Button(g.Text("+1")),
//	            Span(Class("wert"), g.Text("0")),
//	        ),
//	    )
//	}
//
// There are two ways to use a component:
//
//   - c.Box(...)  — registers the component with a per-page Collector. Document
//     then emits each component's CSS and JS exactly once, with the CSP nonce.
//     This is the recommended path for real pages.
//   - Box(...)    — the standalone form: emits the scoped CSS and JS inline,
//     right at the component. No deduplication. Good for quick demos.
//
// In both cases the component is wrapped in a container with a hash-derived
// data-scope attribute:
//
//   - CSS selectors are rewritten to match only inside that container.
//   - JS receives the container as `root`, so it can only touch its own
//     component — never another one, even with the same class names.
//
// Deliberately dependency-free: no runtime template, no CSS framework, no
// build step. Just Go.
package scope

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"strings"

	. "github.com/LukasLow/dreego/v0/pkg/dom"
	g "github.com/LukasLow/dreego/v0/pkg/dom"
)

// cssPart and jsPart are marker types. Box recognises them among its arguments
// and pulls them out of the body, instead of rendering them in place. They
// still implement dom.View so they can be passed as arguments; used
// outside Box, they render as a plain <style> / <script> without scoping.
type cssPart string
type jsPart string

func (teil cssPart) Render(w io.Writer) error {
	_, err := io.WriteString(w, "<style>"+escapeClosingTag(string(teil), "style")+"</style>")
	return err
}

func (teil jsPart) Render(w io.Writer) error {
	_, err := io.WriteString(w, "<script>"+escapeClosingTag(string(teil), "script")+"</script>")
	return err
}

// CSS marks a component's stylesheet. Box scopes it automatically.
func CSS(stylesheet string) g.View { return cssPart(stylesheet) }

// JS marks a component's client script. Box scopes it to the component root,
// which is passed to the script as the variable `root`.
func JS(script string) g.View { return jsPart(script) }

// Box wraps the component's parts in a scope container, emitting its CSS and JS
// inline (no deduplication). Use the Collector's Box method for dedupe.
func Box(parts ...g.View) g.View {
	css, js, body := splitParts(parts)
	scopeID := shortHash(css + "\x00" + js)

	kinder := []g.View{g.Attr("data-scope", scopeID)}

	if strings.TrimSpace(css) != "" {
		scoped := rewriteScoped(css, `[data-scope="`+scopeID+`"]`)
		kinder = append(kinder, StyleEl(g.Raw(escapeClosingTag(scoped, "style"))))
	}

	kinder = append(kinder, g.Group(body))

	if strings.TrimSpace(js) != "" {
		kinder = append(kinder, Script(g.Raw(escapeClosingTag(buildScript(js, scopeID), "script"))))
	}

	return Div(kinder...)
}

// splitParts separates a component's parts into its stylesheet, its script and
// its body nodes, in argument order.
func splitParts(parts []g.View) (string, string, []g.View) {
	var stylesheet strings.Builder
	var script strings.Builder
	var body []g.View

	for _, teil := range parts {
		switch wert := teil.(type) {
		case cssPart:
			stylesheet.WriteString(string(wert))
			stylesheet.WriteByte('\n')
		case jsPart:
			script.WriteString(string(wert))
			script.WriteByte('\n')
		default:
			body = append(body, teil)
		}
	}

	return stylesheet.String(), script.String(), body
}

// shortHash returns a short, stable id derived from the component's CSS and JS.
// The same component always yields the same id, so its styles and script can be
// recognised and emitted once across instances.
func shortHash(quelle string) string {
	summe := sha256.Sum256([]byte(quelle))
	return hex.EncodeToString(summe[:6])
}
