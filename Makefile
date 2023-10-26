GO := $(shell which go)
PATH := $(dir $(GO)):$(PATH)
SHELL := env PATH=$(PATH) /bin/bash
NAME := azion

ifeq (, $(GO))
$(error "No go binary found in your system, please install go 1.17 before continuing")
endif

GOPATH ?= $(shell $(GO) env GOPATH)
GOBIN ?= $(GOPATH)/bin
GOSEC ?= $(GOBIN)/gosec
GOLINT ?= $(GOBIN)/golint
GOFMT ?= $(GOBIN)/gofmt
RELOAD ?= $(GOBIN)/CompileDaemon

# Variables for token endpoints
ENVFILE ?= ./env/prod

BIN := azion
# Version Info
BIN_VERSION=$(shell git describe --tags)
# The variables with $$ should be sourced from an envfile
LDFLAGS=-X github.com/aziontech/azion-cli/pkg/cmd/version.BinVersion=$(BIN_VERSION) \
		-X github.com/aziontech/azion-cli/pkg/constants.StorageApiURL=$$STORAGE_URL \
		-X github.com/aziontech/azion-cli/pkg/constants.AuthURL=$$AUTH_URL \
		-X github.com/aziontech/azion-cli/pkg/constants.ApiURL=$$API_URL \
		-X github.com/aziontech/azion-cli/pkg/cmd/edge_applications/init.TemplateBranch=$$TEMPLATE_BRANCH \
		-X github.com/aziontech/azion-cli/pkg/cmd/edge_applications/init.TemplateMajor=$$TEMPLATE_MAJOR
LDFLAGS_STRIP=-s -w
NAME_WITH_VERSION=$(NAME)-$(BIN_VERSION)

.PHONY : all
all: deps build

.PHONY : deps
deps: ## verify projects dependencies
	@ $(GO) env -w GOPRIVATE=github.com/aziontech/*
	@ $(GO) mod verify
	@ $(GO) mod tidy

.PHONY : clean
clean: ## delete additional files
	@ rm -rf cover

.PHONY: lint
lint: get-lint-deps ## running GoLint
	@ $(GOBIN)/golangci-lint run ./... --verbose


.PHONY: dev
dev: dev-deps
	$(RELOAD) -build 'make build' -exclude-dir '.git'

.PHONY: dev-deps
dev-deps: 
	$(GO) install github.com/githubnemo/CompileDaemon@v1.4.0

.PHONY: get-lint-deps
get-lint-deps:
	@if [ ! -x $(GOBIN)/golangci-lint ]; then\
		curl -sfL \
		https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOBIN) v1.54.1 ;\
	fi

.PHONY: test
test:
	@ echo Running GO tests
	@ mkdir -p cover
	@$(GO) test -v -failfast -coverprofile "./cover/$(NAME)coverage.out" -coverpkg=./... ./...
	@$(GO) tool cover -html="./cover/$(NAME)coverage.out" -o ./cover/$(NAME)coverage.html
	@$(GO) tool cover -func "./cover/$(NAME)coverage.out"
	@echo Done

.PHONY: docs
docs:
	$(GO) run ./cmd/gen_docs/main.go --doc-path ./docs --file-type md

.PHONY: sec
sec: get-gosec-deps ## running GoSec
	@ -$(GOSEC) ./...

.PHONY: get-gosec-deps
get-gosec-deps:
	@ cd $(GOPATH); \
		$(GO) install -u github.com/securego/gosec/cmd/gosec
		
.PHONY : build
build: ## build application
	@ $(GO) version
	@ source $(ENVFILE) && $(GO) build -ldflags "$(LDFLAGS)" -o ./bin/$(NAME) ./cmd/$(BIN)

.PHONY : cross-build
cross-build: ## cross-compile for all platforms/architectures.
	@ $(GO) version
	set -e;\
	source $(ENVFILE); \
	while read spec; \
	do\
		distro=$$(echo $${spec} | cut -d/ -f1);\
		goarch=$$(echo $${spec} | cut -d/ -f2);\
		arch=$$(echo $${goarch} | sed 's/386/x86_32/g; s/amd64/x86_64/g; s/arm$$/arm32/g');\
		echo $$distro/$$arch;\
		mkdir -p dist/$$distro/$$arch;\
		CGO_ENABLED=0 GOOS=$$distro GOARCH=$$goarch $(GO) build -ldflags "$(LDFLAGS)" -o ./dist/$$distro/$$arch/$(NAME_WITH_VERSION) ./cmd/$(BIN); \
	done < BUILD
