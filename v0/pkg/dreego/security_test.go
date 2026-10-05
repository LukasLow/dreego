package dreego

import (
	"strings"
	"testing"

	g "github.com/LukasLow/dreego/v0/pkg/dom"
)

func cspApp() *App {
	app := NewApp()
	app.Page(Page{Path: "/", Get: func(c *Ctx) g.View { return g.Text("x") }})
	return app
}

func getCSP(t *testing.T, app *App, pfad string) string {
	t.Helper()
	schreiber := aufruf(app, "GET", pfad, "")
	return schreiber.Header().Get("Content-Security-Policy")
}

func TestCSPDefaultStreng(t *testing.T) {
	csp := getCSP(t, cspApp(), "/")
	if !strings.Contains(csp, "'nonce-") {
		t.Fatalf("default soll nonce enthalten: %q", csp)
	}
	if !strings.Contains(csp, "frame-ancestors 'none'") {
		t.Fatalf("default soll gehärtet sein: %q", csp)
	}
}

func TestCSPAppWeit(t *testing.T) {
	app := cspApp()
	app.SetCSP("default-src 'none'; script-src 'self' 'nonce-{nonce}'")
	csp := getCSP(t, app, "/")
	if strings.Contains(csp, "frame-ancestors") {
		t.Fatalf("app-CSP soll den default ersetzen: %q", csp)
	}
	if !strings.Contains(csp, "'nonce-") {
		t.Fatalf("nonce soll ersetzt sein: %q", csp)
	}
}

func TestCSPProSeiteUeberschreibtApp(t *testing.T) {
	app := NewApp()
	app.SetCSP("default-src 'self'")
	app.Page(Page{
		Path:     "/locker",
		Get:      func(c *Ctx) g.View { return g.Text("x") },
		Security: &Security{CSP: "default-src 'self' https://cdn.example"},
	})
	app.Page(Page{Path: "/", Get: func(c *Ctx) g.View { return g.Text("x") }})

	locker := getCSP(t, app, "/locker")
	if !strings.Contains(locker, "cdn.example") {
		t.Fatalf("seiten-CSP fehlt: %q", locker)
	}
	start := getCSP(t, app, "/")
	if strings.Contains(start, "cdn.example") {
		t.Fatalf("app-CSP soll auf / gelten: %q", start)
	}
}

func TestCSPAbschaltbar(t *testing.T) {
	app := cspApp()
	app.SetCSP(CSPOff)
	if csp := getCSP(t, app, "/"); csp != "" {
		t.Fatalf("CSP soll aus sein, war %q", csp)
	}
}

func TestCSPNonceWirdEingesetzt(t *testing.T) {
	app := cspApp()
	csp := getCSP(t, app, "/")
	if strings.Contains(csp, "{nonce}") {
		t.Fatalf("Platzhalter nicht ersetzt: %q", csp)
	}
	if !strings.Contains(csp, "'nonce-") {
		t.Fatalf("nonce fehlt: %q", csp)
	}
}
