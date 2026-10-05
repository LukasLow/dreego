# dreego — Dokumentation

Willkommen. Diese `.docs/`-Sammlung erklärt dreego vollständig: was es ist, wie
es aufgebaut ist, und wie man jede Fähigkeit benutzt.

## Inhalt

| Datei | Thema |
|---|---|
| [overview.md](overview.md) | Was dreego ist, Idee, Abgrenzung zu Dreego/gomponents |
| [getting-started.md](getting-started.md) | Erstes Projekt in fünf Minuten |
| [dom.md](dom.md) | DOM-Schicht `dreego/dom`: `View`, `El`, `Attr`, htmx-Basis |
| [routing.md](routing.md) | Seiten, Methoden, dynamische `[id]`-Routen, API-Routen |
| [pages-and-layouts.md](pages-and-layouts.md) | `Page`, `Layout`, Head-Merge |
| [components.md](components.md) | Komponenten = Funktionen, `c.Box`, Scoped CSS/JS |
| [ui-library.md](ui-library.md) | `dreego/ui`: Button, Card, Badge, Alert, Field |
| [forms.md](forms.md) | `Bind[T]`, Validierung, `Old`, `Errors`, PRG |
| [sessions.md](sessions.md) | Session, Flash, CSRF |
| [security.md](security.md) | Escaping, CSP+Nonce, Header, Cookies, Gzip, Recovery |
| [errors.md](errors.md) | Eigene Fehlerseiten (HTML + JSON) |
| [i18n.md](i18n.md) | Übersetzungen, Locale-Auflösung, Plural |
| [static.md](static.md) | `go:embed`, Content-Type, Cache |
| [markdown.md](markdown.md) | Markdown-Seiten (Rechtstexte) |
| [ratelimit.md](ratelimit.md) | Rate-Limiting-Addon |
| [testing.md](testing.md) | Test-Helfer `dreego-test` |
| [architecture.md](architecture.md) | Paketaufbau, Renderfluss, Design-Regeln |
| [dependencies.md](dependencies.md) | gomponents, Lizenz, Attribution |
| [deployment.md](deployment.md) | Docker, Multi-Site, Environment |

## Grundsätze

1. **Kein Compiler.** Eine Seite ist ein Go-Wert, ein Handler eine Go-Funktion.
   Es gibt kein `.dreego`, kein `generate`.
2. **Gomponents gekapselt.** Nutzer importieren `dreego` und `dreego/dom` —
   nie `maragu.dev/gomponents` direkt.
3. **Sichere Vorgaben.** Gzip, Security-Header, CSP mit Nonce, CSRF,
   verschlüsselte Sessions: an ohne Zutun.
4. **Fehler früh und sichtbar.** Keine stillen Vorgaben; fehlende Konfiguration
   bricht den Start ab.
