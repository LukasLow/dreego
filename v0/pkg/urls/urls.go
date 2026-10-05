// Package urls builds absolute URLs for one or more sites from a base address
// and a path. It is the dreego replacement for the Dreego `urls` helper.
//
// A site is configured once from an environment variable (with a default) and
// then used everywhere, so dev and prod differ only in configuration:
//
//	www := urls.FromEnv("WWW_URL", "https://www.example.com")
//	app := urls.FromEnv("APP_URL", "https://app.example.com")
//
//	www.Abs("/pricing")    // https://www.example.com/pricing
//	app.Abs("/auth/login") // https://app.example.com/auth/login
package urls

import (
	"os"
	"strings"
)

// Site is a named base address. The zero value is unusable; build one with
// New or FromEnv.
type Site struct {
	base string
}

// New returns a Site with a fixed base address.
func New(base string) Site {
	return Site{base: strings.TrimRight(base, "/")}
}

// FromEnv returns a Site whose base comes from the environment variable env,
// falling back to fallback when the variable is empty or unset.
func FromEnv(env string, fallback string) Site {
	if wert := os.Getenv(env); wert != "" {
		return New(wert)
	}
	return New(fallback)
}

// Abs builds an absolute URL: the base plus the path (exactly one slash between).
func (site Site) Abs(pfad string) string {
	if pfad == "" {
		return site.base
	}
	if !strings.HasPrefix(pfad, "/") {
		pfad = "/" + pfad
	}
	return site.base + pfad
}

// Base returns the site's base address.
func (site Site) Base() string { return site.base }
