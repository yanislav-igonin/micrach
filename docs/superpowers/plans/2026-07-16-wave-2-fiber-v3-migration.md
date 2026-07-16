# Wave 2 Fiber v3 Migration Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Migrate micrach from Fiber v2 to the current stable Fiber v3 and compatible HTML template engine while preserving every route, status, template, form, upload path, and browser flow.

**Architecture:** Keep the existing `main.go` application assembly, Fiber controllers, global repositories/configuration, pgx persistence, server-rendered templates, and filesystem uploads. Add only two small test seams in `main.go`, characterize migration-sensitive behavior without external services, then manually adapt the Fiber APIs in one buildable migration slice.

**Tech Stack:** Go 1.26, current stable Fiber v3 (design-time snapshot `v3.4.0`), current stable Fiber HTML v3 (design-time snapshot `v3.0.6`), pgx v5, PostgreSQL 18 for manual smoke only, Docker Compose, `govulncheck` v1.6.0.

## Global Constraints

- Recheck stable Fiber v3, Fiber HTML, and Go requirements through Context7, official release notes, and `go list -m` at implementation start.
- Keep Go 1.26 unless the rechecked stable Fiber release requires a later supported Go line.
- Preserve all existing routes, methods, status codes, templates, view data, form field names, static URLs, upload paths, response flow, and minimal-JS UI.
- Preserve current package boundaries and server-rendered architecture.
- Do not add repository mocks, dependency injection, an application factory, generalized compatibility wrappers, or another service layer.
- Do not add database migrations, schema changes, routes, features, CSS, JavaScript, or UI changes.
- Do not fix pre-existing bugs unless Fiber v3 causes a demonstrated regression.
- `go test ./...` must require no PostgreSQL, Docker, network service, browser, or external process.
- Do not add integration or end-to-end test infrastructure.
- Do not take another dependency major upgrade in Wave 2.
- Use `github.com/gofiber/template/html/v3` unless a demonstrated compile/render incompatibility requires the documented `html/v2` fallback.
- Never commit `.env`, uploads, binaries, credentials, temporary smoke files, or generated artifacts.

## File Map

**Create:**

- `main_test.go`: DB-free static-route and rate-limiter characterization.
- `templates/templates_test.go`: real template load/render compatibility.
- `docs/superpowers/modernization-followups.md`: confirmed remaining major/deprecated candidates, or an explicit no-findings record.

**Modify:**

- `main.go`: two small test seams, Fiber v3 imports, limiter callback, static middleware.
- `controllers/threads_controller.go`: Fiber v3 context signatures, route parameter parsing, request-context propagation, explicit 302 redirects.
- `controllers/captcha_controller.go`: Fiber v3 signature and raw request-context API.
- `controllers/threads_controller_test.go`: multipart characterization and Fiber v3 test types.
- `gateway/controllers.go`: Fiber v3 signature/import.
- `gateway/controllers_test.go`: Fiber v3 import.
- `go.mod`, `go.sum`: Fiber v3 and HTML v3 module graph.
- `README.md`, `AGENTS.md`: Fiber v3 baseline and testing guidance.

**Verify only unless a concrete failure requires an in-scope fix:**

- `.github/workflows/push.yml`: Go line and verification commands remain correct.
- `Dockerfile`: Go builder line remains correct and image still builds.
- `templates/**/*.html`, `static/**`, `migrations/**`, `repositories/**`, `db/**`, `utils/**`: behavior remains unchanged.

---

### Task 1: Add DB-Free Fiber Behavior Characterization

**Files:**

- Create: `main_test.go`
- Create: `templates/templates_test.go`
- Modify: `main.go`
- Modify: `controllers/threads_controller_test.go`

**Interfaces:**

- Consumes: current Fiber v2 `*fiber.Ctx`, `app.Static`, limiter middleware, real templates, `controllers.CreateThread` validation flow.
- Produces: `skipRateLimit(isLocal bool, path string) bool` and `registerStaticRoutes(app *fiber.App)`; DB-free tests that pass unchanged in behavior after their Fiber v3 type/import adaptation.

- [ ] **Step 1: Record the clean Wave 1 baseline**

Run:

```bash
git status --short --branch
go test ./...
go vet ./...
go build ./...
```

Expected: worktree is clean apart from the already committed spec/plan history;
all Go commands exit 0.

- [ ] **Step 2: Write failing tests for the two application seams**

Create `main_test.go` with the v2 imports used before dependency migration:

```go
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
```

Run:

```bash
go test . -run 'Test(RegisterStaticRoutes|SkipRateLimit|LimiterRejectsOrdinaryRemoteRequests)' -count=1
```

Expected: build fails with `undefined: registerStaticRoutes` and
`undefined: skipRateLimit`.

- [ ] **Step 3: Add the smallest application seams and use them in production assembly**

Add above `main()` in `main.go`:

```go
func skipRateLimit(isLocal bool, path string) bool {
	isRequestForStatic := strings.Contains(path, "/static") ||
		strings.Contains(path, "/uploads") ||
		strings.Contains(path, "/captcha")
	return isRequestForStatic || isLocal
}

func registerStaticRoutes(app *fiber.App) {
	app.Static("/uploads", "./uploads")
	app.Static("/static", "./static")
}
```

Replace the inline limiter predicate with:

```go
Next: func(c *fiber.Ctx) bool {
	return skipRateLimit(c.IsFromLocal(), c.Path())
},
```

Replace the two direct `app.Static` calls with:

```go
registerStaticRoutes(app)
```

Run:

```bash
go fmt ./...
go test . -run 'Test(RegisterStaticRoutes|SkipRateLimit|LimiterRejectsOrdinaryRemoteRequests)' -count=1
```

Expected: all three tests pass. Existing production limiter and static behavior
is unchanged.

- [ ] **Step 4: Write a DB-free invalid multipart controller test**

Add imports to `controllers/threads_controller_test.go`:

```go
"bytes"
"mime/multipart"
"net/http"

"micrach/config"
```

Append:

```go
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
```

Run:

```bash
go test ./controllers -run TestCreateThreadRejectsInvalidMultipartBeforeDatabase -count=1
```

Expected: PASS without an initialized database pool. The form includes both
keys, so it does not exercise the pre-existing missing-field panic.

- [ ] **Step 5: Add real template load/render characterization**

Create `templates/templates_test.go`:

```go
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
```

Run:

```bash
go test ./templates -run TestHTMLTemplatesLoadAndRender -count=1
go test ./...
go vet ./...
```

Expected: all commands pass with Fiber v2 and no external service.

- [ ] **Step 6: Commit the characterization slice**

```bash
git add main.go main_test.go controllers/threads_controller_test.go templates/templates_test.go
git diff --cached --check
git commit -m "test: characterize Fiber behavior"
```

Expected: focused commit containing only DB-free tests and the two small test
seams.

---

### Task 2: Migrate Fiber and HTML Template APIs

**Files:**

- Modify: `go.mod`
- Modify: `go.sum`
- Modify: `main.go`
- Modify: `main_test.go`
- Modify: `controllers/threads_controller.go`
- Modify: `controllers/captcha_controller.go`
- Modify: `controllers/threads_controller_test.go`
- Modify: `gateway/controllers.go`
- Modify: `gateway/controllers_test.go`
- Modify: `templates/templates_test.go`

**Interfaces:**

- Consumes: Task 1 characterization and existing repository methods accepting `context.Context`.
- Produces: Fiber v3 handlers using `fiber.Ctx`; explicit `strconv.Atoi` parameter errors; explicit 302 redirects; static v3 middleware; HTML v3 engine; no Fiber v2 direct dependency.

- [ ] **Step 1: Recheck and record current stable migration targets**

Use Context7 and official release pages from the design spec, then run:

```bash
go list -m -json github.com/gofiber/fiber/v3@latest
go list -m -json github.com/gofiber/template/html/v3@latest
```

Expected at plan-writing time:

```text
github.com/gofiber/fiber/v3 v3.4.0, GoVersion 1.25.0
github.com/gofiber/template/html/v3 v3.0.6, GoVersion 1.25.0
```

Requirements before continuing:

- both resolved versions are stable releases;
- both support Go 1.26;
- release notes contain no migration blocker affecting micrach;
- if either version differs from the snapshot, record the exact selected
  version in the implementation report and use it in every later command.

- [ ] **Step 2: Add the selected v3 modules**

With the design-time versions shown here, run:

```bash
go get github.com/gofiber/fiber/v3@v3.4.0
go get github.com/gofiber/template/html/v3@v3.0.6
```

If Step 1 found newer stable versions, replace only the two version suffixes
with those exact recorded stable versions.

Expected: commands exit 0. Fiber v2 remains temporarily because source imports
have not yet changed.

- [ ] **Step 3: Move application assembly and its tests to Fiber v3**

In `main.go`, replace imports:

```go
"github.com/gofiber/fiber/v3"
"github.com/gofiber/fiber/v3/middleware/compress"
"github.com/gofiber/fiber/v3/middleware/limiter"
"github.com/gofiber/fiber/v3/middleware/recover"
"github.com/gofiber/fiber/v3/middleware/static"
"github.com/gofiber/template/html/v3"
```

Replace `registerStaticRoutes` with:

```go
func registerStaticRoutes(app *fiber.App) {
	app.Get("/uploads/*", static.New("./uploads"))
	app.Get("/static/*", static.New("./static"))
}
```

Change the limiter callback only:

```go
Next: func(c fiber.Ctx) bool {
	return skipRateLimit(c.IsFromLocal(), c.Path())
},
```

Keep middleware order, `Max: 50`, route order, engine functions, app config,
and listen behavior unchanged.

In `main_test.go`, replace Fiber and limiter imports with `/v3`, then change all
test handlers and callbacks from `func(c *fiber.Ctx)` to `func(c fiber.Ctx)`.
Do not pass an explicit `fiber.TestConfig` to `app.Test`.

In `templates/templates_test.go`, replace the HTML import with:

```go
"github.com/gofiber/template/html/v3"
```

Run:

```bash
go test . ./templates -count=1
```

Expected at this intermediate point: main/template packages compile and their
tests pass; other packages may still fail until their v2 handler imports are
migrated in the next steps.

- [ ] **Step 4: Migrate thread handlers without changing business flow**

In `controllers/threads_controller.go`, replace the Fiber import with:

```go
"github.com/gofiber/fiber/v3"
```

Change signatures exactly:

```go
func GetThreads(c fiber.Ctx) error
func GetThread(c fiber.Ctx) error
func CreateThread(c fiber.Ctx) error
func UpdateThread(c fiber.Ctx) error
```

Delete every:

```go
ctx := c.UserContext()
```

Pass `c` directly as `context.Context` at every former `ctx` call site. This
includes repository calls, `db.Pool.Acquire`, `conn.Begin`, deferred
`tx.Rollback`, and `tx.Commit`. Representative exact forms:

```go
threads, err := repositories.Posts.Get(c, limit, offset)
conn, err := db.Pool.Acquire(c)
tx, err := conn.Begin(c)
rollbackErr := tx.Rollback(c)
if err := tx.Commit(c); err != nil {
```

Replace both removed `ParamsInt` calls with error-preserving parsing:

```go
threadID, err := strconv.Atoi(c.Params("threadID"))
if err != nil {
	return c.Status(fiber.StatusNotFound).Render("pages/404", nil)
}
```

Replace the create-thread redirect with:

```go
path := "/" + strconv.Itoa(threadID)
return c.Redirect().Status(fiber.StatusFound).To(path)
```

Replace the reply redirect with:

```go
path := "/" + strconv.Itoa(threadID) + "#" + strconv.Itoa(postID)
return c.Redirect().Status(fiber.StatusFound).To(path)
```

Do not use Fiber v3's default redirect, because it is 303. Do not use
`fiber.Params[int]`, because its fallback value would change invalid-ID flow.

- [ ] **Step 5: Migrate CAPTCHA and gateway handlers**

In `controllers/captcha_controller.go`, change the import and signature to:

```go
"github.com/gofiber/fiber/v3"

func GetCaptcha(c fiber.Ctx) error {
```

Replace the old raw fasthttp access:

```go
c.RequestCtx().SetContentType("image/png")
```

Do not use `c.Context().SetContentType`; Fiber v3 `Context()` now returns a
standard `context.Context`.

In `gateway/controllers.go`, change the import and signature to:

```go
"github.com/gofiber/fiber/v3"

func Ping(c fiber.Ctx) error {
```

Keep header selection, JSON, and statuses unchanged.

- [ ] **Step 6: Migrate controller and gateway tests**

In `controllers/threads_controller_test.go` and
`gateway/controllers_test.go`, change imports from Fiber v2 to Fiber v3.

Change the test view signature to the v3-compatible spelling:

```go
func (testViews) Render(out io.Writer, name string, _ any, _ ...string) error {
```

Keep every `app.Test(req)` call without an explicit config. No test may
initialize `db.Pool` or start an external service.

- [ ] **Step 7: Normalize the module graph and prove old APIs are gone**

Run:

```bash
go mod tidy
go fmt ./...
rg -n 'gofiber/fiber/v2|gofiber/template/html/v2|\*fiber\.Ctx|UserContext\(|ParamsInt\(|app\.Static\(|c\.Redirect\([^)]' --glob '*.go' go.mod
```

Expected: final `rg` returns no matches. `go.mod` has direct requirements for
Fiber `/v3` and HTML `/v3`; old Fiber v2 and HTML v2 direct requirements are
absent.

Run:

```bash
go test ./...
go vet ./...
go build ./...
```

Expected: all commands pass without PostgreSQL or Docker. The real template
render test proves HTML v3 compatibility, so fallback is not used.

If the HTML v3 compile/render test fails specifically inside the engine, record
the exact failure, restore `github.com/gofiber/template/html/v2@v2.1.3`, repeat
`go mod tidy`, and rerun this entire step. Do not fallback for an unrelated
application failure.

- [ ] **Step 8: Review and commit the Fiber v3 migration**

Run:

```bash
git diff -- go.mod go.sum main.go main_test.go controllers gateway templates/templates_test.go
git diff --check
```

Review requirements:

- no route or middleware order changed;
- both redirects specify `fiber.StatusFound`;
- invalid IDs fail before DB access;
- no form, upload, repository, template, or CSS behavior was redesigned.

Commit:

```bash
git add go.mod go.sum main.go main_test.go controllers gateway templates/templates_test.go
git commit -m "build: migrate to Fiber v3"
```

---

### Task 3: Audit Dependencies and Document the New Baseline

**Files:**

- Create: `docs/superpowers/modernization-followups.md`
- Modify: `go.mod` only if a newly available compatible patch/minor release is selected
- Modify: `go.sum` only with an approved module change
- Modify: `README.md`
- Modify: `AGENTS.md`

**Interfaces:**

- Consumes: passing Fiber v3 graph from Task 2.
- Produces: reviewed direct dependencies, accurate Fiber v3 docs, and an explicit confirmed follow-up record.

- [ ] **Step 1: Audit direct and transitive modules**

Run:

```bash
go list -u -m -json all > /tmp/micrach-wave2-modules.json
jq -r 'select(.Main != true and .Indirect != true) | [.Path, .Version, (.Update.Version // "current"), (.GoVersion // "unknown")] | @tsv' /tmp/micrach-wave2-modules.json
go mod graph | sort > /tmp/micrach-wave2-module-graph.txt
go mod why -m github.com/gofiber/fiber/v3
go mod why -m github.com/gofiber/template/html/v3
```

Expected from the design-time audit: `captcha v1.1.0`, `imaging v1.6.2`,
`pgx/v5 v5.10.0`, and `godotenv v1.5.1` report `current`.

If a direct dependency has a newer patch/minor release, inspect its official
release notes, select its exact version, run `go get module@version`, then
repeat `go mod tidy`, `go test ./...`, and `go vet ./...`. Do not run a blanket
`go get -u`, and do not accept another major.

- [ ] **Step 2: Audit deprecated and remaining major APIs**

Run:

```bash
go vet ./...
go run golang.org/x/vuln/cmd/govulncheck@v1.6.0 ./...
rg -n 'Deprecated:|gofiber/fiber/v2|gofiber/template/html/v2|io/ioutil|context[.]TO[D]O' --glob '*.go' --glob 'go.mod' .
```

Expected: vet and vulnerability scan exit 0; source scan finds no old Fiber,
deprecated I/O, or broad placeholder contexts. Comments unrelated to an imported
deprecated API do not become Wave 2 work.

- [ ] **Step 3: Write the confirmed follow-up record**

If the audit matches the expected no-findings state, create
`docs/superpowers/modernization-followups.md`:

```markdown
# Modernization Follow-ups

Updated: 2026-07-16

Wave 2 completed the confirmed Fiber v3 and template-engine major upgrades.
The final direct-module, deprecated-API, `go vet`, and `govulncheck` audits found
no remaining major-version or deprecated-API candidate requiring a scheduled
follow-up.

Recheck this file during the next deliberate dependency refresh. Do not add
speculative architecture cleanup or ordinary patch/minor updates here.
```

If Step 1 or Step 2 proves a real deferred candidate, replace the no-findings
paragraph with a concrete entry containing its package/API, current version or
usage, candidate version/replacement, official evidence URL, and precise reason
for deferral. Do not write placeholder markers, guesses, or unrelated
refactoring ideas.

- [ ] **Step 4: Update README and AGENTS**

In `README.md`, change the stack item exactly:

```markdown
- Fiber v3
```

In `AGENTS.md`, change the project overview to:

```markdown
micrach is a small, server-rendered, single-board imageboard. It uses Go,
Fiber v3, Fiber HTML templates, PostgreSQL through `pgx`, and filesystem-backed
image uploads. It intentionally uses little browser-side JavaScript.
```

Add this testing rule under `Testing Expectations` in `AGENTS.md`:

```markdown
Routine `go test ./...` must remain DB-free. Fiber handler and middleware tests
must use in-memory apps and temporary directories; PostgreSQL-backed behavior
belongs to the explicit manual smoke pass unless a later design approves an
integration-test harness.
```

Keep Go 1.26, PostgreSQL 18, Compose, Docker, CI, and all other documentation
unchanged unless the rechecked Fiber minimum forced a Go-line change.

- [ ] **Step 5: Verify and commit dependency documentation**

Run:

```bash
go mod tidy
go test ./...
go vet ./...
rg -n 'Fiber v2|gofiber/fiber/v2|gofiber/template/html/v2' README.md AGENTS.md go.mod
rg -n 'gofiber/fiber/v2|gofiber/template/html/v2' --glob '*.go' .
git diff --check
```

Expected: tests/vet pass; `rg` returns no stale baseline; module files remain
stable after tidy.

Commit:

```bash
git add go.mod go.sum README.md AGENTS.md docs/superpowers/modernization-followups.md
git commit -m "docs: record Fiber v3 baseline"
```

If `go.mod` and `go.sum` did not change in this task, omit them from `git add`.

---

### Task 4: Run the Automated Wave 2 Gate

**Files:**

- Verify only; modify only the file responsible for a concrete in-scope failure.

**Interfaces:**

- Consumes: completed migration and documentation tasks.
- Produces: fresh automated evidence with no external service used by Go tests.

- [ ] **Step 1: Verify repository and module hygiene**

Run:

```bash
git status --short --branch
git diff --check
git ls-files .env uploads micrach '*.out' '*.test'
go mod tidy
git diff --exit-code -- go.mod go.sum
go list -m all | rg 'github.com/gofiber/(fiber/v3|template/html/v3)'
```

Expected:

- worktree is clean;
- no forbidden artifact is tracked;
- tidy changes nothing;
- Fiber v3 and HTML v3 appear in the graph.

Run:

```bash
if go list -m all | rg -q 'github.com/gofiber/fiber/v2|github.com/gofiber/template/html/v2'; then
  echo 'obsolete direct/graph dependency remains'
  exit 1
fi
```

Expected: exit 0 for the primary HTML v3 path. Do not reject
`github.com/gofiber/template/v2` when it is a legitimate transitive dependency
of HTML v3; it is a different module from `template/html/v2`.

- [ ] **Step 2: Run the full Go and vulnerability gate**

Run:

```bash
go fmt ./...
test -z "$(gofmt -l .)"
go test ./...
go vet ./...
go build ./...
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /tmp/micrach-linux-amd64 .
go run golang.org/x/vuln/cmd/govulncheck@v1.6.0 ./...
```

Expected: every command exits 0; no test starts or requires PostgreSQL, Docker,
network access, a browser, or another external process; no reachable
vulnerability remains unhandled.

- [ ] **Step 3: Build and inspect the container**

Run:

```bash
docker build --tag micrach:wave2 .
docker image inspect micrach:wave2 --format '{{.Config.User}} {{json .Config.Entrypoint}}'
docker run --rm --entrypoint sh micrach:wave2 -c 'touch /app/uploads/write-test && test -f /app/uploads/write-test'
```

Expected inspect output:

```text
micrach ["./micrach"]
```

Expected: image build and writable-upload check exit 0. No image is published.

- [ ] **Step 4: Review the final code diff against migration constraints**

Run:

```bash
git diff origin/master...HEAD --stat
git diff origin/master...HEAD -- main.go controllers gateway go.mod README.md AGENTS.md
rg -n 'gofiber/fiber/v2|\*fiber\.Ctx|UserContext\(|ParamsInt\(|app\.Static\(' --glob '*.go' .
git status --short --branch
```

Expected:

- final `rg` returns no matches;
- no repository, schema, migration, template, CSS, or JavaScript redesign is in
  the diff;
- branch contains focused, buildable commits and a clean worktree.

---

### Task 5: Run the Manual Application Smoke Matrix

**Files:**

- Verify only; use disposable runtime state and `/tmp` artifacts.

**Interfaces:**

- Consumes: verified Wave 2 binary, local Compose PostgreSQL, existing routes and forms.
- Produces: recorded operational evidence for catalog, forms, redirects, CAPTCHA, uploads, limiter, `sage`, bump limit, archival, static files, and responsive pages. This task creates no automated integration test.

- [ ] **Step 1: Start PostgreSQL and create an isolated smoke database**

Run:

```bash
docker compose up -d db
docker compose ps
SMOKE_DB="micrach_wave2_$(date +%s)"
docker compose exec -T db createdb -U micrach "$SMOKE_DB"
printf '%s\n' "$SMOKE_DB" > /tmp/micrach-wave2-smoke-db
```

Expected: Compose `db` is healthy and the uniquely named database is created.
Do not run `docker compose down -v`; existing developer data must not be
deleted.

- [ ] **Step 2: Generate disposable upload fixtures**

Run on macOS:

```bash
cp static/images/errors/404.png /tmp/micrach-wave2.png
sips -s format jpeg static/images/errors/404.png --out /tmp/micrach-wave2.jpg
cp /tmp/micrach-wave2.png /tmp/micrach-wave2-oversize.png
dd if=/dev/zero bs=1048576 count=4 >> /tmp/micrach-wave2-oversize.png 2>/dev/null
file /tmp/micrach-wave2.png /tmp/micrach-wave2.jpg
```

Expected: `file` identifies one PNG and one JPEG; oversize fixture exceeds 3
MiB. These remain outside the repository.

- [ ] **Step 3: Start the application with smoke-specific limits**

Run:

```bash
SMOKE_DB="$(cat /tmp/micrach-wave2-smoke-db)"
ENV=development \
PORT=3100 \
POSTGRES_URL="postgres://micrach:micrach@localhost:5432/${SMOKE_DB}?pool_max_conns=5" \
IS_DB_SEEDED=false \
IS_CAPTCHA_ACTIVE=false \
IS_RATE_LIMITER_ENABLED=false \
GATEWAY_URL= \
THREADS_MAX_COUNT=2 \
THREAD_BUMP_LIMIT=3 \
go run . >/tmp/micrach-wave2-app.log 2>&1 &
APP_PID=$!
printf '%s\n' "$APP_PID" > /tmp/micrach-wave2-app.pid
```

Poll without a long blocking sleep:

```bash
for attempt in 1 2 3 4 5 6 7 8 9 10; do
  if curl -fsS http://localhost:3100/ >/dev/null; then break; fi
  sleep 1
done
curl -fsS http://localhost:3100/ >/dev/null
```

Expected: migrations apply to the isolated DB and catalog returns 200.

- [ ] **Step 4: Verify catalog, errors, static files, and CAPTCHA media**

Run:

```bash
curl -sS -o /tmp/micrach-wave2-catalog.html -w '%{http_code}\n' http://localhost:3100/
curl -sS -o /tmp/micrach-wave2-page0.html -w '%{http_code}\n' 'http://localhost:3100/?page=0'
curl -sS -o /tmp/micrach-wave2-invalid.html -w '%{http_code}\n' http://localhost:3100/invalid
curl -sS -o /tmp/micrach-wave2-css -w '%{http_code}\n' http://localhost:3100/static/styles/index.css
```

Expected status sequence:

```text
200
404
404
200
```

Open the catalog in a browser at narrow and wide widths. Expected: form,
catalog, pagination area, header, and footer have no visual regression.

CAPTCHA generation uses a live ID embedded in the form only when CAPTCHA is
enabled. Restart the application later in Step 8 with CAPTCHA enabled and
verify the returned image reports `Content-Type: image/png`.

- [ ] **Step 5: Verify thread/reply redirects, uploads, `sage`, bump limit, and archival**

Create thread A with PNG and JPEG:

```bash
curl -sS -D /tmp/micrach-wave2-thread-a.headers -o /tmp/micrach-wave2-thread-a.body \
  -F 'title=thread A' -F 'text=first thread' \
  -F 'files=@/tmp/micrach-wave2.png;type=image/png' \
  -F 'files=@/tmp/micrach-wave2.jpg;type=image/jpeg' \
  http://localhost:3100/
rg -n '^HTTP/|^Location:' /tmp/micrach-wave2-thread-a.headers
```

Expected: `302` and `Location: /<threadAID>`. Save the numeric ID:

```bash
THREAD_A="$(awk '/^Location:/ {gsub("/|\\r", "", $2); print $2}' /tmp/micrach-wave2-thread-a.headers)"
test -n "$THREAD_A"
```

Create thread B the same way without files, then add one normal reply to A:

```bash
curl -sS -D /tmp/micrach-wave2-thread-b.headers -o /dev/null \
  -F 'title=thread B' -F 'text=second thread' http://localhost:3100/
THREAD_B="$(awk '/^Location:/ {gsub("/|\\r", "", $2); print $2}' /tmp/micrach-wave2-thread-b.headers)"

curl -sS -D /tmp/micrach-wave2-reply-a.headers -o /dev/null \
  -F 'text=normal reply' "http://localhost:3100/${THREAD_A}"
rg -n '^HTTP/|^Location:' /tmp/micrach-wave2-reply-a.headers
```

Expected: reply returns `302`; `Location` matches
`/<threadAID>#<newPostID>`. Catalog shows A bumped above B.

Add a `sage` reply to B:

```bash
curl -sS -D /tmp/micrach-wave2-reply-b.headers -o /dev/null \
  -F 'text=sage reply' -F 'sage=on' "http://localhost:3100/${THREAD_B}"
```

Expected: `302` with fragment; B does not move above A.

Create thread C:

```bash
curl -sS -D /tmp/micrach-wave2-thread-c.headers -o /dev/null \
  -F 'title=thread C' -F 'text=third thread' http://localhost:3100/
rg -n '^HTTP/|^Location:' /tmp/micrach-wave2-thread-c.headers
```

Expected: `302`; because `THREADS_MAX_COUNT=2`, the oldest active thread is
archived. At this point B is the oldest active thread because A was bumped and
B's `sage` reply did not bump it. Attempting to reply to B returns 400 with the
existing `THREAD IS ARCHIVED` page.

Add another normal reply to A. With `THREAD_BUMP_LIMIT=3`, A now reaches the
limit and must not move above C. Verify C remains first in the catalog.

Create thread D. Expected: A is now the oldest active thread and becomes
archived; replying to A returns 400 with the existing `THREAD IS ARCHIVED`
page. This separately proves bump-limit and archival behavior.

Inspect thread A in the browser. Expected: both original images and thumbnails
load from unchanged `/uploads/<threadID>/...` paths.

- [ ] **Step 6: Verify validation boundaries without external test infrastructure**

Run:

```bash
curl -sS -o /tmp/micrach-wave2-oversize.body -w '%{http_code}\n' \
  -F 'title=oversize' -F 'text=oversize upload' \
  -F 'files=@/tmp/micrach-wave2-oversize.png;type=image/png' \
  http://localhost:3100/

curl -sS -o /tmp/micrach-wave2-four-files.body -w '%{http_code}\n' \
  -F 'title=four files' -F 'text=four uploads' \
  -F 'files=@/tmp/micrach-wave2.png;type=image/png' \
  -F 'files=@/tmp/micrach-wave2.png;type=image/png' \
  -F 'files=@/tmp/micrach-wave2.png;type=image/png' \
  -F 'files=@/tmp/micrach-wave2.png;type=image/png' \
  http://localhost:3100/
```

Expected: oversize request returns 400 with the existing file-size message;
four valid small files remain accepted and return 302. If Fiber rejects the
request before application validation due solely to changed v3 multipart body
enforcement, capture status/log evidence and add only the smallest explicit
`fiber.Config.BodyLimit` needed to preserve the verified Wave 1 contract, then
rerun Tasks 2-5 gates.

- [ ] **Step 7: Verify limiter behavior from a non-loopback client**

Stop the current app, then restart it with rate limiting enabled and other
smoke settings unchanged:

```bash
kill "$(cat /tmp/micrach-wave2-app.pid)"
wait "$(cat /tmp/micrach-wave2-app.pid)" 2>/dev/null || true
SMOKE_DB="$(cat /tmp/micrach-wave2-smoke-db)"
ENV=development PORT=3100 \
POSTGRES_URL="postgres://micrach:micrach@localhost:5432/${SMOKE_DB}?pool_max_conns=5" \
IS_DB_SEEDED=false IS_CAPTCHA_ACTIVE=false IS_RATE_LIMITER_ENABLED=true \
GATEWAY_URL= \
THREADS_MAX_COUNT=2 THREAD_BUMP_LIMIT=3 \
go run . >/tmp/micrach-wave2-app.log 2>&1 &
printf '%s\n' "$!" > /tmp/micrach-wave2-app.pid
```

Local requests are intentionally bypassed, so the DB-free unit test is the
authoritative 429 check for remote requests. Manually confirm repeated local
catalog requests and `/static`, `/uploads`, and `/captcha` paths remain
unlimited. Do not weaken `IsFromLocal` or trust proxy headers merely to force a
manual 429.

- [ ] **Step 8: Verify pagination with the same disposable database**

Restart with CAPTCHA and limiter disabled and `THREADS_MAX_COUNT=20`:

```bash
kill "$(cat /tmp/micrach-wave2-app.pid)"
wait "$(cat /tmp/micrach-wave2-app.pid)" 2>/dev/null || true
SMOKE_DB="$(cat /tmp/micrach-wave2-smoke-db)"
ENV=development PORT=3100 \
POSTGRES_URL="postgres://micrach:micrach@localhost:5432/${SMOKE_DB}?pool_max_conns=5" \
IS_DB_SEEDED=false IS_CAPTCHA_ACTIVE=false IS_RATE_LIMITER_ENABLED=false \
GATEWAY_URL= THREADS_MAX_COUNT=20 THREAD_BUMP_LIMIT=3 \
go run . >/tmp/micrach-wave2-app.log 2>&1 &
printf '%s\n' "$!" > /tmp/micrach-wave2-app.pid
```

Create eleven small text-only threads:

```bash
for number in 1 2 3 4 5 6 7 8 9 10 11; do
  curl -fsS -o /dev/null \
    -F "title=pagination ${number}" -F "text=page fixture ${number}" \
    http://localhost:3100/
done
curl -sS -o /tmp/micrach-wave2-page2.html -w '%{http_code}\n' 'http://localhost:3100/?page=2'
curl -sS -o /tmp/micrach-wave2-page3.html -w '%{http_code}\n' 'http://localhost:3100/?page=3'
```

Expected: page 2 returns 200 and contains catalog content; page 3 returns 404;
pagination links navigate between pages 1 and 2.

- [ ] **Step 9: Verify CAPTCHA and responsive browser behavior**

Restart exactly:

```bash
kill "$(cat /tmp/micrach-wave2-app.pid)"
wait "$(cat /tmp/micrach-wave2-app.pid)" 2>/dev/null || true
SMOKE_DB="$(cat /tmp/micrach-wave2-smoke-db)"
ENV=development PORT=3100 \
POSTGRES_URL="postgres://micrach:micrach@localhost:5432/${SMOKE_DB}?pool_max_conns=5" \
IS_DB_SEEDED=false IS_CAPTCHA_ACTIVE=true IS_RATE_LIMITER_ENABLED=false \
GATEWAY_URL= THREADS_MAX_COUNT=20 THREAD_BUMP_LIMIT=3 \
go run . >/tmp/micrach-wave2-app.log 2>&1 &
printf '%s\n' "$!" > /tmp/micrach-wave2-app.pid
```

Open the catalog in a browser, then:

1. confirm CAPTCHA image loads with `Content-Type: image/png`;
2. submit a wrong answer and confirm 400 with `INVALID CAPTCHA`;
3. solve a refreshed CAPTCHA and confirm valid thread creation returns 302;
4. inspect catalog and thread pages at narrow and wide viewport widths;
5. confirm no template, CSS, static, upload, or browser-flow regression.

This remains a manual smoke pass. Do not add browser automation or CAPTCHA
test hooks.

- [ ] **Step 10: Clean disposable state and record results**

Run:

```bash
if test -f /tmp/micrach-wave2-app.pid; then
  kill "$(cat /tmp/micrach-wave2-app.pid)" 2>/dev/null || true
fi
SMOKE_DB="$(cat /tmp/micrach-wave2-smoke-db)"
docker compose exec -T db dropdb -U micrach "$SMOKE_DB"
rm -f /tmp/micrach-wave2-* /tmp/micrach-linux-amd64
git status --short --branch
```

Expected: isolated database and `/tmp` fixtures are removed; existing Compose
database/volume and `.env` remain untouched; Git worktree is clean.

Record in the pull request:

- selected Fiber, HTML, and Go versions;
- HTML v3 success or exact fallback reason;
- all automated commands and results;
- manual status/redirect/static/upload/CAPTCHA/limiter/`sage`/bump/archive and
  responsive results;
- confirmed follow-ups or explicit no-findings result;
- statement that no DB-backed/integration/E2E tests were added.
