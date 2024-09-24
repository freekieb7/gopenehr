#!make
include .env
.DEFAULT_GOAL = help

.PHONY: help
help:
	@echo "------------------------------------------------"
	@echo "|    Usage: make ACTION [ARGS]                 |"
	@echo "|----------------------------------------------|"
	@echo "|    Example:                                  |"
	@echo "|        make update packages=\"code24/*\"       |"
	@echo "------------------------------------------------"
	@echo ""
	@echo "Common commands:"
	@cat $(MAKEFILE_LIST) | grep -E '^[a-zA-Z_-]+:.*?##' | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
#
#.PHONY: setup
#setup: stop clean login pull build start install setup-db seed-db gen-aql-data ## Setup project
#
.PHONY: pull
pull: ## Pull Docker images
	docker compose pull

.PHONY: build
build: ## Build Docker images
	docker compose build

.PHONY: up
up: ## Starts Docker images
	docker compose up --remove-orphans --detach --wait

.PHONY: down
down: ## Stops Docker images
	docker compose down --remove-orphans --volumes

.PHONY: restart
restart: ## Restarts Docker images
	docker compose restart

.PHONY: migrate-up
migrate-up: ## Restarts Docker images
	@docker compose run --rm migrate -path=/migrations/ -database ${DB_CONN_STR} up

.PHONY: migrate-down
migrate-down: ## Restarts Docker images
	@docker compose run --rm migrate -path=/migrations/ -database ${DB_CONN_STR} down

.PHONY: gen-antlr
gen-antlr: ## Generate ANTLR AQL code
	docker build . -t antlr/antlr4 --platform linux/amd64 -f aql/Dockerfile
	docker run --rm -u $(id -u ${USER}):$(id -g ${USER}) -v `pwd`/aql:/work antlr/antlr4 -Dlanguage=Go AqlLexer.g4 -o gen -package gen
	docker run --rm -u $(id -u ${USER}):$(id -g ${USER}) -v `pwd`/aql:/work antlr/antlr4 -Dlanguage=Go AqlParser.g4 -o gen -package gen
#
#.PHONY: gen-aql-data
#gen-aql-data: ## Generate Tables and Data required for AQL
#	docker compose run --rm web-tools bash -c "cd openPlatform/tests; php gen-aql-data.php"
#
#.PHONY: seed-db
#seed-db: ## Seed database with EHR / Composition data
#	docker compose run --rm web-tools bash -c "cd tests; php seed.php --rounds=10"
#
#.PHONY: bash-web
#bash-web-tools: ## Bash shell into web-tools
#	docker compose run --rm web-tools bash
#
#.PHONY: bash-db
#bash-db: ## Bash shell into db (pre: must be up and running)
#	docker compose exec db bash
#
#.PHONY: clean
#clean: ## Removes vendor folder
#	rm -rf vendor
