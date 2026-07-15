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

- Go (the module declares Go 1.15; CI builds with Go 1.17)
- [Fiber v2](https://github.com/gofiber/fiber)
- PostgreSQL via `pgx`
- Plain HTML templates and CSS

## Run locally

### Prerequisites

- Go 1.17 or newer
- PostgreSQL
- A PostgreSQL database that the configured user can create tables in

### Setup

1. Create the database. With standard local PostgreSQL tooling:

   ```sh
   createdb micrach
   ```

2. Copy the example configuration:

   ```sh
   cp .env.example .env
   ```

3. Edit `.env`. At minimum, make `POSTGRES_URL` match your local database.
   For normal development, set `IS_DB_SEEDED=false`; when true, sample rows are
   inserted every time the application starts.

4. Start the application from the repository root:

   ```sh
   go run .
   ```

5. Open [http://localhost:3000](http://localhost:3000), or the port configured
   by `PORT`.

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
# Run with automatic restart; requires nodemon.
make dev

# Format, test, vet, and build.
go fmt ./...
go test ./...
go vet ./...
go build .

# Reproduce the Linux CI build.
GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build .
```

`tests/load.sh` is an optional Vegeta load test. Its checked-in target is a
historical public URL, so change it to an environment you own before running
the script. Generated load-test files are ignored by Git.

## Containers and deployment

The `Dockerfile` packages a prebuilt `micrach` binary together with templates,
static assets, and migrations. Build the Linux binary before building the
image. The checked-in `docker-compose.yml` is production-oriented: it expects
Docker Swarm, an existing `web` overlay network, Traefik, a published image,
and host-mounted upload storage. It is not a standalone local PostgreSQL setup.

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
