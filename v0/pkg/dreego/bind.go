package dreego

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// Form carries the submitted values and any field errors, so a handler can
// re-render a form with the user's input preserved and messages attached.
//
// It is the Go replacement for Dreego's c.Old(...) and c.Errors(...).
type Form struct {
	values map[string]string
	errors map[string][]string
	labels map[string]string
}

// newForm reads the request's form values (query and body).
func newForm(c *Ctx) (*Form, error) {
	form := &Form{
		values: map[string]string{},
		errors: map[string][]string{},
		labels: map[string]string{},
	}

	err_parse := c.r.ParseForm()
	if err_parse != nil {
		return form, err_parse
	}

	for name, liste := range c.r.Form {
		if len(liste) == 0 {
			continue
		}
		form.values[name] = liste[0]
	}

	return form, nil
}

// Old returns the value the user submitted for name (empty if none). Use it to
// pre-fill a form after an invalid submission.
func (form *Form) Old(name string) string { return form.values[name] }

// Errors returns the messages for a field (nil if the field is fine).
func (form *Form) Errors(name string) []string { return form.errors[name] }

// HasErrors reports whether any field failed validation.
func (form *Form) HasErrors() bool { return len(form.errors) > 0 }

// SetError attaches a message to a field. Use it for domain errors that are not
// plain validation (e.g. a rejected login).
func (form *Form) SetError(name string, message string) {
	form.errors[name] = append(form.errors[name], message)
}

// label returns the human label for a field, falling back to the field name.
func (form *Form) label(name string) string {
	if beschriftung := form.labels[name]; beschriftung != "" {
		return beschriftung
	}
	return name
}

// Bind reads a form into a typed struct T and validates it.
//
// Struct tags drive it:
//
//	type LoginForm struct {
//	    Email string `form:"email" validate:"required,email" label:"E-Mail"`
//	}
//
// Returns the filled struct, the Form (with Old values and Errors), and an
// error only for a broken request or a misused struct — never for a validation
// failure (that is part of the Form, not an error).
//
//	in, form, err_bind := dreego.Bind[LoginForm](c)
//	if err_bind != nil { /* programming/request error */ }
//	if form.HasErrors() { /* re-render with form.Old / form.Errors */ }
func Bind[T any](c *Ctx) (T, *Form, error) {
	var leer T

	form, err_form := newForm(c)
	if err_form != nil {
		return leer, form, err_form
	}

	typ := reflect.TypeOf(leer)
	wert := reflect.ValueOf(&leer).Elem()

	if typ.Kind() != reflect.Struct {
		return leer, form, fmt.Errorf("dreego.Bind: %T ist keine Struktur", leer)
	}

	for index := 0; index < typ.NumField(); index++ {
		feldTyp := typ.Field(index)
		feldWert := wert.Field(index)

		if !feldWert.CanSet() {
			continue
		}

		name := feldTyp.Tag.Get("form")
		if name == "" {
			name = strings.ToLower(feldTyp.Name[:1]) + feldTyp.Name[1:]
		}

		form.labels[name] = feldTyp.Tag.Get("label")

		roh := form.values[name]
		if err_setz := setField(feldWert, roh); err_setz != nil {
			return leer, form, err_setz
		}

		validateField(form, name, roh, feldTyp.Tag.Get("validate"))
	}

	return leer, form, nil
}

// setField writes a raw string into a struct field, converting by kind.
func setField(feld reflect.Value, roh string) error {
	switch feld.Kind() {
	case reflect.String:
		feld.SetString(roh)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if roh == "" {
			feld.SetInt(0)
			return nil
		}
		zahl, err_konv := strconv.ParseInt(roh, 10, 64)
		if err_konv != nil {
			return fmt.Errorf("dreego.Bind: %q ist keine ganze Zahl", roh)
		}
		feld.SetInt(zahl)
	case reflect.Float32, reflect.Float64:
		if roh == "" {
			feld.SetFloat(0)
			return nil
		}
		zahl, err_konv := strconv.ParseFloat(roh, 64)
		if err_konv != nil {
			return fmt.Errorf("dreego.Bind: %q ist keine Zahl", roh)
		}
		feld.SetFloat(zahl)
	case reflect.Bool:
		feld.SetBool(roh == "on" || roh == "true" || roh == "1")
	default:
		return fmt.Errorf("dreego.Bind: Feldtyp %s wird nicht unterstützt", feld.Kind())
	}
	return nil
}
