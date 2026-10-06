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

---

## Vorschläge (noch nicht entschieden)

### Asset-Bundling — `/assets/bundle.css` + `bundle.js`
**Wunsch:** CSS und JS aller Komponenten sollen beim Serverstart bekannt sein und
als ein (minifizierter) Bundle ausgeliefert werden. Heute entstehen Assets erst
zur Request-Zeit über den per-Request-`Collector`; es gibt keine Aufzählung.

**Warum es liegt:** Kein Build-Schritt (bewusste Philosophie), Assets teils
dynamisch (`c.Style`), JS-Minify ist nicht trivial.

**Möglichkeiten (Stufenmodell):**
1. **Registry beim Start** — Assets auf Package-Ebene (`var heroCSS = scope.CSS(…)`),
   `init()` füllt eine prozessweite Registry; `app.Assets("/assets")` liefert
   `bundle.css`/`bundle.js`, einmalig minifiziert + gecacht. Lockert die
   Co-Location-Regel (Assets lexikalisch neben der Komponente).
2. **Per-Page-Cache (Collector-Modus)** — nach dem ersten Render kennt der
   Collector die genutzten Assets; content-hashed als `/assets/<hash>.css`
   gecacht und per `<link>` referenziert. Kein Regelbruch, kein Over-Ship;
   Cold-Start beim ersten Hit.
3. **Codegen / `go/ast`-Scan** — Tool extrahiert `scope.CSS/JS`-String-Literale,
   scoped + minifiziert, schreibt Manifest/Bundle. Echte Minifier, statisch
   prüfbar; führt Build-Schritt ein, nur Literale (~90 %), `Sprintf`-CSS fällt raus.

**Vorbilder:** Astro (Inseln + Critical CSS), htmx/Turbo (HTML-over-the-wire),
Qwik (Resumability). Bundling vs. Splitting: heute Splitting pro Route,
Vendor/App-Trennung, Content-Hash + langes Caching, HTTP/2-Multiplexing.
CSP erlaubt same-origin Bundles bereits (`script-src/style-src 'self'`).

### SSE + htmx-Integration
**SSE (Server-Seite):** htmx hat eine SSE-**Extension** (`sse-connect`,
`sse-swap`), aber nur zum **Empfangen**; den `text/event-stream`-Erzeuger liefert
htmx nicht — das ist dreegos Aufgabe. dreego hat kein SSE. Fix: zero-dependency
`c.SSE()` über stdlib (`http.Flusher`), `Content-Type: text/event-stream`, plus
Keep-alive. WebSocket würde eine Dependency brechen → bewusst SSE statt WS.

**htmx-Integration (nur serverseitig möglich):**
- Fragment-Modus: bei `HX-Request` nur den Body rendern (Layout überspringen),
  z. B. `c.IsPartial()`.
- Response-Header: `HX-Trigger`, `HX-Redirect`, `HX-Location`, `HX-Push-Url`,
  `HX-Retarget`.
- Re-Hydration nach Swap: `buildScript` (`scope/script.go`, `window.__dreego`-Guard)
  läuft pro Komponente nur EINMAL — nachgeladene Fragmente initialisieren sich
  nicht neu. Fragment-Navigation braucht Re-Init bei DOM-Mutation.

**Weitere mögliche Addons:** Streaming-HTML (chunked statt Vollpuffer in
`documentNode`), Speculation-Rules-Helfer, View-Transitions-Beispiel.

### Sicherheits-Härtung (zero-dependency)
Aus Code-Review von `middleware.go`, `session.go`, `csrf.go`, `redirect.go`,
`host.go`, `api.go`. Vorhanden ist schon viel Gutes (strikte CSP+Nonce,
AES-256-GCM-Sessions, constant-time CSRF, `MaxBytesReader` in `api.go`, Recovery).

- **HSTS fehlt:** `SecurityHeaders` setzt kein `Strict-Transport-Security` →
  SSL-Strip. Nur bei echtem TLS senden (nicht in Dev).
  `max-age=31536000; includeSubDomains; preload`.
- **Origin-Check als Defense-in-Depth:** `csrfValid()` prüft nur das Token;
  zusätzlich `Origin`/`Sec-Fetch-Site` für unsafe Methods verlangen.
- **Offener Redirect:** `Redirect(url, code)` validiert die Ziel-URL nicht
  (`//evil.com`, `\r\n`). Relative Pfade erlauben, absolute gegen Allowlist,
  `\r\n` ausschließen.
- **Kein Body-Limit für HTML-POSTs:** nur `api.go` nutzt `MaxBytesReader`; der
  Page-Pfad (`servePage`) liest unbegrenzt in den RAM (DoS). Limit für Pages setzen.
- **Server-Timeouts fehlen:** `Start()` (`host.go`) baut `http.Server` ohne
  `ReadHeaderTimeout`/`ReadTimeout`/`WriteTimeout`/`IdleTimeout` → Slowloris.
  Defaults setzen bzw. `SetTimeouts` anbieten.
- **`Cache-Control: no-store`** auf authentifizierten Seiten (Session/Flash).
- **Header-Feinschliff:** `Cross-Origin-Opener-Policy: same-origin`,
  `Cross-Origin-Resource-Policy: same-origin`; `Permissions-Policy` um
  `payment=(), usb=()` ergänzen.
- **`SameSite=Strict`** für die App-Session erwägen (Default ist `Lax`),
  alternativ `Lax` + Origin-Check.
- **Timing:** CSRF-Token-Erzeugung vs. Nicht-Existenz (niedrige Priorität).

### Page-Options: Proof-of-Work (opt-in)
Idee: eine **aktiv aktivierbare** Page-Option, bei der der Client einen Hash
löst (Proof-of-Work) als Spamschutz (z. B. Registrierungs-/Kontaktformular),
ohne Captcha. Nonce + Difficulty aus dem Server, Client rechnet, Server prüft.
Bewusst opt-in pro Seite. Zu klären: Difficulty/Zeitbudget, Barrierefreiheit,
Replay-Schutz, Wechselwirkung mit Rate-Limit.

### Cookie-Banner + Cookie-Konsistenz
Idee: `CookiePolicy`/Cookie mit Options (`true`/`false` pro Kategorie) so
erweitern, dass dreego eine **Default-Cookie-Banner-Komponente** mitliefert, die
mit den tatsächlich gesetzten Cookies **im Sync** bleibt (gesetztes Cookie ↔
angezeigter Banner ↔ Consent). Ziel: ein Cookie „hat immer alles richtig", und
Banner/Setzen driften nie auseinander. TDDDG-Konformität (Consent vor nicht
essenziellen Cookies) als Leitplanke. Zu entscheiden: Umfang (nur essenziell +
Statistik?), Consent-Persistenz, Opt-in/Opt-out.

### `d.htmx` / `d.tailwind`
- **`d.htmx`:** dünne Helfer (Attribute wie `hx-get`, `hx-post`, `hx-target`,
  `hx-swap`, `hx-trigger`) plus serverseitige Header/Fragment-Modus (siehe SSE +
  htmx). Klein und passend.
- **`d.tailwind`:** schwierig — bräuchte entweder das Tailwind-CLI im Modul oder
  eine eigene Tailwind-Nachbildung (alle Utilities selbst parsen), um beim
  Serverstart `assets/tailwind.css` zu erzeugen. Widerspricht dem
  Utility-CSS-Verbot (AGENTS.md) und dem „no build step"-Prinzip. Eher
  nicht bauen; als bewusste Alternative/Entscheidung festhalten.

### Weitere Ideen
- `/ready`-Endpunkt, Logging-Middleware, `SafeURL`/`SafeScript` (bereits oben).
- CSP-Report-only-Modus + Report-Endpoint (`report-to`) für Rollout.
- Signed URLs / CSRF für Static-Cache-Invalidierung.
- Health-/Metrics-Addon (Prometheus-Textformat), zero-dependency.
- `d.i18n`-Lücken: Plural-Regeln, Datums-/Zahlenformat (de-DE).

