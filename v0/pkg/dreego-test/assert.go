package dreegotest

import (
	"strings"
	"testing"
)

// AssertStatus fails the test when the status code differs.
func (r *Response) AssertStatus(t *testing.T, erwartet int) {
	t.Helper()
	if r.Status != erwartet {
		t.Fatalf("Status = %d, erwartet %d\nBody: %s", r.Status, erwartet, r.Body)
	}
}

// AssertContains fails the test when the body does not contain the substring.
func (r *Response) AssertContains(t *testing.T, teil string) {
	t.Helper()
	if !strings.Contains(r.Body, teil) {
		t.Fatalf("Body enthält %q nicht\nBody: %s", teil, r.Body)
	}
}

// AssertNotContains fails the test when the body contains the substring.
func (r *Response) AssertNotContains(t *testing.T, teil string) {
	t.Helper()
	if strings.Contains(r.Body, teil) {
		t.Fatalf("Body enthält %q, sollte nicht\nBody: %s", teil, r.Body)
	}
}

// AssertHeader fails the test when a header value differs.
func (r *Response) AssertHeader(t *testing.T, name string, erwartet string) {
	t.Helper()
	if got := r.Header.Get(name); got != erwartet {
		t.Fatalf("Header %s = %q, erwartet %q", name, got, erwartet)
	}
}

// AssertLocation fails the test when the Location header differs.
func (r *Response) AssertLocation(t *testing.T, erwartet string) {
	t.Helper()
	if got := r.Header.Get("Location"); got != erwartet {
		t.Fatalf("Location = %q, erwartet %q", got, erwartet)
	}
}

// Status returns the status code (chainable convenience).
func (r *Response) StatusCode() int { return r.Status }
