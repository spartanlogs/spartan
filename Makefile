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

# General development tasks
doc:
	@godoc -http=:6060 -index

fmt:
	@go fmt $$(go list ./... | grep -v /vendor/)

generate:
	@go generate $$(go list ./... | grep -v /vendor/)

# Testing
test:
	@go test $$(go list ./... | grep -v /vendor/)

test_race:
	@go test -race $$(go list ./... | grep -v /vendor/)

coverage:
	@go test -cover $$(go list ./... | grep -v /vendor/)

benchmark:
	@echo "Running tests..."
	@go test -bench=. $$(go list ./... | grep -v /vendor/)

# Code Checks
codecheck: lint vet gosimple staticcheck

# https://github.com/golang/lint
# go get github.com/golang/lint/golint
lint:
	@golint $$(go list ./... | grep -v /vendor/)

vet:
	@go vet $$(go list ./... | grep -v /vendor/)

# https://github.com/dominikh/go-tools/tree/master/cmd/gosimple
# go get honnef.co/go/tools/cmd/gosimple
gosimple:
	@gosimple $$(go list ./... | grep -v /vendor/)

# https://github.com/dominikh/go-tools/tree/master/cmd/staticcheck
# go get honnef.co/go/tools/cmd/staticcheck
staticcheck:
	@staticcheck $$(go list ./... | grep -v /vendor/)

install-code-checks:
	go get -u github.com/golang/lint/golint
	go get -u honnef.co/go/tools/cmd/gosimple
	go get -u honnef.co/go/tools/cmd/staticcheck

# Building the application
app:
	GOBIN="$(GOBIN)" go install -v -ldflags "$(LDFLAGS)" -tags '$(BUILDTAGS)' ./cmd/spartan
