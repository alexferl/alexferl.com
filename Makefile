.PHONY: dev audit fmt lint pre-commit run tidy update-deps

.DEFAULT: help
help:
	@echo "make dev"
	@echo "	setup development environment"
	@echo "make audit"
	@echo "	conduct quality checks"
	@echo "make fmt"
	@echo "	fix code format issues"
	@echo "make lint"
	@echo "	run lint checks"
	@echo "make pre-commit"
	@echo "	run pre-commit hooks"
	@echo "make run"
	@echo "	run application"
	@echo "make tidy"
	@echo "	clean and tidy dependencies"
	@echo "make update-deps"
	@echo "	update dependencies"

GOLANGCI_LINT := go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.10.1

check-pre-commit:
ifeq (, $(shell which pre-commit))
	$(error "pre-commit not in $(PATH), pre-commit (https://pre-commit.com) is required")
endif

dev: check-pre-commit
	pre-commit install

audit:
	go mod verify
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...

fmt:
	$(GOLANGCI_LINT) fmt

lint:
	$(GOLANGCI_LINT) run

pre-commit: check-pre-commit
	pre-commit run --all-files

run:
	go run . -local=true

tidy:
	go mod tidy -v

update-deps: tidy
	go get -u ./...
