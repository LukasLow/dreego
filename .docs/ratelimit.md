# Rate-Limiting (`dreego/ratelimit`)

Ein Token-Bucket-Limiter als Middleware — pro App oder pro Route.

## App-weit

```go
import (
	"time"
	"github.com/LukasLow/dreego/v0/pkg/ratelimit"
)

limiter := ratelimit.New(ratelimit.Config{Burst: 120, Every: time.Minute})
app.Use(limiter.Middleware)
```

## Pro Route

```go
var Login = dreego.Page{
	Path: "/login",
	Get:  getLogin,
	Post: postLogin,
	Middleware: []func(http.Handler) http.Handler{
		ratelimit.New(ratelimit.Config{Burst: 10, Every: time.Minute}).Middleware,
	},
}
```

## Konfiguration

| Feld | Bedeutung |
|---|---|
| `Burst` | Anzahl Anfragen, die sofort durchgehen |
| `Every` | Zeitfenster, über das der Burst nachgefüllt wird |
| `Key` | optionale Funktion: Bucket-Schlüssel (Default: Client-IP) |
| `TrustedProxies` | Adressen, die über `X-Forwarded-For` die echte IP setzen dürfen |

`Burst <= 0` oder `Every <= 0` deaktiviert das Limit.

## Verhalten

- Über dem Limit: **429** mit `Retry-After`.
- Der Bucket füllt sich **kontinuierlich** nach (kein starres Zeitfenster).
- Der Schlüssel ist die Client-IP. **Nur ein vertrauenswürdiger Proxy** darf
  `X-Forwarded-For`/`X-Real-IP` setzen — ein Client kann die IP nicht fälschen.

## Grenze

Der Zustand liegt im Speicher: das reicht für **einen** Prozess. Für mehrere
Instanzen braucht man einen geteilten Store (Redis) — außerhalb von v0.1.

## Fehlermeldung

Die 429-Antwort nutzt das Fehlerseiten-Modul (HTML im Layout bzw. JSON für
`/api`). Siehe [errors.md](errors.md).
