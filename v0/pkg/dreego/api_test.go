package dreego

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	g "maragu.dev/gomponents"
)

func TestAPIRouteJSON(t *testing.T) {
	app := NewApp()
	app.Page(Page{
		Path: "/api/ping",
		API: func(c *Ctx) error {
			return c.JSON(200, map[string]any{"pong": true})
		},
	})

	schreiber := httptest.NewRecorder()
	app.ServeHTTP(schreiber, httptest.NewRequest("GET", "/api/ping", nil))

	if ct := schreiber.Header().Get("Content-Type"); !strings.HasPrefix(ct, "application/json") {
		t.Fatalf("content-type = %q", ct)
	}
	if !strings.Contains(schreiber.Body.String(), `"pong":true`) {
		t.Fatalf("body = %s", schreiber.Body.String())
	}
}

func TestAPIRouteXML(t *testing.T) {
	app := NewApp()
	type antwort struct {
		XMLName struct{} `xml:"ping"`
		Wert    string   `xml:"wert"`
	}
	app.Page(Page{
		Path: "/sitemap.xml",
		API:  func(c *Ctx) error { return c.XML(200, antwort{Wert: "hallo"}) },
	})

	schreiber := httptest.NewRecorder()
	app.ServeHTTP(schreiber, httptest.NewRequest("GET", "/sitemap.xml", nil))

	if ct := schreiber.Header().Get("Content-Type"); !strings.HasPrefix(ct, "application/xml") {
		t.Fatalf("content-type = %q", ct)
	}
	if !strings.Contains(schreiber.Body.String(), "<wert>hallo</wert>") {
		t.Fatalf("body = %s", schreiber.Body.String())
	}
}

func TestAPIRouteKeinLayout(t *testing.T) {
	app := NewApp()
	app.Page(Page{
		Path: "/api/roh",
		API:  func(c *Ctx) error { return c.Write(200, "text/plain", []byte("nur-text")) },
	})

	schreiber := httptest.NewRecorder()
	app.ServeHTTP(schreiber, httptest.NewRequest("GET", "/api/roh", nil))

	body := schreiber.Body.String()
	if strings.Contains(body, "<!DOCTYPE html>") || strings.Contains(body, "<html") {
		t.Fatalf("API-Route darf kein HTML-Layout rendern: %s", body)
	}
	if body != "nur-text" {
		t.Fatalf("body = %q", body)
	}
}

func TestAPIBind(t *testing.T) {
	app := NewApp()
	type eingabe struct {
		Name string `json:"name"`
	}
	var empfangen eingabe
	app.Page(Page{
		Path: "/api/echo",
		API: func(c *Ctx) error {
			if err := c.Bind(&empfangen); err != nil {
				return err
			}
			return c.JSON(200, empfangen)
		},
	})

	anfrage := httptest.NewRequest("GET", "/api/echo", strings.NewReader(`{"name":"lukas"}`))
	schreiber := httptest.NewRecorder()
	app.ServeHTTP(schreiber, anfrage)

	if empfangen.Name != "lukas" {
		t.Fatalf("Bind füllte nicht: %+v", empfangen)
	}
}

func TestAPIMitParam(t *testing.T) {
	app := NewApp()
	app.Page(Page{
		Path: "/api/projekt/[id]",
		API: func(c *Ctx) error {
			return c.JSON(200, map[string]string{"id": c.Param("id")})
		},
	})

	schreiber := httptest.NewRecorder()
	app.ServeHTTP(schreiber, httptest.NewRequest("GET", "/api/projekt/p7", nil))

	if !strings.Contains(schreiber.Body.String(), `"id":"p7"`) {
		t.Fatalf("param in API fehlt: %s", schreiber.Body.String())
	}
}

func TestAPIFehlerIst500(t *testing.T) {
	app := NewApp()
	app.Page(Page{
		Path: "/api/kaputt",
		API:  func(c *Ctx) error { return http.ErrHandlerTimeout },
	})

	schreiber := httptest.NewRecorder()
	app.ServeHTTP(schreiber, httptest.NewRequest("GET", "/api/kaputt", nil))

	if schreiber.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d", schreiber.Code)
	}
}

func TestAPIPostBrauchtCSRF(t *testing.T) {
	app := NewApp() // CSRF on, kein Store
	app.Page(Page{
		Path:    "/api/x",
		Methods: []string{"POST"},
		API:     func(c *Ctx) error { return c.JSON(200, map[string]bool{"ok": true}) },
	})

	anfrage := httptest.NewRequest("POST", "/api/x", strings.NewReader("{}"))
	schreiber := httptest.NewRecorder()
	app.ServeHTTP(schreiber, anfrage)

	if schreiber.Code != http.StatusInternalServerError {
		t.Fatalf("POST ohne Store soll 500 sein, war %d", schreiber.Code)
	}
}

func TestAPIPageRendertNicht(t *testing.T) {
	// Sicherstellen, dass eine normale Seite weiterhin HTML rendert, eine API
	// aber nicht — beide nebeneinander.
	app := NewApp()
	app.Page(Page{Path: "/seite", Get: func(c *Ctx) g.Node { return g.Text("HTML") }})
	app.Page(Page{Path: "/api", API: func(c *Ctx) error { return c.Write(200, "text/plain", []byte("API")) }})

	seite := httptest.NewRecorder()
	app.ServeHTTP(seite, httptest.NewRequest("GET", "/seite", nil))
	if !strings.Contains(seite.Body.String(), "<!DOCTYPE html>") {
		t.Fatalf("normale Seite soll HTML sein: %s", seite.Body.String())
	}

	api := httptest.NewRecorder()
	app.ServeHTTP(api, httptest.NewRequest("GET", "/api", nil))
	if strings.Contains(api.Body.String(), "<html") {
		t.Fatalf("API soll kein HTML sein: %s", api.Body.String())
	}
}
