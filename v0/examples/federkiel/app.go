package main

import (
	"embed"

	"github.com/LukasLow/dreego/v0/pkg/dreego"
	"github.com/LukasLow/dreego/v0/pkg/urls"
)

// publicFS hält CSS, JS und das Logo (go:embed).
//
//go:embed public/*
var publicFS embed.FS

// site ist die eigene Adresse der (fiktiven) Federkiel-Seite.
// Über die Env-Variable umstellbar, damit dev/prod nur Konfiguration sind.
var site = urls.FromEnv("FEDERKIEL_URL", "http://localhost:4000")

// Register hängt alle Seiten an die App.
func Register(app *dreego.App) {
	app.Page(Home)
	app.Page(Features)
	app.Page(Pricing)
	app.Page(Login)
	app.Page(Dashboard)
	app.Page(Logout)
	app.Page(Impressum)
}
