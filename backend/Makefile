.PHONY: build run dev test test-integration docker-run

build:
	go build -v -o ./bin/main ./cmd/server/main.go

run: build
	./bin/main

dev:
	air

test:
	go test -v ./...

test-integration:
	go test -v ./tests/integration/...

docker-run:
	docker-compose up -d