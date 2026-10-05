// Command portal ist die zweite eigenständige Site der Demo-Landschaft:
// Anmeldung und Konto. Sie läuft auf einem EIGENEN Port (4001) und verlinkt
// über den urls-Helfer auf die öffentliche Fensterbank-Site (4000) — genau wie
// Statuna www und app getrennt betreibt.
//
//	go run ./sites/portal -port 4001
package main

import (
	"embed"
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/LukasLow/dreego/v0/pkg/dreego"
	"github.com/LukasLow/dreego/v0/pkg/ratelimit"
	"github.com/LukasLow/dreego/v0/pkg/urls"
)

//go:embed public/*
var publicFS embed.FS

// public ist die Adresse der öffentlichen Site (Fensterbank).
// markt ist die eigene Adresse dieses Portals. Beide über Env umstellbar.
var (
	public = urls.FromEnv("FENSTERBANK_URL", "http://localhost:4000")
	markt  = urls.FromEnv("PORTAL_URL", "http://localhost:4001")
)

func main() {
	portFlag := flag.String("port", "", "Listen-Port (Default 4001)")
	flag.Parse()

	app := dreego.NewApp()
	app.Static(publicFS, "public")

	secret := os.Getenv("DREEGO_SESSION_SECRET")
	if len(secret) < 32 {
		secret = "portal-dev-secret-mindestens-32-zeichen-la"
	}
	store, err_store := dreego.NewCookieStore([]byte(secret))
	if err_store != nil {
		log.Fatalf("dreego: Session-Store: %v", err_store)
	}
	app.SetSessionStore(store)

	app.Use(ratelimit.New(ratelimit.Config{Burst: 120, Every: time.Minute}).Middleware)

	RegistriereFehlerseiten(app)
	Register(app)

	port := *portFlag
	if port == "" {
		port = os.Getenv("DREEGO_PORT")
	}
	if port == "" {
		port = "4001"
	}

	log.Printf("portal laeuft auf http://localhost:%s", port)
	err_listen := http.ListenAndServe(":"+port, app.Handler())
	if err_listen != nil {
		log.Fatal(err_listen)
	}
}
