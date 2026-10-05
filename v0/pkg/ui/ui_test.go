package ui

import (
	"bytes"
	"net/http/httptest"
	"strings"
	"testing"

	g "maragu.dev/gomponents"

	"github.com/LukasLow/dreego/v0/pkg/dreego"
)

func testCtx(t *testing.T) *dreego.Ctx {
	t.Helper()
	return dreego.NewCtx(httptest.NewRequest("GET", "/", nil))
}

func render(t *testing.T, knoten g.Node) string {
	t.Helper()
	var puffer bytes.Buffer
	err_render := knoten.Render(&puffer)
	if err_render != nil {
		t.Fatalf("render: %v", err_render)
	}
	return puffer.String()
}

func TestButtonPrimary(t *testing.T) {
	html := render(t, Button(testCtx(t), "/pricing", "Preise", Primary))
	if !strings.Contains(html, `href="/pricing"`) {
		t.Fatalf("href fehlt: %s", html)
	}
	if !strings.Contains(html, "Preise") {
		t.Fatalf("label fehlt: %s", html)
	}
	if !strings.Contains(html, "u-btn-primary") {
		t.Fatalf("primary-klasse fehlt: %s", html)
	}
}

func TestButtonGhost(t *testing.T) {
	html := render(t, Button(testCtx(t), "/x", "X", Ghost))
	if !strings.Contains(html, "u-btn-ghost") {
		t.Fatalf("ghost-klasse fehlt: %s", html)
	}
}

func TestCardTitelUndBody(t *testing.T) {
	html := render(t, Card(testCtx(t), "Projekte", g.Text("Inhalt")))
	if !strings.Contains(html, "<h3>Projekte</h3>") || !strings.Contains(html, "Inhalt") {
		t.Fatalf("card inhalt fehlt: %s", html)
	}
}

func TestBadgeTones(t *testing.T) {
	for tone, klasse := range map[Tone]string{
		Info:    "u-badge-info",
		Success: "u-badge-success",
		Warning: "u-badge-warning",
		Danger:  "u-badge-danger",
	} {
		html := render(t, Badge(testCtx(t), "x", tone))
		if !strings.Contains(html, klasse) {
			t.Fatalf("tone %s fehlt: %s", tone, html)
		}
	}
}

func TestAlertRoleJeTone(t *testing.T) {
	danger := render(t, Alert(testCtx(t), "kaputt", Danger))
	if !strings.Contains(danger, `role="alert"`) {
		t.Fatalf("danger soll role=alert haben: %s", danger)
	}
	info := render(t, Alert(testCtx(t), "info", Info))
	if !strings.Contains(info, `role="status"`) {
		t.Fatalf("info soll role=status haben: %s", info)
	}
}

func TestFieldMitForm(t *testing.T) {
	form := dreego.NewForm(
		map[string]string{"email": "kaputt"},
		map[string][]string{"email": {"keine E-Mail"}},
	)
	html := render(t, Field(testCtx(t), form, "email", "E-Mail", "email"))

	if !strings.Contains(html, `for="email"`) || !strings.Contains(html, `id="email"`) {
		t.Fatalf("label/input verknüpfung fehlt: %s", html)
	}
	if !strings.Contains(html, `value="kaputt"`) {
		t.Fatalf("old-wert fehlt: %s", html)
	}
	if !strings.Contains(html, "keine E-Mail") {
		t.Fatalf("fehlermeldung fehlt: %s", html)
	}
}

func TestFieldOhneForm(t *testing.T) {
	html := render(t, Field(testCtx(t), nil, "name", "Name", "text"))
	if !strings.Contains(html, `value=""`) {
		t.Fatalf("ohne form soll value leer sein: %s", html)
	}
	if strings.Contains(html, `value="`) && !strings.Contains(html, `value=""`) {
		t.Fatalf("kein fremder value erwartet: %s", html)
	}
}

func TestScopingVorhanden(t *testing.T) {
	html := render(t, Button(testCtx(t), "/x", "X", Primary))
	if !strings.Contains(html, "data-scope=") {
		t.Fatalf("kein scoping: %s", html)
	}
}
