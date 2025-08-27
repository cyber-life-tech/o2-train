IMPORTS_REVISER_VERSION ?= latest
imports-deps:
	go install -mod=mod github.com/incu6us/goimports-reviser/v3@$(IMPORTS_REVISER_VERSION)

imports:
	@echo "Running imports"
	goimports-reviser -rm-unused -set-alias -format ./...

GOLANG_CI_LINT_VERSION ?= latest
lint-deps:
	go install -mod=mod github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANG_CI_LINT_VERSION)

lint:
	@echo "Running linter"
	golangci-lint run