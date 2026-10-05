# Komponenten (CSS + JS in einer Funktion)

Eine Komponente ist eine Go-Funktion, die einen `d.View` zurückgibt. Mit `c.Box`
bekommt sie HTML, CSS und JavaScript **in derselben Funktion** — gescopet,
dedupliziert, mit CSP-Nonce.

## Beispiel

```go
func Zaehler(c *dreego.Ctx) d.View {
	return c.Box(
		dreego.CSS(`
.zaehler { display: inline-flex; gap: 14px; border: 2px solid #0080ff; padding: 12px }
.zaehler .wert { font-size: 22px; font-weight: 700 }
`),
		dreego.JS(`
var knopf = root.querySelector("button");
var wert = root.querySelector(".wert");
knopf.addEventListener("click", function () {
    wert.textContent = String(Number(wert.textContent) + 1);
});
`),
		d.Div(d.Class("zaehler"),
			d.Button(d.Text("+1")),
			d.Span(d.Class("wert"), d.Text("0")),
		),
	)
}
```

Verwenden:

```go
func getSeite(c *dreego.Ctx) d.View {
	return dreego.Div( /* … */ Zaehler(c) )
}
```

## Was `c.Box` tut

1. **Hasht** CSS+JS → eine `data-scope`-ID.
2. **Präfixt** jeden CSS-Selektor mit `[data-scope="…"]` → kein Leck.
3. **Dedupliziert**: gleiche Komponente mehrfach = CSS/JS nur einmal.
4. **Kapselt JS**: das Skript bekommt den Komponenten-Root als `root`; es kann
   fremde Elemente nicht treffen, auch bei gleichen Klassennamen.
5. **Läuft einmal**: mehrere Instanzen teilen sich ein Skript (Guard).
6. **Sichert**: `</script>` im Komponenten-JS wird entschärft; jede `<style>`/
   `<script>` trägt den CSP-Nonce.

## Regeln

- **Varianten sind Props, keine Forks.** Ein neues Aussehen ist ein neuer Wert
  eines Parameters, nicht eine kopierte Komponente.
- **Kein `style="…"`-Attribut.** dreegos CSP (`style-src 'nonce-…'`) blockiert
  Inline-Style-Attribute. Layout gehört in eine gescopte CSS-Klasse. (Die
  Komponenten-Bibliothek hält sich daran; ein Browser-Check hat genau diesen
  Fehler einmal aufgedeckt.)
- **`CSS`/`JS` sind Nodes**, keine Strings mit Sonderbehandlung.

## Wann eine eigene Komponente?

- Wird an **mehr als einer** Stelle verwendet → eigene Funktion.
- Nur einmal → ruhig inline in der Seite.

## Beispiel-Dateien

- `v0/examples/fensterbank/shop.go` — Komponenten aus `dreego/ui`.
- `v0/examples/www/components/hero.go` — Komponente mit CSS **und** JS.
- `scope/` — die Mechanik dahinter.
