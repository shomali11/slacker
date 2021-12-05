.PHONY: test format help lint 

.DEFAULT_GOAL := help

GO := go
GOPATH := $(shell go env GOPATH)
GOPATH_BIN := $(GOPATH)/bin
GOLANGCI_LINT := $(GOPATH_BIN)/golangci-lint
GO_PACKAGES := $(shell go list ./... | grep -v vendor)
UNAME := $(shell uname)

all: lint test

help:
	@echo "Slacker Makefile"
	@echo "test    - Test slacker"
	@echo "bench   - Bench slacker"
	@echo "race    - Race slacker"
	@echo "cover   - Cover slacker"
	@echo "format  - Format code using golangci-lint"
	@echo "help    - Prints help message"
	@echo "lint    - Lint code using golangci-lint"

format:
	@echo "Formatting..."
	@$(GO) fmt $(GO_PACKAGES)
	@$(GOLANGCI_LINT) run --fix --issues-exit-code 0 > /dev/null 2>&1
	@echo "Code formatted"

lint:
	@echo "Linting..."
	@$(GO) vet $(GO_PACKAGES)
	@$(GOLANGCI_LINT) run
	@echo "No errors found"

vendor:
	@echo "Tidy up go.mod..."
	@$(GO) mod tidy
	@echo "Vendoring..."
	@$(GO) mod vendor
	@echo "Done!"

install-golangcilint:
	@echo "Installing golangci-lint..."
	@curl -sSfL \
	 	https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | \
	 	sh -s -- -b $(GOPATH_BIN) v1.43.0
	@echo "Installed successfully"

test:
	@echo "Testing..."
	@$(GO) test -cover -race .

bench:
	@echo "Benching..."
	@$(GO) test -bench=. .

race:
	@echo "Racing..."
	@$(GO) test -v -race .

cover:
	@echo "Show coverage"
	@./script/coverage

# coveralls:
	# @echo "Show coverage and push to coveralls.io"
	# ./script/coverage --coveralls
