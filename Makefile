PATH := /usr/local/go/bin:$(PATH)
SHELL := env PATH=$(PATH) /bin/bash
GO := $(shell which go)
NAME := azion

GOPATH ?= $(shell $(GO) env GOPATH)
GOBIN ?= $(GOPATH)/bin
GOSEC ?= $(GOBIN)/gosec
GOLINT ?= $(GOBIN)/golint
GOFMT ?= $(GOBIN)/gofmt
RELOAD ?= $(GOBIN)/CompileDaemon
BUILD_DEBUG_VERSION ?= false

# Version Info
BIN_VERSION=$(shell git describe --tags)
LDFLAGS=-X github.com/aziontech/azion-cli/cmd.BinVersion=$(BIN_VERSION)

.PHONY : deps
deps: ## verify projects dependencies
	@ $(GO) env -w GOPRIVATE=github.com/aziontech/*
	@ $(GO) mod verify
	@ $(GO) mod tidy

.PHONY: lint
lint: get-lint-deps ## running GoLint
	@ $(GOBIN)/golangci-lint run ./...

.PHONY: get-lint-deps
get-lint-deps:
	@if [ ! -x $(GOBIN)/golangci-lint ]; then\
		curl -sfL \
		https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOBIN) v1.31.0 ;\
	fi

.PHONY : build
build: ## build application code
	@ $(GO) version
	@ $(GO) build -ldflags '$(LDFLAGS)' -o ./bin/$(NAME)