package controllers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/shareed2k/goth_fiber"

	"github.com/slavik-o/go-web-app/database"
	"github.com/slavik-o/go-web-app/views"
)

func GetLogin(c *fiber.Ctx) error {
	return Render(c, views.Login())
}

func GetAuthProvider(c *fiber.Ctx) error {
	return goth_fiber.BeginAuthHandler(c)
}

func GetAuthProviderCallback(c *fiber.Ctx) error {
	authUser, err := goth_fiber.CompleteUserAuth(c)
	if err != nil {
		return c.Redirect("/login")
	}

	db := database.Connect()

	user, err := db.UpsertAuthUser(authUser)
	if err != nil {
		return c.Redirect("/login")
	}

	c.Cookie(&fiber.Cookie{
		Name:    "user",
		Value:   user.ID,
		Expires: time.Now().Add(7 * 24 * time.Hour),
	})

	return c.Redirect("/")
}

func GetLogout(c *fiber.Ctx) error {
	c.ClearCookie("user")

	return c.Redirect("/login")
}
