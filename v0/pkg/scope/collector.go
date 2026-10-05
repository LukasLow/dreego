package scope

import (
	"crypto/rand"
	"encoding/base64"
	"strings"
	"sync"

	g "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

// asset is one component's scoped CSS and JS, keyed by its scope id.
type asset struct {
	css string
	js  string
}

// Collector gathers the assets of one page (one HTTP response), so every
// component's CSS and JS is emitted exactly once — no matter how many times the
// component appears on the page. It also carries the CSP nonce for the response.
//
// One Collector per request. Safe for concurrent use.
type Collector struct {
	mu     sync.Mutex
	nonce  string
	assets map[string]asset
	order  []string

	// critical holds CSS that must be inlined in <head> before any external
	// stylesheet, so the first paint already has the base colours and font.
	critical     []string
	criticalSeen map[string]bool
}

// New returns a fresh, empty Collector — one per request.
func New() *Collector {
	return &Collector{assets: map[string]asset{}, criticalSeen: map[string]bool{}}
}

// SetNonce sets the CSP nonce that every emitted <style> and <script> carries.
func (c *Collector) SetNonce(nonce string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.nonce = nonce
}

// Nonce returns the configured CSP nonce (empty if none was set).
func (c *Collector) Nonce() string {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.nonce
}

// NewNonce returns a fresh random nonce suitable for a Content-Security-Policy.
// The server puts it into the CSP header and passes the same value to SetNonce.
// On the (catastrophic) failure of crypto/rand it returns "" — the caller
// should treat an empty nonce as a hard error, not as "no nonce".
func NewNonce() string {
	var roh [16]byte
	_, err_lesen := rand.Read(roh[:])
	if err_lesen != nil {
		return ""
	}
	return base64.RawStdEncoding.EncodeToString(roh[:])
}

func (c *Collector) add(id string, eintrag asset) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, vorhanden := c.assets[id]; vorhanden {
		return
	}
	c.assets[id] = eintrag
	c.order = append(c.order, id)
}

// Box registers the component's CSS and JS with the collector and returns its
// scope container. The assets themselves are emitted once by Styles and Scripts.
func (c *Collector) Box(parts ...g.Node) g.Node {
	css, js, body := splitParts(parts)
	scopeID := shortHash(css + "\x00" + js)
	c.add(scopeID, asset{css: css, js: js})
	return Div(g.Attr("data-scope", scopeID), g.Group(body))
}

// Styles renders every collected stylesheet once, scoped and nonce-tagged.
func (c *Collector) Styles() g.Node {
	c.mu.Lock()
	defer c.mu.Unlock()

	var teile []g.Node
	for _, id := range c.order {
		eintrag := c.assets[id]
		if strings.TrimSpace(eintrag.css) == "" {
			continue
		}
		scoped := rewriteScoped(eintrag.css, `[data-scope="`+id+`"]`)
		teile = append(teile, StyleEl(c.withNonce(g.Raw(escapeClosingTag(scoped, "style")))))
	}
	return g.Group(teile)
}

// Scripts renders every collected script once, scoped and nonce-tagged.
func (c *Collector) Scripts() g.Node {
	c.mu.Lock()
	defer c.mu.Unlock()

	var teile []g.Node
	for _, id := range c.order {
		eintrag := c.assets[id]
		if strings.TrimSpace(eintrag.js) == "" {
			continue
		}
		teile = append(teile, Script(c.withNonce(g.Raw(escapeClosingTag(buildScript(eintrag.js, id), "script")))))
	}
	return g.Group(teile)
}

// withNonce prepends a nonce attribute when one is set. Caller holds the lock.
func (c *Collector) withNonce(inhalt g.Node) g.Node {
	if c.nonce == "" {
		return inhalt
	}
	return g.Group{g.Attr("nonce", c.nonce), inhalt}
}

// AddCritical registers CSS that is inlined in <head> before any external
// stylesheet. Use it for the base colours and font so the first paint is already
// themed — no white flash while /public/style.css loads. Duplicate blocks are
// ignored.
func (c *Collector) AddCritical(css string) {
	if strings.TrimSpace(css) == "" {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.criticalSeen[css] {
		return
	}
	c.criticalSeen[css] = true
	c.critical = append(c.critical, css)
}

// Critical renders the inlined critical CSS, one nonce-tagged <style> block.
// An empty result renders nothing.
func (c *Collector) Critical() g.Node {
	c.mu.Lock()
	defer c.mu.Unlock()

	if len(c.critical) == 0 {
		return g.Group(nil)
	}

	var baukasten strings.Builder
	for _, block := range c.critical {
		baukasten.WriteString(block)
		baukasten.WriteByte('\n')
	}

	return StyleEl(c.withNonce(g.Raw(escapeClosingTag(baukasten.String(), "style"))))
}
