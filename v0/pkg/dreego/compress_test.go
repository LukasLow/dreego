package dreego

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http/httptest"
	"strings"
	"testing"
	"testing/fstest"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func gzipApp() *App {
	app := NewApp()
	app.Page(Page{Path: "/", Get: func(c *Ctx) g.Node {
		return h.P(g.Text(strings.Repeat("Hallo Welt. ", 200)))
	}})
	return app
}

func TestGzipKomprimiertHTML(t *testing.T) {
	anfrage := httptest.NewRequest("GET", "/", nil)
	anfrage.Header.Set("Accept-Encoding", "gzip")
	schreiber := httptest.NewRecorder()

	gzipApp().ServeHTTP(schreiber, anfrage)

	antwort := schreiber.Result()
	if antwort.Header.Get("Content-Encoding") != "gzip" {
		t.Fatalf("kein gzip: %q", antwort.Header.Get("Content-Encoding"))
	}
	// Body ist gzip -> entpacken und vergleichen.
	leser, err_gzip := gzip.NewReader(antwort.Body)
	if err_gzip != nil {
		t.Fatalf("gzip lesen: %v", err_gzip)
	}
	entpackt, _ := io.ReadAll(leser)
	if !strings.Contains(string(entpackt), "Hallo Welt.") {
		t.Fatalf("entpackter Inhalt fehlt")
	}
}

func TestOhneGzipHeaderUnkomprimiert(t *testing.T) {
	anfrage := httptest.NewRequest("GET", "/", nil)
	// Kein Accept-Encoding.
	schreiber := httptest.NewRecorder()

	gzipApp().ServeHTTP(schreiber, anfrage)

	if schreiber.Header().Get("Content-Encoding") == "gzip" {
		t.Fatal("ohne Accept-Encoding darf nicht komprimiert werden")
	}
}

func TestGzipQNullVerbietet(t *testing.T) {
	anfrage := httptest.NewRequest("GET", "/", nil)
	anfrage.Header.Set("Accept-Encoding", "gzip;q=0")
	schreiber := httptest.NewRecorder()

	gzipApp().ServeHTTP(schreiber, anfrage)

	if schreiber.Header().Get("Content-Encoding") == "gzip" {
		t.Fatal("gzip;q=0 verbietet Kompression")
	}
}

func TestVaryHeaderImmer(t *testing.T) {
	anfrage := httptest.NewRequest("GET", "/", nil)
	anfrage.Header.Set("Accept-Encoding", "gzip")
	schreiber := httptest.NewRecorder()

	gzipApp().ServeHTTP(schreiber, anfrage)

	if !strings.Contains(schreiber.Header().Get("Vary"), "Accept-Encoding") {
		t.Fatalf("Vary fehlt: %q", schreiber.Header().Get("Vary"))
	}
}

func TestAcceptGzipParsing(t *testing.T) {
	faelle := map[string]bool{
		"gzip":              true,
		"gzip, deflate, br": true,
		"br, gzip;q=0.5":    true,
		"*":                 true,
		"gzip;q=0":          false,
		"deflate, br":       false,
		"":                  false,
	}
	for header, erwartet := range faelle {
		if got := acceptsGzip(header); got != erwartet {
			t.Fatalf("acceptsGzip(%q) = %v, erwartet %v", header, got, erwartet)
		}
	}
}

func TestNichtKomprimierbareInhalte(t *testing.T) {
	app := NewApp()
	app.Static(fstest.MapFS{
		"public/bild.png": &fstest.MapFile{Data: []byte{0x89, 0x50, 0x4E, 0x47, 1, 2, 3}},
	}, "public")
	anfrage := httptest.NewRequest("GET", "/public/bild.png", nil)
	anfrage.Header.Set("Accept-Encoding", "gzip")
	schreiber := httptest.NewRecorder()
	app.ServeHTTP(schreiber, anfrage)

	if schreiber.Header().Get("Content-Encoding") == "gzip" {
		t.Fatal("PNG darf nicht gzip-komprimiert werden")
	}
}

// Regression: eine Antwort, die WriteHeader VOR dem Body schreibt (c.JSON,
// http.Error), muss trotzdem Content-Encoding: gzip tragen. Sonst zeigt der
// Browser Binärmüll.
func TestGzipBeiJSONUndFehler(t *testing.T) {
	app := NewApp()
	app.SetSessionStore(testStore(t))
	app.Page(Page{
		Path: "/api/x",
		API:  func(c *Ctx) error { return c.JSON(200, map[string]string{"a": "b"}) },
	})

	apiAnfrage := httptest.NewRequest("GET", "/api/x", nil)
	apiAnfrage.Header.Set("Accept-Encoding", "gzip")
	apiSchreiber := httptest.NewRecorder()
	app.ServeHTTP(apiSchreiber, apiAnfrage)

	if apiSchreiber.Header().Get("Content-Encoding") != "gzip" {
		t.Fatalf("API: Content-Encoding fehlt, Body waere Müll")
	}
	assertGzipLesbar(t, apiSchreiber)

	fehlerAnfrage := httptest.NewRequest("GET", "/gibtsnicht", nil)
	fehlerAnfrage.Header.Set("Accept-Encoding", "gzip")
	fehlerSchreiber := httptest.NewRecorder()
	app.ServeHTTP(fehlerSchreiber, fehlerAnfrage)

	if fehlerSchreiber.Code != 404 {
		t.Fatalf("404 erwartet, war %d", fehlerSchreiber.Code)
	}
	if fehlerSchreiber.Header().Get("Content-Encoding") != "gzip" {
		t.Fatalf("404: Content-Encoding fehlt")
	}
	assertGzipLesbar(t, fehlerSchreiber)
}

func assertGzipLesbar(t *testing.T, schreiber *httptest.ResponseRecorder) {
	t.Helper()
	leser, err_gzip := gzip.NewReader(bytes.NewReader(schreiber.Body.Bytes()))
	if err_gzip != nil {
		t.Fatalf("body ist kein gültiges gzip: %v", err_gzip)
	}
	entpackt, _ := io.ReadAll(leser)
	if len(entpackt) == 0 {
		t.Fatalf("entpackter Body ist leer")
	}
}
