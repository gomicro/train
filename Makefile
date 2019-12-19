# System Setup
SHELL = bash

# Go Stuff
GOCMD=go
GOLINTCMD=golint
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOLIST=$(GOCMD) list
GOVET=$(GOCMD) vet
GOTEST=$(GOCMD) test -v $(shell $(GOCMD) list ./... | grep -v /vendor/)
GOFMT=$(GOCMD) fmt
CGO_ENABLED ?= 0
GOOS ?= $(shell uname -s | tr '[:upper:]' '[:lower:]')

# General Vars
APP := $(shell basename $(PWD) | tr '[:upper:]' '[:lower:]')

GIT_COMMIT_HASH ?= $(shell git rev-parse --short HEAD)
BUILD_VERSION := dev-$(GIT_COMMIT_HASH)


.PHONY: all
all: test build

.PHONY: build
build: ## Build the project
	CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) $(GOBUILD) -ldflags "-X github.com/gomicro/train/cmd.Version=$(BUILD_VERSION)" -o $(APP) .

.PHONY: clean
clean:  ## Clean out all generated files
	-@$(GOCLEAN)

.PHONY: help
help:  ## Show This Help
	@for line in $$(cat Makefile | grep "##" | grep -v "grep" | sed  "s/:.*##/:/g" | sed "s/\ /!/g"); do verb=$$(echo $$line | cut -d ":" -f 1); desc=$$(echo $$line | cut -d ":" -f 2 | sed "s/!/\ /g"); printf "%-30s--%s\n" "$$verb" "$$desc"; done

.PHONY: install
install:  ## Installs the binary
	$(GOCMD) install

.PHONY: test
test: unit_test  ## Runs all available tests

.PHONY: unit_test
unit_test:  ## Run unit tests
	$(GOTEST)

.PHONY: fmt
fmt: ## Run gofmt
	@echo "checking formatting..."
	@$(GOFMT) $(shell $(GOLIST) ./... | grep -v '/vendor/')

.PHONY: vet
vet: ## Run go vet
	@echo "vetting..."
	@$(GOVET) $(shell $(GOLIST) ./... | grep -v '/vendor/')

.PHONY: lint
lint: ## Run golint
	@echo "linting..."
	@$(GOLINTCMD) -set_exit_status $(shell $(GOLIST) ./... | grep -v '/vendor/')
