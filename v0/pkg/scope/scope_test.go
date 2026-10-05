package scope

import (
	"bytes"
	"strings"
	"testing"

	. "github.com/LukasLow/dreego/v0/pkg/dom"
	g "github.com/LukasLow/dreego/v0/pkg/dom"
)

func render(t *testing.T, knoten g.View) string {
	t.Helper()
	var puffer bytes.Buffer
	err := knoten.Render(&puffer)
	if err != nil {
		t.Fatalf("render: %v", err)
	}
	return puffer.String()
}

func TestBoxScopesSelectors(t *testing.T) {
	css := ".hero h1 { color: red }"
	html := render(t, Box(CSS(css), Section(Class("hero"), H1(g.Text("Hallo")))))

	id := shortHash(css + "\n\x00")
	if !strings.Contains(html, `data-scope="`+id+`"`) {
		t.Fatalf("scope attribute fehlt: %s", html)
	}
	if !strings.Contains(html, `[data-scope="`+id+`"] .hero h1`) {
		t.Fatalf("selektor nicht gescoped: %s", html)
	}
	if strings.Contains(html, "\n.hero h1 {") {
		t.Fatalf("unscopeter selektor noch vorhanden: %s", html)
	}
}

func TestBoxKeepsAtRules(t *testing.T) {
	css := "@keyframes spin { from { transform: rotate(0) } to { transform: rotate(360deg) } }\n.t { animation: spin 1s }"
	html := render(t, Box(CSS(css), Div(Class("t"))))

	if !strings.Contains(html, "@keyframes spin { from { transform: rotate(0) }") {
		t.Fatalf("keyframes veraendert: %s", html)
	}
	if !strings.Contains(html, "] .t") {
		t.Fatalf("regel nach keyframes nicht gescoped: %s", html)
	}
}

func TestBoxScopesInsideMedia(t *testing.T) {
	css := "@media (max-width: 640px) { .hero { font-size: 20px } }"
	html := render(t, Box(CSS(css), Div(Class("hero"))))

	if !strings.Contains(html, "@media (max-width: 640px) {") {
		t.Fatalf("media prelude veraendert: %s", html)
	}
	if !strings.Contains(html, "] .hero{ font-size: 20px }") {
		t.Fatalf("regel in media nicht gescoped: %s", html)
	}
}

func TestBoxSelectorList(t *testing.T) {
	css := ".a, .b:hover { color: red }"
	html := render(t, Box(CSS(css), Div(Class("a"))))

	if !strings.Contains(html, `] .a, `) || !strings.Contains(html, `] .b:hover`) {
		t.Fatalf("selektorliste falsch gescoped: %s", html)
	}
}

func TestSameCssSameScope(t *testing.T) {
	erst := render(t, Box(CSS(".x { color: red }"), Div(Class("x"))))
	zweit := render(t, Box(CSS(".x { color: red }"), Div(Class("x"))))
	if erst != zweit {
		t.Fatalf("gleiches CSS ergibt nicht gleiches Ergebnis")
	}
}

func TestBoxEmitsScript(t *testing.T) {
	html := render(t, Box(
		JS(`root.querySelector("button").addEventListener("click", function () {});`),
		Div(Button(g.Text("Klick"))),
	))

	if !strings.Contains(html, "<script>") {
		t.Fatalf("script tag fehlt: %s", html)
	}
	if !strings.Contains(html, "(function (root) {") {
		t.Fatalf("root-wrapper fehlt: %s", html)
	}
	if !strings.Contains(html, "window.__dreego") {
		t.Fatalf("einmal-guard fehlt: %s", html)
	}
}

func TestBoxEscapesClosingScriptTag(t *testing.T) {
	html := render(t, Box(
		JS(`var böse = "</script><img src=x onerror=alert(1)>";`),
		Div(g.Text("x")),
	))

	if strings.Contains(html, "</script><img") {
		t.Fatalf("closing script tag nicht escaped: %s", html)
	}
	if !strings.Contains(html, `<\/script>`) {
		t.Fatalf("escaped form erwartet: %s", html)
	}
}

func TestBoxScriptCarriesScopeID(t *testing.T) {
	html := render(t, Box(JS(`1;`), Div(g.Text("x"))))
	id := shortHash("\x00" + "1;\n")
	if !strings.Contains(html, `data-scope="`+id+`"`) {
		t.Fatalf("scope id fehlt: %s", html)
	}
	if !strings.Contains(html, `var id = "`+id+`"`) {
		t.Fatalf("script scope id fehlt: %s", html)
	}
}

func TestBoxWithoutCssOrJsIsPlainDiv(t *testing.T) {
	html := render(t, Box(Div(Class("nur-body"))))
	if strings.Contains(html, "<style>") || strings.Contains(html, "<script>") {
		t.Fatalf("kein style/script erwartet: %s", html)
	}
	if !strings.Contains(html, `data-scope=`) {
		t.Fatalf("scope attribut erwartet: %s", html)
	}
}
