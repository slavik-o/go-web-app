package handlers

import (
	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
)

func Render(ctx *fiber.Ctx, cmp templ.Component) error {
	ctx.Set("Content-Type", "text/html")
	return cmp.Render(ctx.Context(), ctx.Response().BodyWriter())
}
