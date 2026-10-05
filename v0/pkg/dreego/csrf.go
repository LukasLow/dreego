package dreego

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

// csrfSessionKey is where the CSRF token lives inside the (HttpOnly) session.
const csrfSessionKey = "csrf_token"

// csrfFieldName is the form field name the token is submitted under.
const csrfFieldName = "_csrf"

// unsafeMethods need CSRF protection.
func isUnsafeMethod(method string) bool {
	switch method {
	case "POST", "PUT", "PATCH", "DELETE":
		return true
	default:
		return false
	}
}

// csrfToken returns the session's CSRF token, creating one on first use.
func (c *Ctx) csrfToken() string {
	if !c.csrfEnabled || c.store == nil {
		return ""
	}

	if vorhanden := c.sessionGet(csrfSessionKey); vorhanden != "" {
		return vorhanden
	}

	neuer := neuerCSRFToken()
	if neuer == "" {
		return ""
	}
	c.sessionSet(csrfSessionKey, neuer)
	return neuer
}

// csrfValid reports whether an unsafe request carries the right token.
func (c *Ctx) csrfValid() bool {
	if c.store == nil {
		return false
	}

	erwartet := c.sessionGet(csrfSessionKey)
	if erwartet == "" {
		// No token issued yet: nothing can be valid. This blocks a POST that
		// never saw a GET, which is the safe default.
		return false
	}

	geliefert := c.r.FormValue(csrfFieldName)
	if geliefert == "" {
		geliefert = c.r.Header.Get("X-CSRF-Token")
	}

	return subtle.ConstantTimeCompare([]byte(erwartet), []byte(geliefert)) == 1
}

// CSRFInput renders the hidden field a form must include. It is a Node, so
// there is no |raw anymore:
//
//	Form(method("post"), c.CSRFInput(), …)
func (c *Ctx) CSRFInput() g.Node {
	token := c.csrfToken()
	if token == "" {
		return g.Group(nil)
	}
	return h.Input(h.Type("hidden"), h.Name(csrfFieldName), h.Value(token))
}

// neuerCSRFToken returns a fresh 128-bit random token.
func neuerCSRFToken() string {
	var roh [16]byte
	_, err_lesen := rand.Read(roh[:])
	if err_lesen != nil {
		return ""
	}
	return base64.RawURLEncoding.EncodeToString(roh[:])
}
