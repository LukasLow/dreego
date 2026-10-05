package dreego

import (
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func flashApp(t *testing.T) *App {
	t.Helper()
	app := NewApp()
	app.SetSessionStore(testStore(t))
	app.Page(Page{
		Path: "/flash",
		Get: func(c *Ctx) g.Node {
			return h.P(g.Text("msg=" + c.FlashGet("fehler")))
		},
		Post: func(c *Ctx) g.Node {
			c.Flash("fehler", "kaputt")
			return Redirect("/flash", 303)
		},
	})
	return app
}

// TestFlashUeberRedirect prueft den echten Ablauf: GET holt Session+Token,
// POST mit gueltigem Token setzt die Flash und antwortet mit Redirect, der
// naechste GET liest die Flash genau einmal.
func TestFlashUeberRedirect(t *testing.T) {
	app := flashApp(t)

	// 1. GET: Session-Cookie besorgen.
	erste := httptest.NewRecorder()
	app.ServeHTTP(erste, httptest.NewRequest("GET", "/flash", nil))
	cookies := erste.Result().Cookies()
	if len(cookies) == 0 {
		t.Fatal("kein Session-Cookie beim GET")
	}
	sessionCookie := cookies[len(cookies)-1]

	// Token aus der Session lesen.
	tokenAnfrage := httptest.NewRequest("GET", "/flash", nil)
	tokenAnfrage.AddCookie(sessionCookie)
	token, _ := app.store.Get(tokenAnfrage, csrfSessionKey)
	if token == "" {
		t.Fatal("kein CSRF-Token in der Session")
	}

	// 2. POST mit gueltigem Token setzt die Flash.
	form := url.Values{csrfFieldName: {token}}
	postAnfrage := httptest.NewRequest("POST", "/flash", strings.NewReader(form.Encode()))
	postAnfrage.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	postAnfrage.AddCookie(sessionCookie)
	post := httptest.NewRecorder()
	app.ServeHTTP(post, postAnfrage)

	if post.Code != 303 {
		t.Fatalf("POST status = %d, erwartet 303", post.Code)
	}
	// Das Cookie aus der POST-Antwort traegt jetzt die Flash.
	postCookies := post.Result().Cookies()
	if len(postCookies) == 0 {
		t.Fatal("POST hat kein Cookie gesetzt")
	}
	flashCookie := postCookies[len(postCookies)-1]

	// 3. GET liest die Flash.
	mitFlash := httptest.NewRequest("GET", "/flash", nil)
	mitFlash.AddCookie(flashCookie)
	zweite := httptest.NewRecorder()
	app.ServeHTTP(zweite, mitFlash)
	if !strings.Contains(zweite.Body.String(), "msg=kaputt") {
		t.Fatalf("flash nicht gelesen: %s", zweite.Body.String())
	}

	// Der Lesen-Request hat die Flash verbraucht und ein neues Cookie gesetzt.
	leerCookies := zweite.Result().Cookies()
	leerCookie := flashCookie
	if len(leerCookies) > 0 {
		leerCookie = leerCookies[len(leerCookies)-1]
	}

	// 4. Nochmal GET: Meldung ist verbraucht.
	dritte := httptest.NewRecorder()
	dritteAnfrage := httptest.NewRequest("GET", "/flash", nil)
	dritteAnfrage.AddCookie(leerCookie)
	app.ServeHTTP(dritte, dritteAnfrage)
	if strings.Contains(dritte.Body.String(), "msg=kaputt") {
		t.Fatalf("flash nicht verbraucht: %s", dritte.Body.String())
	}
}

func TestSessionValRoundtrip(t *testing.T) {
	store := testStore(t)
	schreiber := httptest.NewRecorder()
	anfrage := httptest.NewRequest("GET", "/", nil)
	_ = store.Set(schreiber, anfrage, "session_token", "abc")

	cookie := schreiber.Result().Cookies()[0]
	naechste := httptest.NewRequest("GET", "/", nil)
	naechste.AddCookie(cookie)

	wert, _ := store.Get(naechste, "session_token")
	if wert != "abc" {
		t.Fatalf("session token = %q", wert)
	}
}
