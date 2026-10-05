package scope

import (
	"bytes"
	"io"

	g "github.com/LukasLow/dreego/v0/pkg/dom"
)

// documentNode renders a full HTML document. It renders the body FIRST so that
// every component registers its assets, then writes <head> with, in order:
//
//  1. the caller's head content (meta, title, external stylesheet links),
//  2. the inlined critical CSS (base colours + font — painted immediately),
//  3. the collected component styles.
//
// The scripts go at the end of <body>.
type documentNode struct {
	c    *Collector
	lang string
	head g.View
	body g.View
}

func (dokument documentNode) Render(w io.Writer) error {
	// Body first: rendering it registers all component assets in the collector.
	var puffer bytes.Buffer
	err_body := dokument.body.Render(&puffer)
	if err_body != nil {
		return err_body
	}

	var ausgabe bytes.Buffer

	ausgabe.WriteString(`<!DOCTYPE html><html lang="`)
	ausgabe.WriteString(escapeAttribute(dokument.lang))
	ausgabe.WriteString(`"><head>`)

	err_head := dokument.head.Render(&ausgabe)
	if err_head != nil {
		return err_head
	}

	err_critical := dokument.c.Critical().Render(&ausgabe)
	if err_critical != nil {
		return err_critical
	}

	err_styles := dokument.c.Styles().Render(&ausgabe)
	if err_styles != nil {
		return err_styles
	}

	ausgabe.WriteString(`</head><body>`)
	ausgabe.Write(puffer.Bytes())

	err_scripts := dokument.c.Scripts().Render(&ausgabe)
	if err_scripts != nil {
		return err_scripts
	}

	ausgabe.WriteString(`</body></html>`)

	_, err_schreiben := ausgabe.WriteTo(w)
	return err_schreiben
}

// Document renders a full HTML page with deduplicated component assets.
//
// The body is rendered first so components register themselves; their styles are
// placed in <head> and their scripts at the end of <body>. Each component's CSS
// and JS appears exactly once, carrying the collector's CSP nonce.
//
// Critical base CSS set with Collector.AddCritical is inlined right after head,
// before external stylesheets, to avoid a white flash on first paint.
//
//	c.AddCritical(`:root { color-scheme: light } body { background: #fffdf7 }`)
//	page := scope.Document(c, "de", Head(TitleEl(g.Text("…"))), Body(zaehler(c)))
func Document(c *Collector, lang string, head g.View, body g.View) g.View {
	return documentNode{c: c, lang: lang, head: head, body: body}
}
