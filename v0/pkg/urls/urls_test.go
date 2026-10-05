package urls

import (
	"os"
	"testing"
)

func TestAbsMitSchraegstrich(t *testing.T) {
	site := New("https://example.com")
	if got := site.Abs("/pricing"); got != "https://example.com/pricing" {
		t.Fatalf("got %q", got)
	}
}

func TestAbsOhneSchraegstrich(t *testing.T) {
	site := New("https://example.com/")
	if got := site.Abs("pricing"); got != "https://example.com/pricing" {
		t.Fatalf("got %q", got)
	}
}

func TestAbsLeer(t *testing.T) {
	site := New("https://example.com")
	if got := site.Abs(""); got != "https://example.com" {
		t.Fatalf("got %q", got)
	}
}

func TestFromEnvOverride(t *testing.T) {
	t.Setenv("TEST_SITE_URL", "http://localhost:9000")
	site := FromEnv("TEST_SITE_URL", "https://prod.example")
	if got := site.Abs("/x"); got != "http://localhost:9000/x" {
		t.Fatalf("got %q", got)
	}
}

func TestFromEnvFallback(t *testing.T) {
	os.Unsetenv("TEST_SITE_URL_LEER")
	site := FromEnv("TEST_SITE_URL_LEER", "https://prod.example")
	if got := site.Abs("/x"); got != "https://prod.example/x" {
		t.Fatalf("got %q", got)
	}
}
