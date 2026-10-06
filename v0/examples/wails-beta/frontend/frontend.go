// Package frontend is the dreego site for the Wails window: the page, its layout
// and its static assets. It renders the same UI as the Wails vanilla template,
// but the markup is Go (dreego dom nodes) and the client is plain JavaScript —
// no TypeScript and no bundler.
package frontend

import (
	"embed"

	"github.com/LukasLow/dreego/v0/pkg/dreego"
)

// staticFS holds the template's static assets (style.css, logos, background
// images). They are served by dreego at /static/… — see NewApp.
//
//go:embed static
var staticFS embed.FS

// NewApp builds the dreego application and returns it. It implements
// http.Handler, so main.go hands it straight to Wails' asset server.
func NewApp() *dreego.App {
	app := dreego.NewApp()
	app.Static(staticFS, "static")
	app.Page(Index)
	return app
}
