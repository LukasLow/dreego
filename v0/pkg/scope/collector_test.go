package scope

import (
	"strings"
	"testing"

	. "github.com/LukasLow/dreego/v0/pkg/dom"
	g "github.com/LukasLow/dreego/v0/pkg/dom"
)

func TestCollectorDedupesCss(t *testing.T) {
	c := New()
	seite := Document(c, "de", Head(TitleEl(g.Text("x"))),
		Body(
			c.Box(CSS(".a { color: red }"), Div(Class("a"))),
			c.Box(CSS(".a { color: red }"), Div(Class("a"))), // gleiche Komponente nochmal
		),
	)
	html := render(t, seite)

	id := shortHash(".a { color: red }\n\x00")
	if anzahl := strings.Count(html, `[data-scope="`+id+`"] .a`); anzahl != 1 {
		t.Fatalf("CSS soll genau 1x vorkommen, ist %dx: %s", anzahl, html)
	}
	if anzahl := strings.Count(html, `data-scope="`+id+`"`); anzahl != 3 {
		// 1x im <style> plus 2x im jeweiligen Komponenten-Container.
		t.Fatalf("scope-Attribut erwartet 3x (1 style + 2 container), ist %dx", anzahl)
	}
}

func TestCollectorNonce(t *testing.T) {
	c := New()
	c.SetNonce("abc123")
	seite := Document(c, "de", Head(TitleEl(g.Text("x"))),
		Body(c.Box(
			CSS(".a { color: red }"),
			JS(`root.querySelector(".a");`),
			Div(Class("a")),
		)),
	)
	html := render(t, seite)

	if !strings.Contains(html, `<style nonce="abc123">`) {
		t.Fatalf("style ohne nonce: %s", html)
	}
	if !strings.Contains(html, `<script nonce="abc123">`) {
		t.Fatalf("script ohne nonce: %s", html)
	}
	if strings.Count(html, `nonce="abc123"`) != 2 {
		t.Fatalf("nonce soll 2x vorkommen: %s", html)
	}
}

func TestNewNonceUnique(t *testing.T) {
	erst := NewNonce()
	zweit := NewNonce()
	if erst == "" || zweit == "" {
		t.Fatal("nonce darf nicht leer sein")
	}
	if erst == zweit {
		t.Fatal("nonces sollen unterschiedlich sein")
	}
}

func TestDocumentOrdersHeadStylesScripts(t *testing.T) {
	c := New()
	seite := Document(c, "de", Head(TitleEl(g.Text("Titel"))),
		Body(c.Box(CSS(".a{}"), JS(`1;`), Div(Class("a"), g.Text("HALLOMARKE")))),
	)
	html := render(t, seite)

	if !strings.HasPrefix(html, "<!DOCTYPE html>") {
		t.Fatalf("doctype fehlt: %s", html)
	}
	iTitel := strings.Index(html, "Titel")
	iStyle := strings.Index(html, "<style")
	ihBody := strings.Index(html, "<body>")
	iMarke := strings.Index(html, "HALLOMARKE")
	iScript := strings.Index(html, "<script")

	if !(iTitel < iStyle && iStyle < ihBody && ihBody < iMarke && iMarke < iScript) {
		t.Fatalf("reihenfolge falsch (titel/style/body/marke/script): %d %d %d %d %d\n%s",
			iTitel, iStyle, ihBody, iMarke, iScript, html)
	}
}

func TestCollectorKeepsNonceAcrossRenders(t *testing.T) {
	c := New()
	c.SetNonce("n1")
	if c.Nonce() != "n1" {
		t.Fatalf("nonce nicht gesetzt: %q", c.Nonce())
	}
}

func TestCollectorCriticalInline(t *testing.T) {
	c := New()
	c.SetNonce("abc")
	c.AddCritical(`body { background: #fffdf7 }`)
	c.AddCritical(`body { background: #fffdf7 }`) // Duplikat

	// head-Parameter ist Kopf-INHALT (kein <head>-Element); Document baut <head>.
	seite := Document(c, "de", TitleEl(g.Text("x")), Body(g.Text("inhalt")))
	html := render(t, seite)

	if !strings.Contains(html, "body { background: #fffdf7 }") {
		t.Fatalf("kritisches CSS fehlt: %s", html)
	}
	if anzahl := strings.Count(html, "body { background: #fffdf7 }"); anzahl != 1 {
		t.Fatalf("kritisches CSS soll 1x vorkommen, ist %dx", anzahl)
	}
	iCritical := strings.Index(html, "background: #fffdf7")
	iHeadEnde := strings.Index(html, "</head>")
	if iCritical < 0 || iCritical > iHeadEnde {
		t.Fatalf("kritisches CSS nicht im head: %s", html)
	}
}

func TestCollectorOhneCritical(t *testing.T) {
	c := New()
	seite := Document(c, "de", TitleEl(g.Text("x")), Body(g.Text("x")))
	html := render(t, seite)
	if strings.Contains(html, "color-scheme") {
		t.Fatalf("ohne AddCritical kein Extra-Style: %s", html)
	}
}
