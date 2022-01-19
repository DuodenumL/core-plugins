.PHONY: deps binary build test mock cloc unit-test

REPO_PATH := github.com/projecteru2/core-plugins
REVISION := $(shell git rev-parse HEAD || unknown)
BUILTAT := $(shell date +%Y-%m-%dT%H:%M:%S)
VERSION := $(shell git describe --tags $(shell git rev-list --tags --max-count=1))
GO_LDFLAGS ?= -X $(REPO_PATH)/version.REVISION=$(REVISION) \
			  -X $(REPO_PATH)/version.BUILTAT=$(BUILTAT) \
			  -X $(REPO_PATH)/version.VERSION=$(VERSION)
ifneq ($(KEEP_SYMBOL), 1)
	GO_LDFLAGS += -s
endif

deps:
	env GO111MODULE=on go mod download
	env GO111MODULE=on go mod vendor

binary:
	CGO_ENABLED=0 go build -ldflags "$(GO_LDFLAGS)" -gcflags=all=-G=3 -o bin/cpumem cpumem/cpumem.go && \
	CGO_ENABLED=0 go build -ldflags "$(GO_LDFLAGS)" -gcflags=all=-G=3 -o bin/volume volume/volume.go && \
	CGO_ENABLED=0 go build -ldflags "$(GO_LDFLAGS)" -gcflags=all=-G=3 -o bin/storage storage/storage.go

build: deps binary

test: deps unit-test

.ONESHELL:

cloc:
	cloc --exclude-dir=vendor,3rdmocks,mocks,tools,gen --not-match-f=test .

unit-test:
	go vet `go list ./... | grep -v '/vendor/' | grep -v '/tools'` && \
	go test -race -timeout 240s -count=1 -cover ./cpumem/models/... \
	./cpumem/schedule/... \
	./cpumem/types/...

lint:
	golangci-lint run
