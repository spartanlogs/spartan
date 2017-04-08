# variable definitions
NAME := spartan
DESC := Data processor for logging and other messages
VERSION := $(shell git describe --tags --always --dirty)
GOVERSION := $(shell go version)
BUILDTIME := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
BUILDER := $(shell echo "`git config user.name` <`git config user.email`>")
CGO_ENABLED ?= 1
PWD := $(shell pwd)
GOBIN := $(PWD)/bin

ifeq ($(shell uname -o), Cygwin)
PWD := $(shell cygpath -w -a `pwd`)
GOBIN := $(PWD)\bin
endif

PROJECT_URL := "https://github.com/lfkeitel/$(NAME)"
BUILDTAGS ?= dball
LDFLAGS := -X 'main.version=$(VERSION)' \
			-X 'main.buildTime=$(BUILDTIME)' \
			-X 'main.builder=$(BUILDER)' \
			-X 'main.goversion=$(GOVERSION)'

.PHONY: all doc fmt alltests test coverage benchmark lint vet app dist clean docker

all: test app

# development tasks
doc:
	@godoc -http=:6060 -index

fmt:
	@go fmt $$(go list ./src/...)

alltests: test lint vet

test:
	@go test -race $$(go list ./src/...)

coverage:
	@go test -cover $$(go list ./src/...)

benchmark:
	@echo "Running tests..."
	@go test -bench=. $$(go list ./src/...)

# https://github.com/golang/lint
# go get github.com/golang/lint/golint
lint:
	@golint ./src/...

vet:
	@go vet $$(go list ./src/...)

app:
	GOBIN="$(GOBIN)" go install -v -ldflags "$(LDFLAGS)" -tags '$(BUILDTAGS)' ./cmd/spartan
