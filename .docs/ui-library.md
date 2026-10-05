# Komponenten-Bibliothek `dreego/ui`

Fertige, gestylte Bausteine. Alle nehmen den `Ctx` und Props, sind gescopet und
variantenreich. Import:

```go
import "github.com/LukasLow/dreego/v0/pkg/ui"
```

## Button

```go
ui.Button(c, "/pricing", "Preise ansehen", ui.Primary)
ui.Button(c, "/features", "Mehr erfahren", ui.Ghost)
```

- `Primary` — gefüllt.
- `Ghost` — umrandet, transparent.

## Card

```go
ui.Card(c, "Projekte", d.P(g.Text("Alles an einem Ort.")))
```

Eine Box mit Titel und beliebigem Body-Inhalt. Füllt ihren Grid-/Flex-Platz
(`height: 100%`).

## Badge

```go
ui.Badge(c, "Aktiv", ui.Success)
```

Töne: `Info` (blau), `Success` (grün), `Warning` (orange), `Danger` (rot).

## Alert

```go
ui.Alert(c, "Bitte Eingaben prüfen.", ui.Danger)
```

Wie Badge, aber als Meldungsbox. `Warning`/`Danger` setzen `role="alert"`
(assertiv), sonst `role="status"` — für Screenreader korrekt.

## Field

```go
ui.Field(c, form, "email", "E-Mail", "email")
```

Ein Formularfeld mit:
- assoziiertem `<label>` (`for`/`id`) — Barrierefreiheit,
- `Old`-Wert aus dem Formular,
- Fehlermeldungen (`role="alert"`).

`form` ist ein `*dreego.Form` (aus `dreego.Bind`) oder `nil` für ein frisches
Formular. Siehe [forms.md](forms.md).

## Eigene Komponenten ergänzen

Die Bibliothek ist offen. Eine neue Komponente ist einfach eine Funktion nach
demselben Muster:

```go
func Kpi(c *dreego.Ctx, wert string, label string) d.View {
	return c.Box(
		dreego.CSS(`.kpi { text-align: center } .kpi b { font-size: 32px; color: #0b6b4a }`),
		d.Div(d.Class("kpi"),
			d.B(d.Text(wert)),
			d.Span(d.Text(label)),
		),
	)
}
```

Regel bleibt: über `c.Box`, Varianten als Props.
