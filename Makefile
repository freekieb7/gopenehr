DATABASE_URL=postgres://gopenehr:gopenehrpass@localhost:5432/gopenehr


.PHONY: lint
lint:
# 	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $(go env GOPATH)/bin v2.6.1
	~/go/bin/golangci-lint run

.PHONY: run
run:
	DATABASE_URL=${DATABASE_URL} go run -mod=vendor ./... serve