package dreego

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

// Style registers per-request CSS declarations under a generated class name and
// returns that class name. The rule is emitted once per page in a nonce-tagged
// <style> block. Use it for values that must not live in a style="…" attribute,
// which the built-in CSP (style-src 'nonce-…') blocks:
//
//	cls := c.Style("width: " + strconv.Itoa(percent) + "%")
//	Div(Class("bar " + cls))
//
// The same declarations always yield the same class name, so repeated use is
// deduplicated. An empty declarations string returns "" (no rule, no class).
func (c *Ctx) Style(declarations string) string {
	if strings.TrimSpace(declarations) == "" {
		return ""
	}

	name := "d-" + kurzerHash(declarations)
	if c.collector != nil {
		c.collector.AddPageCSS("." + name + " {" + declarations + "}")
	}
	return name
}

// kurzerHash returns a short, stable id for a declarations string.
func kurzerHash(quelle string) string {
	summe := sha256.Sum256([]byte(quelle))
	return hex.EncodeToString(summe[:4])
}
