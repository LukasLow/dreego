package dreego

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	g "github.com/LukasLow/dreego/v0/pkg/dom"
)

func aufruf(app *App, methode string, pfad string, form string) *httptest.ResponseRecorder {
	var leib io.Reader
	if form != "" {
		leib = strings.NewReader(form)
	}
	anfrage := httptest.NewRequest(methode, pfad, leib)
	if form != "" {
		anfrage.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	schreiber := httptest.NewRecorder()
	app.ServeHTTP(schreiber, anfrage)
	return schreiber
}

func TestPageRouting(t *testing.T) {
	app := NewApp()
	app.Page(Page{
		Path: "/hallo",
		Get:  func(c *Ctx) g.View { return g.H1(g.Text("Hi")) },
	})

	antwort := aufruf(app, "GET", "/hallo", "")
	if antwort.Code != http.StatusOK {
		t.Fatalf("status = %d", antwort.Code)
	}
	if !strings.Contains(antwort.Body.String(), "Hi") {
		t.Fatalf("body: %s", antwort.Body.String())
	}
}

func TestUnbekannteRouteIst404(t *testing.T) {
	app := NewApp()
	antwort := aufruf(app, "GET", "/gibts-nicht", "")
	if antwort.Code != http.StatusNotFound {
		t.Fatalf("status = %d, erwartet 404", antwort.Code)
	}
}

func TestFalscheMethodeIst405(t *testing.T) {
	app := NewApp()
	app.Page(Page{Path: "/x", Get: func(c *Ctx) g.View { return g.Text("x") }})

	antwort := aufruf(app, "POST", "/x", "")
	if antwort.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, erwartet 405", antwort.Code)
	}
	if allow := antwort.Header().Get("Allow"); !strings.Contains(allow, "GET") {
		t.Fatalf("Allow = %q", allow)
	}
}

func TestRedirectIstNode(t *testing.T) {
	app := NewApp()
	app.Page(Page{Path: "/alt", Get: func(c *Ctx) g.View {
		return Redirect("/neu", 303)
	}})

	antwort := aufruf(app, "GET", "/alt", "")
	if antwort.Code != http.StatusSeeOther {
		t.Fatalf("status = %d, erwartet 303", antwort.Code)
	}
	if ziel := antwort.Header().Get("Location"); ziel != "/neu" {
		t.Fatalf("Location = %q", ziel)
	}
	if antwort.Body.Len() != 0 {
		t.Fatalf("redirect soll keinen Body haben: %q", antwort.Body.String())
	}
}

func TestCSPNonceHeader(t *testing.T) {
	app := NewApp()
	app.Page(Page{Path: "/", Get: func(c *Ctx) g.View { return g.Text("x") }})

	antwort := aufruf(app, "GET", "/", "")
	csp := antwort.Header().Get("Content-Security-Policy")
	if !strings.Contains(csp, "'nonce-") {
		t.Fatalf("kein nonce in CSP: %q", csp)
	}
	if !strings.Contains(csp, "script-src") {
		t.Fatalf("keine script-src in CSP: %q", csp)
	}
}

func TestCtxReichtRequestUndNonce(t *testing.T) {
	app := NewApp()
	var gesehenNonce string
	var gesehenPfad string
	app.Page(Page{Path: "/c", Get: func(c *Ctx) g.View {
		gesehenNonce = c.Nonce()
		gesehenPfad = c.Request().URL.Path
		return g.Text("ok")
	}})

	aufruf(app, "GET", "/c", "")
	csp := ""
	_ = csp
	if gesehenNonce == "" {
		t.Fatal("nonce im ctx leer")
	}
	if gesehenPfad != "/c" {
		t.Fatalf("req-pfad = %q", gesehenPfad)
	}
}

func TestMethodenAuswahl(t *testing.T) {
	nurGet := Page{Path: "/g", Get: func(c *Ctx) g.View { return g.Text("g") }}
	if nurGet.handlerFor("POST") != nil {
		t.Fatal("POST soll nil sein")
	}
	if nurGet.handlerFor("GET") == nil {
		t.Fatal("GET soll da sein")
	}

	beide := Page{Path: "/b", Get: func(c *Ctx) g.View { return g.Text("g") }, Post: func(c *Ctx) g.View { return g.Text("p") }}
	if beide.handlerFor("POST") == nil {
		t.Fatal("POST soll da sein")
	}
	if erlaubt := beide.allowedMethods(); !strings.Contains(erlaubt, "POST") {
		t.Fatalf("allowedMethods = %q", erlaubt)
	}
}
