// Package wails adapts a dreego application to a Wails v3 desktop window.
//
// It is deliberately tiny and Wails-free: it returns an http.Handler that Wails
// serves in-process (application.AssetOptions.Handler). No window, no service,
// no lifecycle, and no TCP listener — the application owns all of that.
//
//	app := dreego.NewApp()
//	app.Page(page)
//
//	handler := wails.Handler(app)
//	wailsApp := application.New(application.Options{
//	    Assets: application.AssetOptions{Handler: handler},
//	})
//
// The one thing it adds over app.Handler() is a Content-Security-Policy that a
// Wails webview can actually satisfy: Wails injects its runtime and bindings as
// inline scripts, so the strict nonce-only default would block them. See
// DesktopCSP.
package wails

import (
	"net/http"

	"github.com/LukasLow/dreego/v0/pkg/dreego"
)

// DesktopCSP is the Content-Security-Policy a Wails webview needs.
//
// Wails injects its runtime (/wails/runtime.js) and the generated bindings as
// scripts; those are same-origin ('self') plus inline ('unsafe-inline'). This is
// less strict than dreego's web default (nonce-only) — that is the honest price
// of the Wails runtime. The policy still forbids framing, plugins, and foreign
// origins.
//
// If you drive everything through same-origin files and nonces, set your own
// policy with app.SetCSP(...) and ignore this constant.
const DesktopCSP = "default-src 'self'; " +
	"script-src 'self' 'unsafe-inline'; " +
	"style-src 'self' 'unsafe-inline'; " +
	"img-src 'self' data:; " +
	"font-src 'self'; " +
	"connect-src 'self'; " +
	"base-uri 'none'; " +
	"form-action 'self'; " +
	"frame-ancestors 'none'; " +
	"object-src 'none'"

// Handler returns the in-process handler for a Wails asset server.
//
// It applies DesktopCSP unless the app already declares its own policy (via
// app.SetCSP or per page), so an explicit choice is never overwritten.
func Handler(app *dreego.App) http.Handler {
	if !app.HasSecurity() {
		app.SetCSP(DesktopCSP)
	}
	return app.Handler()
}
