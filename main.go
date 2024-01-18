package main

import (
	"log"
	"os"

	"github.com/slavik-o/go-web-app/controllers"
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
	app.Get("/", middlewares.AuthRequired, controllers.GetIndex)

	app.Get("/login", controllers.GetLogin)
	app.Get("/logout", controllers.GetLogout)
	app.Get("/auth/:provider", controllers.GetAuthProvider)
	app.Get("/auth/:provider/callback", controllers.GetAuthProviderCallback)

	// Run
	log.Fatal(app.Listen(os.Getenv("BIND_URL")))
}
