# Markdown-Seiten

Für Rechtstexte (Impressum, Datenschutz, AGB) schreibt man Markdown und rendert
es zu Nodes. Vorbild: Dreegos `<body lang="md">`.

## Benutzen

```go
import "github.com/LukasLow/dreego/v0/pkg/markdown"

const impressum = `
# Impressum

Angaben gemäß § 5 DDG.

## Kontakt

E-Mail: [hallo@example.de](mailto:hallo@example.de)
`

func getImpressum(c *dreego.Ctx) d.View {
	return d.Main(g.Group(markdown.ToNodes(impressum)))
}
```

## Unterstützter Umfang

- Überschriften `#` … `######`
- Absätze
- Listen (`-`, `*`, `+`)
- **fett**, *kursiv*, `Code`
- Links `[Text](url)`
- Blockzitate `>`
- Horizontale Linie `---`

## Sicherheit

- **Kein Roh-HTML.** HTML im Markdown wird escaped (`<script>` wird zu
  `&lt;script&gt;`).
- **Sichere Link-Schemata.** Nur `http`, `https`, `mailto`, `tel` und relative
  Links werden gerendert. `javascript:` & Co. werden verworfen — das Label bleibt
  als Text.

## Bewusst nicht enthalten

Tabellen, Fußnoten, verschachtelte Listen, HTML-Blöcke, Syntax-Highlighting.
Für diese reicht der Umfang von v0.1 nicht; sie wären eine spätere Erweiterung.
