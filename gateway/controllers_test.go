package gateway

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"

	"micrach/config"
)

func TestPingUsesLastAuthorizationHeader(t *testing.T) {
	previousConfig := config.App
	t.Cleanup(func() {
		config.App = previousConfig
	})
	config.App.Gateway.ApiKey = "configured-key"

	app := fiber.New()
	app.Get("/api/ping", Ping)

	req := httptest.NewRequest("GET", "/api/ping", nil)
	req.Header.Add("Authorization", "wrong-key")
	req.Header.Add("Authorization", "configured-key")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("test ping request: %v", err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status %d when final Authorization header matches, got %d", fiber.StatusOK, resp.StatusCode)
	}
}
