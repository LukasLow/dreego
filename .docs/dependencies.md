# Abhängigkeiten und Lizenz

## dreego selbst

dreego steht unter der **MPL-2.0** (Mozilla Public License 2.0, siehe `LICENSE`).
Schwaches Copyleft: Änderungen **an dreego selbst** müssen offen bleiben, deine
Anwendung, die dreego *benutzt*, aber **nicht**.

**MPL-2.0 + MIT ist kompatibel:** Der gomponents-Teil bleibt MIT, dreegos MPL
gilt nur für dreegos eigene Dateien. Der MIT-Hinweis muss erhalten bleiben (siehe
`THIRD_PARTY_NOTICES.md`).

## gomponents

dreego baut auf **gomponents** auf.

| | |
|---|---|
| Kanonischer Modulpfad | `maragu.dev/gomponents` (Go verlangt ihn in Imports) |
| Echte Quelle (GitHub) | `github.com/maragudk/gomponents`, Tag **`v1.3.0`** |
| Lizenz | **MIT** — Copyright (c) **Maragu ApS** |

dreego pinnt die Quelle explizit in `v0/go.mod`:

```
require maragu.dev/gomponents v1.3.0
replace maragu.dev/gomponents => github.com/maragudk/gomponents v1.3.0
```

**Warum beides?** Go lehnt ein Modul ab, dessen deklarierter Pfad vom
`require` abweicht (`module declares its path as: maragu.dev/gomponents`). Der
`require` muss also den kanonischen Pfad nennen — der `replace` zeigt den Build
auf das echte GitHub-Repo und den Release-Tag.

**Ehrliche Einschränkung:** `replace`-Zeilen gelten nur für das **Hauptmodul**.
Für einen **Nutzer**, der dreego importiert, wird `maragu.dev/gomponents`
weiterhin über den kanonischen Pfad aufgelöst — der seinerseits auf GitHub
zeigt (Vanity-Pfad). Wer den GitHub-Pfad *erzwingen* will, müsste gomponents
forken und dessen Modulpfad ändern; das wäre unüblich.

**Was MIT bedeutet für dich:**
- ✅ Verwenden, ändern, weitergeben, verkaufen.
- ✅ Kompatibel mit dreegos MPL-2.0 (schwaches Copyleft).
- ✅ **Keine Copyleft-Pflicht** — deine Anwendung muss **nicht** offen sein.
- ⚠️ Der Copyright-Hinweis und Lizenztext müssen bleiben.

Der vollständige Lizenztext von gomponents ist in **`THIRD_PARTY_NOTICES.md`**
reproduziert (das verlangt die Lizenz).

## Warum gomponents nicht geforkt wird

dreego **forkt gomponents nicht** als ganzes Repo. Stattdessen übernimmt dreego
nur den kleinen relevanten Teil (im Kern: der `Node`-Typ und die HTML-Helfer) in
die eigene Form — über das Paket `dom`. So fühlt sich nichts wie eine fremde
Bibliothek an, die danebensteht.

gomponents bleibt dabei als **Starthilfe und Inspiration** genannt (MIT,
Attribution in `THIRD_PARTY_NOTICES.md`). dreego ist **kein** Fork im engeren
Sinn: Aufbau und Umfang sind eigenständig, die Abhängigkeit ist eine schmale
Basis (~10 %).

## Warum Nutzer gomponents nie direkt importieren

dreego kapselt gomponents im Paket **`dreego/dom`**. Ein normales Projekt
schreibt:

```go
import (
	"github.com/LukasLow/dreego/v0/pkg/dreego"
	d "github.com/LukasLow/dreego/v0/pkg/dom"
)
```

und benutzt `d.Div`, `d.Text`, `d.Class`, `d.Map` — nie
`maragu.dev/gomponents`. So bleibt die Abhängigkeit eine **Implementierungs-
details** von dreego; wenn dreego sie je austauscht, ändert sich für Nutzer
nichts.

## Keine weiteren Laufzeit-Abhängigkeiten

Der dreego-Kern hat außer gomponents **keine** externen Laufzeit-Abhängigkeiten.
Routing, Session, CSRF, i18n, Markdown, Scoped CSS, Static, Gzip und Fehlerseiten
sind reine Standardbibliothek.
