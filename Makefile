
EXECUTABLES=git go docker

_=$(foreach exec,$(EXECUTABLES), \
	$(if $(shell which $(exec)), ok, $(error "No $(exec) in PATH")))

IMAGE=ivan1993spb/snake-bot

IMAGE_GOLANG=golang:1.21.5-alpine3.19
IMAGE_ALPINE=alpine:3.19

BINARY_NAME=snake-bot
VERSION=$(shell git describe --tags --abbrev=0 2>/dev/null || echo v0.0.0)
BUILD=$(shell git rev-parse --short HEAD)

LDFLAGS=-ldflags "-s -w -X main.Version=$(VERSION) -X main.Build=$(BUILD)"
DOCKER_BUILD_ARGS=\
 --build-arg VERSION=$(VERSION) \
 --build-arg BUILD=$(BUILD) \
 --build-arg IMAGE_GOLANG=$(IMAGE_GOLANG) \
 --build-arg IMAGE_ALPINE=$(IMAGE_ALPINE)

default: build

docker/build:
	@docker build $(DOCKER_BUILD_ARGS) -t $(IMAGE):$(VERSION) .
	@docker tag $(IMAGE):$(VERSION) $(IMAGE):latest
	@echo "Build $(BUILD) tagged $(IMAGE):$(VERSION)"
	@echo "Build $(BUILD) tagged $(IMAGE):latest"

docker/push:
	@echo "Push build $(BUILD) with tag $(IMAGE):$(VERSION)"
	@docker push $(IMAGE):$(VERSION)
	@echo "Push build $(BUILD) with tag $(IMAGE):latest"
	@docker push $(IMAGE):latest

build:
	@go build $(LDFLAGS) -v -o $(BINARY_NAME) ./cmd/snake-bot

install:
	@go install $(LDFLAGS) -v ./cmd/snake-bot

coverprofile:
	@go test -coverprofile=coverage.out ./...
	@go tool cover -func=coverage.out
	@go tool cover -html=coverage.out
