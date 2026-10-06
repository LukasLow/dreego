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

### Sicherheits-Härtung — Erweiterungen
Ergänzend zur Liste oben (COOP/CORP/COEP, Trusted Types, Host-Header-Prüfung,
security.txt, Session-Fixation/-Revocation, Content-Type-Pflicht,
`MaxHeaderBytes`, Trusted-Proxy-IP, Audit-Log, HMAC-signierter State):

- **COOP/CORP/COEP:** `Cross-Origin-Opener-Policy: same-origin`,
  `Cross-Origin-Resource-Policy: same-origin` (ggf. `same-site`), optional
  `Cross-Origin-Embedder-Policy: require-corp` → Prozess-Isolation (Spectre).
- **Trusted Types:** `require-trusted-types-for 'script'` in der CSP, wo JS DOM
  schreibt — macht DOM-XSS strukturell unmöglich (Browser-Support beachten).
- **Host-Header-Validierung:** erlaubte Hosts prüfen → Host-Header-Injection &
  Cache-Poisoning hinter Proxy verhindern.
- **Content-Type-Pflicht** für POST/PUT/PATCH: fehlend/falsch → ablehnen.
- **Trusted-Proxy-IP:** Client-IP nur aus `X-Forwarded-For` lesen, wenn
  RemoteAddr im Trusted-Set liegt (sonst IP-Spoofing im Rate-Limit).
- **Session-Fixation:** Session bei Login/Privilegienwechsel rotieren.
- **Audit-Log:** einheitliches Log für sicherheitsrelevante Events
  (Login, Consent, Löschung) → DSGVO-Nachweis.
- **security.txt:** Helfer für `/.well-known/security.txt`.

### Fehlende Kern-Fähigkeiten
- **Server-Session-Store:** heute nur Cookie mit hartem 4-KB-Limit; optionaler
  In-Memory/DB-Store → beliebig groß + „überall abmelden" (Revocation).
- **Uploads / Multipart:** fehlt komplett — Größenlimit je Feld, Streaming in
  Temp/Storage, sichere Dateinamen, MIME-Prüfung.
- **CORS-Helfer** (`app.SetCORS`) für `API:`-Endpunkte (preflight, Origins,
  Credentials).
- **ETag / Conditional GET (304)** für Static und gerenderte Seiten.
- **Graceful `Shutdown(ctx)`** neben `Close()`.
- **Request-ID + strukturiertes Logging** (Observability).
- **Streaming-HTML:** `documentNode` puffert alles → chunked Rendering.

### Dev-Ergonomie
- **CLI `dreego new`:** Scaffold im Statuna-Layout (`APP_/COMP_/PAGE_`,
  www/app/link, Store, Layouts).
- **OpenAPI aus `API:`-Seiten:** Typen/Validierung sind da → Spec + Client.
- **Config-Loader `env`:** typisiert, fehlende Pflichtwerte brechen den Start ab.
- **Meta-/OG-Helfer:** Title, Description, Open Graph, Twitter-Card, JSON-LD.
- **robots / sitemap** aus registrierten Seiten erzeugen.
- **i18n:** Plural-Regeln + `de-DE` Datums-/Zahlenformat.
- **`/ready`-Endpunkt**, Logging-Middleware, `SafeURL`/`SafeScript` (oben).
- **Health-/Metrics-Addon** (Prometheus-Textformat), zero-dependency.
- **CSP report-only** + Report-Endpoint (`report-to`) für Rollout.

### Animationen (View Transitions statt Animation-Runtime)
**Frage:** Svelte hat `svelte/transition` und `svelte/animate` (fade, fly, slide,
`animate:flip` in keyed each). Direkt übernehmbar ist das **nicht**: Svelte
braucht dafür seinen reaktiven Runtime + Keyed-Diffing, um „rein/raus" zu
erkennen. dreego hat keinen Client-Runtime und kein Reconciliation — der Server
liefert HTML, der Browser fügt Knoten ein/aus.

**Was dreego-idiomatisch geht (zero-dep):**
1. **CSS-Transitions/-Keyframes** — scoped über `scope.CSS`; deckt Fade, Slide,
   Hover, Entrance/Exit, `@starting-style` + `transition-behavior: allow-discrete`.
2. **View Transitions API** (der große Hebel) — browsernativ, MPA **und**
   Fragment-Swaps: `@view-transition { navigation: auto }` bzw.
   `document.startViewTransition(...)`; mit htmx kombinierbar
   (`hx-swap="outerHTML transition:true"`). SPA-artige Übergänge **ohne** Framework.
3. **FLIP-Helfer** (≈40 Zeilen JS als `scope.JS`-Komponente): First–Last–Invert–Play
   über `getBoundingClientRect` — Svelte-`animate:flip`-Äquivalent für Listen,
   das der Browser ohne Runtime ausführt.
4. **Scroll-driven Animations** (`animation-timeline: scroll()/view()`) — rein CSS,
   sehr effektvoll, kein JS.
5. **Web Animations API / Mini-Helfer** für komplexe Fälle (statt Motion-One-Dep).

**Zu entscheiden:** ein `dreego/anim`-Addon (CSS-Keyframes-Helfer + optionaler
FLIP-Helfer) vs. nur Doku/Beispiel. Bewusst **keine** Svelte-artige
Animation-Runtime (würde die Philosophie brechen). Prüfen: `prefers-reduced-motion`
respektieren.

**Demo-Erkenntnis (gebaut, nicht committet):** Zwei Varianten zeigen die Grenze —
- **Basis (reines HTML/CSS):** Staggered Entrance (`@keyframes`+`delay`), Hover
  Lift (`transform`+`shadow`), Progress/`box-shadow`-Pulse, CSS-only Accordion
  (`:checked`+`max-height`). Alles 0 JS, aber **Exit-Animationen** sind zäh (nur
  über `@starting-style` + `transition-behavior: allow-discrete`).
- **Fancy:** FLIP-Reorder (First–Last–Invert–Play über `getBoundingClientRect` +
  `element.animate`, echtes `animate:flip`-Äquivalent ohne Runtime), View
  Transitions (Morphing), Cursor-Glow (CSS-Variablen + `mousemove`),
  Scroll-Reveal (`animation-timeline: view()`), Endlos-Marquee.
- **Fazit der Demo:** Je „fancy“ (Morphing mehrerer Elemente, Federn/Springs,
  orchestrierte Timelines), desto mehr wünscht man sich eine Lib/Micro-Runtime.

### Plugin-System mit Micro-Runtime
**Idee:** dreego bleibt zero-dep im Kern, aber Plugins dürfen optional eine
**Micro-Runtime** mitliefern (z. B. Animationen, Signals, htmx). Jedes Plugin
registriert sich beim Start und kann eigene Assets (JS/CSS), Routen und
`scope`-Erweiterungen beitragen.

**Zu klären:**
- Plugin-Interface: `Register(app)` + `Assets()`, `Middleware()`, `Head()`,
  Version/Abhängigkeiten.
- Asset-Beitrag: koppelt an Asset-Bundling (siehe oben) — Plugin-JS wird Teil
  des Bundles oder ein eigener Chunk.
- Runtime-Grenzen: eine mitgelieferte Runtime darf den Kern nicht zwingen; sie
  lädt nur, wenn das Plugin sie aktiv braucht (sonst ist dreegos „no runtime“-Vorteil weg).
- Sicherheit: Plugins sind Code mit voller App-Rechte — Vertrauensmodell nötig.
- Reihenfolge/Lebenszyklus, Doppel-Registrierung, Deaktivieren pro App.

### Signals / Reaktivität (Micro-Runtime)
**Idee:** feingranulare Reaktivität — ändert sich eine Variable, aktualisieren
sich **alle** Stellen im Frontend, die sie verwenden. Klassisch eine
Micro-Runtime (`signal()` + Effekt + DOM-Bindung, wie SolidJS/Svelte-Runes,
Vorbild Preact Signals).

**Warum es liegt:** bricht das Kernmodell. dreego ist server-gerendert ohne
Runtime und ohne Keyed-Reconciliation; Signals brauchen genau das. Als **Opt-in
Plugin** denkbar, aber bewusst nicht im Kern.

**Zu klären:** Größe der Runtime, Interop mit `scope.JS`-Komponenten,
SSR/Re-Hydration (Server-Wert → Client-Signal), Bundle-Kosten.

### „Ist dreego dann auf React/Svelte-Niveau?“ — Bestandsaufnahme
Nach Signals/Animationen/Plugin-System fehlt zum Gleichstand immer noch:
- **Reaktives Komponenten-/State-Modell** (Signals/Props/Events) — React/Svelte
  Kern, bei dreego nur als Plugin.
- **Client-seitiges Routing** + Daten-/State-Management (Stores, Suspense/Loading).
- **Komponenten-Lifecycle & Composition** (Fragmente, Slots statt Go-Args).
- **Ökosystem:** Form-Bibliotheken, UI-Kits, DevTools, SSR-Frameworks (Next/SvelteKit).
- **Öffnen/Schließen von Schleifen:** dreego ist bewusst **Hypermedia-first**
  (Server rendert), nicht ein Client-Renderer. „Gleichstand“ wäre nur über
  zugekaufte Runtimes erreichbar — und würde den Kernvorteil (kein Build, kein
  Runtime, kleine Payload) aufgeben.

**Einordnung:** dreego ist **nicht** „React/Svelte kleiner“, sondern eine andere
Klasse (server-gerenderte Hypermedia wie Rails+Hotwire/Django+htmx). Die
Feature-Parität entsteht durch **Client-Additionen**, nicht durch den Kern.

