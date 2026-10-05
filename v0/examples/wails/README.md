# wails-dreego (example)

A Wails v3 desktop window whose frontend is plain **dreego** — no `.dreego`
files, no generator, no npm, no Vite, no local HTTP server.

- The window is served by `github.com/LukasLow/dreego/v0/plg/wails`, which is a
  tiny adapter: it returns an `http.Handler` for Wails'
  `AssetOptions.Handler`, with a Content-Security-Policy a Wails webview can
  satisfy.
- The application owns the Wails window, service, bindings and lifecycle.

## Layout

```
v0/examples/wails/
├── main.go                 application.New + wails.Handler(app) + service + window
├── go.mod / go.sum         module wails-dreego (replace -> ../../..)
├── Taskfile.yml            bindings / build / run / dev
├── build/                  Wails build assets
├── app/
│   ├── greeter.go          application-owned Wails service
│   ├── bindings/           generated TypeScript bindings (committed)
│   ├── static/bindings/    generated JavaScript bindings (committed, served)
│   └── web/                dreego pages (pure Go)
│       ├── routes.go       Register(app)
│       ├── shell.go        layout + shared CSS
│       ├── home.go         GET /
│       └── about.go        GET /about
└── README.md
```

## Prerequisites

- Go 1.27+
- The `wails3` CLI (v3.0.0-beta.27)
- On Linux: GTK4 + WebKitGTK development headers (Wails is CGO).

## Run

```sh
# straight run (no build assets needed):
go run .

# or via Wails:
task run        # platform task (go build into bin/)
task dev        # rebuild + restart on change
task package    # platform package (.app / .deb / …)
```

## Bindings

Wails generates typed Go-to-client bindings from the **application-owned
services** (`app/greeter.go`, registered in `main.go`). Two sets are committed:

| Set | Path | Purpose |
|---|---|---|
| TypeScript | `app/bindings/<module>/<pkg>/…` | typing / editor |
| JavaScript | `app/static/bindings/<module>/<pkg>/…` | runtime, served at `/bindings/…` |

The page imports the runtime binding dynamically from
`/bindings/wails-dreego/app/greeterservice.js` — it only resolves in the
webview. A page with no interaction needs no bindings at all.

## CSP

`wails.Handler` applies a Wails-suitable CSP (allows the injected runtime and
inline bindings, still forbids framing/plugins/foreign origins). If the app sets
its own CSP (`app.SetCSP(...)`), that is respected and not overwritten.
