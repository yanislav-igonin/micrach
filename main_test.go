package main

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"os"
	"path/filepath"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/limiter"

	"micrach/config"
	"micrach/controllers"
	"micrach/utils"
)

type mainTestViews struct{}

func (mainTestViews) Load() error { return nil }

func (mainTestViews) Render(out io.Writer, name string, _ any, _ ...string) error {
	_, err := io.WriteString(out, name)
	return err
}

func TestRegisterStaticRoutes(t *testing.T) {
	t.Chdir(t.TempDir())

	files := map[string]string{
		"static/styles/app.css": "body { color: black; }",
		"uploads/42/o/7.png":    "png-body",
	}
	for name, content := range files {
		if err := os.MkdirAll(filepath.Dir(name), 0o755); err != nil {
			t.Fatalf("create fixture directory: %v", err)
		}
		if err := os.WriteFile(name, []byte(content), 0o644); err != nil {
			t.Fatalf("write fixture: %v", err)
		}
	}

	app := fiber.New()
	registerStaticRoutes(app)
	app.Get("/static/fallback", func(c fiber.Ctx) error {
		return c.SendString("fallback")
	})

	tests := []struct {
		name     string
		path     string
		wantCode int
		wantBody string
	}{
		{name: "checked-in static file", path: "/static/styles/app.css", wantCode: fiber.StatusOK, wantBody: files["static/styles/app.css"]},
		{name: "generated upload", path: "/uploads/42/o/7.png", wantCode: fiber.StatusOK, wantBody: files["uploads/42/o/7.png"]},
		{name: "missing file continues routing", path: "/static/fallback", wantCode: fiber.StatusOK, wantBody: "fallback"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := app.Test(httptest.NewRequest(http.MethodGet, tt.path, nil))
			if err != nil {
				t.Fatalf("request %s: %v", tt.path, err)
			}
			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("read response: %v", err)
			}
			if resp.StatusCode != tt.wantCode || string(body) != tt.wantBody {
				t.Fatalf("response = (%d, %q), want (%d, %q)", resp.StatusCode, body, tt.wantCode, tt.wantBody)
			}
		})
	}
}

func TestSkipRateLimit(t *testing.T) {
	tests := []struct {
		name    string
		isLocal bool
		path    string
		want    bool
	}{
		{name: "ordinary remote request", path: "/", want: false},
		{name: "local request", isLocal: true, path: "/", want: true},
		{name: "static", path: "/static/styles/index.css", want: true},
		{name: "upload", path: "/uploads/1/o/1.png", want: true},
		{name: "captcha", path: "/captcha/id", want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := skipRateLimit(tt.isLocal, tt.path); got != tt.want {
				t.Fatalf("skipRateLimit(%v, %q) = %v, want %v", tt.isLocal, tt.path, got, tt.want)
			}
		})
	}
}

func TestLimiterRejectsOrdinaryRemoteRequests(t *testing.T) {
	app := fiber.New()
	app.Use(limiter.New(limiter.Config{
		Next: func(c fiber.Ctx) bool {
			return skipRateLimit(c.IsFromLocal(), c.Path())
		},
		Max: 2,
	}))
	app.Get("/", func(c fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	for requestNumber, want := range []int{fiber.StatusOK, fiber.StatusOK, fiber.StatusTooManyRequests} {
		resp, err := app.Test(httptest.NewRequest(http.MethodGet, "/", nil))
		if err != nil {
			t.Fatalf("request %d: %v", requestNumber+1, err)
		}
		resp.Body.Close()
		if resp.StatusCode != want {
			t.Fatalf("request %d status = %d, want %d", requestNumber+1, resp.StatusCode, want)
		}
	}
}

func TestMultipartBodyLimitAllowsApplicationFileSizeValidation(t *testing.T) {
	previousConfig := config.App
	t.Cleanup(func() {
		config.App = previousConfig
	})
	config.App.IsCaptchaActive = false

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	if err := writer.WriteField("title", "oversize"); err != nil {
		t.Fatal(err)
	}
	if err := writer.WriteField("text", "oversize upload"); err != nil {
		t.Fatal(err)
	}
	header := make(textproto.MIMEHeader)
	header.Set("Content-Disposition", `form-data; name="files"; filename="oversize.png"`)
	header.Set("Content-Type", "image/png")
	part, err := writer.CreatePart(header)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := part.Write(bytes.Repeat([]byte("x"), utils.FILE_SIZE_IN_BYTES+(1<<20))); err != nil {
		t.Fatal(err)
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}

	app := fiber.New(fiber.Config{Views: mainTestViews{}, BodyLimit: multipartBodyLimit})
	app.Post("/", controllers.CreateThread)
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
