#!/usr/bin/env make -f

SHELL := /bin/bash

APP := git-rid-of-keys
OUTPUT := pre-commit

HAS_GOMETALINTER := $(shell command -v gometalinter;)
HAS_DEP := $(shell command -v dep;)

default: all

.PHONY: all
all: bootstrap deps fmt vet lint build

.PHONY: bootstrap
bootstrap:
ifndef HAS_GOMETALINTER
	go get -u github.com/alecthomas/gometalinter
endif
ifndef HAS_DEP
	go get -u github.com/golang/dep/cmd/dep
endif

.PHONY: fmt
fmt:
		go fmt ./...

.PHONY: vet
vet:
		go vet ./...

.PHONY: lint
lint:
		@gometalinter --install
		gometalinter          \
		--enable-gc           \
		--deadline 40s        \
		--exclude bindata     \
		--exclude .pb.        \
		--exclude vendor      \
		--skip vendor         \
		--disable-all         \
		--enable=errcheck     \
		--enable=goconst      \
		--enable=gofmt        \
		--enable=golint       \
		--enable=gosimple     \
		--enable=ineffassign  \
		--enable=gotype       \
		--enable=misspell     \
		--enable=vet          \
		--enable=vetshadow    \
		--no-vendored-linters \
		./...

.PHONY: deps
deps: bootstrap
		dep ensure

.PHONY: build
build: lint deps
ifeq (,$(wildcard bin))
	mkdir bin
endif
		go build -o bin/$(OUTPUT) cmd/$(APP)/main.go

.PHONY: install-global
install:
ifeq (,$(wildcard /etc/git/hooks))
	sudo mkdir -p /etc/git/hooks
endif
		sudo ln -sf $(PWD)/bin/$(OUTPUT) /etc/git/hooks/

.PHONY: uninstall-global
uninstall:
		sudo rm /etc/git/hooks/$(OUTPUT)

.PHONY: clean
clean:
ifneq (,$(wildcard bin))
	@rm -rf bin && echo "bin/ removed!"
endif
ifneq (,$(wildcard vendor))
	@rm -rf vendor && echo "vendor/ removed!"
endif
		@echo "cleaned!"

