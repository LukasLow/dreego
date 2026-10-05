// Package i18n is dreego's small translation layer. No build step, no CLI:
// catalogs live as Go maps or JSON files (via go:embed).
//
//	bundle := i18n.New(i18n.Config{
//	    Default: "de",
//	    Messages: map[string]map[string]string{
//	        "de": {"login.title": "Anmelden"},
//	        "en": {"login.title": "Sign in"},
//	    },
//	})
//	bundle.T("en", "login.title")          // "Sign in"
//	bundle.T("de", "welcome", i18n.Args{"name": "Lukas"})
//	bundle.Tn("de", "mails", 5, nil)       // uses "mails.one" / "mails.other"
package i18n

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"path"
	"sort"
	"strings"
)

// Args are the placeholder values for {name} substitution.
type Args map[string]any

// Config holds the catalogs and the fallback locale.
type Config struct {
	// Default is the fallback locale (e.g. "de").
	Default string
	// Messages maps locale -> key -> text.
	Messages map[string]map[string]string
}

// Bundle resolves translations. Safe for concurrent reads.
type Bundle struct {
	standard string
	messages map[string]map[string]string
}

// New builds a bundle from a config.
func New(cfg Config) *Bundle {
	if cfg.Messages == nil {
		cfg.Messages = map[string]map[string]string{}
	}
	return &Bundle{standard: cfg.Default, messages: cfg.Messages}
}

// Open loads every *.json file from dir inside the filesystem. Each file is
// named after its locale (de.json, en.json). The JSON is a flat key/value map.
//
//	//go:embed locales
//	var localesFS embed.FS
//	bundle, err := i18n.Open(localesFS, "locales", "de")
func Open(dateiSystem fs.FS, dir string, standard string) (*Bundle, error) {
	eintraege, err_lesen := fs.ReadDir(dateiSystem, dir)
	if err_lesen != nil {
		return nil, err_lesen
	}

	messages := map[string]map[string]string{}
	for _, eintrag := range eintraege {
		if eintrag.IsDir() || path.Ext(eintrag.Name()) != ".json" {
			continue
		}
		sprache := strings.TrimSuffix(eintrag.Name(), ".json")
		pfad := path.Join(dir, eintrag.Name())

		daten, err_datei := fs.ReadFile(dateiSystem, pfad)
		if err_datei != nil {
			return nil, err_datei
		}

		katalog := map[string]string{}
		err_json := json.Unmarshal(daten, &katalog)
		if err_json != nil {
			return nil, fmt.Errorf("i18n: %s: %w", pfad, err_json)
		}
		messages[sprache] = katalog
	}

	return New(Config{Default: standard, Messages: messages}), nil
}

// Has reports whether a locale catalog was loaded.
func (b *Bundle) Has(locale string) bool {
	_, vorhanden := b.messages[locale]
	return vorhanden
}

// Default returns the fallback locale.
func (b *Bundle) Default() string { return b.standard }

// Locales lists the loaded locales, sorted.
func (b *Bundle) Locales() []string {
	sprachen := make([]string, 0, len(b.messages))
	for sprache := range b.messages {
		sprachen = append(sprachen, sprache)
	}
	sort.Strings(sprachen)
	return sprachen
}

// T resolves key in locale, with {name} substitution from args.
//
// Order: the requested locale, then the default locale. If neither has the key,
// the key itself is returned — visible, never an empty string.
//
//	bundle.T("en", "login.title")
//	bundle.T("de", "welcome", i18n.Args{"name": "Lukas"})
func (b *Bundle) T(locale string, key string, args ...Args) string {
	var werte Args
	if len(args) > 0 {
		werte = args[0]
	}

	text, gefunden := b.lookup(locale, key)
	if !gefunden {
		return key
	}
	return substitute(text, werte)
}

// Tn resolves a plural key: "key.one" for n == 1, otherwise "key.other". If the
// plural variants are missing, the base key is used as fallback.
func (b *Bundle) Tn(locale string, key string, anzahl int, args ...Args) string {
	variante := key + ".other"
	if anzahl == 1 {
		variante = key + ".one"
	}

	if _, gefunden := b.lookup(locale, variante); !gefunden {
		variante = key // fallback: base key
	}

	werte := Args{}
	if len(args) > 0 {
		for k, v := range args[0] {
			werte[k] = v
		}
	}
	werte["n"] = anzahl

	return b.T(locale, variante, werte)
}

// lookup finds a key: try the locale, then the default locale.
func (b *Bundle) lookup(locale string, key string) (string, bool) {
	if katalog, vorhanden := b.messages[locale]; vorhanden {
		if text, da := katalog[key]; da {
			return text, true
		}
	}
	if locale != b.standard {
		if katalog, vorhanden := b.messages[b.standard]; vorhanden {
			if text, da := katalog[key]; da {
				return text, true
			}
		}
	}
	return "", false
}

// substitute replaces {name} placeholders with values from args.
func substitute(text string, args Args) string {
	if len(args) == 0 || !strings.Contains(text, "{") {
		return text
	}
	var bau strings.Builder
	rest := text
	for {
		start := strings.IndexByte(rest, '{')
		if start < 0 {
			bau.WriteString(rest)
			break
		}
		ende := strings.IndexByte(rest[start:], '}')
		if ende < 0 {
			bau.WriteString(rest)
			break
		}
		ende += start

		bau.WriteString(rest[:start])
		name := rest[start+1 : ende]

		if wert, da := args[name]; da {
			bau.WriteString(fmt.Sprint(wert))
		} else {
			// Unknown placeholder: keep it visible rather than dropping it.
			bau.WriteString(rest[start : ende+1])
		}
		rest = rest[ende+1:]
	}
	return bau.String()
}
