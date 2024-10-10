NO_COLOR=\033[0m
OK_COLOR=\033[32;01m
GOLANGCI_CMD := $(shell command -v golangci-lint 2> /dev/null)
PKGS := $(shell go list ./... | grep -v /vendor/ | grep -v /internal/mock)
ALL_PACKAGES := $(shell go list ./... | grep -v /vendor/ | grep -v /internal/mock)

check-golangci:
ifndef GOLANGCI_CMD
	$(error "Please install golangci linters from https://golangci-lint.run/usage/install/")
endif

lint: check-golangci fmt
	@echo -e "$(OK_COLOR)==> linting projects$(NO_COLOR)..."
	@golangci-lint run --fix

fmt:
	@go fmt $(ALL_PACKAGES)

test:
	@go test -coverprofile=test_coverage.out $(PKGS)
	go tool cover -html=test_coverage.out -o test_coverage.html
	rm test_coverage.out
	@echo -e "$(OK_COLOR)==> Open test_coverage.html file on your web browser for detailed coverage$(NO_COLOR)..."

deps:
	go mod tidy && go mod vendor

run:
	go run cmd/api/application.go

migration:
	@latest_migration=$$(ls toolkit/db/migrations/*.go | grep -v 'migrate.go' | sort | tail -n 1); \
	latest_version=$$(basename $$latest_migration .go | sed 's/[^0-9]*//g'); \
	new_version=$$(printf "%04d" $$((10#$$latest_version + 1))); \
	new_migration_file="toolkit/db/migrations/$$new_version.go"; \
	cp toolkit/db/migrations/0000.go $$new_migration_file; \
	echo "Creating migration file: $$new_migration_file with version $$new_version"; \
	sed -i '' -e 's/0000/'$$new_version'/g' -e '/^\/\/.*$$/d' $$new_migration_file; \
	echo -e "$(OK_COLOR)==> New migration created: $$new_migration_file $(NO_COLOR)"
