.PHONY: help fmt build test race vet check bench

help:
	@echo 'Targets: fmt build test race vet check bench'
	@echo 'build checks packages; add a binary target when cmd/loglens exists.'

fmt:
	gofmt -w .

build:
	go build ./...

test:
	go test ./...

race:
	go test -race ./...

vet:
	go vet ./...

check: build test race vet

bench:
	go test -run '^$$' -bench . -benchmem ./...
