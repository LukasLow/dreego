package dreego

import (
	"regexp"
	"strconv"
	"strings"
)

// validateField applies a `validate:"…"` tag to one field and records messages.
// Rules are comma separated: required, email, min=N, max=N, len=N.
func validateField(form *Form, name string, roh string, regeln string) {
	if regeln == "" {
		return
	}

	beschriftung := form.label(name)
	wert := strings.TrimSpace(roh)

	for _, regel := range strings.Split(regeln, ",") {
		regel = strings.TrimSpace(regel)
		if regel == "" {
			continue
		}

		switch {
		case regel == "required":
			if wert == "" {
				form.SetError(name, beschriftung+" ist ein Pflichtfeld.")
			}

		case regel == "email":
			if wert != "" && !emailPattern.MatchString(wert) {
				form.SetError(name, beschriftung+" ist keine gültige E-Mail-Adresse.")
			}

		case strings.HasPrefix(regel, "min="):
			grenze := zahlNach(regel)
			if grenze >= 0 && len([]rune(wert)) < grenze {
				form.SetError(name,
					beschriftung+" muss mindestens "+strconv.Itoa(grenze)+" Zeichen haben.")
			}

		case strings.HasPrefix(regel, "max="):
			grenze := zahlNach(regel)
			if grenze >= 0 && len([]rune(wert)) > grenze {
				form.SetError(name,
					beschriftung+" darf höchstens "+strconv.Itoa(grenze)+" Zeichen haben.")
			}

		case strings.HasPrefix(regel, "len="):
			grenze := zahlNach(regel)
			if grenze >= 0 && len([]rune(wert)) != grenze {
				form.SetError(name,
					beschriftung+" muss genau "+strconv.Itoa(grenze)+" Zeichen haben.")
			}
		}
	}
}

// zahlNach liest die Zahl hinter einem "=" (z. B. "min=3" -> 3). -1 bei Fehler.
func zahlNach(regel string) int {
	_, rechts, gefunden := strings.Cut(regel, "=")
	if !gefunden {
		return -1
	}
	zahl, err_konv := strconv.Atoi(strings.TrimSpace(rechts))
	if err_konv != nil {
		return -1
	}
	return zahl
}

// emailPattern is intentionally permissive — it rejects the obvious mistakes
// without pretending to prove an address is deliverable.
var emailPattern = regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`)
