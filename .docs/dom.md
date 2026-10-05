# DOM-Schicht `dreego/dom`

`dom` ist dreegos eigenes HTML-Modell: HTML als Go-Werte, ohne Template-Sprache
und ohne externe Bibliothek.

```go
import d "github.com/LukasLow/dreego/v0/pkg/dom"

d.Div(d.Class("karte"),
	d.H1(d.Text("Hallo")),
	d.P(d.Text("Welt")),
)
// <div class="karte"><h1>Hallo</h1><p>Welt</p></div>
```

## `View` — der Typ

```go
type View interface {
	Render(w io.Writer) error
}
```

Alles ist eine `View`: Element, Attribut, Text, Gruppe. Eine Komponente ist eine
Funktion, die eine `View` zurückgibt:

```go
func Hero(c *dreego.Ctx) dreego.View { return d.Section(d.H1(d.Text("Hi"))) }
```

**Warum `View` und nicht `Node`?** Klarer und selbsterklärend. „Node" ist
DOM-Jargon; „View" sagt, was es ist. Im Zweifel gewinnt der verständlichere Name.

## Kernfunktionen

| Funktion | Zweck |
|---|---|
| `El(name, children…)` | Element mit Namen |
| `Attr(name, value…)` | Attribut (ohne Wert = boolesch) |
| `Text(s)` / `Textf(…)` | HTML-escaped Text |
| `Raw(s)` / `Rawf(…)` | unescaped Text (nur für Vertrauenswürdiges) |
| `Group{…}` | mehrere Views als eine |
| `Map(slice, f)` | Daten → Views |
| `If(cond, view)` / `Iff(cond, f)` | bedingt rendern |
| `Doctype(sibling)` | `<!doctype html>` voranstellen |

## Elemente und Attribute

Es gibt alle gängigen HTML-Funktionen: `Div`, `H1` … `H6`, `A`, `Form`, `Input`,
`Table`, `Meta`, `Link`, `TitleEl`, … und Attribute wie `Class`, `Href`, `Src`,
`Type`, `Value`, `Checked`, `Disabled`, …

Sie werden von **`internal/gen_dom`** erzeugt (ein Autoren-Werkzeug, kein
Build-Schritt). Neu erzeugen:

```sh
go run ./internal/gen_dom <quellverzeichnis> ./pkg/dom
```

## Attribute sind Knoten (wichtig für htmx)

Jedes Attribut ist eine `View`. Darum kann eine spätere htmx-Schicht einfach
weitere Attribut-Knoten liefern:

```go
d.Div(
	d.Class("liste"),
	htmx.Get("/status"),        // hx-get="/status"
	htmx.Target("#liste"),      // hx-target="#liste"
	htmx.Trigger("every 5s"),   // hx-trigger="every 5s"
)
```

Kein Sonderfall, kein Extra-Modell — htmx lebt im normalen Attribut-Modell.

## Sicherheit

- `Text` escaped automatisch (`<` → `&lt;` usw.).
- `Attr` escaped den Wert.
- `Raw` ist das bewusste Opt-out.
- Void-Elemente (`br`, `img`, `input`, …) erhalten korrekt keinen Schluss-Tag.
