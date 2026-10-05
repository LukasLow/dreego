package dreego

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	g "github.com/LukasLow/dreego/v0/pkg/dom"
)

func TestParamEinSegment(t *testing.T) {
	app := NewApp()
	app.Page(Page{
		Path: "/projekte/[id]",
		Get:  func(c *Ctx) g.View { return g.Text("id=" + c.Param("id")) },
	})

	schreiber := httptest.NewRecorder()
	app.ServeHTTP(schreiber, httptest.NewRequest("GET", "/projekte/abc-123", nil))

	if !strings.Contains(schreiber.Body.String(), "id=abc-123") {
		t.Fatalf("param nicht gelesen: %s", schreiber.Body.String())
	}
}

func TestParamZweiSegmente(t *testing.T) {
	app := NewApp()
	app.Page(Page{
		Path: "/projekte/[id]/aufgaben/[aufgabe]",
		Get: func(c *Ctx) g.View {
			return g.Text(c.Param("id") + ":" + c.Param("aufgabe"))
		},
	})

	schreiber := httptest.NewRecorder()
	app.ServeHTTP(schreiber, httptest.NewRequest("GET", "/projekte/p1/aufgaben/a9", nil))

	if !strings.Contains(schreiber.Body.String(), "p1:a9") {
		t.Fatalf("mehrere params falsch: %s", schreiber.Body.String())
	}
}

func TestCatchAllRest(t *testing.T) {
	app := NewApp()
	app.Page(Page{
		Path: "/dateien/[...pfad]",
		Get:  func(c *Ctx) g.View { return g.Text("rest=" + c.Param("pfad")) },
	})

	schreiber := httptest.NewRecorder()
	app.ServeHTTP(schreiber, httptest.NewRequest("GET", "/dateien/a/b/c.txt", nil))

	if !strings.Contains(schreiber.Body.String(), "rest=a/b/c.txt") {
		t.Fatalf("catch-all falsch: %s", schreiber.Body.String())
	}
}

func TestFesteUndDynamischeRoutenNeben(t *testing.T) {
	app := NewApp()
	app.Page(Page{Path: "/projekte/neu", Get: func(c *Ctx) g.View { return g.Text("NEU") }})
	app.Page(Page{Path: "/projekte/[id]", Get: func(c *Ctx) g.View { return g.Text("ID=" + c.Param("id")) }})

	// Fester Pfad gewinnt (Go-1.22-Mux bevorzugt den spezifischeren Treffer).
	fest := httptest.NewRecorder()
	app.ServeHTTP(fest, httptest.NewRequest("GET", "/projekte/neu", nil))
	if !strings.Contains(fest.Body.String(), "NEU") {
		t.Fatalf("fester Pfad: %s", fest.Body.String())
	}

	dyn := httptest.NewRecorder()
	app.ServeHTTP(dyn, httptest.NewRequest("GET", "/projekte/x1", nil))
	if !strings.Contains(dyn.Body.String(), "ID=x1") {
		t.Fatalf("dynamischer Pfad: %s", dyn.Body.String())
	}
}

func TestUnbekannterPfad404(t *testing.T) {
	app := NewApp()
	app.Page(Page{Path: "/a/[id]", Get: func(c *Ctx) g.View { return g.Text("x") }})

	schreiber := httptest.NewRecorder()
	app.ServeHTTP(schreiber, httptest.NewRequest("GET", "/voellig/anders", nil))

	if schreiber.Code != http.StatusNotFound {
		t.Fatalf("status = %d", schreiber.Code)
	}
}

func TestToServeMuxPattern(t *testing.T) {
	faelle := map[string]string{
		"/projekte/[id]":     "/projekte/{id}",
		"/a/[x]/b/[y]":       "/a/{x}/b/{y}",
		"/dateien/[...pfad]": "/dateien/{pfad...}",
		"/fest":              "/fest",
		"/projekte/[id]/":    "/projekte/{id}/",
	}
	for eingabe, erwartet := range faelle {
		if got := toServeMuxPattern(eingabe); got != erwartet {
			t.Fatalf("toServeMuxPattern(%q) = %q, erwartet %q", eingabe, got, erwartet)
		}
	}
}
