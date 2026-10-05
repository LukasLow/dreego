package dreego

import (
	"strings"
	"testing"

	g "github.com/LukasLow/dreego/v0/pkg/dom"
	"github.com/LukasLow/dreego/v0/pkg/scope"
)

func renderSeite(t *testing.T, seite g.View) string {
	t.Helper()
	var bau strings.Builder
	err_render := seite.Render(&bau)
	if err_render != nil {
		t.Fatalf("render: %v", err_render)
	}
	return bau.String()
}

func TestStyleErzeugtNonceKlasse(t *testing.T) {
	collector := scope.New()
	collector.SetNonce("abc123")
	ctx := &Ctx{collector: collector}

	klasse := ctx.Style("width: 42%")
	if klasse == "" {
		t.Fatal("leerer Klassenname")
	}

	seite := scope.Document(collector, "de", g.Head(g.TitleEl(g.Text("x"))),
		g.Body(g.Div(g.Class(klasse))))

	html := renderSeite(t, seite)

	if !strings.Contains(html, "."+klasse+" {width: 42%}") {
		t.Fatalf("Regel fehlt: %s", html)
	}
	if !strings.Contains(html, `<style nonce="abc123">`) {
		t.Fatalf("Page-CSS ohne nonce: %s", html)
	}
	if !strings.Contains(html, `class="`+klasse+`"`) {
		t.Fatalf("Klasse nicht am Element: %s", html)
	}
}

func TestStyleLeerGibtLeer(t *testing.T) {
	ctx := &Ctx{collector: scope.New()}
	if ctx.Style("   ") != "" {
		t.Fatal("leere Deklaration soll \"\" ergeben")
	}
}

func TestStyleDedupliziert(t *testing.T) {
	ctx := &Ctx{collector: scope.New()}
	erste := ctx.Style("width: 42%")
	zweite := ctx.Style("width: 42%")
	if erste != zweite {
		t.Fatalf("gleiche Deklaration, verschiedene Namen: %q / %q", erste, zweite)
	}

	html := renderSeite(t, scope.Document(ctx.collector, "de", g.Head(g.TitleEl(g.Text("x"))), g.Body(g.Text("x"))))
	if anzahl := strings.Count(html, "width: 42%"); anzahl != 1 {
		t.Fatalf("Regel soll 1x vorkommen, ist %dx", anzahl)
	}
}
