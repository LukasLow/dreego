# Sicherheit

dreego ist auf sichere Vorgaben ausgelegt: die wichtigen Schutzmaßnahmen sind
**an ohne Zutun**.

## Output-Safety (XSS)

- `d.Text(…)` escaped HTML automatisch.
- `d.Raw(…)` ist das bewusste Opt-out — nur für vertrauenswürdigen Inhalt.
- Komponenten-CSS/JS werden über `escapeClosingTag` entschärft (`</script>` wird
  zu `<\/script>`), damit Komponenten-JS nicht aus dem Tag ausbrechen kann.

## Content-Security-Policy + Nonce

Jede Antwort trägt eine strenge CSP:

```
default-src 'self';
script-src 'self' 'nonce-<zufall>';
style-src  'self' 'nonce-<zufall>';
img-src 'self' data:;
base-uri 'none'; form-action 'self'; frame-ancestors 'none'; object-src 'none'
```

- `'self'` erlaubt eigene statische Dateien (`/public/app.js`).
- Der **Nonce** erlaubt dreego's eigene inline `<style>`/`<script>` (Scoped
  CSS/JS, Critical CSS).
- **Inline-Style-Attribute** (`style="…"`) sind **blockiert** — deshalb nie
  verwenden (siehe [components.md](components.md)).
- Ein fehlender Nonce (crypto/rand-Ausfall) → **500**, keine Seite mit
  blockierten Assets.

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

## Bewusst offen / spätere Phase

- Ein **Audit-Log** für Security-Ereignisse fehlt noch.
- **Subresource Integrity (SRI)** für externe Assets fehlt.
- **Wails/Desktop** als eigener Adapter.
