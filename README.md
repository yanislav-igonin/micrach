# micrach

micrach is a small, server-rendered imageboard written in Go. It provides one
board where visitors can create threads, reply with text or images, bump or
`sage` threads, and browse the thread catalog without a client-side SPA.

The project began as a Go learning project and as an experiment in building a
lightweight imageboard with very little JavaScript. The former public instance
is offline because it required more moderation than the maintainer could
provide; the source remains available for local use and development.

## Features

- Server-rendered HTML using Fiber templates
- Thread creation, replies, catalog pagination, bump limits, and `sage`
- JPEG and PNG uploads with generated thumbnails
- Optional CAPTCHA and request rate limiting
- Automatic PostgreSQL migrations on startup
- Automatic archival when the active-thread limit is reached
- Optional registration with an external board gateway

## Stack

- Go 1.26
- Fiber v2
- PostgreSQL 18 for local development
- PostgreSQL access through [`pgx`](https://github.com/jackc/pgx)
- Plain HTML templates and CSS

## Run locally

### Prerequisites

- Go 1.26
- Docker with Docker Compose
- `nodemon` when using automatic reload

### Setup

From the repository root:

```sh
cp .env.example .env
docker compose up -d db
go run .
```

The example configuration matches the PostgreSQL 18 service started by
Compose. Compose runs PostgreSQL only; the micrach application runs on the
host. Edit `.env` when you need different local settings. For normal
development, keep `IS_DB_SEEDED=false`; when true, sample rows are inserted
every time the application starts.

For automatic reload, install `nodemon`, start PostgreSQL as shown above, and
run:

```sh
make dev
```

Open [http://localhost:3000](http://localhost:3000), or the port configured by
`PORT`.

Startup connects to PostgreSQL, applies missing files from `migrations/`, and
creates the ignored local `uploads/` directory. Uploaded originals and
thumbnails are stored on disk, so keep that directory persistent when running
the service in a container.

## Configuration

`.env` is loaded automatically and is ignored by Git.

| Variable | Default in code | Purpose |
| --- | --- | --- |
| `ENV` | `release` | Runtime mode. `production` enables CSS filename rewriting at startup. |
| `PORT` | `3000` | HTTP listen port. |
| `POSTGRES_URL` | `postgresql://localhost/micrach` | PostgreSQL connection URL. |
| `IS_DB_SEEDED` | `false` | Insert sample posts at startup. Avoid enabling repeatedly. |
| `IS_RATE_LIMITER_ENABLED` | `true` | Enable Fiber request limiting. |
| `THREADS_MAX_COUNT` | `50` | Maximum active threads before old threads are archived. |
| `THREAD_BUMP_LIMIT` | `500` | Reply count after which a thread no longer bumps. |
| `IS_CAPTCHA_ACTIVE` | `true` | Require CAPTCHA for new threads and replies. |
| `GATEWAY_URL` | empty | External gateway base URL. Empty disables gateway integration. |
| `GATEWAY_API_KEY` | empty | Shared key used for gateway registration and `/api/ping`. |
| `GATEWAY_BOARD_ID` | empty | Board identifier sent to the gateway. |
| `GATEWAY_BOARD_URL` | empty | Parsed from the environment but not currently used by the gateway client. |
| `GATEWAY_BOARD_DESCRIPTION` | empty | Board description sent to the gateway. |

See [.env.example](.env.example) for a complete example.

## Development commands

```sh
# Format, test, vet, and build.
go fmt ./...
go test ./...
go vet ./...
go build .

# Reproduce the Linux CI build.
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /tmp/micrach-linux-amd64 .
```

`tests/load.sh` is an optional Vegeta load test. Its checked-in target is a
historical public URL, so change it to an environment you own before running
the script. Generated load-test files are ignored by Git.

## Containers and CI

The multi-stage `Dockerfile` builds micrach in a pinned Go builder image and
runs it as a non-root user in a small runtime image. Verify the image locally:

```sh
docker build -t micrach:local .
```

CI checks formatting, tests, vet, native and Linux builds, reachable
vulnerabilities, and Docker image construction. It has read-only repository
permissions and does not publish or deploy images.

## Project layout

```text
controllers/   HTTP handlers and request validation flow
repositories/  PostgreSQL queries and view data types
db/            connection pool and migration runner
config/        environment parsing
gateway/       optional external gateway integration
templates/     server-rendered pages, components, and head fragments
static/        CSS, icons, and error images
uploads/       runtime uploads (created locally, ignored by Git)
migrations/    ordered PostgreSQL schema changes
tests/         optional load-test scripts
```

## Screenshot

![micrach catalog](https://user-images.githubusercontent.com/15846431/237027105-481f2d55-541b-4c3b-8dd0-db464c25102f.jpeg)

## License

[MIT](LICENSE)
