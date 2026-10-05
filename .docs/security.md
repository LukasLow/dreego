# Sicherheit

dreego ist auf sichere Vorgaben ausgelegt: die wichtigen Schutzmaßnahmen sind
**an ohne Zutun**.

## Output-Safety (XSS)

- `d.Text(…)` escaped HTML automatisch.
- `d.Raw(…)` ist das bewusste Opt-out — nur für vertrauenswürdigen Inhalt.
- Komponenten-CSS/JS werden über `escapeClosingTag` entschärft (`</script>` wird
  zu `<\/script>`), damit Komponenten-JS nicht aus dem Tag ausbrechen kann.

## Content-Security-Policy + Nonce

Jede Antwort trägt standardmäßig eine strenge CSP. Sie ist **deklarativ** und
**überschreibbar** — wie ein `Page`-Feld, in drei Ebenen:

**1. Eingebauter Default** (streng, mit Nonce):

```
default-src 'self';
script-src 'self' 'nonce-<zufall>';
style-src  'self' 'nonce-<zufall>';
img-src 'self' data:;
base-uri 'none'; form-action 'self'; frame-ancestors 'none'; object-src 'none'
```

**2. App-weit** über `app.SetCSP(...)` (oder `app.SetSecurity(Security{CSP: …})`):

```go
app.SetCSP("default-src 'self'; script-src 'self' 'nonce-{nonce}' https://cdn.example")
```

**3. Pro Seite** über das `Security`-Feld an `Page`:

```go
var Report = dreego.Page{
	Path: "/report",
	Security: &dreego.Security{
		CSP: "default-src 'self'; script-src 'self' 'nonce-{nonce}' https://trusted.example",
	},
	Get: getReport,
}
```

**Platzhalter:** `{nonce}` wird pro Request durch den echten Nonce ersetzt.

**Abschalten:** `CSPOff` ("off") auf App- oder Seitenebene.

**Reihenfolge:** Seite → App → Default. So bleibt der Default sicher, aber jede
echte Seite kann ihre Policy deklarieren (z. B. eine Seite mit externem Widget).

Ein fehlender Nonce (crypto/rand-Ausfall) → **500**, keine Seite mit
blockierten Assets.

Für Inline-Style-**Attribute** (`style="…"`) gilt weiter: vermeiden — sie werden
von `style-src 'nonce-…'` blockiert. Layout in eine gescopte CSS-Klasse.

## Security-Header (immer an)

| Header | Wert |
|---|---|
| `X-Content-Type-Options` | `nosniff` |
| `X-Frame-Options` | `DENY` |
| `Referrer-Policy` | `strict-origin-when-cross-origin` |
| `Permissions-Policy` | `geolocation=(), microphone=(), camera=()` |

## CSRF

Synchronizer-Token in der Session, automatisch geprüft. Details:
[sessions.md](sessions.md).

## Session-Cookies

Verschlüsselt (AES-256-GCM), HttpOnly, SameSite=Lax, Secure bei TLS. Größenlimit
4096 Bytes.

## Gzip

Textantworten (HTML, CSS, JS, SVG, JSON) schrumpfen auf ~⅓. Bilder, Fonts und
bereits Komprimiertes werden **nicht** gzip't. `Vary: Accept-Encoding` wird
immer gesetzt. Die Entscheidung fällt beim `WriteHeader` — sonst bekämen
Antworten mit vorzeitigem Header (JSON, Fehler) einen komprimierten Body **ohne**
`Content-Encoding` und der Browser zeigte Binärmüll.

Abschaltbar für Vergleiche: `app.SetCompress(false)`.

## Rate-Limiting

Als Addon: `dreego/ratelimit`. Siehe [ratelimit.md](ratelimit.md).

## Recovery

Ein Panic in einem Handler wird zu einer generischen `500`-Seite; die Ursache
geht nur ins Server-Log, nie an den Client.

## Fehler nie still

- Fehlende Konfiguration bricht den Start ab (`log.Fatal` nur in `main`).
- Session-Schreibfehler → gemeldete `500`, kein stiller Logout.
- Keine stillen Standardwerte.

## Was ein Nutzer abschalten/ändern kann

| Sicherheits-Ding | Steuerung |
|---|---|
| CSRF | `app.SetCSRF(false)` |
| Gzip | `app.SetCompress(false)` |
| Session-Cookies | keinen Store setzen |
| **CSP** | `app.SetCSP(...)` app-weit, `Page.Security` pro Seite, `CSPOff` aus |
| Security-Header (nosniff, DENY, …) | noch **fest** (spätere Option) |

Sichere Vorgaben sind **Defaults**, keine Fesseln: jede lässt sich bewusst
ändern — mit klarer Ansage im Code.

## Bewusst offen / spätere Phase

- Ein **Audit-Log** für Security-Ereignisse fehlt noch.
- **Subresource Integrity (SRI)** für externe Assets fehlt.
- **Wails/Desktop** als eigener Adapter (braucht die CSP-Steuerung, die es nun gibt).
