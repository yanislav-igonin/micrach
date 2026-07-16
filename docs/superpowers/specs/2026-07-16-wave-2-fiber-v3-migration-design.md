# Wave 2 Fiber v3 Migration Design

Date: 2026-07-16

## Goal

Migrate micrach from Fiber v2 to the current stable Fiber v3 and refresh the
compatible dependency baseline without changing application behavior,
architecture, routes, templates, persistence, uploads, or browser flow.

## Entry Baseline

Wave 1 is merged and green. The starting baseline is:

- Go `1.26.0` with toolchain `go1.26.5`;
- Fiber `v2.52.14`;
- Fiber HTML template engine `v2.1.3`;
- pgx `v5.10.0`;
- passing `go test ./...`, `go vet ./...`, and `go build ./...`;
- local PostgreSQL through Docker Compose and host-run application workflow.

At design time, the current stable migration targets are Fiber `v3.4.0` and
Fiber HTML `v3.0.6`. Fiber v3 requires Go 1.25 or later, so Go remains on 1.26.
These are a dated snapshot, not permanent pins. Implementation must recheck
Context7, official release notes, and `go list -m` before changing dependencies.

## Scope

Wave 2 is one focused pull request containing:

- Fiber v2 to Fiber v3 migration;
- Fiber v3-compatible HTML template engine migration;
- required handler, context, middleware, static, redirect, and test API changes;
- DB-free unit characterization for migration-sensitive behavior;
- compatible patch and minor dependency refresh after Fiber migration;
- documentation updates and a confirmed modernization follow-up list;
- automated verification and manual application smoke testing.

## Constraints

- Preserve all existing routes, methods, status codes, templates, view data,
  form field names, static URLs, upload paths, response flow, and minimal-JS UI.
- Preserve current package boundaries and server-rendered architecture.
- Do not redesign handlers or repositories, add dependency injection, introduce
  an application factory, or add generalized compatibility wrappers.
- Do not add database migrations, schema changes, new routes, features, or UI
  changes.
- Do not fix pre-existing bugs unless Fiber v3 causes a demonstrated regression.
- Do not add tests that require PostgreSQL, Docker, network services, or other
  external infrastructure.
- Do not add an integration or end-to-end test harness.
- Do not combine other major dependency upgrades with Fiber v3.
- Never commit `.env`, uploads, binaries, credentials, or generated artifacts.

## Chosen Migration Strategy

Use a manual, narrow migration. Do not run `fiber migrate` in the working tree.
The migration CLI may be run against a disposable copy for comparison only;
its output is not authoritative and must not be copied without review.

Implementation order:

1. Add or refine DB-free characterization tests for migration-sensitive APIs.
2. Change Fiber and template dependencies and adapt compile-time API changes.
3. Make only the smallest changes needed to preserve runtime behavior.
4. Refresh compatible dependencies and normalize the module graph.
5. Update documentation and record confirmed follow-ups.
6. Run the full automated gate and manual smoke matrix.

This approach is preferred over CLI-driven rewriting because the Fiber surface
is small and status/static behavior requires deliberate review. Compatibility
wrappers are rejected because they would add lasting structure for a one-time
migration.

## Dependency Selection

### Go and Fiber

Keep the module, CI, and Docker builder on Go 1.26 because Fiber v3's current
minimum is Go 1.25. Upgrade them only if the rechecked stable Fiber release has
a higher requirement.

Replace `github.com/gofiber/fiber/v2` with the current stable
`github.com/gofiber/fiber/v3`. At implementation start, verify that the selected
release is stable, not a beta or release candidate.

### Template Engine

The primary target is the current stable
`github.com/gofiber/template/html/v3`. Before committing the dependency graph,
compile the application and load/render the real templates through Fiber v3.

If HTML v3 has a demonstrated incompatibility with the selected Fiber v3,
retain the documented Fiber v3-compatible `html/v2` engine as a controlled
fallback. Record the exact incompatibility and repeat the complete verification
gate. Documentation inconsistency alone is not enough to trigger fallback.

### Other Dependencies

After Fiber v3 is working:

- run `go list -u -m` and review direct dependencies individually;
- accept only compatible patch or minor releases;
- do not add unrelated direct dependencies;
- do not take another major upgrade;
- run `go mod tidy`, inspect `go.mod`, `go.sum`, and `go mod graph`;
- remove obsolete direct Fiber v2 and HTML v2 requirements when HTML v3 is used;
- tolerate a lower-major transitive module only when the selected upstream
  package legitimately requires it.

At design time, the non-Fiber direct dependencies have no available updates.
Implementation must still recheck them rather than assuming that remains true.

## Fiber API Migration

### Handler Contexts

Change all handlers and middleware callbacks from:

```go
func(c *fiber.Ctx) error
```

to:

```go
func(c fiber.Ctx) error
```

Fiber v3 `fiber.Ctx` implements `context.Context`. Remove `c.UserContext()` and
pass `c` directly to repository and pgx methods that accept `context.Context`.
Do not retain a second request context variable unless a specific API requires
one.

### Route Parameters

Fiber v3 removes `ParamsInt`. Do not replace it with
`fiber.Params[int](c, "threadID", 0)`, because an invalid parameter would become
zero and reach the repository instead of returning the current early 404.

Parse with `strconv.Atoi(c.Params("threadID"))` and preserve the existing 404
rendering on conversion failure. Catalog query parsing already uses
`strconv.Atoi` and remains unchanged.

### Redirects

Fiber v3 changes the default redirect from 302 to 303 and replaces the old
redirect call with a builder. Both successful write flows must set 302
explicitly:

```go
return c.Redirect().Status(fiber.StatusFound).To(path)
```

Thread creation continues redirecting to `/:threadID`. Reply creation continues
redirecting to `/:threadID#<postID>`.

### Static Files

Fiber v3 removes `app.Static`. Register the existing public prefixes through
`github.com/gofiber/fiber/v3/middleware/static`:

- `/static/*` serves `./static`;
- `/uploads/*` serves `./uploads`.

Preserve prefix stripping, file bodies, status behavior, dynamic route
precedence, and directory-browsing defaults. Unit tests must validate the exact
route form selected during implementation before the old registrations are
removed.

### Middleware

Move recover, limiter, and compress imports to Fiber v3 while preserving their
registration order:

1. recover;
2. optional limiter;
3. compress;
4. static and application routes as currently ordered.

Change limiter `Next` to `func(c fiber.Ctx) bool`. Preserve:

- `Max: 50`;
- local-request bypass;
- `/static`, `/uploads`, and `/captcha` bypass;
- current default expiration, keying, headers, and limit response.

Small unexported helpers are allowed for exact DB-free testing of static route
registration and the limiter skip predicate. Do not turn these helpers into a
new application-construction layer.

### Multipart Forms and Uploads

Keep `c.MultipartForm()` and `c.SaveFile()`; both remain available in Fiber v3.
Do not migrate forms to the Binder API without a compile-time requirement.

Fiber v3 changed multipart body-limit enforcement. Compare observed v2 and v3
behavior during manual smoke testing. Do not change the configured/effective
limit unless a migration regression is demonstrated. If needed, add the
smallest explicit `BodyLimit` setting that preserves the verified v2 contract,
and document the evidence.

### Templates and Tests

Keep template root, extension, helper names, template names, bindings, and
render calls unchanged. Update only engine imports or signatures required by
the selected template release.

`app.Test(req)` remains valid without an explicit test configuration. Do not add
an empty `fiber.TestConfig`, because its semantics differ from Fiber's omitted
default configuration.

## DB-Free Unit Test Safety Net

All routine tests must run with no PostgreSQL, Docker, network service, or
external process.

Update existing controller and gateway tests for Fiber v3 and cover:

- invalid catalog page still renders 404;
- invalid thread ID still renders 404 before repository access;
- invalid multipart thread input returns 400 before repository access;
- gateway ping authentication and status behavior remain unchanged.

Add focused tests where practical:

- load the real templates with the configured helper functions;
- serve files for `/static/*` and `/uploads/*` from temporary directories;
- prove static prefixes map to the intended filesystem paths;
- prove a missing static file does not mask application routes;
- prove the limiter reaches 429 for ordinary remote requests;
- prove local, static, uploads, and CAPTCHA requests bypass the limiter.

Tests may use `fiber.App`, `httptest`, `t.TempDir`, and in-memory middleware
state. Do not mock repositories or build an integration test seam merely to
automate successful database-backed flows.

## Error and Regression Handling

- Preserve existing 400, 404, 429, and 500 behavior and error templates.
- Treat compilation failures, changed statuses, changed route precedence,
  broken template rendering, static mapping changes, and changed browser flow
  as migration regressions.
- Reproduce a suspected regression with a focused unit check or manual request
  before changing behavior.
- Apply the smallest Fiber-compatibility fix and rerun the full gate.
- Record unrelated defects and deferred redesigns as follow-ups.
- If the stable Fiber v3 migration cannot pass the acceptance gate, revert the
  isolated migration rather than weakening the Wave 1 baseline.

## Documentation and Follow-Ups

Update README and AGENTS guidance to state Go 1.26 and Fiber v3, while keeping
local development, CI, Docker, and PostgreSQL instructions aligned with the
actual implementation.

Create `docs/superpowers/modernization-followups.md`. Include only candidates
confirmed after the final module, vulnerability, and API audit. Each entry must
name the package or API, current and candidate state, why it is deferred, and
the evidence source. Do not add speculative cleanup ideas.

The pull request must report selected versions, compatibility decisions,
dependency/fallback results, commands actually run, manual smoke results, and
any remaining follow-ups.

## Automated Acceptance Gate

Run fresh commands after the final dependency graph is settled:

```sh
go fmt ./...
test -z "$(gofmt -l .)"
go mod tidy
git diff --exit-code -- go.mod go.sum
go test ./...
go vet ./...
go build ./...
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /tmp/micrach-linux-amd64 .
go run golang.org/x/vuln/cmd/govulncheck@v1.6.0 ./...
git diff --check
docker build --tag micrach:wave2 .
```

Requirements:

- every command exits zero;
- `go mod tidy` leaves no unexplained diff;
- no reachable vulnerability remains unhandled;
- the Docker build uses the selected Go line and succeeds;
- no forbidden generated or local artifact is tracked.

No command in `go test ./...` may require PostgreSQL or another external
service.

## Manual Smoke Matrix

Use the disposable local Compose PostgreSQL and run the application on the
host. This is an operational acceptance pass, not an automated test suite.

Verify:

- catalog renders, including empty and populated states;
- `/?page=0` and out-of-range pages render 404;
- thread page and missing/invalid thread IDs behave unchanged;
- valid thread creation returns 302 and its exact location;
- valid reply returns 302 with the post fragment;
- CAPTCHA image is PNG and valid/invalid answers retain their behavior;
- JPEG and PNG uploads honor verified limits;
- multiple uploads, originals, and thumbnails use unchanged paths;
- pagination links and page boundaries work;
- limiter returns 429 and still bypasses local/static/uploads/CAPTCHA requests;
- `sage` does not bump the thread;
- normal reply bump, bump limit, and archival behavior remain unchanged;
- `/static` and `/uploads` assets render correctly;
- catalog and thread pages have no narrow or wide viewport regression.

Record actual results in the pull request. Do not replace failures with an
integration harness inside this wave.

## Non-Goals

- Handler or repository architecture cleanup.
- Repository mocks, service interfaces, application factories, or DI.
- Automated PostgreSQL, upload-flow, or browser E2E tests.
- New features, routes, templates, CSS, JavaScript, or UI redesign.
- Database schema or migration changes.
- Deployment, publishing, proxy, cache, queue, or observability changes.
- Unrelated dependency major upgrades.

## References

- [Fiber v3 migration guide](https://docs.gofiber.io/next/whats_new/)
- [Fiber v3.4.0 release](https://github.com/gofiber/fiber/releases/tag/v3.4.0)
- [Fiber static middleware](https://docs.gofiber.io/next/middleware/static/)
- [Fiber HTML v3.0.6 release](https://github.com/gofiber/template/releases/tag/html%2Fv3.0.6)
