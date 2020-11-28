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
	go build -o example/$(PROJECT) ./example/...
	$(CURDIR)/example/$(PROJECT)

.PHONY: docker
docker:
	docker build -t $(PROJECT) .

.PHONY: run-docker
run-docker:
	docker run -it --rm -e PORT=5000 -p 8000:5000 $(PROJECT)
