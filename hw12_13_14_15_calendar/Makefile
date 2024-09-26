LOCAL_BIN:=$(CURDIR)/bin
CALENDAR_APP:=$(LOCAL_BIN)/calendar
DOCKER_IMG="calendar:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o $(CALENDAR_APP) -ldflags "$(LDFLAGS)" -mod vendor ./cmd/calendar/.

run: build
	$(CALENDAR_APP) --config ./configs/calendar/config.yaml

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f ./build/Dockerfile .

run-img: build-img
	docker run -p 8080:8080 -p 50051:50051 $(DOCKER_IMG)

version: build
	$(CALENDAR_APP) version

test:
	go test --race ./internal/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.59.1

lint: install-lint-deps
	golangci-lint run

install-generate-proto:
	GOBIN=$(LOCAL_BIN) go install github.com/bufbuild/buf/cmd/buf@v1.29.0
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.32
	GOBIN=$(LOCAL_BIN) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.22.0
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.22.0

generate: install-generate-proto
	buf generate $(CURDIR)/api

.PHONY: build run build-img run-img version test lint generate
