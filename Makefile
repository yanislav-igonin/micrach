.PHONY: run dev infra-up infra-down

run:
	go run .

dev:
	nodemon --exec go run . --signal SIGTERM

infra-up:
	docker compose up -d db

infra-down:
	docker compose down
