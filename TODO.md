# TODO — dreego (v0)

Offene, bewusst zurückgestellte Punkte. **Nicht** Teil der aktuellen v0.0.x.

Lebendes Dokument: Punkte wandern hierher, wenn sie bewusst später entschieden
werden; erledigte Punkte werden entfernt.

---

## Erledigt (v0.0.2) — entfernt

- ~~Multi-Site (mehrere Ports, ein Prozess)~~ → **erledigt**: `dreego.Start` /
  `Wait` / `Addr` / `Close` (`pkg/dreego/host.go`) + `examples/multisite`.
- ~~Session-Cookie `MaxAge`/`Expires`~~ → **erledigt**: `CookiePolicy.MaxAge`
  (time.Duration) + `Days`/`Hours`/`Minutes`; setzt Max-Age und Expires.
- ~~CSP & dynamische Inline-Styles~~ → **erledigt**: `c.Style(...)` erzeugt eine
  nonce'd Klassenregel über den Collector.

---

---

## Zurückgestellt (spätere Phase)

### Markdown: Roh-HTML & Tabellen
**Warum es liegt:** Statunas Rechtsseiten (`<body lang="md">`: Impressum,
Datenschutz, AGB, AVV) mischen Markdown mit HTML-Blöcken (`<p>`, `<div>`) und
brauchen Tabellen. Das neue dreego gibt Roh-HTML als Text aus (`<script>` wird
`&lt;script&gt;`) und kennt **keine** Tabellen.

**Warum bewusst später:** Roh-HTML-Durchlass ist eine Sicherheitsentscheidung
(XSS-Fläche) — das will man nicht nebenbei ändern.

**Möglichkeiten (zu entscheiden):**
1. Tabellen nativ ergänzen (klein, sicher, reines Markdown).
2. Opt-in Roh-HTML wie im alten Dreego (`lang="md"`) — braucht eine klare,
   eng begrenzte Regel.
3. Kein Markdown für Rechtsseiten — in reinem Go-HTML (dom) bauen.

### Named Slots für Komponenten
Statunas Komponenten nutzen benannte Slots (`{#slot nummer}`, `{#slot text}`),
z. B. `HowItWorksCards`, `FaqItem`. Das neue Modell kennt keine Slots.
Aktuell: Kindinhalt als `d.View`-Argument übergeben.
Zu entscheiden: reicht das, oder braucht dreego ein echtes Slot-System?

### SVG-Helfer — ERLEDIGT (v0.0.3)
`pkg/dom/svg.go` ergänzt `Circle`, `Rect`, `Path`, `G`, `Line`, `Polygon`,
`Polyline`, `Ellipse`, `Defs`, `Use`, `TextPath` und die Attribute `ViewBox`,
`Cx`, `Cy`, `R`, `D`, `Fill`, `Stroke`, `StrokeWidth`, `StrokeLinecap`,
`StrokeLinejoin`, `StrokeDasharray`, `Transform`, `AriaLabel`, …

### `/ready`-Endpunkt
Es gibt nur `/health`. Monitoring/Readiness braucht `/ready`.
Fix: wie `/health` fest verdrahten oder konfigurierbar machen.

### Logging-Middleware
Das alte Dreego hatte `SetLogging(true)`. Das neue hat nur `app.Use`.
Fix: kleine, optionale Request-Log-Middleware als Beispiel/Addon.

### Öffentliches `SafeURL` / `SafeScript`
Das alte Dreego exportierte `SafeText/SafeAttr/SafeURL/SafeScript`. Das neue
escaped über `d.Text`/`d.Attr`; für URL-/Script-Kontexte fehlt eine öffentliche
Hilfe (Markdown prüft Links intern selbst).
Zu entscheiden, ob überhaupt gebraucht.
