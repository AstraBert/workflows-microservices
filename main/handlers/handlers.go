package handlers

import (
	"ecommerce/templates"

	"github.com/gofiber/fiber/v2"
)

func HomeRoute(c *fiber.Ctx) error {
	home := templates.Home()
	c.Set("Content-Type", "text/html")
	return home.Render(c.Context(), c.Response().BodyWriter())
}

func RegisterUser(c *fiber.Ctx) error {
	return nil
}

func LoginUser(c *fiber.Ctx) error {
	return nil
}

func LogoutUser(c *fiber.Ctx) error {
	return nil
}
