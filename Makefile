SHELL := /usr/bin/env bash
.SHELLFLAGS := -eu -o pipefail -c

GO ?= go
PKGS := ./pkg/dmesg
COVEROUT := coverage.out

.PHONY: help fmt-check vet test test coverage-report build ci

help:
	@echo "Available targets:"
	@echo "  make fmt-check  - ensure Go files are gofmt-formatted"
	@echo "  make vet        - run go vet"
	@echo "  make test - run unit tests and write coverage profile"
	@echo "  make coverage-report - print function coverage from profile"
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
	$(GO) test -coverprofile=$(COVEROUT) $(PKGS)

coverage-report: test
	$(GO) tool cover -func=$(COVEROUT)

build:
	$(GO) build cmd/go-dmesg.go

ci: fmt-check vet coverage-report build
