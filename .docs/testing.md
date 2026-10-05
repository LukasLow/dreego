# Test-Helfer `dreego-test`

Testet eine dreego-App **ohne Netzwerk**: der Client stellt Anfragen direkt
gegen den Handler, merkt sich Cookies (Session/CSRF) und liefert ein `Response`
mit Assertions.

```go
import dreegotest "github.com/LukasLow/dreego/v0/pkg/dreego-test"

func TestStartseite(t *testing.T) {
	app := dreego.NewApp()
	app.Page(Home)

	client := dreegotest.New(app.Handler())
	resp := client.Get("/")

	resp.AssertStatus(t, 200)
	resp.AssertContains(t, "Hallo")
	resp.AssertNotContains(t, "Fehler")
}
```

## Client-Methoden

| Methode | Wirkung |
|---|---|
| `client.Get(pfad)` | GET, folgt **nicht** automatisch Redirects |
| `client.PostForm(pfad, map[string]string{…})` | POST urlencoded |
| `client.PostJSON(pfad, jsonString)` | POST mit JSON-Body |

Cookies werden zwischen den Aufrufen gehalten — so funktionieren Session und
CSRF wie im Browser (GET holt das Token, POST trägt es).

## Response-Assertions

| Methode | prüft |
|---|---|
| `resp.AssertStatus(t, 200)` | Statuscode |
| `resp.AssertContains(t, "…")` | Body enthält Text |
| `resp.AssertNotContains(t, "…")` | Body enthält Text nicht |
| `resp.AssertHeader(t, "Content-Type", "…")` | genau ein Header |
| `resp.AssertLocation(t, "/")` | `Location`-Header |

## Felder

`resp.Status`, `resp.Header` (`http.Header`), `resp.Body` (String).

## Warum kein Netzwerk

`httptest` genügt: der Handler wird direkt aufgerufen. Das ist schnell,
deterministisch und braucht keinen Port. Für Cookies sorgt ein `cookiejar`.

## Kein eigenes Modul

`dreego-test` liegt unter **`v0/pkg/dreego-test`** und gehört zum selben Modul
wie dreego — **kein** separater Tag. Es wird also mit `v0/v0.x.y` zusammen
versioniert.
