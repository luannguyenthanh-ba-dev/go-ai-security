# Go AI Security

Backend service for AI-security workloads built on Go, following Domain-Driven Design (DDD) and Clean Architecture. The codebase is organised into isolated modules (`auth`, `users`) and infrastructure packages for persistence, caching, and shared tooling.

## Quick Start

```bash
# 1. Clone the repo
git clone https://github.com/luannguyenthanh-ba-dev/go-ai-security.git
cd go-ai-security

# 2. Install dependencies
go mod download

# 3. Copy the sample environment and adjust values
cp local.env .env

# 4. Start MongoDB (default: mongodb://admin:admin_password@localhost:27017/)
# 5. Start Redis   (default: redis://localhost:6379)

# 6. Run the API
go run ./cmd
```

The service listens on `:8080` by default, or on the port defined in `PORT`. Health check: `GET http://localhost:8080/health`.

## Architecture At A Glance

- **Domain (`internal/<module>/domain`)**: Pure business entities, domain errors, and events.
- **Use Case (`internal/<module>/usecase`)**: Application services that orchestrate domain logic via injected interfaces.
- **Repository (`internal/<module>/repository`)**: Infrastructure implementations (MongoDB, Redis, etc.).
- **Delivery (`internal/<module>/delivery/http`)**: Gin handlers, routing, and middleware wiring.
- **DTOs (`internal/<module>/dto`)**: External request/response models kept separate from domain entities.
- **Shared Infrastructure (`pkg`, `infrastructure`)**: Logger (`zap`), middleware, cache client abstractions, and response helpers.

```
.
├── cmd/                 # Application entrypoint with DI wiring and swagger registration
├── config/              # Environment loading, MongoDB & Redis connection factories
├── docs/                # Generated Swagger/OpenAPI artefacts
├── infrastructure/      # Cache implementations (Redis)
├── internal/
│   ├── auth/            # Authentication bounded context (JWT, refresh tokens)
│   └── users/           # User registration, profile, password management
└── pkg/                 # Logger, middleware, shared types, response helpers
```

## Features

- JWT-based authentication with refresh token rotation (`/auth/login`, `/auth/tokens/refresh`).
- User registration, profile updates, password change, and self-service endpoints (protected by JWT middleware).
- MongoDB persistence and Redis-backed cache abstraction.
- Centralised logging with `zap` and structured error responses.
- OpenAPI/Swagger documentation served under `/swagger/index.html`.

## Configuration

All runtime configuration is handled through `.env` (loaded via `viper`). Key variables:

| Variable | Description | Default (see `local.env`) |
| --- | --- | --- |
| `APP_NAME` | Application name for logging | `go-ai-security` |
| `PORT` | HTTP port | `8080` |
| `APP_ENV` | `development` or `production` | `development` |
| `MONGO_URI` | MongoDB connection string | `mongodb://admin:admin_password@localhost:27017/...` |
| `MONGO_DATABASE` | MongoDB database name | `ai-security` |
| `REDIS_HOST` / `REDIS_PORT` / `REDIS_DB` | Redis connection | `localhost` / `6379` / `0` |
| `PASSWORD_HASH_SALT_ROUNDS` | bcrypt cost factor | `10` |
| `JWT_SECRET` | HMAC secret for tokens | `go` |
| `JWT_EXPIRES_IN` | Access token TTL in seconds | `300` (example) |

> **Note:** `JWT_EXPIRES_IN` must be an integer representing seconds. Update `local.env` accordingly when creating your own `.env`.

## Generating & Serving Swagger

Swagger specs are generated with [`swag`](https://github.com/swaggo/swag). Regenerate the docs whenever handler annotations change:

```bash
# Install swag (once)
go install github.com/swaggo/swag/cmd/swag@latest

# Generate swagger files into ./docs (leverages //go:generate directive)
go generate ./cmd
# or call swag directly
swag init -g cmd/main.go -o ./docs
```

While the server is running, browse the docs at `http://localhost:8080/swagger/index.html`. The underlying JSON can be fetched from `/docs/swagger.json`.

## Development Commands

```bash
# Run unit tests
go test ./...

# Format go modules
go fmt ./...

# Build binary
go build -o bin/go-ai-security ./cmd

# Lint (if golangci-lint installed)
golangci-lint run
```

## Module Overview

- `internal/auth`: Auth use cases, JWT service (`NewJWTService`), Redis-backed token revocation, Gin handlers for login and refresh.
- `internal/users`: User CRUD use cases, MongoDB repository, DTO mappings, protected profile endpoints.
- `pkg/middleware`: JWT authentication middleware used by protected routes.
- `infrastructure/cache`: Cache client factory with Redis implementation and typed error helpers.

## Contributing

1. Fork the repo & create a branch (`git checkout -b feature/my-change`).
2. Implement your feature or bug fix with accompanying tests.
3. Ensure `go test ./...` and (if available) `golangci-lint run` pass.
4. Regenerate Swagger docs if handler contracts changed.
5. Submit a PR describing the change and testing performed.

## License

Distributed under the MIT License. See `LICENSE` for details.
