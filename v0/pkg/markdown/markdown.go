// Package markdown renders a small, safe subset of Markdown to gomponents
// nodes. It exists so legal pages (Impressum, Datenschutz, AGB) can be written
// as text and rendered as HTML — the dreego replacement for Dreego's
// `<body lang="md">`.
//
// The subset is deliberately small: headings, paragraphs, lists, emphasis,
// links, inline code, blockquotes and horizontal rules. Everything is escaped;
// no raw HTML passes through. That keeps the output safe by default.
package markdown

import (
	"strings"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

// ToNodes parses Markdown text into a slice of gomponents nodes.
func ToNodes(quelle string) []g.Node {
	zeilen := strings.Split(strings.ReplaceAll(quelle, "\r\n", "\n"), "\n")

	var knoten []g.Node
	var absatz []string

	abschlussAbsatz := func() {
		if len(absatz) == 0 {
			return
		}
		text := strings.Join(absatz, " ")
		knoten = append(knoten, h.P(inline(text)))
		absatz = absatz[:0]
	}

	for index := 0; index < len(zeilen); index++ {
		zeile := strings.TrimRight(zeilen[index], " \t")
		gekuerzt := strings.TrimSpace(zeile)

		// Leerzeile beendet den Absatz.
		if gekuerzt == "" {
			abschlussAbsatz()
			continue
		}

		// Horizontale Linie.
		if gekuerzt == "---" || gekuerzt == "***" {
			abschlussAbsatz()
			knoten = append(knoten, h.Hr())
			continue
		}

		// Überschrift: # … ######
		if stufe := ueberschriftStufe(gekuerzt); stufe > 0 {
			abschlussAbsatz()
			text := strings.TrimSpace(gekuerzt[stufe+1:])
			knoten = append(knoten, ueberschrift(stufe, text))
			continue
		}

		// Liste: - / * / +
		if istListenpunkt(gekuerzt) {
			abschlussAbsatz()
			eintraege, verbraucht := sammleListe(zeilen, index)
			knoten = append(knoten, h.Ul(g.Map(eintraege, func(eintrag string) g.Node {
				return h.Li(inline(eintrag))
			})))
			index += verbraucht - 1
			continue
		}

		// Zitat: >
		if strings.HasPrefix(gekuerzt, "> ") {
			abschlussAbsatz()
			knoten = append(knoten, h.BlockQuote(inline(strings.TrimPrefix(gekuerzt, "> "))))
			continue
		}

		absatz = append(absatz, gekuerzt)
	}

	abschlussAbsatz()
	return knoten
}

// sichereSchemas sind die Link-Schemata, die gerendert werden. Alles andere
// (insbesondere javascript:, data:, vbscript:) wird verworfen und nur als Text
// ausgegeben. Relative Links (ohne Schema) sind immer erlaubt.
var sichereSchemas = map[string]bool{
	"http":   true,
	"https":  true,
	"mailto": true,
	"tel":    true,
	"":       true, // relativ / fragment
}

// sichererLink meldet, ob ein Link-Ziel gerendert werden darf.
func sichererLink(ziel string) bool {
	gekuerzt := strings.TrimSpace(ziel)
	doppelpunkt := strings.IndexByte(gekuerzt, ':')
	if doppelpunkt < 0 {
		return true // relativ oder Fragment
	}

	// Nur als Schema zählen, wenn vor dem Doppelpunkt keine '/', '?' oder '#'
	// stehen (sonst ist es ein relativer Pfad mit Doppelpunkt).
	schema := strings.ToLower(gekuerzt[:doppelpunkt])
	if strings.ContainsAny(schema, "/?#") {
		return true
	}

	return sichereSchemas[schema]
}

func ueberschriftStufe(zeile string) int {
	stufe := 0
	for stufe < len(zeile) && zeile[stufe] == '#' {
		stufe++
	}
	if stufe == 0 || stufe > 6 {
		return 0
	}
	if stufe < len(zeile) && zeile[stufe] != ' ' {
		return 0
	}
	return stufe
}

func ueberschrift(stufe int, text string) g.Node {
	inhalt := inline(text)
	switch stufe {
	case 1:
		return h.H1(inhalt)
	case 2:
		return h.H2(inhalt)
	case 3:
		return h.H3(inhalt)
	case 4:
		return h.H4(inhalt)
	case 5:
		return h.H5(inhalt)
	default:
		return h.H6(inhalt)
	}
}

func istListenpunkt(zeile string) bool {
	if len(zeile) < 2 {
		return false
	}
	if zeile[0] != '-' && zeile[0] != '*' && zeile[0] != '+' {
		return false
	}
	return zeile[1] == ' '
}

func sammleListe(zeilen []string, start int) ([]string, int) {
	var eintraege []string
	index := start
	for index < len(zeilen) {
		gekuerzt := strings.TrimSpace(zeilen[index])
		if !istListenpunkt(gekuerzt) {
			break
		}
		eintraege = append(eintraege, strings.TrimSpace(gekuerzt[2:]))
		index++
	}
	return eintraege, index - start
}

// inline parses inline Markdown (**bold**, *italic*, `code`, [text](url)).
func inline(text string) g.Node {
	return g.Group(inlineNodes(text))
}

func inlineNodes(text string) []g.Node {
	var knoten []g.Node
	puffer := strings.Builder{}

	absaugen := func() {
		if puffer.Len() > 0 {
			knoten = append(knoten, g.Text(puffer.String()))
			puffer.Reset()
		}
	}

	index := 0
	for index < len(text) {
		rest := text[index:]

		// Inline-Code `…`
		if rest[0] == '`' {
			if ende := strings.IndexByte(rest[1:], '`'); ende >= 0 {
				absaugen()
				knoten = append(knoten, h.Code(g.Text(rest[1:1+ende])))
				index += ende + 2
				continue
			}
		}

		// Fett **…**
		if strings.HasPrefix(rest, "**") {
			if ende := strings.Index(rest[2:], "**"); ende >= 0 {
				absaugen()
				knoten = append(knoten, h.Strong(g.Text(rest[2:2+ende])))
				index += ende + 4
				continue
			}
		}

		// Kursiv *…*
		if rest[0] == '*' {
			if ende := strings.IndexByte(rest[1:], '*'); ende >= 0 {
				absaugen()
				knoten = append(knoten, h.Em(g.Text(rest[1:1+ende])))
				index += ende + 2
				continue
			}
		}

		// Link [text](url)
		if rest[0] == '[' {
			if textEnde := strings.IndexByte(rest, ']'); textEnde > 0 {
				nachBracket := rest[textEnde+1:]
				if strings.HasPrefix(nachBracket, "(") {
					if urlEnde := strings.IndexByte(nachBracket, ')'); urlEnde > 0 {
						absaugen()
						label := rest[1:textEnde]
						ziel := nachBracket[1:urlEnde]
						if sichererLink(ziel) {
							knoten = append(knoten, h.A(h.Href(ziel), g.Text(label)))
						} else {
							// Unsicheres Schema (z.B. javascript:): nur als Text.
							knoten = append(knoten, g.Text(label))
						}
						index += textEnde + 1 + urlEnde + 1
						continue
					}
				}
			}
		}

		puffer.WriteByte(rest[0])
		index++
	}

	absaugen()
	return knoten
}
