# Formulare: Binding und Validierung

## Eingabe als Struktur (das `in`)

```go
type LoginForm struct {
	Email string `form:"email" validate:"required,email" label:"E-Mail"`
	Team  string `form:"team"  validate:"required,min=2" label:"Teamname"`
}
```

Tags:
- `form:"name"` — Feldname im Formular (Default: kleingeschriebener Feldname).
- `validate:"…"` — Regeln (siehe unten).
- `label:"…"` — Anzeigename in Fehlermeldungen.

## Binden und prüfen

```go
func postLogin(c *dreego.Ctx) d.View {
	in, form, err_bind := dreego.Bind[LoginForm](c)
	if err_bind != nil {
		// kaputte Anfrage oder falsch getaggte Struktur — kein Validierungsfehler
		return c.Redirect("/login", 303)
	}
	if form.HasErrors() {
		// ungültig: Formular erneut zeigen (PRG), Meldung per Flash
		c.Flash("login_error", "Bitte Eingaben prüfen.")
		return c.Redirect("/login", 303)
	}
	// gültig: `in` benutzen
	return c.Redirect("/dashboard", 303)
}
```

`Bind[T]` gibt drei Dinge zurück:
1. **T** — die gefüllte Struktur (`in`).
2. **\*Form** — Werte und Fehler.
3. **error** — nur für echte Fehler (kaputte Anfrage, falsche Tags). Eine
   fehlgeschlagene Validierung ist **kein** Go-Fehler, sondern landet im `Form`.

## Das Formular rendern

```go
func loginForm(c *dreego.Ctx, form *dreego.Form) d.View {
	return d.Form(
		d.Attr("method", "post"), d.Attr("action", "/login"),
		c.CSRFInput(),                       // CSRF-Token (Pflicht bei POST)
		ui.Field(c, form, "email", "E-Mail", "email"),
		ui.Field(c, form, "team",  "Team", "text"),
		d.Button(d.Attr("type", "submit"), d.Text("Absenden")),
	)
}
```

- `form.Old("email")` — der zuletzt eingegebene Wert.
- `form.Errors("email")` — Liste der Meldungen.
- `form.SetError("email", "…")` — eigene (fachliche) Meldung anhängen.

## Validierungsregeln

| Regel | Wirkung |
|---|---|
| `required` | darf nicht leer sein |
| `email` | muss wie eine E-Mail aussehen |
| `min=N` | mindestens N Zeichen |
| `max=N` | höchstens N Zeichen |
| `len=N` | genau N Zeichen |

Mehrere Regeln durch Komma: `validate:"required,email"`.

## Meldungen und Labels

Die Meldung nutzt das `label`-Tag: bei `label:"E-Mail"` lautet sie
„E-Mail ist keine gültige E-Mail-Adresse." Ohne `label` wird der Feldname
verwendet.

## Post/Redirect/Get (PRG)

Nach einem POST nie direkt rendern, sondern umleiten (`303`). Die Meldung reist
per Flash — siehe [sessions.md](sessions.md). So verhindert man Doppel-Absenden
beim Neuladen.

```go
c.Flash("login_error", "Bitte prüfen.")
return c.Redirect("/login", 303)
```
