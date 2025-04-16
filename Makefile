GOCMD=go
BINARY_DIR=bin

.PHONY: deps
deps: ## Install all Go dependencies
	$(GOCMD) mod tidy

.PHONY: clean
clean: ## Clean up all build artifacts
	$(GOCMD) clean
	rm -rf $(BINARY_DIR)

.PHONY: build
build: build-server build-cli ## Build all

.PHONY: build-server
build-server: ## Build server application
	mkdir -p $(BINARY_DIR)
	$(GOCMD) build -o $(BINARY_DIR)/server ./cmd/server/

.PHONY: build-cli
build-cli: ## Build CLI application
	mkdir -p $(BINARY_DIR)
	$(GOCMD) build -o $(BINARY_DIR)/rejoinderoo ./cmd/cli/

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
