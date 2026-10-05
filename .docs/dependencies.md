# Abhängigkeiten und Lizenz

## dreego selbst

dreego steht unter der **MIT-Lizenz** (siehe `LICENSE`). Du darfst es verwenden,
ändern, weitergeben und kommerziell nutzen — nur der Copyright-Hinweis und der
Lizenztext müssen erhalten bleiben.

## gomponents

dreego baut auf **gomponents** auf.

| | |
|---|---|
| Modul | `maragu.dev/gomponents` |
| Echte Quelle | https://github.com/maragudk/gomponents |
| Version (gepinnt) | `v1.3.0` |
| Lizenz | **MIT** — Copyright (c) Maragu AG |

**Zur Domain:** `maragu.dev` ist die schöne Modul-Adresse; sie zeigt auf das
GitHub-Repo `maragudk/gomponents`. Die kanonische Import-ID ist bewusst die
Domain (so empfiehlt es Go), nicht die GitHub-URL.

**Was MIT bedeutet für dich:**
- ✅ Verwenden, ändern, weitergeben, verkaufen.
- ✅ Kompatibel mit dreegos MIT-Lizenz.
- ✅ **Keine Copyleft-Pflicht** — deine Anwendung muss **nicht** offen sein.
- ⚠️ Der Copyright-Hinweis und Lizenztext müssen bleiben.

Der vollständige Lizenztext von gomponents ist in **`THIRD_PARTY_NOTICES.md`**
reproduziert (das verlangt die Lizenz).

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
