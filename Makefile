SHELL := /usr/bin/env bash
.SHELLFLAGS := -eu -o pipefail -c

GO ?= go
PKGS := ./...

.PHONY: help fmt-check vet test build ci

help:
	@echo "Available targets:"
	@echo "  make fmt-check  - ensure Go files are gofmt-formatted"
	@echo "  make vet        - run go vet"
	@echo "  make test       - run unit tests"
	@echo "  make build      - build all packages"
	@echo "  make ci         - run all PR CI checks"

fmt-check:
	@unformatted="$$(gofmt -l . | grep -E '\\.go$$' || true)"; \
	if [[ -n "$$unformatted" ]]; then \
		echo "The following files are not gofmt-formatted:"; \
		echo "$$unformatted"; \
		exit 1; \
	fi

vet:
	$(GO) vet $(PKGS)

test:
	$(GO) test $(PKGS)

build:
	$(GO) build $(PKGS)

ci: fmt-check vet test build
