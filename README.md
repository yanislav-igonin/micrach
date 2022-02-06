# micrach
Go tryout. Single board imageboard.

## Motivation
Tired of a fucking SPAs, so I decided to write this simple and lightweight imageboard (microboard, actually) with as less JS as possible.  
Also, I wanted to learn Go, so I decided to write this in Go.

Maybe I'll add an api and create a SPA with the Next.js framework in the future, but right now I'm trying to make this as simple as possible.

## Prerequisites
1. Go 1.13+.
2. PostgreSQL.
3. Create `.env` file from `.env.example`, change env vars to your needs.
4. Create db with the name that is specified in `.env` file in `POSTGRES_URL` env var.

## Run
Just run:
```sh
go run main.go
```

**In development** I prefer to run it with [fresh](https://github.com/gravityblast/fresh) for live reload.