package main

import (
	"embed"

	"github.com/LukasLow/dreego/v0/pkg/urls"
)

// publicFS hält statische Dateien (go:embed). Sie werden unter /public/* ausgeliefert.
//
//go:embed public/*
var publicFS embed.FS

// Die Site-Adressen. Dev/Prod unterscheiden sich nur über die Env-Variablen.
var (
	siteApp = urls.FromEnv("DREEGO_APP_URL", "http://localhost:8081")
	siteWWW = urls.FromEnv("DREEGO_WWW_URL", "http://localhost:8080")
)
