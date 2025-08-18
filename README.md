# OpsMastery.v5

OpsMastery.v5 is a modular, production-ready platform for real-time chat, video/audio calling (WebRTC), and RESTful business operations. It is built in Go, uses JWT authentication, and is designed for easy integration with modern frontends (e.g., Next.js).

## Project Structure

```
OpsMastery.v5/
├── cmd/
│   ├── api/                # Main REST API server (users, tickets, auth, etc.)
│   ├── chat-service/       # Real-time chat microservice (WebSocket)
│   └── webrtc-service/     # WebRTC signaling microservice (WebSocket)
├── internal/
│   ├── chat_service/       # Chat logic
│   ├── webrtc_service/     # WebRTC signaling logic
│   ├── handlers/           # REST API handlers
│   ├── models/             # Data models
│   ├── database/           # DB connection and migrations
│   ├── middleware/         # JWT and role middleware
│   ├── server/             # REST API server setup
│   └── utils/              # Utility functions (JWT, email, etc.)
├── docker-compose.yml      # Multi-service orchestration
├── Dockerfile*             # Dockerfiles for each service
├── go.mod, go.sum          # Go module files
└── README.md
```

## Services Overview

### 1. REST API (`cmd/api`)
- Handles authentication, user management, tickets, and other business logic.
- JWT-secured endpoints.
- Integrates with PostgreSQL.

### 2. Chat Service (`cmd/chat-service`)
- Real-time chat via WebSocket.
- JWT authentication for all connections.
- Broadcasts messages to all connected clients.
- Persists chat messages in the database.

### 3. WebRTC Signaling Service (`cmd/webrtc-service`)
- Handles signaling for video/audio calls (SDP/ICE exchange).
- JWT authentication for all connections.
- Designed for integration with WebRTC clients (e.g., browser, mobile).

## Frontend Integration

- Designed to be consumed by any modern frontend (e.g., Next.js).
- REST API: Use `fetch`/`axios` with JWT in cookies or headers.
- Chat/WebRTC: Connect via WebSocket with JWT as a query param.

## Getting Started

### Prerequisites
- [Go](https://golang.org/doc/install) >= 1.20
- [Docker](https://docs.docker.com/get-docker/)
- [Node.js](https://nodejs.org/) (for frontend, optional)

### Clone the Repository

```sh
git clone https://github.com/gibbyDev/OpsMastery.v5.git
cd OpsMastery.v5
```

### Running with Docker Compose

```sh
docker compose up --build
```
- This will start the REST API, chat service, WebRTC signaling service, and PostgreSQL database.
- Services are exposed on:
	- REST API: `http://localhost:8080`
	- Chat: `ws://localhost:5000/ws`
	- WebRTC: `ws://localhost:4000/ws`

### Running Locally (without Docker)

1. Start PostgreSQL (see `docker-compose.yml` for env vars).
2. Run each service in a separate terminal:
	 ```sh
	 go run cmd/api/main.go
	 go run cmd/chat-service/main.go
	 go run cmd/webrtc-service/main.go
	 ```

### Environment Variables
- See `.env.example` or `docker-compose.yml` for required variables (DB connection, JWT secret, etc.).

## Development & Testing

- All code is in Go, organized for easy extension.
- Unit and integration tests are in `internal/*/tests`.
- Frontend integration examples available upon request.

## Contributing

Pull requests and issues are welcome! Please open an issue for major changes.

## License

MIT

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

## MakeFile

Run build make command with tests
```bash
make all
```

Build the application
```bash
make build
```

Run the application
```bash
make run
```
Create DB container
```bash
make docker-run
```

Shutdown DB Container
```bash
make docker-down
```

DB Integrations Test:
```bash
make itest
```

Live reload the application:
```bash
make watch
```

Run the test suite:
```bash
make test
```

Clean up binary from the last build:
```bash
make clean
```
