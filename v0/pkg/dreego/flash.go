package dreego

// Flash stores a one-shot message in the session. The next request reads it
// with FlashGet and the message is consumed.
//
// Typical Post/Redirect/Get use:
//
//	c.Flash("login_error", "Bitte Eingaben prüfen.")
//	return dreego.Redirect("/login", 303)
func (c *Ctx) Flash(key string, message string) {
	c.sessionSet("flash_"+key, message)
}

// FlashGet returns and clears the flash message for key (empty if none).
func (c *Ctx) FlashGet(key string) string {
	wert := c.sessionGet("flash_" + key)
	if wert == "" {
		return ""
	}
	c.sessionDel("flash_" + key)
	return wert
}

// SessionVal reads a value from the session.
func (c *Ctx) SessionVal(key string) string {
	return c.sessionGet(key)
}

// SetSessionVal writes a value into the session.
func (c *Ctx) SetSessionVal(key string, val string) {
	c.sessionSet(key, val)
}

// DelSessionVal removes a value from the session.
func (c *Ctx) DelSessionVal(key string) {
	c.sessionDel(key)
}

// DestroySession clears the whole session (logout).
func (c *Ctx) DestroySession() {
	for key := range c.session {
		delete(c.session, key)
	}
	c.sessionDirty = true
}
