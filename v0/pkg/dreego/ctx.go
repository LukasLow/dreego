package dreego

import (
	"net/http"

	g "maragu.dev/gomponents"

	"github.com/LukasLow/dreego/v0/pkg/i18n"
	"github.com/LukasLow/dreego/v0/pkg/scope"
)

// Ctx is the per-request context. It carries what security needs — the
// response writer (to set cookies), the request (to read cookies and forms),
// the CSP nonce, the asset collector, the session store and whether CSRF is
// enabled — plus the current page.
//
// The session is loaded once per request into memory. All reads and writes go
// through this map, and the whole session is written back as ONE cookie at the
// end of the request. This avoids the classic bug where several Set-Cookie
// headers each carry a stale copy of the session and overwrite each other.
//
// One Ctx per request.
type Ctx struct {
	w         http.ResponseWriter
	r         *http.Request
	nonce     string
	collector *scope.Collector
	page      Page

	store       *CookieStore
	csrfEnabled bool

	locale string
	bundle *i18n.Bundle

	session      map[string]string
	sessionDirty bool
	sessionError error
}

// newCtx loads the session for a request.
func newCtx(w http.ResponseWriter, r *http.Request, seite Page, store *CookieStore, csrf bool, nonce string, collector *scope.Collector, locale string, bundle *i18n.Bundle) *Ctx {
	session := map[string]string{}
	if store != nil {
		session = store.Load(r)
	}

	return &Ctx{
		w:           w,
		r:           r,
		nonce:       nonce,
		collector:   collector,
		page:        seite,
		store:       store,
		csrfEnabled: csrf,
		locale:      locale,
		bundle:      bundle,
		session:     session,
	}
}

// Request returns the underlying *http.Request.
func (c *Ctx) Request() *http.Request { return c.r }

// Response returns the underlying http.ResponseWriter.
func (c *Ctx) Response() http.ResponseWriter { return c.w }

// Nonce returns the CSP nonce for this response.
func (c *Ctx) Nonce() string { return c.nonce }

// Collector returns the per-page asset collector.
func (c *Ctx) Collector() *scope.Collector { return c.collector }

// AddCritical registers CSS that is inlined in <head> before external
// stylesheets, so the first paint already carries the base colours and font —
// no white flash while an external stylesheet loads.
func (c *Ctx) AddCritical(css string) { c.collector.AddCritical(css) }

// Box is a component's scope container, registered with the request collector.
// Components use it so their CSS and JS is emitted once per page, with the CSP
// nonce:
//
//	func Card(c *dreego.Ctx) g.Node {
//	    return c.Box(scope.CSS(`.card { … }`), Div(Class("card"), …))
//	}
func (c *Ctx) Box(parts ...g.Node) g.Node { return c.collector.Box(parts...) }

// Nav returns the active navigation marker declared on the page.
func (c *Ctx) Nav() string { return c.page.Nav }

// FormValue returns the first value for the named form field (query or body).
func (c *Ctx) FormValue(name string) string { return c.r.FormValue(name) }

// Param returns a dynamic path value from a [name] segment:
//
//	var Project = dreego.Page{Path: "/projekte/[id]", Get: getProject}
//	func getProject(c *dreego.Ctx) g.Node { id := c.Param("id") }
func (c *Ctx) Param(name string) string { return c.r.PathValue(name) }

// Redirect returns a Node that redirects the client. Typical use is the
// Post/Redirect/Get pattern. It is a method so it reads naturally on the ctx:
//
//	return c.Redirect("/auth/login", 303)
func (c *Ctx) Redirect(url string, code int) g.Node { return Redirect(url, code) }

// sessionGet reads a value from the in-memory session.
func (c *Ctx) sessionGet(key string) string { return c.session[key] }

// sessionSet writes a value into the in-memory session.
func (c *Ctx) sessionSet(key string, val string) {
	if val == "" {
		delete(c.session, key)
	} else {
		c.session[key] = val
	}
	c.sessionDirty = true
}

// sessionDel removes a key from the in-memory session.
func (c *Ctx) sessionDel(key string) {
	delete(c.session, key)
	c.sessionDirty = true
}

// flushSession writes the session back as one cookie, if anything changed.
// It must run before any header is written.
//
// A write failure (e.g. the session exceeds the cookie size limit) is reported
// on the ctx; the caller surfaces it as a deliberate error rather than losing
// the login silently.
func (c *Ctx) flushSession() {
	if !c.sessionDirty || c.store == nil {
		return
	}
	err_save := c.store.Save(c.w, c.r, c.session)
	if err_save != nil {
		c.sessionError = err_save
		return
	}
	c.sessionDirty = false
}
