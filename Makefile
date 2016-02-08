.PHONY: all run docker

all: check test install

check: goimports govet

goimports:
	@echo checking go imports...
	@goimports -d .

govet:
	@echo checking go vet...
	@go tool vet .

test:
	@go get
	@go test -v ./...

clean:
	@-rm -v "vropsbot" 2>/dev/null
	@-rm -v "vropsbot.db" 2>/dev/null
	@-rm -v "$(GOPATH)/bin/vropsbot" 2>/dev/null
	@-docker rmi vropsbot 2>/dev/null

build:
	@echo build vropsbot
	@go build github.com/bruceadowns/vropsbot

install:
	@echo install vropsbot
	@go install github.com/bruceadowns/vropsbot

docker:
	@echo build docker image
	@docker build --no-cache --file docker/Dockerfile --tag vropsbot .
