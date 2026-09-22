.PHONY: help db-up db-down migrate-up migrate-down run test lint build tidy

BACKEND_DIR := backend
MIGRATIONS_DIR := $(BACKEND_DIR)/migrations
POSTGRES_PORT ?= 5432
DATABASE_URL ?= postgres://hairhistory:hairhistory@localhost:$(POSTGRES_PORT)/hairhistory?sslmode=disable

export POSTGRES_PORT
export DATABASE_URL

help:
	@echo "db-up        start the local Postgres container"
	@echo "db-down      stop the local Postgres container"
	@echo "migrate-up   apply all migrations"
	@echo "migrate-down roll back the last migration"
	@echo "run          run the API server"
	@echo "test         run the Go test suite"
	@echo "lint         run gofmt and go vet"
	@echo "build        compile all Go packages"
	@echo "tidy         sync go.mod and go.sum"

db-up:
	docker compose up -d db
	@until docker compose exec -T db pg_isready -U hairhistory -d hairhistory >/dev/null 2>&1; do sleep 1; done
	@echo "postgres is ready on port $(POSTGRES_PORT)"

db-down:
	docker compose down

migrate-up:
	migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" up

migrate-down:
	migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" down 1

run:
	cd $(BACKEND_DIR) && go run ./cmd/api

test:
	cd $(BACKEND_DIR) && go test ./...

lint:
	@test -z "$$(cd $(BACKEND_DIR) && gofmt -l ./cmd ./internal)" || \
		(echo "gofmt needed:" && cd $(BACKEND_DIR) && gofmt -l ./cmd ./internal && exit 1)
	cd $(BACKEND_DIR) && go vet ./...

build:
	cd $(BACKEND_DIR) && go build ./...

tidy:
	cd $(BACKEND_DIR) && go mod tidy
