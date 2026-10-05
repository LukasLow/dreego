package dreego

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	g "github.com/LukasLow/dreego/v0/pkg/dom"
)

func TestContextSetGetData(t *testing.T) {
	anfrage := httptest.NewRequest("GET", "/?x=1", nil)
	c := NewCtx(anfrage)

	c.Set("nav", "start")
	if c.Get("nav") != "start" {
		t.Fatalf("Get = %q", c.Get("nav"))
	}
	if c.Data("nav") != "start" {
		t.Fatalf("Data = %v", c.Data("nav"))
	}

	c.Delete("nav")
	if c.Get("nav") != "" {
		t.Fatalf("nach Delete noch da: %q", c.Get("nav"))
	}
}

func TestContextQuery(t *testing.T) {
	c := NewCtx(httptest.NewRequest("GET", "/?email=a@b.de", nil))
	if c.Query("email") != "a@b.de" {
		t.Fatalf("Query = %q", c.Query("email"))
	}
}

func TestContextUsesRequestContext(t *testing.T) {
	c := NewCtx(httptest.NewRequest("GET", "/", nil))
	if c.Context() == nil {
		t.Fatal("Context() ist nil")
	}
}

func TestBindHAengtFormAnCtx(t *testing.T) {
	type eingabe struct {
		Email string `form:"email" validate:"required,email"`
	}
	anfrage := httptest.NewRequest("POST", "/", strings.NewReader("email=kaputt"))
	anfrage.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	c := NewCtx(anfrage)

	_, _, err_bind := Bind[eingabe](c)
	if err_bind != nil {
		t.Fatalf("bind: %v", err_bind)
	}

	// Old/Err arbeiten jetzt über den Ctx.
	if c.Old("email") != "kaputt" {
		t.Fatalf("Old = %q", c.Old("email"))
	}
	if c.Err("email") == "" {
		t.Fatal("Err liefert keine Meldung")
	}
	if c.Form() == nil {
		t.Fatal("Form() ist nil")
	}
}

func TestContextImHandler(t *testing.T) {
	app := NewApp()
	app.Page(Page{Path: "/", Get: func(c *Ctx) g.View {
		c.Set("nav", "home")
		return g.Text("nav=" + c.Get("nav"))
	}})
	schreiber := aufruf(app, "GET", "/", "")
	if !strings.Contains(schreiber.Body.String(), "nav=home") {
		t.Fatalf("ctx im handler: %s", schreiber.Body.String())
	}
}

var _ = http.StatusOK
