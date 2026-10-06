// Command wails-dreego is the Wails v3 desktop template for the new dreego.
//
// The frontend is rendered by dreego: dreego's *App is itself an http.Handler,
// so Wails' asset server calls it in-process — no adapter, no bundler and no TCP
// listener. The window's markup lives in the frontend package (dreego pages) and
// the client is plain JavaScript emitted inline; the generated bindings are
// served from frontend/static/bindings. The GreetService and the "time" event are
// the same as the vanilla template.
package main

import (
	"log"
	"time"

	"wails-dreego/frontend"

	"github.com/wailsapp/wails/v3/pkg/application"
)

func init() {
	// Register a typed custom event; the binding generator exposes it to the
	// frontend. The dreego page listens for "time" and shows it in the footer.
	application.RegisterEvent[string]("time")
}

func main() {
	// The dreego app owns the pages. No listener is started: Wails serves every
	// request from this handler in-process.
	dreegoApp := frontend.NewApp()

	wailsApp := application.New(application.Options{
		Name:        "Wails Dreego",
		Description: "A dreego frontend rendered in a native Wails v3 window",
		Services: []application.Service{
			application.NewService(&GreetService{}),
		},
		Assets: application.AssetOptions{
			Handler:        dreegoApp.Handler(),
			DisableLogging: true,
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	wailsApp.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:  "Wails Dreego",
		Width:  1000,
		Height: 618,
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInset,
		},
		BackgroundColour: application.NewRGB(6, 7, 15),
		URL:              "/",
	})

	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				wailsApp.Event.Emit("time", time.Now().Format(time.RFC1123))
			case <-wailsApp.Context().Done():
				return
			}
		}
	}()

	if err := wailsApp.Run(); err != nil {
		log.Fatal(err)
	}
}
