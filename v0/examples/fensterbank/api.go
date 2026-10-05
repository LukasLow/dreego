package main

import (
	"github.com/LukasLow/dreego/v0/pkg/dreego"
)

// RegisterAPI hängt die Nicht-HTML-Endpunkte an die App. Sie liegen neben den
// HTML-Seiten im selben Mux — kein Konflikt, solange die Pfade unterschiedlich
// sind.
func RegisterAPI(app *dreego.App) {
	app.Page(Ping)
	app.Page(SortenJSON)
	app.Page(SorteJSON)
}

// Ping — einfacher JSON-Endpunkt.
var Ping = dreego.Page{
	Path: "/api/ping",
	API: func(c *dreego.Ctx) error {
		return c.JSON(200, map[string]any{"pong": true, "app": "fensterbank"})
	},
}

// SortenJSON — Liste als JSON.
var SortenJSON = dreego.Page{
	Path: "/api/sorten",
	API: func(c *dreego.Ctx) error {
		return c.JSON(200, sortenListe)
	},
}

// SorteJSON — DYNAMISCHE API-Route /api/sorten/[id].
var SorteJSON = dreego.Page{
	Path: "/api/sorten/[id]",
	API: func(c *dreego.Ctx) error {
		s, gefunden := findeSorte(c.Param("id"))
		if !gefunden {
			return c.JSON(404, map[string]string{"error": "nicht gefunden"})
		}
		return c.JSON(200, s)
	},
}
