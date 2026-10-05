# Seiten und Layouts

## `Page` — die Felder

```go
var Shop = dreego.Page{
	Path:       "/shop",              // die URL (einzige Quelle)
	Layout:     Marketing,            // Hülle; ohne Angabe: defaultLayout
	Nav:        "shop",               // aktiver Navigationspunkt (c.Nav())
	Head:       headShop,             // Kopf-Inhalt (Titel, Meta)
	Get:        getShop,              // GET-Handler -> d.Node
	Post:       postShop,             // optional
	API:        nil,                  // optional: JSON/XML statt HTML
	Methods:    nil,                  // optional: explizite Methodenliste
	Middleware: nil,                  // optional: nur für diese Seite
}
```

## Der `Head`

`Head` liefert Kopf-**Inhalt** (kein `<head>`-Element — das baut das Dokument):

```go
func headShop(c *dreego.Ctx) []d.Node {
	return []d.Node{
		d.TitleEl(d.Text("Shop — Beispiel")),
		d.Meta(d.Name("description"), d.Content("…")),
	}
}
```

Der Titel kommt in die Kopfzeile, aber dreego setzt ihn **nicht** automatisch
ins Dokument-Titel-Element — das macht das Layout über `{head}`.

## Layouts

Ein Layout ist eine Funktion `(c, head, body) Node`:

```go
func Marketing(c *dreego.Ctx, head d.Node, body d.Node) d.Node {
	kopf := d.Group([]d.Node{
		d.Meta(d.Charset("utf-8")),
		d.Meta(d.Name("viewport"), d.Content("width=device-width, initial-scale=1")),
		d.Meta(d.Name("theme-color"), d.Content("#fffdf7")),
		head,
	})

	return c.Document("de", kopf, d.Group([]d.Node{
		header(c),
		d.Main(body),
		footer(c),
	}))
}
```

- `c.Document(lang, kopf, body)` baut das volle Dokument: `<!DOCTYPE html>`,
  `<html lang>`, `<head>` (mit gesammelten Styles), `<body>` (mit Scripts).
- **Reihenfolge im `<head>`:** Kopf-Inhalt → kritisches CSS → Komponenten-Styles.
- **Scripts** landen ans Ende von `<body>`.

## Kein weißer Blitz (Critical CSS)

Ein `<link rel="stylesheet">` blockiert das Rendern — der Browser malt weiß, bis
die Datei da ist. dreego löst das, indem das Basiscss **inline** im Kopf steht:

```go
func Shell(c *dreego.Ctx, head d.Node, body d.Node) d.Node {
	c.AddCritical(meinInlineCSS)   // sofort gefärbt, kein Blitz
	…
}
```

Das komplette Basis-Stylesheet inline + `theme-color` + `color-scheme` = der
erste Paint ist bereits thematisiert. (Vorbild: Dreegos Generator schreibt
`<style>` ebenfalls inline.)

## Default-Layout

Ohne `Layout` nutzt dreego ein minimales Dokument (utf-8, viewport) — korrekt,
aber schmucklos. Für echte Seiten immer ein Layout angeben.
