// Command messdient startet dieselbe Federkiel-Site konfigurierbar, damit man
// Gzip an/aus und den Port frei waehlen kann (fuer Messreihen):
//
//	go run ./sites/federkiel -port 4000             # mit Gzip (Standard)
//	go run ./sites/federkiel -port 4001 -nogzip     # ohne Gzip
package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/LukasLow/dreego/v0/pkg/dreego"
	"github.com/LukasLow/dreego/v0/pkg/ratelimit"
)

func main() {
	portFlag := flag.String("port", "", "Listen-Port (ueberschreibt DREEGO_PORT, Default 4000)")
	nogzip := flag.Bool("nogzip", false, "Gzip-Kompression abschalten (fuer Vergleichsmessung)")
	flag.Parse()

	app := dreego.NewApp()
	if *nogzip {
		app.SetCompress(false)
	}

	app.Static(publicFS, "public")

	secret := os.Getenv("DREEGO_SESSION_SECRET")
	if len(secret) < 32 {
		secret = "federkiel-dev-secret-mindestens-32-zeichen-lang"
	}
	store, err_store := dreego.NewCookieStore([]byte(secret))
	if err_store != nil {
		log.Fatalf("dreego: Session-Store: %v", err_store)
	}
	app.SetSessionStore(store)

	app.Use(ratelimit.New(ratelimit.Config{Burst: 120, Every: time.Minute}).Middleware)

	Register(app)

	port := *portFlag
	if port == "" {
		port = os.Getenv("DREEGO_PORT")
	}
	if port == "" {
		port = "4000"
	}

	log.Printf("federkiel laeuft auf http://localhost:%s (gzip=%v)", port, !*nogzip)
	err_listen := http.ListenAndServe(":"+port, app.Handler())
	if err_listen != nil {
		log.Fatal(err_listen)
	}
}
