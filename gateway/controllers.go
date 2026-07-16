package gateway

import (
	"micrach/config"

	"github.com/gofiber/fiber/v3"
)

func Ping(c fiber.Ctx) error {
	headerValues := c.GetReqHeaders()["Authorization"]
	headerKey := ""
	if len(headerValues) > 0 {
		headerKey = headerValues[len(headerValues)-1]
	}
	if config.App.Gateway.ApiKey != headerKey {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}
	return c.JSON(fiber.Map{
		"message": "pong",
	})
}
