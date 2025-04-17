# Backend Guide

## Architecture

The backend is built with Go and uses the following technologies:

- Gin Framework for REST API
- GORM for database operations
- JWT for authentication
- Swagger for API documentation

## Directory Structure

```bash
admin-backend/
├── config/       # Configuration files
├── controller/   # HTTP request handlers
├── middleware/   # Custom middleware
├── model/       # Database models
├── service/     # Business logic
└── docs/        # Swagger documentation
```

## Environment Variables

```bash
# Database Configuration
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=password
DB_NAME=i18n_flow

# JWT Configuration
JWT_SECRET=your-secret-key
JWT_EXPIRE_HOURS=24
```

## API Endpoints

The backend provides the following main API groups:

- `/api/auth` - Authentication endpoints
- `/api/projects` - Project management
- `/api/translations` - Translation management
- `/api/languages` - Language management
