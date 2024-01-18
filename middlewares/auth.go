package middlewares

import (
	"github.com/gofiber/fiber/v2"

	"github.com/slavik-o/go-web-app/database"
)

func AuthRequired(c *fiber.Ctx) error {
	if c.Cookies("user") == "" {
		return c.Redirect("/login")
	}

	db := database.Connect()

	user, err := db.FindUserByID(c.Cookies("user"))
	if err != nil {
		return c.Redirect("/login")
	}

	c.Locals("user", user)

	return c.Next()
}
