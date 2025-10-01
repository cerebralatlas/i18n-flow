# i18n-flow Backend

A powerful backend service for internationalization management in application development. Built with Go 1.23, Gin, GORM, and Redis.

## Overview

i18n-flow is a comprehensive backend system designed to streamline the internationalization (i18n) process for software applications. It provides a centralized platform for managing translation keys, values, and languages across multiple projects through a RESTful API.

## Features

- **Project Management**: Create and organize multiple translation projects with role-based access control
- **Project Member Management**: Invite users to projects with viewer, editor, or owner roles
- **Multi-language Support**: Add, edit, and manage different languages and locales
- **Translation Management**: Centralize all translations with context support
- **API-driven Architecture**: RESTful API for seamless integration
- **CLI Integration**: Command-line tool support for automated translation workflows
- **Batch Operations**: Efficiently handle bulk translations
- **Export/Import**: Flexible data interchange in various formats
- **Dashboard**: Statistics for monitoring translation progress
- **Authentication**: JWT-based auth for admin interface and API key auth for CLI tools
- **User Management**: Admin functionality for creating and managing system users
- **Role-based Permissions**: Fine-grained access control for projects and translations

## Tech Stack

- **Language**: Go 1.23
- **Web Framework**: Gin 1.9.1
- **ORM**: GORM 1.30.0
- **Database**: MySQL
- **Documentation**: Swagger/OpenAPI
- **Authentication**: JWT & API keys
- **Caching**: Redis
- **Logging**: Zap with lumberjack rotation

## API Endpoints

### Authentication & User

- `POST /api/login`: Admin login
- `POST /api/refresh`: Refresh JWT token
- `GET /api/user/info`: Get current user info
- `POST /api/user/change-password`: Change current user password
- `GET /api/cli/auth`: Validate CLI API key

### User Management (Admin)

- `POST /api/users`: Create new user
- `GET /api/users`: List all users with pagination
- `GET /api/users/:id`: Get user details
- `PUT /api/users/:id`: Update user information
- `POST /api/users/:id/reset-password`: Reset user password
- `DELETE /api/users/:id`: Delete user
- `GET /api/user-projects/:user_id`: Get user's projects

### Projects

- `POST /api/projects`: Create project
- `GET /api/projects`: List projects with pagination
- `GET /api/projects/accessible`: Get accessible projects for current user
- `GET /api/projects/detail/:id`: Get project details
- `PUT /api/projects/update/:id`: Update project
- `DELETE /api/projects/delete/:id`: Delete project
- `GET /api/projects/:project_id/members`: Get project members
- `POST /api/projects/:project_id/members`: Add project member
- `PUT /api/projects/:project_id/members/:user_id`: Update member role
- `DELETE /api/projects/:project_id/members/:user_id`: Remove project member
- `GET /api/projects/:project_id/members/:user_id/permission`: Check user permission in project

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

### Permissions & Roles

The system implements a role-based access control system with three permission levels:

**Project Roles:**

- **Viewer**: Can view projects, translations, and export data
- **Editor**: Viewer permissions + can create, update, and delete translations
- **Owner**: Editor permissions + can manage project settings and members

**System Roles:**

- **Admin**: Full system access, can manage users, languages, and all projects
- **User**: Standard user with project-specific permissions

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
   
   ADMIN_USERNAME=admin
   ADMIN_PASSWORD=your_password
   
   REDIS_HOST=localhost
   REDIS_PORT=6379
   REDIS_PASSWORD=
   REDIS_DB=0
   REDIS_PREFIX=i18n_flow:
   
   LOG_LEVEL=info           # debug, info, warn, error, fatal
   LOG_FORMAT=console       # console, json
   LOG_OUTPUT=both          # console, file, both
   LOG_DIR=logs
   LOG_DATE_FORMAT=2006-01-02
   LOG_MAX_SIZE=100         # MB
   LOG_MAX_AGE=7            # days
   LOG_MAX_BACKUPS=5
   LOG_COMPRESS=true
   LOG_ENABLE_CONSOLE=true
   ```

4. Start the server:

   ```bash
   # Using air for hot-reload during development
   air
   
   # Or run directly
   go run cmd/server/main.go
   
   # Or build and run the binary
   go build -o i18n-flow ./cmd/server
   ./i18n-flow
   ```

   The server will start on port 8080 by default. The database tables will be automatically created using GORM AutoMigrate, along with default admin user and language seed data.

5. Access the Swagger documentation:

   ```
   http://localhost:8080/swagger/index.html
   ```

### Initial Access

On first run, the system creates a default admin user based on your environment variables:

- Username: Value from `ADMIN_USERNAME` in .env file (defaults to "admin")
- Password: Value from `ADMIN_PASSWORD` in .env file (defaults to "admin123")

The system also creates 20 default languages (English, Chinese, Japanese, Korean, French, German, Spanish, etc.).

**Important:** It's highly recommended to change the default password after the first login through the admin interface.

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
- `/migrations`: Database migration files (GORM AutoMigrate is used)
- `/utils`: Shared utility functions
- `/tests`: Test files
- `/logs`: Application log files (created at runtime)

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

### Recent Updates

The system has been enhanced with comprehensive project member management functionality:

- **Project Member Handler**: Complete CRUD operations for project members
- **Role-based Access Control**: Three permission levels (viewer, editor, owner)
- **User Management**: Admin interface for creating and managing system users
- **Accessible Projects**: Users can view projects they have access to
- **Permission Checking**: Real-time permission validation for project operations

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
