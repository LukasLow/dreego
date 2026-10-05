// Package main — Demo-Site für dreego.
//
//	cd sites/www && go run .      # http://localhost:8080
//
// Zeigt: Page (Route), Layout, Ctx, in (Form-Binding), Redirect, Session, CSRF,
// Flash, Static (go:embed), urls und scoped CSS/JS.
package main

import (
	"log"
	"net/http"
	"os"

	"github.com/LukasLow/dreego/v0/pkg/dreego"
)

// Register hängt alle Seiten an die App. Die URL kommt nur aus Page.Path —
// der Dateiname ist bedeutungslos (kein file-based Routing).
func Register(app *dreego.App) {
	app.Page(Home)
	app.Page(Login)
	app.Page(Dashboard)
	app.Page(Logout)
	app.Page(Legal)
}

func main() {
	app := dreego.NewApp()

	// Static: eingebettete Dateien unter /public/*.
	app.Static(publicFS, "public")

	// Session: verschlüsseltes Cookie. Ohne Store wären Session/Flash/CSRF inert.
	// Das Secret MUSS gesetzt sein — ein fest eingebauter Fallback wäre ein
	// öffentlich bekanntes AES-Passwort (Session-Fälschung). Im Dev-Betrieb
	// erlaubt DREEGO_DEV=1 ein bewusst unsicheres Standardgeheimnis.
	secret := os.Getenv("DREEGO_SESSION_SECRET")
	if len(secret) < 32 {
		if os.Getenv("DREEGO_DEV") != "1" {
			log.Fatal("dreego: DREEGO_SESSION_SECRET fehlt oder ist kürzer als 32 Zeichen " +
				"(für Dev bewusst DREEGO_DEV=1 setzen)")
		}
		secret = "nur-fuer-dev-nicht-in-produktion-verwenden-32b"
	}
	store, err_store := dreego.NewCookieStore([]byte(secret))
	if err_store != nil {
		log.Fatalf("dreego: Session-Store: %v", err_store)
	}
	app.SetSessionStore(store)

	Register(app)

	adresse := "localhost:8080"
	log.Printf("dreego www laeuft auf http://%s", adresse)
	err_listen := http.ListenAndServe(adresse, app.Handler())
	if err_listen != nil {
		log.Fatal(err_listen)
	}
}
