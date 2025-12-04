DATABASE_URL := postgres://gopenehr:gopenehrpass@localhost:5432/gopenehr
LOG_LEVEL := DEBUG
# KAFKA_BROKERS := localhost:9092

APP_VARIABLES := \
	DATABASE_URL=$(DATABASE_URL) \
	LOG_LEVEL=$(LOG_LEVEL) \
# 	KAFKA_BROKERS=$(KAFKA_BROKERS)

.PHONY: up
up:
	docker compose up -d --wait --remove-orphans

.PHONY: down
down:
	docker compose down --remove-orphans

.PHONY: lint
lint:
# 	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $(go env GOPATH)/bin v2.6.1
	~/go/bin/golangci-lint run

.PHONY: test
test:
	$(APP_VARIABLES) go test -mod=vendor ./... -v -coverprofile=coverage.out

.PHONY: migrate-up
migrate-up:
	$(APP_VARIABLES) go run -mod=vendor ./... migrate up

.PHONY: migrate-down
migrate-down:
	$(APP_VARIABLES) go run -mod=vendor ./... migrate down

.PHONY: serve
serve:
	$(APP_VARIABLES) go run -mod=vendor ./... serve

.PHONY: seed
seed:
ifeq ($(count),)
	$(APP_VARIABLES) go run -mod=vendor ./... seed
else 
	$(APP_VARIABLES) go run -mod=vendor ./... seed $(count) 
endif

.PHONY: aql-gen
aql-gen:
	docker run --rm -u $(id -u ${USER}):$(id -g ${USER}) --volume `pwd`/internal/openehr/aql:/work antlr/antlr4 -Dlanguage=Go AQL.g4 -o gen -package gen

.PHONY: pprof
pprof:
	go test -test.benchmem -cpuprofile cpu.prof -memprofile mem.prof -bench BenchmarkCompositionUnmarshal ./internal/openehr/rmv2/...
	go tool pprof -http localhost:8080 mem.prof
	