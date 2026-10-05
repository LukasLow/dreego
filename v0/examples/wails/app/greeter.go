package app

import (
	"context"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// GreeterService is an application-owned Wails service. dreego never wraps or
// validates it; Wails generates the typed Go-to-client bindings from its
// exported methods. State belongs here, guarded explicitly.
type GreeterService struct{}

// NewGreeterService builds the service.
func NewGreeterService() *GreeterService {
	return &GreeterService{}
}

// ServiceStartup satisfies the Wails service lifecycle.
func (g *GreeterService) ServiceStartup(_ context.Context, _ application.ServiceOptions) error {
	return nil
}

// ServiceShutdown satisfies the Wails service lifecycle.
func (g *GreeterService) ServiceShutdown() error {
	return nil
}

// Greet returns a greeting for the given name. Empty input falls back to
// "world". The boundary value is validated here rather than trusted.
func (g *GreeterService) Greet(name string) string {
	if name == "" {
		name = "world"
	}
	return "Hello, " + name + "!"
}
