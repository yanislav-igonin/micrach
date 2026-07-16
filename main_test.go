package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

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
	app.Get("/static/fallback", func(c *fiber.Ctx) error {
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
		Next: func(c *fiber.Ctx) bool {
			return skipRateLimit(c.IsFromLocal(), c.Path())
		},
		Max: 2,
	}))
	app.Get("/", func(c *fiber.Ctx) error {
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
