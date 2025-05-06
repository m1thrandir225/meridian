# Meridian (Slack Integration)

## Prerequisites

- Docker
- Docker Compose
- Go 1.24+ (for local development)

## Structure of the project

Each entry point of the services is defined under `/cmd/{service}/main.go`.

All of the internal logic for each service is under `/internal/{service}/`

All of the `Dockerfile` and `compose.yml` configurations for each service are under `/deployments/{service}`

Each service has it's own database, along with it's own migrations, to make it easier to run migrations each service has a seperate container called `{service}_migrate` which runs the CLI-tool `golang-migrate` which makes sure the database is up to date when running the containers.

Each entity has public getters and private setters except the Aggregate Root.

Common shared packages and tools are placed under the `pkg` folder.

## Running The Project

1. Clone the repository

2. Start the services using the `Makefile`

The makefile is quite configured check out all of the commands using:

```bash
make help
```

To run all of the currently defined services:

```bash
make docker-up
```

3. To stop the running containers

To stop the containers:

```bash
make docker-stop
```

or to stop and remove the containers:

```
make docker-down
```

## Environment Setup

By default when running the `Makefile` it will generate a basic `.env` file by
default. If you want to configure it there is an example of what the project
needs under `/deployments/.env.example`.

Here are the required variables by the services:

### Messaging Service Environment Variables

- `MESSAGING_POSTGReES_USER`: username for the PostgreSQL instance of the messaging service
- `MESSAGING_POSTGRES_PASSWORD`: password for the PostgreSQL instance of the messaging service
- `MESSAGING_DB_NAME`: name of the database for the messaging service
- `MESSAGING_DB_PORT`: port for the PostgreSQL instance of the messaging service
- `MESSAGING_REDIS_PORT`: port for the Redis instace of the messaging service
- `MESSAGING_HTTP_PORT`: The http port the messaging service is going to be running on
- `MESSAGING_DB_URL`: The url that the messaging service is going to use to connect to the PostgreSQL instance (currently using the PGX driver)
- `MESSAGING_KAFKA_BROKERS`: a list of kafka brokers for the messaging service
- `MESSAGING_KAFKA_DEFAULT_TOPIC`: the default topic for the messaging service

## API Documentation

Here are the exposed API endpoints by the services:

### Messaging Service

#### CreateChannel

- **Endpont**: `/api/v1/channels`
- **Method**: `POST`
- **Description**: Create a channel and become the owner of it
- **Body**:
  ```json
  {
    "name": "string",
    "topic": "string",
    "creator_user_id": "string" //The UUID of the user from the identity service
  }
  ```
- **Response Status**: `201 Created`
- **Response Content**:
  ```json
  {
    "id": "string",
    "name": "string",
    "topic": "string",
    "creator_user_id": "string", //UUID of the user that created the channel
    "creation_time": "string", //ISO 8601 timestamp
    "last_message_time": "string", //ISO 8601 timestamp
    "is_archived": "boolean",
    "members_count": "integer" //Number of members in the channel
  }
  ```

#### JoinChannel

- **Endpont**: `/api/v1/channels/:channelId/join`
- **Method**: `POST`
- **Description**: Join a channel
- **Body**:
  ```json
  {
    "user_id": "string" //The UUID of the user from the identity service
  }
  ```
- **Response Status**: `204 No Content`

#### SendMessage

- **Endpont**: `/api/v1/channels/:channelId/messages`
- **Method**: `POST`
- **Description**: Send a message to a channel
- **Body**:
  ```json
  {
    "content_text": "string",
    "is_integration_message": "boolean", //is the message sent from an integration service
    "sender_user_id": "string", //The UUID of the user or integration id
    "parent_message_id": "string" //If the message is a reply the ID of the parent message
  }
  ```
- **Response Status**: `200 OK`
- **Response Content**:
  ```json
  {
    "id": "string",
    "channel_id": "string",
    "sender_user_id": "string",
    "integration_id": "string",
    "content_text": "string",
    "timestamp": "string", //ISO 8601 timestamp
    "parent_message_id": "string" //If the message is a reply the ID of the parent message
  }
  ```

#### GetChannel

- **Endpont**: `/api/v1/channels/:channelId`
- **Method**: `GET`
- **Description**: Get a channels details by their ID
- **Response Status**: `200 OK`
- **Response Content**:
  ```json
  {
    "id": "string",
    "name": "string",
    "topic": "string",
    "creator_user_id": "string", //UUID of the user that created the channel
    "creation_time": "string", //ISO 8601 timestamp
    "last_message_time": "string", //ISO 8601 timestamp
    "is_archived": "boolean",
    "members_count": "integer" //Number of members in the channel
  }
  ```

#### ArchiveChannel

- **Endpont**: `/api/v1/channels/:channelId/archive`
- **Method**: `PUT`
- **Description**: Archive a channel
- **Response Status**: `204 No Content`

#### Unarchive Channel

- **Endpont**: `/api/v1/channels/:channelId/unarchive`
- **Method**: `PUT`
- **Description**: Unarchive a channel
- **Response Status**: `204 No Content`

#### List Messages

- **Endpont**: `/api/v1/channels/:channelId/messages`
- **Method**: `Get`
- **Description**: List all messages for the given channelId
- **Response Status**: `200 OK`
- **Response Body**:

  ```json
  {
      "messages": {
        "id": "string",
        "channel_id": "string",
        "sender_user_id": "string",
        "integration_id": "string",
        "content_text": "string",
        "timestamp": "string", //ISO 8601 timestamp
        "parent_message_id": "string" //If the message is a reply the ID of the parent message
      }[]
  }
  ```

#### AddReaction

- **Endpont**: `/api/v1/channels/:channelId/messages/:messageId/reactions`
- **Method**: `POST`
- **Description**: Add a reaction to a message
- **Body**:
  ```json
  {
    "user_id": "string",
    "reaction_type": "string"
  }
  ```
- **Response Status**: `200 OK`
- **Response Content**:
  ```json
  {
    "id": "string",
    "message_id": "string",
    "user_id": "string",
    "reaction_type": "sting",
    "timestamp": "string"
  }
  ```

#### RemoveReaction

- **Endpont**: `/api/v1/channels/:channelId/messages/:messageId/reactions`
- **Method**: `DELETE`
- **Description**: Remove a reaction from a message
- **Body**:
  ```json
  {
    "user_id": "string",
    "reaction_type": "string"
  }
  ```
- **Response Status**: `204 No Content`
