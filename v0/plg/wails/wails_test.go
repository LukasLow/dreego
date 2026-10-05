package wails

import (
	"net/http/httptest"
	"strings"
	"testing"

	d "github.com/LukasLow/dreego/v0/pkg/dom"
	"github.com/LukasLow/dreego/v0/pkg/dreego"
)

func testApp() *dreego.App {
	app := dreego.NewApp()
	app.Page(dreego.Page{
		Path: "/",
		Get:  func(c *dreego.Ctx) d.View { return d.H1(d.Text("Hallo Wails")) },
	})
	return app
}

func TestHandlerLiefertWailsCSP(t *testing.T) {
	app := testApp()
	handler := Handler(app)

	schreiber := httptest.NewRecorder()
	handler.ServeHTTP(schreiber, httptest.NewRequest("GET", "/", nil))

	csp := schreiber.Header().Get("Content-Security-Policy")
	if !strings.Contains(csp, "'unsafe-inline'") {
		t.Fatalf("Wails-CSP fehlt (unsafe-inline): %q", csp)
	}
	if !strings.Contains(csp, "frame-ancestors 'none'") {
		t.Fatalf("CSP soll weiterhin gehärtet sein: %q", csp)
	}
}

func TestHandlerRespektiertEigeneCSP(t *testing.T) {
	app := testApp()
	app.SetCSP("default-src 'none'")
	handler := Handler(app)

	schreiber := httptest.NewRecorder()
	handler.ServeHTTP(schreiber, httptest.NewRequest("GET", "/", nil))

	if csp := schreiber.Header().Get("Content-Security-Policy"); csp != "default-src 'none'" {
		t.Fatalf("eigene CSP wurde überschrieben: %q", csp)
	}
}

func TestHandlerRendertSeite(t *testing.T) {
	handler := Handler(testApp())
	schreiber := httptest.NewRecorder()
	handler.ServeHTTP(schreiber, httptest.NewRequest("GET", "/", nil))

	if schreiber.Code != 200 {
		t.Fatalf("status = %d", schreiber.Code)
	}
	if !strings.Contains(schreiber.Body.String(), "Hallo Wails") {
		t.Fatalf("body: %s", schreiber.Body.String())
	}
}
