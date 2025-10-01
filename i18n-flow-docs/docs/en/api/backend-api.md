# Backend API Reference

i18n Flow provides a complete RESTful API to interact with the system. This document details the available API endpoints, parameters, and response formats.

## API Basics

### Base URL

The base URL for all API requests is:

```
http://your-i18n-flow-server.com/api
```

### Authentication

The i18n Flow API supports two authentication methods:

1. **JWT Authentication**: For frontend user access
2. **API Key Authentication**: For CLI tools and other automated integrations

#### JWT Authentication

For requests requiring user authentication, add the JWT token to the HTTP header:

```
Authorization: Bearer <jwt-token>
```

#### API Key Authentication

For CLI tools and automated integrations, add the API key to the HTTP header:

```
X-API-Key: <api-key>
```

The API key is set in the backend `.env` file as `CLI_API_KEY`.

### Response Format

All API responses use JSON format and follow this structure:

**Successful Response**:

```json
{
  "success": true,
  "data": {
    // Response data
  }
}
```

**Error Response**:

```json
{
  "success": false,
  "error": {
    "code": "ERROR_CODE",
    "message": "Error description",
    "details": {} // Optional additional error information
  }
}
```

### Pagination

Endpoints supporting pagination use the following query parameters:

- `page`: Page number, defaults to 1
- `limit`: Items per page, defaults to 20

Paginated responses include the following metadata:

```json
{
  "success": true,
  "data": [...],
  "meta": {
    "total": 100,
    "page": 1,
    "limit": 20,
    "pages": 5
  }
}
```

## Authentication API

### Login

Obtain a JWT token.

**Request**:

```
POST /login
```

**Request Body**:

```json
{
  "username": "admin",
  "password": "password"
}
```

**Response**:

```json
{
  "success": true,
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_at": "2023-06-01T12:00:00Z",
    "user": {
      "id": 1,
      "username": "admin",
      "email": "admin@example.com",
      "role": "admin"
    }
  }
}
```

### Refresh Token

Use a refresh token to obtain a new JWT token.

**Request**:

```
POST /refresh
```

**Request Body**:

```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Response**:

```json
{
  "success": true,
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_at": "2023-06-01T12:00:00Z"
  }
}
```

### Verify CLI API Key

Verify if a CLI API key is valid.

**Request**:

```
GET /cli/auth
```

**Request Headers**:

```
X-API-Key: your-api-key
```

**Response**:

```json
{
  "success": true,
  "data": {
    "valid": true
  }
}
```

## Project API

### Create Project

Create a new project.

**Request**:

```
POST /projects
```

**Request Headers**:

```
Authorization: Bearer <jwt-token>
```

**Request Body**:

```json
{
  "name": "My Project",
  "description": "Project description",
  "slug": "my-project"
}
```

**Response**:

```json
{
  "success": true,
  "data": {
    "id": 1,
    "name": "My Project",
    "description": "Project description",
    "slug": "my-project",
    "created_at": "2023-01-01T00:00:00Z",
    "updated_at": "2023-01-01T00:00:00Z"
  }
}
```

### Get Project List

Get a list of all projects, with pagination.

**Request**:

```
GET /projects?page=1&limit=20
```

**Request Headers**:

```
Authorization: Bearer <jwt-token>
```

**Response**:

```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "name": "My Project",
      "description": "Project description",
      "slug": "my-project",
      "created_at": "2023-01-01T00:00:00Z",
      "updated_at": "2023-01-01T00:00:00Z",
      "stats": {
        "languages": 3,
        "keys": 150,
        "completion": 75
      }
    }
    // More projects...
  ],
  "meta": {
    "total": 10,
    "page": 1,
    "limit": 20,
    "pages": 1
  }
}
```

### Get Project Details

Get detailed information for a single project.

**Request**:

```
GET /projects/detail/:id
```

**Request Headers**:

```
Authorization: Bearer <jwt-token>
```

**Response**:

```json
{
  "success": true,
  "data": {
    "id": 1,
    "name": "My Project",
    "description": "Project description",
    "slug": "my-project",
    "created_at": "2023-01-01T00:00:00Z",
    "updated_at": "2023-01-01T00:00:00Z",
    "stats": {
      "languages": 3,
      "keys": 150,
      "completion": 75,
      "language_stats": [
        {
          "language_code": "en",
          "language_name": "English",
          "total_keys": 150,
          "translated_keys": 150,
          "completion": 100
        },
        {
          "language_code": "zh-CN",
          "language_name": "Simplified Chinese",
          "total_keys": 150,
          "translated_keys": 125,
          "completion": 83
        },
        {
          "language_code": "ja",
          "language_name": "Japanese",
          "total_keys": 150,
          "translated_keys": 75,
          "completion": 50
        }
      ]
    }
  }
}
```

### Update Project

Update an existing project's information.

**Request**:

```
PUT /projects/update/:id
```

**Request Headers**:

```
Authorization: Bearer <jwt-token>
```

**Request Body**:

```json
{
  "name": "My Project (Updated)",
  "description": "Updated project description",
  "slug": "my-project-updated"
}
```

**Response**:

```json
{
  "success": true,
  "data": {
    "id": 1,
    "name": "My Project (Updated)",
    "description": "Updated project description",
    "slug": "my-project-updated",
    "created_at": "2023-01-01T00:00:00Z",
    "updated_at": "2023-01-02T00:00:00Z"
  }
}
```

### Delete Project

Delete a project and all its related data.

**Request**:

```
DELETE /projects/delete/:id
```

**Request Headers**:

```
Authorization: Bearer <jwt-token>
```

**Response**:

```json
{
  "success": true,
  "data": {
    "message": "Project successfully deleted"
  }
}
```

## Project Member Management API

### Get Project Members

Get a list of all members in a project.

**Request**:

```
GET /projects/{project_id}/members
```

**Request Headers**:

```
Authorization: Bearer <jwt-token>
```

**Response**:

```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "user_id": 2,
      "username": "john_doe",
      "email": "john@example.com",
      "role": "editor",
      "joined_at": "2024-01-15T10:30:00Z"
    }
  ]
}
```

### Add Project Member

Add a new member to a project.

**Request**:

```
POST /projects/{project_id}/members
```

**Request Headers**:

```
Authorization: Bearer <jwt-token>
Content-Type: application/json
```

**Request Body**:

```json
{
  "user_id": 2,
  "role": "editor"
}
```

**Response**:

```json
{
  "success": true,
  "data": {
    "id": 1,
    "project_id": 1,
    "user_id": 2,
    "role": "editor",
    "created_at": "2024-01-15T10:30:00Z"
  }
}
```

### Update Member Role

Update a project member's role.

**Request**:

```
PUT /projects/{project_id}/members/{user_id}
```

**Request Headers**:

```
Authorization: Bearer <jwt-token>
Content-Type: application/json
```

**Request Body**:

```json
{
  "role": "owner"
}
```

**Response**:

```json
{
  "success": true,
  "data": {
    "id": 1,
    "project_id": 1,
    "user_id": 2,
    "role": "owner",
    "updated_at": "2024-01-15T10:30:00Z"
  }
}
```

### Remove Project Member

Remove a member from a project.

**Request**:

```
DELETE /projects/{project_id}/members/{user_id}
```

**Request Headers**:

```
Authorization: Bearer <jwt-token>
```

**Response**:

```json
{
  "success": true,
  "data": {
    "message": "Member successfully removed from project"
  }
}
```

### Check User Permission

Check if a user has specific permissions in a project.

**Request**:

```
GET /projects/{project_id}/members/{user_id}/permission?required_role=editor
```

**Request Headers**:

```
Authorization: Bearer <jwt-token>
```

**Response**:

```json
{
  "success": true,
  "data": {
    "has_permission": true
  }
}
```

### Get User's Projects

Get all projects that a specific user has access to.

**Request**:

```
GET /user-projects/{user_id}
```

**Request Headers**:

```
Authorization: Bearer <jwt-token>
```

**Response**:

```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "name": "My Project",
      "slug": "my-project",
      "description": "Project description",
      "role": "owner",
      "member_count": 5,
      "language_count": 3,
      "created_at": "2024-01-01T00:00:00Z"
    }
  ]
}
```

## Language API

### Get Language List

Get a list of all supported languages.

**Request**:

```
GET /languages
```

**Request Headers**:

```
Authorization: Bearer <jwt-token>
```

**Response**:

```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "name": "English",
      "code": "en",
      "locale": "en-US",
      "is_rtl": false,
      "created_at": "2023-01-01T00:00:00Z",
      "updated_at": "2023-01-01T00:00:00Z"
    },
    {
      "id": 2,
      "name": "Simplified Chinese",
      "code": "zh-CN",
      "locale": "zh-Hans",
      "is_rtl": false,
      "created_at": "2023-01-01T00:00:00Z",
      "updated_at": "2023-01-01T00:00:00Z"
    }
    // More languages...
  ]
}
```

### Create Language

Add a new supported language.

**Request**:

```
POST /languages
```

**Request Headers**:

```
Authorization: Bearer <jwt-token>
```

**Request Body**:

```json
{
  "name": "Spanish",
  "code": "es",
  "locale": "es-ES",
  "is_rtl": false
}
```

**Response**:

```json
{
  "success": true,
  "data": {
    "id": 3,
    "name": "Spanish",
    "code": "es",
    "locale": "es-ES",
    "is_rtl": false,
    "created_at": "2023-01-02T00:00:00Z",
    "updated_at": "2023-01-02T00:00:00Z"
  }
}
```

### Update Language

Update an existing language's information.

**Request**:

```
PUT /languages/:id
```

**Request Headers**:

```
Authorization: Bearer <jwt-token>
```

**Request Body**:

```json
{
  "name": "Spanish (Spain)",
  "locale": "es-ES",
  "is_rtl": false
}
```

**Response**:

```json
{
  "success": true,
  "data": {
    "id": 3,
    "name": "Spanish (Spain)",
    "code": "es",
    "locale": "es-ES",
    "is_rtl": false,
    "created_at": "2023-01-02T00:00:00Z",
    "updated_at": "2023-01-03T00:00:00Z"
  }
}
```

### Delete Language

Delete a language and all its related translations.

**Request**:

```
DELETE /languages/:id
```

**Request Headers**:

```
Authorization: Bearer <jwt-token>
```

**Response**:

```json
{
  "success": true,
  "data": {
    "message": "Language successfully deleted"
  }
}
```

## Translation API

### Create Translation

Add a new translation.

**Request**:

```
POST /translations
```

**Request Headers**:

```
Authorization: Bearer <jwt-token>
```

**Request Body**:

```json
{
  "project_id": 1,
  "key": "common.buttons.save",
  "description": "Save button text",
  "translations": [
    {
      "language_id": 1,
      "value": "Save"
    },
    {
      "language_id": 2,
      "value": "保存"
    }
  ]
}
```

**Response**:

```json
{
  "success": true,
  "data": {
    "key_id": 1,
    "key": "common.buttons.save",
    "description": "Save button text",
    "project_id": 1,
    "translations": [
      {
        "id": 1,
        "language_id": 1,
        "language_code": "en",
        "value": "Save"
      },
      {
        "id": 2,
        "language_id": 2,
        "language_code": "zh-CN",
        "value": "保存"
      }
    ],
    "created_at": "2023-01-05T00:00:00Z",
    "updated_at": "2023-01-05T00:00:00Z"
  }
}
```

### Batch Create Translations

Batch add multiple translations.

**Request**:

```
POST /translations/batch
```

**Request Headers**:

```
Authorization: Bearer <jwt-token>
```

**Request Body**:

```json
{
  "project_id": 1,
  "items": [
    {
      "key": "common.buttons.save",
      "description": "Save button text",
      "translations": [
        {
          "language_id": 1,
          "value": "Save"
        },
        {
          "language_id": 2,
          "value": "保存"
        }
      ]
    },
    {
      "key": "common.buttons.cancel",
      "description": "Cancel button text",
      "translations": [
        {
          "language_id": 1,
          "value": "Cancel"
        },
        {
          "language_id": 2,
          "value": "取消"
        }
      ]
    }
  ]
}
```

**Response**:

```json
{
  "success": true,
  "data": {
    "created": 2,
    "updated": 0,
    "failed": 0,
    "items": [
      {
        "key_id": 1,
        "key": "common.buttons.save",
        "success": true
      },
      {
        "key_id": 2,
        "key": "common.buttons.cancel",
        "success": true
      }
    ]
  }
}
```

### Get Translations by Project

Get all translations for a project.

**Request**:

```
GET /translations/by-project/:project_id
```

**Request Headers**:

```
Authorization: Bearer <jwt-token>
```

**Response**:

```json
{
  "success": true,
  "data": [
    {
      "key_id": 1,
      "key": "common.buttons.save",
      "description": "Save button text",
      "project_id": 1,
      "translations": [
        {
          "id": 1,
          "language_id": 1,
          "language_code": "en",
          "value": "Save"
        },
        {
          "id": 2,
          "language_id": 2,
          "language_code": "zh-CN",
          "value": "保存"
        }
      ],
      "created_at": "2023-01-05T00:00:00Z",
      "updated_at": "2023-01-05T00:00:00Z"
    }
    // More translations...
  ]
}
```

### Get Translation Matrix

Get a project's translation matrix (table view).

**Request**:

```
GET /translations/matrix/by-project/:project_id
```

**Request Headers**:

```
Authorization: Bearer <jwt-token>
```

**Response**:

```json
{
  "success": true,
  "data": {
    "languages": [
      {
        "id": 1,
        "code": "en",
        "name": "English"
      },
      {
        "id": 2,
        "code": "zh-CN",
        "name": "Simplified Chinese"
      }
    ],
    "keys": [
      {
        "id": 1,
        "key": "common.buttons.save",
        "description": "Save button text",
        "values": {
          "1": "Save",
          "2": "保存"
        }
      },
      {
        "id": 2,
        "key": "common.buttons.cancel",
        "description": "Cancel button text",
        "values": {
          "1": "Cancel",
          "2": "取消"
        }
      }
    ]
  }
}
```

### Update Translation

Update a translation value.

**Request**:

```
PUT /translations/:id
```

**Request Headers**:

```
Authorization: Bearer <jwt-token>
```

**Request Body**:

```json
{
  "value": "Save Changes"
}
```

**Response**:

```json
{
  "success": true,
  "data": {
    "id": 2,
    "language_id": 2,
    "language_code": "zh-CN",
    "key_id": 1,
    "value": "Save Changes",
    "updated_at": "2023-01-06T00:00:00Z"
  }
}
```

### Delete Translation Key

Delete a translation key and all its related translation values.

**Request**:

```
DELETE /translations/:key_id
```

**Request Headers**:

```
Authorization: Bearer <jwt-token>
```

**Response**:

```json
{
  "success": true,
  "data": {
    "message": "Translation key successfully deleted"
  }
}
```

## CLI API

### Get Translations (CLI)

Get translations for a project, for CLI tool use.

**Request**:

```
GET /cli/translations?project_id=my-project&locales=en,zh-CN
```

**Request Headers**:

```
X-API-Key: your-api-key
```

**Response**:

```json
{
  "success": true,
  "data": {
    "project": {
      "id": 1,
      "name": "My Project",
      "slug": "my-project"
    },
    "languages": [
      {
        "id": 1,
        "code": "en",
        "name": "English"
      },
      {
        "id": 2,
        "code": "zh-CN",
        "name": "Simplified Chinese"
      }
    ],
    "translations": {
      "en": {
        "common.buttons.save": "Save",
        "common.buttons.cancel": "Cancel"
      },
      "zh-CN": {
        "common.buttons.save": "保存",
        "common.buttons.cancel": "取消"
      }
    }
  }
}
```

### Push Translation Keys (CLI)

Push new translation keys from the CLI tool.

**Request**:

```
POST /cli/keys
```

**Request Headers**:

```
X-API-Key: your-api-key
```

**Request Body**:

```json
{
  "project_id": "my-project",
  "keys": [
    "common.buttons.save",
    "common.buttons.cancel",
    "common.buttons.edit"
  ]
}
```

**Response**:

```json
{
  "success": true,
  "data": {
    "added": 1,
    "existing": 2,
    "failed": 0,
    "keys": [
      {
        "key": "common.buttons.save",
        "status": "existing"
      },
      {
        "key": "common.buttons.cancel",
        "status": "existing"
      },
      {
        "key": "common.buttons.edit",
        "status": "added"
      }
    ]
  }
}
```

## Import/Export API

### Export Project Translations

Export translations for a project.

**Request**:

```
GET /exports/project/:project_id?format=json&locales=en,zh-CN
```

**Request Headers**:

```
Authorization: Bearer <jwt-token>
```

**Response**:

File download, format depends on the `format` parameter (json, xlsx, csv).

### Import Project Translations

Import translations for a project.

**Request**:

```
POST /imports/project/:project_id
```

**Request Headers**:

```
Authorization: Bearer <jwt-token>
Content-Type: multipart/form-data
```

**Request Parameters**:

- `file`: File to import
- `format`: File format (json, xlsx, csv)
- `locale`: Language code (only for single language imports)
- `merge_strategy`: Merge strategy (overwrite, keep_existing)

**Response**:

```json
{
  "success": true,
  "data": {
    "imported": 50,
    "updated": 25,
    "skipped": 5,
    "project_id": 1,
    "locale": "zh-CN"
  }
}
```

## Dashboard API

### Get System Statistics

Get system statistics for dashboard display.

**Request**:

```
GET /dashboard/stats
```

**Request Headers**:

```
Authorization: Bearer <jwt-token>
```

**Response**:

```json
{
  "success": true,
  "data": {
    "projects": {
      "total": 5,
      "recent": [
        {
          "id": 1,
          "name": "My Project",
          "slug": "my-project",
          "updated_at": "2023-01-06T00:00:00Z"
        }
        // More projects...
      ]
    },
    "languages": {
      "total": 6,
      "list": [
        {
          "code": "en",
          "name": "English",
          "projects_count": 5
        },
        {
          "code": "zh-CN",
          "name": "Simplified Chinese",
          "projects_count": 4
        }
        // More languages...
      ]
    },
    "translations": {
      "total_keys": 500,
      "total_translations": 2500,
      "completion_rate": 83.5
    },
    "recent_activity": [
      {
        "type": "translation_update",
        "project_id": 1,
        "project_name": "My Project",
        "user": "admin",
        "timestamp": "2023-01-06T12:30:00Z",
        "details": {
          "key": "common.buttons.save",
          "language": "zh-CN"
        }
      }
      // More activities...
    ]
  }
}
```

## Error Codes

Here are common error codes that the system may return:

| Code                            | Description                        |
| ------------------------------- | ---------------------------------- |
| `AUTH_INVALID_CREDENTIALS`      | Invalid authentication credentials |
| `AUTH_TOKEN_EXPIRED`            | Authentication token expired       |
| `AUTH_TOKEN_INVALID`            | Invalid authentication token       |
| `AUTH_INSUFFICIENT_PERMISSIONS` | Insufficient permissions           |
| `RESOURCE_NOT_FOUND`            | Requested resource not found       |
| `RESOURCE_ALREADY_EXISTS`       | Resource already exists            |
| `VALIDATION_ERROR`              | Request data validation error      |
| `SERVER_ERROR`                  | Server internal error              |

## Rate Limiting

The API implements rate limiting to prevent abuse. Limits are as follows:

- Authentication endpoints: 10 requests per minute per IP address
- Other endpoints: 100 requests per minute per API key or user

Exceeding the limit will return HTTP status 429 with the following response:

```json
{
  "success": false,
  "error": {
    "code": "RATE_LIMIT_EXCEEDED",
    "message": "Rate limit exceeded, please try again later",
    "details": {
      "retry_after": 30
    }
  }
}
```

## User Management API

### Create User

Create a new system user (admin only).

**Request**:

```
POST /users
```

**Request Headers**:

```
Authorization: Bearer <jwt-token>
Content-Type: application/json
```

**Request Body**:

```json
{
  "username": "newuser",
  "email": "newuser@example.com",
  "password": "securepassword123",
  "role": "user"
}
```

**Response**:

```json
{
  "success": true,
  "data": {
    "id": 3,
    "username": "newuser",
    "email": "newuser@example.com",
    "role": "user",
    "status": "active",
    "created_at": "2024-01-15T10:30:00Z"
  }
}
```

### Get User List

Get a list of all system users (admin only).

**Request**:

```
GET /users
```

**Request Headers**:

```
Authorization: Bearer <jwt-token>
```

**Response**:

```json
{
  "success": true,
  "data": {
    "users": [
      {
        "id": 1,
        "username": "admin",
        "email": "admin@example.com",
        "role": "admin",
        "status": "active",
        "created_at": "2024-01-01T00:00:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 10,
      "total": 1,
      "total_pages": 1
    }
  }
}
```

### Get User Details

Get detailed information about a specific user.

**Request**:

```
GET /users/{id}
```

**Request Headers**:

```
Authorization: Bearer <jwt-token>
```

**Response**:

```json
{
  "success": true,
  "data": {
    "id": 1,
    "username": "admin",
    "email": "admin@example.com",
    "role": "admin",
    "status": "active",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

### Update User

Update user information (admin only).

**Request**:

```
PUT /users/{id}
```

**Request Headers**:

```
Authorization: Bearer <jwt-token>
Content-Type: application/json
```

**Request Body**:

```json
{
  "username": "updateduser",
  "email": "updated@example.com",
  "role": "admin"
}
```

**Response**:

```json
{
  "success": true,
  "data": {
    "id": 1,
    "username": "updateduser",
    "email": "updated@example.com",
    "role": "admin",
    "updated_at": "2024-01-15T10:30:00Z"
  }
}
```

### Reset User Password

Reset a user's password (admin only).

**Request**:

```
POST /users/{id}/reset-password
```

**Request Headers**:

```
Authorization: Bearer <jwt-token>
```

**Response**:

```json
{
  "success": true,
  "data": {
    "message": "Password reset successful",
    "temporary_password": "temp123456"
  }
}
```

### Delete User

Delete a user from the system (admin only).

**Request**:

```
DELETE /users/{id}
```

**Request Headers**:

```
Authorization: Bearer <jwt-token>
```

**Response**:

```json
{
  "success": true,
  "data": {
    "message": "User successfully deleted"
  }
}
```

### Change Current User Password

Allow users to change their own password.

**Request**:

```
POST /user/change-password
```

**Request Headers**:

```
Authorization: Bearer <jwt-token>
Content-Type: application/json
```

**Request Body**:

```json
{
  "current_password": "oldpassword123",
  "new_password": "newpassword456"
}
```

**Response**:

```json
{
  "success": true,
  "data": {
    "message": "Password changed successfully"
  }
}
```
