// Command gen_dom generates dreego's own dom package files (elements.go,
// attributes.go) from a source of HTML helper code. One-off authoring tool:
// run it, review the output, commit the result. It does NOT run at build time.
//
// Usage:
//
//	go run ./internal/gen_dom <inputdir> <outputdir>
//
// The input is a directory with elements.go and attributes.go in the shape of
// plain forwarding functions (func Name(...) T { return El("x", …) }).
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

var funcRe = regexp.MustCompile(`(?m)^func ([A-Za-z0-9]+)\(([^)]*)\) (\S+) \{\n\treturn (.*)\n\}`)

var elRe = regexp.MustCompile(`El\("([a-z0-9]+)"`)
var attrRe = regexp.MustCompile(`Attr\("([a-z0-9-]+)"(?:,\s*v)?\)`)

type attribut struct {
	htmlName string
	mitWert  bool
}

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "usage: gen_dom <inputdir> <outputdir>")
		os.Exit(1)
	}
	elemente := map[string]string{}    // Go-Name -> HTML-Name
	attribute := map[string]attribut{} // Go-Name -> Attribut

	lese(filepath.Join(os.Args[1], "elements.go"), func(name, body string) {
		if treffer := elRe.FindStringSubmatch(body); treffer != nil {
			elemente[name] = treffer[1]
		}
	})
	lese(filepath.Join(os.Args[1], "attributes.go"), func(name, body string) {
		if treffer := attrRe.FindStringSubmatch(body); treffer != nil {
			attribute[name] = attribut{
				htmlName: treffer[1],
				mitWert:  strings.Contains(body, ","), // "Attr("x", v)"
			}
		}
	})

	for _, name := range []string{"CiteEl", "DataAttr", "FormEl", "LabelEl", "StyleAttr", "TitleAttr"} {
		delete(elemente, name)
		delete(attribute, name)
	}

	schreibeElemente(filepath.Join(os.Args[2], "elements.go"), elemente)
	schreibeAttribute(filepath.Join(os.Args[2], "attributes.go"), attribute)
	fmt.Printf("elemente=%d attribute=%d\n", len(elemente), len(attribute))
}

func lese(pfad string, pro func(name, body string)) {
	daten, err_lesen := os.ReadFile(pfad)
	if err_lesen != nil {
		fmt.Fprintln(os.Stderr, "lesen:", err_lesen)
		os.Exit(1)
	}
	for _, treffer := range funcRe.FindAllStringSubmatch(string(daten), -1) {
		pro(treffer[1], strings.TrimSpace(treffer[4]))
	}
}

func schreibeElemente(pfad string, elemente map[string]string) {
	var b strings.Builder
	kopf(&b, "HTML-Element-Helfer — erzeugt von internal/gen_dom.")
	for _, goName := range sortiert(elemente) {
		htmlName := elemente[goName]
		fmt.Fprintf(&b, "// %s rendert ein <%s>-Element.\n", goName, htmlName)
		fmt.Fprintf(&b, "func %s(children ...View) View {\n\treturn El(%q, children...)\n}\n\n", goName, htmlName)
	}
	schreibe(pfad, b.String())
}

func schreibeAttribute(pfad string, attribute map[string]attribut) {
	var b strings.Builder
	kopf(&b, "HTML-Attribut-Helfer — erzeugt von internal/gen_dom.")

	namen := make([]string, 0, len(attribute))
	for name := range attribute {
		namen = append(namen, name)
	}
	sort.Strings(namen)

	for _, goName := range namen {
		at := attribute[goName]
		if at.mitWert {
			fmt.Fprintf(&b, "// %s setzt das %s-Attribut.\n", goName, at.htmlName)
			fmt.Fprintf(&b, "func %s(v string) View {\n\treturn Attr(%q, v)\n}\n\n", goName, at.htmlName)
		} else {
			fmt.Fprintf(&b, "// %s setzt das boolesche %s-Attribut.\n", goName, at.htmlName)
			fmt.Fprintf(&b, "func %s() View {\n\treturn Attr(%q)\n}\n\n", goName, at.htmlName)
		}
	}
	schreibe(pfad, b.String())
}

func sortiert(menge map[string]string) []string {
	namen := make([]string, 0, len(menge))
	for name := range menge {
		namen = append(namen, name)
	}
	sort.Strings(namen)
	return namen
}

func kopf(b *strings.Builder, beschreibung string) {
	b.WriteString("// Code erzeugt von internal/gen_dom. NICHT von Hand editieren.\n")
	fmt.Fprintf(b, "// %s\n\npackage dom\n\n", beschreibung)
}

func schreibe(pfad, inhalt string) {
	err_schreiben := os.WriteFile(pfad, []byte(inhalt), 0o644)
	if err_schreiben != nil {
		fmt.Fprintln(os.Stderr, "schreiben:", err_schreiben)
		os.Exit(1)
	}
}
