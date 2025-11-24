DATABASE_URL := postgres://gopenehr:gopenehrpass@localhost:5432/gopenehr
MIGRATIONS_DIR := ./internal/database/migrations

APP_VARIABLES := \
	DATABASE_URL=$(DATABASE_URL) \
	MIGRATIONS_DIR=$(MIGRATIONS_DIR)

.PHONY: up
up:
	docker compose up -d --wait

.PHONY: down
down:
	docker compose down

.PHONY: lint
lint:
# 	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $(go env GOPATH)/bin v2.6.1
	~/go/bin/golangci-lint run

.PHONY: migrate
migrate:
	$(APP_VARIABLES) go run -mod=vendor ./... migrate up

.PHONY: run
run:
	$(APP_VARIABLES) go run -mod=vendor ./... serve

.PHONY: aql-gen
aql-gen:
	docker run --rm -u $(id -u ${USER}):$(id -g ${USER}) --volume `pwd`/internal/openehr/aql:/work antlr/antlr4 -Dlanguage=Go AQL.g4 -o gen -package gen
