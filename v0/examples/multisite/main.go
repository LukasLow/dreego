// Command multisite is a demo of several apps in ONE process, each on its own
// port — the shape Statuna uses for www / app / link (separate subdomains,
// separate session cookies, cross-links by absolute URL).
//
// Start returns at once; the servers run in the background. Wait blocks and
// collects the errors, so a crash on one port is reported instead of aborting
// the process silently.
//
//	go run ./v0/examples/multisite
//	# http://localhost:4100  öffentliche Site
//	# http://localhost:4101  Konto-Site
package main

import (
	"log"
	"os"

	d "github.com/LukasLow/dreego/v0/pkg/dom"
	"github.com/LukasLow/dreego/v0/pkg/dreego"
	"github.com/LukasLow/dreego/v0/pkg/urls"
)

// public und konto sind die beiden Site-Adressen. Dev und prod unterscheiden
// sich nur über die Umgebungsvariable (siehe urls.FromEnv).
var (
	public = urls.FromEnv("PUBLIC_URL", "http://localhost:4100")
	konto  = urls.FromEnv("KONTO_URL", "http://localhost:4101")
)

// publicApp zeigt die öffentlichen Seiten. Ohne Formulare → CSRF aus.
func publicApp() *dreego.App {
	app := dreego.NewApp()
	app.SetCSRF(false)

	app.Page(dreego.Page{
		Path: "/",
		Get: func(c *dreego.Ctx) d.View {
			return c.Document("de", d.TitleEl(d.Text("Öffentlich")),
				d.Body(
					d.H1(d.Text("Öffentliche Site")),
					d.P(d.Text("Diese Site läuft auf ihrem eigenen Port.")),
					d.P(d.A(d.Href(konto.Abs("/konto")), d.Text("Zum Konto (andere Site)"))),
				),
			)
		},
	})
	return app
}

// kontoApp zeigt die Konto-Seite. Eigener Handler, eigenes Session-Cookie.
func kontoApp() *dreego.App {
	app := dreego.NewApp()

	app.Page(dreego.Page{
		Path: "/konto",
		Get: func(c *dreego.Ctx) d.View {
			return c.Document("de", d.TitleEl(d.Text("Konto")),
				d.Body(
					d.H1(d.Text("Konto-Site")),
					d.P(d.Text("Anderer Port, andere Session — wie app.statuna.com.")),
					d.P(d.A(d.Href(public.Abs("/")), d.Text("Zur öffentlichen Site"))),
				),
			)
		},
	})
	return app
}

func main() {
	// Ohne Store wäre Session/CSRF inert. Das Secret muss gesetzt sein; im Dev
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
	store.SetCookiePolicy(dreego.CookiePolicy{MaxAge: dreego.Days(30)})

	public_app := publicApp()
	konto_app := kontoApp()
	konto_app.SetSessionStore(store)

	// Beide Server starten sofort und laufen im Hintergrund.
	public_server := dreego.Start(":4100", public_app.Handler())
	konto_server := dreego.Start(":4101", konto_app.Handler())

	log.Printf("öffentlich auf %s, konto auf %s", public.Base(), konto.Base())

	// Wait blockiert und meldet den ersten Fehler.
	if err_public := public_server.Wait(); err_public != nil {
		log.Fatalf("public-Server beendet: %v", err_public)
	}
	if err_konto := konto_server.Wait(); err_konto != nil {
		log.Fatalf("konto-Server beendet: %v", err_konto)
	}
}
