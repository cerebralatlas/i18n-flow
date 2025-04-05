# i18n-flow Backend

A powerful backend service for internationalization management in application development.

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

- **Language**: Go (Golang)
- **Web Framework**: Gin
- **ORM**: GORM
- **Database**: MySQL/MariaDB
- **Documentation**: Swagger/OpenAPI
- **Authentication**: JWT & API keys

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

- Go 1.16+
- MySQL/MariaDB
- Git

### Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/yourusername/i18n-flow.git
   cd i18n-flow/admin-backend
   ```

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
   ```

4. Start the server:

   ```bash
   go run main.go
   ```

   The server will start on port 8080 by default. The database tables will be automatically created.

5. Access the Swagger documentation:

   ```
   http://localhost:8080/swagger/index.html
   ```

### Initial Access

On first run, the system creates a default admin user:

- Username: `admin`
- Password: `admin123`

It's highly recommended to change this password after the first login.

## Development

### API Documentation

The API documentation is automatically generated with Swagger. You can access it at:

```
http://localhost:8080/swagger/index.html
```

### Project Structure

- `/controller`: HTTP request handlers
- `/service`: Business logic layer
- `/model`: Data models and database operations
- `/middleware`: Request processing middleware
- `/config`: Application configuration
- `/docs`: Auto-generated API documentation

### Running Tests

```bash
go test ./...
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

# Run container
docker run -p 8080:8080 --env-file .env i18n-flow-backend
```

## License

This project is licensed under the MIT License.
