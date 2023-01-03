PHONY: build

build: gomod build-server build-tools

build-server:
	go build -o bin/post ./server/

build-tools:
	go build -o bin/postctl postctl/main.go

gomod:
	go mod tidy

debug: build
	DEBUG=true ./bin/post
