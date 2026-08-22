package main

import (
	"context"

	"SmokeLab/packages/engine"
)

// App struct
type App struct {
	ctx             context.Context
	greetingService *engine.GreetingService
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		greetingService: engine.NewGreetingService(),
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return a.greetingService.Greet(name)
}
