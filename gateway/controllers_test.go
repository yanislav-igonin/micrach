package gateway

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestPingIsPublic(t *testing.T) {
	app := fiber.New()
	app.Get("/api/ping", Ping)

	req := httptest.NewRequest("GET", "/api/ping", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("test ping request: %v", err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status %d, got %d", fiber.StatusOK, resp.StatusCode)
	}

	var body map[string]bool
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if !body["ok"] {
		t.Fatalf("expected ok=true, got %+v", body)
	}
}
