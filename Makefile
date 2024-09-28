tests:
	go test -v -count=1 ./internal/...
.PHONY: tests

bench:
	go test -bench=. ./internal/...
.PHONY: bench

bench-mem:
	go test -bench=. ./internal/... -benchmem
.PHONY: bench-mem

lint:
	@go clean -cache
	golangci-lint run ./...
.PHONY: lint

install-lint:
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
.PHONY: install-lint

build:
	@go build -v -o app cmd/application/main.go
.PHONY: build

build-dockerfile:
	docker build -t user-service -f Dockerfile .
.PHONY: build-dockerfile

run-dockerfile:
	docker run -p 8080:8080 user-service
.PHONY:run-dockerfile

compose-up:
	docker-compose up
.PHONY: compose-up

env:
	@include configs/config.env
.PHONY: env

migrate-up:
	goose -dir internal/migrations postgres "postgresql://postgres:qwerty@127.0.0.1:5432?sslmode=disable" up
.PHONY: migrate-up

migrate-down:
	goose -dir internal/migrations postgres "postgresql://postgres:qwerty@127.0.0.1:5432?sslmode=disable" down
.PHONY: migrate-down

all-tests:
	make tests && make bench && make lint
.PHONY: all-tests

coverage:
	@go test -cover ./... -coverprofile=./coverage.out
	@go tool cover -func ./coverage.out -o ./coverage.tool.out
.PHONY: coverage