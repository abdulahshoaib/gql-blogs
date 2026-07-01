.PHONY: build run generate test

build:
	go build -o bin/server ./...

run:
	go run ./...

generate:
	go run github.com/99designs/gqlgen generate

test:
	go test ./...
