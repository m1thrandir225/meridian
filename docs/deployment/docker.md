# Docker Deployment Guide

## Overview

Meridian is designed to run seamlessly with Docker and Docker Compose, providing a complete containerized environment for development and testing.

## Prerequisites

- **Docker**: Version 20.10 or higher
- **Docker Compose**: Version 2.0 or higher
- **System Requirements**: 8GB+ RAM, 20GB+ disk space
- **Network**: Ports 80, 443, 8080-8084, 5432-5435, 6379-6381, 9092, 2181

## Quick Start

### 1. Clone Repository

```bash
git clone https://github.com/your-org/meridian.git
cd meridian
```

### 2. Setup Environment

```bash
make setup
```

This command will:

- Generate PASETO security keys
- Create a comprehensive `.env` file with all service configurations
- Set up the environment for Docker deployment

### 3. Update Security Keys

Edit `deployments/.env` and replace the placeholder PASETO keys with the generated ones:

```bash
# Replace these placeholder values with the generated keys
IDENTITY_PASETO_PRIVATE_KEY=YOUR_PRIVATE_KEY_HERE
IDENTITY_PASETO_PUBLIC_KEY=YOUR_PUBLIC_KEY_HERE
```

### 4. Build and Start Services

```bash
make docker-build
make docker-up
```

Note: When using Docker the migrations are all ran automatically as sidecart containers.

## Available Commands

### Setup Commands

| Command              | Description                                |
| -------------------- | ------------------------------------------ |
| `make setup`         | Complete setup: generate keys and env file |
| `make generate-keys` | Generate PASETO security keys              |

### Build Commands

| Command                         | Description                   |
| ------------------------------- | ----------------------------- |
| `make build`                    | Build all service binaries    |
| `make build-one service=<name>` | Build specific service binary |

### Docker Commands

| Command             | Description                    |
| ------------------- | ------------------------------ |
| `make docker-build` | Build Docker images            |
| `make docker-up`    | Start all services             |
| `make docker-down`  | Stop and remove services       |
| `make docker-stop`  | Stop services without removing |
| `make docker-logs`  | Follow service logs            |
| `make docker-ps`    | List running services          |

### Database Commands

| Command                                               | Description                  |
| ----------------------------------------------------- | ---------------------------- |
| `make migrate-up service=<name>`                      | Apply migrations for service |
| `make migrate-down service=<name>`                    | Rollback migrations          |
| `make migrate-status service=<name>`                  | Check migration status       |
| `make migrate-create service=<name> name=<migration>` | Create new migration         |

### Development Commands

| Command      | Description            |
| ------------ | ---------------------- |
| `make tidy`  | Tidy Go module files   |
| `make fmt`   | Format Go source code  |
| `make lint`  | Lint Go source code    |
| `make test`  | Run Go tests           |
| `make clean` | Remove build artifacts |

## Service Architecture

### Core Services

- **Identity Service** (Port 8080): User authentication and management
- **Messaging Service** (Port 8081): Real-time messaging and channels
- **Integration Service** (Port 8082): Third-party integrations and webhooks
- **Analytics Service** (Port 8083): Usage analytics and metrics

### Infrastructure Services

- **PostgreSQL**: Primary database for all services
- **Redis**: Caching and session storage
- **Apache Kafka**: Event streaming and message queue
- **Zookeeper**: Kafka coordination service

### Frontend Services

- **Frontend** (Port 3000): Main web application
- **Landing** (Port 3001): Marketing landing page

### Reverse Proxy

- **Traefik**: Load balancer and reverse proxy with automatic SSL
  - Configuration: `deployments/traefik/traefik.yml`
  - Dynamic configuration: `deployments/traefik/dynamic.yml`
  - Service discovery: Automatic service discovery through Docker labels

## Configuration

### Environment Variables

The `make setup` command generates a comprehensive `.env` file with all necessary configurations:

```bash
# Global Environment Settings
ENVIRONMENT=development
LOG_LEVEL=info

# Identity Service Configuration
IDENTITY_POSTGRES_USER=root
IDENTITY_POSTGRES_PASSWORD=secret
IDENTITY_DB_NAME=identity_db
IDENTITY_DB_PORT=5432
IDENTITY_REDIS_PORT=6379
IDENTITY_HTTP_PORT=8080
IDENTITY_GRPC_PORT=9090
IDENTITY_DB_URL=postgres://root:secret@identity_postgres:5432/identity_db?sslmode=disable
IDENTITY_REDIS_URL=redis://identity_redis:6379
IDENTITY_KAFKA_BROKERS=kafka:9092
IDENTITY_KAFKA_DEFAULT_TOPIC=meridian.identity.events
IDENTITY_PASETO_PRIVATE_KEY=YOUR_PRIVATE_KEY_HERE
IDENTITY_PASETO_PUBLIC_KEY=YOUR_PUBLIC_KEY_HERE

# Similar configurations for other services...
```

### Service-Specific Configurations

Each service has its own database, Redis instance, and Kafka topic for isolation and scalability.

#### Traefik Service Discovery

Traefik automatically discovers services through Docker labels and provides:

- **Automatic SSL**: Let's Encrypt certificate generation
- **Load Balancing**: Round-robin load balancing across service instances
- **Health Checks**: Automatic health check monitoring
- **Service Routing**: Path-based routing to different services

Configuration files:

- `deployments/traefik/traefik.yml` - Main Traefik configuration
- `deployments/traefik/dynamic.yml` - Dynamic configuration for middleware and TLS

## Database Management

### Migration Commands

```bash
# Apply migrations for Identity service
make migrate-up service=identity

# Apply migrations for Messaging service
make migrate-up service=messaging

# Apply migrations for Integration service
make migrate-up service=integration

# Apply migrations for Analytics service
make migrate-up service=analytics
```

### Creating New Migrations

```bash
# Create migration for Identity service
make migrate-create service=identity name=add_user_preferences

# Create migration for Messaging service
make migrate-create service=messaging name=add_message_attachments
```

**Note** When running through docker all of the migrations are automatically run on startup, as sidecart containers.

## Monitoring and Logs

### View Service Logs

```bash
# Follow all service logs
make docker-logs

# Follow specific service logs
make docker-logs service=identity
make docker-logs service=messaging
make docker-logs service=integration
make docker-logs service=analytics
```

### Health Checks

Each service provides health check endpoints:

- Identity: `http://localhost:8080/health`
- Messaging: `http://localhost:8081/health`
- Integration: `http://localhost:8082/health`
- Analytics: `http://localhost:8083/health`

### Metrics Endpoints

- Identity: `http://localhost:8080/metrics`
- Messaging: `http://localhost:8081/metrics`
- Integration: `http://localhost:8082/metrics`
- Analytics: `http://localhost:8083/metrics`

## Troubleshooting

### Common Issues

#### Port Conflicts

If you encounter port conflicts, check what's running on the required ports:

```bash
# Check ports in use
netstat -tulpn | grep :8080
netstat -tulpn | grep :8081
netstat -tulpn | grep :8082
netstat -tulpn | grep :8083
```

#### Database Connection Issues

```bash
# Check database containers
make docker-ps

# View database logs
make docker-logs service=identity_postgres
```

#### Service Startup Issues

```bash
# Check service status
make docker-ps

# View detailed logs
make docker-logs service=identity
```

### Reset Environment

To completely reset the environment:

```bash
# Stop and remove all containers and volumes
make docker-down options=-v

# Clean build artifacts
make clean

# Rebuild and restart
make build
make docker-up
```

## Development Workflow

### Typical Development Session

```bash
# 1. Start the day
make docker-up

# 2. View logs for debugging
make docker-logs

# 3. Make code changes and rebuild
make docker-build service=identity
make docker-up service=identity

# 4. End the day
make docker-down
```

### Hot Reloading

TODO: Currently no service supports Hot Reloading, this is a plan for the future.

## Security Considerations

### PASETO Keys

- Generate unique keys for each environment
- Never commit keys to version control
- Rotate keys regularly in production

### Database Security

- Use strong passwords in production
- Enable SSL connections
- Restrict network access

### Network Security

- Use reverse proxy for SSL termination
- Implement proper firewall rules
- Monitor network traffic

## Performance Tuning

### Resource Allocation

Adjust Docker Compose resource limits based on your system:

```yaml
services:
  identity:
    deploy:
      resources:
        limits:
          memory: 512M
          cpus: "0.5"
```

### Database Optimization

- Monitor query performance
- Add appropriate indexes
- Configure connection pooling

### Caching Strategy

- Use Redis for session storage
- Implement application-level caching
- Monitor cache hit rates
