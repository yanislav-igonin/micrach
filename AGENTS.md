# Repository Guidelines

## Project Structure & Module Organization

`main.go` wires the Fiber application, middleware, routes, database startup, and template engine. Keep HTTP handlers in `controllers/`, persistence logic and data types in `repositories/`, database setup in `db/`, and environment parsing in `config/`. Gateway integration belongs in `gateway/`; shared helpers live in `utils/` and `files/`. Server-rendered HTML is organized under `templates/{pages,components,head}`, while browser assets live in `static/`. Add ordered PostgreSQL changes to `migrations/` using the existing `<number>-<description>.sql` pattern. Load-testing scripts belong in `tests/`.

## Build, Test, and Development Commands

- `cp .env.example .env` creates local configuration; adjust `POSTGRES_URL` before starting.
- `go run .` starts the service and applies migrations against the configured PostgreSQL database.
- `make dev` runs the app through `nodemon` for automatic restarts.
- `go build .` produces the local `micrach` binary.
- `GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build .` matches the CI Linux build.
- `go test ./...` runs all Go tests.
- `go vet ./...` performs standard static checks.
- `sh tests/load.sh` runs the optional Vegeta load test; review its target URL first.

## Coding Style & Naming Conventions

Use `gofmt` (tabs for Go indentation) and run `go fmt ./...` before submitting. Keep package names short and lowercase. Use `PascalCase` for exported identifiers, `camelCase` for unexported identifiers, and descriptive handler names such as `CreateThread`. Keep controllers thin: validation and HTTP concerns stay in handlers, while SQL and persistence behavior stay in repositories. Follow existing kebab-case names for CSS, templates, and migration descriptions.

## Testing Guidelines

The repository currently has no unit tests or enforced coverage threshold. Add tests beside the code as `*_test.go`, using `TestXxx` and table-driven cases where practical. Avoid shared production databases; isolate repository tests with disposable test data or a dedicated database. Run `go test ./...` and `go vet ./...` for every change. Use Vegeta only for explicit performance work.

## Commit & Pull Request Guidelines

History mixes terse imperative subjects with `feat:` and `fix:` prefixes. Prefer focused Conventional Commit subjects, for example `fix: reject oversized uploads`. Keep each commit buildable. Pull requests should explain behavior changes, configuration or migration impact, and verification commands. Link relevant issues and include screenshots for template or CSS changes. Never commit `.env`, uploads, generated load-test artifacts, or credentials.
