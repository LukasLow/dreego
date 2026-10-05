package i18n

import (
	"strings"
	"testing"
)

func testBundle(t *testing.T) *Bundle {
	t.Helper()
	return New(Config{
		Default: "de",
		Messages: map[string]map[string]string{
			"de": {
				"login.title": "Anmelden",
				"welcome":     "Hallo {name}",
				"mails":       "{n} Nachrichten",
				"mails.one":   "{n} Nachricht",
				"mails.other": "{n} Nachrichten",
			},
			"en": {
				"login.title": "Sign in",
				"welcome":     "Hello {name}",
			},
		},
	})
}

func TestTUebersetzt(t *testing.T) {
	b := testBundle(t)
	if got := b.T("de", "login.title"); got != "Anmelden" {
		t.Fatalf("de = %q", got)
	}
	if got := b.T("en", "login.title"); got != "Sign in" {
		t.Fatalf("en = %q", got)
	}
}

func TestPlatzhalter(t *testing.T) {
	b := testBundle(t)
	got := b.T("de", "welcome", Args{"name": "Lukas"})
	if got != "Hallo Lukas" {
		t.Fatalf("platzhalter nicht ersetzt: %q", got)
	}
}

func TestFallbackAufDefault(t *testing.T) {
	b := testBundle(t)
	// "mails" fehlt in en -> Fallback auf Default (de).
	got := b.T("en", "mails", Args{"n": 3})
	if got != "3 Nachrichten" {
		t.Fatalf("fallback auf default fehlt: %q", got)
	}
}

func TestEnHatEigenenText(t *testing.T) {
	b := testBundle(t)
	got := b.T("en", "welcome", Args{"name": "Lukas"})
	if got != "Hello Lukas" {
		t.Fatalf("en-text falsch: %q", got)
	}
}

func TestFehlenderKeySichtbar(t *testing.T) {
	b := testBundle(t)
	got := b.T("de", "gibts.nicht")
	if got != "gibts.nicht" {
		t.Fatalf("fehlender key soll sichtbar sein: %q", got)
	}
}

func TestUnbekannteLocaleNutztDefault(t *testing.T) {
	b := testBundle(t)
	got := b.T("fr", "login.title")
	if got != "Anmelden" {
		t.Fatalf("unbekannte locale soll default nutzen: %q", got)
	}
}

func TestPluralEin(t *testing.T) {
	b := testBundle(t)
	got := b.Tn("de", "mails", 1, nil)
	if got != "1 Nachricht" {
		t.Fatalf("singular falsch: %q", got)
	}
}

func TestPluralMehrere(t *testing.T) {
	b := testBundle(t)
	got := b.Tn("de", "mails", 5, nil)
	if got != "5 Nachrichten" {
		t.Fatalf("plural falsch: %q", got)
	}
}

func TestPluralOhneVariantenNutztBasis(t *testing.T) {
	b := testBundle(t)
	// en hat keine mails.one/other -> Fallback auf Default de.
	got := b.Tn("en", "mails", 3, nil)
	if !strings.Contains(got, "3") {
		t.Fatalf("plural fallback falsch: %q", got)
	}
}
