# dreego Demo — Multi-Stage-Build, ein statisches Binary.
#
# Das Go-Modul liegt unter v0/ (Modulpfad github.com/LukasLow/dreego/v0).
# Die Beispiele liegen unter v0/examples/<name>.
#
# Build:  docker build -t dreego-demo .
# Run:    docker run -p 4000:4000 -e DREEGO_SESSION_SECRET=<32+ Zeichen> dreego-demo
#
# Welche Site gebaut wird, steuert das Build-Arg CMD (Default: fensterbank).

FROM golang:1.24-alpine AS bauen
WORKDIR /src

# Nur die Modul-Dateien zuerst (Layer-Cache für Abhängigkeiten).
COPY v0/go.mod v0/go.sum ./
RUN go mod download

# Quellcode.
COPY v0/ .

ARG CMD=./examples/fensterbank
RUN CGO_ENABLED=0 go build -o /dreego-demo ${CMD}

FROM alpine:3.20
RUN adduser -D -H -u 10001 dreego
USER dreego
COPY --from=bauen /dreego-demo /dreego-demo

EXPOSE 4000
ENTRYPOINT ["/dreego-demo"]
