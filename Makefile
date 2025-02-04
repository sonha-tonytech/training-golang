APP_NAME = my-golang-app
BUILD_FLAGS = -ldflags="-s -w"
BUILD_DIR = ./build
GO = go
CONFIG_FILE = dbconfig.yml
ENVIRONMENT = development
MIGRATION_DIR := migrations
VERSION_FILE := migrations/version.txt

.DEFAULT_GOAL := help

build:
	@mkdir -p $(BUILD_DIR)
	$(GO) build $(BUILD_FLAGS) -o $(BUILD_DIR)/$(APP_NAME) main.go

start:
	reflex -r '\.go$$' -s -- sh -c '$(GO) run main.go'

test:
	$(GO) test ./... -v

clean:
	@rm -rf $(BUILD_DIR)

current-version:
	@if [ -f $(VERSION_FILE) ]; then \
		echo "Current version: $$(cat $(VERSION_FILE))"; \
	else \
		echo "No version found in $(VERSION_FILE)"; \
	fi

update-version:
	@echo "SQL migrate status output:"
	@sql-migrate status -config=$(CONFIG_FILE)
	@echo "Extracting version..."
	@sql-migrate status -config=$(CONFIG_FILE) | grep 'up.sql' | grep -v 'no' | tail -n 1 | sed 's/^| *//' | cut -d'_' -f1 > migrations/version.txt
	@if [ -s migrations/version.txt ]; then \
		echo "Version updated to: $$(cat migrations/version.txt)"; \
	else \
		echo "Error: Version not found!"; exit 1; \
	fi

db-migrate-up:
	@sql-migrate up -env=$(ENVIRONMENT) -config=$(CONFIG_FILE)
	@make update-version

db-migrate-down:
	@sql-migrate down -env=$(ENVIRONMENT) -config=$(CONFIG_FILE)
	@make update-version

sql-migrate: 
	$(GO) run ./migrations/migration.go

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'
