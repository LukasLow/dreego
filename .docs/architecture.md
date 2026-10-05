# Architektur

## Paketaufbau

Das Go-Modul liegt unter **`v0/`** (Modulpfad `github.com/LukasLow/dreego/v0`).
So passen später Adapter und weitere Phasen sauber darunter.

```
github.com/LukasLow/dreego/          (Repo-Root)
├── .docs/                            Dokumentation
├── LICENSE, README.md, THIRD_PARTY_NOTICES.md
├── Dockerfile, docker-compose.yml
└── v0/                               ← Go-Modul (go.mod hier)
    ├── go.mod                        module github.com/LukasLow/dreego/v0
    ├── pkg/
    │   ├── dreego/   App, Page, Ctx, Session, CSRF, Error-Pages, Static,
    │   │             Gzip, Recovery, Redirect, i18n-Anbindung  (package dreego)
    │   ├── dom/      eigene DOM-Schicht: View, El, Attr, Text, Map + generierte Elemente/Attribute
    │   ├── scope/    Scoped CSS/JS + Collector (Dedupe, Nonce)
    │   ├── ui/       Komponenten-Bibliothek (Button, Card, Badge, Alert, Field)
    │   ├── i18n/     Übersetzungen (Katalog, T, Tn)
    │   ├── markdown/ Markdown -> Nodes
    │   ├── ratelimit/ Token-Bucket-Middleware
    │   ├── urls/     Basis-Adressen je Site
    │   └── dreego-test/  Test-Client + Assertions
    └── examples/     Demo-Sites (fensterbank, portal, www, federkiel) — main
```

Spätere Adapter (z. B. für Desktop-Hosts) liegen in `v0/plg/` daneben.

**Tag-Konvention:** Da `go.mod` in `v0/` liegt, heißt der Release-Tag
`v0/v0.0.1`. Der Import ist `github.com/LukasLow/dreego/v0/pkg/dreego`.

## Renderfluss einer Anfrage

```
HTTP-Request
  └─ Handler()
      ├─ Recover            (Panic -> 500-Seite)
      ├─ Compress           (Gzip, entschieden beim WriteHeader)
      ├─ SecurityHeaders    (nosniff, DENY, …)
      ├─ [app.Middleware…]  (app.Use, z. B. ratelimit)
      └─ dispatch
          ├─ /health        (immer verfügbar)
          └─ mux.ServeHTTP
              ├─ Seite gefunden
              │   └─ [Page.Middleware…]
              │       ├─ CSP + Nonce setzen
              │       ├─ Ctx bauen (Session laden, Locale auflösen)
              │       ├─ Methode prüfen (405)
              │       ├─ CSRF prüfen (403)
              │       ├─ Handler -> Node
              │       ├─ Session als EIN Cookie schreiben
              │       ├─ Redirect? -> Status+Location, kein Body
              │       └─ Layout(c, head, body) -> Dokument rendern
              ├─ Static gefunden -> ausliefern
              └─ sonst -> 404-Seite (Fehlermodul)
```

## Dokument-Reihenfolge

`c.Document` rendert den **Body zuerst** (damit Komponenten sich mit CSS/JS beim
Collector anmelden), dann:

```
<!DOCTYPE html><html lang><head>
  Kopf-Inhalt (Titel/Meta des Layouts + der Seite)
  Kritisches CSS (AddCritical, inline — kein weißer Blitz)
  Gesammelte Komponenten-Styles (je einmal, mit Nonce)
</head><body>
  Body-Markup
  Gesammelte Komponenten-Skripte (je einmal, mit Nonce)
</body></html>
```

## Design-Regeln (bindend)

1. **Eine Datei, ein Job** — kleine, klar benannte Dateien.
2. **Namen erklären ihren Job** — keine Ein-Buchstaben-Namen (außer Receiver,
   Laufvariablen).
3. **Fehler am Entstehungsort** — nie still verschlucken.
4. **Keine stillen Vorgaben** — fehlende Konfiguration bricht den Start ab.
5. **Kein `log.Fatal` in Lib-Paketen** — nur in `main`.
6. **Komponenten über `c.Box`** — nie `style="…"`-Attribute (CSP).
7. **Varianten als Props, nie als Fork.**

## Warum ServeMux statt eigener Router

Go 1.22 bringt Methoden- und Pfad-Patterns (`GET /x/{id}`) im Standard-mux. dreego
übersetzt nur `[id]` → `{id}` und liest `r.PathValue`. Das spart eine
Abhängigkeit und ist korrekt getestet. Ein Sonderfall: `/` wird zu `/{$}`, damit
unbekannte Pfade nicht die Startseite rendern.

## Warum gomponents gekapselt

Ein Nutzer soll mit **dreego** arbeiten, nicht mit einer fremden Bibliothek. Das
`dom`-Paket re-exportiert gomponents (MIT, siehe
[dependencies.md](dependencies.md)) unter dem dreego-Namen. Der Kern selbst
importiert gomponents direkt — das ist intern und geht Nutzer nichts an.
