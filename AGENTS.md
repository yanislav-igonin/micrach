# AGENTS.md

Guidance for coding agents working in this repository. Keep changes focused,
preserve existing behavior unless the task explicitly changes it, and verify
all claims with fresh commands.

## Project Overview

micrach is a small, server-rendered, single-board imageboard. It uses Go,
Fiber v2, Fiber HTML templates, PostgreSQL through `pgx`, and filesystem-backed
image uploads. It intentionally uses little browser-side JavaScript.

The service supports thread creation, replies, JPEG/PNG attachments,
thumbnails, CAPTCHA, rate limiting, pagination, bump limits, `sage`, thread
archival, and optional external gateway registration.

## Repository Map

- `main.go`: application startup, middleware, static mounts, routes, and server.
- `controllers/`: Fiber handlers, request parsing, HTTP responses, and form flow.
- `repositories/`: SQL queries, persistence types, and template view data.
- `db/`: PostgreSQL pool initialization and startup migration runner.
- `config/`: environment parsing and defaults.
- `gateway/`: optional gateway registration and authenticated ping endpoint.
- `utils/`: post validation, upload directories, and image thumbnail handling.
- `files/`: shared filesystem read helpers.
- `build/`: production CSS cache-busting rewrite.
- `templates/{pages,components,head}`: server-rendered HTML.
- `static/`: checked-in CSS, icons, and error images.
- `uploads/`: generated runtime files; ignored by Git and never commit them.
- `migrations/`: ordered PostgreSQL changes named `<number>-<description>.sql`.
- `tests/`: optional Vegeta load test, not the unit-test suite.

## Runtime Flow and Routes

Startup loads `.env`, initializes configuration and the DB pool, applies
missing migrations, optionally seeds data, performs production CSS rewriting,
creates `uploads/`, then starts Fiber.

| Method | Path | Handler / purpose |
| --- | --- | --- |
| `GET` | `/` | Render paginated thread catalog. |
| `POST` | `/` | Validate and create a thread. |
| `GET` | `/:threadID` | Render one thread and its replies. |
| `POST` | `/:threadID` | Validate and add a reply. |
| `GET` | `/captcha/:captchaID` | Render CAPTCHA PNG. |
| `GET` | `/api/ping` | Gateway health endpoint; registered only when `GATEWAY_URL` is set. |

Static content is served from `/static`; generated files are served from
`/uploads`. Thread and reply creation spans PostgreSQL plus filesystem writes,
so review both paths when changing upload behavior.

## Local Setup and Commands

Use Go 1.26. From the repository root, copy the development configuration,
start the PostgreSQL 18 service, and run micrach on the host:

```sh
cp .env.example .env
docker compose up -d db
go run .
```

The example `POSTGRES_URL` matches the Compose service. `make dev` provides
automatic reload and requires `nodemon`.

Common checks:

```sh
go fmt ./...
go test ./...
go vet ./...
go build .
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /tmp/micrach-linux-amd64 .
```

`sh tests/load.sh` requires Vegeta and targets a historical public URL by
default; change the target to a controlled environment before use. Do not run
the load test as routine verification.

## Configuration and Generated State

- `.env` is local and must never be committed.
- `POSTGRES_URL` selects the database; startup applies migrations automatically.
- `IS_DB_SEEDED=true` adds sample data on every startup, not only once.
- `ENV=production` renames CSS files and rewrites templates in place at startup.
  Do not use this mode in a development checkout unless that mutation is wanted.
- `GATEWAY_URL=""` disables gateway routes and registration.
- Uploads are limited to four JPEG/PNG files, 3 MiB each. Titles allow 100
  runes; post text allows 1,000 runes.
- Runtime uploads, built binaries, `.env`, and load-test artifacts are ignored
  and must remain untracked.

## Architecture and Coding Conventions

- Keep handlers thin: HTTP parsing/status/rendering in `controllers/`; SQL and
  persistence behavior in `repositories/`.
- Keep DB startup and migration mechanics in `db/`, environment parsing in
  `config/`, gateway behavior in `gateway/`, and reusable helpers in `utils/`
  or `files/`.
- Follow existing package naming: short, lowercase names. Use `PascalCase` for
  exported identifiers and `camelCase` for unexported identifiers.
- Format Go with `gofmt`/`go fmt`; do not hand-format with spaces.
- Follow existing kebab-case naming for templates, CSS, and migration
  descriptions. Preserve current server-rendered approach and minimal-JS goal.
- Use parameterized SQL. Keep related DB operations transactional where the
  existing workflow requires atomicity.
- Avoid unrelated refactors. If a task exposes adjacent problems, report them
  separately unless they block the requested work.

## Change Guide

- Route or request behavior: update `main.go` and/or `controllers/`, then cover
  status codes, redirects, validation errors, and template data.
- Persistence behavior: update `repositories/`; add a migration for schema
  changes rather than editing an already-applied migration.
- Schema changes: add the next numeric file in `migrations/`. Keep filename
  shape exactly `<number>-<description>.sql`; the runner parses it directly.
- Page or component changes: edit `templates/` and matching `static/styles/`.
  Check catalog, thread, error, empty, and pagination states as relevant.
- Upload changes: inspect controller transaction flow, `utils/`, repository file
  records, original/thumbnail paths, validation limits, and cleanup behavior.
- Config changes: update `config/config.go`, `.env.example`, and the README
  config table together.
- Deployment changes:
  - `Dockerfile` builds the application in a pinned Go builder stage and runs it
    as a non-root user in a small runtime image.
  - `docker-compose.yml` starts local PostgreSQL only; run micrach on the host.
  - CI verifies the project and container build but does not publish or deploy.

## Testing Expectations

Go unit tests are checked in across the application packages; there is no
coverage threshold. Add `*_test.go` files beside changed code, use `TestXxx`,
and prefer table-driven cases when practical. Repository tests must use
disposable test data or a dedicated database, never a shared production
database.

Minimum verification for code changes:

1. Run `go fmt ./...` and review the resulting diff.
2. Run `go test ./...`.
3. Run `go vet ./...`.
4. Run `go build .` (or the Linux CI build for deployment-sensitive work).
5. For templates/CSS, start the app and manually inspect affected responsive
   states; include screenshots in the PR.

For documentation-only changes, inspect the rendered Markdown, verify all local
links and commands against the repository, and review `git diff --check`.

## Git, Commits, and Pull Requests

- Inspect `git status` before editing. Preserve user changes and unrelated files.
- Never commit `.env`, credentials, uploads, binaries, or generated load-test
  outputs.
- Prefer focused Conventional Commit subjects, such as
  `fix: reject oversized uploads` or `docs: expand setup guide`.
- Keep each commit buildable. Do not rewrite history or discard local changes
  unless explicitly requested.
- PRs should describe behavior, config/migration impact, verification commands,
  linked issues, and screenshots for UI changes.

## Agent Completion Checklist

- Requested scope is fully addressed; unrelated behavior is untouched.
- New behavior has tests where practical.
- Documentation and `.env.example` match config/code changes.
- No secret, upload, binary, or generated artifact is included.
- Fresh verification commands pass, and final report names commands actually run.
