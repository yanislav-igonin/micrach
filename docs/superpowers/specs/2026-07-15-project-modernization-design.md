# Project Modernization Design

Date: 2026-07-15

## Goal

Modernize micrach in controlled stages without changing application behavior or
reworking its architecture. First establish a supported toolchain, current
dependencies, reproducible local infrastructure, useful CI, and a minimal test
safety net. Migrate from Fiber v2 to Fiber v3 only after that baseline is
stable.

## Current State

- `go.mod` declares Go 1.15, while CI uses Go 1.17.
- Fiber is v2.31.0 and pgx is v4.15.0.
- The Dockerfile packages a prebuilt binary in BusyBox; it does not build the
  application.
- The GitHub Actions workflow builds, publishes an image, and deploys it to an
  obsolete Docker Swarm environment over SSH.
- `docker-compose.yml` describes that obsolete Swarm, Traefik, and VPS setup.
- The CodeSee workflow and its action versions are obsolete.
- No checked-in Go tests provide a regression safety net.
- Parts of the code use deprecated `ioutil` APIs, broad `context.TODO()` calls,
  and unchecked transaction completion errors.

## Scope and Constraints

- Preserve routes, templates, HTTP status behavior, database schema, upload
  layout, and minimal-JavaScript approach.
- Preserve the current package boundaries and application architecture.
- Do not combine dependency work with handler or repository redesign.
- Do not add infrastructure that the application does not need. PostgreSQL is
  the only Docker Compose service required now.
- Run the application directly on the host with `go run .` or `make dev`.
- Keep a multi-stage Dockerfile as a reproducible application build artifact,
  but do not publish or deploy its image in CI.
- Do not add or modify database migrations unless an upgrade proves that a
  compatibility migration is unavoidable. Such a change requires separate
  review before implementation.

## Strategy

Use two modernization waves. Wave 1 stays on the latest compatible Fiber v2 to
avoid mixing framework API migration with the toolchain, database driver, CI,
container, and test changes. Wave 2 performs the Fiber v3 migration as a
separate change after Wave 1 is green.

Exact dependency versions must be rechecked against Context7 and official
release notes when each wave begins. Major upgrades remain explicit; automated
dependency tooling must not silently merge them.

## Wave 1: Stable Modern Baseline

### Go and Dependencies

- Set the supported Go line to Go 1.26. Use the latest available Go 1.26 patch
  release in CI and the Docker builder image; keep documentation aligned.
- Upgrade Fiber from v2.31.0 to the latest stable Fiber v2 release.
- Upgrade the Fiber HTML template engine to the latest Fiber v2-compatible
  release.
- Migrate `github.com/jackc/pgx/v4` and standalone
  `github.com/jackc/pgtype` usage to `github.com/jackc/pgx/v5` packages.
- Replace `pgtype.Int4Array` with the pgx v5-compatible representation selected
  during implementation, preserving query results and ordering.
- Upgrade other direct dependencies to current versions compatible with the
  selected Go and Fiber versions.
- Run `go mod tidy` and review every direct and indirect dependency change.
- Upgrade dependencies in coherent groups so failures can be attributed to a
  specific change.

### Compatibility Cleanup

- Replace deprecated `io/ioutil` calls with `io` and `os` equivalents.
- Replace broad `context.TODO()` usage where a request, startup, or explicit
  operation context is available.
- Check and propagate transaction `Commit` errors. Preserve safe rollback
  behavior.
- Limit changes to those needed for supported APIs, correct error handling, and
  testability. Do not refactor package ownership or business flow.

### Local Development Infrastructure

- Replace the Swarm-oriented Compose file with a local PostgreSQL service.
- Give PostgreSQL a persistent named volume and a healthcheck.
- Use development-only credentials and defaults consistent with
  `.env.example`.
- Do not add the application, Redis, a proxy, or other optional services to
  Compose.
- Document the normal workflow: start PostgreSQL with Compose, then run micrach
  on the host with `go run .` or `make dev` using nodemon.

### Dockerfile

- Replace the prebuilt-binary BusyBox image with a multi-stage Dockerfile.
- Build a static Linux binary in a pinned Go 1.26 builder image.
- Copy only the binary, templates, static assets, and migrations into a small
  runtime image.
- Run the application as a non-root user and preserve writable upload storage.
- Do not add a container healthcheck until the application has a suitable
  health endpoint. Adding such an endpoint is outside this modernization.

### CI and Automation

- Replace the existing build, publish, and deploy workflow with a read-only CI
  workflow.
- Run formatting verification, `go test ./...`, `go vet ./...`, `go build .`,
  the Linux static build, and `govulncheck ./...`.
- Build the Docker image as verification, without login, push, registry tags,
  or deployment.
- Use current stable GitHub-maintained action releases and Go module caching.
- Remove SSH deployment, Docker Swarm service update, registry publishing, and
  their workflow configuration.
- Remove the obsolete CodeSee workflow.
- Add Dependabot configuration for Go modules and GitHub Actions. Updates open
  reviewable pull requests and are not auto-merged.

### Test Safety Net

Add focused tests beside changed code, favoring table-driven cases where
practical:

- configuration parsing and defaults;
- post and upload validation boundaries;
- pure migration, filesystem, and repository helpers where they can be tested
  without a shared database;
- Fiber handler smoke tests for key status and rendering paths that do not
  require production data;
- pgx v5 array conversion or scanning behavior introduced by the migration.

Tests requiring PostgreSQL must use disposable data and the local development
database only when explicitly invoked. Routine unit tests must not depend on a
developer's shared database.

### Documentation

- Update README prerequisites, setup, development commands, container notes,
  and supported versions.
- Make `.env.example` safe for local development: no implicit external gateway
  dependency and no repeated seed behavior by default.
- Update AGENTS guidance when its Go, CI, Docker, or Compose statements become
  stale.

## Wave 2: Fiber v3 Migration

Wave 2 starts only after Wave 1 passes all automated and manual checks.

- Recheck the current stable Fiber v3, template engine, and minimum Go
  requirements through Context7 and official release notes.
- Keep Go 1.26 if supported; otherwise upgrade to the required supported Go
  line as an explicit part of Wave 2.
- Change Fiber imports from v2 to v3 and use a v3-compatible template engine.
- Adapt handler context signatures, middleware configuration, static serving,
  form and body parsing, redirects, and other changed APIs.
- Preserve routes, response codes, template data, forms, and browser behavior.
- Exercise catalog, thread view, thread creation, replies, CAPTCHA, pagination,
  rate limiting, `sage`, archival, JPEG/PNG uploads, and thumbnail generation.
- Fix only migration regressions. Do not include architectural cleanup.
- Refresh compatible dependencies again after Fiber v3, run `go mod tidy`, and
  repeat vulnerability scanning.
- Record any remaining major-version candidates or deprecated APIs for later
  work instead of expanding this migration.

## Ongoing Maintenance

- Keep the project on a supported Go release line and consume patch releases
  promptly through reviewed dependency automation.
- Review Dependabot pull requests for Go modules and GitHub Actions.
- Perform periodic dependency refreshes rather than another multi-year jump.
- Isolate future major upgrades in their own pull requests and do not mix them
  with unrelated refactoring.

## Verification and Acceptance

Each wave must provide fresh evidence for all applicable checks:

1. Formatting produces no diff.
2. `go test ./...` passes.
3. `go vet ./...` passes.
4. `go build .` passes.
5. The Linux static build passes.
6. `govulncheck ./...` reports no unhandled reachable vulnerability.
7. `go mod tidy` produces no unexplained diff.
8. `git diff --check` passes.
9. The Docker image builds successfully.
10. Compose PostgreSQL reaches healthy state.
11. The host-run application connects to Compose PostgreSQL and applies
    migrations to a clean database.
12. Manual smoke tests cover affected pages, forms, responsive states, CAPTCHA,
    pagination, rate limiting, uploads, thumbnails, `sage`, and archival.
13. No `.env`, upload, binary, credential, or generated artifact is tracked.

Wave 1 is complete only when its checks pass and existing behavior remains
unchanged. Wave 2 cannot begin before that point. A failed major upgrade is
rolled back as its own isolated change; it must not destabilize the Wave 1
baseline.

## Explicit Non-Goals

- Handler or repository architecture cleanup.
- New routes, UI redesign, or feature work.
- Database schema redesign.
- New caching, queueing, proxy, or observability services.
- Containerized local execution of the application.
- Image publishing, GitHub Container Registry integration, VPS deployment,
  Docker Swarm, or Traefik configuration.
