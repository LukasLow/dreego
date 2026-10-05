package dreego

import (
	"net/http"

	"github.com/LukasLow/dreego/v0/pkg/scope"
)

// NewCtx builds a bare request context with a fresh asset collector and nonce,
// but no session store. It is meant for tests and for rendering a component
// outside the running server (e.g. building a fragment).
//
// In normal request handling the framework calls newCtx itself.
func NewCtx(r *http.Request) *Ctx {
	nonce := scope.NewNonce()
	collector := scope.New()
	collector.SetNonce(nonce)

	return &Ctx{
		r:         r,
		nonce:     nonce,
		collector: collector,
		locale:    "de",
		session:   map[string]string{},
		data:      map[string]any{},
	}
}

// NewForm builds a Form from explicit values and errors. Use it when you need a
// form's Old/Errors without an incoming request (tests, pre-filled forms).
func NewForm(values map[string]string, errors map[string][]string) *Form {
	form := &Form{
		values: map[string]string{},
		errors: map[string][]string{},
		labels: map[string]string{},
	}
	for name, wert := range values {
		form.values[name] = wert
	}
	for name, meldungen := range errors {
		form.errors[name] = meldungen
	}
	return form
}
