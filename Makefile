PATH := /usr/local/go/bin:$(PATH)
SHELL := env PATH=$(PATH) /bin/bash
GO := $(shell which go)
NAME := azioncli

GOPATH ?= $(shell $(GO) env GOPATH)
GOBIN ?= $(GOPATH)/bin
GOSEC ?= $(GOBIN)/gosec
GOLINT ?= $(GOBIN)/golint
GOFMT ?= $(GOBIN)/gofmt
RELOAD ?= $(GOBIN)/CompileDaemon

# Variables for token endpoints
ENVFILE ?= ./env/stage

# Version Info
BIN_VERSION=$(shell git describe --tags)
# The variables with $$ should be sourced from an envfile
LDFLAGS=-X github.com/aziontech/azion-cli/cmd/version.BinVersion=$(BIN_VERSION) \
		-X github.com/aziontech/azion-cli/pkg/token.AuthEndpoint=$$AUTH_URL \
		-X github.com/aziontech/azion-cli/cmd/edge_services/requests.ApiUrl=$$API_URL
LDFLAGS_STRIP=-s -w
NAME_WITH_VERSION=$(NAME)-$(BIN_VERSION)


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

.PHONY: test
test:
	@$(GO) test -v ./...

.PHONY: sec
sec: get-gosec-deps ## running GoSec
	@ -$(GOSEC) ./...

.PHONY: get-gosec-deps
get-gosec-deps:
	@ cd $(GOPATH); \
		$(GO) get -u github.com/securego/gosec/cmd/gosec
		
.PHONY : build
build: ## build application
	@ $(GO) version
	@ source $(ENVFILE) && $(GO) build -ldflags "$(LDFLAGS)" -o ./bin/$(NAME)

.PHONY : cross-build
cross-build: ## cross-compile for all platforms/architectures.
	@ $(GO) version
	set -e;\
	source $(ENVFILE); \
	while read spec; \
	do\
		distro=$$(echo $${spec} | cut -d/ -f1);\
		goarch=$$(echo $${spec} | cut -d/ -f2);\
		arch=$$(echo $${goarch} | sed 's/386/x86_32/g; s/amd64/x86_64/g; s/arm$$/arm32/g;');\
		mkdir -p dist/$$distro/$$arch;\
		CGO_ENABLED=0 GOOS=$$distro GOARCH=$$goarch $(GO) build -ldflags "$(LDFLAGS)" -o ./dist/$$distro/$$arch/$(NAME_WITH_VERSION); \
	done < BUILD
