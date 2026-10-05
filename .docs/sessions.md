# Session, Flash, CSRF

## Session-Store anlegen

```go
store, err_store := dreego.NewCookieStore([]byte(secret))  // secret >= 32 Bytes
if err_store != nil {
	log.Fatal(err_store)
}
app.SetSessionStore(store)
```

Ohne Store sind Session, Flash und CSRF inert. Das Secret **muss** gesetzt und
geheim sein — ein fest eingebautes Standardgeheimnis wäre ein öffentliches
Passwort.

Eigenschaften:
- **Verschlüsselt** (AES-256-GCM): Inhalt und Integrität geschützt.
- **Ein Write pro Request**: die Session wird einmal geladen, im Speicher
  geändert und als **ein** Cookie geschrieben. So überschreiben sich mehrere
  `Set-Cookie`-Header nicht (ein Fehler, der sonst die Anmeldung still killt).
- **HttpOnly**, **SameSite=Lax**, **Secure** bei TLS/trusted Proxy.
- Größenlimit 4096 Bytes (Cookie-Grenze). Zu groß → Serverfehler, nicht stiller
  Verlust.
- **Lebensdauer** über `CookiePolicy.MaxAge` (Default `0` = Session-Cookie).
  Positive Werte setzen **Max-Age und Expires**. Dauer-Helfer:
  `dreego.Days(30)`, `dreego.Hours(12)`, `dreego.Minutes(90)` (Go kennt keine
  `time.Day`-Konstante).

```go
store.SetCookiePolicy(dreego.CookiePolicy{
    SameSite: http.SameSiteLaxMode,
    MaxAge:   dreego.Days(30),   // 30 Tage „eingeloggt bleiben"
})
```

## Werte lesen/schreiben

```go
c.SetSessionVal("konto_email", in.Email)
email := c.SessionVal("konto_email")
c.DelSessionVal("konto_email")
c.DestroySession()   // logout: alles löschen
```

## Flash (Einmal-Nachrichten)

```go
// nach einem POST:
c.Flash("login_error", "Bitte Eingaben prüfen.")
return c.Redirect("/login", 303)

// im GET:
meldung := c.FlashGet("login_error")   // liest UND löscht
```

Genau einmal sichtbar — ideal für Post/Redirect/Get.

## CSRF

Standardmäßig **an**. Unsichere Methoden (POST/PUT/PATCH/DELETE) müssen das
Token tragen, sonst `403`.

Im Formular:

```go
d.Form(
	d.Attr("method", "post"),
	c.CSRFInput(),     // <- das verborgene Token-Feld
	…
)
```

- Token lebt in der (HttpOnly-)Session; ein GET legt es an.
- Prüfung mit konstantzeit-Vergleich (kein Timing-Leak).
- Alternative: Header `X-CSRF-Token` (für Fetch/JSON-POSTs).
- Ohne Store kann CSRF nicht funktionieren → **500** statt stiller Durchlass.
- Abschaltbar für reine Marketing-Apps ohne Formulare:
  `app.SetCSRF(false)`.

## Cookie-Richtlinie

```go
store.SetCookiePolicy(dreego.CookiePolicy{
	SameSite: http.SameSiteStrictMode,
	Secure:   true,   // z. B. hinter TLS-Proxy
})
```

`HttpOnly` und `Path=/` können nicht abgeschaltet werden. Hinter einem
TLS-terminierenden Proxy:

```go
store.SetTrustedProxies([]string{"10.0.0.1"})   // nutzt X-Forwarded-Proto
```

Nur ein **vertrauenswürdiger** Proxy darf `Secure` über den Header setzen — ein
beliebiger Client kann es nicht erzwingen.
