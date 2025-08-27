# Integration Service Documentation

## Overview

The Integration Service manages third-party service connections, webhooks, and bot integrations. It provides secure API token management and enables external services to send messages to channels through webhook endpoints.

## Domain Model

### Bounded Context: Third-party Integrations

The Integration service manages third-party service connections and API tokens within its own bounded context.

### Aggregate Root: Integration

The `Integration` aggregate manages third-party service connections and API tokens. It encapsulates integration data and behavior, ensuring consistency and business rules.

```go
type Integration struct {
    ID               IntegrationID
    ServiceName      string
    CreatorUserID    UserIDRef
    HashedAPIToken   APIToken
    TokenLookupHash  string
    CreatedAt        time.Time
    IsRevoked        bool
    TargetChannelIDs []ChannelIDRef
    Version          int64
}
```

### Value Objects

| Value Object    | Purpose                                     | Implementation                  |
| --------------- | ------------------------------------------- | ------------------------------- |
| `IntegrationID` | Unique integration identifier               | UUID v7 format                  |
| `APIToken`      | Secure API token storage                    | Hashed with lookup optimization |
| `UserIDRef`     | Reference to user from Identity service     | String UUID reference           |
| `ChannelIDRef`  | Reference to channel from Messaging service | String UUID reference           |

### Domain Events

- `IntegrationRegistered` - New integration created
- `APITokenRevoked` - Token revoked
- `APITokenUpvoked` - Token reactivated
- `IntegrationTargetChannelsUpdated` - Target channels modified

### Commands

- `RegisterIntegration` - Register new integration
- `RevokeToken` - Revoke API token
- `UpvokeIntegration` - Reactivate integration
- `UpdateIntegration` - Update integration settings

## API Reference

### REST API Endpoints

Base URL: `/api/v1/integrations`

#### Integration Management

| Method | Endpoint          | Description              | Auth Required |
| ------ | ----------------- | ------------------------ | ------------- |
| POST   | `/`               | Register new integration | Yes           |
| GET    | `/`               | List user's integrations | Yes           |
| GET    | `/validate-token` | Validate API token       | No            |
| PUT    | `/:id/revoke`     | Revoke integration       | Yes           |
| PUT    | `/:id/upvoke`     | Reactivate integration   | Yes           |
| PUT    | `/:id`            | Update integration       | Yes           |

#### Webhook Endpoints

| Method | Endpoint            | Description              | Auth Required |
| ------ | ------------------- | ------------------------ | ------------- |
| POST   | `/webhook/message`  | Send message via webhook | API Token     |
| POST   | `/webhook/callback` | Handle callback messages | API Token     |

### Request/Response Examples

#### Register Integration

```http
POST /api/v1/integrations
Authorization: Bearer v4.local.xxx...
X-User-ID: 01234567-89ab-cdef-0123-456789abcdef
Content-Type: application/json

{
  "service_name": "GitHub Bot",
  "target_channel_ids": [
    "11234567-89ab-cdef-0123-456789abcdef",
    "21234567-89ab-cdef-0123-456789abcdef"
  ]
}
```

**Response (200):**

```json
{
  "id": "71234567-89ab-cdef-0123-456789abcdef",
  "serviceName": "GitHub Bot",
  "creatorUserId": "01234567-89ab-cdef-0123-456789abcdef",
  "apiToken": "mrd_1234567890abcdef...",
  "targetChannelIds": [
    "11234567-89ab-cdef-0123-456789abcdef",
    "21234567-89ab-cdef-0123-456789abcdef"
  ],
  "createdAt": "2024-01-15T15:30:00Z",
  "isRevoked": false
}
```

#### List Integrations

```http
GET /api/v1/integrations
Authorization: Bearer v4.local.xxx...
X-User-ID: 01234567-89ab-cdef-0123-456789abcdef
```

**Response (200):**

```json
[
  {
    "id": "71234567-89ab-cdef-0123-456789abcdef",
    "serviceName": "GitHub Bot",
    "creatorUserId": "01234567-89ab-cdef-0123-456789abcdef",
    "targetChannelIds": [
      "11234567-89ab-cdef-0123-456789abcdef",
      "21234567-89ab-cdef-0123-456789abcdef"
    ],
    "createdAt": "2024-01-15T15:30:00Z",
    "isRevoked": false
  }
]
```

#### Send Webhook Message

```http
POST /api/v1/integrations/webhook/message
Authorization: Bearer mrd_1234567890abcdef...
Content-Type: application/json

{
  "content_text": "🚀 New deployment successful!",
  "target_channel_id": "11234567-89ab-cdef-0123-456789abcdef",
  "metadata": {
    "source": "GitHub Actions",
    "deployment_id": "dep_123456"
  }
}
```

**Response (201):**

```json
{
  "messageId": "81234567-89ab-cdef-0123-456789abcdef",
  "channelId": "11234567-89ab-cdef-0123-456789abcdef",
  "sent": true,
  "timestamp": "2024-01-15T17:00:00Z"
}
```

#### Validate API Token

```http
GET /api/v1/integrations/validate-token
Authorization: Bearer mrd_1234567890abcdef...
```

**Response (200):**

```json
{
  "valid": true,
  "integrationId": "71234567-89ab-cdef-0123-456789abcdef",
  "serviceName": "GitHub Bot",
  "targetChannelIds": [
    "11234567-89ab-cdef-0123-456789abcdef",
    "21234567-89ab-cdef-0123-456789abcdef"
  ]
}
```

#### Revoke Integration

```http
PUT /api/v1/integrations/71234567-89ab-cdef-0123-456789abcdef/revoke
Authorization: Bearer v4.local.xxx...
X-User-ID: 01234567-89ab-cdef-0123-456789abcdef
```

**Response (202):** No content

#### Update Integration

```http
PUT /api/v1/integrations/71234567-89ab-cdef-0123-456789abcdef
Authorization: Bearer v4.local.xxx...
X-User-ID: 01234567-89ab-cdef-0123-456789abcdef
Content-Type: application/json

{
  "service_name": "Updated GitHub Bot",
  "target_channel_ids": [
    "11234567-89ab-cdef-0123-456789abcdef"
  ]
}
```

**Response (200):**

```json
{
  "id": "71234567-89ab-cdef-0123-456789abcdef",
  "serviceName": "Updated GitHub Bot",
  "creatorUserId": "01234567-89ab-cdef-0123-456789abcdef",
  "targetChannelIds": ["11234567-89ab-cdef-0123-456789abcdef"],
  "createdAt": "2024-01-15T15:30:00Z",
  "isRevoked": false
}
```

### gRPC Services

**Note**: gRPC services are for internal inter-service communication only. External integrations should use HTTP REST APIs.

#### IntegrationService

```protobuf
service IntegrationService {
  rpc RegisterBot(RegisterBotRequest) returns (RegisterBotResponse);
  rpc ValidateIntegrationToken(ValidateTokenRequest) returns (ValidateTokenResponse);
}
```

## Domain Logic

### Integration Registration

1. Validate service name and target channels
2. Generate secure API token with prefix
3. Hash token for secure storage
4. Create Integration aggregate with domain events
5. Persist to database
6. Publish `IntegrationRegistered` event

### Webhook Message Processing

1. Validate API token and permissions
2. Check target channel access
3. Forward message to Messaging service
4. Create domain events
5. Publish `MessageSent` event

### Token Security

- **Token Format**: `mrd_` prefix followed by secure random string
- **Storage**: Hashed with bcrypt for security
- **Lookup**: Optimized hash for quick token validation
- **Revocation**: Soft delete with timestamp tracking

## Events

### Published Events

#### IntegrationRegisteredEvent

```json
{
  "eventType": "IntegrationRegistered",
  "aggregateId": "71234567-89ab-cdef-0123-456789abcdef",
  "version": 1,
  "serviceName": "GitHub Bot",
  "creatorUserID": "01234567-89ab-cdef-0123-456789abcdef",
  "targetChannelIDs": [
    "11234567-89ab-cdef-0123-456789abcdef",
    "21234567-89ab-cdef-0123-456789abcdef"
  ],
  "timestamp": "2024-01-15T15:30:00Z"
}
```

#### APITokenRevokedEvent

```json
{
  "eventType": "APITokenRevoked",
  "aggregateId": "71234567-89ab-cdef-0123-456789abcdef",
  "version": 2,
  "integrationID": "71234567-89ab-cdef-0123-456789abcdef",
  "revokedAt": "2024-01-15T16:00:00Z"
}
```

## Infrastructure

### Technology Stack

- **Language**: Go 1.24+
- **Framework**: Gin (HTTP), gRPC (Inter-service)
- **Database**: PostgreSQL
- **Cache**: Redis
- **Message Queue**: Apache Kafka

### Ports

- **HTTP**: 8082
- **gRPC**: 9092

### Configuration

#### Environment Variables

| Variable                    | Description                  | Default | Required |
| --------------------------- | ---------------------------- | ------- | -------- |
| `INTEGRATION_HTTP_PORT`     | HTTP server port             | `:8082` | Yes      |
| `INTEGRATION_GRPC_PORT`     | gRPC server port             | `9092`  | Yes      |
| `INTEGRATION_DB_URL`        | PostgreSQL connection string | -       | Yes      |
| `INTEGRATION_REDIS_URL`     | Redis connection string      | -       | Yes      |
| `INTEGRATION_KAFKA_BROKERS` | Kafka broker addresses       | -       | Yes      |
| `MESSAGING_GRPC_URL`        | Messaging service gRPC URL   | -       | Yes      |

### Database Schema

#### Entity Relationship Diagram

```mermaid
erDiagram
    integrations {
        UUID id PK
        VARCHAR service_name
        VARCHAR creator_user_id
        VARCHAR api_token_hash
        VARCHAR token_lookup_hash UK
        TEXT[] target_channel_ids
        TIMESTAMPTZ created_at
        BOOLEAN is_revoked
        BIGINT version
    }
```

#### Integrations Table

```sql
CREATE TABLE integrations (
    id UUID PRIMARY KEY,
    service_name VARCHAR(255) NOT NULL,
    creator_user_id VARCHAR(255) NOT NULL,
    api_token_hash VARCHAR(255) NOT NULL,
    token_lookup_hash VARCHAR(255) NOT NULL UNIQUE,
    target_channel_ids TEXT[] NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    is_revoked BOOLEAN NOT NULL DEFAULT FALSE,
    version BIGINT NOT NULL DEFAULT 1
);

-- Indexes
CREATE INDEX idx_integrations_creator ON integrations(creator_user_id);
CREATE INDEX idx_integrations_lookup_hash ON integrations(token_lookup_hash);
```

## Error Handling

### Common Error Codes

| Code                      | Description                  | HTTP Status |
| ------------------------- | ---------------------------- | ----------- |
| `ErrIntegrationNotFound`  | Integration doesn't exist    | 404         |
| `ErrInvalidAPIToken`      | Invalid or expired API token | 401         |
| `ErrIntegrationRevoked`   | Integration has been revoked | 403         |
| `ErrChannelAccessDenied`  | No access to target channel  | 403         |
| `ErrInvalidServiceName`   | Invalid service name format  | 400         |
| `ErrDuplicateIntegration` | Integration already exists   | 409         |

### Error Response Format

```json
{
  "error": "Invalid API token"
}
```

## Monitoring & Health Checks

### Health Endpoint

```http
GET /health
```

Response:

```json
{
  "status": "healthy",
  "service": "integration"
}
```

### Metrics

- Active integrations count
- Webhook message throughput
- Token validation rates
- Error rates by endpoint
- Database query performance
