.PHONY: fmt vet test test-race check build run

GOCACHE ?= /tmp/viger-go-cache
export GOCACHE
GOMODCACHE ?= /tmp/viger-go-mod
export GOMODCACHE

fmt:
	gofmt -w ./cmd ./internal

vet:
	go vet ./...

test:
	go test ./...

test-race:
	go test -race ./...

check:
	test -z "$$(gofmt -l ./cmd ./internal)"
	go vet ./...
	go test -race ./...

build:
	mkdir -p bin
	go build -trimpath -o bin/viger-api ./cmd/server

run:
	go run ./cmd/server
