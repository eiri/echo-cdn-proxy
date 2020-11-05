.DEFAULT_GOAL := all

PROJECT := echo-cdn-proxy

all: deps build test

.PHONY: deps
deps:
	go get -t ./...

.PHONY: build
build:
	go build ./...

.PHONY: test
test:
	go test -v ./...

.PHONY: run
run:
	go build -o $(PROJECT) ./example/...
	$(CURDIR)/$(PROJECT)
