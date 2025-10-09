.PROXY: run
run:
	go run ./cmd/main.go

.PHONY: up
up:
	docker-compose up -d

.PHONY: down
down:
	docker-compose down

.PHONY: air
air:
	~/go/bin/air -c .air.toml

.PHONY: dump
dump:
	docker compose exec -u postgres db pg_dump postgres --schema-only --clean --if-exists > var/pg/dump2.sql

.PHONY: buf-gen
buf-gen:
	docker run --rm --volume "$(pwd):/workspace" --workdir /workspace bufbuild/buf generate

.PROXY: aql-gen
aql-gen:
	docker run --rm -u $(id -u ${USER}):$(id -g ${USER}) --volume `pwd`/internal/aql:/work antlr/antlr4 -Dlanguage=Go AQL.g4 -o gen -package gen

.PHONY: bench
bench:
	go build -o aql-bench ./cmd/bench/

.PHONY: bench-simple
bench-simple: bench
	./aql-bench --suite cmd/bench/simple-suite.json --verbose

.PHONY: bench-production
bench-production: bench
	./aql-bench --suite cmd/bench/production-suite.json --verbose

.PHONY: bench-quick
bench-quick: bench
	./aql-bench --suite cmd/bench/production-suite.json --tags "quick,count" --verbose

.PHONY: bench-junit
bench-junit: bench
	./aql-bench --format junit --output test-results.xml --verbose