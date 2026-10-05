package dreego

import (
	"net/http"
	"net/http/httptest"
	"testing"

	g "maragu.dev/gomponents"

	"github.com/LukasLow/dreego/v0/pkg/i18n"
)

func i18nTestApp() *App {
	bundle := i18n.New(i18n.Config{
		Default: "de",
		Messages: map[string]map[string]string{
			"de": {"hallo": "Hallo"},
			"en": {"hallo": "Hello"},
			"fr": {"hallo": "Bonjour"},
		},
	})
	app := NewApp()
	app.SetI18n(bundle)
	app.SetDefaultLocale("de")
	return app
}

func TestLocaleAusCookie(t *testing.T) {
	app := i18nTestApp()
	anfrage := httptest.NewRequest("GET", "/", nil)
	anfrage.AddCookie(&http.Cookie{Name: "lang", Value: "en"})
	if got := app.resolveLocale(anfrage); got != "en" {
		t.Fatalf("cookie-locale = %q", got)
	}
}

func TestLocaleAusAcceptLanguage(t *testing.T) {
	app := i18nTestApp()
	anfrage := httptest.NewRequest("GET", "/", nil)
	anfrage.Header.Set("Accept-Language", "fr-FR,fr;q=0.9,en;q=0.5")
	if got := app.resolveLocale(anfrage); got != "fr" {
		t.Fatalf("accept-language = %q", got)
	}
}

func TestLocaleQValue(t *testing.T) {
	app := i18nTestApp()
	anfrage := httptest.NewRequest("GET", "/", nil)
	anfrage.Header.Set("Accept-Language", "en;q=0.3, fr;q=0.9")
	if got := app.resolveLocale(anfrage); got != "fr" {
		t.Fatalf("q-value = %q", got)
	}
}

func TestLocaleDefault(t *testing.T) {
	app := i18nTestApp()
	anfrage := httptest.NewRequest("GET", "/", nil)
	if got := app.resolveLocale(anfrage); got != "de" {
		t.Fatalf("default = %q", got)
	}
}

func TestLocaleUnbekannteSpracheCookieIgnoriert(t *testing.T) {
	app := i18nTestApp()
	anfrage := httptest.NewRequest("GET", "/", nil)
	anfrage.AddCookie(&http.Cookie{Name: "lang", Value: "xx"}) // nicht im Bundle
	if got := app.resolveLocale(anfrage); got != "de" {
		t.Fatalf("unbekanntes cookie-locale soll default sein: %q", got)
	}
}

func TestCtxTImHandler(t *testing.T) {
	app := i18nTestApp()
	var gesehen string
	var gesehenLocale string
	app.Page(Page{Path: "/x", Get: func(c *Ctx) g.Node {
		gesehen = c.T("hallo")
		gesehenLocale = c.Locale()
		return g.Text(gesehen)
	}})

	anfrage := httptest.NewRequest("GET", "/x", nil)
	anfrage.Header.Set("Accept-Language", "en")
	schreiber := httptest.NewRecorder()
	app.ServeHTTP(schreiber, anfrage)

	if gesehenLocale != "en" {
		t.Fatalf("ctx locale = %q", gesehenLocale)
	}
	if gesehen != "Hello" {
		t.Fatalf("ctx.T = %q", gesehen)
	}
}
