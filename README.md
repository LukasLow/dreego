# dreego

Ein schlankes Go-Webframework für serverseitig gerenderte Seiten. Kein Compiler,
keine Template-Sprache — eine Seite ist ein Go-Wert, eine Komponente eine
Go-Funktion, CSS und JS stehen in derselben Funktion.

```go
package main

import (
	"log"
	"net/http"

	"github.com/LukasLow/dreego/v0/pkg/dreego"
	d "github.com/LukasLow/dreego/v0/pkg/dom"
)

var Start = dreego.Page{
	Path: "/",
	Get: func(c *dreego.Ctx) d.Node {
		return c.Box(
			dreego.CSS(`h1 { color: #0080ff }`),
			d.H1(d.Text("Hallo dreego")),
		)
	},
}

func main() {
	app := dreego.NewApp()
	app.Page(Start)

	err := http.ListenAndServe(":8080", app.Handler())
	if err != nil {
		log.Fatal(err)
	}
}
```

Kein Build-Schritt. `go run .` und die Seite läuft.

## Philosophie

1. **Reines Go.** Routing, Layouts, Komponenten, Formulare, Sessions, CSRF,
   i18n, Markdown, Static — alles Go-Standardbibliothek und ein gekapseltes
   Darstellungs-Paket. Kein `.dreego`, kein `generate`.
2. **Sichere Vorgaben.** Gzip, Security-Header, CSP mit Nonce, CSRF, AES-
   verschlüsselte Sessions: an ohne Zutun.
3. **Lesbar für Menschen und Vorlese-Tools.** `d.H1(d.Text("…"))` statt
   `<h1>…</h1>`; der Compiler prüft, `gofmt` räumt auf.

## Was drin ist

| Bereich | Fähigkeit |
|---|---|
| Routing | feste + dynamische `[id]`-Routen, Methoden, Catch-all |
| Seiten | `Page`, `Layout`, Head-Merge, Critical CSS |
| Komponenten | HTML + CSS + JS in einer Funktion, gescopet, dedupliziert |
| UI-Bibliothek | Button, Card, Badge, Alert, Field |
| Formulare | `Bind[T]`, Validierung, `Old`, `Errors`, PRG |
| Sicherheit | Session (AES-256-GCM), CSRF, CSP+Nonce, Header, Gzip, Recovery |
| Fehlerseiten | eigene 404/500 als Seite oder JSON |
| API-Routen | `c.JSON`, `c.XML`, `c.Write`, `c.Bind` |
| i18n | Kataloge aus Go/JSON, `Accept-Language`/Cookie, Plural |
| Statisch | `go:embed`, eigene MIME-Map |
| Markdown | Rechtstexte, sicher (kein Roh-HTML) |
| Rate-Limiting | Token-Bucket-Addon |
| Test-Helfer | `dreego-test`: Client + Assertions ohne Netzwerk |

## Aufbau

Das Go-Modul liegt unter **`v0/`**:

```
dreego/                            Repo-Root
├── .docs/                         ausführliche Dokumentation
├── LICENSE, README.md, THIRD_PARTY_NOTICES.md
└── v0/                            Go-Modul (github.com/LukasLow/dreego/v0)
    ├── pkg/dreego/                Framework (package dreego)
    ├── pkg/dom/                   re-exportiert gomponents — nie direkt importieren
    ├── pkg/{scope,ui,i18n,markdown,ratelimit,urls}/
    ├── pkg/dreego-test/            Test-Client + Assertions
    └── examples/                  Demo-Sites (fensterbank :4000, portal :4001)
```

## Dokumentation

Ausführlich unter [`./.docs/`](./.docs/index.md): Überblick, erste Schritte,
Routing, Seiten/Layouts, Komponenten, UI-Bibliothek, Formulare, Sessions,
Sicherheit, Fehlerseiten, i18n, Static, Markdown, Rate-Limiting, Architektur,
Abhängigkeiten, Deployment.

## Beispiele starten

```sh
cd v0/examples/fensterbank && DREEGO_DEV=1 go run . -port 4000
cd v0/examples/portal       && DREEGO_DEV=1 go run . -port 4001
```

Oder per Docker:

```sh
docker compose up -d --build
# http://localhost:4000  öffentliche Seite
# http://localhost:4001  Anmeldung/Konto
```

## Lizenz

**MPL-2.0** (Mozilla Public License 2.0, siehe [LICENSE](LICENSE)) — schwaches
Copyleft: dreego selbst (und Änderungen daran) bleibt offen, aber deine
Anwendung, die dreego *benutzt*, muss **nicht** offen sein.

dreego baut auf **gomponents** (https://github.com/maragudk/gomponents,
**MIT**, gepinnt `v1.3.0`) — die Attribution steht in
[THIRD_PARTY_NOTICES.md](THIRD_PARTY_NOTICES.md).
