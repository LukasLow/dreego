package dreego

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"io"
	"log"
	"net/http"
)

// maxAPIBody is the read limit for API request bodies (1 MiB, like Dreego).
const maxAPIBody = 1 << 20

// ErrAPIBodyTooLarge is returned by Bind when the request body exceeds the limit.
var ErrAPIBodyTooLarge = errors.New("dreego: API-Body zu groß (max 1 MiB)")

// serveAPI runs an API page: CSRF for unsafe methods, session flush, then the
// handler. No layout and no HTML render happen.
func (app *App) serveAPI(w http.ResponseWriter, r *http.Request, seite Page) {
	// API responses need no CSP nonce and no asset collector.
	ctx := newCtx(w, r, seite, app.store, app.csrf, "", nil, app.resolveLocale(r), app.i18n)

	// CSRF applies to unsafe methods unless the app disabled it.
	if app.csrf && isUnsafeMethod(r.Method) {
		if app.store == nil {
			app.serveError(w, r, http.StatusInternalServerError)
			return
		}
		if !ctx.csrfValid() {
			app.serveError(w, r, http.StatusForbidden)
			return
		}
	}

	err_api := seite.API(ctx)

	// The session is written before any error response.
	ctx.flushSession()

	if err_api != nil {
		log.Printf("dreego: API-Fehler für %s: %v", r.URL.Path, err_api)
		// Only send an error if nothing was written yet. Handlers that already
		// wrote a status own the response.
		app.serveError(w, r, http.StatusInternalServerError)
		return
	}
}

// JSON marshals v, sets the status and the JSON content type.
func (c *Ctx) JSON(status int, v any) error {
	daten, err_marshal := json.Marshal(v)
	if err_marshal != nil {
		return err_marshal
	}
	c.w.Header().Set("Content-Type", "application/json; charset=utf-8")
	c.w.WriteHeader(status)
	_, err_schreiben := c.w.Write(daten)
	return err_schreiben
}

// XML marshals v, sets the status and the XML content type. An XML header is
// prepended, as is conventional.
func (c *Ctx) XML(status int, v any) error {
	daten, err_marshal := xml.Marshal(v)
	if err_marshal != nil {
		return err_marshal
	}
	c.w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	c.w.WriteHeader(status)
	_, err_kopf := c.w.Write([]byte(xml.Header))
	if err_kopf != nil {
		return err_kopf
	}
	_, err_schreiben := c.w.Write(daten)
	return err_schreiben
}

// Write sends a raw body with an explicit content type (Dreego's c.Write).
func (c *Ctx) Write(status int, contentType string, body []byte) error {
	if contentType != "" {
		c.w.Header().Set("Content-Type", contentType)
	}
	c.w.WriteHeader(status)
	_, err_schreiben := c.w.Write(body)
	return err_schreiben
}

// Bind reads the request body as JSON into target, limited to 1 MiB.
func (c *Ctx) Bind(target any) error {
	begrenzt := http.MaxBytesReader(c.w, c.r.Body, maxAPIBody)
	daten, err_lesen := io.ReadAll(begrenzt)
	if err_lesen != nil {
		var zuGross *http.MaxBytesError
		if errors.As(err_lesen, &zuGross) {
			return ErrAPIBodyTooLarge
		}
		return err_lesen
	}
	return json.Unmarshal(daten, target)
}
