# go-fx-clean-architecture

A production-style Go service scaffold running **HTTP (Gin)** and **gRPC** servers side by side, wired together with **Uber FX** dependency injection. It demonstrates a clean, layered architecture (controller/handler → service → repository) backed by **GORM** and PostgreSQL.

## Features

- **Uber FX** DI — every package exposes a `Module`, composed in `cmd/main.go`.
- **Dual transport** — Gin REST API and gRPC service share the same service layer.
- **GORM** data access with read-replica support (via `dbresolver`).
- **cleanenv + godotenv** configuration from environment / `.env`.
- **zap** structured logging.
- **protoc** code generation with a `make gen-protos` target.
- **Docker** + **docker-compose** for local development.

## Architecture

```
HTTP request ─▶ controllers/  ┐
                              ├─▶ services/ ─▶ repositories/ ─▶ GORM ─▶ PostgreSQL
gRPC request ─▶ handlers/     ┘
```

- **controllers** (Gin) and **handlers** (gRPC) are thin: extract input, call a service, shape the response.
- **services** hold business logic, validation, and orchestration; they return `(data, statusCode, error)`.
- **repositories** own all queries and accept `opts ...QueryOption` for optional relationship preloading.
- Layers depend on **interfaces** (`IUserService`, `IUserRepository`), so transports and persistence stay decoupled.

## Project layout

```
cmd/
  main.go                  # fx.New(config, internal, app).Run()
config/
  config.go                # env config (cleanenv + godotenv)
  module.go                # fx provider
app/
  module.go                # logger + handlers + controllers + servers
  controllers/             # HTTP (Gin) controllers + route registration
  handlers/                # gRPC handlers + service registration
  servers/                 # gRPC + HTTP server lifecycle (fx hooks)
  interceptors/            # gRPC metadata interceptor
internal/
  module.go                # db + repositories + services
  dbconnections/           # GORM connection + auto-migrate
  repositories/            # data access + QueryOption helpers
  services/                # business logic
  models/                  # GORM models
  dto/                     # request/response DTOs (per domain)
pkg/
  db/                      # GORM connect (driver + replicas)
  logger/                  # zap logger provider
protos/
  user/user.proto          # proto source
genprotos/
  user/                    # generated code (never edit by hand)
Dockerfile
docker-compose.yaml
Makefile
.env.sample
```

## Getting started

### Prerequisites

- Go 1.24+
- Docker & Docker Compose (for the easy path)
- `protoc`, `protoc-gen-go`, `protoc-gen-go-grpc` (only if regenerating protos)

### Run with Docker

```bash
cp .env.sample .env
make docker-up        # starts PostgreSQL + the app
```

### Run locally

```bash
cp .env.sample .env   # point DB_* at a reachable PostgreSQL
make run
```

The HTTP server listens on `:8080`, the gRPC server on `:9090` (with reflection enabled).

## Configuration

All configuration is read from environment variables (or `.env`). See `.env.sample`:

| Variable           | Default                    | Description                     |
| ------------------ | -------------------------- | ------------------------------- |
| `APP_NAME`         | `golang-project-structure` | Application name                |
| `APP_ENV`          | `development`              | Environment                     |
| `APP_DEBUG`        | `true`                     | Debug mode (Gin + zap)          |
| `HTTP_PORT`        | `8080`                     | HTTP (Gin) port                 |
| `GRPC_PORT`        | `9090`                     | gRPC port                       |
| `DB_DRIVER`        | `postgres`                 | Database driver                 |
| `DB_HOST`          | `localhost`                | Database host                   |
| `DB_PORT`          | `5432`                     | Database port                   |
| `DB_USER`          | `postgres`                 | Database user                   |
| `DB_PASSWORD`      | `postgres`                 | Database password               |
| `DB_NAME`          | `app`                      | Database name                   |
| `DB_SSLMODE`       | `disable`                  | SSL mode                        |
| `DB_REPLICA_HOSTS` | _(empty)_                  | Comma-separated replica hosts   |

## HTTP API

Base path: `/api/v1`

| Method   | Path          | Description     |
| -------- | ------------- | --------------- |
| `GET`    | `/health`     | Health check    |
| `POST`   | `/users`      | Create a user   |
| `GET`    | `/users`      | List users      |
| `GET`    | `/users/:id`  | Get a user      |
| `PUT`    | `/users/:id`  | Update a user   |
| `DELETE` | `/users/:id`  | Delete a user   |

Example:

```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Ada Lovelace","email":"ada@example.com"}'

curl "http://localhost:8080/api/v1/users?page=1&limit=20"
```

## gRPC API

The `user.UserService` mirrors the REST endpoints (`Create`, `GetByID`, `List`, `Update`, `Delete`). With reflection enabled you can explore it via [grpcurl](https://github.com/fullstorydev/grpcurl):

```bash
grpcurl -plaintext localhost:9090 list
grpcurl -plaintext -d '{"name":"Ada","email":"ada@example.com"}' \
  localhost:9090 user.UserService/Create
```

## Code generation

Proto sources live in `protos/<domain>/`; generated Go lands in `genprotos/<domain>/`.

```bash
make gen-protos folder=user
```

## Make targets

| Target                | Description                          |
| --------------------- | ------------------------------------ |
| `make run`            | Run the app locally                  |
| `make build`          | Build the binary to `bin/server`     |
| `make tidy`           | `go mod tidy`                        |
| `make docker-up`      | Start app + dependencies             |
| `make docker-down`    | Stop containers                      |
| `make docker-rebuild` | Rebuild and restart containers       |
| `make gen-protos`     | Generate Go code from `.proto` files |

## Adding a new domain

1. Add the model in `internal/models/` and the DTOs in `internal/dto/<domain>/`.
2. Add the repository (`I<Domain>Repository`) and register it in `internal/repositories/module.go`.
3. Add the service (`I<Domain>Service`) and register it in `internal/services/module.go`.
4. Add a controller in `app/controllers/` (HTTP) and/or a handler in `app/handlers/` (gRPC), and register each in its `module.go`.
5. For gRPC, define `protos/<domain>/*.proto` and run `make gen-protos folder=<domain>`.
