# OAuth2 Authentication System

A lightweight OAuth2 authentication system powered by Golang (backend). Built with Clean Architecture, it supports access & refresh token management and is designed for easy integration and scalability.

## Features

- OAuth2 authentication flow (login, logout, token refresh)
- Access & refresh token management
- Google OAuth2 integration
- Clean Architecture (separation of concerns)
- Easy to integrate with other services
- Docker support for development

## Flow Chart
![alt text](image.png)

## Folder Structure

```
.
├── cmd/                # Entry point for the application
├── delivery/           # HTTP delivery layer
├── dto/                # Data transfer objects
├── middleware/         # Middleware (auth, etc.)
├── repository/         # Data access layer
├── service/            # Business logic (Google OAuth, etc.)
├── usecase/            # Application use cases
├── utils/              # Utility functions (config, JWT, etc.)
├── sqlmigrations/      # Database migration scripts
├── docker/             # Docker configuration
├── docs/               # Documentation & Swagger files
├── main.go             # Main application file
├── go.mod, go.sum      # Go dependencies
├── dbconfig.yml        # Database configuration
├── keyfile.json        # Google OAuth2 credentials
└── README.md           # Project documentation
```

## Getting Started

### Prerequisites

- Go 1.18+
- Docker & Docker Compose

### Backend Setup

```bash
# Clone the repository
git clone https://github.com/johnquangdev/oauth2.git

# Run database and redis with Docker Compose
docker-compose --env-file .env -f docker/docker-compose.yaml up -d

# Run migrations (if needed)
# Example: psql -U <user> -d <db> -f sqlmigrations/1_init.sql

# Start backend server
go run .
```

## API Endpoints

- `POST /login` - User login
- `POST /refresh` - Refresh access token
- `GET /profile` - Get user profile
- `POST /logout` - User logout
- `GET /auth/google` - Google OAuth2 login

> Xem chi tiết trong file `docs/swagger.yaml` hoặc `swagger.json`.


