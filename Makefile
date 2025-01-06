APP_NAME = my-golang-app
BUILD_FLAGS = -ldflags="-s -w"
BUILD_DIR = ./build
GO = go
CONFIG_FILE = dbconfig.yml

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

db-migrate:
	sql-migrate up -env=production -config=$(CONFIG_FILE)

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'
