package dreego

import (
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	g "github.com/LukasLow/dreego/v0/pkg/dom"
)

// beweist: app.Use bindet eigenes Logging ein, das NICHT in dreego liegt.
func TestUseBindetEigenesLoggingEin(t *testing.T) {
	var geloggt []string

	logging := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			geloggt = append(geloggt, r.Method+" "+r.URL.Path)
			next.ServeHTTP(w, r)
		})
	}

	app := NewApp()
	app.Use(logging)
	app.Page(Page{Path: "/x", Get: func(c *Ctx) g.View { return g.Text("ok") }})

	schreiber := httptest.NewRecorder()
	app.ServeHTTP(schreiber, httptest.NewRequest("GET", "/x", nil))

	if len(geloggt) != 1 || geloggt[0] != "GET /x" {
		t.Fatalf("eigenes Logging lief nicht: %v", geloggt)
	}
}

// beweist: roher net/http-Redirect über c.Response()/c.Request() funktioniert.
func TestRoherHTTPRedirect(t *testing.T) {
	app := NewApp()
	app.Page(Page{Path: "/alt", Get: func(c *Ctx) g.View {
		http.Redirect(c.Response(), c.Request(), "/neu", http.StatusMovedPermanently)
		return g.Text("") // nie gerendert
	}})

	schreiber := httptest.NewRecorder()
	app.ServeHTTP(schreiber, httptest.NewRequest("GET", "/alt", nil))

	if schreiber.Code != http.StatusMovedPermanently {
		t.Fatalf("status = %d", schreiber.Code)
	}
	if ziel := schreiber.Header().Get("Location"); ziel != "/neu" {
		t.Fatalf("Location = %q", ziel)
	}
}

// beweist: dreego lässt sich als Handler in einen eigenen http.ServeMux hängen,
// neben eigenen Routen (z. B. einer JSON-API).
func TestMountInEigenemMux(t *testing.T) {
	app := NewApp()
	app.Page(Page{Path: "/seite", Get: func(c *Ctx) g.View { return g.Text("dreego-seite") }})

	mux := http.NewServeMux()
	mux.HandleFunc("/api/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"pong":true}`))
	})
	mux.Handle("/", app.Handler()) // dreego als Unterbaum

	// Eigene API-Route:
	api := httptest.NewRecorder()
	mux.ServeHTTP(api, httptest.NewRequest("GET", "/api/ping", nil))
	if !strings.Contains(api.Body.String(), "pong") {
		t.Fatalf("eigene API-Route: %s", api.Body.String())
	}

	// dreego-Seite daneben:
	seite := httptest.NewRecorder()
	mux.ServeHTTP(seite, httptest.NewRequest("GET", "/seite", nil))
	if !strings.Contains(seite.Body.String(), "dreego-seite") {
		t.Fatalf("dreego-Seite im Mux: %s", seite.Body.String())
	}

	log.Printf("mux-mount ok")
}
