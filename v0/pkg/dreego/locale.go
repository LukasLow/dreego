package dreego

import (
	"net/http"
	"sort"
	"strconv"
	"strings"

	g "github.com/LukasLow/dreego/v0/pkg/dom"

	"github.com/LukasLow/dreego/v0/pkg/i18n"
)

// SetI18n attaches a translation bundle. Once set, c.T and c.Tn resolve
// messages, and the locale comes from the request (see resolveLocale).
func (app *App) SetI18n(bundle *i18n.Bundle) { app.i18n = bundle }

// SetDefaultLocale sets the locale used when the request gives no usable hint.
// It overrides the bundle's own default for request resolution.
func (app *App) SetDefaultLocale(locale string) { app.localeDefault = locale }

// resolveLocale picks the request locale:
//
//  1. the "lang" cookie (an explicit user choice),
//  2. the Accept-Language header (best quality match against the bundle),
//  3. the app's default locale.
func (app *App) resolveLocale(r *http.Request) string {
	// 1. Cookie
	if cookie, err_cookie := r.Cookie("lang"); err_cookie == nil && cookie.Value != "" {
		if app.i18n == nil || app.i18n.Has(cookie.Value) {
			return cookie.Value
		}
	}

	// 2. Accept-Language
	if app.i18n != nil {
		if gewaehlt := bestAcceptLanguage(r.Header.Get("Accept-Language"), app.i18n.Locales()); gewaehlt != "" {
			return gewaehlt
		}
	}

	// 3. Default
	if app.localeDefault != "" {
		return app.localeDefault
	}
	if app.i18n != nil {
		return app.i18n.Default()
	}
	return "de"
}

// bestAcceptLanguage returns the best available locale for an Accept-Language
// header, or "" when none matches. It respects q-values.
func bestAcceptLanguage(header string, verfuegbar []string) string {
	if header == "" || len(verfuegbar) == 0 {
		return ""
	}

	// Available locales as a set, and their primary subtags ("de-DE" -> "de").
	menge := map[string]bool{}
	for _, sprache := range verfuegbar {
		menge[sprache] = true
	}

	type kandidat struct {
		sprache string
		q       float64
		index   int
	}
	var kandidaten []kandidat

	for index, teil := range strings.Split(header, ",") {
		teil = strings.TrimSpace(teil)
		if teil == "" {
			continue
		}
		sprache, q := parseAcceptLanguage(teil)
		if sprache == "" {
			continue
		}
		if q <= 0 {
			continue
		}

		if menge[sprache] {
			kandidaten = append(kandidaten, kandidat{sprache, q, index})
			continue
		}
		// "de-DE" -> "de" primary match
		if strich := strings.IndexByte(sprache, '-'); strich > 0 {
			haupt := sprache[:strich]
			if menge[haupt] {
				kandidaten = append(kandidaten, kandidat{haupt, q, index})
			}
		}
		// "*" -> any available
		if sprache == "*" {
			sort.Strings(verfuegbar)
			kandidaten = append(kandidaten, kandidat{verfuegbar[0], q, index})
		}
	}

	if len(kandidaten) == 0 {
		return ""
	}

	sort.SliceStable(kandidaten, func(i, j int) bool {
		if kandidaten[i].q != kandidaten[j].q {
			return kandidaten[i].q > kandidaten[j].q
		}
		return kandidaten[i].index < kandidaten[j].index
	})

	return kandidaten[0].sprache
}

// parseAcceptLanguage splits one entry "de-DE;q=0.8" into locale and q.
func parseAcceptLanguage(teil string) (string, float64) {
	name, params, hatParams := strings.Cut(teil, ";")
	name = strings.TrimSpace(name)
	q := 1.0
	if hatParams {
		for _, param := range strings.Split(params, ";") {
			param = strings.TrimSpace(param)
			if wert, gefunden := strings.CutPrefix(param, "q="); gefunden {
				if zahl, err_parse := strconv.ParseFloat(strings.TrimSpace(wert), 64); err_parse == nil {
					q = zahl
				}
			}
		}
	}
	return name, q
}

// Locale returns the locale chosen for this request.
func (c *Ctx) Locale() string { return c.locale }

// T translates a key in the request locale.
func (c *Ctx) T(key string, args ...i18n.Args) string {
	if c.bundle == nil {
		return key
	}
	return c.bundle.T(c.locale, key, args...)
}

// Text translates a key and returns it as a dom view, so it can be used
// directly in element position:
//
//	H1(c.Text("start.titel"))
func (c *Ctx) Text(key string, args ...i18n.Args) g.View {
	return g.Text(c.T(key, args...))
}

// Tn translates a plural key in the request locale (n == 1 -> ".one").
func (c *Ctx) Tn(key string, anzahl int, args ...i18n.Args) string {
	if c.bundle == nil {
		return key
	}
	return c.bundle.Tn(c.locale, key, anzahl, args...)
}
