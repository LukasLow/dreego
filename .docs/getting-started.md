# Erste Schritte

## Voraussetzungen

Go 1.22 oder neuer (`go version`). dreego nutzt `http.ServeMux`-Pattern wie
`{id}` und `r.PathValue`, die seit Go 1.22 dabei sind.

## Projekt anlegen

```
mkdir meine-app && cd meine-app
go mod init beispiel.de/meine-app
go get github.com/LukasLow/dreego/v0
```

## Eine erste Seite

`main.go`:

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
	Get: func(c *dreego.Ctx) d.View {
		return c.Box(
			dreego.CSS(`h1 { color: #0080ff }`),
			d.H1(d.Text("Hallo dreego")),
		)
	},
}

func main() {
	app := dreego.NewApp()
	app.Page(Start)

	log.Println("http://localhost:8080")
	err_listen := http.ListenAndServe(":8080", app.Handler())
	if err_listen != nil {
		log.Fatal(err_listen)
	}
}
```

```
go run .
# http://localhost:8080
```

Das war es. Kein Generator, kein Build-Schritt, keine Template-Datei.

## Was passiert ist

- `dreego.Page` beschreibt die Route `/` (nur `Path` bestimmt die URL).
- `c.Box` meldet eine Komponente an: das CSS wird gehasht, gescopet und pro
  Seite einmal inline in `<head>` ausgegeben.
- `app.Handler()` liefert den fertigen `http.Handler` mit Sicherheits-Headern,
  Gzip und Recovery.

## Nächste Schritte

- [routing.md](routing.md) — mehrere Seiten, Methoden, `[id]`-Routen.
- [components.md](components.md) — eigene Komponenten mit CSS und JS.
- [ui-library.md](ui-library.md) — fertige Bausteine benutzen.
- [deployment.md](deployment.md) — als Container starten.
