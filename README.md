# Authentication System

A full-stack authentication system with a Go backend and a simple HTML/CSS/JS frontend.

## Project Structure

- `cmd/` - Contains the main application entry point
- `internal/` - Internal application code
  - `config/` - Database and configuration
  - `controller/` - Request handlers
  - `delivery/` - HTTP delivery layer
  - `domain/` - Domain models and interfaces
  - `repository/` - Database repositories
  - `services/` - Business logic services
  - `usecases/` - Application use cases
- `pkg/` - Shared packages
  - `helper/` - Helper functions
  - `middleware/` - Middleware components
- `frontend/` - Frontend code
  - `css/` - Stylesheets
  - `js/` - JavaScript files
  - `img/` - Image assets

## Features

- User registration and login
- OAuth login with Google, GitHub, and Facebook (frontend UI)
- JWT token authentication
- Token refresh mechanism
- User profile management
- Account deletion
- Rate limiting
- Responsive frontend UI

## Getting Started

### Prerequisites

- Go 1.19+
- PostgreSQL
- Redis
- Docker and Docker Compose (optional)

### Environment Setup

1. Copy the example environment file and update with your configuration:
   ```
   cp .env.example .env
   ```

2. Update the values in `.env` to match your environment.

### Running with Docker Compose

The easiest way to run the application is with Docker Compose:

```
docker-compose up -d
```

Or using the provided Makefile:

```
make up_build
```

This will:
1. Build the backend Go binary
2. Start the backend service in a Docker container
3. Start the frontend service in a Docker container
4. Start PostgreSQL and Redis containers

You can then access:
- Frontend: http://localhost:3000
- Backend API: http://localhost:8080

### Running Locally

1. Start PostgreSQL and Redis separately
2. Run the backend:
   ```
   go run cmd/main.go
   ```
3. Run the frontend:
   ```
   cd frontend
   go run server.go
   ```

Or use the provided batch file (Windows):
```
run_auth_system.bat
```

## API Endpoints

### Public Endpoints

- `POST /auth/register` - Register a new user
- `POST /auth/login` - Login and get JWT token
- `GET /auth/google` - Initiate Google OAuth login (frontend UI only)
- `GET /auth/github` - Initiate GitHub OAuth login (frontend UI only)
- `GET /auth/facebook` - Initiate Facebook OAuth login (frontend UI only)
- `POST /auth/oauth/callback` - OAuth callback handler (frontend UI only)

### Protected Endpoints (requires authentication)

- `POST /user/logout` - Logout and invalidate token
- `GET /user/profile` - Get user profile
- `DELETE /user/delete` - Delete user account
- `POST /user/refresh` - Refresh access token

## Frontend

The frontend is a simple HTML/CSS/JavaScript application that interacts with the backend API.

- **Pages**:
  - Home page (`index.html`)
  - Login page (`login.html`)
  - Registration page (`register.html`)
  - User dashboard (`dashboard.html`)

- **Configuration**: The frontend can be configured by editing `js/config.js` to point to the correct backend API URL.

## Security Features

- JWT token-based authentication
- OAuth authentication (frontend UI)
- Password hashing
- Token refresh mechanism
- Rate limiting
- Token blacklisting
- CSRF protection for OAuth flows

## License

This project is licensed under the MIT License.
