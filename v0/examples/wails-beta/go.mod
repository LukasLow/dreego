module wails-dreego

go 1.25.0

require (
	github.com/LukasLow/dreego v0.0.2
	github.com/wailsapp/wails/v3 v3.0.0-beta.27
)

require (
	github.com/adrg/xdg v0.5.3 // indirect
	github.com/coder/websocket v1.8.14 // indirect
	github.com/go-ole/go-ole v1.3.0 // indirect
	github.com/godbus/dbus/v5 v5.2.2 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	golang.org/x/sys v0.46.0 // indirect
)

// This example lives inside the dreego repository, so it builds against the
// checkout it sits in (three levels up: v0/examples/wails-beta -> repo root).
replace github.com/LukasLow/dreego => ../../..
