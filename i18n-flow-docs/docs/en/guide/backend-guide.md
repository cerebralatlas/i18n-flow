# Backend Usage Guide

The i18n Flow backend is a powerful API server developed in Go that handles all translation data and provides interfaces for the frontend and CLI tool. This guide will help you understand how to configure and use the backend service.

## Configuring the Backend

### Environment Configuration

The backend service uses a `.env` file for configuration. A default configuration template is available in the `.env.example` file, which you need to copy and modify:

```bash
cp .env.example .env
```

Here are the main configuration options:

```properties
# Database Configuration
DB_USERNAME=i18nflow           # Database username
DB_PASSWORD=password           # Database password
DB_HOST=localhost              # Database host
DB_PORT=3306                   # Database port
DB_NAME=i18n_flow              # Database name

# JWT Authentication Configuration
JWT_SECRET=your_secure_jwt_secret                 # JWT secret key
JWT_EXPIRATION_HOURS=24                           # JWT expiration time (hours)
JWT_REFRESH_SECRET=your_secure_refresh_secret     # JWT refresh secret key
JWT_REFRESH_EXPIRATION_HOURS=168                  # JWT refresh expiration time (hours)

# CLI API Configuration
CLI_API_KEY=your_secure_api_key                   # API key for CLI tool

# Server Configuration
SERVER_PORT=8080               # Server listening port
SERVER_HOST=0.0.0.0            # Server listening address
```

### Database Setup

The backend uses MySQL as the database. Make sure you've created the appropriate database:

```sql
CREATE DATABASE i18n_flow CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE USER 'i18nflow'@'localhost' IDENTIFIED BY 'password';
GRANT ALL PRIVILEGES ON i18n_flow.* TO 'i18nflow'@'localhost';
FLUSH PRIVILEGES;
```

::: tip
On first startup, the system will automatically create the required database tables.
:::

## Starting the Service

Starting the backend service is simple:

```bash
# Navigate to the backend directory
cd admin-backend

# Download dependencies
go mod download

# Start the service
go run main.go
```

The service will start on `http://localhost:8080` by default. You can view the API documentation at `http://localhost:8080/swagger/index.html`.

## API Authentication

The i18n Flow backend uses two authentication methods:

1. **JWT Authentication**: For frontend user access
2. **API Key Authentication**: For CLI tool access

### JWT Authentication

To obtain a JWT token:

```http
POST /api/login

{
  "username": "admin",
  "password": "admin123"
}
```

Successful response:

```json
{
  "token": "eyJhbGciOiJI...",
  "refresh_token": "eyJhbGciOiJI...",
  "expires_at": "2023-03-01T12:00:00Z"
}
```

For subsequent requests, add the token to the Authorization header:

```
Authorization: Bearer eyJhbGciOiJI...
```

### API Key Authentication

The CLI tool uses a preset API key for authentication, configured in the `.env` file (`CLI_API_KEY`).

For CLI requests, add the following header:

```
X-API-Key: your_secure_api_key
```

## Main API Endpoints

### Project Management

- `POST /api/projects` - Create a project
- `GET /api/projects` - Get project list
- `GET /api/projects/detail/:id` - Get project details
- `PUT /api/projects/update/:id` - Update a project
- `DELETE /api/projects/delete/:id` - Delete a project

### Language Management

- `GET /api/languages` - Get language list
- `POST /api/languages` - Create a language
- `PUT /api/languages/:id` - Update a language
- `DELETE /api/languages/:id` - Delete a language

### Translation Management

- `POST /api/translations` - Create a translation
- `POST /api/translations/batch` - Batch create translations
- `GET /api/translations/by-project/:project_id` - Get all translations for a project
- `GET /api/translations/matrix/by-project/:project_id` - Get translation matrix for a project
- `GET /api/translations/:id` - Get translation details
- `PUT /api/translations/:id` - Update a translation
- `DELETE /api/translations/:id` - Delete a translation
- `POST /api/translations/batch-delete` - Batch delete translations

### Import/Export

- `GET /api/exports/project/:project_id` - Export project translations
- `POST /api/imports/project/:project_id` - Import project translations

### CLI Tool Integration

- `GET /api/cli/translations` - Get translations (for CLI use)
- `POST /api/cli/keys` - Push translation keys (for CLI use)

### Dashboard

- `GET /api/dashboard/stats` - Get system statistics

## Data Models

The i18n Flow backend uses the following main data models:

### Project

```go
type Project struct {
    ID          uint      `json:"id" gorm:"primaryKey"`
    Name        string    `json:"name" gorm:"not null"`
    Description string    `json:"description"`
    Slug        string    `json:"slug" gorm:"uniqueIndex;not null"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}
```

### Language

```go
type Language struct {
    ID        uint      `json:"id" gorm:"primaryKey"`
    Name      string    `json:"name" gorm:"not null"`
    Code      string    `json:"code" gorm:"uniqueIndex;not null"`
    Locale    string    `json:"locale" gorm:"uniqueIndex"`
    IsRTL     bool      `json:"is_rtl" gorm:"default:false"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

### Translation Key

```go
type Key struct {
    ID          uint      `json:"id" gorm:"primaryKey"`
    ProjectID   uint      `json:"project_id" gorm:"not null"`
    Name        string    `json:"name" gorm:"not null"`
    Description string    `json:"description"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}
```

### Translation Value

```go
type Translation struct {
    ID         uint      `json:"id" gorm:"primaryKey"`
    KeyID      uint      `json:"key_id" gorm:"not null"`
    LanguageID uint      `json:"language_id" gorm:"not null"`
    Value      string    `json:"value"`
    CreatedAt  time.Time `json:"created_at"`
    UpdatedAt  time.Time `json:"updated_at"`
}
```

## Logging and Monitoring

The i18n Flow backend uses structured logging to record system activities. Logs are output to standard out by default, but you can configure them to be redirected to a file:

```bash
# Output logs to a file
go run main.go > i18n-flow.log 2>&1
```

## Backup and Recovery

Regular database backups are a good practice:

```bash
# Backup the database
mysqldump -u i18nflow -p i18n_flow > i18n_flow_backup.sql

# Restore the database
mysql -u i18nflow -p i18n_flow < i18n_flow_backup.sql
```

## Performance Optimization

For large projects, consider the following optimizations:

1. Enable database connection pooling
2. Configure appropriate caching strategies
3. Use a reverse proxy (like Nginx) for handling static resources

These configurations can be adjusted in the `.env` file for production environments.

## Troubleshooting

### Database Connection Issues

If you encounter database connection problems:

1. Check if the database credentials are correct
2. Confirm that the database service is running
3. Check network connections and firewall settings

### API Error Responses

System error responses have the following format:

```json
{
  "error": "Error message",
  "code": "ERROR_CODE",
  "details": {} // Optional error details
}
```
