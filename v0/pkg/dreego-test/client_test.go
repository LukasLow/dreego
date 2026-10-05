package dreegotest

import (
	"net/http"
	"testing"
	"testing/fstest"

	d "github.com/LukasLow/dreego/v0/pkg/dom"
	"github.com/LukasLow/dreego/v0/pkg/dreego"
)

func staticFS() fstest.MapFS {
	return fstest.MapFS{
		"public/hallo.txt": &fstest.MapFile{Data: []byte("statisch")},
	}
}

func testApp() *dreego.App {
	app := dreego.NewApp()
	app.SetCSRF(false)
	app.Page(dreego.Page{
		Path: "/",
		Get:  func(c *dreego.Ctx) d.Node { return d.H1(d.Text("Hallo Welt")) },
	})
	app.Page(dreego.Page{
		Path: "/weg",
		Get:  func(c *dreego.Ctx) d.Node { return c.Redirect("/", 303) },
	})
	app.Page(dreego.Page{
		Path: "/echo",
		API:  func(c *dreego.Ctx) error { return c.JSON(200, map[string]string{"ok": "ja"}) },
	})
	app.Static(staticFS(), "public")
	return app
}

func TestGetUndAssert(t *testing.T) {
	client := New(testApp().Handler())
	resp := client.Get("/")
	resp.AssertStatus(t, 200)
	resp.AssertContains(t, "Hallo Welt")
}

func TestRedirectWirdNichtVerfolgt(t *testing.T) {
	client := New(testApp().Handler())
	resp := client.Get("/weg")
	resp.AssertStatus(t, 303)
	resp.AssertLocation(t, "/")
}

func TestAPIJSON(t *testing.T) {
	client := New(testApp().Handler())
	resp := client.Get("/echo")
	resp.AssertStatus(t, 200)
	resp.AssertHeader(t, "Content-Type", "application/json; charset=utf-8")
	resp.AssertContains(t, `"ok":"ja"`)
}

func TestPostForm(t *testing.T) {
	app := dreego.NewApp()
	app.SetCSRF(false)
	app.Page(dreego.Page{
		Path: "/form",
		Post: func(c *dreego.Ctx) d.Node { return d.Text("empfangen:" + c.FormValue("name")) },
	})
	client := New(app.Handler())
	resp := client.PostForm("/form", map[string]string{"name": "Lukas"})
	resp.AssertStatus(t, 200)
	resp.AssertContains(t, "empfangen:Lukas")
}

func TestNichtGefunden(t *testing.T) {
	client := New(testApp().Handler())
	resp := client.Get("/gibtsnicht")
	resp.AssertStatus(t, http.StatusNotFound)
}

func TestStatic(t *testing.T) {
	client := New(testApp().Handler())
	resp := client.Get("/public/hallo.txt")
	resp.AssertStatus(t, 200)
	resp.AssertContains(t, "statisch")
}
