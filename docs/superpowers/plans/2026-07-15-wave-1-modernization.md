# Wave 1 Modernization Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Establish a supported Go 1.26 baseline with current Fiber v2 and pgx v5 dependencies, local PostgreSQL infrastructure, build-only CI, a reproducible container image, and focused regression tests without changing application behavior or architecture.

**Architecture:** Preserve the existing Fiber controllers, global configuration and repository objects, PostgreSQL schema, server-rendered templates, and filesystem upload layout. Apply compatibility changes in reviewable slices: characterize behavior first, update the HTTP/toolchain dependencies, migrate pgx with explicit context propagation, modernize deprecated APIs, then replace development and CI infrastructure. Fiber v3 remains a separate Wave 2 implementation after this plan is fully green.

**Tech Stack:** Go 1.26.5, Fiber v2.52.14, Fiber HTML template v2.1.3, pgx v5.10.0, PostgreSQL 18.4, Docker Compose, GitHub Actions, Dependabot, and `govulncheck` v1.6.0.

## Global Constraints

- Preserve all routes, templates, HTTP statuses, database schema, upload paths, and minimal-JavaScript behavior.
- Preserve existing package boundaries; do not redesign handlers or repositories.
- Do not add database migrations, Redis, a proxy, a queue, or another service.
- Run the application on the host with `go run .` or `make dev`; Compose runs PostgreSQL only.
- Keep Fiber on v2 throughout Wave 1.
- Do not publish container images or deploy from CI.
- Never use a production or shared PostgreSQL database for tests.
- Do not commit `.env`, uploads, binaries, credentials, or generated test artifacts.

## File Map

**Create:**

- `config/config_test.go`: configuration default and parsing characterization.
- `utils/utils_test.go`: post and upload validation boundaries.
- `controllers/threads_controller_test.go`: DB-free HTTP error-path smoke tests.
- `repositories/files_repository_test.go`: pgx v5 integer-array conversion.
- `db/db_test.go`: PostgreSQL undefined-table error classification.
- `.dockerignore`: exclude local and generated state from image context.
- `.github/dependabot.yml`: reviewed Go module and Actions updates.

**Modify:**

- `go.mod`, `go.sum`: Go line and dependency graph.
- `main.go`: Fiber template v2 import, startup context, and contextual startup calls.
- `controllers/threads_controller.go`: request context propagation and checked transaction completion.
- `repositories/posts_repository.go`: pgx v5 imports, contextual methods, and rows lifecycle.
- `repositories/files_repository.go`: pgx v5 imports, contextual methods, and native integer-array encoding.
- `repositories/seeds.go`: explicit startup context.
- `db/db.go`: pgx v5 pool, startup context, rows lifecycle, and SQLSTATE handling.
- `build/css.go`, `files/get_files_in_folder.go`, `gateway/connect.go`: replace deprecated `ioutil`; close gateway response bodies.
- `Dockerfile`: multi-stage non-root application image.
- `docker-compose.yml`: PostgreSQL-only local infrastructure.
- `.env.example`: safe local defaults matching Compose.
- `Makefile`: host application and infrastructure commands.
- `.github/workflows/push.yml`: read-only verification CI.
- `README.md`, `AGENTS.md`: supported versions and workflows.

**Delete:**

- `.github/workflows/codesee-arch-diagram.yml`: obsolete CodeSee integration.

---

### Task 1: Add Behavior Characterization Tests

**Files:**

- Create: `config/config_test.go`
- Create: `utils/utils_test.go`
- Create: `controllers/threads_controller_test.go`

**Interfaces:**

- Consumes: existing unexported configuration helpers, `utils.ValidatePost`, `controllers.GetThreads`, and `controllers.GetThread`.
- Produces: regression tests that must pass before and after dependency migration.

- [x] **Step 1: Record the pre-change baseline**

Run:

```bash
go test ./...
go vet ./...
go build ./...
```

Expected: all commands exit 0; packages report no test files.

- [x] **Step 2: Add configuration characterization tests**

Create `config/config_test.go`:

```go
package config

import "testing"

func TestGetValueOrDefaultBoolean(t *testing.T) {
	tests := []struct {
		name         string
		value        string
		defaultValue bool
		want         bool
	}{
		{name: "default true", defaultValue: true, want: true},
		{name: "default false", defaultValue: false, want: false},
		{name: "explicit true", value: "true", want: true},
		{name: "explicit false", value: "false", defaultValue: true, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getValueOrDefaultBoolean(tt.value, tt.defaultValue); got != tt.want {
				t.Fatalf("getValueOrDefaultBoolean() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetValueOrDefaultInt(t *testing.T) {
	tests := []struct {
		name         string
		value        string
		defaultValue int
		want         int
	}{
		{name: "default", defaultValue: 3000, want: 3000},
		{name: "configured", value: "8080", defaultValue: 3000, want: 8080},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getValueOrDefaultInt(tt.value, tt.defaultValue); got != tt.want {
				t.Fatalf("getValueOrDefaultInt() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestDefaultConfigs(t *testing.T) {
	for _, key := range []string{
		"ENV", "PORT", "IS_DB_SEEDED", "IS_RATE_LIMITER_ENABLED",
		"THREADS_MAX_COUNT", "THREAD_BUMP_LIMIT", "IS_CAPTCHA_ACTIVE",
		"GATEWAY_URL", "GATEWAY_API_KEY", "GATEWAY_BOARD_ID",
		"GATEWAY_BOARD_URL", "GATEWAY_BOARD_DESCRIPTION", "POSTGRES_URL",
	} {
		t.Setenv(key, "")
	}

	app := getAppConfig()
	if app.Env != "release" || app.Port != 3000 || app.IsDbSeeded {
		t.Fatalf("unexpected app defaults: %+v", app)
	}
	if !app.IsRateLimiterEnabled || app.ThreadsMaxCount != 50 || app.ThreadBumpLimit != 500 || !app.IsCaptchaActive {
		t.Fatalf("unexpected feature defaults: %+v", app)
	}
	if db := getDbConfig(); db.Url != "postgresql://localhost/micrach" {
		t.Fatalf("unexpected database default: %+v", db)
	}
}
```

- [x] **Step 3: Add post and upload validation tests**

Create `utils/utils_test.go`:

```go
package utils

import (
	"mime/multipart"
	"net/textproto"
	"strings"
	"testing"

	"micrach/repositories"
)

func fileHeader(contentType string, size int64) *multipart.FileHeader {
	header := make(textproto.MIMEHeader)
	header.Set("Content-Type", contentType)
	return &multipart.FileHeader{Filename: "upload", Header: header, Size: size}
}

func TestValidatePost(t *testing.T) {
	validPNG := fileHeader("image/png", 1)
	validJPEG := fileHeader("image/jpeg", FILE_SIZE_IN_BYTES)

	tests := []struct {
		name  string
		title string
		text  string
		files []*multipart.FileHeader
		want  string
	}{
		{name: "text only", text: "hello"},
		{name: "file only", files: []*multipart.FileHeader{validPNG}},
		{name: "empty", want: repositories.InvalidTextOrFilesErrorMessage},
		{name: "title over limit", title: strings.Repeat("я", 101), text: "ok", want: repositories.InvalidTitleLengthErrorMessage},
		{name: "text over limit", text: strings.Repeat("я", 1001), want: repositories.InvalidTextLengthErrorMessage},
		{name: "four files", files: []*multipart.FileHeader{validPNG, validPNG, validJPEG, validJPEG}},
		{name: "five files", files: []*multipart.FileHeader{validPNG, validPNG, validPNG, validPNG, validPNG}, want: repositories.InvalidFilesLengthErrorMessage},
		{name: "unsupported content type", files: []*multipart.FileHeader{fileHeader("image/gif", 1)}, want: repositories.InvalidFileExtErrorMessage},
		{name: "size at limit", files: []*multipart.FileHeader{validJPEG}},
		{name: "size over limit", files: []*multipart.FileHeader{fileHeader("image/jpeg", FILE_SIZE_IN_BYTES+1)}, want: repositories.InvalidFileSizeErrorMessage},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidatePost(tt.title, tt.text, tt.files); got != tt.want {
				t.Fatalf("ValidatePost() = %q, want %q", got, tt.want)
			}
		})
	}
}
```

- [x] **Step 4: Add DB-free Fiber error-path smoke tests**

Create `controllers/threads_controller_test.go`:

```go
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
```

- [x] **Step 5: Run the new safety net**

Run:

```bash
go fmt ./...
go test ./...
go vet ./...
```

Expected: all tests pass and vet exits 0.

- [x] **Step 6: Commit the safety net**

```bash
git add config/config_test.go utils/utils_test.go controllers/threads_controller_test.go
git commit -m "test: characterize legacy behavior"
```

---

### Task 2: Upgrade Go, Fiber v2, and Non-Database Dependencies

**Files:**

- Modify: `go.mod`
- Modify: `go.sum`
- Modify: `main.go`

**Interfaces:**

- Consumes: characterization tests from Task 1.
- Produces: Go 1.26 module baseline using Fiber v2.52.14 and HTML template v2.1.3; no route or handler API change.

- [x] **Step 1: Set the module Go line**

Run:

```bash
go mod edit -go=1.26.0
go mod edit -toolchain=go1.26.5
```

Expected: `go.mod` contains:

```text
go 1.26.0

toolchain go1.26.5
```

- [x] **Step 2: Upgrade direct non-database dependencies deliberately**

Run:

```bash
go get github.com/gofiber/fiber/v2@v2.52.14
go get github.com/gofiber/template/html/v2@v2.1.3
go get github.com/dchest/captcha@v1.1.0
go get github.com/joho/godotenv@v1.5.1
go get github.com/disintegration/imaging@v1.6.2
```

Expected: each command exits 0; Fiber remains on major version 2.

- [x] **Step 3: Switch the template import to its current module path**

In `main.go`, replace:

```go
"github.com/gofiber/template/html"
```

with:

```go
"github.com/gofiber/template/html/v2"
```

Keep initialization unchanged:

```go
engine := html.New("./templates", ".html")
```

- [x] **Step 4: Normalize the module graph**

Run:

```bash
go mod tidy
go list -m all
```

Expected direct requirements:

```text
github.com/dchest/captcha v1.1.0
github.com/disintegration/imaging v1.6.2
github.com/gofiber/fiber/v2 v2.52.14
github.com/gofiber/template/html/v2 v2.1.3
github.com/joho/godotenv v1.5.1
```

The old root `github.com/gofiber/template` requirement must be absent.

- [x] **Step 5: Verify behavior on the new HTTP dependency baseline**

Run:

```bash
go test ./...
go vet ./...
go build ./...
```

Expected: all commands exit 0.

- [x] **Step 6: Commit the toolchain and Fiber v2 baseline**

```bash
git add go.mod go.sum main.go
git commit -m "build: update Go and Fiber v2"
```

---

### Task 3: Migrate pgx v4 to v5 and Propagate Contexts

**Files:**

- Create: `repositories/files_repository_test.go`
- Create: `db/db_test.go`
- Modify: `go.mod`
- Modify: `go.sum`
- Modify: `main.go`
- Modify: `controllers/threads_controller.go`
- Modify: `repositories/posts_repository.go`
- Modify: `repositories/files_repository.go`
- Modify: `repositories/seeds.go`
- Modify: `db/db.go`

**Interfaces:**

- Consumes: `fiber.Ctx.UserContext() context.Context`; startup uses `context.Background()`.
- Produces: repository and DB operations accepting `context.Context`; pgx v5 transaction types; `toInt32s([]int) []int32` for PostgreSQL `integer[]` parameters.

- [x] **Step 1: Write the integer-array conversion test**

Create `repositories/files_repository_test.go`:

```go
package repositories

import (
	"reflect"
	"testing"
)

func TestToInt32s(t *testing.T) {
	tests := []struct {
		name string
		in   []int
		want []int32
	}{
		{name: "nil", in: nil, want: nil},
		{name: "empty", in: []int{}, want: []int32{}},
		{name: "values", in: []int{1, 2, 2147483647}, want: []int32{1, 2, 2147483647}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := toInt32s(tt.in); !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("toInt32s() = %#v, want %#v", got, tt.want)
			}
		})
	}
}
```

- [x] **Step 2: Run the focused test and confirm the missing helper**

Run:

```bash
go test ./repositories -run TestToInt32s -count=1
```

Expected: build failure containing `undefined: toInt32s`.

- [x] **Step 3: Add the minimal array conversion helper**

Add to `repositories/files_repository.go`:

```go
func toInt32s(values []int) []int32 {
	if values == nil {
		return nil
	}

	result := make([]int32, len(values))
	for i, value := range values {
		result[i] = int32(value)
	}
	return result
}
```

Use it directly in the query:

```go
rows, err := Db.Pool.Query(ctx, sql, toInt32s(postIDs))
```

Do not keep `pgtype.Int4Array` or the standalone `github.com/jackc/pgtype` import.

- [x] **Step 4: Write SQLSTATE classification tests**

Create `db/db_test.go`:

```go
package db

import (
	"errors"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
)

func TestIsUndefinedTable(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{name: "undefined table", err: &pgconn.PgError{Code: "42P01"}, want: true},
		{name: "other postgres error", err: &pgconn.PgError{Code: "23505"}},
		{name: "ordinary error", err: errors.New("boom")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isUndefinedTable(tt.err); got != tt.want {
				t.Fatalf("isUndefinedTable() = %v, want %v", got, tt.want)
			}
		})
	}
}
```

- [x] **Step 5: Upgrade pgx and replace imports**

Run:

```bash
go get github.com/jackc/pgx/v5@v5.10.0
```

Replace imports exactly:

```go
"github.com/jackc/pgx/v4"
"github.com/jackc/pgx/v4/pgxpool"
```

with:

```go
"github.com/jackc/pgx/v5"
"github.com/jackc/pgx/v5/pgconn"
"github.com/jackc/pgx/v5/pgxpool"
```

Only import `pgconn` in `db/db.go` and `db/db_test.go`.

- [x] **Step 6: Define the contextual DB and repository signatures**

Change signatures to this exact set:

```go
// db/db.go
func Init(ctx context.Context)
func Migrate(ctx context.Context)
func runMigration(ctx context.Context, mid int, mname, msql string)
func getDbMigrations(ctx context.Context) MigrationsMap

// repositories/files_repository.go
func (r *FilesRepository) Create(ctx context.Context, f File) error
func (r *FilesRepository) GetByPostIDs(ctx context.Context, postIDs []int) (map[int][]File, error)
func (r *FilesRepository) CreateInTx(ctx context.Context, tx pgx.Tx, f File) (int, error)

// repositories/posts_repository.go
func (r *PostsRepository) Get(ctx context.Context, limit, offset int) ([]Post, error)
func (r *PostsRepository) GetThreadsCount(ctx context.Context) (int, error)
func (r *PostsRepository) Create(ctx context.Context, p Post) (int, error)
func (r *PostsRepository) GetThreadByPostID(ctx context.Context, ID int) ([]Post, error)
func (r *PostsRepository) GetIfThreadIsArchived(ctx context.Context, ID int) (bool, error)
func (r *PostsRepository) CreateInTx(ctx context.Context, tx pgx.Tx, p Post) (int, error)
func (r *PostsRepository) GetOldestThreadUpdatedAt(ctx context.Context) (time.Time, error)
func (r *PostsRepository) ArchiveThreadsFrom(ctx context.Context, t time.Time) error
func (r *PostsRepository) GetThreadPostsCount(ctx context.Context, id int) (int, error)
func (r *PostsRepository) BumpThreadInTx(ctx context.Context, tx pgx.Tx, id int) error

// repositories/seeds.go
func seedDb(ctx context.Context)
func Seed(ctx context.Context)
```

Replace every `context.TODO()` passed to pgx with the method's `ctx`. Pass `ctx`
through internal repository calls such as:

```go
filesMap, err := Files.GetByPostIDs(ctx, postIDs)
```

- [x] **Step 7: Propagate startup and request contexts at call sites**

In `main.go`:

```go
ctx := context.Background()
config.Init()
db.Init(ctx)
db.Migrate(ctx)
defer db.Pool.Close()

if config.App.IsDbSeeded {
	repositories.Seed(ctx)
}
```

At the beginning of every controller that reaches a repository, add:

```go
ctx := c.UserContext()
```

Pass `ctx` as the first argument to every repository call, `db.Pool.Acquire`,
`conn.Begin`, `tx.Rollback`, and `tx.Commit`.

- [x] **Step 8: Make rows lifecycle and migration error handling correct**

After every successful `Query`, add:

```go
defer rows.Close()
```

Remove pre-iteration `rows.Err()` checks. After each `for rows.Next()` loop,
repository methods returning a slice or map must return the final error:

```go
if err := rows.Err(); err != nil {
	return nil, err
}
```

In `getDbMigrations`, which preserves the existing panic-on-startup-error policy,
use this exact check before returning the map:

```go
if err := rows.Err(); err != nil {
	log.Panicln(err)
}
```

In `runMigration`, record the migration with `Exec`, not `Query`:

```go
_, err = Pool.Exec(ctx, `INSERT INTO migrations (id, name) VALUES ($1, $2)`, mid, mname)
if err != nil {
	log.Panicln(err)
}
```

Replace string matching with SQLSTATE classification:

```go
func isUndefinedTable(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "42P01"
}
```

In `getDbMigrations`, treat only that code as an empty migration map:

```go
rows, err := Pool.Query(ctx, `SELECT id, name FROM migrations`)
if err != nil {
	if isUndefinedTable(err) {
		return make(MigrationsMap)
	}
	log.Panicln(err)
}
defer rows.Close()
```

- [x] **Step 9: Check transaction completion in both write handlers**

Import `errors` and `github.com/jackc/pgx/v5` in
`controllers/threads_controller.go`. Replace each bare rollback defer with:

```go
defer func() {
	rollbackErr := tx.Rollback(ctx)
	if rollbackErr != nil && !errors.Is(rollbackErr, pgx.ErrTxClosed) {
		log.Println("error:", rollbackErr)
	}
}()
```

Replace each unchecked commit with:

```go
if err := tx.Commit(ctx); err != nil {
	log.Println("error:", err)
	return c.Status(fiber.StatusInternalServerError).Render("pages/500", nil)
}
```

- [x] **Step 10: Normalize and verify the pgx v5 migration**

Run:

```bash
go mod tidy
go fmt ./...
go test ./...
go vet ./...
go build ./...
rg 'jackc/(pgx/v4|pgtype)|context\.TODO' --glob '*.go'
```

Expected: build and checks pass; final `rg` returns no matches.

- [x] **Step 11: Commit the database driver migration**

```bash
git add go.mod go.sum main.go controllers/threads_controller.go repositories db
git commit -m "refactor(db): migrate to pgx v5"
```

---

### Task 4: Replace Deprecated I/O APIs

**Files:**

- Modify: `build/css.go`
- Modify: `files/get_files_in_folder.go`
- Modify: `gateway/connect.go`

**Interfaces:**

- Consumes: existing file paths, CSS rewrite behavior, and gateway registration request.
- Produces: identical behavior through `os` and `io`, with gateway response bodies closed.

- [x] **Step 1: Replace deprecated filesystem calls**

In `build/css.go`, remove `io/ioutil` and use:

```go
files, err := os.ReadDir(dir)
files, err := os.ReadDir("static/styles")
htmlTemplate, err := os.ReadFile(htmlTemplatePath)
err = os.WriteFile(htmlTemplatePath, htmlTemplate, 0644)
```

In `files/get_files_in_folder.go`, remove `io/ioutil` and use:

```go
files, err := os.ReadDir(fullFolderPath)
file, err := os.ReadFile(fullFilePath)
```

Preserve existing return values, ordering, permissions, and error policy.

- [x] **Step 2: Replace gateway body reading and close the response**

In `gateway/connect.go`, replace `io/ioutil` with `io`, then use:

```go
resp, err := client.Do(req)
if err != nil {
	log.Panicln(err)
}
defer resp.Body.Close()

body, err := io.ReadAll(resp.Body)
if err != nil {
	log.Panicln(err)
}
```

- [x] **Step 3: Verify deprecated APIs are gone**

Run:

```bash
go fmt ./...
go test ./...
go vet ./...
go build ./...
rg 'io/ioutil|ioutil\.' --glob '*.go'
```

Expected: Go checks pass; final `rg` returns no matches.

- [x] **Step 4: Commit the compatibility cleanup**

```bash
git add build/css.go files/get_files_in_folder.go gateway/connect.go
git commit -m "refactor: replace deprecated I/O APIs"
```

---

### Task 5: Replace Legacy Compose with Local PostgreSQL

**Files:**

- Modify: `docker-compose.yml`
- Modify: `.env.example`
- Modify: `Makefile`

**Interfaces:**

- Consumes: application `POSTGRES_URL` configuration.
- Produces: PostgreSQL at `localhost:5432`, database/user/password `micrach`, persistent volume `micrach-postgres`, and host-run development commands.

- [x] **Step 1: Replace Compose with the PostgreSQL-only definition**

Replace `docker-compose.yml` with:

```yaml
services:
  db:
    image: postgres:18.4-alpine3.23
    environment:
      POSTGRES_DB: micrach
      POSTGRES_USER: micrach
      POSTGRES_PASSWORD: micrach
    ports:
      - "5432:5432"
    volumes:
      - micrach-postgres:/var/lib/postgresql
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U micrach -d micrach"]
      interval: 2s
      timeout: 5s
      retries: 15
      start_period: 5s

volumes:
  micrach-postgres:
```

The `/var/lib/postgresql` mount intentionally follows the PostgreSQL 18 image
layout.

- [x] **Step 2: Make local environment defaults safe and compatible**

Replace `.env.example` with:

```dotenv
ENV=development
PORT=3000
IS_DB_SEEDED=false
IS_RATE_LIMITER_ENABLED=true
THREADS_MAX_COUNT=50
POSTGRES_URL=postgres://micrach:micrach@localhost:5432/micrach?pool_max_conns=5
THREAD_BUMP_LIMIT=500
IS_CAPTCHA_ACTIVE=true
GATEWAY_URL=
GATEWAY_API_KEY=
GATEWAY_BOARD_ID=
GATEWAY_BOARD_URL=
GATEWAY_BOARD_DESCRIPTION=
```

- [x] **Step 3: Add explicit host and infrastructure commands**

Replace `Makefile` with:

```make
.PHONY: run dev infra-up infra-down

run:
	go run .

dev:
	nodemon --exec go run . --signal SIGTERM

infra-up:
	docker compose up -d db

infra-down:
	docker compose down
```

- [x] **Step 4: Validate Compose syntax and database health**

Run:

```bash
docker compose config
docker compose up -d db
docker compose ps
```

Expected: configuration is valid; `db` transitions to `healthy`.

- [x] **Step 5: Verify clean startup against disposable local data**

Run with a copied local environment:

```bash
if test ! -e .env; then cp .env.example .env; fi
go run .
```

Expected: logs report database and migrations online; application listens on
port 3000. Stop the process after startup verification. An existing `.env` is
never overwritten or removed, and `.env` is never staged.

- [x] **Step 6: Commit local infrastructure**

```bash
git add docker-compose.yml .env.example Makefile
git commit -m "build: add local PostgreSQL compose"
```

---

### Task 6: Build a Multi-Stage Non-Root Image

**Files:**

- Modify: `Dockerfile`
- Create: `.dockerignore`

**Interfaces:**

- Consumes: Go module, templates, static assets, migrations, and writable `/app/uploads`.
- Produces: local `micrach:wave1` Linux image; no registry interaction.

- [x] **Step 1: Replace the Dockerfile**

Replace `Dockerfile` with:

```dockerfile
FROM golang:1.26.5-alpine3.24 AS build

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/micrach .

FROM alpine:3.24.1

RUN apk add --no-cache ca-certificates \
    && addgroup -S micrach \
    && adduser -S -G micrach micrach \
    && mkdir -p /app/uploads \
    && chown -R micrach:micrach /app

WORKDIR /app
COPY --from=build /out/micrach ./micrach
COPY --chown=micrach:micrach templates/ templates/
COPY --chown=micrach:micrach static/ static/
COPY --chown=micrach:micrach migrations/ migrations/

USER micrach
ENTRYPOINT ["./micrach"]
```

- [x] **Step 2: Exclude local state from build context**

Create `.dockerignore`:

```text
.git
.github
.env
.env.*
!.env.example
docs
tests
uploads
micrach
*.out
*.test
```

- [x] **Step 3: Build and inspect the image**

Run:

```bash
docker build --tag micrach:wave1 .
docker image inspect micrach:wave1 --format '{{.Config.User}} {{json .Config.Entrypoint}}'
```

Expected inspect output:

```text
micrach ["./micrach"]
```

- [x] **Step 4: Confirm upload storage is writable by the runtime user**

Run:

```bash
docker run --rm --entrypoint sh micrach:wave1 -c 'touch /app/uploads/write-test && test -f /app/uploads/write-test'
```

Expected: exit 0.

- [x] **Step 5: Commit the image build**

```bash
git add Dockerfile .dockerignore
git commit -m "build: add multi-stage image"
```

---

### Task 7: Replace Deployment Automation with Verification CI

**Files:**

- Modify: `.github/workflows/push.yml`
- Delete: `.github/workflows/codesee-arch-diagram.yml`
- Create: `.github/dependabot.yml`

**Interfaces:**

- Consumes: `go.mod`, `go.sum`, Dockerfile, and repository tests.
- Produces: read-only checks for pushes and pull requests; reviewed dependency update pull requests.

- [x] **Step 1: Replace the workflow with read-only CI**

Replace `.github/workflows/push.yml` with:

```yaml
name: CI

on:
  push:
  pull_request:

permissions:
  contents: read

jobs:
  verify:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout
        uses: actions/checkout@v6

      - name: Set up Go
        uses: actions/setup-go@v6
        with:
          go-version: "1.26.5"
          cache: true
          cache-dependency-path: go.sum

      - name: Verify formatting
        run: |
          files="$(gofmt -l .)"
          test -z "$files" || { printf '%s\n' "$files"; exit 1; }

      - name: Test
        run: go test ./...

      - name: Vet
        run: go vet ./...

      - name: Build
        run: go build ./...

      - name: Build Linux binary
        run: CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /tmp/micrach-linux-amd64 .

      - name: Install govulncheck
        run: go install golang.org/x/vuln/cmd/govulncheck@v1.6.0

      - name: Scan reachable vulnerabilities
        run: govulncheck ./...

      - name: Build container image
        run: docker build --tag micrach:ci .
```

Do not add registry login, image push, artifact upload, SSH, environment secrets,
Swarm, or deployment steps.

- [x] **Step 2: Remove CodeSee**

Delete:

```text
.github/workflows/codesee-arch-diagram.yml
```

- [x] **Step 3: Add Dependabot**

Create `.github/dependabot.yml`:

```yaml
version: 2
updates:
  - package-ecosystem: gomod
    directory: /
    schedule:
      interval: weekly
    open-pull-requests-limit: 5

  - package-ecosystem: github-actions
    directory: /
    schedule:
      interval: weekly
    open-pull-requests-limit: 5
```

- [x] **Step 4: Validate workflow syntax and execute its commands locally**

Run:

```bash
test -z "$(gofmt -l .)"
go test ./...
go vet ./...
go build ./...
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /tmp/micrach-linux-amd64 .
go run golang.org/x/vuln/cmd/govulncheck@v1.6.0 ./...
docker build --tag micrach:ci .
```

Expected: every command exits 0. If `govulncheck` reports a reachable
vulnerability, upgrade or document the exact blocking dependency before this
task can pass; do not suppress findings generically.

- [x] **Step 5: Commit CI and dependency automation**

```bash
git add .github
git commit -m "ci: replace legacy deployment workflow"
```

---

### Task 8: Align Documentation and Agent Guidance

**Files:**

- Modify: `README.md`
- Modify: `AGENTS.md`

**Interfaces:**

- Consumes: completed Wave 1 commands and versions.
- Produces: one accurate setup and verification path for humans and coding agents.

- [x] **Step 1: Update README stack and prerequisites**

State these exact baselines:

```markdown
- Go 1.26
- Fiber v2
- PostgreSQL 18 for local development
```

Remove claims that CI uses Go 1.17 and that Docker requires a prebuilt binary.

- [x] **Step 2: Replace local setup with the host-run workflow**

Document this command sequence:

```sh
cp .env.example .env
docker compose up -d db
go run .
```

Document automatic reload separately:

```sh
make dev
```

State that Compose runs PostgreSQL only and that nodemon is required by
`make dev`.

- [x] **Step 3: Replace container and CI documentation**

Document local image verification:

```sh
docker build -t micrach:local .
```

State that CI checks formatting, tests, vet, native and Linux builds,
reachable vulnerabilities, and Docker image construction. State explicitly
that CI does not publish or deploy images.

- [x] **Step 4: Update AGENTS guidance**

Replace stale Go 1.15/1.17 constraints with Go 1.26. Replace the production
Docker and Swarm description with:

```markdown
- `Dockerfile` builds the application in a pinned Go builder stage and runs it
  as a non-root user in a small runtime image.
- `docker-compose.yml` starts local PostgreSQL only; run micrach on the host.
- CI verifies the project and container build but does not publish or deploy.
```

Keep all existing architecture, upload, migration, testing, and secret-handling
rules that remain accurate.

- [x] **Step 5: Verify documentation**

Run:

```bash
rg -n 'Go 1\.15|Go 1\.17|prebuilt|Swarm|Traefik|docker\.pkg\.github|ghcr\.io' README.md AGENTS.md .github Dockerfile docker-compose.yml
git diff --check
```

Expected: no stale deployment or toolchain matches; `git diff --check` exits 0.

- [x] **Step 6: Commit documentation**

```bash
git add README.md AGENTS.md
git commit -m "docs: document modern development workflow"
```

---

### Task 9: Run the Wave 1 Acceptance Gate

**Files:**

- Verify only; modify a file only to fix a concrete failing check within Wave 1 scope.

**Interfaces:**

- Consumes: all Wave 1 tasks.
- Produces: fresh evidence that Wave 1 is safe to merge and that Wave 2 may be planned separately.

- [x] **Step 1: Verify repository hygiene**

Run:

```bash
git status --short
git diff --check
git ls-files .env uploads micrach '*.out' '*.test'
```

Expected: clean worktree after task commits; no forbidden local artifacts are
tracked.

- [x] **Step 2: Run the full Go gate**

Run:

```bash
go fmt ./...
test -z "$(gofmt -l .)"
go mod tidy
git diff --exit-code -- go.mod go.sum
go test ./...
go vet ./...
go build ./...
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /tmp/micrach-linux-amd64 .
go run golang.org/x/vuln/cmd/govulncheck@v1.6.0 ./...
```

Expected: every command exits 0; module files remain unchanged after tidy;
no reachable vulnerability remains unhandled.

- [x] **Step 3: Run the container and infrastructure gate**

Run:

```bash
docker compose config
docker compose up -d db
docker compose ps
docker build --tag micrach:wave1 .
docker run --rm --entrypoint sh micrach:wave1 -c 'touch /app/uploads/write-test && test -f /app/uploads/write-test'
```

Expected: PostgreSQL is healthy, the image builds, and the runtime user can
write uploads.

- [x] **Step 4: Run manual application smoke tests**

With `.env` copied from `.env.example`, run `go run .` and verify:

```text
GET  /                         catalog renders
GET  /?page=0                  404 renders
POST /                         valid thread redirects
GET  /:threadID                thread renders
POST /:threadID                valid reply redirects with fragment
GET  /captcha/:captchaID       PNG renders when CAPTCHA is enabled
```

Also verify JPEG and PNG upload limits, thumbnail creation, pagination, rate
limiting, `sage`, bump limit, and archival behavior. Inspect catalog and thread
pages at narrow and wide viewport widths; Wave 1 must introduce no visual
change.

- [x] **Step 5: Remove disposable local state**

Run:

```bash
docker compose down
rm -f /tmp/micrach-linux-amd64
git status --short --branch
```

Expected: no generated or local-only file appears in Git status. Leave any
pre-existing or locally created ignored `.env` untouched.

## Wave 2 Entry Gate

Do not implement Fiber v3 in this plan. Start a separate spec-derived plan only
after Task 9 passes. That follow-up must recheck the current Fiber v3, template
engine, and Go requirements through Context7 and official release notes, then
cover handler context signatures, static middleware, form parsing, redirects,
CAPTCHA, uploads, pagination, rate limiting, `sage`, archival, and responsive
smoke tests. Architecture cleanup remains excluded.
