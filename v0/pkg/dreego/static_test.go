package dreego

import (
	"net/http/httptest"
	"testing"
	"testing/fstest"
)

func TestStaticLiefertDatei(t *testing.T) {
	app := NewApp()
	app.Static(fstest.MapFS{
		"public/design.css": &fstest.MapFile{Data: []byte("body{}")},
	}, "public")

	anfrage := httptest.NewRequest("GET", "/public/design.css", nil)
	schreiber := httptest.NewRecorder()
	app.ServeHTTP(schreiber, anfrage)

	if schreiber.Code != 200 {
		t.Fatalf("status = %d", schreiber.Code)
	}
	if typ := schreiber.Header().Get("Content-Type"); typ != "text/css; charset=utf-8" {
		t.Fatalf("content-type = %q", typ)
	}
	if schreiber.Body.String() != "body{}" {
		t.Fatalf("body = %q", schreiber.Body.String())
	}
}

func TestStaticOhnePraefix(t *testing.T) {
	app := NewApp()
	app.Static(fstest.MapFS{
		"favicon.svg": &fstest.MapFile{Data: []byte("<svg/>")},
	}, "")

	schreiber := httptest.NewRecorder()
	app.ServeHTTP(schreiber, httptest.NewRequest("GET", "/favicon.svg", nil))

	if schreiber.Code != 200 {
		t.Fatalf("status = %d", schreiber.Code)
	}
	if typ := schreiber.Header().Get("Content-Type"); typ != "image/svg+xml" {
		t.Fatalf("content-type = %q", typ)
	}
}

func TestStaticLiefertFontBinaer(t *testing.T) {
	// Beweist, dass der Static-Handler Binärdateien (z.B. Schriften) korrekt
	// ausliefert — der Inhalt geht unverändert durch, nur der Typ wird gesetzt.
	roh := []byte{0x77, 0x4F, 0x46, 0x32, 0x00, 0x01, 0x02, 0x03}
	app := NewApp()
	app.Static(fstest.MapFS{
		"public/fonts/x.woff2": &fstest.MapFile{Data: roh},
	}, "public")

	schreiber := httptest.NewRecorder()
	app.ServeHTTP(schreiber, httptest.NewRequest("GET", "/public/fonts/x.woff2", nil))

	if schreiber.Code != 200 {
		t.Fatalf("status = %d", schreiber.Code)
	}
	if typ := schreiber.Header().Get("Content-Type"); typ != "font/woff2" {
		t.Fatalf("content-type = %q", typ)
	}
	if schreiber.Body.String() != string(roh) {
		t.Fatalf("binärer Inhalt verändert")
	}
}

func TestStaticUnbekanntIst404(t *testing.T) {
	app := NewApp()
	schreiber := httptest.NewRecorder()
	app.ServeHTTP(schreiber, httptest.NewRequest("GET", "/nichts", nil))
	if schreiber.Code != 404 {
		t.Fatalf("status = %d", schreiber.Code)
	}
}
