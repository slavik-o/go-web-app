package main

import (
	"log"
	"os"
	"time"

	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/encryptcookie"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/google"
	"github.com/shareed2k/goth_fiber"

	"github.com/slavik-o/go-web-app/views"
)

func Render(ctx *fiber.Ctx, cmp templ.Component) error {
	ctx.Set("Content-Type", "text/html")
	return cmp.Render(ctx.Context(), ctx.Response().BodyWriter())
}

func AuthRequired(c *fiber.Ctx) error {
	if c.Cookies("user") == "" {
		return c.Redirect("/login")
	}

	return c.Next()
}

func main() {
	// Env
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// OAuth
	goth.UseProviders(
		google.New(
			os.Getenv("GOOGLE_CLIENT_ID"),
			os.Getenv("GOOGLE_CLIENT_SECRET"),
			os.Getenv("GOOGLE_CALLBACK_URL"),
			"profile", "email", "openid",
		),
	)

	// App
	app := fiber.New()

	// Assets
	app.Static("/", "./assets")

	// Middleware
	app.Use(logger.New())

	app.Use(encryptcookie.New(encryptcookie.Config{
		Key: os.Getenv("COOKIE_SECRET"),
	}))

	// Routes
	app.Get("/", AuthRequired, func(c *fiber.Ctx) error {
		return Render(c, views.Index())
	})

	app.Get("/login", func(c *fiber.Ctx) error {
		return Render(c, views.Login())
	})

	app.Get("/auth/:provider", func(c *fiber.Ctx) error {
		return goth_fiber.BeginAuthHandler(c)
	})

	app.Get("/auth/:provider/callback", func(c *fiber.Ctx) error {
		user, err := goth_fiber.CompleteUserAuth(c)
		if err != nil {
			return c.Redirect("/login")
		}

		c.Cookie(&fiber.Cookie{
			Name:    "user",
			Value:   user.UserID,
			Expires: time.Now().Add(7 * 24 * time.Hour),
		})

		return c.Redirect("/")
	})

	app.Get("/logout", func(c *fiber.Ctx) error {
		c.ClearCookie("user")

		return c.Redirect("/login")
	})

	// Run
	log.Fatal(app.Listen(os.Getenv("BIND_URL")))
}
