package dreego

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func TestRedirectKlemmtUngueltigenCode(t *testing.T) {
	app := NewApp()
	app.Page(Page{Path: "/r", Get: func(c *Ctx) g.Node { return Redirect("/ziel", 0) }})

	schreiber := httptest.NewRecorder()
	app.ServeHTTP(schreiber, httptest.NewRequest("GET", "/r", nil))

	if schreiber.Code != http.StatusSeeOther {
		t.Fatalf("status = %d, erwartet 303", schreiber.Code)
	}
}

func TestCSRFohneStoreFehlerNicht403(t *testing.T) {
	// CSRF an, aber kein Store: unsicherer Request muss als Serverfehler
	// enden (laut), nicht als 403 (irreführend) — und nicht durchgelassen.
	app := NewApp()
	app.Page(Page{Path: "/p", Post: func(c *Ctx) g.Node { return g.Text("ok") }})

	anfrage := httptest.NewRequest("POST", "/p", strings.NewReader(""))
	anfrage.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	schreiber := httptest.NewRecorder()
	app.ServeHTTP(schreiber, anfrage)

	if schreiber.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, erwartet 500", schreiber.Code)
	}
}

func TestCSRFausOhneStore(t *testing.T) {
	app := NewApp()
	app.SetCSRF(false)
	app.Page(Page{Path: "/p", Post: func(c *Ctx) g.Node { return g.Text("ok") }})

	anfrage := httptest.NewRequest("POST", "/p", strings.NewReader(""))
	anfrage.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	schreiber := httptest.NewRecorder()
	app.ServeHTTP(schreiber, anfrage)

	if schreiber.Code != http.StatusOK {
		t.Fatalf("status = %d, erwartet 200", schreiber.Code)
	}
}

func TestGrosseSessionMeldetFehler(t *testing.T) {
	app := NewApp()
	app.SetSessionStore(testStore(t))
	app.Page(Page{Path: "/gross", Get: func(c *Ctx) g.Node {
		c.SetSessionVal("riesig", strings.Repeat("x", 8000))
		return h.P(g.Text("ok"))
	}})

	schreiber := httptest.NewRecorder()
	app.ServeHTTP(schreiber, httptest.NewRequest("GET", "/gross", nil))

	if schreiber.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, erwartet 500 (Session zu groß)", schreiber.Code)
	}
}

func TestFormActionUndCSP(t *testing.T) {
	app := NewApp()
	app.Page(Page{Path: "/", Get: func(c *Ctx) g.Node { return g.Text("x") }})

	schreiber := httptest.NewRecorder()
	app.ServeHTTP(schreiber, httptest.NewRequest("GET", "/", nil))

	csp := schreiber.Header().Get("Content-Security-Policy")
	for _, teil := range []string{"base-uri", "form-action", "frame-ancestors", "object-src", "img-src", "'self'"} {
		if !strings.Contains(csp, teil) {
			t.Fatalf("CSP ohne %s: %q", teil, csp)
		}
	}
}

func TestCSRFTokenHeaderAlternative(t *testing.T) {
	app := csrfApp(t)
	cookie := holeCSRFCookie(t, app)

	sessionAnfrage := httptest.NewRequest("GET", "/form", nil)
	sessionAnfrage.AddCookie(cookie)
	token, _ := app.store.Get(sessionAnfrage, csrfSessionKey)

	anfrage := httptest.NewRequest("POST", "/form", strings.NewReader(url.Values{}.Encode()))
	anfrage.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	anfrage.Header.Set("X-CSRF-Token", token)
	anfrage.AddCookie(cookie)

	schreiber := httptest.NewRecorder()
	app.ServeHTTP(schreiber, anfrage)
	if schreiber.Code != http.StatusOK {
		t.Fatalf("status = %d, erwartet 200", schreiber.Code)
	}
}
