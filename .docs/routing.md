# Routing

## Eine Route ist ein Wert

```go
var Pricing = dreego.Page{
	Path: "/pricing",
	Get:  getPricing,
}

func getPricing(c *dreego.Ctx) d.View {
	return d.H1(d.Text("Preise"))
}
```

Registriert wird zentral:

```go
app := dreego.NewApp()
app.Page(Pricing)
```

Es gibt **kein** datei-basiertes Routing: der Dateiname ist bedeutungslos, nur
`Path` zählt.

## HTTP-Methoden

```go
var Login = dreego.Page{
	Path: "/login",
	Get:  getLogin,   // zeigt das Formular
	Post: postLogin,  // verarbeitet es (PRG)
}
```

Eine nicht erlaubte Methode ergibt `405` mit `Allow`-Header — als HTML-Seite
(siehe [errors.md](errors.md)), nicht als nackter Text.

## Dynamische Segmente `[id]`

Ein Segment in eckigen Klammern wird zum Pfad-Parameter:

```go
var Project = dreego.Page{
	Path: "/projekte/[id]",
	Get: func(c *dreego.Ctx) d.View {
		id := c.Param("id")          // "abc-123"
		return d.H1(d.Text("Projekt " + id))
	},
}
```

Mehrere Segmente:

```go
Path: "/projekte/[id]/aufgaben/[aufgabe]"
c.Param("id")       // erstes
c.Param("aufgabe")  // zweites
```

Catch-all (Rest des Pfads):

```go
Path: "/dateien/[...pfad]"
c.Param("pfad")     // "a/b/c.txt"
```

Innere Mechanik: `[id]` wird zu Go-1.22-`{id}`, `[...rest]` zu `{rest...}`. Der
vorgelagerte `http.ServeMux` extrahiert die Werte, dreego liest sie über
`r.PathValue`.

## Feste und dynamische Routen nebeneinander

```go
app.Page(dreego.Page{Path: "/projekte/neu", Get: getNew})   // gewinnt bei /projekte/neu
app.Page(dreego.Page{Path: "/projekte/[id]", Get: getByID}) // sonst
```

Der spezifischere feste Pfad hat Vorrang.

## Root-Sonderfall

`Path: "/"` wird intern zu `/{$}` — sonst würde der Go-Mux `/` als Catch-all
behandeln und jede unbekannte URL die Startseite rendern. Unbekannte Pfade
ergeben so korrekt `404`.

## API-Routen

Statt HTML kann eine Route JSON/XML liefern:

```go
var Ping = dreego.Page{
	Path: "/api/ping",
	API: func(c *dreego.Ctx) error {
		return c.JSON(200, map[string]any{"pong": true})
	},
}
```

`Page.API` rendert **kein** Layout. Siehe [forms.md](forms.md) für `c.Bind` und
[security.md](security.md) für CSRF auf API-POSTs.

## Mehrere Sites (mehrere Ports)

Wie Statuna (www/app/link): jede Site ist eine eigene `App` mit eigenem Handler.
Cross-Site-Links baut man über `dreego/urls`:

```go
var public = urls.FromEnv("PUBLIC_URL", "http://localhost:4000")
h.A(h.Href(public.Abs("/login")), d.Text("Anmelden"))
```

Beispiel: `v0/examples/fensterbank` (:4000) und `v0/examples/portal` (:4001).
