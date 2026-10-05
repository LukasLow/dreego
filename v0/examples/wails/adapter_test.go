package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	webapp "wails-dreego/app/web"

	"github.com/LukasLow/dreego/v0/pkg/dreego"
	"github.com/LukasLow/dreego/v0/plg/wails"
)

// TestAdapterServesPages proves the dreego pages render through the
// listener-free Wails adapter: literal GET routes only, no localhost server.
func TestAdapterServesPages(t *testing.T) {
	app := dreego.NewApp()
	webapp.Register(app)
	handler := wails.Handler(app)

	for _, tc := range []struct{ path, want string }{
		{"/", "dreego frontend"},
		{"/about", "Clear ownership"},
	} {
		schreiber := httptest.NewRecorder()
		handler.ServeHTTP(schreiber, httptest.NewRequest(http.MethodGet, tc.path, nil))
		if schreiber.Code != http.StatusOK {
			t.Fatalf("GET %s status = %d, want 200", tc.path, schreiber.Code)
		}
		if !strings.Contains(schreiber.Body.String(), tc.want) {
			t.Fatalf("GET %s body missing %q", tc.path, tc.want)
		}
	}
}

// TestCSPAllowsWailsRuntime guards the reason the adapter exists: the CSP must
// permit the injected Wails runtime (inline scripts).
func TestCSPAllowsWailsRuntime(t *testing.T) {
	app := dreego.NewApp()
	webapp.Register(app)

	schreiber := httptest.NewRecorder()
	wails.Handler(app).ServeHTTP(schreiber, httptest.NewRequest(http.MethodGet, "/", nil))

	csp := schreiber.Header().Get("Content-Security-Policy")
	if !strings.Contains(csp, "'unsafe-inline'") {
		t.Fatalf("CSP must allow the Wails runtime (inline): %q", csp)
	}
}

// TestUnknownPathIs404 keeps the adapter honest about its boundary.
func TestUnknownPathIs404(t *testing.T) {
	app := dreego.NewApp()
	webapp.Register(app)

	schreiber := httptest.NewRecorder()
	wails.Handler(app).ServeHTTP(schreiber, httptest.NewRequest(http.MethodGet, "/nope", nil))
	if schreiber.Code != http.StatusNotFound {
		t.Fatalf("GET /nope status = %d, want 404", schreiber.Code)
	}
}
