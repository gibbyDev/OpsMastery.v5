<!-- # OpsMastery.v5 — Copilot Instructions

## Big Picture Architecture
- Modular, microservices-based platform for real-time chat, ticketing, user management, and WebRTC video/audio calls.
- Major services:
  - **API Service** (`cmd/api/`): REST endpoints for users, tickets, authentication, and business logic.
  - **Chat Service** (`cmd/chat-service/`): WebSocket-based real-time chat, JWT-authenticated, persists messages to DB.
  - **WebRTC Service** (`cmd/webrtc-service/`): WebSocket signaling for video/audio calls, JWT-authenticated.
  - **Database** (`psql_bp` in Docker): PostgreSQL, schema managed via GORM models in `internal/models/`.
- Data flows through REST and WebSocket endpoints, with JWT for authentication and session management.

## Developer Workflows
- **Build/Run:** Use `make build`, `make run`, or `docker-compose up` for service orchestration. Each service has its own Dockerfile.
- **Seeding:** Run `make seed` or use the `seed` Docker service to inject test data (see `dev/seed.go`).
- **Testing:** Use `make test` for unit tests, `make itest` for integration tests (see `internal/database/database_test.go`).
- **Migrations:** Models are auto-migrated on backend startup (`database.Init()` in main entrypoints).
- **Debugging:** Enable GORM debug mode in handlers for SQL inspection. Use Postman/curl for API/WebSocket testing.

## Project-Specific Conventions
- **JWT Auth:** Access tokens are returned in JSON, refresh tokens as HttpOnly cookies. All protected routes require `Authorization: Bearer <access_token>`.
- **Case Normalization:** Emails, usernames, and names are normalized to lowercase before DB writes and queries.
- **Chat Model:** Supports group, ticket, and private chats. `Chat` model uses many-to-many `Users`, optional `TicketID`, and `IsPrivate` flag.
- **Search:** User search endpoint (`/api/v1/users/search?q=...`) matches on lowercase username, email, or name.
- **OAuth:** Google/GitHub OAuth via `goth` library, with provider setup in `internal/oauth/oauth.go` and handlers in `internal/handlers/authHandler.go`.
- **WebSocket:** Chat and WebRTC services require JWT token as query param for connection. Messages are broadcast and persisted.

## Integration Points & External Dependencies
- **Docker Compose:** All services and the database are orchestrated via `docker-compose.yml`. Environment variables are loaded from `.env`.
- **Database:** PostgreSQL, accessed via GORM. Models in `internal/models/` define schema and relationships.
- **Frontend:** Expected to connect via REST (API) and WebSocket (chat/webrtc) endpoints. CORS is enabled for `http://localhost:3000`.
- **OAuth Providers:** Google/GitHub credentials are set via environment variables and loaded in Docker.

## Key Files & Directories
- `cmd/` — Service entrypoints (api, chat-service, webrtc-service, seed)
- `internal/models/` — GORM models for all entities
- `internal/handlers/` — REST and WebSocket handlers
- `internal/oauth/oauth.go` — OAuth provider setup
- `internal/server/routes.go` — Route registration for all endpoints
- `dev/seed.go` — Database seeding logic
- `Makefile` — Build, test, seed, and run commands
- `docker-compose.yml` — Service orchestration
- `.env` — Environment variable configuration

## Examples
- **User Search:** `GET /api/v1/users/search?q=alice` returns users matching 'alice' in username, email, or name (case-insensitive).
- **Chat WebSocket:** Connect to `ws://localhost:5000/ws?token=<access_token>` for real-time messaging.
- **OAuth Login:** Redirect to `/api/v1/auth/google` or `/api/v1/auth/github` to start OAuth flow.

---
**Feedback:** If any section is unclear or missing, please specify which workflows, conventions, or integration points need more detail. -->
