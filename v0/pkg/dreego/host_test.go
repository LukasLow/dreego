package dreego

import (
	"io"
	"net/http"
	"testing"

	g "github.com/LukasLow/dreego/v0/pkg/dom"
)

func TestHostStartServeClose(t *testing.T) {
	app := NewApp()
	app.Page(Page{Path: "/", Get: func(c *Ctx) g.View { return g.Text("hallo") }})

	server := Start("127.0.0.1:0", app.Handler())
	if server.Addr() == "" {
		t.Fatal("kein Listener")
	}

	antwort, err_get := http.Get("http://" + server.Addr() + "/")
	if err_get != nil {
		t.Fatalf("get: %v", err_get)
	}

	defer func() { _ = antwort.Body.Close() }()

	leib, err_lesen := io.ReadAll(antwort.Body)
	if err_lesen != nil {
		t.Fatalf("lesen: %v", err_lesen)
	}

	if antwort.StatusCode != http.StatusOK {
		t.Fatalf("status = %d", antwort.StatusCode)
	}
	if string(leib) == "" {
		t.Fatal("leerer Body")
	}

	err_close := server.Close()
	if err_close != nil {
		t.Fatalf("close: %v", err_close)
	}
	if err_wait := server.Wait(); err_wait != nil {
		t.Fatalf("wait nach Close = %v, erwartet nil", err_wait)
	}
}

func TestHostBindFehlerKommtAusWait(t *testing.T) {
	// Zweimal derselbe Port: der zweite Start scheitert — Wait meldet es,
	// statt den Prozess hart zu beenden.
	erste := Start("127.0.0.1:0", http.NotFoundHandler())
	defer func() { _ = erste.Close() }()

	adresse := erste.Addr()
	if adresse == "" {
		t.Fatal("kein Listener")
	}

	zweite := Start(adresse, http.NotFoundHandler())
	if err_wait := zweite.Wait(); err_wait == nil {
		t.Fatal("erwartet Bind-Fehler aus Wait")
	}
}
