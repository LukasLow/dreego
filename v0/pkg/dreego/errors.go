package dreego

import (
	"encoding/json"
	"net/http"

	g "github.com/LukasLow/dreego/v0/pkg/dom"
	"strings"

	"github.com/LukasLow/dreego/v0/pkg/scope"
)

// ErrorFunc renders the complete response body for one status code. It has a
// working Ctx (with a fresh nonce), so it can call c.Document to render inside
// the site's layout — exactly like a normal page.
type ErrorFunc func(c *Ctx) g.View

// SetErrorPage registers a custom page for a status code (404, 403, 405, 429,
// 500, …):
//
//	app.SetErrorPage(404, func(c *dreego.Ctx) g.View {
//	    return c.Document("de", nil, errorBox(c, "Nicht gefunden", "…"))
//	})
//
// A status without an entry uses a small built-in page — never the bare Go text.
func (app *App) SetErrorPage(status int, handler ErrorFunc) {
	if app.errors == nil {
		app.errors = map[int]ErrorFunc{}
	}
	app.errors[status] = handler
}

// serveError writes an error response: JSON for API/JSON clients, HTML (custom
// page or built-in fallback) otherwise. It never leaks internal detail.
func (app *App) serveError(w http.ResponseWriter, r *http.Request, status int) {
	if wantsJSON(r) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(status)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"error":  http.StatusText(status),
			"status": status,
		})
		return
	}

	nonce := scope.NewNonce()
	collector := scope.New()
	collector.SetNonce(nonce)

	ctx := &Ctx{
		w:         w,
		r:         r,
		nonce:     nonce,
		collector: collector,
		locale:    app.resolveLocale(r),
		bundle:    app.i18n,
	}

	var inhalt g.View
	if handler, vorhanden := app.errors[status]; vorhanden {
		inhalt = handler(ctx)
	}
	if inhalt == nil {
		inhalt = fallbackErrorDocument(status)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	_ = inhalt.Render(w)
}

// wantsJSON reports whether the client asked for JSON: an /api path, or an
// Accept header that prefers application/json without text/html.
func wantsJSON(r *http.Request) bool {
	if strings.HasPrefix(r.URL.Path, "/api") {
		return true
	}
	accept := r.Header.Get("Accept")
	return strings.Contains(accept, "application/json") && !strings.Contains(accept, "text/html")
}

// fallbackErrorDocument is the built-in page used when no custom one is set.
// It is minimal, valid and safe.
func fallbackErrorDocument(status int) g.View {
	return g.Group([]g.View{
		g.Raw("<!DOCTYPE html>"),
		g.HTML(g.Lang("de"),
			g.Head(
				g.Meta(g.Charset("utf-8")),
				g.Meta(g.Name("viewport"), g.Content("width=device-width, initial-scale=1")),
				g.TitleEl(g.Text(http.StatusText(status))),
			),
			g.Body(
				g.H1(g.Text(itoa(status)+" "+http.StatusText(status))),
				g.P(g.A(g.Href("/"), g.Text("Zur Startseite"))),
			),
		),
	})
}

// itoa is a tiny int->string without importing strconv here.
func itoa(zahl int) string {
	if zahl == 0 {
		return "0"
	}
	var ziffern []byte
	for zahl > 0 {
		ziffern = append([]byte{byte('0' + zahl%10)}, ziffern...)
		zahl /= 10
	}
	return string(ziffern)
}
