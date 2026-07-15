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
