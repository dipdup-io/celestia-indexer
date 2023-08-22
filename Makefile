-include .env
export $(shell sed 's/=.*//' .env)

indexer:
	cd cmd/indexer && go run . -c ../../build/dipdup.yml

api:
	cd cmd/api && go run . -c ../../build/dipdup.yml

build:
	docker-compose up -d -- build

lint:
	golangci-lint run

test:
	go test -p 8 ./...

## adr: Generate ADR from template. Must set NUM and TITLE parameters.
adr:
	@echo "Generating ADR"
	@cp adr/adr-template.md adr/adr-$(NUM)-$(TITLE).md

mock:
	go generate ./internal/storage

.PHONY: indexer api build lint test adr