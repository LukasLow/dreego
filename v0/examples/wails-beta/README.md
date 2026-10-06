# wails-beta — dreego in a Wails v3 window

A Wails v3 desktop application whose window is rendered by dreego. It is the
dreego counterpart of the Wails `vanilla` template: **same UI, same logic** (a
greeting round-trip to Go and a live "time" event), but the frontend is dreego
Go values emitting plain JavaScript — **no TypeScript, no npm, no Vite, and no
TCP listener**.

dreego's `*App` is itself an `http.Handler`, so Wails' asset server calls it
in-process; nothing listens on a port.

## Layout

```
wails-beta/
├── main.go                    application.New + dreego handler + GreetService
├── greetservice.go            Wails service (GreetService.Greet)
├── go.mod / go.sum            module wails-dreego; replace dreego -> ../../..
├── Taskfile.yml
├── build/                     Wails build assets (no frontend/vite steps)
│   └── config.yml             dev_mode: rebuild + restart (no dev server)
└── frontend/                  the dreego site (package frontend)
    ├── frontend.go            NewApp: dreego app + static mount
    ├── layout.go              Shell: the page <head>/<body> layout
    ├── page.go                Index: the page (header, title, greet, footer, toast)
    ├── client.go              the client module (plain JS, port of main.ts)
    └── static/
        ├── style.css          the vanilla Neo-Night stylesheet
        ├── wails.png, javascript.svg, bg-*.jpg
        └── bindings/          generated Wails bindings (served at /static/bindings/…)
```

The page's markup lives in `frontend/page.go` (dreego `d.Div`, `d.H1`, …); its
client script is inlined as `<script type="module" nonce="…">`. The Wails
bindings are generated into `frontend/static/bindings` and served by dreego at
`/static/bindings/…`.

## Build & run

The examples build against the checkout they live in (`replace
github.com/LukasLow/dreego => ../../..`), so no extra setup is needed.

```sh
go build .                    # quick compile
task darwin:build             # full Wails build (icons, bindings, binary)
task darwin:run               # build the .app dev bundle and launch it
task ios:run                  # build and launch in the iOS simulator
task dev                      # wails3 dev: rebuild + restart on change
```

There is no `npm install` and no `wails3 dev` frontend server: a build regenerates
the JS bindings (`wails3 generate bindings -d frontend/static/bindings -b
-noevents ./...`) and then compiles the Go binary.

## Client JavaScript

`frontend/client.go` is the plain-JavaScript port of the vanilla `main.ts`. It
imports the Wails runtime and the generated binding as ES modules:

```js
import { Events, WML } from "/wails/runtime.js";
import { GreetService } from "/static/bindings/wails-dreego/index.js";
```

`/wails/runtime.js` is served by Wails itself; the binding is a dreego static
asset. There is no bundler — the browser loads them directly.
