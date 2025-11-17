## go-todo-service

A production-style REST API for managing todo tasks with user authentication, implemented in Go using a clean architecture approach.

### Features
- JWT-based authentication with signup and login endpoints
- Task CRUD restricted to the authenticated user
- BCrypt password hashing and short-lived access tokens
- Structured JSON logging with request IDs and panic recovery
- PostgreSQL persistence with SQL migrations
- Docker/Docker Compose for local development
- OpenAPI specification, Makefile-powered workflows, and service-layer unit tests

### Project Structure
```
cmd/api              # Application entrypoint
internal/config      # Environment configuration
internal/domain      # Domain models and errors
internal/handlers    # HTTP handlers, routes, middleware
internal/repository  # Persistence interfaces and PostgreSQL implementations
internal/service     # Business logic (auth/tasks)
pkg                  # Shared utilities (jwt, logger, password, uuid)
migrations           # SQL migrations
```

### Requirements
- Go 1.24+
- Docker & Docker Compose (for containerised setup)
- PostgreSQL 15+ (if running without Docker)

### Configuration
The application is configured via environment variables (defaults shown):

| Variable | Default | Description |
| --- | --- | --- |
| `SERVER_PORT` | `8080` | API listen port |
| `DB_HOST` | `postgres` | Database host |
| `DB_PORT` | `5432` | Database port |
| `DB_USER` | `todo` | Database user |
| `DB_PASSWORD` | `todo` | Database password |
| `DB_NAME` | `todo` | Database name |
| `DB_SSL_MODE` | `disable` | PostgreSQL SSL mode |
| `JWT_SECRET` | _required_, ≥32 chars | Secret used to sign JWTs |
| `JWT_TTL_MINUTES` | `15` | Access-token lifetime |

### Running with Docker Compose
```bash
docker compose up --build
```
This starts both the API and PostgreSQL. The initial migration runs automatically through the mounted `001_init.up.sql` file. The API logs every request with an `X-Request-ID` to aid correlation across services.

### Running Locally
1. Start PostgreSQL and apply the migration in `migrations/001_init.up.sql`.
2. Export the required environment variables (particularly `JWT_SECRET`).
3. Build and run:
```bash
go build -o bin/go-todo-service ./cmd/api
JWT_SECRET=supersecret-supersecret-supersecret!! DB_HOST=localhost DB_USER=todo DB_PASSWORD=todo DB_NAME=todo ./bin/go-todo-service
```
   For convenience you can also use the provided `Makefile`:
```bash
make build
JWT_SECRET=supersecret-supersecret-supersecret!! make run
```

### Database Migrations
The project ships with SQL migrations under `migrations/`. If you use [golang-migrate](https://github.com/golang-migrate/migrate):
```bash
migrate -path migrations -database "postgres://todo:todo@localhost:5432/todo?sslmode=disable" up
```

### API Usage
Refer to `openapi.yaml` for the full specification. Sample curl commands:

```bash
# Signup
curl -X POST http://localhost:8080/auth/signup \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"secret123"}'

# Login
TOKEN=$(curl -s -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"secret123"}' | jq -r '.token')

# Create task
curl -X POST http://localhost:8080/tasks \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title":"Pay bills","description":"Pay electricity"}'
```

### Testing
Service layer tests live under `internal/service/...`. Run all tests with:
```bash
go test ./...
# or
make test
```

### Documentation
- **Swagger UI**: http://localhost:8080/docs
- **OpenAPI**: `openapi.yaml` (also available at `/docs/openapi.{json,yaml}`)

Error responses include the request identifier where available:
```json
{"error":"invalid credentials","request_id":"4c6f..." }
```

### License
MIT
