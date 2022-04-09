package gateway

import (
	"micrach/config"

	"github.com/gofiber/fiber/v2"
)

func Ping(c *fiber.Ctx) error {
	headerKey := c.GetReqHeaders()["Authorization"]
	if config.App.Gateway.ApiKey != headerKey {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}
	return c.JSON(fiber.Map{
		"message": "pong",
	})
}
