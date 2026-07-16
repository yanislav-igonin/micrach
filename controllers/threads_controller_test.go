package controllers

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"

	"micrach/config"
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

func TestCreateThreadRejectsInvalidMultipartBeforeDatabase(t *testing.T) {
	previousConfig := config.App
	t.Cleanup(func() {
		config.App = previousConfig
	})
	config.App.IsCaptchaActive = false

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	if err := writer.WriteField("title", ""); err != nil {
		t.Fatal(err)
	}
	if err := writer.WriteField("text", ""); err != nil {
		t.Fatal(err)
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}

	app := fiber.New(fiber.Config{Views: testViews{}})
	app.Post("/", CreateThread)

	req := httptest.NewRequest(http.MethodPost, "/", &body)
	req.Header.Set(fiber.HeaderContentType, writer.FormDataContentType())
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("status = %d, want %d", resp.StatusCode, fiber.StatusBadRequest)
	}
}
