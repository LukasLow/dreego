package main

import (
	d "github.com/LukasLow/dreego/v0/pkg/dom"

	"github.com/LukasLow/dreego/v0/pkg/dreego"
	"github.com/LukasLow/dreego/v0/pkg/scope"
	"github.com/LukasLow/dreego/v0/pkg/ui"
)

// Shop — eine Seite, die die Komponenten-Bibliothek dreego/ui verwendet:
// Cards, Badges und ein Button. Kein handgeschriebenes CSS/HTML pro Baustein.
var Shop = dreego.Page{
	Path:   "/shop",
	Layout: Shell,
	Nav:    "shop",
	Head:   func(c *dreego.Ctx) []d.View { return []d.View{d.TitleEl(d.Text("Shop — Fensterbank"))} },
	Get:    getShop,
}

func getShop(c *dreego.Ctx) d.View {
	return d.Group([]d.View{
		c.Box(
			scope.CSS(`
.shop { padding: 34px 0 }
.shop h1 { font-size: 34px; margin: 0 0 8px }
.shop .lead { color: #3d5c4e; margin: 0 0 24px }
.shop .raster { display: grid; grid-template-columns: 1fr 1fr; gap: 16px }
.shop .aktionen { margin-top: 24px; display: flex; gap: 12px; flex-wrap: wrap }
@media (max-width: 640px) { .shop h1 { font-size: 27px } .shop .raster { grid-template-columns: 1fr } }
`),
			d.Section(d.Class("shop"),
				d.H1(d.Text("Shop")),
				d.P(d.Class("lead"), d.Text("Beispielseite mit der Komponenten-Bibliothek dreego/ui.")),
				d.Div(d.Class("raster"),
					ui.Card(c, "Efeutute", d.Group([]d.View{
						d.P(d.Text("Wächst auch bei wenig Licht.")),
						ui.Badge(c, "Pflegeleicht", ui.Success),
					})),
					ui.Card(c, "Monstera", d.Group([]d.View{
						d.P(d.Text("Braucht etwas mehr Licht.")),
						ui.Badge(c, "Beliebt", ui.Info),
					})),
					ui.Card(c, "Bogenhanf", d.Group([]d.View{
						d.P(d.Text("Extrem genügsam.")),
						ui.Badge(c, "Wenig Wasser", ui.Info),
					})),
					ui.Card(c, "Zamioculcas", d.Group([]d.View{
						d.P(d.Text("Verträgt Trockenheit lange.")),
						ui.Badge(c, "Anfänger", ui.Success),
					})),
				),
				d.Div(d.Class("aktionen"),
					ui.Button(c, "/kontakt", "Anfragen", ui.Primary),
					ui.Button(c, "/sorten", "Alle Sorten", ui.Ghost),
				),
			),
		),
	})
}
