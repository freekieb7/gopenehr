DATABASE_URL := postgres://gopenehr:gopenehrpass@localhost:5432/gopenehr
LOG_LEVEL := DEBUG
OAUTH_TRUSTED_ISSUERS := https://dev-3pgtw461x5f8dd7k.us.auth0.com
OAUTH_AUDIENCE := http://localhost:3000

APP_VARIABLES := \
	DATABASE_URL=$(DATABASE_URL) \
	LOG_LEVEL=$(LOG_LEVEL) \
# 	OAUTH_TRUSTED_ISSUERS=$(OAUTH_TRUSTED_ISSUERS) \
	OAUTH_AUDIENCE=$(OAUTH_AUDIENCE)

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

.PHONY: migrate-up
migrate-up:
	$(APP_VARIABLES) go run -mod=vendor ./... migrate up

.PHONY: migrate-down
migrate-down:
	$(APP_VARIABLES) go run -mod=vendor ./... migrate down

.PHONY: serve
serve:
	$(APP_VARIABLES) go run -mod=vendor ./... serve

.PHONY: aql-gen
aql-gen:
	docker run --rm -u $(id -u ${USER}):$(id -g ${USER}) --volume `pwd`/internal/openehr/aql:/work antlr/antlr4 -Dlanguage=Go AQL.g4 -o gen -package gen
