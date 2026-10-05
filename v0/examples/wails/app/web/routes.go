// Package web is the dreego frontend of the Wails app: pages rendered
// in-process and served by the Wails adapter. No .dreego files, no generator —
// pages are plain Go.
package web

import (
	"github.com/LukasLow/dreego/v0/pkg/dreego"
)

// Register adds every page to the application.
func Register(app *dreego.App) {
	app.Page(Home)
	app.Page(About)
}
