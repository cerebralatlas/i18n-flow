# i18n-flow Backend

A powerful backend service for internationalization management in application development. Built with Go 1.23, Gin, GORM, and Redis.

## Overview

i18n-flow is a comprehensive backend system designed to streamline the internationalization (i18n) process for software applications. It provides a centralized platform for managing translation keys, values, and languages across multiple projects through a RESTful API.

## Features

- **Project Management**: Create and organize multiple translation projects
- **Multi-language Support**: Add, edit, and manage different languages and locales
- **Translation Management**: Centralize all translations with context support
- **API-driven Architecture**: RESTful API for seamless integration
- **CLI Integration**: Command-line tool support for automated translation workflows
- **Batch Operations**: Efficiently handle bulk translations
- **Export/Import**: Flexible data interchange in various formats
- **Dashboard**: Statistics for monitoring translation progress
- **Authentication**: JWT-based auth for admin interface and API key auth for CLI tools

## Tech Stack

- **Language**: Go 1.23
- **Web Framework**: Gin 1.9.1
- **ORM**: GORM 1.30.0
- **Database**: MySQL/MariaDB, SQLite (for testing)
- **Documentation**: Swagger/OpenAPI
- **Authentication**: JWT & API keys
- **Caching**: Redis
- **Logging**: Zap with lumberjack rotation

## API Endpoints

### Authentication

- `POST /api/login`: Admin login
- `POST /api/refresh`: Refresh JWT token
- `GET /api/cli/auth`: Validate CLI API key

### Projects

- `POST /api/projects`: Create project
- `GET /api/projects`: List projects with pagination
- `GET /api/projects/detail/:id`: Get project details
- `PUT /api/projects/update/:id`: Update project
- `DELETE /api/projects/delete/:id`: Delete project

### Languages

- `GET /api/languages`: List languages
- `POST /api/languages`: Create language
- `PUT /api/languages/:id`: Update language
- `DELETE /api/languages/:id`: Delete language

### Translations

- `POST /api/translations`: Create translation
- `POST /api/translations/batch`: Batch create translations
- `GET /api/translations/by-project/:project_id`: Get project translations
- `GET /api/translations/matrix/by-project/:project_id`: Get translation matrix
- `GET /api/translations/:id`: Get translation details
- `PUT /api/translations/:id`: Update translation
- `DELETE /api/translations/:id`: Delete translation
- `POST /api/translations/batch-delete`: Batch delete translations
- `GET /api/exports/project/:project_id`: Export project translations
- `POST /api/imports/project/:project_id`: Import project translations

### CLI Tool Integration

- `GET /api/cli/translations`: Get translations for CLI
- `POST /api/cli/keys`: Push new translation keys from CLI

### Dashboard

- `GET /api/dashboard/stats`: Get system statistics

## Getting Started

### Prerequisites

- Go 1.23+
- MySQL
- Git
- Redis (for caching)

### Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/yourusername/i18n-flow.git
   cd i18n-flow/admin-backend
   ```

   Or download the latest release from the releases page.

2. Install dependencies:

   ```bash
   go mod tidy
   ```

3. Configure environment:

   ```bash
   cp .env.example .env
   ```

   Edit `.env` file with your database credentials and settings:

   ```
   DB_DRIVER=mysql           # mysql
   DB_USERNAME=your_db_user
   DB_PASSWORD=your_db_password
   DB_HOST=localhost
   DB_PORT=3306
   DB_NAME=i18n_flow
   
   JWT_SECRET=your_secure_jwt_secret
   JWT_EXPIRATION_HOURS=24
   JWT_REFRESH_SECRET=your_secure_refresh_secret
   JWT_REFRESH_EXPIRATION_HOURS=168
   
   CLI_API_KEY=your_secure_api_key
   
   ADMIN_NAME=your_name
   ADMIN_PASSWORD=your_password
   
   REDIS_HOST=localhost
   REDIS_PORT=6379
   REDIS_PASSWORD=
   REDIS_DB=0
   
   LOG_LEVEL=info           # debug, info, warn, error
   LOG_FILE=./logs/app.log
   ```

4. Start the server:

   ```bash
   # Using air for hot-reload during development
   air
   
   # Or run directly
   go run cmd/server/main.go
   
   # Or build and run the binary
   go build -o i18n-flow
   ./i18n-flow
   ```

   The server will start on port 8080 by default. The database tables will be automatically created.

5. Access the Swagger documentation:

   ```
   http://localhost:8080/swagger/index.html
   ```

### Initial Access

On first run, the system creates a default admin user based on your environment variables:

- Username: Value from `ADMIN_NAME` in .env file
- Password: Value from `ADMIN_PASSWORD` in .env file

It's highly recommended to change this password after the first login through the admin interface.

## Development

### API Documentation

The API documentation is automatically generated with Swagger. You can access it at:

```
http://localhost:8080/swagger/index.html
```

### Hot Reload Development

This project uses [Air](https://github.com/cosmtrek/air) for hot reloading during development. The configuration is in `.air.toml`.

To install Air:

```bash
go install github.com/cosmtrek/air@latest
```

Then run:

```bash
air
```

### Project Structure

- `/cmd/server`: Application entry point
- `/internal/api/handlers`: HTTP request handlers
- `/internal/api/middleware`: Request processing middleware
- `/internal/api/response`: Response formatting
- `/internal/api/routes`: API route definitions
- `/internal/config`: Application configuration
- `/internal/container`: Dependency injection container
- `/internal/domain`: Core domain models and interfaces
- `/internal/repository`: Data access layer
- `/internal/service`: Business logic layer
- `/internal/utils`: Utility functions
- `/docs`: Auto-generated API documentation
- `/migrations`: Database migration files
- `/utils`: Shared utility functions
- `/tests`: Test files

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage report
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

## CLI Tool Integration

The backend provides dedicated endpoints for CLI tool integration, allowing automated workflows:

- Authentication via API key
- Pulling translations for specific projects and languages
- Pushing new translation keys
- Checking API connection status

## Docker Deployment

```bash
# Build image
docker build -t i18n-flow-backend .

# Run container with MySQL
docker run -p 8080:8080 --env-file .env i18n-flow-backend

# Or run with Docker Compose (recommended for production)
# Create a docker-compose.yml file with MySQL and Redis services
docker-compose up -d
```

### Environment Variables for Docker

When running in Docker, you may want to adjust these environment variables:

```
DB_HOST=mysql  # Use service name from docker-compose
REDIS_HOST=redis  # Use service name from docker-compose
```

## License

This project is licensed under the MIT License.
