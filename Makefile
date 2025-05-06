SERVICES := identity messaging integration notification
DEFAULT_SERVICE := messaging
GO := go
GO_VERSION := $(shell $(GO) version)
GOPATH := $(shell $(GO) env GOPATH)
GOBIN := $(GOPATH)/bin
LDFLAGS := "-w -s" # Strip debug information
DOCKER_COMPOSE := docker compose
COMPOSE_FILE := deployments/compose.yml
COMPOSE_ENV_FILE := deployments/.env
MIGRATE := $(GOBIN)/migrate
MIGRATION_PATH_PATTERN = internal/%/infrastructure/persistence/migrations
DB_URL_ENV_VAR_PATTERN = $(shell echo $(1) | tr '[:lower:]' '[:upper:]')_DB_URL
GOLANGCI_LINT := $(GOBIN)/golangci-lint

$(COMPOSE_ENV_FILE):
	@echo "# Messaging Service" >> $(COMPOSE_ENV_FILE)
	@echo "MESSAGING_POSTGRES_USER=root" >> $(COMPOSE_ENV_FILE)
	@echo "MESSAGING_POSTGRES_PASSWORD=secret" >> $(COMPOSE_ENV_FILE)
	@echo "MESSAGING_DB_NAME=messaging_db" >> $(COMPOSE_ENV_FILE)
	@echo "MESSAGING_DB_PORT=5433" >> $(COMPOSE_ENV_FILE)
	@echo "MESSAGING_REDIS_PORT=6380" >> $(COMPOSE_ENV_FILE)
	@echo "MESSAGING_HTTP_PORT=8081" >> $(COMPOSE_ENV_FILE)
	@echo "MESSAGING_DB_URL=pgx5://root:secret@messaging_postgres:5432/messaging_db?sslmode=disable" >> $(COMPOSE_ENV_FILE)
	@echo "MESSAGING_KAFKA_BROKERS=kafka:9092" >> $(COMPOSE_ENV_FILE)
	@echo "MESSAGING_KAFKA_DEFAULT_TOPIC=meridian.messaging.events" >> $(COMPOSE_ENV_FILE)
	@echo "" >> $(COMPOSE_ENV_FILE)


.PHONY: build
build: tidy ## Build all service binaries
	@echo "Building Go binaries..."
	@$(foreach service,$(SERVICES), \
		echo " > Building $(service)..."; \
		CGO_ENABLED=0 $(GO) build -ldflags=$(LDFLAGS) -o ./bin/$(service) ./cmd/$(service); \
	)
	@echo "Build complete. Binaries in ./bin/"

.PHONY: build-one
build-one: tidy ## Build a specific service binary (e.g., make build-one service=messaging)
	@$(if $(service),,$(error Please specify the service, e.g., make build-one service=messaging))
	@echo "Building $(service) binary..."
	@CGO_ENABLED=0 $(GO) build -ldflags=$(LDFLAGS) -o ./bin/$(service) ./cmd/$(service)
	@echo "Build complete for $(service). Binary in ./bin/$(service)"

.PHONY: tidy
tidy: ## Tidy Go module files
	@echo "Running go mod tidy..."
	@$(GO) mod tidy
	@$(GO) mod vendor # Optional: if you vendor dependencies

.PHONY: fmt
fmt: ## Format Go source code
	@echo "Formatting Go code..."
	@go fmt ./...

.PHONY: lint
lint: ## Lint Go source code using golangci-lint
	@if ! command -v $(GOLANGCI_LINT) &> /dev/null; then \
		echo "golangci-lint not found. Installing..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
	fi
	@echo "Running golangci-lint..."
	@$(GOLANGCI_LINT) run ./...

.PHONY: test
test: tidy ## Run Go tests for all modules
	@echo "Running Go tests..."
	@$(GO) test -v -race -cover ./...
.PHONY: docker-build
docker-build: $(COMPOSE_ENV_FILE) ## Build Docker images using Docker Compose
	@echo "Building Docker images defined in $(COMPOSE_FILE)..."
	@$(DOCKER_COMPOSE) -f $(COMPOSE_FILE) --env-file $(COMPOSE_ENV_FILE) build $(service) # Build specific or all

.PHONY: docker-up
docker-up: $(COMPOSE_ENV_FILE) ## Start services using Docker Compose (detached mode)
	@echo "Starting services via Docker Compose..."
	@$(DOCKER_COMPOSE) -f $(COMPOSE_FILE) --env-file $(COMPOSE_ENV_FILE) up -d --remove-orphans $(service) # Start specific or all

.PHONY: docker-down
docker-down: $(COMPOSE_ENV_FILE) ## Stop and remove Docker Compose services
	@echo "Stopping and removing services via Docker Compose..."
	@$(DOCKER_COMPOSE) -f $(COMPOSE_FILE) --env-file $(COMPOSE_ENV_FILE) down $(options) # Pass options like -v

.PHONY: docker-stop
docker-stop: $(COMPOSE_ENV_FILE) ## Stop Docker Compose services without removing them
	@echo "Stopping services via Docker Compose..."
	@$(DOCKER_COMPOSE) -f $(COMPOSE_FILE) --env-file $(COMPOSE_ENV_FILE) stop $(service)

.PHONY: docker-logs
docker-logs: $(COMPOSE_ENV_FILE) ## Follow logs from Docker Compose services
	@echo "Following logs (Ctrl+C to stop)..."
	@$(DOCKER_COMPOSE) -f $(COMPOSE_FILE) --env-file $(COMPOSE_ENV_FILE) logs -f $(service) # Follow specific or all

.PHONY: docker-ps
docker-ps: $(COMPOSE_ENV_FILE) ## List running Docker Compose services
	@echo "Listing running services..."
	@$(DOCKER_COMPOSE) -f $(COMPOSE_FILE) --env-file $(COMPOSE_ENV_FILE) ps


get_db_url_var_name = $(call DB_URL_ENV_VAR_PATTERN,$(1))
get_migration_path = $(subst %,$(1),$(MIGRATION_PATH_PATTERN))


ifneq (,$(wildcard $(COMPOSE_ENV_FILE)))
    include $(COMPOSE_ENV_FILE)
    export $(shell sed 's/=.*//' $(COMPOSE_ENV_FILE))
endif

.PHONY: migrate-create
migrate-create: ## Create new migration files (e.g., make migrate-create service=messaging name=add_foo)
	@$(if $(name),,$(error Please specify migration name, e.g., make migrate-create name=my_migration))
	@$(eval SERVICE_NAME := $(if $(service),$(service),$(DEFAULT_SERVICE)))
	@$(eval MIGRATION_DIR := $(call get_migration_path,$(SERVICE_NAME)))
	@echo "Creating migration '$(name)' for service '$(SERVICE_NAME)' in $(MIGRATION_DIR)..."
	@mkdir -p $(MIGRATION_DIR)
	@$(MIGRATE) create -ext sql -dir $(MIGRATION_DIR) -seq $(name)

.PHONY: migrate-up
migrate-up: ## Apply all pending UP migrations (e.g., make migrate-up service=messaging)
	@$(eval SERVICE_NAME := $(if $(service),$(service),$(DEFAULT_SERVICE)))
	@$(eval MIGRATION_DIR := $(call get_migration_path,$(SERVICE_NAME)))
	@$(eval DB_URL_VAR := $(call get_db_url_var_name,$(SERVICE_NAME)))
	@$(if $(MIGRATION_DIR),,$(error Migration directory not found for $(SERVICE_NAME)))
	@$(if $(DB_URL_VAR),,$(error DB URL environment variable name not found for $(SERVICE_NAME)))
	@$(if $(value $(DB_URL_VAR)),,$(error Environment variable $(DB_URL_VAR) is not set))
	@echo "Applying UP migrations for service '$(SERVICE_NAME)' from $(MIGRATION_DIR)..."
	@$(MIGRATE) -path $(MIGRATION_DIR) -database '$(value $(DB_URL_VAR))' up $(steps) # Pass optional steps

.PHONY: migrate-down
migrate-down: ## Roll back migrations (e.g., make migrate-down service=messaging steps=1)
	@$(eval SERVICE_NAME := $(if $(service),$(service),$(DEFAULT_SERVICE)))
	@$(eval MIGRATION_DIR := $(call get_migration_path,$(SERVICE_NAME)))
	@$(eval DB_URL_VAR := $(call get_db_url_var_name,$(SERVICE_NAME)))
	@$(if $(MIGRATION_DIR),,$(error Migration directory not found for $(SERVICE_NAME)))
	@$(if $(DB_URL_VAR),,$(error DB URL environment variable name not found for $(SERVICE_NAME)))
	@$(if $(value $(DB_URL_VAR)),,$(error Environment variable $(DB_URL_VAR) is not set))
	@$(eval STEPS_ARG := $(if $(steps),$(steps),1)) # Default to 1 step if not provided
	@echo "Rolling back $(STEPS_ARG) DOWN migration(s) for service '$(SERVICE_NAME)' from $(MIGRATION_DIR)..."
	@read -p "WARNING: Rolling back migrations. Are you sure? (y/N) " answer && [ $${answer:-N} = y ] || (echo "Aborted." && exit 1)
	@$(MIGRATE) -path $(MIGRATION_DIR) -database '$(value $(DB_URL_VAR))' down $(STEPS_ARG)

.PHONY: migrate-status
migrate-status: ## Check migration status (e.g., make migrate-status service=messaging)
	@$(eval SERVICE_NAME := $(if $(service),$(service),$(DEFAULT_SERVICE)))
	@$(eval MIGRATION_DIR := $(call get_migration_path,$(SERVICE_NAME)))
	@$(eval DB_URL_VAR := $(call get_db_url_var_name,$(SERVICE_NAME)))
	@$(if $(MIGRATION_DIR),,$(error Migration directory not found for $(SERVICE_NAME)))
	@$(if $(DB_URL_VAR),,$(error DB URL environment variable name not found for $(SERVICE_NAME)))
	@$(if $(value $(DB_URL_VAR)),,$(error Environment variable $(DB_URL_VAR) is not set))
	@echo "Checking migration status for service '$(SERVICE_NAME)'..."
	@$(MIGRATE) -path $(MIGRATION_DIR) -database '$(value $(DB_URL_VAR))' version
.PHONY: clean
clean: ## Remove build artifacts
	@echo "Cleaning build artifacts..."
	@rm -rf ./bin
	@$(GO) clean

.PHONY: help
help: ## Display this help screen
	@echo "Available commands:"
	@grep -h -E '^[a-zA-Z_-]+:.*##' $(MAKEFILE_LIST) | sed -e 's/\(.*\):.*##[ \t]*\(.*\)/  \1|\2/' | column -t -s '|' | sort

# --- Default Target ---
.DEFAULT_GOAL := help
