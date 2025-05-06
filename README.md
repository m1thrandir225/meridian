# Meridian (Slack Integration)

## Prerequisites

- Docker
- Docker Compose
- Go 1.24+ (if you watn to run it without docker)

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
