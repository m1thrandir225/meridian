# Meridian Services Overview

## Table of Contents

- [Architecture Overview](#architecture-overview)
- [Identity Service](#identity-service)
- [Messaging Service](#messaging-service)
- [Integration Service](#integration-service)
- [Analytics Service](#analytics-service)
- [Service Communication](#service-communication)
- [Infrastructure Dependencies](#infrastructure-dependencies)

## Architecture Overview

Meridian is a distributed microservices-based Slack-like chat application built using Domain Driven Design (DDD) principles. The system consists of four core services that communicate through Kafka event streaming and gRPC calls.

### Core Services

1. **Identity Service** - User authentication and management
2. **Messaging Service** - Real-time messaging and channels
3. **Integration Service** - Third-party integrations and webhooks
4. **Analytics Service** - Data analytics and metrics tracking

### Technology Stack

- **Language**: Go 1.24+
- **Framework**: Gin (HTTP), gRPC (Inter-service communication)
- **Database**: PostgreSQL (per service)
- **Cache**: Redis (per service)
- **Message Queue**: Apache Kafka
- **Frontend**: Vue.js 3 + TypeScript
- **Deployment**: Docker + Docker Compose
- **Reverse Proxy**: Traefik

---

## Identity Service

### Purpose

Manages user authentication, authorization, and user profile information for the entire platform.

### Core Features

- User registration and authentication
- JWT token generation and validation (PASETO tokens)
- User profile management
- Password management with secure hashing
- Admin user management

### Domain Model

- **Aggregate Root**: `User`
- **Value Objects**: `UserID`, `Username`, `Email`, `PasswordHash`
- **Entities**: `RefreshToken`

### API Endpoints

**Base URL**: `/api/v1/auth`

| Method | Endpoint             | Description                   |
| ------ | -------------------- | ----------------------------- |
| POST   | `/register`          | Register a new user           |
| POST   | `/login`             | Authenticate user             |
| GET    | `/validate-token`    | Validate authentication token |
| POST   | `/refresh-token`     | Refresh authentication token  |
| GET    | `/me`                | Get current user profile      |
| PUT    | `/me/update-profile` | Update user profile           |
| PUT    | `/me/password`       | Change user password          |
| DELETE | `/me`                | Delete user account           |

### gRPC Services

- `ValidateToken` - Token validation for other services
- `GetUserByID` - User information retrieval
- `GetUsers` - Bulk user information retrieval

### Infrastructure

- **Database**: PostgreSQL
- **Cache**: Redis (session storage)
- **Message Queue**: Kafka (user events)
- **Ports**:
  - HTTP: 8080
  - gRPC: 9090

### Events Published

- `UserRegistered`
- `UserAuthenticated`
- `UserProfileUpdated`
- `UserPasswordUpdated`
- `UserDeleted`

---

## Messaging Service

### Purpose

Handles real-time messaging, channel management, and WebSocket connections for live chat functionality.

### Core Features

- Channel creation and management
- Real-time messaging with WebSocket support
- Message reactions and threading
- Channel invitations and member management
- Message persistence and history
- Bot integration support

### Domain Model

- **Aggregate Root**: `Channel`
- **Entities**: `Message`, `Member`, `ChannelInvite`, `Reaction`
- **Value Objects**: `MessageContent`

### API Endpoints

**Base URL**: `/api/v1/messages`

| Method | Endpoint                                  | Description                                  |
| ------ | ----------------------------------------- | -------------------------------------------- |
| GET    | `/ws`                                     | WebSocket connection for real-time messaging |
| GET    | `/channels/`                              | Get user's channels                          |
| POST   | `/channels/`                              | Create a new channel                         |
| GET    | `/channels/:id`                           | Get channel details                          |
| POST   | `/channels/:id/join`                      | Join a channel                               |
| PUT    | `/channels/:id/archive`                   | Archive a channel                            |
| PUT    | `/channels/:id/unarchive`                 | Unarchive a channel                          |
| POST   | `/channels/:id/bots`                      | Add bot to channel                           |
| POST   | `/channels/:id/invites`                   | Create channel invite                        |
| GET    | `/channels/:id/invites`                   | Get channel invites                          |
| GET    | `/channels/:id/messages`                  | Get channel messages                         |
| POST   | `/channels/:id/messages`                  | Send message to channel                      |
| PUT    | `/channels/:id/messages/:msgId/reactions` | Add reaction                                 |
| DELETE | `/channels/:id/messages/:msgId/reactions` | Remove reaction                              |
| POST   | `/invites/accept`                         | Accept channel invite                        |
| DELETE | `/invites/:id`                            | Deactivate channel invite                    |

### gRPC Services

- `SendMessage` - External message sending
- `RegisterBot` - Bot registration for integrations

### Real-time Features

- **WebSocket Support**: Real-time message delivery
- **Message Broadcasting**: Live updates to all channel members
- **Typing Indicators**: Real-time typing status
- **Connection Management**: Automatic reconnection handling

### Infrastructure

- **Database**: PostgreSQL
- **Cache**: Redis (real-time data)
- **Message Queue**: Kafka (message events)
- **WebSocket**: Real-time communication
- **Ports**:
  - HTTP: 8081
  - gRPC: 9091

### Events Published

- `ChannelCreated`
- `UserJoinedChannel`
- `UserLeftChannel`
- `MessageSent`
- `ReactionAdded`
- `ReactionRemoved`
- `ChannelArchived`
- `ChannelInviteCreated`

---

## Integration Service

### Purpose

Manages third-party integrations, webhooks, and external service connections.

### Core Features

- API token management for integrations
- Webhook endpoint handling
- Integration registration and management
- Secure token generation and validation
- Channel targeting for integrations

### Domain Model

- **Aggregate Root**: `Integration`
- **Value Objects**: `IntegrationID`, `APIToken`, `UserIDRef`, `ChannelIDRef`

### API Endpoints

**Base URL**: `/api/v1/integrations`

| Method | Endpoint            | Description                    |
| ------ | ------------------- | ------------------------------ |
| POST   | `/`                 | Register new integration       |
| POST   | `/upvoke`           | Reactivate revoked integration |
| DELETE | `/revoke`           | Revoke integration token       |
| PUT    | `/:id`              | Update integration             |
| DELETE | `/:id`              | Delete integration             |
| GET    | `/`                 | List user integrations         |
| POST   | `/webhook/message`  | Webhook message endpoint       |
| POST   | `/callback/message` | Callback message endpoint      |

### gRPC Services

- `ValidateAPIToken` - Token validation for external services
- `GetIntegration` - Integration information retrieval

### Security Features

- **Token Hashing**: Secure API token storage
- **Lookup Hash**: Fast token validation
- **Revocation System**: Token lifecycle management
- **Channel Targeting**: Scoped message delivery

### Infrastructure

- **Database**: PostgreSQL
- **Cache**: Redis (token caching)
- **Message Queue**: Kafka (integration events)
- **Ports**:
  - HTTP: 8082
  - gRPC: 9092

### Events Published

- `IntegrationRegistered`
- `APITokenRevoked`
- `APITokenUpvoked`
- `IntegrationTargetChannelsUpdated`

---

## Analytics Service

### Purpose

Collects, processes, and provides analytics data about user behavior, message patterns, and system usage.

### Core Features

- Real-time metrics tracking
- User activity monitoring
- Channel usage analytics
- Message volume analysis
- Growth metrics and reporting
- Event-driven data collection

### Domain Model

- **Aggregate Root**: `Analytics`
- **Entities**: `AnalyticsMetric`, `UserActivity`, `ChannelActivity`
- **Value Objects**: `AnalyticsID`, `MetricID`

### API Endpoints

**Base URL**: `/api/v1/analytics`

| Method | Endpoint            | Description              |
| ------ | ------------------- | ------------------------ |
| GET    | `/metrics`          | Get system metrics       |
| GET    | `/user-growth`      | User growth analytics    |
| GET    | `/message-volume`   | Message volume analytics |
| GET    | `/channel-activity` | Channel activity metrics |
| GET    | `/active-users`     | Active user statistics   |

### Analytics Capabilities

- **Real-time Tracking**: Live metric collection
- **Historical Analysis**: Time-series data aggregation
- **User Behavior**: Activity patterns and engagement
- **System Performance**: Usage statistics and trends
- **Custom Metrics**: Extensible metric system

### Event Processing

Listens to events from all services:

- User registration events
- Message sent events
- Channel activity events
- Reaction events
- Login/logout events

### Infrastructure

- **Database**: PostgreSQL (time-series optimized)
- **Cache**: Redis (real-time aggregation)
- **Message Queue**: Kafka (event consumption)
- **Ports**:
  - HTTP: 8083

### Events Consumed

- `UserRegistered`
- `MessageSent`
- `ChannelCreated`
- `UserJoinedChannel`
- `ReactionAdded`

---

## Service Communication

### Inter-service Communication

1. **Synchronous**: gRPC for real-time data needs
2. **Asynchronous**: Kafka events for eventual consistency
3. **Caching**: Redis for performance optimization

### Communication Patterns

- **Identity ↔ Messaging**: User validation via gRPC
- **Integration ↔ Messaging**: Message sending via gRPC
- **Analytics ← All Services**: Event consumption via Kafka
- **All Services**: Cross-cutting events via Kafka

### Event Flow
