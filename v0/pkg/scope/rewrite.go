package scope

import "strings"

// rewriteScoped prefixes every rule selector in the stylesheet with the scope
// prefix, so the CSS applies only inside the scope container.
//
// Handled:
//   - selector lists (comma separated)
//   - nested conditional at-rules: @media, @supports, @container, @layer,
//     @document, @scope
//   - comments and strings
//
// Copied verbatim (selectors there must NOT be prefixed):
//   - @keyframes, @font-face, @page, @property, @counter-style
func rewriteScoped(css string, prefix string) string {
	var ausgabe strings.Builder
	rewriteBlock(&ausgabe, css, prefix)
	return ausgabe.String()
}

func rewriteBlock(ausgabe *strings.Builder, quelle string, prefix string) {
	index := 0
	laenge := len(quelle)

	for index < laenge {
		zeichen := quelle[index]

		if isWhitespace(zeichen) {
			ausgabe.WriteByte(zeichen)
			index++
			continue
		}

		// Kommentar: unverändert übernehmen.
		if zeichen == '/' && index+1 < laenge && quelle[index+1] == '*' {
			ende := strings.Index(quelle[index+2:], "*/")
			if ende < 0 {
				ausgabe.WriteString(quelle[index:])
				return
			}
			absolut := index + 2 + ende + 2
			ausgabe.WriteString(quelle[index:absolut])
			index = absolut
			continue
		}

		// Prelude lesen: bis '{', ';' oder '}'.
		start := index
		for index < laenge && quelle[index] != '{' && quelle[index] != ';' && quelle[index] != '}' {
			index++
		}

		if index >= laenge {
			ausgabe.WriteString(quelle[start:index])
			return
		}

		switch quelle[index] {
		case ';':
			ausgabe.WriteString(quelle[start : index+1])
			index++

		case '}':
			// Sollte auf Blockebene nicht vorkommen; sicherheitshalber ausgeben.
			ausgabe.WriteString(quelle[start:index])
			ausgabe.WriteByte('}')
			index++

		case '{':
			prelude := quelle[start:index]
			block, naechster := readBlock(quelle, index)
			gekuerzt := strings.TrimSpace(prelude)

			if strings.HasPrefix(gekuerzt, "@") {
				ausgabe.WriteString(prelude)
				ausgabe.WriteByte('{')
				if atRuleHasNestedRules(gekuerzt) {
					rewriteBlock(ausgabe, block, prefix)
				} else {
					ausgabe.WriteString(block)
				}
				ausgabe.WriteByte('}')
			} else {
				ausgabe.WriteString(prefixSelectorList(prelude, prefix))
				ausgabe.WriteByte('{')
				ausgabe.WriteString(block)
				ausgabe.WriteByte('}')
			}
			index = naechster
		}
	}
}

// readBlock liefert den Inhalt zwischen einem '{' und dem zugehörigen '}' sowie
// den Index direkt hinter dem schließenden '}'. Strings und Kommentare werden
// übersprungen, damit Klammern darin nicht zählen.
func readBlock(quelle string, oeffnendeKlammer int) (string, int) {
	tiefe := 0
	index := oeffnendeKlammer

	for index < len(quelle) {
		zeichen := quelle[index]
		switch zeichen {
		case '"', '\'':
			index = skipString(quelle, index)
			continue
		case '/':
			if index+1 < len(quelle) && quelle[index+1] == '*' {
				ende := strings.Index(quelle[index+2:], "*/")
				if ende < 0 {
					return quelle[oeffnendeKlammer+1:], len(quelle)
				}
				index = index + 2 + ende + 2
				continue
			}
		case '{':
			tiefe++
		case '}':
			tiefe--
			if tiefe == 0 {
				return quelle[oeffnendeKlammer+1 : index], index + 1
			}
		}
		index++
	}
	return quelle[oeffnendeKlammer+1:], len(quelle)
}

// skipString liefert den Index hinter einem String-Literal, das bei start
// beginnt.
func skipString(quelle string, start int) int {
	quote := quelle[start]
	index := start + 1
	for index < len(quelle) {
		if quelle[index] == '\\' {
			index += 2
			continue
		}
		if quelle[index] == quote {
			return index + 1
		}
		index++
	}
	return len(quelle)
}

var nestedAtRules = map[string]bool{
	"@media":     true,
	"@supports":  true,
	"@container": true,
	"@layer":     true,
	"@document":  true,
	"@scope":     true,
}

// atRuleHasNestedRules sagt, ob in dieser At-Rule weitere Regeln mit Selektoren
// stehen (dann wird hineinrekurriert) oder nur Deklarationen/Schlüsselbilder
// (dann bleibt der Block unverändert).
func atRuleHasNestedRules(prelude string) bool {
	name := prelude
	if grenze := strings.IndexAny(prelude, " \t\r\n("); grenze >= 0 {
		name = prelude[:grenze]
	}
	return nestedAtRules[strings.ToLower(name)]
}

// prefixSelectorList präfixt jeden Selektor einer Selektorliste mit dem Scope.
func prefixSelectorList(prelude string, prefix string) string {
	selektoren := splitTopLevel(prelude, ',')
	var teile []string
	for _, selektor := range selektoren {
		gekuerzt := strings.TrimSpace(selektor)
		if gekuerzt == "" {
			continue
		}
		teile = append(teile, prefix+" "+gekuerzt)
	}
	return strings.Join(teile, ", ")
}

// splitTopLevel trennt bei trenner, ignoriert trenner aber innerhalb von
// Klammern, eckigen Klammern und Strings (z. B. :not(a, b) oder [title="a,b"]).
func splitTopLevel(quelle string, trenner byte) []string {
	var teile []string
	tiefe := 0
	start := 0
	index := 0

	for index < len(quelle) {
		zeichen := quelle[index]
		switch zeichen {
		case '"', '\'':
			index = skipString(quelle, index)
			continue
		case '(', '[':
			tiefe++
		case ')', ']':
			if tiefe > 0 {
				tiefe--
			}
		case trenner:
			if tiefe == 0 {
				teile = append(teile, quelle[start:index])
				start = index + 1
			}
		}
		index++
	}
	teile = append(teile, quelle[start:])
	return teile
}

func isWhitespace(zeichen byte) bool {
	return zeichen == ' ' || zeichen == '\t' || zeichen == '\n' ||
		zeichen == '\r' || zeichen == '\f'
}
