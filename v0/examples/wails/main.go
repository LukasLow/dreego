package main

import (
	"log"

	"wails-dreego/app"
	webapp "wails-dreego/app/web"

	"github.com/LukasLow/dreego/v0/pkg/dreego"
	"github.com/LukasLow/dreego/v0/plg/wails"
	"github.com/wailsapp/wails/v3/pkg/application"
)

// main wires dreego (render + assets) into Wails. dreego owns no window, no
// service, no lifecycle and starts no HTTP listener; the adapter only serves
// rendered pages and assets in-process. The application owns every Wails
// decision explicitly.
func main() {
	dreegoApp := dreego.NewApp()
	webapp.Register(dreegoApp)

	wailsApp := application.New(application.Options{
		Name:        "Wails Dreego",
		Description: "A dreego frontend rendered in a native Wails v3 window",
		Assets: application.AssetOptions{
			Handler: wails.Handler(dreegoApp),
		},
		Services: []application.Service{
			application.NewService(app.NewGreeterService()),
		},
	})

	wailsApp.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:     "Wails Dreego",
		URL:       "/",
		Width:     720,
		Height:    560,
		MinWidth:  360,
		MinHeight: 420,
	})

	err_run := wailsApp.Run()
	if err_run != nil {
		log.Fatal(err_run)
	}
}
