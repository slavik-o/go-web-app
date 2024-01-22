package main

import (
	"log"
	"os"

	"github.com/slavik-o/go-web-app/handlers"
	"github.com/slavik-o/go-web-app/middlewares"
)

// Main
func main() {
	// Env
	LoadEnv()

	// OAuth
	SetupOAuth()

	// App
	app := NewApp()

	// Routes
	app.Get("/", middlewares.AuthRequired, handlers.GetIndex)

	app.Get("/login", handlers.GetLogin)
	app.Get("/logout", handlers.GetLogout)
	app.Get("/auth/:provider", handlers.GetAuthProvider)
	app.Get("/auth/:provider/callback", handlers.GetAuthProviderCallback)

	// Run
	log.Fatal(app.Listen(os.Getenv("BIND_URL")))
}
