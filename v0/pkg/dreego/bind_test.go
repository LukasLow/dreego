package dreego

import (
	"net/http/httptest"
	"strings"
	"testing"
)

type testForm struct {
	Email string `form:"email" validate:"required,email" label:"E-Mail"`
	Team  string `form:"team" validate:"required,min=3" label:"Team"`
	Alter int    `form:"alter"`
}

func ctxMitForm(body string) *Ctx {
	anfrage := httptest.NewRequest("POST", "/x", strings.NewReader(body))
	anfrage.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return &Ctx{r: anfrage}
}

func TestBindFuelltWerte(t *testing.T) {
	c := ctxMitForm("email=a@b.de&team=ACME&alter=7")
	in, form, err_bind := Bind[testForm](c)
	if err_bind != nil {
		t.Fatalf("bind: %v", err_bind)
	}
	if in.Email != "a@b.de" || in.Team != "ACME" || in.Alter != 7 {
		t.Fatalf("falsch gefuellt: %+v", in)
	}
	if form.HasErrors() {
		t.Fatalf("unerwartete Fehler: %+v", form.errors)
	}
}

func TestBindValidierung(t *testing.T) {
	c := ctxMitForm("email=kaputt&team=x")
	_, form, err_bind := Bind[testForm](c)
	if err_bind != nil {
		t.Fatalf("bind: %v", err_bind)
	}
	if !form.HasErrors() {
		t.Fatal("Fehler erwartet")
	}
	if len(form.Errors("email")) == 0 {
		t.Fatal("email-Fehler fehlt")
	}
	if len(form.Errors("team")) == 0 {
		t.Fatal("team-Fehler fehlt")
	}
}

func TestOldBehaeltEingabe(t *testing.T) {
	c := ctxMitForm("email=kaputt&team=x")
	_, form, _ := Bind[testForm](c)
	if form.Old("email") != "kaputt" {
		t.Fatalf("Old = %q", form.Old("email"))
	}
	if form.Old("gibtsnicht") != "" {
		t.Fatal("unbekanntes Feld soll leer sein")
	}
}

func TestPflichtfeldLeer(t *testing.T) {
	c := ctxMitForm("")
	_, form, _ := Bind[testForm](c)
	if len(form.Errors("email")) == 0 {
		t.Fatal("required-Fehler fuer email fehlt")
	}
}

func TestLabelInMeldung(t *testing.T) {
	c := ctxMitForm("email=kaputt&team=x")
	_, form, _ := Bind[testForm](c)
	meldungen := form.Errors("email")
	if len(meldungen) == 0 {
		t.Fatal("keine Meldung")
	}
	if meldungen[0][:6] != "E-Mail" {
		t.Fatalf("Label fehlt in Meldung: %q", meldungen[0])
	}
}
