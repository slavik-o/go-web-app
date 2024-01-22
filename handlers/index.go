package handlers

import (
	"github.com/gofiber/fiber/v2"

	"github.com/slavik-o/go-web-app/models"
	"github.com/slavik-o/go-web-app/views"
)

func GetIndex(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)

	return Render(c, views.Index(user.FirstName))
}
