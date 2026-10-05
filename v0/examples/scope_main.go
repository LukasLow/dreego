// Command demo zeigt die dreego-Addons (scope: CSS + JS, Dedupe, CSP-Nonce).
//
// Ausfuehren:
//
//	go run ./example            # Server auf http://localhost:8080
//	go run ./example -html      # HTML nach stdout (ohne Server)
package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	. "github.com/LukasLow/dreego/v0/pkg/dom"
	g "github.com/LukasLow/dreego/v0/pkg/dom"

	"github.com/LukasLow/dreego/v0/pkg/scope"
)

// hero — CSS neben dem HTML.
func hero(c *scope.Collector) g.View {
	return c.Box(
		scope.CSS(`
.hero { text-align: center; padding: 40px 0 }
.hero h1 { font-size: 40px; margin: 0 0 12px; color: #111 }
.hero p { font-size: 18px; color: #333 }
@media (max-width: 640px) { .hero h1 { font-size: 30px } }
`),
		Section(Class("hero"),
			H1(g.Text("Ist alles bereit?")),
			P(g.Text("Scoped CSS kommt direkt aus der Komponente.")),
		),
	)
}

// karte — dieselbe Klasse .hero, bleibt aber unberuehrt (beweist Scoping).
func karte(c *scope.Collector) g.View {
	return c.Box(
		scope.CSS(`
.hero { text-align: left; padding: 12px; border: 2px dashed #0080ff }
`),
		Section(Class("hero"),
			P(g.Text("Diese Karte nutzt dieselbe Klasse .hero — bleibt aber unberuehrt.")),
		),
	)
}

// zaehler — HTML + CSS + JS in EINER Funktion. Dank Collector wird das CSS nur
// einmal ausgegeben, auch wenn die Komponente zweimal auf der Seite steht.
func zaehler(c *scope.Collector) g.View {
	return c.Box(
		scope.CSS(`
.zaehler { display: inline-flex; align-items: center; gap: 14px;
           border: 2px solid #0080ff; border-radius: 12px; padding: 12px 16px }
.zaehler .wert { font-size: 22px; font-weight: 700; min-width: 36px; text-align: center }
.zaehler button { font: inherit; font-weight: 700; padding: 6px 14px; border: none;
                  border-radius: 8px; background: #0080ff; color: #000; cursor: pointer }
`),
		scope.JS(`
var plus = root.querySelector('[data-tat="plus"]');
var minus = root.querySelector('[data-tat="minus"]');
var wert = root.querySelector('.wert');

plus.addEventListener('click', function () {
    wert.textContent = String(Number(wert.textContent) + 1);
});
minus.addEventListener('click', function () {
    wert.textContent = String(Number(wert.textContent) - 1);
});
`),
		Div(Class("zaehler"),
			Button(g.Attr("data-tat", "minus"), g.Text("−")),
			Span(Class("wert"), g.Text("0")),
			Button(g.Attr("data-tat", "plus"), g.Text("+")),
		),
	)
}

func seite(c *scope.Collector) g.View {
	body := Body(
		hero(c),
		karte(c),
		Section(H2(g.Text("Zähler (JS in derselben Komponente)")), zaehler(c)),
		// Zweite Instanz: das CSS wird trotzdem nur EINMAL ausgegeben.
		zaehler(c),
	)

	head := Head(
		Meta(Charset("utf-8")),
		Meta(Name("viewport"), Content("width=device-width, initial-scale=1")),
		TitleEl(g.Text("dreego — scope: CSS + JS")),
	)

	return scope.Document(c, "de", head, body)
}

func main() {
	alsHTML := flag.Bool("html", false, "Seite als HTML nach stdout schreiben und beenden")
	flag.Parse()

	if *alsHTML {
		sammler := scope.New()
		sammler.SetNonce("demo-nonce")
		err_render := seite(sammler).Render(os.Stdout)
		if err_render != nil {
			log.Fatal(err_render)
		}
		return
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		nonce := scope.NewNonce()

		// Eine echte CSP: nur Skripte/Stile mit diesem Nonce duerfen laufen.
		w.Header().Set("Content-Security-Policy",
			"default-src 'self'; script-src 'nonce-"+nonce+"'; style-src 'nonce-"+nonce+"'")
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		sammler := scope.New()
		sammler.SetNonce(nonce)

		err_render := seite(sammler).Render(w)
		if err_render != nil {
			log.Printf("render: %v", err_render)
		}
	})

	adresse := "localhost:8080"
	log.Printf("dreego demo laeuft auf http://%s", adresse)
	err_listen := http.ListenAndServe(adresse, nil)
	if err_listen != nil {
		log.Fatal(err_listen)
	}
}
