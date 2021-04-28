.PHONY: all

all: lint test build

build:
	CGO_ENABLED=0 go build -o git-telegram-bot -ldflags "-s -w" .

docker:
	docker build . -t git-telegram-bot

test:
	go test -v -p 1 -race ./...

lint:
	golangci-lint run ./...