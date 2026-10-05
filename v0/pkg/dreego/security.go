package dreego

import "strings"

// Security declares a page's (or the whole app's) response security policy.
//
// The zero value means "use dreego's strict default". Set it on the App for a
// site-wide policy, or on a single Page to override just that page — the same
// declarative style as the other Page fields.
//
//	var Upload = dreego.Page{
//	    Path: "/upload",
//	    Security: &dreego.Security{
//	        CSP: "default-src 'self'; script-src 'self' 'nonce-{nonce}' https://trusted.example",
//	    },
//	    Get:  getUpload,
//	}
type Security struct {
	// CSP is the Content-Security-Policy sent with the page.
	//
	//   - ""          -> the built-in strict default (StrictCSP)
	//   - "off"       -> no CSP at all for this page/app
	//   - anything else -> used verbatim, with {nonce} replaced by the
	//     per-request nonce (e.g. "script-src 'self' 'nonce-{nonce}'")
	CSP string
}

// StrictCSP is dreego's default policy: only same-origin assets plus inline
// code that carries the per-request nonce.
const StrictCSP = "default-src 'self'; " +
	"script-src 'self' 'nonce-{nonce}'; " +
	"style-src 'self' 'nonce-{nonce}'; " +
	"img-src 'self' data:; " +
	"base-uri 'none'; " +
	"form-action 'self'; " +
	"frame-ancestors 'none'; " +
	"object-src 'none'"

// CSPOff disables the Content-Security-Policy for a page or the whole app.
const CSPOff = "off"

// DefaultSecurity returns the built-in strict policy.
func DefaultSecurity() Security { return Security{CSP: StrictCSP} }

// SetSecurity sets the app-wide security policy. Pages without their own
// Security use this; without it, the strict default applies.
func (app *App) SetSecurity(s Security) { app.security = &s }

// SetCSP is a convenience for the common case: it sets only the CSP app-wide.
func (app *App) SetCSP(csp string) { app.SetSecurity(Security{CSP: csp}) }

// securityFor resolves the effective policy: page -> app -> strict default.
func (app *App) securityFor(seite Page) Security {
	if seite.Security != nil {
		return *seite.Security
	}
	if app.security != nil {
		return *app.security
	}
	return DefaultSecurity()
}

// contentSecurityPolicy renders the policy for one response: it fills the
// {nonce} placeholder and returns "" when the policy is disabled.
func (app *App) contentSecurityPolicy(seite Page, nonce string) string {
	csp := app.securityFor(seite).CSP
	if csp == "" {
		csp = StrictCSP
	}
	if csp == CSPOff {
		return ""
	}
	return strings.ReplaceAll(csp, "{nonce}", nonce)
}
