package dreego

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	g "github.com/LukasLow/dreego/v0/pkg/dom"
)

func csrfApp(t *testing.T) *App {
	t.Helper()
	app := NewApp()
	app.SetSessionStore(testStore(t))
	app.Page(Page{
		Path: "/form",
		Get: func(c *Ctx) g.View {
			return g.Form(c.CSRFInput(), g.Button(g.Text("ok")))
		},
		Post: func(c *Ctx) g.View { return g.Text("durchgelassen") },
	})
	return app
}

func holeCSRFCookie(t *testing.T, app *App) *http.Cookie {
	t.Helper()
	schreiber := httptest.NewRecorder()
	app.ServeHTTP(schreiber, httptest.NewRequest("GET", "/form", nil))
	cookies := schreiber.Result().Cookies()
	if len(cookies) == 0 {
		t.Fatal("kein Session-Cookie beim GET")
	}
	return cookies[0]
}

func TestCSRFBlockiertPostOhneToken(t *testing.T) {
	app := csrfApp(t)
	anfrage := httptest.NewRequest("POST", "/form", strings.NewReader(""))
	anfrage.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	anfrage.AddCookie(holeCSRFCookie(t, app))

	schreiber := httptest.NewRecorder()
	app.ServeHTTP(schreiber, anfrage)

	if schreiber.Code != http.StatusForbidden {
		t.Fatalf("status = %d, erwartet 403", schreiber.Code)
	}
}

func TestCSRFFormularHatToken(t *testing.T) {
	app := csrfApp(t)
	schreiber := httptest.NewRecorder()
	app.ServeHTTP(schreiber, httptest.NewRequest("GET", "/form", nil))

	if !strings.Contains(schreiber.Body.String(), `name="_csrf"`) {
		t.Fatalf("kein CSRF-Feld: %s", schreiber.Body.String())
	}
}

func TestCSRFLaesstGueltigenTokenDurch(t *testing.T) {
	app := csrfApp(t)
	cookie := holeCSRFCookie(t, app)

	// Token aus der Session lesen.
	sessionAnfrage := httptest.NewRequest("GET", "/form", nil)
	sessionAnfrage.AddCookie(cookie)
	token, _ := app.store.Get(sessionAnfrage, csrfSessionKey)
	if token == "" {
		t.Fatal("kein Token in der Session")
	}

	form := url.Values{csrfFieldName: {token}}
	anfrage := httptest.NewRequest("POST", "/form", strings.NewReader(form.Encode()))
	anfrage.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	anfrage.AddCookie(cookie)

	schreiber := httptest.NewRecorder()
	app.ServeHTTP(schreiber, anfrage)

	if schreiber.Code != http.StatusOK {
		t.Fatalf("status = %d, erwartet 200", schreiber.Code)
	}
	if !strings.Contains(schreiber.Body.String(), "durchgelassen") {
		t.Fatalf("body: %s", schreiber.Body.String())
	}
}

func TestCSRFAbgeschaltet(t *testing.T) {
	app := NewApp()
	app.SetCSRF(false)
	app.Page(Page{Path: "/frei", Post: func(c *Ctx) g.View { return g.Text("ok") }})

	anfrage := httptest.NewRequest("POST", "/frei", strings.NewReader(""))
	anfrage.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	schreiber := httptest.NewRecorder()
	app.ServeHTTP(schreiber, anfrage)

	if schreiber.Code != http.StatusOK {
		t.Fatalf("status = %d", schreiber.Code)
	}
}
