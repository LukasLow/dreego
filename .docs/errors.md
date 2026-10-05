# Fehlerseiten

dreego rendert Fehler als **richtige Seiten** — im Layout der Site, nicht als
nackter Text. Für API-Wege kommt JSON.

## Eigene Seite registrieren

```go
app.SetErrorPage(404, func(c *dreego.Ctx) d.View {
	return c.Document("de",
		d.TitleEl(d.Text("Nicht gefunden")),
		fehlerInhalt(c, 404, "Nicht gefunden", "Diese Seite gibt es nicht."),
	)
})
app.SetErrorPage(500, func(c *dreego.Ctx) d.View { … })
```

Pro Status (404, 403, 405, 429, 500, …) ein Handler. Er bekommt einen voll
funktionsfähigen `Ctx` (mit Nonce) und gibt das **komplette Dokument** zurück —
so kann er dasselbe Layout wie die Seite benutzen.

## Ohne eigene Seite

Rendert dreego eine minimalistische, aber korrekte HTML-Seite (Titel, Status,
Link nach Hause) — **nie** der nackte Go-Text, nie eine leere Antwort.

## JSON für API

Bei `/api/…` oder `Accept: application/json` liefert dreego JSON statt HTML:

```json
{"error":"Not Found","status":404}
```

## Betroffene Stellen

Alle früheren `http.Error`/`http.NotFound`-Aufrufe sind ersetzt:

- unbekannte Route → 404
- falsche Methode → 405 (mit `Allow`)
- CSRF-Fehler → 403
- fehlender Store bei aktivem CSRF → 500
- Panic → 500
- Static nicht gefunden → 404

## 500 zeigt nie Details

Eine `500`-Seite enthält **keine** internen Fehlermeldungen — nur einen
generischen Text. Die Ursache steht im Server-Log. (Per Test abgesichert.)

## Beispiel

- `v0/examples/fensterbank/errors.go` — 404/500 im Fensterbank-Look (übersetzt).
- `v0/examples/portal/errors.go` — dasselbe im Portal-Look.
