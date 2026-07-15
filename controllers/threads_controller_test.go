package controllers

import (
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

type testViews struct{}

func (testViews) Load() error { return nil }

func (testViews) Render(out io.Writer, name string, _ interface{}, _ ...string) error {
	_, err := io.WriteString(out, name)
	return err
}

func TestGetThreadsRejectsInvalidPage(t *testing.T) {
	app := fiber.New(fiber.Config{Views: testViews{}})
	app.Get("/", GetThreads)

	resp, err := app.Test(httptest.NewRequest("GET", "/?page=0", nil))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != fiber.StatusNotFound {
		t.Fatalf("status = %d, want %d", resp.StatusCode, fiber.StatusNotFound)
	}
}

func TestGetThreadRejectsInvalidID(t *testing.T) {
	app := fiber.New(fiber.Config{Views: testViews{}})
	app.Get("/:threadID", GetThread)

	resp, err := app.Test(httptest.NewRequest("GET", "/invalid", nil))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != fiber.StatusNotFound {
		t.Fatalf("status = %d, want %d", resp.StatusCode, fiber.StatusNotFound)
	}
}
