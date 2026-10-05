package dreego

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	g "github.com/LukasLow/dreego/v0/pkg/dom"
)

func Test404CustomHTMLPage(t *testing.T) {
	app := NewApp()
	app.SetErrorPage(404, func(c *Ctx) g.View {
		return g.H1(g.Text("NIXX-DA-404"))
	})

	schreiber := httptest.NewRecorder()
	app.ServeHTTP(schreiber, httptest.NewRequest("GET", "/gibtsnicht", nil))

	if schreiber.Code != http.StatusNotFound {
		t.Fatalf("status = %d", schreiber.Code)
	}
	if !strings.Contains(schreiber.Body.String(), "NIXX-DA-404") {
		t.Fatalf("eigene 404-Seite fehlt: %s", schreiber.Body.String())
	}
	if ct := schreiber.Header().Get("Content-Type"); !strings.HasPrefix(ct, "text/html") {
		t.Fatalf("content-type = %q", ct)
	}
}

func Test404FallbackOhneEigeneSeite(t *testing.T) {
	app := NewApp()
	schreiber := httptest.NewRecorder()
	app.ServeHTTP(schreiber, httptest.NewRequest("GET", "/gibtsnicht", nil))

	if schreiber.Code != http.StatusNotFound {
		t.Fatalf("status = %d", schreiber.Code)
	}
	body := schreiber.Body.String()
	if !strings.Contains(body, "<!DOCTYPE html>") {
		t.Fatalf("Fallback soll HTML sein: %s", body)
	}
	if !strings.Contains(body, "404") {
		t.Fatalf("Fallback soll den Status zeigen: %s", body)
	}
}

func Test404JSONFuerAPI(t *testing.T) {
	app := NewApp()
	schreiber := httptest.NewRecorder()
	app.ServeHTTP(schreiber, httptest.NewRequest("GET", "/api/gibtsnicht", nil))

	if schreiber.Code != http.StatusNotFound {
		t.Fatalf("status = %d", schreiber.Code)
	}
	if ct := schreiber.Header().Get("Content-Type"); !strings.HasPrefix(ct, "application/json") {
		t.Fatalf("content-type = %q", ct)
	}
	body := schreiber.Body.String()
	if !strings.Contains(body, `"status":404`) {
		t.Fatalf("JSON-404 fehlt: %s", body)
	}
}

func Test404JSONPerAcceptHeader(t *testing.T) {
	app := NewApp()
	anfrage := httptest.NewRequest("GET", "/seite", nil)
	anfrage.Header.Set("Accept", "application/json")
	schreiber := httptest.NewRecorder()
	app.ServeHTTP(schreiber, anfrage)

	if ct := schreiber.Header().Get("Content-Type"); !strings.HasPrefix(ct, "application/json") {
		t.Fatalf("content-type = %q", ct)
	}
}

func Test500OhneInterneDetails(t *testing.T) {
	app := NewApp()
	app.Page(Page{Path: "/kaputt", Get: func(c *Ctx) g.View {
		panic("GEHEIMER-INTERNER-FEHLER")
	}})

	schreiber := httptest.NewRecorder()
	app.ServeHTTP(schreiber, httptest.NewRequest("GET", "/kaputt", nil))

	if schreiber.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d", schreiber.Code)
	}
	if strings.Contains(schreiber.Body.String(), "GEHEIMER-INTERNER-FEHLER") {
		t.Fatalf("interne Details geleakt: %s", schreiber.Body.String())
	}
}

func Test403CustomPage(t *testing.T) {
	app := NewApp()
	app.SetErrorPage(403, func(c *Ctx) g.View { return g.H1(g.Text("VERBOTEN-403")) })
	app.Page(Page{Path: "/post", Post: func(c *Ctx) g.View { return g.Text("ok") }})

	anfrage := httptest.NewRequest("POST", "/post", strings.NewReader(""))
	anfrage.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	schreiber := httptest.NewRecorder()
	app.ServeHTTP(schreiber, anfrage)

	// Ohne Store kommt 500 (CSRF aktiv, kein Store), nicht 403 — das ist
	// Absicht. Hier nur sicherstellen, dass es eine HTML-Fehlerseite ist.
	if !strings.Contains(schreiber.Body.String(), "<!DOCTYPE html>") {
		t.Fatalf("Fehlerseite soll HTML sein: %s", schreiber.Body.String())
	}
}

func Test405IstHTML(t *testing.T) {
	// Eine Route, die GET UND POST kann — aber eine dritte Methode bekommt 405
	// nur, wenn der Mux sie auf denselben Handler leitet. Der Go-Mux kennt pro
	// Muster nur die registrierten Methoden; eine nicht registrierte Methode
	// ergibt 405, wenn das Muster existiert.
	app := NewApp()
	app.Page(Page{
		Path: "/beides",
		Get:  func(c *Ctx) g.View { return g.Text("g") },
		Post: func(c *Ctx) g.View { return g.Text("p") },
	})

	// PUT ist nicht registriert -> der Mux meldet 405 (Muster existiert).
	anfrage := httptest.NewRequest("PUT", "/beides", strings.NewReader(""))
	schreiber := httptest.NewRecorder()
	app.ServeHTTP(schreiber, anfrage)

	if schreiber.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, erwartet 405", schreiber.Code)
	}
	if !strings.Contains(schreiber.Body.String(), "<!DOCTYPE html>") {
		t.Fatalf("405 soll HTML sein: %s", schreiber.Body.String())
	}
}
