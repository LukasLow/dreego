package scope

import "strings"

// buildScript wraps a component's JavaScript so that
//
//  1. it runs once per component, even if the component appears several times
//     on the page (guarded by the scope id on window.__dreego);
//  2. it runs for every instance of the component, passing the scope root as
//     `root` and catching errors so one component cannot break the page;
//  3. it waits for the DOM before touching the document.
//
// The scope id is a hex string, so embedding it in the JS string literal is
// safe.
func buildScript(js string, scopeID string) string {
	// Defuse a closing script tag inside the component's own JS, so it cannot
	// break out of the <script> element.
	js = escapeClosingTag(js, "script")

	// Remove one trailing newline so the injected code keeps its indentation.
	js = strings.TrimRight(js, "\n")

	return `(function () {
  var dreego = window.__dreego = window.__dreego || {};
  var id = "` + scopeID + `";
  if (dreego[id]) { return; }
  dreego[id] = true;

  function start() {
    var roots = document.querySelectorAll('[data-scope="' + id + '"]');
    for (var i = 0; i < roots.length; i++) {
      try {
        (function (root) {
` + js + `
        })(roots[i]);
      } catch (error) {
        if (window.console) { console.error("[dreego] " + id, error); }
      }
    }
  }

  if (document.readyState === "loading") {
    document.addEventListener("DOMContentLoaded", start);
  } else {
    start();
  }
})();`
}

// escapeClosingTag replaces a closing tag such as "</script" with "<\/script".
// In JavaScript a backslash before a slash is harmless ("<\/script>" equals
// "</script>" as a string, and regex literals stay valid). In CSS, this only
// ever appears in hostile input, so neutralising it is the safe default.
// The comparison is case-insensitive.
func escapeClosingTag(code string, tag string) string {
	nadel := "</" + tag
	var ausgabe strings.Builder
	rest := code

	for {
		position := indexFold(rest, nadel)
		if position < 0 {
			ausgabe.WriteString(rest)
			break
		}
		ausgabe.WriteString(rest[:position])
		ausgabe.WriteString("<\\/")
		ausgabe.WriteString(rest[position+2 : position+2+len(tag)])
		rest = rest[position+2+len(tag):]
	}

	return ausgabe.String()
}

// indexFold liefert den Index des ersten Vorkommens von nadel in text, ohne
// Rücksicht auf Groß- und Kleinschreibung. -1, wenn es nicht vorkommt.
func indexFold(text string, nadel string) int {
	if len(nadel) == 0 || len(nadel) > len(text) {
		return -1
	}
	for start := 0; start <= len(text)-len(nadel); start++ {
		if strings.EqualFold(text[start:start+len(nadel)], nadel) {
			return start
		}
	}
	return -1
}

// escapeAttribute escapes a value that is written into an HTML attribute.
func escapeAttribute(wert string) string {
	ersetzer := strings.NewReplacer(
		"&", "&amp;",
		"<", "&lt;",
		">", "&gt;",
		`"`, "&quot;",
	)
	return ersetzer.Replace(wert)
}
