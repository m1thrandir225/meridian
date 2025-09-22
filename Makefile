SERVICES := identity messaging integration analytics
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

# Generate comprehensive .env file for all services
$(COMPOSE_ENV_FILE):
	@echo "Generating comprehensive .env file..."
	@echo "# ========================================" > $(COMPOSE_ENV_FILE)
	@echo "# Meridian Microservices Environment File" >> $(COMPOSE_ENV_FILE)
	@echo "# ========================================" >> $(COMPOSE_ENV_FILE)
	@echo "" >> $(COMPOSE_ENV_FILE)
	@echo "# Global Environment Settings" >> $(COMPOSE_ENV_FILE)
	@echo "ENVIRONMENT=development" >> $(COMPOSE_ENV_FILE)
	@echo "LOG_LEVEL=info" >> $(COMPOSE_ENV_FILE)
	@echo "" >> $(COMPOSE_ENV_FILE)
	@echo "# ========================================" >> $(COMPOSE_ENV_FILE)
	@echo "# Identity Service Configuration" >> $(COMPOSE_ENV_FILE)
	@echo "# ========================================" >> $(COMPOSE_ENV_FILE)
	@echo "IDENTITY_POSTGRES_USER=root" >> $(COMPOSE_ENV_FILE)
	@echo "IDENTITY_POSTGRES_PASSWORD=secret" >> $(COMPOSE_ENV_FILE)
	@echo "IDENTITY_DB_NAME=identity_db" >> $(COMPOSE_ENV_FILE)
	@echo "IDENTITY_DB_PORT=5432" >> $(COMPOSE_ENV_FILE)
	@echo "IDENTITY_REDIS_PORT=6379" >> $(COMPOSE_ENV_FILE)
	@echo "IDENTITY_HTTP_PORT=8080" >> $(COMPOSE_ENV_FILE)
	@echo "IDENTITY_GRPC_PORT=9090" >> $(COMPOSE_ENV_FILE)
	@echo "IDENTITY_DB_URL=postgres://root:secret@identity_postgres:5432/identity_db?sslmode=disable" >> $(COMPOSE_ENV_FILE)
	@echo "IDENTITY_DB_URL_MIGRATE=pgx5://root:secret@identity_postgres:5432/identity_db?sslmode=disable" >> $(COMPOSE_ENV_FILE)
	@echo "IDENTITY_REDIS_URL=redis://identity_redis:6379" >> $(COMPOSE_ENV_FILE)
	@echo "IDENTITY_KAFKA_BROKERS=kafka:9092" >> $(COMPOSE_ENV_FILE)
	@echo "IDENTITY_KAFKA_DEFAULT_TOPIC=meridian.identity.events" >> $(COMPOSE_ENV_FILE)
	@echo "IDENTITY_PASETO_PRIVATE_KEY=YOUR_PRIVATE_KEY_HERE" >> $(COMPOSE_ENV_FILE)
	@echo "IDENTITY_PASETO_PUBLIC_KEY=YOUR_PUBLIC_KEY_HERE" >> $(COMPOSE_ENV_FILE)
	@echo "AUTH_TOKEN_VALIDITY_MINUTES=60" >> $(COMPOSE_ENV_FILE)
	@echo "REFRESH_TOKEN_VALIDITY_MINUTES=1440" >> $(COMPOSE_ENV_FILE)
	@echo "INTEGRATION_GRPC_URL=integration:9091" >> $(COMPOSE_ENV_FILE)
	@echo "IDENTITY_ENVIRONMENT=development" >> $(COMPOSE_ENV_FILE)
	@echo "IDENTITY_LOG_LEVEL=info" >> $(COMPOSE_ENV_FILE)
	@echo "" >> $(COMPOSE_ENV_FILE)
	@echo "# ========================================" >> $(COMPOSE_ENV_FILE)
	@echo "# Messaging Service Configuration" >> $(COMPOSE_ENV_FILE)
	@echo "# ========================================" >> $(COMPOSE_ENV_FILE)
	@echo "MESSAGING_POSTGRES_USER=root" >> $(COMPOSE_ENV_FILE)
	@echo "MESSAGING_POSTGRES_PASSWORD=secret" >> $(COMPOSE_ENV_FILE)
	@echo "MESSAGING_DB_NAME=messaging_db" >> $(COMPOSE_ENV_FILE)
	@echo "MESSAGING_DB_PORT=5433" >> $(COMPOSE_ENV_FILE)
	@echo "MESSAGING_REDIS_PORT=6380" >> $(COMPOSE_ENV_FILE)
	@echo "MESSAGING_HTTP_PORT=8081" >> $(COMPOSE_ENV_FILE)
	@echo "MESSAGING_GRPC_PORT=9091" >> $(COMPOSE_ENV_FILE)
	@echo "MESSAGING_DB_URL=postgres://root:secret@messaging_postgres:5433/messaging_db?sslmode=disable" >> $(COMPOSE_ENV_FILE)
	@echo "MESSAGING_DB_URL_MIGRATE=pgx5://root:secret@messaging_postgres:5433/messaging_db?sslmode=disable" >> $(COMPOSE_ENV_FILE)
	@echo "MESSAGING_REDIS_URL=redis://messaging_redis:6380" >> $(COMPOSE_ENV_FILE)
	@echo "MESSAGING_KAFKA_BROKERS=kafka:9092" >> $(COMPOSE_ENV_FILE)
	@echo "MESSAGING_KAFKA_DEFAULT_TOPIC=meridian.messaging.events" >> $(COMPOSE_ENV_FILE)
	@echo "IDENTITY_GRPC_URL=identity:9090" >> $(COMPOSE_ENV_FILE)
	@echo "INTEGRATION_GRPC_URL=integration:9091" >> $(COMPOSE_ENV_FILE)
	@echo "MESSAGING_ENVIRONMENT=development" >> $(COMPOSE_ENV_FILE)
	@echo "MESSAGING_LOG_LEVEL=info" >> $(COMPOSE_ENV_FILE)
	@echo "" >> $(COMPOSE_ENV_FILE)
	@echo "# ========================================" >> $(COMPOSE_ENV_FILE)
	@echo "# Integration Service Configuration" >> $(COMPOSE_ENV_FILE)
	@echo "# ========================================" >> $(COMPOSE_ENV_FILE)
	@echo "INTEGRATION_POSTGRES_USER=root" >> $(COMPOSE_ENV_FILE)
	@echo "INTEGRATION_POSTGRES_PASSWORD=secret" >> $(COMPOSE_ENV_FILE)
	@echo "INTEGRATION_DB_NAME=integration_db" >> $(COMPOSE_ENV_FILE)
	@echo "INTEGRATION_DB_PORT=5434" >> $(COMPOSE_ENV_FILE)
	@echo "INTEGRATION_REDIS_PORT=6381" >> $(COMPOSE_ENV_FILE)
	@echo "INTEGRATION_HTTP_PORT=8082" >> $(COMPOSE_ENV_FILE)
	@echo "INTEGRATION_GRPC_PORT=9092" >> $(COMPOSE_ENV_FILE)
	@echo "INTEGRATION_DB_URL=postgres://root:secret@integration_postgres:5434/integration_db?sslmode=disable" >> $(COMPOSE_ENV_FILE)
	@echo "INTEGRATION_DB_URL_MIGRATE=pgx5://root:secret@integration_postgres:5434/integration_db?sslmode=disable" >> $(COMPOSE_ENV_FILE)
	@echo "INTEGRATION_REDIS_URL=redis://integration_redis:6381" >> $(COMPOSE_ENV_FILE)
	@echo "INTEGRATION_KAFKA_BROKERS=kafka:9092" >> $(COMPOSE_ENV_FILE)
	@echo "INTEGRATION_KAFKA_DEFAULT_TOPIC=meridian.integration.events" >> $(COMPOSE_ENV_FILE)
	@echo "MESSAGING_GRPC_URL=messaging:9091" >> $(COMPOSE_ENV_FILE)
	@echo "INTEGRATION_ENVIRONMENT=development" >> $(COMPOSE_ENV_FILE)
	@echo "INTEGRATION_LOG_LEVEL=info" >> $(COMPOSE_ENV_FILE)
	@echo "" >> $(COMPOSE_ENV_FILE)
	@echo "# ========================================" >> $(COMPOSE_ENV_FILE)
	@echo "# Analytics Service Configuration" >> $(COMPOSE_ENV_FILE)
	@echo "# ========================================" >> $(COMPOSE_ENV_FILE)
	@echo "ANALYTICS_POSTGRES_USER=root" >> $(COMPOSE_ENV_FILE)
	@echo "ANALYTICS_POSTGRES_PASSWORD=secret" >> $(COMPOSE_ENV_FILE)
	@echo "ANALYTICS_DB_NAME=analytics_db" >> $(COMPOSE_ENV_FILE)
	@echo "ANALYTICS_DB_PORT=5435" >> $(COMPOSE_ENV_FILE)
	@echo "ANALYTICS_HTTP_PORT=8084" >> $(COMPOSE_ENV_FILE)
	@echo "ANALYTICS_DB_URL=postgres://root:secret@analytics_postgres:5432/analytics_db?sslmode=disable" >> $(COMPOSE_ENV_FILE)
	@echo "ANALYTICS_DB_URL_MIGRATE=pgx5://root:secret@analytics_postgres:5432/analytics_db?sslmode=disable" >> $(COMPOSE_ENV_FILE)
	@echo "ANALYTICS_KAFKA_BROKERS=kafka:9092" >> $(COMPOSE_ENV_FILE)
	@echo "ANALYTICS_CONSUMER_GROUP=analytics-service" >> $(COMPOSE_ENV_FILE)
	@echo "ANALYTICS_ENVIRONMENT=development" >> $(COMPOSE_ENV_FILE)
	@echo "ANALYTICS_LOG_LEVEL=info" >> $(COMPOSE_ENV_FILE)
	@echo "" >> $(COMPOSE_ENV_FILE)
	@echo "# ========================================" >> $(COMPOSE_ENV_FILE)
	@echo "# Frontend Configuration" >> $(COMPOSE_ENV_FILE)
	@echo "# ========================================" >> $(COMPOSE_ENV_FILE)
	@echo "FRONTEND_PORT=3000" >> $(COMPOSE_ENV_FILE)
	@echo "VITE_API_BASE_URL=http://localhost:8080" >> $(COMPOSE_ENV_FILE)
	@echo "VITE_WS_BASE_URL=ws://localhost:8081" >> $(COMPOSE_ENV_FILE)
	@echo "" >> $(COMPOSE_ENV_FILE)
	@echo "Environment file generated successfully at $(COMPOSE_ENV_FILE)"
	@echo "⚠️  Remember to update the PASETO keys with actual values!"

.PHONY: generate-keys
generate-keys: ## Generate PASETO security keys using the keygen script
	@echo "🔑 Generating PASETO V4 security keys..."
	@$(GO) run scripts/keygen.go
	@echo ""
	@echo " Copy the generated keys to your .env file:"
	@echo "   - Replace IDENTITY_PASETO_PRIVATE_KEY with the private key"
	@echo "   - Replace IDENTITY_PASETO_PUBLIC_KEY with the public key"

.PHONY: generate-proto
generate-proto: ## Generate gRPC Go code from protobuf files
	@echo " Generating gRPC Go code from protobuf files..."
	@echo "Generating messaging service protobuf..."
	@protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		internal/messaging/infrastructure/api/messaging.proto
	@echo "Generating integration service protobuf..."
	@protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		internal/integration/infrastructure/api/integration.proto
	@echo "Generating identity service protobuf..."
	@protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		internal/identity/infrastructure/api/identity.proto
	@echo "✅ gRPC code generation complete!"

.PHONY: setup
setup: generate-keys $(COMPOSE_ENV_FILE) ## Complete setup: generate keys and env file
	@echo " Setup complete!"
	@echo ""
	@echo "Next steps:"
	@echo "1. Update the PASETO keys in $(COMPOSE_ENV_FILE)"
	@echo "2. Run 'make docker-up' to start services"
	@echo "3. Run 'make migrate-up service=identity' to apply migrations"

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

get_migrate_db_url_var_name = $(shell echo $(1) | tr '[:lower:]' '[:upper:]')_DB_URL_MIGRATE
derive_migrate_url = $(shell echo $(1) | sed -e 's/^postgres:\/\//pgx5:\/\//')

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
	@$(eval BASE_DB_URL_VAR := $(call get_db_url_var_name,$(SERVICE_NAME)))
	@$(eval MIGRATE_DB_URL_VAR := $(call get_migrate_db_url_var_name,$(SERVICE_NAME)))
	@$(if $(MIGRATION_DIR),,$(error Migration directory not found for $(SERVICE_NAME)))
	@$(if $(value $(BASE_DB_URL_VAR)),,$(error Environment variable $(BASE_DB_URL_VAR) is not set))
	@$(eval EFFECTIVE_MIGRATE_URL := $(or $(value $(MIGRATE_DB_URL_VAR)),$(call derive_migrate_url,$(value $(BASE_DB_URL_VAR)))))
	@echo "Applying UP migrations for service '$(SERVICE_NAME)' from $(MIGRATION_DIR)..."
	@$(DOCKER_COMPOSE) -f $(COMPOSE_FILE) --env-file $(COMPOSE_ENV_FILE) run --rm $(SERVICE_NAME)_migrate -path /migrations -database '$(EFFECTIVE_MIGRATE_URL)' up $(steps)

.PHONY: migrate-down
migrate-down: ## Roll back migrations (e.g., make migrate-down service=messaging steps=1)
	@$(eval SERVICE_NAME := $(if $(service),$(service),$(DEFAULT_SERVICE)))
	@$(eval MIGRATION_DIR := $(call get_migration_path,$(SERVICE_NAME)))
	@$(eval BASE_DB_URL_VAR := $(call get_db_url_var_name,$(SERVICE_NAME)))
	@$(eval MIGRATE_DB_URL_VAR := $(call get_migrate_db_url_var_name,$(SERVICE_NAME)))
	@$(if $(MIGRATION_DIR),,$(error Migration directory not found for $(SERVICE_NAME)))
	@$(if $(value $(BASE_DB_URL_VAR)),,$(error Environment variable $(BASE_DB_URL_VAR) is not set))
	@$(eval STEPS_ARG := $(if $(steps),$(steps),-all))
	@$(eval EFFECTIVE_MIGRATE_URL := $(or $(value $(MIGRATE_DB_URL_VAR)),$(call derive_migrate_url,$(value $(BASE_DB_URL_VAR)))))
	@echo "Rolling back $(STEPS_ARG) DOWN migration(s) for service '$(SERVICE_NAME)' from $(MIGRATION_DIR)..."
	@read -p "WARNING: Rolling back migrations. Are you sure? (y/N) " answer && [ $${answer:-N} = y ] || (echo "Aborted." && exit 1)
	@$(DOCKER_COMPOSE) -f $(COMPOSE_FILE) --env-file $(COMPOSE_ENV_FILE) run --rm $(SERVICE_NAME)_migrate -path /migrations -database '$(EFFECTIVE_MIGRATE_URL)' down $(STEPS_ARG)

.PHONY: migrate-status
migrate-status: ## Check migration status (e.g., make migrate-status service=messaging)
	@$(eval SERVICE_NAME := $(if $(service),$(service),$(DEFAULT_SERVICE)))
	@$(eval MIGRATION_DIR := $(call get_migration_path,$(SERVICE_NAME)))
	@$(eval BASE_DB_URL_VAR := $(call get_db_url_var_name,$(SERVICE_NAME)))
	@$(eval MIGRATE_DB_URL_VAR := $(call get_migrate_db_url_var_name,$(SERVICE_NAME)))
	@$(if $(MIGRATION_DIR),,$(error Migration directory not found for $(SERVICE_NAME)))
	@$(if $(value $(BASE_DB_URL_VAR)),,$(error Environment variable $(BASE_DB_URL_VAR) is not set))
	@$(eval EFFECTIVE_MIGRATE_URL := $(or $(value $(MIGRATE_DB_URL_VAR)),$(call derive_migrate_url,$(value $(BASE_DB_URL_VAR)))))
	@echo "Checking migration status for service '$(SERVICE_NAME)'..."
	@$(DOCKER_COMPOSE) -f $(COMPOSE_FILE) --env-file $(COMPOSE_ENV_FILE) run --rm $(SERVICE_NAME)_migrate -path /migrations -database '$(EFFECTIVE_MIGRATE_URL)' version

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
