# Echo Casbin DDD App

A robust Go web application built with Echo framework, Casbin for authorization, and following Domain-Driven Design (DDD) principles.

## Features

- **Authentication**: JWT-based authentication for users and API key authentication for services
- **Authorization**: Role-based access control using Casbin
- **API Documentation**: Swagger/OpenAPI documentation
- **WebSocket Support**: Real-time communication
- **Clean Architecture**: Following DDD principles with clear separation of concerns
- **Database**: PostgreSQL with GORM ORM
- **Docker Support**: Easy deployment with Docker and Docker Compose

## Architecture

This application follows Domain-Driven Design (DDD) principles with a clean architecture approach:

- **Domain Layer**: Contains the core business logic, entities, and repository interfaces
- **Application Layer**: Contains use cases that orchestrate the domain logic
- **Infrastructure Layer**: Contains implementations of repositories and external services
- **Interface Layer**: Contains HTTP handlers, middleware, and WebSocket handlers

## Prerequisites

- Go 1.23 or higher
- PostgreSQL
- Docker and Docker Compose (optional)

## Installation

### Using Docker Compose

1. Clone the repository:
   ```bash
   git clone https://github.com/example/echo-casbin-ddd-app.git
   cd echo-casbin-ddd-app
   ```

2. Start the application:
   ```bash
   docker-compose up -d
   ```

3. Run database migrations:
   ```bash
   docker-compose exec app /app/echo-casbin-ddd-app --migrate
   ```

### Manual Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/example/echo-casbin-ddd-app.git
   cd echo-casbin-ddd-app
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Set up environment variables (see `.env.example` for required variables)

4. Run the application:
   ```bash
   go run cmd/main.go --migrate
   ```

## Usage

### API Endpoints

The application provides the following API endpoints:

- **User Management**:
  - `POST /v1/users/register`: Register a new user
  - `POST /v1/users/login`: Login a user
  - `GET /v1/users/:id`: Get user details
  - `PUT /v1/users/:id`: Update user details
  - `POST /v1/users/:id/change-password`: Change user password
  - `GET /v1/users`: List users

- **API Client Management**:
  - `POST /v1/api-clients`: Create a new API client
  - `GET /v1/api-clients/:id`: Get API client details
  - `PUT /v1/api-clients/:id`: Update API client details
  - `DELETE /v1/api-clients/:id`: Delete an API client
  - `POST /v1/api-clients/:id/regenerate-key`: Regenerate API key
  - `GET /v1/api-clients`: List API clients

### Authentication

- **JWT Authentication**: For user authentication, include the JWT token in the `Authorization` header:
  ```
  Authorization: Bearer <token>
  ```

- **API Key Authentication**: For API client authentication, include the API key in the header specified in the configuration (default: `X-API-Key`):
  ```
  X-API-Key: <api-key>
  ```

### API Documentation

The API documentation is available at:

- Swagger UI: `http://localhost:8080/swagger/`
- OpenAPI JSON: `http://localhost:8080/swagger/doc.json`
- OpenAPI YAML: `http://localhost:8080/swagger/doc.yaml`

## Development

### Project Structure

```
.
├── cmd/                  # Application entry points
├── docs/                 # Documentation and Swagger files
├── internal/             # Internal packages
│   ├── application/      # Application layer (use cases)
│   ├── domain/           # Domain layer (entities, repositories)
│   ├── infrastructure/   # Infrastructure layer (implementations)
│   └── interfaces/       # Interface layer (handlers, middleware)
├── pkg/                  # Shared packages
├── web/                  # Static web files
├── Dockerfile            # Docker configuration
├── docker-compose.yml    # Docker Compose configuration
└── go.mod                # Go module definition
```

### Running Tests

```bash
go test ./...
```

### Generating Swagger Documentation

```bash
swag init -g cmd/main.go -o docs
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.
