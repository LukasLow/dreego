// Package dreego is a small Go web framework built on gomponents.
//
// There is no compiler and no template language: a page is a Go value, a
// handler is a Go function, and rendering is plain Go. The framework owns the
// web plumbing — routing, per-request context, layouts, form binding, security
// defaults — while the view comes from gomponents plus the scope addon
// (component-scoped CSS and JS).
//
// A minimal application:
//
//	app := dreego.NewApp()
//	app.Page(dreego.Page{
//	    Path: "/",
//	    Get: func(c *dreego.Ctx) g.Node {
//	        return c.Box(scope.CSS(`h1 { color: red }`), H1(g.Text("Hallo")))
//	    },
//	})
//	http.ListenAndServe(":8080", app)
package dreego

import (
	"log"
	"net/http"

	"github.com/LukasLow/dreego/v0/pkg/i18n"
	"github.com/LukasLow/dreego/v0/pkg/scope"
)

// App is one web application: a set of pages, middleware and security defaults.
// It implements http.Handler.
type App struct {
	pages      map[string]Page
	mux        *http.ServeMux
	middleware []func(http.Handler) http.Handler
	statics    []staticMount

	store *CookieStore
	csrf  bool

	// errors holds custom error pages per status code (SetErrorPage).
	errors map[int]ErrorFunc

	// i18n holds the translation bundle, localeDefault the fallback locale.
	i18n          *i18n.Bundle
	localeDefault string

	// noCompress disables the built-in gzip middleware (default: compress on).
	noCompress bool
}

// NewApp returns an empty application with secure defaults: CSRF on, built-in
// security headers, gzip compression and panic recovery.
func NewApp() *App {
	app := &App{
		pages: map[string]Page{},
		mux:   http.NewServeMux(),
		csrf:  true,
	}
	// Catch-all: unmatched paths first try embedded static files, then render the
	// app's 404 page. Method mismatches on registered patterns are still turned
	// into 405 by the ServeMux itself.
	app.mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if app.serveStatic(w, r) {
			return
		}
		app.serveError(w, r, http.StatusNotFound)
	})
	return app
}

// SetCompress turns the built-in gzip middleware on or off. It is on by
// default; turn it off only for a controlled comparison.
func (app *App) SetCompress(an bool) { app.noCompress = !an }

// SetSessionStore sets the session store. Without one, sessions, flash and CSRF
// are inert. The store also validates at startup.
func (app *App) SetSessionStore(store *CookieStore) { app.store = store }

// SetCSRF turns CSRF protection on or off for the whole app. It is on by
// default. An app without forms and without a cookie (e.g. a pure marketing
// site) may turn it off.
func (app *App) SetCSRF(an bool) { app.csrf = an }

// Page registers a page. The URL comes only from Page.Path — the file the page
// lives in has no meaning (no file-based routing).
//
// Dynamic segments use the Dreego-style brackets: [id] becomes a path value,
// [...rest] captures the remainder. Read them in the handler with c.Param(name).
func (app *App) Page(page Page) {
	app.pages[page.Path] = page

	// Register the pattern WITHOUT a method prefix, so a wrong method reaches
	// our own handler and becomes a proper 405 (with Allow) instead of falling
	// through to the "/" catch-all. The method is checked inside servePage.
	muster := toServeMuxPattern(page.Path)
	app.mux.HandleFunc(muster, func(w http.ResponseWriter, r *http.Request) {
		app.servePage(w, r, page)
	})
}

// Use appends middleware. Middleware runs in registration order around every
// request, before the page handler.
func (app *App) Use(mw func(http.Handler) http.Handler) {
	app.middleware = append(app.middleware, mw)
}

// Handler returns the full handler: built-in middleware (recovery, security
// headers) plus the app's own middleware, around the page dispatcher.
func (app *App) Handler() http.Handler {
	var kern http.Handler = http.HandlerFunc(app.dispatch)

	// User middleware wraps the dispatcher, in registration order.
	for index := len(app.middleware) - 1; index >= 0; index-- {
		kern = app.middleware[index](kern)
	}

	// Built-in middleware is outermost: recovery first, then compression (unless
	// disabled), then security headers.
	kette := SecurityHeaders(kern)
	if !app.noCompress {
		kette = Compress(kette)
	}
	return app.recoverMiddleware(kette)
}

// ServeHTTP lets the App be used directly as an http.Handler.
func (app *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	app.Handler().ServeHTTP(w, r)
}

// dispatch handles the health endpoint and otherwise delegates to the mux.
// Unmatched paths are caught by a "/" handler that renders the app's 404 page;
// method mismatches produce 405 inside the ServeMux.
func (app *App) dispatch(w http.ResponseWriter, r *http.Request) {
	// Health endpoint, always available (Dreego parity).
	if r.URL.Path == "/health" {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
		return
	}

	app.mux.ServeHTTP(w, r)
}

// servePage runs one page: context, CSRF, method, session, render or redirect.
func (app *App) servePage(w http.ResponseWriter, r *http.Request, seite Page) {
	// API pages respond with JSON/XML/etc. — no layout, no HTML render.
	if seite.API != nil {
		app.serveAPI(w, r, seite)
		return
	}
	// Everything below runs inside the page, so per-page middleware only wraps
	// this route's handling. The context is built inside, so middleware that
	// replaces w/r is respected.
	seitenHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nonce := scope.NewNonce()
		if nonce == "" {
			// crypto/rand failed: refuse rather than serve pages whose assets
			// the CSP would silently block.
			app.serveError(w, r, http.StatusInternalServerError)
			return
		}
		collector := scope.New()
		collector.SetNonce(nonce)

		// Security default: a strict CSP. 'self' allows our own static files
		// (e.g. /public/app.js), the nonce allows our inline component JS/CSS.
		// Third-party scripts and inline code without a nonce are still blocked.
		w.Header().Set("Content-Security-Policy",
			"default-src 'self'; "+
				"script-src 'self' 'nonce-"+nonce+"'; "+
				"style-src 'self' 'nonce-"+nonce+"'; "+
				"img-src 'self' data:; "+
				"base-uri 'none'; "+
				"form-action 'self'; "+
				"frame-ancestors 'none'; "+
				"object-src 'none'")

		ctx := newCtx(w, r, seite, app.store, app.csrf, nonce, collector, app.resolveLocale(r), app.i18n)

		// Method first: a POST to a GET-only page is 405, not a CSRF error.
		handler := seite.handlerFor(r.Method)
		if handler == nil {
			w.Header().Set("Allow", seite.allowedMethods())
			app.serveError(w, r, http.StatusMethodNotAllowed)
			return
		}

		// CSRF: unsafe methods must carry the right token. Without a store there
		// is nowhere to keep a token — fail loudly, never half-guard.
		if app.csrf && isUnsafeMethod(r.Method) && app.store == nil {
			app.serveError(w, r, http.StatusInternalServerError)
			return
		}

		if app.csrf && isUnsafeMethod(r.Method) && !ctx.csrfValid() {
			app.serveError(w, r, http.StatusForbidden)
			return
		}

		// Safe methods: make sure a token exists, so any rendered form has one.
		if app.csrf && !isUnsafeMethod(r.Method) {
			_ = ctx.csrfToken()
		}

		body := handler(ctx)

		// Write the session as ONE cookie before anything goes to the client.
		ctx.flushSession()

		// A failed session write (e.g. too large) must not look like success.
		if ctx.sessionError != nil {
			log.Printf("dreego: Session-Schreibfehler für %s: %v", r.URL.Path, ctx.sessionError)
			app.serveError(w, r, http.StatusInternalServerError)
			return
		}

		// A redirect is itself a Node: set the status and stop, no body.
		if ziel, istRedirect := body.(redirect); istRedirect {
			w.Header().Set("Location", ziel.url)
			w.WriteHeader(ziel.code)
			return
		}

		layout := seite.Layout
		if layout == nil {
			layout = defaultLayout
		}

		dokument := layout(ctx, seite.headNode(ctx), body)

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		err_render := dokument.Render(w)
		if err_render != nil {
			// The response may already be partially written; report at least.
			log.Printf("dreego: render fehlgeschlagen für %s: %v", r.URL.Path, err_render)
		}
	})

	// Per-page middleware wraps the page handling, in registration order.
	gefuehrt := seitenHandler
	for index := len(seite.Middleware) - 1; index >= 0; index-- {
		gefuehrt = seite.Middleware[index](gefuehrt).ServeHTTP
	}
	gefuehrt(w, r)
}
