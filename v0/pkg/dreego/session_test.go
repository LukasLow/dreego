package dreego

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func testStore(t *testing.T) *CookieStore {
	t.Helper()
	store, err_store := NewCookieStore([]byte("0123456789abcdef0123456789abcdef"))
	if err_store != nil {
		t.Fatalf("store: %v", err_store)
	}
	return store
}

func TestSecretZuKurz(t *testing.T) {
	_, err := NewCookieStore([]byte("kurz"))
	if err != ErrSecretTooShort {
		t.Fatalf("erwartet ErrSecretTooShort, bekam %v", err)
	}
}

func TestCookieMaxAge(t *testing.T) {
	store := testStore(t)
	store.SetCookiePolicy(CookiePolicy{MaxAge: Days(30)})

	schreiber := httptest.NewRecorder()
	anfrage := httptest.NewRequest("GET", "/", nil)
	err_set := store.Set(schreiber, anfrage, "user", "lukas")
	if err_set != nil {
		t.Fatalf("set: %v", err_set)
	}

	cookie := schreiber.Result().Cookies()[0]
	if cookie.MaxAge != 30*24*60*60 {
		t.Fatalf("MaxAge = %d, erwartet 2592000", cookie.MaxAge)
	}
	if cookie.Expires.IsZero() {
		t.Fatal("Expires soll gesetzt sein")
	}
}

func TestCookieOhneMaxAgeIstSessionCookie(t *testing.T) {
	store := testStore(t)

	schreiber := httptest.NewRecorder()
	anfrage := httptest.NewRequest("GET", "/", nil)
	_ = store.Set(schreiber, anfrage, "user", "lukas")

	cookie := schreiber.Result().Cookies()[0]
	if cookie.MaxAge != 0 || !cookie.Expires.IsZero() {
		t.Fatalf("ohne MaxAge muss das Cookie ein Session-Cookie sein")
	}
}

func TestSessionRoundTrip(t *testing.T) {
	store := testStore(t)
	schreiber := httptest.NewRecorder()
	anfrage := httptest.NewRequest("GET", "/", nil)

	err_set := store.Set(schreiber, anfrage, "user", "lukas")
	if err_set != nil {
		t.Fatalf("set: %v", err_set)
	}

	// Cookie aus der Antwort zurück in die nächste Anfrage.
	antwort := schreiber.Result()
	cookies := antwort.Cookies()
	if len(cookies) == 0 {
		t.Fatal("kein Cookie gesetzt")
	}

	naechste := httptest.NewRequest("GET", "/", nil)
	naechste.AddCookie(cookies[0])

	wert, err_get := store.Get(naechste, "user")
	if err_get != nil {
		t.Fatalf("get: %v", err_get)
	}
	if wert != "lukas" {
		t.Fatalf("wert = %q", wert)
	}
}

func TestCookieIstHttpOnlyUndVerschluesselt(t *testing.T) {
	store := testStore(t)
	schreiber := httptest.NewRecorder()
	anfrage := httptest.NewRequest("GET", "/", nil)
	_ = store.Set(schreiber, anfrage, "user", "geheim")

	cookie := schreiber.Result().Cookies()[0]
	if !cookie.HttpOnly {
		t.Fatal("HttpOnly fehlt")
	}
	if cookie.Value == "geheim" || contains(cookie.Value, "geheim") {
		t.Fatalf("Cookie enthält Klartext: %q", cookie.Value)
	}
	if cookie.SameSite != http.SameSiteLaxMode {
		t.Fatalf("SameSite = %v", cookie.SameSite)
	}
}

func TestManipuliertesCookieWirdAbgelehnt(t *testing.T) {
	store := testStore(t)
	anfrage := httptest.NewRequest("GET", "/", nil)
	anfrage.AddCookie(&http.Cookie{Name: store.Policy().Name, Value: "boese-daten"})

	_, err_get := store.Get(anfrage, "user")
	if err_get == nil {
		t.Fatal("manipuliertes Cookie soll abgelehnt werden")
	}
}

func TestDestroyLoeschtSession(t *testing.T) {
	store := testStore(t)
	schreiber := httptest.NewRecorder()
	anfrage := httptest.NewRequest("GET", "/", nil)
	_ = store.Set(schreiber, anfrage, "a", "1")
	_ = store.Set(schreiber, anfrage, "b", "2")
	_ = store.Destroy(schreiber, anfrage)

	// Das LETZTE Cookie ist das Destroy-Cookie: es soll leer sein.
	cookies := schreiber.Result().Cookies()
	cookie := cookies[len(cookies)-1]
	naechste := httptest.NewRequest("GET", "/", nil)
	naechste.AddCookie(cookie)

	wert, _ := store.Get(naechste, "a")
	if wert != "" {
		t.Fatalf("nach Destroy noch Wert da: %q", wert)
	}
}

func contains(text string, teil string) bool {
	for index := 0; index+len(teil) <= len(text); index++ {
		if text[index:index+len(teil)] == teil {
			return true
		}
	}
	return false
}
