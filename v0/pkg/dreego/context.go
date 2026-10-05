package dreego

import (
	"context"
)

// The request context (state) layer. dreego's Ctx is the per-request "bag" a
// handler and its components share: a handler sets a value (e.g. the active nav
// item), a component reads it. This mirrors the old Dreego Context.Set/Get/Data.
//
// Keys are plain strings; values are any. Keep keys few and obvious — a typed
// struct on the handler is often clearer than many loose keys.

// Context returns the standard library context of the request. It carries
// cancellation and deadlines and is what you pass to database or mail calls.
func (c *Ctx) Context() context.Context { return c.r.Context() }

// Set stores a value on the request context.
func (c *Ctx) Set(key string, value any) {
	if c.data == nil {
		c.data = map[string]any{}
	}
	c.data[key] = value
}

// Data returns a stored value as any (nil if absent).
func (c *Ctx) Data(key string) any {
	if c.data == nil {
		return nil
	}
	return c.data[key]
}

// Get returns a stored value as string ("" if absent or not a string).
func (c *Ctx) Get(key string) string {
	wert, _ := c.Data(key).(string)
	return wert
}

// Delete removes a stored value.
func (c *Ctx) Delete(key string) {
	if c.data == nil {
		return
	}
	delete(c.data, key)
}

// Query returns a URL query parameter (e.g. "?email=x" -> Query("email")).
func (c *Ctx) Query(name string) string { return c.r.URL.Query().Get(name) }

// Old returns the value the client submitted for a form field, so a form can be
// re-rendered with the user's input after a validation error (see Bind/Form).
func (c *Ctx) Old(name string) string {
	if c.form == nil {
		return ""
	}
	return c.form.Old(name)
}

// Err returns the first validation message for a field, or "".
func (c *Ctx) Err(name string) string {
	if c.form == nil {
		return ""
	}
	meldungen := c.form.Errors(name)
	if len(meldungen) == 0 {
		return ""
	}
	return meldungen[0]
}

// SetForm attaches a parsed form (with Old/Errors) to the context, so Old and
// Err work after Bind without threading the form through every component.
func (c *Ctx) SetForm(form *Form) { c.form = form }

// Form returns the attached form, or nil.
func (c *Ctx) Form() *Form { return c.form }
