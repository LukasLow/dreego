# dreego Demo — Multi-Stage-Build, ein statisches Binary.
#
# Das Go-Modul liegt im Repo-Root (module github.com/LukasLow/dreego).
# Die Beispiele liegen unter v0/examples/<name>.
#
# Build:  docker build -t dreego-demo .
# Run:    docker run -p 4000:4000 -e DREEGO_SESSION_SECRET=<32+ Zeichen> dreego-demo
#
# Welche Site gebaut wird, steuert das Build-Arg CMD (Default: fensterbank).

FROM golang:1.24-alpine AS bauen
WORKDIR /src

# Nur die Modul-Dateien zuerst (Layer-Cache für Abhängigkeiten).
COPY go.mod ./
RUN go mod download || true

# Quellcode.
COPY . .

ARG CMD=./v0/examples/fensterbank
RUN CGO_ENABLED=0 go build -o /dreego-demo ${CMD}

FROM alpine:3.20
RUN adduser -D -H -u 10001 dreego
USER dreego
COPY --from=bauen /dreego-demo /dreego-demo

EXPOSE 4000
ENTRYPOINT ["/dreego-demo"]
