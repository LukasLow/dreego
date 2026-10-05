// Command fensterbank ist eine Multi-Seiten-Demo auf Basis von dreego.
//
// Zeigt: viele feste Seiten + dynamische [id]-Routen + eine API-Route, alles in
// EINER App auf EINEM Port — sie geraten sich nicht in die Quere.
//
//	go run ./sites/fensterbank -port 4000           # viele HTML-Seiten
//	go run ./sites/fensterbank -port 4001 -apiOnly  # nur die API-Seiten
package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/LukasLow/dreego/v0/pkg/dreego"
	"github.com/LukasLow/dreego/v0/pkg/i18n"
	"github.com/LukasLow/dreego/v0/pkg/ratelimit"
)

func main() {
	portFlag := flag.String("port", "", "Listen-Port (Default 4000)")
	apiOnly := flag.Bool("apiOnly", false, "nur die API-Routen registrieren")
	flag.Parse()

	app := dreego.NewApp()
	app.Static(publicFS, "public")

	secret := os.Getenv("DREEGO_SESSION_SECRET")
	if len(secret) < 32 {
		secret = "fensterbank-dev-secret-mindestens-32-zeichen"
	}
	store, err_store := dreego.NewCookieStore([]byte(secret))
	if err_store != nil {
		log.Fatalf("dreego: Session-Store: %v", err_store)
	}
	app.SetSessionStore(store)

	app.Use(ratelimit.New(ratelimit.Config{Burst: 120, Every: time.Minute}).Middleware)

	// i18n: Kataloge aus dem Embed laden (kein Buildstep).
	bundle, err_i18n := i18n.Open(localeFS, "locales", "de")
	if err_i18n != nil {
		log.Fatalf("fensterbank: i18n: %v", err_i18n)
	}
	app.SetI18n(bundle)
	app.SetDefaultLocale("de")

	RegistriereFehlerseiten(app)

	if *apiOnly {
		RegisterAPI(app)
	} else {
		RegisterPages(app)
		RegisterAPI(app)
	}

	port := *portFlag
	if port == "" {
		port = os.Getenv("DREEGO_PORT")
	}
	if port == "" {
		port = "4000"
	}

	log.Printf("fensterbank laeuft auf http://localhost:%s (apiOnly=%v)", port, *apiOnly)
	err_listen := http.ListenAndServe(":"+port, app.Handler())
	if err_listen != nil {
		log.Fatal(err_listen)
	}
}
