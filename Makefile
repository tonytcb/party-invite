
SHELL := /bin/bash

.PHONY: help
## help: shows this help message
help:
	@ echo "Usage: make [target]"
	@ sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

## up: starts the application exposing its HTTP port
up: clean
	docker-compose up app
	docker-compose down

## clean: clean up all docker containers
clean:
	docker-compose down
	docker ps -aq | xargs docker stop | xargs docker rm

## test: runs all tests
test:
	docker-compose run test

## lint: runs linter for a given directory, specified via PACKAGE variable
lint:
	@ if [ -z "$(PACKAGE)" ]; then echo >&2 please set directory via variable PACKAGE; exit 2; fi
	@ docker run  --rm -v "`pwd`:/workspace:cached" -w "/workspace/$(PACKAGE)" golangci/golangci-lint:latest golangci-lint run

## lint-all: runs linter for all packages
lint-all:
	@ docker run  --rm -v "`pwd`:/workspace:cached" -w "/workspace/." golangci/golangci-lint:latest golangci-lint run


## unit-test-ci: runs all tests on CI, with no Docker
unit-test-ci:
	@go generate
	@go test -v -race ./...

## build-image-prod: build a docker image ready for production
build-image-prod:
	docker build -t party-invite-prod .
	## docker container run --name party-invite-prod -p 8080:8080 party-invite-prod