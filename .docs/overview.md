# Überblick

**dreego** ist ein schlankes Go-Webframework für serverseitig gerenderte Seiten
(SSR) — ohne Compiler, ohne Template-Sprache. Eine Seite ist ein Go-Wert, ein
Handler eine Go-Funktion, HTML entsteht aus Go-Aufrufen.

Es entstand aus Dreego (`.dreego`-Dateien, die zu Go transpiliert werden). dreego
behält die Idee — Datei-Routing, Layouts, Komponenten, Sicherheit — wirft aber
den Generator weg: alles ist normales Go.

## Die Idee in einem Satz

> Eine Komponente ist eine Go-Funktion, die einen `Node` zurückgibt — mit HTML,
> CSS und JavaScript in derselben Funktion, sicher und typgeprüft.

## Bausteine

- **Dom:** HTML wird über `dreego/dom` geschrieben (Wrapper um gomponents).
- **Page:** eine Route als Wert (`Path`, `Get`, `Post`, `API`, `Layout`).
- **Layout:** eine Funktion, die Kopf und Body umschließt.
- **Ctx:** der Request-Kontext (`c.Param`, `c.T`, `c.JSON`, `c.Box`, …).
- **scope:** Scoped CSS/JS pro Komponente, dedupliziert, mit CSP-Nonce.
- **ui:** fertige Komponenten (Button, Card, Badge, Alert, Field).
- **i18n:** Übersetzungen aus Go-Maps oder JSON via `go:embed`.

## Was dreego (bewusst) NICHT hat

- Keinen `.dreego`-Transpiler und kein `generate`.
- Keine TypeScript- oder Lua-Client-Sprache.
- Keinen Desktop-Adapter — spätere Phase.
- Kein Tailwind, keine Utility-CSS-Frameworks.

## Abgrenzung

| | Dreego | gomponents | dreego |
|---|---|---|---|
| Schreibweise | `.dreego`-Dateien | Go (`Node`) | Go (`Node`) |
| Build-Schritt | Generator nötig | keiner | keiner |
| CSS/JS in der Komponente | datei-basiert | nein | ja |
| Sicherheit | eingebaut | nein | eingebaut |
| Nutzer importiert gomponents | – | direkt | nie (gekapselt) |

Siehe auch [dependencies.md](dependencies.md) zur gomponents-Lizenz (MIT).
