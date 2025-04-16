# 🎟️ Book Microservice

**Book** is a sleek and efficient microservice built in Go, designed to handle movie ticket bookings in a modern distributed system.

It exposes a gRPC API for clients to submit booking requests, writes bookings to a SQLite datastore, and publishes booking events to Kafka—making it a great plug-and-play component in a cinema or ticketing platform.

Whether you're reserving one seat or a hundred, **Book** is got your back. It's fully containerized with Docker, supports config-driven environments, and includes a functional test suite to keep your booking logic bulletproof.

## Features

  - 🎬 **gRPC API** — Fast, typed, and scalable endpoint for booking movie tickets.
  - 💾 **SQLite Storage** — Lightweight, file-based persistence with transactional integrity.
  - 🧠 **Business Logic Layer** — Validates booking data.
  - 🧵 **Kafka Integration** — Publishes booking events to a Kafka topic for downstream consumers.
  - 🧪 **Functional Test Suite** — Covers end-to-end booking flows with full gRPC client testing.
  - ⚙️ **Configurable by Environment** — Load configs dynamically via env var `CONFIG_PATH`, supporting `local`, `dev`, `test`, `prod`, and `custom` setups.
  - 🐳 **Dockerized** — Easily build and run in isolated container, ready for deployment or testing.
  - 📜 **Migrations CLI** — Handy built-in migrator for applying SQLite schema migrations via a CLI command.
  - 🪵 **Structured Logging** — Context-rich logs using slog, configurable log modes (like `silent`, `local`, etc.).

## Prerequisites

  - **Go 1.21+** installed
  - **Docker + Docker Compose**
  - **make** installed
  - (Optional) Kafka running locally or remotely for event publishing

## Setup

### 1. Clone the Repository

```
git clone https://github.com/bookamovie/book
```

### 2. Environment Configuration

The service is controlled using two environment variables:

#### `CONFIG_PATH`

Specifies which config file to load. Available options:

| Value                 | Description             |
|-----------------------|-------------------------|
| `config/local.yaml`   | For local development   |
| `config/dev.yaml`     | For development server  |
| `config/prod.yaml`    | For production          |
| `config/custom.yaml`  | For custom setups       |

#### `LOG_MODE`

Defines the logging format:

| Value     | Description                        |
|-----------|------------------------------------|
| `local`   | Human-readable logs with colors    |
| `dev`     | For dev environment                |
| `prod`    | For prod environment               |
| `silent`  | Suppresses logs (great for testing)|

### 🐳 3.1. Docker Workflow (via `make`)

This section explains how to manage the Docker container lifecycle and execute tasks such as building, running, and stopping the Docker container, all through `make` commands.

#### Build the Docker Image

```
make docker ACTION=build
```

#### Run the Container

```
make docker ACTION=run CONFIG_PATH=config/local.yaml LOG_MODE=local 
```

#### Stop the Container

```
make docker ACTION=stop
```

#### Start a Stopped Container

```
make docker ACTION=start
```

#### Remove Container & Image

```
make docker ACTION=remove
```

#### Run the database migrator inside the container

```
make docker ACTION=exec EXEC=migrate STORAGE=storage/db.sqlite MIGRATIONS=migrations/sqlite
```

### 🖥️ 3.2. Local Usage (via `make`)

This section describes how to run the project and migrate the database locally using `make`. It includes commands for running the application locally and for performing database migrations.

#### Run the Project Locally

```
make run CONFIG_PATH=config/local.yaml LOG_MODE=local
```

#### Run the Database Migrator Locally

```
make migrate STORAGE=storage/db.sqlite MIGRATIONS=migrations/sqlite
```

### 🧪 4. Testing (via `make`)

This section explains how to run the tests for the project using `make`. It includes commands for running all tests or specific tests based on their type.

#### Run All Tests

```
make test
```

#### Run Specific Tests

If you want to run specific tests, you can specify it with the `TYPE` variable. For example, to run the **Functional** tests:

```
make test TYPE=functional
```

## API Reference

This microservice exposes a single gRPC method through the `Book` service. The following describes the structure of the API.

#### `BookRequest`

```proto
message BookRequest {
  Cinema cinema = 1;
  Movie movie = 2;
  Session session = 3;
}
```

#### `Cinema`

```proto
message Cinema {
  string name = 1;
  string location = 2;
}
```

#### `Movie`

```proto
message Movie {
  string title = 1;
  string genre = 2;
  string country = 3;
  google.protobuf.Timestamp premier = 4;
  google.protobuf.Duration duration = 5;
}
```

#### `Session`

```proto
message Session {
  int32 screen = 1;
  int32 seat = 2;
  google.protobuf.Timestamp date = 3;
}
```

Here's an example of what a `BookRequest` might look like when interacting with the service using gRPC:

```json
{
  "cinema": {
    "name": "IMAX Central",
    "location": "Downtown"
  },
  "movie": {
    "title": "Inception",
    "genre": "Sci-Fi",
    "country": "USA",
    "premier": "2010-07-16T00:00:00Z",
    "duration": "7200s"
  },
  "session": {
    "screen": 2,
    "seat": 12,
    "date": "2025-04-16T19:00:00Z"
  }
}
```

#### `BookResponse`

```proto
message BookResponse {
  Order order = 1;
}
```

#### `Order`

```proto
message Order {
  string ticket = 1;
}
```

Here's an example of what a `BookResponse` might look like when interacting with the service using gRPC:

```json
{
  "order": {
    "ticket": "abc123xyz"
  }
}
```

## Author

[**@xoticdsign**](https://github.com/xoticdsign). Crafted with care as a part of a pet project focused on clean architecture and gRPC microservices.

## License

[**MIT**](https://mit-license.org/)
