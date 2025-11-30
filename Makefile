DATABASE_URL := postgres://gopenehr:gopenehrpass@localhost:5432/gopenehr
LOG_LEVEL := DEBUG

APP_VARIABLES := \
	DATABASE_URL=$(DATABASE_URL) \
	LOG_LEVEL=$(LOG_LEVEL) \

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

.PHONY: pprof
pprof:
	go test -test.benchmem -cpuprofile cpu.prof -memprofile mem.prof -bench BenchmarkCompositionUnmarshal ./internal/openehr/rmv2/...
	go tool pprof -http localhost:8080 mem.prof
	