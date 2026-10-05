package markdown

import (
	"bytes"
	"strings"
	"testing"

	g "maragu.dev/gomponents"
)

func render(t *testing.T, knoten []g.Node) string {
	t.Helper()
	var puffer bytes.Buffer
	err := g.Group(knoten).Render(&puffer)
	if err != nil {
		t.Fatalf("render: %v", err)
	}
	return puffer.String()
}

func TestUeberschriften(t *testing.T) {
	html := render(t, ToNodes("# Titel\n\n## Untertitel"))
	if !strings.Contains(html, "<h1>Titel</h1>") {
		t.Fatalf("h1 fehlt: %s", html)
	}
	if !strings.Contains(html, "<h2>Untertitel</h2>") {
		t.Fatalf("h2 fehlt: %s", html)
	}
}

func TestAbsatz(t *testing.T) {
	html := render(t, ToNodes("Hallo Welt"))
	if !strings.Contains(html, "<p>Hallo Welt</p>") {
		t.Fatalf("absatz fehlt: %s", html)
	}
}

func TestListe(t *testing.T) {
	html := render(t, ToNodes("- eins\n- zwei\n- drei"))
	if !strings.Contains(html, "<ul><li>eins</li><li>zwei</li><li>drei</li></ul>") {
		t.Fatalf("liste falsch: %s", html)
	}
}

func TestInlineFettKursivCode(t *testing.T) {
	html := render(t, ToNodes("Das ist **fett**, *kursiv* und `code`."))
	if !strings.Contains(html, "<strong>fett</strong>") {
		t.Fatalf("fett fehlt: %s", html)
	}
	if !strings.Contains(html, "<em>kursiv</em>") {
		t.Fatalf("kursiv fehlt: %s", html)
	}
	if !strings.Contains(html, "<code>code</code>") {
		t.Fatalf("code fehlt: %s", html)
	}
}

func TestLink(t *testing.T) {
	html := render(t, ToNodes("[Datenschutz](/legal/datenschutz)"))
	if !strings.Contains(html, `<a href="/legal/datenschutz">Datenschutz</a>`) {
		t.Fatalf("link fehlt: %s", html)
	}
}

func TestJavaScriptLinkWirdVerworfen(t *testing.T) {
	html := render(t, ToNodes("[klick](javascript:alert(1))"))
	if strings.Contains(html, "javascript:") {
		t.Fatalf("javascript: durchgelassen: %s", html)
	}
	if !strings.Contains(html, "klick") {
		t.Fatalf("label soll als Text bleiben: %s", html)
	}
}

func TestEscapingKeinRohHTML(t *testing.T) {
	html := render(t, ToNodes("<script>alert(1)</script>"))
	if strings.Contains(html, "<script>") {
		t.Fatalf("roh-HTML durchgelassen: %s", html)
	}
	if !strings.Contains(html, "&lt;script&gt;") {
		t.Fatalf("escaping fehlt: %s", html)
	}
}

func TestZitatUndLinie(t *testing.T) {
	html := render(t, ToNodes("> Zitat\n\n---"))
	if !strings.Contains(html, "<blockquote>Zitat</blockquote>") {
		t.Fatalf("zitat fehlt: %s", html)
	}
	if !strings.Contains(html, "<hr") {
		t.Fatalf("linie fehlt: %s", html)
	}
}
