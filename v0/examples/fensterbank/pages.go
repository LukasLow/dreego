package main

import (
	"strings"

	d "github.com/LukasLow/dreego/v0/pkg/dom"

	"github.com/LukasLow/dreego/v0/pkg/dreego"
	"github.com/LukasLow/dreego/v0/pkg/scope"
)

// RegisterPages hängt alle HTML-Seiten an die App. Jede Seite hat einen eigenen
// Pfad — sie geraten sich nicht in die Quere.
func RegisterPages(app *dreego.App) {
	app.Page(Start)
	app.Page(Sorten)
	app.Page(Sorte)
	app.Page(Pflege)
	app.Page(Kontakt)
	app.Page(Ueber)
	app.Page(Shop)
}

// Start — die Startseite.
var Start = dreego.Page{
	Path:   "/",
	Layout: Shell,
	Nav:    "start",
	Head:   func(c *dreego.Ctx) []d.View { return []d.View{d.TitleEl(d.Text("Fensterbank — Start"))} },
	Get:    getStart,
}

func getStart(c *dreego.Ctx) d.View {
	sorten := []string{"Efeutute", "Bogenhanf", "Zamioculcas", "Monstera"}

	return d.Group([]d.View{
		c.Box(
			scope.CSS(`
.hero { text-align: center; padding: 44px 16px 30px }
.hero h1 { font-size: 40px; margin: 0 0 10px }
.hero h1 em { color: #0b6b4a; font-style: normal }
.hero p { font-size: 19px; color: #3d5c4e; max-width: 560px; margin: 0 auto }
`),
			d.Section(d.Class("hero"),
				d.H1(c.Text("start.titel")),
				d.P(c.Text("start.lead")),
			),
		),
		d.Section(
			d.H2(d.Text(c.T("start.beliebt"))),
			d.Ul(d.Map(sorten, func(name string) d.View {
				// Jede Sorte verlinkt auf ihre dynamische [sorte]-Seite.
				return d.Li(d.A(d.Href("/sorten/"+slug(name)), d.Text(name)))
			})),
		),
	})
}

// Sorten — Übersicht, verlinkt auf die dynamischen Detailseiten.
var Sorten = dreego.Page{
	Path:   "/sorten",
	Layout: Shell,
	Nav:    "sorten",
	Head:   func(c *dreego.Ctx) []d.View { return []d.View{d.TitleEl(d.Text("Sorten — Fensterbank"))} },
	Get:    getSorten,
}

func getSorten(c *dreego.Ctx) d.View {
	return d.Group([]d.View{
		d.H1(d.Text(c.T("sorten.titel"))),
		d.P(d.Text(c.T("sorten.lead"))),
		d.Ul(d.Map(sortenListe, func(s sorte) d.View {
			return d.Li(d.A(d.Href("/sorten/"+s.Slug), d.Text(s.Name)))
		})),
	})
}

// Sorte — DYNAMISCHE Route /sorten/[id]. Der Wert kommt aus c.Param("id").
var Sorte = dreego.Page{
	Path:   "/sorten/[id]",
	Layout: Shell,
	Nav:    "sorten",
	Head:   func(c *dreego.Ctx) []d.View { return []d.View{d.TitleEl(d.Text("Sorte — Fensterbank"))} },
	Get:    getSorte,
}

func getSorte(c *dreego.Ctx) d.View {
	id := c.Param("id")
	s, gefunden := findeSorte(id)
	if !gefunden {
		return d.Group([]d.View{
			d.H1(d.Text(c.T("sorte.nichtgefunden"))),
			d.P(d.Text("Keine Sorte mit der Kennung "), d.Code(d.Text(id))),
			d.P(d.A(d.Href("/sorten"), d.Text(c.T("sorte.zurueck")))),
		})
	}
	return d.Group([]d.View{
		d.H1(d.Text(s.Name)),
		d.P(d.Em(d.Text(s.Latein))),
		d.P(d.Text(s.Text)),
		d.P(d.A(d.Href("/pflege"), d.Text(c.T("sorte.pflege")))),
	})
}

// Pflege — feste Seite.
var Pflege = dreego.Page{
	Path:   "/pflege",
	Layout: Shell,
	Nav:    "pflege",
	Head:   func(c *dreego.Ctx) []d.View { return []d.View{d.TitleEl(d.Text("Pflege — Fensterbank"))} },
	Get: func(c *dreego.Ctx) d.View {
		regeln := []string{
			"Lieber zu wenig als zu viel giessen.",
			"Alle zwei Wochen düngen reicht.",
			"Im Winter weniger Wasser.",
		}
		return d.Group([]d.View{
			d.H1(d.Text("Pflege")),
			d.Ul(d.Map(regeln, func(r string) d.View { return d.Li(d.Text(r)) })),
		})
	},
}

// Kontakt — feste Seite.
var Kontakt = dreego.Page{
	Path:   "/kontakt",
	Layout: Shell,
	Nav:    "kontakt",
	Head:   func(c *dreego.Ctx) []d.View { return []d.View{d.TitleEl(d.Text("Kontakt — Fensterbank"))} },
	Get: func(c *dreego.Ctx) d.View {
		return d.Group([]d.View{
			d.H1(d.Text(c.T("kontakt.titel"))),
			d.P(d.Text(c.T("kontakt.lead")), d.Text(" "),
				d.A(d.Href("mailto:hallo@fensterbank.example"), d.Text("hallo@fensterbank.example"))),
		})
	},
}

// Ueber — feste Seite.
var Ueber = dreego.Page{
	Path:   "/ueber",
	Layout: Shell,
	Head:   func(c *dreego.Ctx) []d.View { return []d.View{d.TitleEl(d.Text("Über — Fensterbank"))} },
	Get: func(c *dreego.Ctx) d.View {
		return d.Group([]d.View{
			d.H1(d.Text("Über")),
			d.P(d.Text("Eine Demo mit dreego: viele Seiten, dynamische Parameter, API.")),
		})
	},
}

type sorte struct {
	Slug   string `json:"slug"`
	Name   string `json:"name"`
	Latein string `json:"latein"`
	Text   string `json:"text"`
}

var sortenListe = []sorte{
	{"efeutute", "Efeutute", "Epipremnum aureum", "Wächst auch bei wenig Licht."},
	{"bogenhanf", "Bogenhanf", "Sansevieria", "Extrem genügsam, kaum giessen."},
	{"zamioculcas", "Zamioculcas", "Zamioculcas zamiifolia", "Verträgt Trockenheit lange."},
	{"monstera", "Monstera", "Monstera deliciosa", "Braucht etwas mehr Licht."},
}

func findeSorte(id string) (sorte, bool) {
	for _, s := range sortenListe {
		if s.Slug == id {
			return s, true
		}
	}
	return sorte{}, false
}

func slug(name string) string {
	return strings.ToLower(name)
}
