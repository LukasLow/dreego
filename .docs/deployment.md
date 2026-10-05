# Deployment

## Als ein Binary

dreego ist reines Go — `go build` ergibt ein Binary ohne Laufzeit-Abhängigkeit
(außer den eingebetteten Assets).

```sh
go build -o dreego-app ./meine-app
```

## Environment

Wichtige Variablen (Namen frei wählbar, hier die Demo-Konvention):

| Variable | Bedeutung |
|---|---|
| `DREEGO_SESSION_SECRET` | Session-Geheimnis, **≥ 32 Zeichen, geheim** |
| `DREEGO_PORT` | Listen-Port (Apps nutzen oft `-port`-Flag) |
| `DREEGO_DEV` | `1` erlaubt ein Dev-Standardgeheimnis (nie in Produktion!) |

**Sicherheitsregel:** Ohne gesetztes `DREEGO_SESSION_SECRET` bricht der Start ab
(kein stilles Standardgeheimnis).

## Docker

`Dockerfile` (Multi-Stage, nonroot):

```dockerfile
FROM golang:1.24-alpine AS bauen
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ARG CMD=./v0/examples/fensterbank
RUN CGO_ENABLED=0 go build -o /dreego-demo ${CMD}

FROM alpine:3.20
RUN adduser -D -H -u 10001 dreego
USER dreego
COPY --from=bauen /dreego-demo /dreego-demo
EXPOSE 4000
ENTRYPOINT ["/dreego-demo"]
```

Welche Site gebaut wird, steuert `--build-arg CMD=./sites/<name>`.

## Mehrere Sites (mehrere Ports)

Eine „Site" ist eine eigene `dreego.NewApp()` mit eigenem Handler und eigenem
Port — wie Statuna (www :3000, app :3001, link :3002). Cross-Links über
`dreego/urls`, dev/prod nur per Environment.

`docker-compose.yml` startet zwei Sites:

```yaml
services:
  fensterbank:
    build: { context: ., args: { CMD: ./v0/examples/fensterbank } }
    command: ["-port", "4000"]
    ports: ["4000:4000"]
    environment:
      DREEGO_SESSION_SECRET: "…≥32 Zeichen…"
      PORTAL_URL: "http://localhost:4001"
  portal:
    build: { context: ., args: { CMD: ./v0/examples/portal } }
    command: ["-port", "4001"]
    ports: ["4001:4001"]
    environment:
      DREEGO_SESSION_SECRET: "…anderes Secret…"
      FENSTERBANK_URL: "http://localhost:4000"
```

```
docker compose up -d --build
# http://localhost:4000  (Fensterbank)
# http://localhost:4001  (Portal)
```

## Hinter einem Reverse-Proxy (Traefik)

- TLS terminiert Traefik → HTTP/2 für den Browser automatisch.
- Session-Cookies `Secure` setzen: `store.SetTrustedProxies([…])` mit den
  Proxy-IPs, damit `X-Forwarded-Proto: https` gilt (und nur von dort).
- Statische Dateien länger cachen lassen (`Cache-Control`) auf Proxy-Ebene.

## Health/Ready

`GET /health` → `200 ok`, immer verfügbar. (Einen konfigurierbaren
`/ready`-Endpunkt gibt es noch nicht — offene Aufgabe.)

## Noch offen

- Dev-Watch (automatischer Neustart bei Änderung) — fehlt.
- Audit-Log für Security-Ereignisse — fehlt.
