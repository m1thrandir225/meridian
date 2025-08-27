# REST API Documentation

## Overview

Meridian provides REST APIs across all services for user management, messaging, integrations, and analytics. All APIs follow RESTful principles with consistent request/response patterns.

## Base URLs

| Service     | Base URL                                    | Port |
| ----------- | ------------------------------------------- | ---- |
| Identity    | `http://localhost:8080/api/v1/auth`         | 8080 |
| Messaging   | `http://localhost:8081/api/v1/messages`     | 8081 |
| Integration | `http://localhost:8082/api/v1/integrations` | 8082 |
| Analytics   | `http://localhost:8083/api/v1/analytics`    | 8083 |

## Authentication

### Bearer Token Authentication

Most endpoints require authentication using Bearer tokens:

```http
Authorization: Bearer v4.local.xxx...
```

### User ID Header

Some endpoints require the authenticated user's ID:

```http
X-User-ID: 01234567-89ab-cdef-0123-456789abcdef
```

## Common Response Patterns

### Success Response

```json
{
  "data": {
    // Response payload
  }
}
```

### Error Response

```json
{
  "error": "Error message description"
}
```

## HTTP Status Codes

| Code | Description           | Usage                                |
| ---- | --------------------- | ------------------------------------ |
| 200  | OK                    | Successful GET, PUT requests         |
| 201  | Created               | Successful POST requests             |
| 202  | Accepted              | Successful DELETE requests           |
| 400  | Bad Request           | Invalid request format or parameters |
| 401  | Unauthorized          | Missing or invalid authentication    |
| 403  | Forbidden             | Insufficient permissions             |
| 404  | Not Found             | Resource doesn't exist               |
| 409  | Conflict              | Resource already exists or conflict  |
| 500  | Internal Server Error | Server-side errors                   |

---

## Identity Service API

### Authentication Endpoints

#### Register User

```http
POST /api/v1/auth/register
Content-Type: application/json

{
  "username": "johndoe",
  "email": "john@example.com",
  "firstName": "John",
  "lastName": "Doe",
  "password": "SecurePassword123!"
}
```

**Response (200):**

```json
{
  "id": "01234567-89ab-cdef-0123-456789abcdef",
  "username": "johndoe",
  "email": "john@example.com",
  "firstName": "John",
  "lastName": "Doe",
  "isAdmin": false
}
```

#### User Login

```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "login": "johndoe",
  "password": "SecurePassword123!"
}
```

**Response (200):**

```json
{
  "user": {
    "id": "01234567-89ab-cdef-0123-456789abcdef",
    "username": "johndoe",
    "email": "john@example.com",
    "firstName": "John",
    "lastName": "Doe",
    "isAdmin": false
  },
  "tokens": {
    "accessToken": "v4.local.xxx...",
    "refreshToken": "v4.local.yyy...",
    "tokenType": "bearer",
    "expiresIn": 3600
  }
}
```

#### Get Current User

```http
GET /api/v1/auth/me
Authorization: Bearer v4.local.xxx...
```

**Response (200):**

```json
{
  "id": "01234567-89ab-cdef-0123-456789abcdef",
  "username": "johndoe",
  "email": "john@example.com",
  "firstName": "John",
  "lastName": "Doe",
  "isAdmin": false
}
```

---

## Messaging Service API

### Channel Management

#### Get User Channels

```http
GET /api/v1/messages/channels/
Authorization: Bearer v4.local.xxx...
X-User-ID: 01234567-89ab-cdef-0123-456789abcdef
```

**Response (200):**

```json
[
  {
    "id": "11234567-89ab-cdef-0123-456789abcdef",
    "name": "general",
    "topic": "General discussion",
    "creatorUserId": "01234567-89ab-cdef-0123-456789abcdef",
    "creationTime": "2024-01-01T10:00:00Z",
    "lastMessageTime": "2024-01-15T14:30:00Z",
    "memberCount": 25,
    "isArchived": false
  }
]
```

#### Create Channel

```http
POST /api/v1/messages/channels/
Authorization: Bearer v4.local.xxx...
X-User-ID: 01234567-89ab-cdef-0123-456789abcdef
Content-Type: application/json

{
  "name": "development",
  "topic": "Development discussions and updates"
}
```

**Response (201):**

```json
{
  "id": "21234567-89ab-cdef-0123-456789abcdef",
  "name": "development",
  "topic": "Development discussions and updates",
  "creatorUserId": "01234567-89ab-cdef-0123-456789abcdef",
  "creationTime": "2024-01-15T15:00:00Z",
  "lastMessageTime": "2024-01-15T15:00:00Z",
  "memberCount": 1,
  "isArchived": false
}
```

#### Send Message

```http
POST /api/v1/messages/channels/11234567-89ab-cdef-0123-456789abcdef/messages
Authorization: Bearer v4.local.xxx...
X-User-ID: 01234567-89ab-cdef-0123-456789abcdef
Content-Type: application/json

{
  "content_text": "Hello, everyone!",
  "is_integration_message": false,
  "parent_message_id": null
}
```

**Response (201):**

```json
{
  "id": "31234567-89ab-cdef-0123-456789abcdef",
  "channelId": "11234567-89ab-cdef-0123-456789abcdef",
  "senderUserId": "01234567-89ab-cdef-0123-456789abcdef",
  "content": {
    "text": "Hello, everyone!",
    "type": "text"
  },
  "createdAt": "2024-01-15T14:30:00Z",
  "reactions": [],
  "user": {
    "username": "johndoe",
    "firstName": "John",
    "lastName": "Doe"
  }
}
```

---

## Integration Service API

### Integration Management

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

---

## Analytics Service API

**Note**: All analytics endpoints require admin authentication.

### Dashboard Analytics

#### Get Dashboard Data

```http
GET /api/v1/analytics/dashboard?timeRange=7d
Authorization: Bearer <admin-token>
```

**Response (200):**

```json
{
  "total_users": 1250,
  "active_users": 450,
  "new_users_today": 15,
  "messages_today": 2847,
  "total_channels": 85,
  "active_channels": 42,
  "average_messages_per_user": 6.3,
  "peak_hour": 14,
  "last_updated": "2024-01-15T17:30:00Z"
}
```

#### Get User Growth

```http
GET /api/v1/analytics/user-growth?startDate=2024-01-01&endDate=2024-01-15&interval=daily
Authorization: Bearer <admin-token>
```

**Response (200):**

```json
[
  {
    "period": "2024-01-01",
    "new_users": 25,
    "total_users": 1200,
    "growth_rate": 2.13
  },
  {
    "period": "2024-01-02",
    "new_users": 18,
    "total_users": 1218,
    "growth_rate": 1.5
  }
]
```

## Error Handling

### Common Error Scenarios

#### Validation Error (400)

```json
{
  "error": "Validation failed",
  "details": {
    "email": "Email format is invalid",
    "password": "Password must be at least 8 characters"
  }
}
```

#### Authentication Error (401)

```json
{
  "error": "Invalid or expired token"
}
```

#### Permission Error (403)

```json
{
  "error": "Insufficient permissions to access this resource"
}
```

#### Not Found Error (404)

```json
{
  "error": "Channel not found"
}
```

#### Conflict Error (409)

```json
{
  "error": "Username is already taken"
}
```

## WebSocket API

WebSocket connections are available at:

```
ws://localhost:8081/api/v1/messages/ws
```

**Note**: WebSocket connections are only available to the frontend application for real-time messaging. External integrations should use HTTP REST APIs for sending messages.

See individual service documentation for detailed WebSocket protocol information.
