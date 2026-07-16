package templates

import (
	"bytes"
	"strings"
	"testing"

	"github.com/gofiber/template/html/v2"

	"micrach/repositories"
)

func TestHTMLTemplatesLoadAndRender(t *testing.T) {
	engine := html.New("../templates", ".html")
	engine.AddFunc("Iterate", Iterate)
	engine.AddFunc("NotNil", NotNil)

	if err := engine.Load(); err != nil {
		t.Fatalf("load templates: %v", err)
	}

	const message = "template compatibility"
	var rendered bytes.Buffer
	if err := engine.Render(&rendered, "pages/400", repositories.BadRequestHtmlData{Message: message}); err != nil {
		t.Fatalf("render template: %v", err)
	}
	if !strings.Contains(rendered.String(), message) {
		t.Fatalf("rendered template does not contain %q", message)
	}
}
