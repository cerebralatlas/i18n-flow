# Backend API Reference

## Authentication

### Login

```http
POST /api/auth/login
```

### Refresh Token

```http
POST /api/auth/refresh
```

## Projects

### Create Project

```http
POST /api/projects
```

### List Projects

```http
GET /api/projects
```

## Translations

### Create Translation

```http
POST /api/translations
```

### Batch Create Translations

```http
POST /api/translations/batch
```

### Get Translation Matrix

```http
GET /api/translations/matrix/by-project/{project_id}
```
