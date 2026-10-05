# Statische Dateien

`go:embed` liefert CSS, JS, Bilder, Fonts — ohne externe Abhängigkeit.

## Einbinden

```go
//go:embed public/*
var publicFS embed.FS

app.Static(publicFS, "public")   // public/favicon.svg -> /public/favicon.svg
```

- Das **Präfix ist zugleich der Ordner** im Dateisystem.
- Leeres Präfix: Dateien liegen direkt an der Wurzel (`favicon.svg` →
  `/favicon.svg`).
- Nur `GET`/`HEAD`; andere Methoden ergeben `405`.
- Reihenfolge: **Seiten gewinnen** vor Static. Ein nicht gefundenes Static →
  404-Fehlerseite.

## Content-Type

dreego bringt eine eigene Erweiterungs-Tabelle mit (`text/css`, `text/javascript`,
`image/svg+xml`, `font/woff2`, …), weil die System-MIME-Datenbank in schlanken
Containern oft fehlt und `.css` sonst als `text/plain` ausgeliefert würde. Erst
danach: System-MIME, dann Content-Sniffing.

## Caching

Static-Antworten tragen `Cache-Control: public, max-age=300`. Für Produktion
längere Werte über den Reverse-Proxy (Traefik/CDN).

## Fonts einbinden (Beispiel)

```css
@font-face {
  font-family: "Atkinson Hyperlegible";
  font-weight: 400;
  font-display: swap;
  src: url("/public/fonts/atkinson-400.woff2") format("woff2");
}
```

Fonts per `<link rel="preload" as="font" crossorigin>` vorladen spart eine
Latenzrunde.

## Kein Build-Schritt

Es gibt keine Asset-Pipeline, kein Bundling, kein Minifizieren. Dateien werden
so ausgeliefert, wie sie im Embed liegen. (Ein Bundling widerspräche dem
„kein Buildstep"-Prinzip; ggf. spätere optionale Phase.)
