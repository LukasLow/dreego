package main

import (
	"embed"
)

//go:embed public/*
var publicFS embed.FS

//go:embed locales/*.json
var localeFS embed.FS
